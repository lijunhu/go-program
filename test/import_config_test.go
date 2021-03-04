package test

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type Rewrite struct {
	Cond  string `bson:"cond" json:"cond"`
	Type  string `bson:"type" json:"type"`   //枚举last，break，redirect，permanent
	Start string `bson:"start" json:"start"` //rewrite前的匹配正则
	End   string `bson:"end" json:"end"`     //rewrite后的重写地址
}

type Rewrites []Rewrite

func (items Rewrites) Len() int { return len(items) }

func (items Rewrites) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

type SortRewritesByType struct {
	Rewrites
}

func (s SortRewritesByType) Less(i, j int) bool {
	return s.Rewrites[i].Type > s.Rewrites[j].Type
}

type PageCache struct {
	Using          bool     `bson:"using" json:"using"`                   //是否启用
	LocalMaxAge    int64    `bson:"localMaxAge" json:"localMaxAge"`       //本地缓存超时时间
	ContentType    string   `bson:"contentType" json:"contentType"`       //请求contentType
	ParamKeyList   []string `bson:"paramKeyList" json:"paramKeyList"`     //页面关键参数key
	CookieKeyList  []string `bson:"cookieKeyList" json:"cookieKeyList"`   //请求Cookie关键key
	NeedContentStr string   `bson:"needContentStr" json:"needContentStr"` //页面必须包含缓存内容
	StaticBackup   string   `bson:"staticBackup" json:"staticBackup"`     //静态兜底信息
	DiffDeveice    bool     `bson:"diffDeveice" json:"diffDeveice"`       //是否区分设备 android、ios、pc
	DiffAddr       bool     `bson:"diffAddr" json:"diffAddr"`             //是否区分地理位置
	ParseTemplate  bool     `bson:"parseTemplate" json:"parseTemplate"`   //是否开启模板解析
	CacheAjax      bool     `bson:"cacheAjax" json:"cacheAjax"`           //是否缓存ajax请求
	DynamicBackup  bool     `bson:"dynamicBackup" json:"dynamicBackup"`   //是否开启动态兜底
	IsSpider       bool     `bson:"isSpider" json:"isSpider"`             //是否蜘蛛单独缓存
}

type IpLocation struct {
	Using bool `bson:"using" json:"using"`
}

type BaseConfig struct {
	Using               bool              `bson:"using" json:"using" `
	TemplateId          bson.ObjectId     `bson:"templateId,omitempty" json:"templateId"`
	TemplateName        string            `bson:"templateName" json:"templateName"`
	MatchType           string            `bson:"matchType" json:"matchType"` //uri匹配类型    equal-全匹配 prefix-前缀匹配 regex-正则匹配
	Uri                 string            `bson:"uri" json:"uri"`
	Hosts               []string          `bson:"hosts" json:"hosts"`
	GroupId             bson.ObjectId     `bson:"groupId,omitempty" json:"groupId"`
	ProxyPassType       string            `bson:"proxyPassType" json:"proxyPassType"` //server-指定服务器 、name-指定名称
	ServerId            bson.ObjectId     `bson:"serverId,omitempty" json:"serverId"`
	IsNameServer        bool              `bson:"isNameServer" json:"isNameServer"`     //是否名字服务
	DefaultNameEnv      string            `bson:"defaultNameEnv" json:"defaultNameEnv"` //默认名字服务器环境
	UpStreamName        string            `bson:"upstreamName" json:"upstreamName"`
	Upstreams           []FeatureUpstream `bson:"upstreams" json:"upstreams"`
	NextUpstreamRules   []string          `bson:"nextUpstreamRules" json:"nextUpstreamRules"`
	NextUpstreamTimeout int64             `bson:"nextUpstreamTimeout" json:"nextUpstreamTimeout"`
	NextUpstreamTries   int64             `bson:"nextUpstreamTries" json:"nextUpstreamTries"`
	ConnectTimeOut      int64             `bson:"connectTimeOut" json:"connectTimeOut"` //连接超时时间 单位ms
	ReadTimeOut         int64             `bson:"readTimeOut" json:"readTimeOut"`       //读取超时时间 单位ms
	SendTimeOut         int64             `bson:"sendTimeOut" json:"sendTimeOut"`       //发送超时时间 单位ms
	CustResHeader       map[string]string `bson:"custResHeader" json:"custResHeader"`
	Tags                []string          `bson:"tags" json:"tags"`
	Rewrites            []Rewrite         `bson:"rewrites" json:"rewrites"`
	SslEnable           bool              `bson:"sslEnable" json:"sslEnable"`                   //是否支持Https
	Http2Enable         bool              `bson:"http2Enable" json:"http2Enable"`               //是否支持Http2
	CertificateId       bson.ObjectId     `bson:"certificateId,omitempty" json:"certificateId"` //证书Id
	CertificateName     string            `bson:"certificateName"json:"certificateName"`        //证书名称
	TargetMethod        string            `bson:"targetMethod" json:"targetMethod"`             //目标方法
}

type FeatureUpstream struct {
	ProxyPassType  string        `bson:"proxyPassType" json:"proxyPassType"` //server-指定服务器 、name-指定名称
	ServerId       bson.ObjectId `bson:"serverId,omitempty" json:"serverId"`
	Rules          []Rules       `bson:"rules" json:"rules"`                   //特征列表
	IsNameServer   bool          `bson:"isNameServer" json:"isNameServer"`     //是否名字服务
	DefaultNameEnv string        `bson:"defaultNameEnv" json:"defaultNameEnv"` //默认名字服务器环境
	UpStreamName   string        `bson:"upstreamName" json:"upstreamName"`
}

type JanusConfig struct {
	Id                     bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	LocationId             string                 `bson:"locationId" json:"locationId"` //该字段只在生成nginx配置时使用
	BaseConfig             BaseConfig             `bson:"baseConfig" json:"baseConfig"`
	IpLocation             IpLocation             `bson:"ipLocation" json:"ipLocation"`
	PageCache              PageCache              `bson:"pageCache" json:"pageCache"`
	Canary                 Canary                 `bson:"canary" json:"canary"`
	ValidateRequest        ValidateRequest        `bson:"validateRequest" json:"validateRequest"`
	TransProtocol          TransProtocol          `bson:"transProtocol" json:"transProtocol"`
	Compress               Compress               `bson:"compress" json:"compress"`
	Decompress             Decompress             `bson:"decompress" json:"decompress"`
	Crypt                  Crypt                  `bson:"crypt" json:"crypt"`
	Decrypt                Decrypt                `bson:"decrypt" json:"decrypt"`
	AntiSpider             AntiSpider             `bson:"antiSpider" json:"antiSpider"`
	Intercept              Intercept              `bson:"intercept" json:"intercept"`                           //拦截策略
	CheckList              CheckList              `bson:"checkList" json:"checkList"`                           //checkList日志
	Exception              Exception              `bson:"exception" json:"exception"`                           //exception
	RequestValidateSandbox RequestValidateSandbox `bson:"requestValidateSandbox" json:"requestValidateSandbox"` //请求校验沙盒
	NeoTransform           NeoTransform           `bson:"neoTransform" json:"neoTransform"`
	CreateUserId           string                 `bson:"createUserId,omitempty" json:"createUserId"`
	CreateUserName         string                 `bson:"createUserName,omitempty" json:"createUserName"`
	CreateTime             time.Time              `bson:"createTime,omitempty" json:"createTime"`        //创建时间
	UpdateTime             time.Time              `bson:"updateTime,omitempty" json:"updateTime"`        //更新时间
	UpdateUserId           string                 `bson:"updateUserId,omitempty"json:"updateUserId"`     //更新人工号
	UpdateUserName         string                 `bson:"updateUserName,omitempty"json:"updateUserName"` //更新人姓名
}

//
type CheckList struct {
	Using          bool   `bson:"using" json:"using"`                //是否启用
	ProductLine    string `bson:"productLine" json:"productLine"`    //产品线
	SamplingRate   int    `bson:"samplingRate"  json:"samplingRate"` //采样率
	RecordHttpBody bool   `bson:"recordHttpBody" json:"recordHttpBody"`
}

type Exception struct {
	Using bool   `bson:"using" json:"using"` //是否启用
	Model string `bson:"model" json:"model"` //产品线
}

type Intercept struct {
	Using              bool                `bson:"using" json:"using"`                           //是否启用
	HashKey            HashKey             `bson:"hashKey" json:"hashKey"`                       //hashKey
	InterceptStrategys []InterceptStrategy `bson:"interceptStrategys" json:"interceptStrategys"` //拦截策略
}
type InterceptStrategy struct {
	InterceptType        string `bson:"interceptType" json:"interceptType"`               //类型    feature-流量特征  ratio-流量比例
	Quota                int64  `bson:"quota" json:"quota"`                               //流量比例
	Type                 string `bson:"type" json:"type"`                                 //类型 0-放行 1-拦截
	Rules                []Rule `bson:"rules" json:"rules"`                               //拦截特征
	StartTime            int64  `bson:"startTime" json:"startTime"`                       //拦截开始时间
	EndTime              int64  `bson:"endTime" json:"endTime"`                           //拦截结束时间
	InterceptRespCode    int64  `bson:"interceptRespCode" json:"interceptRespCode"`       //拦截响应状态吗
	InterceptRespContent string `bson:"interceptRespContent" json:"interceptRespContent"` //拦截响应内容
	ContentType          string `bson:"contentType"  json:"contentType"`                  //contentType
}

type InterceptStrategys []InterceptStrategy

func (items InterceptStrategys) Len() int { return len(items) }

func (items InterceptStrategys) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

//反扒配置
type AntiSpider struct {
	Using           bool            `bson:"using" json:"using"`
	FlowFeature     string          `bson:"flowFeature" json:"flowFeature"`  //流量特征提取模板
	AppUniqueKey    string          `bson:"appUniqueKey"json:"appUniqueKey"` //反爬
	SpiderIntercept SpiderIntercept `bson:"spiderIntercept" json:"spiderIntercept"`
	ImportBlackHole ImportBlackHole `bson:"importBlackHole" json:"importBlackHole"`
	Poison          Poison          `bson:"poison" json:"poison"`
	SyncSpider      bool            `bson:"syncSpider" json:"syncSpider"`
}

//拦截动作
type SpiderIntercept struct {
	Using       bool   `bson:"using" json:"using"`
	Resp        string `bson:"resp" json:"resp"`
	ContentType string `bson:"contentType" json:"contentType"`
}

//导入黑洞
type ImportBlackHole struct {
	Using             bool               `bson:"using" json:"using"`
	AntiSpiderActions []AntiSpiderAction `bson:"antiSpiderActions" json:"antiSpiderActions"`
}

//下毒
type Poison struct {
	Using             bool               `bson:"using" json:"using"`
	AntiSpiderActions []AntiSpiderAction `bson:"antiSpiderActions" json:"antiSpiderActions"`
}

type AntiSpiderAction struct {
	ValidateResults   []string  `bson:"validateResults" json:"validateResults"` // -2-黑名单 1-白名单 0-正常 1判定为爬虫 2-疑似爬虫
	ProxyUpstreamName string    `bson:"proxyUpstreamName" json:"proxyUpstreamName"`
	NameServerEnable  bool      `bson:"nameServerEnable" json:"nameServerEnable"` //是否名字服务
	NameServerEnv     string    `bson:"nameServerEnv" json:"nameServerEnv"`       //名字服务器环境
	ProxyUpstream     string    `bson:"proxyUpstream" json:"proxyUpstream"`       // 跳转Upstream
	Rewrites          []Rewrite `bson:"rewrites" json:"rewrites"`                 //rewrite地址
}

type HashKey struct {
	KeyName string `bson:"keyName" json:"keyName"`
	KeyType string `bson:"keyType" json:"keyType"` // cookie|header|param
}
type Config struct {
	Type         string `bson:"type" json:"type"`                 //泳道策略类型     feature-流量特征  ratio-流量比例
	Quota        int    `bson:"quota" json:"quota"`               //流量比例
	Rule         []Rule `bson:"rule" json:"rule"`                 //流量特征
	CanaryName   string `bson:"canaryName" json:"canaryName"`     //泳道名称
	CanaryHeader string `bson:"canaryHeader" json:"canaryHeader"` //泳道header
}

type Configs []Config

func (configs Configs) Len() int { return len(configs) }

func (configs Configs) Swap(i, j int) {
	configs[i], configs[j] = configs[j], configs[i]
}

type SortConfigsByName struct {
	Configs
}

func (s SortConfigsByName) Less(i, j int) bool {
	return s.Configs[i].CanaryName > s.Configs[j].CanaryName
}

type Rule struct {
	Type    string `bson:"type" json:"type"`
	KeyName string `bson:"keyName" json:"keyName"`
	KeyType string `bson:"keyType" json:"keyType"` // cookie|header|param|ip|city
	Value   string `bson:"value" json:"value"`
}

type Rules []Rule

func (items Rules) Len() int { return len(items) }

func (items Rules) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

type SortRulesByType struct {
	Rules
}

func (s SortRulesByType) Less(i, j int) bool {
	return s.Rules[i].Type > s.Rules[j].Type
}

type Canary struct {
	Using   bool     `bson:"using" json:"using"`     //泳道是否启用
	HashKey HashKey  `bson:"hashKey" json:"hashKey"` // hash标识
	Configs []Config `bson:"configs" json:"configs"` //泳道配置信息
}

type ValidateRequest struct {
	Using               bool   `bson:"using" json:"using"`                             //请求合法性校验
	HeaderValidateType  string `bson:"headerValidateType" json:"headerValidateType"`   //请求头校验类型  elong、tongcheng
	SessionValidateType string `bson:"sessionValidateType" json:"sessionValidateType"` //session校验类型 ForceCheckStrategy,ForceUnCheckStrategy,ForceCheckWithVersionStrategy,ForceCheckWriteCardNoStrategy,ValidateCardNoWithSession
}

type TransProtocol struct {
	Using          bool   `bson:"using" json:"using"`                   //协议转换是否启用
	TargetProtocol string `bson:"targetProtocol" json:"targetProtocol"` //目标协议
}

type Compress struct {
	Using        bool   `bson:"using" json:"using"`               //是否启用
	CompressType string `bson:"compressType" json:"compressType"` //类型
}
type Decompress struct {
	Using          bool   `bson:"using" json:"using"`
	DecompressType string `bson:"decompressType" json:"decompressType"`
}
type Crypt struct {
	Using     bool   `bson:"using" json:"using"`         //是否启用
	CryptType string `bson:"cryptType" json:"cryptType"` //加解密方法
}
type Decrypt struct {
	Using            bool     `bson:"using" json:"using"`
	DecryptType      string   `bson:"decryptType" json:"decryptType"`
	DecryptWhiteList []string `bson:"decryptWhiteList" json:"decryptWhiteList"` //白名单
}

//请求校验沙盒
type RequestValidateSandbox struct {
	Using        bool   `bson:"using" json:"using"`
	ValidateType string `bson:"validateType" json:"validateType"`
}

//中台接入
type NeoTransform struct {
	Using        bool   `bson:"using" json:"using"`
	HttpMethod   string `bson:"httpMethod" json:"httpMethod"`
	ServiceGroup string `bson:"serviceGroup" json:"serviceGroup"`
	Service      string `bson:"service" json:"service"`
	Api          string `bson:"api" json:"api"`
}

type JanusVersion struct {
	Id        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	VersionId string        `json:"versionId" bson:"versionId,omitempty"`
	Name      string        `json:"name" bson:"name"`
	GroupId   bson.ObjectId `json:"groupId" bson:"groupId"`
	//ServerConfig   []JanusAppServer `json:"serverConfig" bson:"serverConfig"`
	ConfigCenter   []JanusConfig `json:"configCenter" bson:"configCenter"`
	LogServer      string        `json:"logServer" bson:"logServer"`                     //根据idc不同logServer也不相同
	CreateTime     time.Time     `json:"createTime" bson:"createTime,omitempty"`         //创建时间
	CreateUserId   string        `json:"createUserId" bson:"createUserId,omitempty"`     //创建人工号
	CreateUserName string        `json:"createUserName" bson:"createUserName,omitempty"` //创建人姓名
	Md5sum         string        `json:"md5sum" bson:"md5sum"`
}

type response struct {
	Code   int          `json:"code"`   //返回code 0 表示成功
	Msg    string       `json:"msg"`    //描述信息
	Result JanusVersion `json:"result"` //结果数据
}

func Test_Import_Config(t *testing.T) {
	req := httplib.Get("http://10.100.202.119:6010/janus-api/api/version/5cfb7443b1e3c1519246ce42?id=5d8cd14fb1e3c1797f297ed3&queryType=0").
		Header("user-token", "5d775164b1e3c13b35ffbb66")
	resp, err := req.Response()
	if err != nil {
		fmt.Println("查询数据错误")
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("http请求错误")
	}
	body := resp.Body
	bodyData, _ := ioutil.ReadAll(body)

	var response response
	_ = json.Unmarshal(bodyData, &response)

	janusVersion := response.Result


	for i, _ := range janusVersion.ConfigCenter {
		//janusVersion.ConfigCenter[i].Id = ""

		configData,_:= json.Marshal(janusVersion.ConfigCenter[i])
		fmt.Println(string(configData))
		/*configReq := httplib.Post("http://10.100.202.119:6010/janus-api/api/uriconfig/5cfb7443b1e3c1519246ce42/config").
			Header("user-token", "5d775164b1e3c13b35ffbb66")
		request, _ := configReq.JSONBody(janusVersion.ConfigCenter[i])
		resp, err = request.DoRequest()
		fmt.Println(resp.StatusCode)*/
	}

}
