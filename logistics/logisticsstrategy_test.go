package logistics

import (
	"testing"

	"github.com/txze/wzkj-common/logistics/kd100"
)

func TestLogisticsContext_QueryLogisticsByNumber(t *testing.T) {

	t.Run("TestLogisticsContext_QueryLogisticsByNumber", func(t *testing.T) {
		kd := kd100.NewKD100(kd100.KD100Config{
			KEY:      "xUwyaYSR8331",
			CUSTOMER: "280D13614D729A608991196DB24BB6C3",
		})

		c := &LogisticsContext{
			strategy: kd,
		}
		got, err := c.QueryLogisticsByNumber("yuantong", "xxxxxx")
		if err != nil {
			t.Errorf("QueryLogisticsByNumber() error = %v", err)
			return
		}
		t.Log(got)
	})
}
