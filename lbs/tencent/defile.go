package tencent

// 通过ip获取地址信息的接口返回数据
type IpResult struct {
	Message string `json:"message"`
	Result  struct {
		AdInfo struct {
			Adcode   int64  `json:"adcode"`
			City     string `json:"city"`
			District string `json:"district"`
			Nation   string `json:"nation"`
			Province string `json:"province"`
		} `json:"ad_info"`
		IP       string `json:"ip"`
		Location struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
	} `json:"result"`
	Status int64 `json:"status"`
}

type LocationResult struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Result    struct {
		AdInfo struct {
			Adcode   string `json:"adcode"`
			City     string `json:"city"`
			CityCode string `json:"city_code"`
			District string `json:"district"`
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			Name       string `json:"name"`
			Nation     string `json:"nation"`
			NationCode string `json:"nation_code"`
			Province   string `json:"province"`
		} `json:"ad_info"`
		Address          string `json:"address"`
		AddressComponent struct {
			City         string `json:"city"`
			District     string `json:"district"`
			Nation       string `json:"nation"`
			Province     string `json:"province"`
			Street       string `json:"street"`
			StreetNumber string `json:"street_number"`
		} `json:"address_component"`
		AddressReference struct {
			BusinessArea struct {
				DirDesc  string  `json:"dir_desc"`
				Distance float64 `json:"distance"`
				ID       string  `json:"id"`
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
				Title string `json:"title"`
			} `json:"business_area"`
			Crossroad struct {
				DirDesc  string  `json:"dir_desc"`
				Distance float64 `json:"distance"`
				ID       string  `json:"id"`
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
				Title string `json:"title"`
			} `json:"crossroad"`
			FamousArea struct {
				DirDesc  string  `json:"dir_desc"`
				Distance float64 `json:"distance"`
				ID       string  `json:"id"`
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
				Title string `json:"title"`
			} `json:"famous_area"`
			LandmarkL1 struct {
				DirDesc  string  `json:"dir_desc"`
				Distance float64 `json:"distance"`
				ID       string  `json:"id"`
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
				Title string `json:"title"`
			} `json:"landmark_l1"`
			LandmarkL2 struct {
				DirDesc  string  `json:"dir_desc"`
				Distance float64 `json:"distance"`
				ID       string  `json:"id"`
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
				Title string `json:"title"`
			} `json:"landmark_l2"`
			Street struct {
				DirDesc  string  `json:"dir_desc"`
				Distance float64 `json:"distance"`
				ID       string  `json:"id"`
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
				Title string `json:"title"`
			} `json:"street"`
			StreetNumber struct {
				DirDesc  string  `json:"dir_desc"`
				Distance float64 `json:"distance"`
				ID       string  `json:"id"`
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
				Title string `json:"title"`
			} `json:"street_number"`
			Town struct {
				DirDesc  string  `json:"dir_desc"`
				Distance float64 `json:"distance"`
				ID       string  `json:"id"`
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
				Title string `json:"title"`
			} `json:"town"`
		} `json:"address_reference"`
		FormattedAddresses struct {
			Recommend string `json:"recommend"`
			Rough     string `json:"rough"`
		} `json:"formatted_addresses"`
		Location struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
	} `json:"result"`
	Status int64 `json:"status"`
}

type RouteResult struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Result    struct {
		Routes []struct {
			Distance    int64     `json:"distance"`
			Duration    int64     `json:"duration"`
			Mode        string    `json:"mode"`
			Polyline    []float64 `json:"polyline"`
			Restriction struct {
				Status int64 `json:"status"`
			} `json:"restriction"`
			Steps []struct {
				AccessorialDesc string  `json:"accessorial_desc"`
				ActDesc         string  `json:"act_desc"`
				DirDesc         string  `json:"dir_desc"`
				Distance        int64   `json:"distance"`
				Instruction     string  `json:"instruction"`
				PolylineIdx     []int64 `json:"polyline_idx"`
				RoadName        string  `json:"road_name"`
			} `json:"steps"`
			Tags     []interface{} `json:"tags"`
			TaxiFare struct {
				Fare int64 `json:"fare"`
			} `json:"taxi_fare"`
			Toll              int64 `json:"toll"`
			TrafficLightCount int64 `json:"traffic_light_count"`
			Waypoints         []struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
				Title string `json:"title"`
			} `json:"waypoints"`
		} `json:"routes"`
	} `json:"result"`
	Status int64 `json:"status"`
}
