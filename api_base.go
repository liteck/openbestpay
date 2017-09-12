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
