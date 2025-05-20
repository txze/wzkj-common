package util_test

import (
	"fmt"
	"testing"

	"github.com/txze/wzkj-common/pkg/util"
)

func TestTime(t *testing.T) {
	fmt.Println(util.DayStart())
	fmt.Println(util.DayEnd())
	fmt.Println(util.WeekStart())
	fmt.Println(util.WeekEnd())
	fmt.Println(util.MonthStart())
	fmt.Println(util.MonthEnd())
}
