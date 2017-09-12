package openbestpay

import (
	"testing"

	"github.com/liteck/logs"
)

//测试 付款码支付
func Test_bestpay_barcode_placeorder(t *testing.T) {
	api := GetApi(BESTPAY_URL_BARCODE_PLACEORDER)
	if err := api.SetBizContent(Biz_bestpay_barcode_placeorder{
		MerchantId:    "043101180050000",
		SubMerchantId: "043101180050009",
		Barcode:       "515665002854886972",
		OrderNo:       "14337346095601",
		OrderReqNo:    "14337346095601",
		OrderDate:     "20150608113649",
		OrderAmt:      1,
		ProductAmt:    1,
		AttachAmt:     0,
		GoodsName:     "你好",
		StoreId:       "201231",
	}, "1"); err != nil {
		logs.Error(err.Error())
		return
	}

	api.Run()
}

//测试 交易查询
func Test_bestpay_queryorder(t *testing.T) {
	api := GetApi(BESTPAY_URL_QUERYORDER)
	if err := api.SetBizContent(Biz_bestpay_queryorder{
		MerchantId: "043101180050000",
		OrderNo:    "14337346095601",
		OrderReqNo: "14337346095601",
		OrderDate:  "20150608113649",
	}, "1"); err != nil {
		logs.Error(err.Error())
		return
	}

	api.Run()
}

//测试 交易退款
func Test_bestpay_commonrefund(t *testing.T) {
	api := GetApi(BESTPAY_URL_COMMONREFUND)
	if err := api.SetBizContent(Biz_bestpay_commonrefund{
		MerchantId:    "043101180050000",
		SubMerchantId: "043101180050009",
	}, "1"); err != nil {
		logs.Error(err.Error())
		return
	}

	api.Run()
}
