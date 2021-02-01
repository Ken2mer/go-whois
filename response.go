package whois

import (
	"bytes"
	"io/ioutil"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

type Response struct {
	Query string
	Host  string

	MediaType string
	Charset   string

	Body []byte
}

func NewResponse(query, host string) *Response {
	return &Response{
		Query:     query,
		Host:      host,
		MediaType: "text/plain",
		Charset:   "utf-8",
	}
}

func (res *Response) String() string {
	enc, _ := charset.Lookup(res.Charset)

	r := transform.NewReader(bytes.NewReader(res.Body), enc.NewDecoder())

	text, err := ioutil.ReadAll(r)
	if err != nil {
		return ""
	}

	return string(text)
}
