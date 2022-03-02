package lib

import (
	"net/url"
	"strconv"
)

const (
	baseURLv2       = "https://app.terraform.io/api/v2"
	pageSize        = 10
	paramPageSize   = "page[size]"
	paramPageNumber = "page[number]"
)

type TfcUrl struct {
	url.URL
}

func NewTfcUrl(path string) TfcUrl {
	newURL, _ := url.Parse(baseURLv2 + path)
	v := url.Values{}
	newURL.RawQuery = v.Encode()
	tfcUrl := TfcUrl{
		URL: *newURL,
	}
	tfcUrl.SetParam(paramPageSize, strconv.Itoa(pageSize))
	return tfcUrl
}

func (t *TfcUrl) SetParam(name, value string) {
	values := t.Query()
	values.Set(name, value)
	t.RawQuery = values.Encode()
}
