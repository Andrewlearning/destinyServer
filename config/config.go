package config

type Config struct {
	Port         string
	WxAppID      string
	WxAppSecret  string
	WxMchID      string // 微信支付商户号
	WxMchAPIKey  string // 微信支付API密钥
	WxNotifyURL  string // 支付回调地址
	PriceFen     int    // 单次分析价格（分），50 = 0.5元
	InitFreeUses int    // 新用户初始免费次数
}

var Cfg = Config{
	Port:         ":8080",
	WxAppID:      "wxYOUR_APP_ID",
	WxAppSecret:  "YOUR_APP_SECRET",
	WxMchID:      "YOUR_MCH_ID",
	WxMchAPIKey:  "YOUR_MCH_API_KEY",
	WxNotifyURL:  "https://your-server.com/api/pay/notify",
	PriceFen:     50,
	InitFreeUses: 1,
}
