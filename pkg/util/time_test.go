package util_test

import (
	"fmt"
	"testing"
	"time"

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

func TestHumanDurationZH(t *testing.T) {
	fmt.Println(util.HumanDurationZH(5 * time.Second))
	fmt.Println(util.HumanDurationZH(5 * time.Minute))
	fmt.Println(util.HumanDurationZH(5*time.Minute + 5*time.Second))
	fmt.Println(util.HumanDurationZH(5 * time.Hour))
	fmt.Println(util.HumanDurationZH(5*time.Hour + time.Minute))
	fmt.Println(util.HumanDurationZH(5*time.Hour + time.Minute + 5*time.Second))
}
