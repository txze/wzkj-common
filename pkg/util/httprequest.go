package util

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/hzxiao/goutil"
)

func HttpGet(uri string, params goutil.Map) (goutil.Map, error) {
	var values = url.Values{}
	for key := range params {
		values.Add(key, params.GetString(key))
	}

	uri = fmt.Sprintf("%s?%s", uri, values.Encode())

	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result goutil.Map
	err = Bytes2S(body, &result)
	return result, err
}

func HttpPost(url string, data goutil.Map) (goutil.Map, error) {
	var body = bytes.NewBuffer([]byte(S2Json(data)))
	res, err := http.Post(url, "application/json", body)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result goutil.Map
	err = Bytes2S(resBody, &result)
	return result, err
}
