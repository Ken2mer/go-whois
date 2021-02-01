package whois

import "fmt"

type Request struct {
	Query string
	Host  string
	URL   string
	Body  []byte
}

func NewRequest(query string) (*Request, error) {
	req := &Request{Query: query}
	if err := req.Prepare(); err != nil {
		return nil, err
	}
	return req, nil
}

func (req *Request) Prepare() error {
	var err error
	if req.Host == "" {
		if req.Host, req.URL, err = Server(req.Query); err != nil {
			return err
		}
	}
	req.Body = []byte(fmt.Sprintf("%s\r\n", req.Query))
	return nil
}
