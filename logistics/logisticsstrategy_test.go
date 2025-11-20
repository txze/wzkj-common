package logistics

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/txze/wzkj-common/logistics/kd100"
	"github.com/txze/wzkj-common/logistics/model"
	"github.com/txze/wzkj-common/logistics/shentong"
)

const (
	AppKey       = ""
	SecretKey    = ""
	ResourceCode = ""
	SourceCode   = ""

	KEY      = ""
	CUSTOMER = ""
	Secret   = ""
)

func TestLogisticsContext_QueryLogisticsByNumber(t *testing.T) {

	t.Run("TestLogisticsContext_QueryLogisticsByNumber", func(t *testing.T) {

		req := model.QueryLogisticsRequest{
			ExpressCode: "shunfeng",
			WaybillNo:   "SF3230939648090",
			Phone:       "13643348872",
		}

		kd := kd100.NewKD100(kd100.KD100Config{
			KEY:      KEY,
			CUSTOMER: CUSTOMER,
			Secret:   Secret,
		})

		c := &LogisticsContext{
			strategy: kd,
		}
		got, err := c.QueryLogisticsByNumber(&req)
		if err != nil {
			t.Errorf("QueryLogisticsByNumber() error = %v", err)
			return
		}
		t.Log(got)
	})

}

func TestLogisticsContext_CreateOrder(t *testing.T) {

	t.Run("TestLogisticsContext_CreateOrder", func(t *testing.T) {
		client := shentong.NewSTOClient(&shentong.Config{
			IsSandbox:    true,
			AppKey:       AppKey,
			SecretKey:    SecretKey,
			ResourceCode: ResourceCode,
			SourceCode:   SourceCode,
		})

		rsp, err := client.CreateOrder(&model.CreateOrderReq{
			OrderNo: strconv.Itoa(int(time.Now().Unix())),
			Sender: model.Sender{
				Name:     "孔",
				Tel:      "",
				Mobile:   "18652161332",
				PostCode: "",
				Country:  "",
				Province: "江苏省",
				City:     "徐州市",
				Area:     "云龙区",
				Town:     "",
				Address:  "生生世世",
			},
			Receiver: model.Receiver{
				Name:     "孔",
				Tel:      "",
				Mobile:   "18652161332",
				PostCode: "",
				Country:  "",
				Province: "江苏省",
				City:     "徐州市",
				Area:     "云龙区",
				Town:     "",
				Address:  "XX街道XX小区XX楼666",
			},
			Cargo: model.Cargo{
				Battery:   "30",
				GoodsType: "小件",
				GoodsName: "生鲜",
				Weight:    20,
			},
		})

		if err != nil {
			t.Errorf("TestLogisticsContext_CreateOrder() error = %v", err)
			return
		}
		t.Log(rsp)
	})

}

func TestLogisticsContext_CancelOrder(t *testing.T) {

	t.Run("TestLogisticsContext_CancelOrder", func(t *testing.T) {
		client := shentong.NewSTOClient(&shentong.Config{
			IsSandbox:    true,
			AppKey:       AppKey,
			SecretKey:    SecretKey,
			ResourceCode: ResourceCode,
			SourceCode:   SourceCode,
		})

		err := client.CancelOrder(&model.CancelOrderReq{
			WaybillNo: "884000392400539",
			Remark:    "不需要了",
		})

		if err != nil {
			t.Errorf("TestLogisticsContext_CancelOrder() error = %v", err)
			return
		}
	})

}

func TestLogisticsContext_QuerySendServiceDetail(t *testing.T) {

	t.Run("TestLogisticsContext_QuerySendServiceDetail", func(t *testing.T) {
		client := shentong.NewSTOClient(&shentong.Config{
			IsSandbox:    true,
			AppKey:       AppKey,
			SecretKey:    SecretKey,
			ResourceCode: ResourceCode,
			SourceCode:   SourceCode,
		})

		rsp, err := client.GetPriceQuote(&model.GetPriceQuoteReq{
			SendName:    "孔先生",
			SendMobile:  "18652161332",
			SendProv:    "江苏省",
			SendCity:    "徐州市",
			SendArea:    "云龙区",
			SendAddress: "万科翡翠之光",
			RecName:     "孔小叔",
			RecMobile:   "18652161332",
			RecProv:     "江苏省",
			RecCity:     "苏州市",
			RecArea:     "吴中区",
			RecAddress:  "尹山湖公园",
			OpenId:      "1111",
			Weight:      "10",
		})

		if err != nil {
			t.Errorf("TestLogisticsContext_CreateOrder() error = %v", err)
			return
		}
		t.Log(rsp)
	})

}

func TestLogisticsContext_ParseAddress(t *testing.T) {
	t.Run("TestLogisticsContext_ParseAddress", func(t *testing.T) {
		kd := kd100.NewKD100(kd100.KD100Config{
			KEY:      KEY,
			CUSTOMER: CUSTOMER,
			Secret:   Secret,
		})

		c := &LogisticsContext{
			strategy: kd,
		}
		got, err := c.ParseAddress("上海市青浦区赵重公路1888号")
		if err != nil {
			t.Errorf("TestLogisticsContext_ParseAddress() error = %v", err)
			return
		}

		data, _ := json.Marshal(got)
		t.Log(string(data))
	})
}

func TestLogisticsContext_SubscribeTracking(t *testing.T) {

	t.Run("TestLogisticsContext_SubscribeTracking", func(t *testing.T) {
		client := shentong.NewSTOClient(&shentong.Config{
			IsSandbox:    true,
			AppKey:       AppKey,
			SecretKey:    SecretKey,
			ResourceCode: ResourceCode,
			SourceCode:   SourceCode,
		})

		err := client.SubscribeTracking(&model.SubscribeTrackingReq{
			ExpressCode: "",
			WaybillNo:   "777031922725111",
			Parameters:  nil,
		})

		if err != nil {
			t.Errorf("TestLogisticsContext_CreateOrder() error = %v", err)
			return
		}
	})

}

func TestLogisticsContext_StoQueryLogisticsByNumber(t *testing.T) {

	t.Run("TestLogisticsContext_StoQueryLogisticsByNumber", func(t *testing.T) {
		client := shentong.NewSTOClient(&shentong.Config{
			IsSandbox:    true,
			AppKey:       AppKey,
			SecretKey:    SecretKey,
			ResourceCode: ResourceCode,
			SourceCode:   SourceCode,
		})

		rsp, err := client.QueryLogistics(&model.QueryLogisticsRequest{
			WaybillNo: "777031922725111",
		})

		if err != nil {
			t.Errorf("TestLogisticsContext_CreateOrder() error = %v", err)
			return
		}
		t.Log(rsp)
	})

}
