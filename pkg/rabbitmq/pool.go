package mq

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go-program/pkg/loader"
	"go-program/pkg/logger"
	"go-program/pkg/trace"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type RabbitMQPool struct {
	pool         chan *ChannelPool
	url          string
	maxConn      int
	maxChannel   int
	connTimeout  time.Duration
	reconnectDur time.Duration
	closed       bool
	mutex        sync.RWMutex
	declareOnce  sync.Once
}

func (p *RabbitMQPool) GetName() string {
	return "RabbitMQPool"
}

func (p *RabbitMQPool) QueueDeclare(ch *amqp.Channel, queueName string) error {
	var err error
	p.declareOnce.Do(
		func() {
			_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
			if err != nil {
				return
			}
		})
	return err
}

func (p *RabbitMQPool) RunLoad() error {
	var (
		err error
	)

	maxChannel := viper.GetInt("RMQ_MAX_CHANNEL")
	if maxChannel == 0 {
		maxChannel = 5 // 默认开启5个channel
	}

	p.pool = make(chan *ChannelPool, maxChannel)

	dsn := viper.GetString("RMQ_DSN")
	if len(dsn) <= 0 {
		return errors.New("miss env RMQ_DSN")
	}

	p.url = dsn

	maxConn := viper.GetInt("RMQ_MAX_CONN")
	if maxConn == 0 {
		maxConn = 1 // 默认1个连接
	}
	p.maxConn = maxConn
	p.maxChannel = maxChannel

	connTimeout := viper.GetInt64("RMQ_CONN_TIMEOUT")
	if connTimeout == 0 {
		connTimeout = 10
	}
	p.connTimeout = time.Duration(connTimeout) * time.Second

	reconnectDur := viper.GetInt64("RMQ_RECONNECT_DURATION")
	if reconnectDur == 0 {
		reconnectDur = 15
	}
	p.reconnectDur = time.Duration(reconnectDur) * time.Second

	// Initialize the pool with connections
	for i := 0; i < p.maxConn; i++ {
		var conn *ChannelPool
		conn, err = p.createConnection()
		if err != nil {
			return err
		}

		p.pool <- conn
	}

	return nil
}

var ins = &RabbitMQPool{}

func init() {
	loader.Register(ins)
}

func GetRabbitMqPool() *RabbitMQPool {
	return ins
}

func NewRabbitMQPool(url string, maxConn, maxChannel int, connTimeout, reconnectDur time.Duration) (pool *RabbitMQPool, err error) {
	pool = &RabbitMQPool{
		pool:         make(chan *ChannelPool, maxConn),
		url:          url,
		maxConn:      maxConn,
		maxChannel:   maxChannel,
		connTimeout:  connTimeout,
		reconnectDur: reconnectDur,
	}

	// Initialize the pool with connections
	for i := 0; i < maxConn; i++ {
		var conn *ChannelPool
		conn, err = pool.createConnection()
		if err != nil {
			return
		}

		pool.pool <- conn
	}

	return pool, nil
}

func (p *RabbitMQPool) createConnection() (*ChannelPool, error) {
	conn, err := amqp.Dial(p.url)
	if err != nil {
		return nil, err
	}

	chPool, err := NewChannelPool(conn, p.maxChannel, p.reconnectDur)
	if err != nil {
		return nil, err
	}
	return chPool, nil
}

func (p *RabbitMQPool) isConnectionOpen(conn *amqp.Connection) bool {
	return conn != nil && !conn.IsClosed()
}

var ErrCloseRabbitMQPool = errors.New("rabbit mq pool is closed")

func (p *RabbitMQPool) Get() (*ChannelPool, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	if p.closed {
		return nil, ErrCloseRabbitMQPool
	}
	select {
	case conn := <-p.pool:
		if p.isConnectionOpen(conn.conn) {
			return conn, nil
		}
		newConn, err := p.createConnection()
		if err != nil {
			p.Put(conn)
			return nil, err
		}
		return newConn, nil
	case <-time.After(p.connTimeout):
		return nil, fmt.Errorf("getting connection timeout")
	}
}

func (p *RabbitMQPool) Put(conn *ChannelPool) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	if p.closed {
		return
	}

	p.pool <- conn
}

func (p *RabbitMQPool) Close() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.closed {
		return nil
	}
	close(p.pool)
	p.closed = true

	for conn := range p.pool {
		if err := conn.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (p *RabbitMQPool) ConsumeMessages(queueName string, handle func(ctx context.Context, delivery amqp.Delivery)) error {
	conn, err := p.Get()
	if err != nil {
		time.Sleep(p.reconnectDur)
		return err
	}
	defer p.Put(conn)

	ch, err := conn.GetChannel()
	if err != nil {
		return err
	}
	defer conn.Put(ch)

	err = p.QueueDeclare(ch, queueName)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		time.Sleep(conn.reconnectDur)
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var traceId string
			tmp, ok := d.Headers["traceId"]
			if ok {
				traceId = tmp.(string)
			} else {
				traceId = uuid.New().String()
			}
			ctx := trace.NewTraceCtxWithTraceID(context.Background(), traceId)
			logger.GetLoggerWithContext(ctx).WithFields(logrus.Fields{
				"body":   string(d.Body),
				"module": "consumer",
			})
			handle(ctx, d)
		}
		log.Println("Consumer closed - reconnecting...")
		forever <- true
	}()
	<-forever
	return nil
}

func (p *RabbitMQPool) PublishMessages(ctx context.Context, queueName string, msg []byte) error {
	conn, err := p.Get()
	if err != nil {
		time.Sleep(p.reconnectDur)
		return err
	}
	defer p.Put(conn)

	ch, err := conn.GetChannel()
	if err != nil {
		return err
	}
	defer conn.Put(ch)

	err = p.QueueDeclare(ch, queueName)
	if err != nil {
		return err
	}

	traceId, err := trace.GetTraceID(ctx)
	if err != nil {
		return err
	}

	err = ch.Publish("", queueName, false, false, amqp.Publishing{
		Headers: map[string]interface{}{
			"traceId": traceId,
		},
		ContentType:  "application/json",
		DeliveryMode: 2,
		Priority:     9,
		Expiration:   "",
		Body:         msg,
	},
	)
	logger.GetLoggerWithContext(ctx).WithFields(logrus.Fields{
		"body":   string(msg),
		"module": "producer",
	})
	if err != nil {
		time.Sleep(conn.reconnectDur)
		return err
	}

	return nil
}

// ChannelPool 代表一个 RabbitMQ 通道池
type ChannelPool struct {
	channels     chan *amqp.Channel
	conn         *amqp.Connection
	reconnectDur time.Duration
	closed       bool
	mutex        sync.RWMutex
}

// NewChannelPool 创建一个新的 ChannelPool
func NewChannelPool(conn *amqp.Connection, maxChannels int, reconnectDur time.Duration) (*ChannelPool, error) {
	pool := &ChannelPool{
		channels:     make(chan *amqp.Channel, maxChannels),
		conn:         conn,
		reconnectDur: reconnectDur,
	}

	if err := pool.initializeChannels(); err != nil {
		return nil, err
	}

	return pool, nil
}

// initializeChannels 初始化通道池
func (p *ChannelPool) initializeChannels() error {
	for i := 0; i < cap(p.channels); i++ {
		ch, err := p.openChannel()
		if err != nil {
			return err
		}
		p.channels <- ch
	}
	return nil
}

// initializeChannels 初始化通道池
func (p *ChannelPool) openChannel() (*amqp.Channel, error) {
	ch, err := p.conn.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil

}

func (p *ChannelPool) isChannelValid(ch *amqp.Channel) bool {
	if ch == nil {
		return false
	}
	err := ch.ExchangeDeclarePassive(
		"amq.direct", // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	return err == nil
}

var ErrClosedChannelPool = errors.New("rabbit mq pool is closed")

// GetChannel 从池中获取一个通道
func (p *ChannelPool) GetChannel() (*amqp.Channel, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	if p.closed {
		return nil, ErrClosedChannelPool
	}
	ch := <-p.channels
	if p.isChannelValid(ch) {
		return ch, nil
	}
	newCh, err := p.openChannel()
	if err != nil {
		p.Put(ch)
		return nil, err
	}
	return newCh, nil
}

// Put 将通道放回池中
func (p *ChannelPool) Put(ch *amqp.Channel) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	if p.closed {
		return
	}
	p.channels <- ch
}

// Close 关闭连接池
func (p *ChannelPool) Close() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.closed = true
	close(p.channels)
	for ch := range p.channels {
		if err := ch.Close(); err != nil {
			return err
		}
	}
	if err := p.conn.Close(); err != nil {
		return err
	}
	return nil
}
