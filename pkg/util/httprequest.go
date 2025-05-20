package util

import (
	"bytes"
	"io"
	"net/http"

	"github.com/txze/goutil"
)

func HttpGet(url string) (goutil.Map, error) {
	res, err := http.Get(url)
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
