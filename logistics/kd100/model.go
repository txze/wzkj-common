package kd100

type ParseAddress struct {
	Code int `json:"code"`
	Data struct {
		Result []Result `json:"result"`
		TaskId string   `json:"taskId"`
	} `json:"data"`
	Message string `json:"message"`
	Time    int    `json:"time"`
	Success bool   `json:"success"`
}

type Result struct {
	Address string        `json:"address"`
	Content string        `json:"content"`
	Mobile  []interface{} `json:"mobile"`
	Name    string        `json:"name"`
	Xzq     *Xzq          `json:"xzq"`
}

type Xzq struct {
	City       string `json:"city"`
	Code       string `json:"code"`
	District   string `json:"district"`
	Fourth     string `json:"fourth"`
	FullName   string `json:"fullName"`
	Level      int    `json:"level"`
	ParentCode string `json:"parentCode"`
	Province   string `json:"province"`
	SubArea    string `json:"subArea"`
}

func (x *Xzq) GetFourth() string {
	return x.Fourth
}

func (x *Xzq) GetDistrict() string {
	return x.District
}

func (x *Xzq) GetCityCode() string {
	return x.Code
}

func (x *Xzq) GetCityName() string {
	return x.City
}

func (x *Xzq) GetProvince() string {
	return x.Province
}

func (x *Xzq) GetSubArea() string {
	return x.SubArea
}
