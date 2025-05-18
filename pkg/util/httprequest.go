package util

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func HttpGet(url string) (M, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result M
	err = Bytes2S(body, &result)
	return result, err
}

func HttpPost(url string, data M) (M, error) {
	var body = bytes.NewBuffer([]byte(S2Json(data)))
	res, err := http.Post(url, "application/json", body)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result M
	err = Bytes2S(resBody, &result)
	return result, err
}
