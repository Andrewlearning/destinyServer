package handler

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"destinyServer/config"
	"destinyServer/store"
	"destinyServer/wechat"
)

// POST /api/pay/create  { "open_id": "..." }
func HandlePayCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResp(w, 405, map[string]string{"message": "method not allowed"})
		return
	}

	var req struct {
		OpenID string `json:"open_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.OpenID == "" {
		jsonResp(w, 400, map[string]string{"message": "invalid request"})
		return
	}

	outTradeNo := fmt.Sprintf("DST%d%s", time.Now().UnixMilli(), req.OpenID[:8])
	clientIP := r.RemoteAddr

	if err := store.CreateOrder(req.OpenID, outTradeNo, config.Cfg.PriceFen); err != nil {
		jsonResp(w, 500, map[string]string{"message": "create order failed"})
		return
	}

	prepayID, err := wechat.UnifiedOrder(req.OpenID, outTradeNo, clientIP, config.Cfg.PriceFen)
	if err != nil {
		jsonResp(w, 500, map[string]string{"message": "wx pay failed: " + err.Error()})
		return
	}

	timeStamp := fmt.Sprintf("%d", time.Now().Unix())
	nonceStr := outTradeNo[:16]
	pkg := "prepay_id=" + prepayID
	paySign := wechat.Sign(map[string]string{
		"appId":     config.Cfg.WxAppID,
		"timeStamp": timeStamp,
		"nonceStr":  nonceStr,
		"package":   pkg,
		"signType":  "MD5",
	})

	jsonResp(w, 200, map[string]any{
		"pay_params": map[string]string{
			"timeStamp": timeStamp,
			"nonceStr":  nonceStr,
			"package":   pkg,
			"signType":  "MD5",
			"paySign":   paySign,
		},
	})
}

// POST /api/pay/notify  (微信支付回调)
func HandlePayNotify(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var notify wechat.NotifyReq
	if err := xml.Unmarshal(body, &notify); err != nil {
		w.Write([]byte("<xml><return_code>FAIL</return_code></xml>"))
		return
	}

	if notify.ReturnCode == "SUCCESS" && notify.ResultCode == "SUCCESS" {
		_ = store.CompleteOrder(notify.OutTradeNo)
	}

	w.Write([]byte("<xml><return_code>SUCCESS</return_code><return_msg>OK</return_msg></xml>"))
}
