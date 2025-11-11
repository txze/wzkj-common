package util

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"strings"

	"github.com/hzxiao/goutil"

	"github.com/txze/wzkj-common/pkg/ierr"
)

// ResponseParser 响应解析器
type ResponseParser struct {
	Content     []byte
	ContentType string
}

// NewResponseParser 创建响应解析器
func NewResponseParser(res *http.Response) (*ResponseParser, error) {
	// 读取响应体
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, ierr.NewIErrorf(ierr.ReadFileError, "读取响应体失败: %w", err)
	}

	// 获取Content-Type
	contentType := res.Header.Get("Content-Type")

	return &ResponseParser{
		Content:     resBody,
		ContentType: contentType,
	}, nil
}

// Parse 解析响应数据
func (rp *ResponseParser) Parse() (goutil.Map, error) {
	// 根据Content-Type选择解析方式
	if rp.isJSON() {
		return rp.parseJSON()
	} else if rp.isXML() {
		return rp.parseXML()
	} else {
		// 自动检测格式
		return rp.autoDetectAndParse()
	}
}

// isJSON 检查是否为JSON格式
func (rp *ResponseParser) isJSON() bool {
	// 根据Content-Type判断
	if strings.Contains(rp.ContentType, "application/json") ||
		strings.Contains(rp.ContentType, "text/json") {
		return true
	}

	// 如果没有Content-Type，尝试解析内容
	if rp.ContentType == "" {
		return json.Valid(rp.Content)
	}

	return false
}

// isXML 检查是否为XML格式
func (rp *ResponseParser) isXML() bool {
	// 根据Content-Type判断
	if strings.Contains(rp.ContentType, "application/xml") ||
		strings.Contains(rp.ContentType, "text/xml") {
		return true
	}

	// 如果没有Content-Type，尝试解析内容
	if rp.ContentType == "" {
		return rp.isValidXML()
	}

	return false
}

// isValidXML 验证是否为有效的XML
func (rp *ResponseParser) isValidXML() bool {
	return xml.Unmarshal(rp.Content, new(interface{})) == nil
}

// parseJSON 解析JSON数据
func (rp *ResponseParser) parseJSON() (goutil.Map, error) {
	var result goutil.Map
	if err := json.Unmarshal(rp.Content, &result); err != nil {
		return nil, ierr.NewIErrorf(ierr.ParseDataFail, "JSON解析失败: %w", err)
	}
	return result, nil
}

type XmlError struct {
	Success   string `xml:"success"`
	ErrorCode string `xml:"errorCode"`
	ErrorMsg  string `xml:"errorMsg"`
}

// parseXML 解析XML数据
func (rp *ResponseParser) parseXML() (goutil.Map, error) {
	var result XmlError
	if err := xml.Unmarshal(rp.Content, &result); err != nil {
		return nil, ierr.NewIErrorf(ierr.ParseDataFail, "XML解析失败: %w", err)
	}
	return goutil.Map{
		"success":   result.Success,
		"errorCode": result.ErrorCode,
		"errorMsg":  result.ErrorMsg,
	}, nil
}

// autoDetectAndParse 自动检测并解析
func (rp *ResponseParser) autoDetectAndParse() (goutil.Map, error) {
	// 首先尝试JSON
	if json.Valid(rp.Content) {
		return rp.parseJSON()
	}

	// 然后尝试XML
	if rp.isValidXML() {
		return rp.parseXML()
	}

	// 如果都不是，返回原始文本
	return nil, ierr.NewIError(ierr.ParseDataFail, "数据解析失败"+string(rp.Content))
}

// GetRawContent 获取原始内容
func (rp *ResponseParser) GetRawContent() string {
	return string(rp.Content)
}

// GetContentType 获取内容类型
func (rp *ResponseParser) GetContentType() string {
	return rp.ContentType
}
