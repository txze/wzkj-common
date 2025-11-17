package model

type Address interface {
	// 省份名称
	GetProvince() string

	// 城市名称
	GetCityName() string

	// 区县名称
	GetDistrict() string

	// 乡 镇 街道等四级行政区名称
	GetFourth() string

	// 详细地址
	GetSubArea() string

	// 城市编码
	GetCityCode() string
}
