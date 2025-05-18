package tencent_test

import (
	"testing"

	"wzkj-common/lbs/tencent"
	"wzkj-common/pkg/util"
)

func TestLocation(t *testing.T) {
	var key = "BKGBZ-MGYC6-KCVST-E5PDV-327ZZ-DMFO7"
	tencent.LBSInit(key)

	// var ipRes tencent.IpResult
	var locationRes tencent.LocationResult
	var err error

	// ipRes, err = tencent.Ip("", "", "")
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// t.Log(util.S2Json(ipRes))

	// ipRes, err = tencent.Ip("220.181.38.148", "", "")
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// t.Log(util.S2Json(ipRes))

	locationRes, err = tencent.Location("23.69795", "113.06269")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(util.S2Json(locationRes))
}
