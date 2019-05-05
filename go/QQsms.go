package main

import (
	qcloudsms "github.com/qichengzx/qcloudsms_go"
	"strings"
)

var (
	appid      string = "1400207787"
	appkey     string = "d1d70b6fc1853bcf107258cdc723b7ea"
	sign       string = "五毛的技术" //发出后自动拼接成 【】
	phoneList  string = "17600697663"
	templateid int    = 324359 //申请的模板ID
)

func main() {
	params := []string{"金石小区"}
	MultiSend(params, strings.Split(phoneList, ","))
}

func MultiSend(params []string, tl []string) {
	opt := qcloudsms.NewOptions(appid, appkey, sign)
	opt.Debug = true

	var Tel []qcloudsms.SMSTel
	st := Tel[:]
	for _, tel := range tl {
		st = append(st, qcloudsms.SMSTel{Nationcode: "86", Mobile: tel})
	}

	var client = qcloudsms.NewClient(opt)

	var sm = qcloudsms.SMSMultiReq{
		Type:   0,
		Params: params,
		Tel:    st,
		//TplID:  324318,
		TplID: 324359,
		Sign:  sign,
	}

	client.SendSMSMulti(sm)
}

func SingleSend() {
	opt := qcloudsms.NewOptions(appid, appkey, sign)
	opt.Debug = true

	var client = qcloudsms.NewClient(opt)

	params := []string{}

	var sm = qcloudsms.SMSSingleReq{
		Type:   0,
		Params: params,
		Tel:    qcloudsms.SMSTel{Nationcode: "86", Mobile: "13436967911"},
		TplID:  324318,
		Sign:   sign,
	}

	client.SendSMSSingle(sm)
}
