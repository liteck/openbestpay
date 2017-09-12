package openbestpay

import (
	"errors"
	"fmt"

	"strings"

	"reflect"

	"github.com/liteck/logs"
	"github.com/liteck/tools"
	"github.com/liteck/tools/httplib"
)

const (
	CAN_NOT_NIL  = "不能为空"
	FORAMT_ERROR = "格式错误"

	// 付款码支付
	BESTPAY_URL_BARCODE_PLACEORDER = "https://webpaywg.bestpay.com.cn/barcode/placeOrder"
	// 交易查询
	BESTPAY_URL_QUERYORDER = "https://webpaywg.bestpay.com.cn/query/queryOrder"
	// 退款
	BESTPAY_URL_COMMONREFUND = "https://webpaywg.bestpay.com.cn/refund/commonRefund"
	// 撤单
	BESTPAY_URL_REVERSE = "https://webpaywg.bestpay.com.cn/reverse/reverse"
)

type bizInterface interface {
	valid() error
	tobe_mac() string
}

type responseInterface interface{}

type Response struct {
	Success   string      `json:"success,omitempty"`
	ErrorCode string      `json:"errorCode,omitempty"`
	ErrorMsg  string      `json:"errorMsg,omitempty"`
	Result    interface{} `json:"result,omitempty"`
}

type ApiHander interface {
	apiMethod() string
	apiName() string
}

var apiRegistry map[string]BestpayApi = map[string]BestpayApi{}

func registerApi(handler ApiHander) {
	apiRegistry[handler.apiMethod()] = BestpayApi{
		apiname:   handler.apiName,
		apimethod: handler.apiMethod,
	}
}

func GetApi(method string) BestpayApi {
	return apiRegistry[method]
}

/**
提供几个公共方法
*/
//分账数据结构
type Ledger map[string]int

func (l Ledger) Set(subMchId string, amount int) error {
	if len(subMchId) == 0 {
		return errors.New("subMchId can not be nil")
	}

	if amount == 0 {
		return errors.New("amount can not be zero")
	}

	l[subMchId] = amount

	return nil
}

func (l Ledger) Get(subMchId string) int {
	return l[subMchId]
}

//设置分账信息.并转化为接口满足的格式
func SetLedgers(total_amt int, legder Ledger) (string, error) {
	/**
	支付时规则
	分账支付规则:分账商户必须是分账支付商户的子商户、
	分账金额必须大于 0 分，最小 分账单位为 1 分、
	单笔参与分账商户个数不能超过 10 个、
	单笔交易参与分账商户只能出现 一次。
	*/

	if n := len(legder); n == 0 {
		return "", errors.New("legder is nil")
	} else if n > 10 {
		return "", errors.New("legder max number is ten")
	}

	//分账商户必须是分账支付商户的子商户、这个暂时无法判断.交给对方去处理.
	//分账金额必须大于 0 分，最小 分账单位为 1 分、
	if total_amt == 0 {
		return "", errors.New("total amount is zero")
	}

	ret := ""
	for subMchId, amount := range legder {
		if amount < 1 {
			return "", errors.New("per legder min amount is 1")
		}
		total_amt -= amount
		ret += fmt.Sprintf("%s:%d|", subMchId, amount)
	}

	if total_amt != 0 {
		return "", errors.New("total amount not equal sum(legder amount)")
	}

	if len(ret) <= 1 {
		return "", errors.New("system error")
	}

	ret = ret[:len(ret)-1]

	return ret, nil
}

/**
BANKID 的字段说明
输入接口响应过来的 bankid.. 返回对应的说明
*/
type BankId struct {
	Category string //类别
	BankId   string //bankId 字段.
	Desc     string //具体说明
}

func GetBankId(bankid string) (bid BankId) {
	switch bankid {
	case "COMPANYACC_3AC":
		bid = BankId{
			Category: "个账余额",
			BankId:   bankid,
			Desc:     "立减优惠",
		}
		break
	case "EPAYTRAVELACC_3AC":
		bid = BankId{
			Category: "个账余额",
			BankId:   bankid,
			Desc:     "翼游账户",
		}
		break
	case "EPAYACC":
		bid = BankId{
			Category: "个账余额",
			BankId:   bankid,
			Desc:     "翼支付账户",
		}
		break
	case "VOUCHER_3AC":
		bid = BankId{
			Category: "个账余额",
			BankId:   bankid,
			Desc:     "营销代金券",
		}
		break
	case "BESTCARDOLD":
		bid = BankId{
			Category: "个账余额",
			BankId:   bankid,
			Desc:     "老翼支付卡",
		}
		break
	case "BESTCARD":
		bid = BankId{
			Category: "个账余额",
			BankId:   bankid,
			Desc:     "新翼支付卡",
		}
		break
	case "EPAYTRAVELCARD_PRE":
		bid = BankId{
			Category: "个账余额",
			BankId:   bankid,
			Desc:     "翼游卡",
		}
		break
	case "EPAYACCWM":
		bid = BankId{
			Category: "个账余额",
			BankId:   bankid,
			Desc:     "翼支付无密账户",
		}
		break
	default:
		bid.BankId = bankid
		//其它的均归类于这里...
		if strings.HasSuffix(bankid, "_RPB2C") ||
			strings.HasSuffix(bankid, "_QB2C") ||
			strings.HasSuffix(bankid, "_Q") {
			bid.Category = "快捷支付"
			bid.Desc = "快捷支付类型"
		} else if strings.HasSuffix(bankid, "_B2B") ||
			strings.HasSuffix(bankid, "B2B") {
			bid.Category = "企业网银"
			bid.Desc = "对公业务类型"
		} else if strings.HasSuffix(bankid, "_B2C") ||
			strings.HasSuffix(bankid, "_C") ||
			strings.HasSuffix(bankid, "_D") {
			bid.Category = "个人网银"
			bid.Desc = "对私业务类型"
		} else {
			bid = BankId{
				Category: "无",
				BankId:   bankid,
				Desc:     "无此信息",
			}
		}
		break
	}

	return bid
}

type BestpayApi struct {
	Key       string //bestpay 针对每个商户申请之后都会有一个秘钥..需要进行配置.
	params    bizInterface
	apiname   func() string
	apimethod func() string
}

func (b *BestpayApi) SetBizContent(biz bizInterface, key string) error {
	if key == "" {
		return errors.New("key is nil")
	}

	b.Key = key

	if err := biz.valid(); err != nil {
		return err
	}

	b.params = biz

	return nil
}

/**
做签名
*/
func (b *BestpayApi) mac() string {
	tobe_mac := b.params.tobe_mac()
	tobe_mac += "KEY=" + b.Key
	logs.Debug(fmt.Sprintf("==[tobe sign]==[%s]", tobe_mac))
	_sign := tools.MD5(tobe_mac)
	return strings.ToUpper(_sign)
}

/**
请求
*/
func (b *BestpayApi) request(m map[string]interface{}) (string, error) {
	url_link := b.apimethod()
	logs.Debug(fmt.Sprintf("==[request params]==[%s]", url_link))
	http_request := httplib.Post(url_link)
	tmp_string := ""
	for k, _ := range m {
		value := fmt.Sprintf("%v", m[k])
		if value != "" {
			http_request.Param(k, value)
			tmp_string = tmp_string + k + "=" + value + "\t"
		}
	}
	logs.Debug(fmt.Sprintf("==[reuest params]==[%s]", tmp_string))
	var string_result string
	if v, err := http_request.String(); err != nil {
		return "", err
	} else {
		string_result = v

	}
	return string_result, nil
}

func (b *BestpayApi) struct_to_map() map[string]interface{} {
	t := reflect.TypeOf(b.params)
	v := reflect.ValueOf(b.params)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		key := t.Field(i).Name
		value := v.Field(i).Interface()
		tag := t.Field(i).Tag.Get("ali")
		if tag != "" {
			if strings.Contains(tag, ",") {
				ps := strings.Split(tag, ",")
				key = ps[0]
			} else {
				key = tag
			}
		}
		data[key] = value
	}
	return data
}

func (b *BestpayApi) Run() error {
	defer logs.Debug("==bestpay api end=====================")
	logs.Debug("==bestpay api start=====================")
	logs.Debug(fmt.Sprintf("==[method]==[%s]:[%s]", b.apiname(), b.apimethod()))

	//做mac签名
	sign := b.mac()
	logs.Debug(fmt.Sprintf("==[sign result]==[%s]", sign))

	//转换下
	m := b.struct_to_map()
	m["mac"] = sign

	result_string := ""
	if v, err := b.request(m); err != nil {
		return err
	} else {
		result_string = v
		logs.Debug(fmt.Sprintf("==[response]==[%s]", result_string))
	}

	return nil
}
