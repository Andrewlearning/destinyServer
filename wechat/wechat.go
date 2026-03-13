package wechat

import (
	"crypto/md5"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"destinyServer/config"
)

type Code2SessionResp struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func Code2Session(code string) (*Code2SessionResp, error) {
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		config.Cfg.WxAppID, config.Cfg.WxAppSecret, code,
	)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result Code2SessionResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		return nil, fmt.Errorf("wx error: %d %s", result.ErrCode, result.ErrMsg)
	}
	return &result, nil
}

type UnifiedOrderReq struct {
	XMLName        xml.Name `xml:"xml"`
	AppID          string   `xml:"appid"`
	MchID          string   `xml:"mch_id"`
	NonceStr       string   `xml:"nonce_str"`
	Sign           string   `xml:"sign"`
	Body           string   `xml:"body"`
	OutTradeNo     string   `xml:"out_trade_no"`
	TotalFee       int      `xml:"total_fee"`
	SpbillCreateIP string   `xml:"spbill_create_ip"`
	NotifyURL      string   `xml:"notify_url"`
	TradeType      string   `xml:"trade_type"`
	OpenID         string   `xml:"openid"`
}

type UnifiedOrderResp struct {
	ReturnCode string `xml:"return_code"`
	ResultCode string `xml:"result_code"`
	PrepayID   string `xml:"prepay_id"`
	ErrCodeDes string `xml:"err_code_des"`
}

type NotifyReq struct {
	ReturnCode string `xml:"return_code"`
	ResultCode string `xml:"result_code"`
	OutTradeNo string `xml:"out_trade_no"`
	Sign       string `xml:"sign"`
}

func Sign(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "sign" && params[k] != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(params[k])
	}
	buf.WriteString("&key=")
	buf.WriteString(config.Cfg.WxMchAPIKey)

	hash := md5.Sum([]byte(buf.String()))
	return strings.ToUpper(fmt.Sprintf("%x", hash))
}

func UnifiedOrder(openID, outTradeNo, clientIP string, totalFen int) (string, error) {
	nonceStr := outTradeNo[:16]
	params := map[string]string{
		"appid":            config.Cfg.WxAppID,
		"mch_id":           config.Cfg.WxMchID,
		"nonce_str":        nonceStr,
		"body":             "职业命运-详细分析",
		"out_trade_no":     outTradeNo,
		"total_fee":        fmt.Sprintf("%d", totalFen),
		"spbill_create_ip": clientIP,
		"notify_url":       config.Cfg.WxNotifyURL,
		"trade_type":       "JSAPI",
		"openid":           openID,
	}
	sign := Sign(params)

	reqBody := UnifiedOrderReq{
		AppID:          config.Cfg.WxAppID,
		MchID:          config.Cfg.WxMchID,
		NonceStr:       nonceStr,
		Sign:           sign,
		Body:           "职业命运-详细分析",
		OutTradeNo:     outTradeNo,
		TotalFee:       totalFen,
		SpbillCreateIP: clientIP,
		NotifyURL:      config.Cfg.WxNotifyURL,
		TradeType:      "JSAPI",
		OpenID:         openID,
	}

	xmlData, _ := xml.Marshal(reqBody)
	resp, err := http.Post(
		"https://api.mch.weixin.qq.com/pay/unifiedorder",
		"application/xml",
		strings.NewReader(string(xmlData)),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result UnifiedOrderResp
	if err := xml.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if result.ReturnCode != "SUCCESS" || result.ResultCode != "SUCCESS" {
		return "", fmt.Errorf("unified order failed: %s", result.ErrCodeDes)
	}
	return result.PrepayID, nil
}
