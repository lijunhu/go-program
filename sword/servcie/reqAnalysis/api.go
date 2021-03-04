package reqAnalysis

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/httplib"
	"io/ioutil"
	"time"

	"gopkg.in/mgo.v2/bson"
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
	Using          bool     `bson:"using" json:"using" web:"启用状态"`                             //是否启用
	LocalMaxAge    int64    `bson:"localMaxAge" json:"localMaxAge" web:"本地缓存时间"`               //本地缓存超时时间
	ContentType    string   `bson:"contentType" json:"contentType" web:"请求 body 类型"`           //请求contentType
	ParamKeyList   []string `bson:"paramKeyList" json:"paramKeyList" web:"页面关键参数"`             //页面关键参数key
	CookieKeyList  []string `bson:"cookieKeyList" json:"cookieKeyList" web:"请求 cookie 关键 key"` //请求Cookie关键key
	NeedContentStr string   `bson:"needContentStr" json:"needContentStr" web:"页面必然存在内容"`       //页面必须包含缓存内容
	StaticBackup   string   `bson:"staticBackup" json:"staticBackup" web:"静态兜底信息"`             //静态兜底信息
	DiffDeveice    bool     `bson:"diffDeveice" json:"diffDeveice" web:"区分设备"`                 //是否区分设备 android、ios、pc
	DiffAddr       bool     `bson:"diffAddr" json:"diffAddr" web:"区分地理位置"`                     //是否区分地理位置
	ParseTemplate  bool     `bson:"parseTemplate" json:"parseTemplate" web:"模板解析"`             //是否开启模板解析
	CacheAjax      bool     `bson:"cacheAjax" json:"cacheAjax" web:"Ajax 缓存"`                  //是否缓存ajax请求
	DynamicBackup  bool     `bson:"dynamicBackup" json:"dynamicBackup" web:"动态兜底"`             //是否开启动态兜底
	IsSpider       bool     `bson:"isSpider" json:"isSpider" web:"单独缓存爬虫"`                     //是否蜘蛛单独缓存
}

type IpLocation struct {
	Using bool `bson:"using" json:"using"`
}

type BaseConfig struct {
	Using                 bool              `bson:"using" json:"using"  web:"启用状态"`                                  //
	TemplateId            bson.ObjectId     `bson:"templateId,omitempty" json:"templateId" web:"模板 Id"`              //
	TemplateName          string            `bson:"templateName" json:"templateName" web:"模板名字"`                     //
	MatchType             string            `bson:"matchType" json:"matchType" web:"uri 匹配类型"`                       //uri匹配类型    equal-全匹配 prefix-前缀匹配 regex-正则匹配
	Uri                   string            `bson:"uri" json:"uri" web:"uri"`                                        //
	Hosts                 []string          `bson:"hosts" json:"hosts" web:"hosts"`                                  //
	GroupId               bson.ObjectId     `bson:"groupId,omitempty" json:"groupId"`                                //
	ConnectTimeOut        int64             `bson:"connectTimeOut" json:"connectTimeOut" web:"连接超时时间"`               //连接超时时间 单位ms
	ReadTimeOut           int64             `bson:"readTimeOut" json:"readTimeOut" web:"读取超时时间"`                     //读取超时时间 单位ms
	SendTimeOut           int64             `bson:"sendTimeOut" json:"sendTimeOut" web:"发送超时时间"`                     //发送超时时间 单位ms
	SslEnable             bool              `bson:"sslEnable" json:"sslEnable" web:"支持 Https"`                       //是否支持Https
	Http2Enable           bool              `bson:"http2Enable" json:"http2Enable" web:"支持 Http2"`                   //是否支持Http2
	CertificateId         bson.ObjectId     `bson:"certificateId,omitempty" json:"certificateId" web:"证书 Id"`        //证书Id
	CertificateName       string            `bson:"certificateName"json:"certificateName" web:"证书名称"`                //证书名称
	CustomResponseHeaders map[string]string `bson:"customResponseHeaders" json:"customResponseHeaders" web:"自定义响应头"` //自定义响应头部
	DefaultUpstream       DefaultUpstream   `bson:"defaultUpstream" json:"defaultUpstream" web:"默认转发"`               //默认转发配置
	FeatureUpstreams      []FeatureUpstream `bson:"featureUpstreams" json:"featureUpstreams" web:"流量特征转发"`           //流量特征转发配置
	RatioUpstream         RatioUpstream     `bson:"ratioUpstream" json:"ratioUpstream" web:"流量比例转发"`                 //流量比例转发配置
}

//转发配置
type DefaultUpstream struct {
	CustomRequestHeaders map[string]string `bson:"customRequestHeaders" json:"customRequestHeaders" web:"自定义请求头部"` //自定义请求头部
	UpstreamName         string            `bson:"upstreamName" json:"upstreamName" web:"转发upstream名称"`            //转发upstream名称
	ProxyPassType        string            `bson:"proxyPassType" json:"proxyPassType" web:"目标服务类型"`                //server-指定服务器 、name-指定名称、neo sidecar
	NeoTransform         NeoTransform      `bson:"neoTransform" json:"neoTransform" web:"neo 转发配置"`                //neo转发配置
	Rewrites             []Rewrite         `bson:"rewrites" json:"rewrites" web:"rewrites"`                        //rewrites
	ServerId             bson.ObjectId     `bson:"serverId,omitempty" json:"serverId" web:"服务器 Id"`                //服务器Id
	EnsureEnv            bool              `json:"ensureEnv" bson:"ensureEnv" web:"指定环境"`                          // 指定环境 聚合的结果
	IdcEnv               map[string]string `json:"idcEnv" bson:"idcEnv"`                                           // 某个 idc 指定的环境
	TemplateEnv          TemplateEnv       `bson:"templateEnv" json:"templateEnv" web:"指定环境模版"`                    // 指定环境模版
	NameServerEnable     bool              `bson:"nameServerEnable" json:"nameServerEnable" web:"是否启用名字服务"`        //是否启用名字服务
}

type TemplateEnv struct {
	Using         bool      `bson:"using" json:"using"  web:"启用状态"`     //
	TemplateName  string    `json:"templateName" bson:"templateName"`   // 模版名称
	TemplateParam []HashKey `json:"templateParam" bson:"templateParam"` // 模版参数
	InjectParam   string    `json:"injectParam" bson:"injectParam"`     // 注入名称
}

//特征转发配置
type FeatureUpstream struct {
	DefaultUpstream `bson:"defaultUpstream" web:"默认转发配置"`         //默认转发配置
	Rules           []Rule `bson:"rules" json:"rules" web:"特征列表"` //特征列表
}

//流量比例转发配置
type RatioUpstream struct {
	HashKey              HashKey               `bson:"hashKey" json:"hashKey" web:"hashKey"`                             //hashKey
	RatioUpstreamConfigs []RatioUpstreamConfig `bson:"ratioUpstreamConfigs" json:"ratioUpstreamConfigs" web:"按流量比例转发配置"` //按流量比例转发配置
}
type RatioUpstreamConfig struct {
	Quota           int `bson:"quota" json:"quota" web:"流量比例"` //流量比例
	DefaultUpstream `bson:"defaultUpstream" web:"默认转发配置"`      //默认转发配置
}

type JanusConfig struct {
	Id              bson.ObjectId   `bson:"_id,omitempty" json:"id"`                           //
	Describe        string          `bson:"describe" json:"describe" web:"描述"`                 //配置描述信息
	LocationId      string          `bson:"locationId" json:"locationId"`                      //该字段只在生成nginx配置时使用
	BaseConfig      BaseConfig      `bson:"baseConfig" json:"baseConfig" web:"基础配置"`           //
	IpLocation      IpLocation      `bson:"ipLocation" json:"ipLocation" web:"地理位置"`           //
	PageCache       PageCache       `bson:"pageCache" json:"pageCache" web:"页面缓存"`             //
	Canary          Canary          `bson:"canary" json:"canary" web:"泳道"`                     //
	Waf             WafConfig       `bson:"waf" json:"waf" web:"waf"`                          //1
	ValidateRequest ValidateRequest `bson:"validateRequest" json:"validateRequest" web:"请求校验"` //
	TransProtocol   TransProtocol   `bson:"transProtocol" json:"transProtocol" web:"协议转换"`     //
	Compress        Compress        `bson:"compress" json:"compress" web:"压缩"`                 //
	Decompress      Decompress      `bson:"decompress" json:"decompress" web:"解压"`             //
	Crypt           Crypt           `bson:"crypt" json:"crypt" web:"加密"`                       //
	//Deprecated: can use ReqAnalysis instead
	Decrypt                Decrypt                `bson:"decrypt" json:"decrypt" web:"解密"`                                   //
	AntiSpider             AntiSpider             `bson:"antiSpider" json:"antiSpider" web:"反爬"`                             //
	Intercept              Intercept              `bson:"intercept" json:"intercept" web:"流量防护"`                             // 拦截策略
	CheckList              CheckList              `bson:"checkList" json:"checkList" web:"checkList"`                        // checkList日志
	Exception              Exception              `bson:"exception" json:"exception" web:"异常处理"`                             //exception
	RequestValidateSandbox RequestValidateSandbox `bson:"requestValidateSandbox" json:"requestValidateSandbox" web:"请求校验沙盒"` //请求校验沙盒
	NeoTransform           NeoTransform           `bson:"neoTransform" json:"neoTransform" web:"neo 转换"`                     //
	XssDefend              XssDefend              `json:"xssDefend" bson:"xssDefend" web:"xss 防护"`                           //
	ParamVerify            ParamVerify            `json:"paramVerify" bson:"paramVerify" web:"参数校验"`                         //
	ReqAnalysis            ReqAnalysis            `json:"reqAnalysis" bson:"reqAnalysis" web:"请求协议解析"`                       //请求协议解析                                 //
	Plugins                int64                  `json:"plugins" bson:"plugins"`                                            // 启用的功能
	CreateUserId           string                 `bson:"createUserId,omitempty" json:"createUserId"`                        //
	CreateUserName         string                 `bson:"createUserName,omitempty" json:"createUserName"`                    //
	CreateTime             time.Time              `bson:"createTime,omitempty" json:"createTime"`                            //创建时间
	TagMd5                 string                 `json:"tagMd5" bson:"tagMd5"`                                              // 打版本的 md5 值
	CurrentMd5             string                 `json:"currentMd5" bson:"currentMd5"`                                      // 当前 md5 值
	UpdateTime             time.Time              `bson:"updateTime,omitempty" json:"updateTime"`                            //更新时间
	UpdateUserId           string                 `bson:"updateUserId,omitempty"json:"updateUserId"`                         //更新人工号
	UpdateUserName         string                 `bson:"updateUserName,omitempty"json:"updateUserName"`                     //更新人姓名
	Delete                 bool                   `json:"delete" bson:"delete"`                                              // 是否删除
	Status                 int                    `json:"status" bson:"status"`                                              // 此条配置信息的状态 1 开启、2 关闭、3 删除
	Lock                   int64                  `json:"lock" bson:"lock"`                                                  // 0 可修改、1 流程审批再用不可修改
}

type JanusRouteWeb struct {
	JanusConfig

	NowMd5    string `json:"nowMd5"`
	NowStatus bool   `json:"nowStatus"`

	Modify     bool   `json:"modify"`     // 是否可修改
	ModifyInfo string `json:"modifyInfo"` // 不可修改原因
}

//动态配置信息
type DynamicConfig struct {
	ConfigId               string                 `bson:"configId" json:"configId"` //配置Id
	BaseConfig             BaseConfig             `bson:"baseConfig" json:"baseConfig"`
	IpLocation             IpLocation             `bson:"ipLocation" json:"ipLocation"`
	PageCache              PageCache              `bson:"pageCache" json:"pageCache"`
	Canary                 Canary                 `bson:"canary" json:"canary"`
	Waf                    WafConfig              `bson:"waf" json:"waf"`
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
	XssDefend              XssDefend              `json:"xssDefend" bson:"xssDefend" web:"xss 防护"`              //
	ParamVerify            ParamVerify            `json:"paramVerify" bson:"paramVerify" web:"参数校验"`            //
	ReqAnalysis            ReqAnalysis            `json:"reqAnalysis" bson:"reqAnalysis" web:"请求协议解析"`          //请求协议解析                                    //请求协议解析
}

//
type CheckList struct {
	Using          bool   `bson:"using" json:"using" web:"是否启用"`               //是否启用
	ProductLine    string `bson:"productLine" json:"productLine" web:"产品线"`    //产品线
	SamplingRate   int    `bson:"samplingRate"  json:"samplingRate" web:"采样率"` //采样率
	RecordHttpBody bool   `bson:"recordHttpBody" json:"recordHttpBody" web:"是否记录请求包体"`
}

type Exception struct {
	Using bool   `bson:"using" json:"using" web:"是否启用"` //是否启用
	Model string `bson:"model" json:"model" web:"产品线"`  //产品线
}

//cc策略
type Intercept struct {
	Using           bool            `bson:"using" json:"using" web:"是否启用"`        //是否启用
	WhiteStrategies []WhiteStrategy `bson:"whiteList" json:"whiteList" web:"白名单"` //白名单-放行策略
	Type            string          `bson:"type" json:"type" web:"限流类型"`          //限流类型   quota-比例  rate-限速
	Quota           QuotaIntercept  `bson:"quota" json:"quota" web:"比例拦截配置"`      //流量比例拦截
	Rate            RateIntercept   `bson:"rate" json:"rate" web:"速率拦截配置"`        //速率拦截
}

//特征白名单
type WhiteStrategy struct {
	StartTime int64  `bson:"startTime" json:"startTime" web:"开始时间"` //开始时间
	EndTime   int64  `bson:"endTime" json:"endTime" web:"结束时间"`     //结束时间
	Rules     []Rule `bson:"rules" json:"rules" web:"规则列表"`         //规则列表
}

//流量比例拦截
type QuotaIntercept struct {
	HashKey   HashKey           `bson:"hashKey" json:"hashKey" web:"比例计算 Key"` //比例计算hashKey
	Quota     int64             `bson:"quota" json:"quota" web:"流量占比"`         //流量占比
	StartTime int64             `bson:"startTime" json:"startTime" web:"开始时间"` //开始时间
	EndTime   int64             `bson:"endTime" json:"endTime" web:"结束时间"`     //结束时间
	Response  InterceptResponse `bson:"response" json:"response" web:"拦截响应结果"` //拦截响应结果
}

//速率拦截
type RateIntercept struct {
	Default DefaultStrategy   `bson:"default" json:"default" web:"默认拦截策略"` //默认拦截策略
	Feature []FeatureStrategy `bson:"feature" json:"feature" web:"特征拦截策略"` //特征拦截策略
}

//默认策略
type DefaultStrategy struct {
	Type         string            `bson:"type" json:"type" web:"限速类型"`                   //限速类型 static-静态  dynamic-动态
	StartTime    int64             `bson:"startTime" json:"startTime" web:"开始时间"`         //开始时间
	EndTime      int64             `bson:"endTime" json:"endTime" web:"结束时间"`             //结束时间
	Response     InterceptResponse `bson:"response" json:"response" web:"响应结果"`           //响应结果
	StaticLimit  StaticLimit       `bson:"staticLimit" json:"staticLimit" web:"静态限速配置"`   //静态限速配置
	DynamicLimit DynamicLimit      `bson:"dynamicLimit" json:"dynamicLimit" web:"动态限速配置"` //动态限速配置
}

//特征拦截策略
type FeatureStrategy struct {
	DefaultStrategy `bson:"defaultStrategy" web:"默认策略"`           //默认策略
	Rules           []Rule `bson:"rules" json:"rules" web:"规则列表"` //规则列表
}

//拦截响应结果
type InterceptResponse struct {
	Code        int    `bson:"code" json:"code" web:"拦截响应状态码"`              //拦截响应状态吗-httpCode
	Content     string `bson:"content" json:"content" web:"响应数据"`           //响应数据
	ContentType string `bson:"contentType" json:"contentType" web:"响应数据类型"` //响应contentType
}

//动态拦截配置
type DynamicLimit struct {
	AverageElapseTime int64 `json:"averageElapseTime" bson:"averageElapseTime" web:"平均消耗时长"` //平均耗时时长 单位：ms
	UpperLimitCount   int64 `json:"upperLimitCount" bson:"upperLimitCount" web:"最大请求数量"`     //请求数量上限 单位：个/s
	LowerLimitCount   int64 `json:"lowerLimitCount" bson:"lowerLimitCount" web:"最小请求数量"`     //请求数量下限 单位：个/s
	ElapseTimeRatio   int64 `json:"elapseTimeRatio" bson:"elapseTimeRatio" web:"耗时超过平均耗时占比"` //耗时超过平均耗时占比
	LoadFactor        int64 `json:"loadFactor" bson:"loadFactor" web:"调整幅度"`                 //增长/缩小比例
}

//静态拦截配置
type StaticLimit struct {
	LimitCount int64 `json:"limitCount" bson:"limitCount" web:"每秒访问量"` //限流速率  单位：个/s
}

//反扒配置
type AntiSpider struct {
	Using                bool                 `bson:"using" json:"using" web:"启用状态"`
	FlowFeature          string               `bson:"flowFeature" json:"flowFeature" web:"流量特征模板"` //流量特征提取模板 default、mapi、hotel-wechat
	FeatureItems         []FeatureItem        `bson:"featureItems" json:"featureItems" web:"请求特征"` //模板提取请求特征
	AppUniqueKey         string               `bson:"appUniqueKey"json:"appUniqueKey" web:"应用标示"`  //反爬
	SpiderIntercept      SpiderIntercept      `bson:"spiderIntercept" json:"spiderIntercept" web:"封禁"`
	ImportBlackHole      ImportBlackHole      `bson:"importBlackHole" json:"importBlackHole" web:"导黑洞"`
	Poison               Poison               `bson:"poison" json:"poison" web:"下毒"`
	SyncSpider           bool                 `bson:"syncSpider" json:"syncSpider" web:"同步神盾局"`
	AntiSpiderStrategies []AntiSpiderStrategy `bson:"antiSpiderStrategies" json:"antiSpiderStrategies" web:"异常流量转发"`
	CheckCode            AntiSpiderCheckCode  `json:"checkCode" bson:"checkCode" web:"验证码"` // 验证码
}

type AntiSpiderCheckCode struct {
	Using       bool   `bson:"using" json:"using" web:"启用状态"`
	RespType    int    `json:"respType" bson:"respType" web:"返回类型"`         // 1 重定向跳转 2 code 码
	Redirect    string `json:"redirect" bson:"redirect" web:"重定向地址"`        //
	Resp        string `bson:"resp" json:"resp" web:"返回数据"`                 // code 码返回数据
	ContentType string `bson:"contentType" json:"contentType" web:"返回数据格式"` // code 码返回数据格式
}

//特征反爬策略组
type AntiSpiderStrategy struct {
	Rules           []Rule          `bson:"rules" json:"rules" web:"特征组"`                       //特征组
	Type            int             `bson:"type" json:"type" web:"策略类型"`                        // 策略类型，1-拦截 2-导流量 0-不处理
	Action          int             `bson:"action" json:"action" web:"反爬动作"`                    //反爬动作 映射到神盾局的反爬动作 0-不处理 封禁-3 流量转发-2
	Result          int             `bson:"result" json:"result" web:"反爬结果"`                    //反爬结果   映射到神盾局的反爬结果
	Resp            string          `bson:"resp" json:"resp" web:"拦截响应结果"`                      //拦截响应结果
	ContentType     string          `bson:"contentType" json:"contentType" web:"相应类型"`          //响应content-type
	DefaultUpstream DefaultUpstream `bson:"defaultUpstream" json:"defaultUpstream" web:"导流量配置"` //导流转发配置
}

type FeatureItem struct {
	From   string `bson:"from" json:"from" web:"请求来源"`        //请求来源 cookie
	Name   string `bson:"name" json:"name" web:"字段名称"`        //字段名称
	Target string `bson:"target" json:"target" web:"映射神盾局字段"` //映射神盾局字段名称
}

//拦截动作
type SpiderIntercept struct {
	Using       bool   `bson:"using" json:"using" web:"启用状态"`
	Resp        string `bson:"resp" json:"resp" web:"返回数据"`
	ContentType string `bson:"contentType" json:"contentType" web:"返回数据格式"`
}

//导入黑洞
type ImportBlackHole struct {
	Using             bool               `bson:"using" json:"using" web:"启用状态"`
	AntiSpiderActions []AntiSpiderAction `bson:"antiSpiderActions" json:"antiSpiderActions" web:"导入黑洞策略"`
}

//下毒
type Poison struct {
	Using             bool               `bson:"using" json:"using" web:"启用状态"`
	AntiSpiderActions []AntiSpiderAction `bson:"antiSpiderActions" json:"antiSpiderActions" web:"下毒策略"`
}

type AntiSpiderAction struct {
	ValidateResults []string `bson:"validateResults" json:"validateResults" web:"反爬策略"` // -2-黑名单 1-白名单 0-正常 1判定为爬虫 2-疑似爬虫
	DefaultUpstream
}

type HashKey struct {
	KeyName string `bson:"keyName" json:"keyName" web:"keyName"`
	KeyType string `bson:"keyType" json:"keyType" web:"keyType"` // cookie|header|param
}

//特征泳道配置
type FeatureCanaryConfig struct {
	Rule         []Rule        `bson:"rule" json:"rule" web:"流量特征"`                      //流量特征
	CanaryId     bson.ObjectId `bson:"canaryId" json:"canaryId" web:"泳道 Id"`             //泳道Id
	CanaryName   string        `bson:"canaryName" json:"canaryName" web:"泳道名称"`          //泳道名称
	CanaryHeader string        `bson:"canaryHeader" json:"canaryHeader" web:"泳道 header"` //泳道header
}

//比例泳道配置
type RatioCanaryConfig struct {
	HashKey                HashKey                 `bson:"hashKey" json:"hashKey"  web:"hashKey"`                             // hash标识
	RatioCanaryConfigItems []RatioCanaryConfigItem `bson:"ratioCanaryConfigItems" json:"ratioCanaryConfigItems" web:"比例配置信息"` //比例配置信息
}

//比例配置详细信息
type RatioCanaryConfigItem struct {
	Quota        int           `bson:"quota" json:"quota" web:"流量比例"`                    //流量比例
	CanaryId     bson.ObjectId `bson:"canaryId" json:"canaryId" web:"泳道 Id"`             //泳道Id
	CanaryName   string        `bson:"canaryName" json:"canaryName" web:"泳道名称"`          //泳道名称
	CanaryHeader string        `bson:"canaryHeader" json:"canaryHeader" web:"泳道 header"` //泳道header
}

type Rule struct {
	Type    string `bson:"type" json:"type"  web:"类型"`
	KeyName string `bson:"keyName" json:"keyName" web:"keyName"`
	KeyType string `bson:"keyType" json:"keyType" web:"keyType"` // cookie|header|param|ip|city
	Value   string `bson:"value" json:"value" web:"value"`
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
	Using                bool                  `bson:"using" json:"using" web:"启用状态"`                                 //泳道是否启用
	FeatureCanaryConfigs []FeatureCanaryConfig `bson:"featureCanaryConfigs" json:"featureCanaryConfigs" web:"特征泳道配置"` //特征泳道配置
	RatioCanaryConfig    RatioCanaryConfig     `bson:"ratioCanaryConfig" json:"ratioCanaryConfig" web:"比例泳道配置"`       //比例泳道配置
}

type ValidateRequest struct {
	Using               bool   `bson:"using" json:"using" web:"启用状态"`                                     //请求合法性校验
	HeaderValidateType  string `bson:"headerValidateType" json:"headerValidateType" web:"请求头校验类型"`        //请求头校验类型  elong、tongcheng
	SessionValidateType string `bson:"sessionValidateType" json:"sessionValidateType" web:"Session 校验类型"` //session校验类型 ForceCheckStrategy,ForceUnCheckStrategy,ForceCheckWithVersionStrategy,ForceCheckWriteCardNoStrategy,ValidateCardNoWithSession
}

type TransProtocol struct {
	Using          bool        `bson:"using" json:"using" web:"启用状态"`                   //协议转换是否启用
	TargetProtocol string      `bson:"targetProtocol" json:"targetProtocol" web:"目标协议"` //目标协议
	TargetMethod   string      `bson:"targetMethod" json:"targetMethod" web:"目标方法"`     //目标方法
	ContentType    string      `bson:"contentType" json:"contentType" web:"数据类型"`       //
	ParamItems     []ParamItem `json:"paramItems" bson:"paramItems" web:"参数映射关系"`       //app请求协议转换参数映射
}

type ParamItem struct {
	Name           string `json:"name" bson:"name" web:"参数名"`
	TargetLocation string `json:"targetLocation" bson:"targetLocation" web:"目标参数位置" ` // 目标参数位置 path|param|header|body
	TargetName     string `json:"targetName" bson:"targetName" web:"目标参数名"`           // 目标参数位置为body时不用填，目标参数位置为path时为uri替换表达式
}

type Compress struct {
	Using        bool   `bson:"using" json:"using" web:"启用状态"`             //是否启用
	CompressType string `bson:"compressType" json:"compressType" web:"类型"` //类型
}
type Decompress struct {
	Using          bool   `bson:"using" json:"using" web:"启用状态"`
	DecompressType string `bson:"decompressType" json:"decompressType" web:"类型"`
}
type Crypt struct {
	Using     bool   `bson:"using" json:"using" web:"启用状态"`         //是否启用
	CryptType string `bson:"cryptType" json:"cryptType" web:"加密方法"` //加解密方法
}

//Deprecated: can use ReqAnalysis instead
type Decrypt struct {
	Using            bool     `bson:"using" json:"using" web:"启用状态"`
	DecryptType      string   `bson:"decryptType" json:"decryptType" web:"解密方法"`
	DecryptWhiteList []string `bson:"decryptWhiteList" json:"decryptWhiteList" web:"白名单"` //白名单
}

//请求校验沙盒
type RequestValidateSandbox struct {
	Using           bool            `bson:"using" json:"using" web:"启用状态"`                         //是否启用
	ValidateType    string          `bson:"validateType" json:"validateType" web:"校验类型"`           //校验类型
	AuthTokenConfig AuthTokenConfig `bson:"authTokenConfig" json:"authTokenConfig" web:"Token 来源"` //授权token来源信息
}

//中台接入
type NeoTransform struct {
	HttpMethod   string `bson:"httpMethod" json:"httpMethod" web:"Http Method"`
	ServiceGroup string `bson:"serviceGroup" json:"serviceGroup" web:"服务分组"`
	Service      string `bson:"service" json:"service" web:"服务器"`
	Api          string `bson:"api" json:"api" web:"Api"`
}

type XssDefend struct {
	Using bool      `bson:"using" json:"using"` //配置是否启用
	Rules []XssRule `json:"rules" bson:"rules" web:"规则"`
}

type XssRule struct {
	Type    string `bson:"type" json:"type" web:"类型"`            // eq ne in
	KeyType string `bson:"keyType" json:"keyType" web:"keyType"` // param|form-data
	KeyName string `bson:"keyName" json:"keyName" web:"keyName"`
	Value   string `bson:"value" json:"value" web:"value"`
	Action  int    `json:"action" bson:"action" web:"行为"` // 1 过滤、2 html 转义、3 js 转义
}

type ParamVerify struct {
	Using       bool                   `bson:"using" json:"using"`                     //配置是否启用
	VerifyRules map[string]VerifyRules `json:"verifyRules" bson:"verifyRules"`         // 参数校验规则 key-path 唯一
	RespCode    int64                  `bson:"respCode" json:"respCode" web:"拦截响应状态码"` //拦截响应状态吗
	Resp        string                 `bson:"resp" json:"resp" web:"返回数据"`
	ContentType string                 `bson:"contentType" json:"contentType" web:"返回数据格式"`
}

type VerifyRules struct {
	KeyName   string `bson:"keyName" json:"keyName" web:"keyName"`
	KeyType   string `bson:"keyType" json:"keyType" web:"keyType"` // param|form-data|header
	Rules     []Rule `json:"rules" bson:"rules"`                   // 规则列表
	ValueType int    `json:"valueType" bson:"valueType"`           // 1 字符串 2 数字
}

//静态配置，用于生成nginx Location配置
type LocationConfig struct {
	Using           bool          `bson:"using" json:"using"`                           //配置是否启用
	LocationId      string        `bson:"locationId" json:"locationId"`                 //该字段只在生成nginx配置时使用
	MatchType       string        `bson:"matchType" json:"matchType"`                   //uri匹配类型    equal-全匹配 prefix-前缀匹配 regex-正则匹配
	Uri             string        `bson:"uri" json:"uri"`                               //location Uri
	Hosts           []string      `bson:"hosts" json:"hosts"`                           //生成配置serverName
	DefaultUpstream `bson:"defaultUpstream"`                                            //默认转发配置
	SslEnable       bool          `bson:"sslEnable" json:"sslEnable"`                   //是否支持Https
	Http2Enable     bool          `bson:"http2Enable" json:"http2Enable"`               //是否支持Http2
	CertificateId   bson.ObjectId `bson:"certificateId,omitempty" json:"certificateId"` //证书Id
	CertificateName string        `bson:"certificateName"json:"certificateName"`        //证书名称
	ConnectTimeOut  int64         `bson:"connectTimeOut" json:"connectTimeOut"`         //连接超时时间 单位ms
	ReadTimeOut     int64         `bson:"readTimeOut" json:"readTimeOut"`               //读取超时时间 单位ms
	SendTimeOut     int64         `bson:"sendTimeOut" json:"sendTimeOut"`               //发送超时时间 单位ms
}

/**
 * 授权token校验模板
 */
type AuthTokenConfig struct {
	Token  FeatureItem `bson:"token" json:"token" web:"Token 特征提取信息"` //token 特征提取信息
	Source FeatureItem `bson:"source" json:"source" web:"请求来源信息"`     //source 请求来源提取信息
}

type WafConfig struct {
	Using       bool   `bson:"using" json:"using" web:"启用状态"`               //配置是否启用
	RespCode    int64  `bson:"respCode" json:"respCode" web:"拦截响应状态码"`      //拦截响应状态吗
	RespContent string `bson:"respContent" json:"respContent" web:"拦截响应内容"` //拦截响应内容
	ContentType string `bson:"contentType"  json:"contentType" web:"响应类型"`  //contentType
}

type ReqAnalysis struct {
	Using    bool   `json:"using" bson:"using" web:"启用状态"`
	Protocol string `json:"protocol" bson:"protocol" web:"协议"`
}
type Response struct {
	Code   int          `json:"code"`   //返回code 0 表示成功
	Msg    string       `json:"msg"`    //描述信息
	Result pageResponse `json:"result"` //结果数据
}

type ObjectResponse struct {
	Code   int         `json:"code"`   //返回code 0 表示成功
	Msg    string      `json:"msg"`    //描述信息
	Result interface{} `json:"result"` //结果数据
}

type ConfigDetailResponse struct {
	Code   int           `json:"code"`   //返回code 0 表示成功
	Msg    string        `json:"msg"`    //描述信息
	Result JanusRouteWeb `json:"result"` //结果数据
}

type pageResponse struct {
	PageData   []RouteList `json:"pageData"`
	TotalCount int         `json:"totalCount"`
}

type RouteList struct {
	Id         bson.ObjectId `json:"id" bson:"_id"`
	BaseConfig struct {
		Using bool     `json:"using" bson:"using"`
		Uri   string   `json:"uri" bson:"uri"`
		Hosts []string `json:"hosts" bson:"hosts"`

		DefaultUpstream struct {
			UpstreamName string `json:"upstreamName" bson:"upstreamName"`
		} `json:"defaultUpstream" bson:"defaultUpstream"`
	} `json:"baseConfig" bson:"baseConfig"`

	Describe       string    `json:"describe" bson:"describe"`
	Plugins        int64     `json:"plugins" bson:"plugins"` // 启用的功能
	Delete         bool      `json:"delete" bson:"delete"`   // 是否删除
	Status         int       `json:"status" bson:"status"`   // 此条配置信息的状态 1 开启、2 关闭、3 删除
	CreateUserName string    `json:"createUserName" bson:"createUserName"`
	CreateUserId   string    `json:"createUserId" bson:"createUserId"`
	UpdateUserId   string    `json:"updateUserId" bson:"updateUserId"`
	UpdateUserName string    `json:"updateUserName" bson:"updateUserName"`
	UpdateTime     time.Time `json:"updateTime" bson:"updateTime"`
	//Lock           int       `json:"lock" bson:"lock"`
	ProcessStatus int `json:"processStatus" bson:"processStatus"`
}

func QueryAllConfigs(url string, headers map[string]string) (baseConfigs []RouteList, err error) {
	request := httplib.Get(url)
	for key, value := range headers {
		request.Header(key, value)
	}
	response, err := request.DoRequest()

	if err != nil {
		return []RouteList{}, err
	}
	defer response.Body.Close()
	var data []byte
	data, err = ioutil.ReadAll(response.Body)

	resp := new(Response)
	err = json.Unmarshal(data, resp)
	if err != nil {
		return []RouteList{}, err
	}
	if resp.Code != 0 {
		return []RouteList{}, errors.New(resp.Msg)
	}
	return resp.Result.PageData, nil
}

func QueryConfigDetail(url string, headers map[string]string) (configDetail JanusRouteWeb, err error) {
	request := httplib.Get(url)
	for key, value := range headers {
		request.Header(key, value)
	}
	response, err := request.DoRequest()

	if err != nil {
		return JanusRouteWeb{}, err
	}
	defer response.Body.Close()
	var data []byte
	data, err = ioutil.ReadAll(response.Body)

	resp := new(ConfigDetailResponse)
	err = json.Unmarshal(data, resp)
	if err != nil {
		return JanusRouteWeb{}, err
	}
	if resp.Code != 0 {
		return JanusRouteWeb{}, errors.New(resp.Msg)
	}
	return resp.Result, nil
}

func UpdateConfig(url string, headers map[string]string, param interface{}) (err error) {
	request := httplib.Put(url)
	for key, value := range headers {
		request.Header(key, value)
	}

	data, err := json.Marshal(param)
	request.Body(data)
	response, err := request.DoRequest()

	if err != nil {
		return err
	}
	defer response.Body.Close()
	data, err = ioutil.ReadAll(response.Body)

	resp := new(ObjectResponse)
	err = json.Unmarshal(data, resp)
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		return errors.New(resp.Msg)
	}
	return nil
}
