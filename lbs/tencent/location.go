package tencent

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const LocationAPIAddress = "https://apis.map.qq.com/ws/geocoder/v1"
const IPApiAddress = "https://apis.map.qq.com/ws/location/v1/ip"
const RouteApiAddress = "https://apis.map.qq.com/ws/direction/v1/driving"

type Service struct {
	Key string
}

var service *Service

func LBSInit(key string) *Service {
	service = &Service{Key: key}
	return service
}

// API地址：https://lbs.qq.com/service/webService/webServiceGuide/webServiceGcoder
// params:
//	lat: 必传
//	lon: 必传
func Location(lat, lng string) (LocationResult, error) { return service.Location(lat, lng) }
func (s *Service) Location(lat, lng string) (LocationResult, error) {
	var params = url.Values{}
	var result LocationResult
	params.Add("location", lat+","+lng)
	params.Add("key", s.Key)
	var requestUrl = LocationAPIAddress + "?" + params.Encode()

	resp, err := http.Get(requestUrl)
	if err != nil {
		return result, nil
	}
	defer resp.Body.Close()

	bys, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(bys, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// API地址：https://lbs.qq.com/service/webService/webServiceGuide/webServiceIp
// params:
// 	ip: 可选，
// 	output: 返回格式：支持JSON/JSONP，默认JSON
//	callback: JSONP方式回调函数
func Ip(ip, output, callback string) (IpResult, error) { return service.Ip(ip, output, callback) }
func (s *Service) Ip(ip, output, callback string) (IpResult, error) {
	var result IpResult
	var params = url.Values{}
	if ip != "" {
		params.Add("ip", ip)
	}
	if output != "" {
		params.Add("output", output)
	}
	if callback != "" {
		params.Add("callback", callback)
	}
	params.Add("key", s.Key)

	var requestUrl = IPApiAddress + "?" + params.Encode()
	resp, err := http.Get(requestUrl)
	if err != nil {
		return result, nil
	}
	defer resp.Body.Close()

	bys, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(bys, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func Route(latFrom, lngFrom, latTo, lngTo string) (RouteResult, error) {
	return service.Route(latFrom, lngFrom, latTo, lngTo)
}
func (s *Service) Route(latFrom, lngFrom, latTo, lngTo string) (RouteResult, error) {
	var route RouteResult
	var params = url.Values{}
	params.Add("from", latFrom+","+lngFrom)
	params.Add("to", latTo+","+lngTo)
	params.Add("key", s.Key)

	var requestUrl = RouteApiAddress + "?" + params.Encode()

	resp, err := http.Get(requestUrl)
	if err != nil {
		return route, nil
	}
	defer resp.Body.Close()

	bys, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return route, err
	}
	err = json.Unmarshal(bys, &route)
	if err != nil {
		return route, err
	}
	return route, nil
}
