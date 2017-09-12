package openbestpay

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

/**
翼支付交易模块
*/

/**
付款码支付
https://webpaywg.bestpay.com.cn/barcode/placeOrder
用户展示付款码.商户扫描用户付款码
*/
type bestpay_barcode_placeorder struct {
	BestpayApi
}

func (a *bestpay_barcode_placeorder) apiMethod() string {
	return BESTPAY_URL_BARCODE_PLACEORDER
}

func (a *bestpay_barcode_placeorder) apiName() string {
	return "付款码支付"
}

type Biz_bestpay_barcode_placeorder struct {
	MerchantId    string `json:"merchantId,omitempty"`        //由翼支付网关平台统一分配 30
	SubMerchantId string `json:"subMerchantId,omitempty"`     //由商户平台自己分配 30
	Barcode       string `json:"barcode,omitempty"`           //商户POS扫描用户客户端条形码 30
	OrderNo       string `json:"orderNo,omitempty"`           //由商户平台提供，支持纯数字、纯字母、字 母+数字组成，全局唯一(如果需要使用条 码退款业务，订单号必须为偶数位) 30
	OrderReqNo    string `json:"orderReqNo,omitempty"`        //同上
	Channel       string `json:"channel,omitempty"`           //默认填:05
	BusiType      string `json:"busiType,omitempty"`          //默认填:0000001
	OrderDate     string `json:"orderDate,omitempty"`         //由商户提供，长度14位，格式 yyyyMMddhhmmss (说明:该时间必须为 )
	OrderAmt      int    `json:"orderAmt,omitempty,string"`   //单位:分。订单总金额 = 产品金额+附加金 额
	ProductAmt    int    `json:"productAmt,omitempty,string"` //单位:分。
	AttachAmt     int    `json:"attachAmt,omitempty,string"`  //单位:分。
	GoodsName     string `json:"goodsName,omitempty"`         //商品信息 256
	StoreId       string `json:"storeId,omitempty"`           //门店号 10
	BackUrl       string `json:"backUrl,omitempty"`           //商户提供的用于异步接收交易返回结果的后 台url，若不需要后台返回，可不填，若需要 后台返回，请保障地址可用 255
	LedgerDetail  string `json:"ledgerDetail,omitempty"`      //商户需要在结算时进行分账情况，需填写此字段，详情见接口说明分账明细 256
	Attach        string `json:"attach,omitempty"`            //商户附加信息 128
	Mac           string `json:"mac,omitempty"`               //采用标准的MD5算法，由商户实现， MD5 加密获得32位大写字符 32
	MchntTmNum    string `json:"mchntTmNum,omitempty"`        //商户自定义终端号 50
	DeviceTmNum   string `json:"deviceTmNum,omitempty"`       //设备终端号 50
	ErpNo         string `json:"erpNo,omitempty"`             //商户营业员 编号 64
	GoodsDetail   []struct {
		GoodsId       string `json:"goodsId,omitempty"`         // 商品的编号
		GoodsName     string `json:"goodsName,omitempty"`       // 商品名称
		Quantity      int    `json:"quantity,omitempty,string"` // 商品数量
		Price         int    `json:"price,omitempty,string"`    // 商品价格
		GoodsCategory string `json:"goodsCategory,omitempty"`   // 商品分类
		Body          string `json:"body,omitempty"`            // 商品描述
	} `json:"goodsDetail,omitempty"` //商品详情，以 json 格式传过来，详见说明 5.2.4 4000

}

//mac 校验域.看起来像是一个请求签名的动作
//返回一个待 mac 的数据
func (b Biz_bestpay_barcode_placeorder) tobe_mac() string {
	tobe_mac := "MERCHANTID=" + b.MerchantId
	tobe_mac += "&ORDERNO=" + b.OrderNo
	tobe_mac += "&ORDERREQNO=" + b.OrderReqNo
	tobe_mac += "&ORDERDATE=" + b.OrderDate
	tobe_mac += "&BARCODE=" + b.Barcode
	tobe_mac += fmt.Sprintf("&%s=%d", "ORDERAMT", b.OrderAmt)

	return tobe_mac
}

func (b Biz_bestpay_barcode_placeorder) valid() error {
	if v := len(b.MerchantId); v == 0 || v > 30 {
		return errors.New("merchantId " + FORAMT_ERROR)
	}

	if v := len(b.SubMerchantId); v > 30 {
		return errors.New("subMerchantId " + FORAMT_ERROR)
	}

	if v := len(b.Barcode); v == 0 || v > 30 {
		return errors.New("barcode " + FORAMT_ERROR)
	}

	if v := len(b.OrderNo); v == 0 || v > 30 || v%2 != 0 {
		return errors.New("orderNo " + FORAMT_ERROR)
	}

	if v := len(b.OrderReqNo); v == 0 || v > 30 || v%2 != 0 {
		return errors.New("orderReqNo " + FORAMT_ERROR)
	}

	b.Channel = "05"
	b.BusiType = "0000001"

	if _, err := time.Parse("20060102150405", b.OrderDate); err != nil {
		return errors.New("orderDate " + FORAMT_ERROR)
	}

	if b.OrderAmt <= 0 {
		return errors.New("orderAmt " + FORAMT_ERROR)
	}

	if b.ProductAmt <= 0 {
		return errors.New("productAmt " + FORAMT_ERROR)
	}

	if b.AttachAmt < 0 {
		return errors.New("attachAmt " + FORAMT_ERROR)
	}

	if b.AttachAmt+b.ProductAmt != b.OrderAmt {
		return errors.New("orderAmt = productAmt + attachAmt")
	}

	if v := len(b.GoodsName); v > 256 {
		return errors.New("goodsName " + FORAMT_ERROR)
	}

	if v := len(b.StoreId); v == 0 || v > 10 {
		return errors.New("storeId " + FORAMT_ERROR)
	}

	if v := len(b.BackUrl); v > 255 {
		return errors.New("backUrl " + FORAMT_ERROR)
	}

	if v := len(b.LedgerDetail); v > 255 {
		return errors.New("ledgerDetail " + FORAMT_ERROR)
	}

	if v := len(b.Attach); v > 128 {
		return errors.New("attach " + FORAMT_ERROR)
	}

	//b.Mac 不做校验..这是一个类似签名的东西

	if v := len(b.MchntTmNum); v > 50 {
		return errors.New("mchntTmNum " + FORAMT_ERROR)
	}

	if v := len(b.DeviceTmNum); v > 50 {
		return errors.New("deviceTmNum " + FORAMT_ERROR)
	}

	if v := len(b.ErpNo); v > 64 {
		return errors.New("erpNo " + FORAMT_ERROR)
	}

	if len(b.GoodsDetail) > 0 {
		if v, err := json.Marshal(&b.GoodsDetail); err != nil {
			return errors.New("goodsDetail " + FORAMT_ERROR)
		} else if len(v) > 4000 {
			return errors.New("goodsDetail too long")
		}
		//TODO:目前文档中并没有看到 这个结构中的字段是否 都是必填.
	}

	return nil
}

type Resp_bestpay_barcode_placeorder struct {
	MerchantId   string `json:"merchantId,omitempty"`      //由翼支付网关平台统一分配 30
	OrderNo      string `json:"orderNo,omitempty"`         //由商户平台提供，支持纯数字、纯字母、字 母+数字组成，全局唯一(如果需要使用条 码退款业务，订单号必须为偶数位) 30
	OrderReqNo   string `json:"orderReqNo,omitempty"`      //同上
	OrderDate    string `json:"orderDate,omitempty"`       //由商户提供，长度14位，格式 yyyyMMddhhmmss (说明:该时间必须为 )
	OurTransNo   string `json:"ourTransNo,omitempty"`      //翼支付生成的内部流水号(用户支付后生成) 30
	TransAmt     int    `json:"transAmt,omitempty,string"` //单位:分。订单总金额 = 产品金额+附加金 额
	TransStatus  string `json:"transStatus,omitempty"`     //A:请求(支付中) B:成功(支付成功) C:失败(订单状态结果)
	EncodeType   string `json:"encodeType,omitempty"`      //1代表MD5; 3代表RSA;9代表CA;默认为1
	Sign         string `json:"sign,omitempty"`            //十六进制
	Coupon       int    `json:"coupon,omitempty,string"`   //单位:分。 订单优惠金额，用户使用代金券或立减的金额，金额为分
	ScValue      int    `json:"scValue,omitempty,string"`  //单位:分。 商户营销优惠成本
	PayerAccount string `json:"payerAccount,omitempty"`    //付款人账 号 30
	PayeeAccount string `json:"payeeAccount,omitempty"`    //收款人账 号 30
	PayChannel   string `json:"payChannel,omitempty"`      //付款明细 30
	ProductDesc  string `json:"productDesc,omitempty"`     //备注
	RefundFlag   string `json:"refundFlag,omitempty"`      //退款标示
	CustomerId   string `json:"customerId,omitempty"`      //客户登陆 账号
	MchntTmNum   string `json:"mchntTmNum,omitempty"`      //商户自定义终端号 50
	DeviceTmNum  string `json:"deviceTmNum,omitempty"`     //设备终端号 50
	Attach       string `json:"attach,omitempty"`          //商户附加信息 128
	TransPhone   string `json:"transPhone,omitempty"`      //商户附加信息 128
}

/**
交易查询
https://webpaywg.bestpay.com.cn/query/queryOrder
交易查询
*/
type bestpay_queryorder struct {
	BestpayApi
}

func (a *bestpay_queryorder) apiMethod() string {
	return BESTPAY_URL_QUERYORDER
}

func (a *bestpay_queryorder) apiName() string {
	return "交易查询"
}

type Biz_bestpay_queryorder struct {
	MerchantId string `json:"merchantId,omitempty"` //由翼支付网关平台统一分配 30
	OrderNo    string `json:"orderNo,omitempty"`    //由商户平台提供，支持纯数字、纯字母、字 母+数字组成，全局唯一(如果需要使用条 码退款业务，订单号必须为偶数位) 30
	OrderReqNo string `json:"orderReqNo,omitempty"` //同上
	OrderDate  string `json:"orderDate,omitempty"`  //由商户提供，长度14位，格式 yyyyMMddhhmmss (说明:该时间必须为 )
	Mac        string `json:"mac,omitempty"`        //采用标准的MD5算法，由商户实现， MD5 加密获得32位大写字符 32
}

//mac 校验域.看起来像是一个请求签名的动作
//返回一个待 mac 的数据
func (b Biz_bestpay_queryorder) tobe_mac() string {
	tobe_mac := "MERCHANTID=" + b.MerchantId
	tobe_mac += "&ORDERNO=" + b.OrderNo
	tobe_mac += "&ORDERREQNO=" + b.OrderReqNo
	tobe_mac += "&ORDERDATE=" + b.OrderDate

	return tobe_mac
}

func (b Biz_bestpay_queryorder) valid() error {
	if v := len(b.MerchantId); v == 0 || v > 30 {
		return errors.New("merchantId " + FORAMT_ERROR)
	}

	if v := len(b.OrderNo); v == 0 || v > 30 || v%2 != 0 {
		return errors.New("orderNo " + FORAMT_ERROR)
	}

	if v := len(b.OrderReqNo); v == 0 || v > 30 || v%2 != 0 {
		return errors.New("orderReqNo " + FORAMT_ERROR)
	}

	if _, err := time.Parse("20060102150405", b.OrderDate); err != nil {
		return errors.New("orderDate " + FORAMT_ERROR)
	}
	//b.Mac 不做校验..这是一个类似签名的东西
	return nil
}

type Resp_bestpay_queryorder struct {
	MerchantId   string `json:"merchantId,omitempty"`      //由翼支付网关平台统一分配 30
	OrderNo      string `json:"orderNo,omitempty"`         //由商户平台提供，支持纯数字、纯字母、字 母+数字组成，全局唯一(如果需要使用条 码退款业务，订单号必须为偶数位) 30
	OrderReqNo   string `json:"orderReqNo,omitempty"`      //同上
	OrderDate    string `json:"orderDate,omitempty"`       //由商户提供，长度14位，格式 yyyyMMddhhmmss (说明:该时间必须为 )
	OurTransNo   string `json:"ourTransNo,omitempty"`      //翼支付生成的内部流水号(用户支付后生成) 30
	TransAmt     int    `json:"transAmt,omitempty,string"` //单位:分。订单总金额 = 产品金额+附加金 额
	TransStatus  string `json:"transStatus,omitempty"`     //A:请求(支付中) B:成功(支付成功) C:失败(订单状态结果)
	EncodeType   string `json:"encodeType,omitempty"`      //1代表MD5; 3代表RSA;9代表CA;默认为1
	Sign         string `json:"sign,omitempty"`            //十六进制
	RefundFlag   string `json:"refundFlag,omitempty"`      //退款标示
	CustomerId   string `json:"customerId,omitempty"`      //客户登陆 账号
	Coupon       int    `json:"coupon,omitempty,string"`   //单位:分。 订单优惠金额，用户使用代金券或立减的金额，金额为分
	ScValue      int    `json:"scValue,omitempty,string"`  //单位:分。 商户营销优惠成本
	PayerAccount string `json:"payerAccount,omitempty"`    //付款人账 号 30
	PayeeAccount string `json:"payeeAccount,omitempty"`    //收款人账 号 30
	PayChannel   string `json:"payChannel,omitempty"`      //付款明细 30
	ProductDesc  string `json:"productDesc,omitempty"`     //备注
}

/**
交易退款
https://webpaywg.bestpay.com.cn/refund/commonRefund
交易退款
*/
type bestpay_commonrefund struct {
	BestpayApi
}

func (a *bestpay_commonrefund) apiMethod() string {
	return BESTPAY_URL_COMMONREFUND
}

func (a *bestpay_commonrefund) apiName() string {
	return "交易退款"
}

type Biz_bestpay_commonrefund struct {
	MerchantId    string `json:"merchantId,omitempty"`      //由翼支付网关平台统一分配 30
	SubMerchantId string `json:"subMerchantId,omitempty"`   //由商户平台自己分配 30
	MerchantPwd   string `json:"merchantPwd,omitempty"`     //商户执行时需填入相应密码 ，又称:交易key
	OldOrderNo    string `json:"oldOrderNo,omitempty"`      //原扣款成功的订单号 30
	OldOrderReqNo string `json:"oldOrderReqNo,omitempty"`   //原扣款成功的请求支付流水号
	RefundReqNo   string `json:"refundReqNo,omitempty"`     //该流水在商户处必须唯一。退款流水 refundReqNo不能和支付流水oldOrderNo 相同。若存在部分退款场景，具体见说明8.5原扣款成功的请求支付流水号
	RefundReqDate string `json:"refundReqDate,omitempty"`   //yyyyMMDD
	TransAmt      int    `json:"transAmt,omitempty,string"` //单位为分，小于等于原订单金额
	LedgerDetail  string `json:"ledgerDetail,omitempty"`    //商户需要在结算时进行分账情况，需填写此字段，详情见接口说明分账明细 256
	Channel       string `json:"channel,omitempty"`         //默认填:05
	Mac           string `json:"mac,omitempty"`             //采用标准的MD5算法，由商户实现， MD5 加密获得32位大写字符 32
	BgUrl         string `json:"bgUrl,omitempty"`           //商户的退款回调地址，当退款受理 255
}

//mac 校验域.看起来像是一个请求签名的动作
//返回一个待 mac 的数据
func (b Biz_bestpay_commonrefund) tobe_mac() string {
	tobe_mac := "MERCHANTID=" + b.MerchantId
	tobe_mac += "&MERCHANTPWD=" + b.MerchantPwd
	tobe_mac += "&OLDORDERNO=" + b.OldOrderNo
	tobe_mac += "&OLDORDERREQNO=" + b.OldOrderReqNo
	tobe_mac += "&REFUNDREQNO=" + b.RefundReqNo
	tobe_mac += "&REFUNDREQDATE=" + b.RefundReqDate
	tobe_mac += "&TRANSAMT=" + fmt.Sprintf("%d", b.TransAmt)
	tobe_mac += "&LEDGERDETAIL=" + b.LedgerDetail
	return tobe_mac
}

func (b Biz_bestpay_commonrefund) valid() error {
	if v := len(b.MerchantId); v == 0 || v > 30 {
		return errors.New("merchantId " + FORAMT_ERROR)
	}

	if v := len(b.SubMerchantId); v > 30 {
		return errors.New("subMerchantId " + FORAMT_ERROR)
	}

	if v := len(b.MerchantPwd); v == 0 || v > 20 {
		return errors.New("merchantPwd " + FORAMT_ERROR)
	}

	if v := len(b.OldOrderNo); v == 0 || v > 30 || v%2 != 0 {
		return errors.New("oldOrderNo " + FORAMT_ERROR)
	}

	if v := len(b.OldOrderReqNo); v == 0 || v > 30 || v%2 != 0 {
		return errors.New("oldOrderReqNo " + FORAMT_ERROR)
	}

	if v := len(b.RefundReqNo); v == 0 || v > 30 || v%2 != 0 {
		return errors.New("refundReqNo " + FORAMT_ERROR)
	}

	if _, err := time.Parse("20060102", b.RefundReqDate); err != nil {
		return errors.New("refundReqDate " + FORAMT_ERROR)
	}

	if b.TransAmt <= 0 {
		return errors.New("transAmt " + FORAMT_ERROR)
	}

	if v := len(b.LedgerDetail); v > 255 {
		return errors.New("ledgerDetail " + FORAMT_ERROR)
	}

	b.Channel = "05"

	if v := len(b.BgUrl); v > 255 {
		return errors.New("bgUrl " + FORAMT_ERROR)
	}

	//b.Mac 不做校验..这是一个类似签名的东西
	return nil
}

type Resp_bestpay_commonrefund struct {
	OldOrderNo  string `json:"oldOrderNo,omitempty"`      //原扣款成功的订单号 30
	RefundReqNo string `json:"refundReqNo,omitempty"`     //该流水在商户处必须唯一。退款流水 refundReqNo不能和支付流水oldOrderNo 相同。若存在部分退款场景，具体见说明8.5原扣款成功的请求支付流水号
	TransAmt    int    `json:"transAmt,omitempty,string"` //单位为分，小于等于原订单金额
	Sign        string `json:"sign,omitempty"`            //十六进制
}

/**
交易撤单
https://webpaywg.bestpay.com.cn/reverse/reverse
交易撤单
*/
type bestpay_reverse struct {
	BestpayApi
}

func (a *bestpay_reverse) apiMethod() string {
	return BESTPAY_URL_REVERSE
}

func (a *bestpay_reverse) apiName() string {
	return "交易撤单"
}

type Biz_bestpay_reverse struct {
	MerchantId    string `json:"merchantId,omitempty"`      //由翼支付网关平台统一分配 30
	SubMerchantId string `json:"subMerchantId,omitempty"`   //由商户平台自己分配 30
	MerchantPwd   string `json:"merchantPwd,omitempty"`     //商户执行时需填入相应密码 ，又称:交易key
	OldOrderNo    string `json:"oldOrderNo,omitempty"`      //原扣款成功的订单号 30
	OldOrderReqNo string `json:"oldOrderReqNo,omitempty"`   //原扣款成功的请求支付流水号
	RefundReqNo   string `json:"refundReqNo,omitempty"`     //该流水在商户处必须唯一。退款流水 refundReqNo不能和支付流水oldOrderNo 相同。若存在部分退款场景，具体见说明8.5原扣款成功的请求支付流水号
	RefundReqDate string `json:"refundReqDate,omitempty"`   //yyyyMMDD
	TransAmt      int    `json:"transAmt,omitempty,string"` //单位为分，小于等于原订单金额
	Channel       string `json:"channel,omitempty"`         //默认填:05
	Mac           string `json:"mac,omitempty"`             //采用标准的MD5算法，由商户实现， MD5 加密获得32位大写字符 32
}

//mac 校验域.看起来像是一个请求签名的动作
//返回一个待 mac 的数据
func (b Biz_bestpay_reverse) tobe_mac() string {
	tobe_mac := "MERCHANTID=" + b.MerchantId
	tobe_mac += "&MERCHANTPWD=" + b.MerchantPwd
	tobe_mac += "&OLDORDERNO=" + b.OldOrderNo
	tobe_mac += "&OLDORDERREQNO=" + b.OldOrderReqNo
	tobe_mac += "&REFUNDREQNO=" + b.RefundReqNo
	tobe_mac += "&REFUNDREQDATE=" + b.RefundReqDate
	tobe_mac += "&TRANSAMT=" + fmt.Sprintf("%d", b.TransAmt)
	return tobe_mac
}

func (b Biz_bestpay_reverse) valid() error {
	if v := len(b.MerchantId); v == 0 || v > 30 {
		return errors.New("merchantId " + FORAMT_ERROR)
	}

	if v := len(b.SubMerchantId); v > 30 {
		return errors.New("subMerchantId " + FORAMT_ERROR)
	}

	if v := len(b.MerchantPwd); v == 0 || v > 20 {
		return errors.New("merchantPwd " + FORAMT_ERROR)
	}

	if v := len(b.OldOrderNo); v == 0 || v > 30 || v%2 != 0 {
		return errors.New("oldOrderNo " + FORAMT_ERROR)
	}

	if v := len(b.OldOrderReqNo); v == 0 || v > 30 || v%2 != 0 {
		return errors.New("oldOrderReqNo " + FORAMT_ERROR)
	}

	if v := len(b.RefundReqNo); v == 0 || v > 30 || v%2 != 0 {
		return errors.New("refundReqNo " + FORAMT_ERROR)
	}

	if _, err := time.Parse("20060102", b.RefundReqDate); err != nil {
		return errors.New("refundReqDate " + FORAMT_ERROR)
	}

	if b.TransAmt <= 0 {
		return errors.New("transAmt " + FORAMT_ERROR)
	}

	b.Channel = "05"

	//b.Mac 不做校验..这是一个类似签名的东西
	return nil
}

type Resp_bestpay_reverse struct {
	OldOrderNo  string `json:"oldOrderNo,omitempty"`      //原扣款成功的订单号 30
	RefundReqNo string `json:"refundReqNo,omitempty"`     //该流水在商户处必须唯一。退款流水 refundReqNo不能和支付流水oldOrderNo 相同。若存在部分退款场景，具体见说明8.5原扣款成功的请求支付流水号
	TransAmt    int    `json:"transAmt,omitempty,string"` //单位为分，小于等于原订单金额
	Sign        string `json:"sign,omitempty"`            //十六进制
}

func init() {
	registerApi(new(bestpay_barcode_placeorder))
	registerApi(new(bestpay_queryorder))
	registerApi(new(bestpay_commonrefund))
	registerApi(new(bestpay_reverse))
}
