package whois

import (
	"fmt"

	"github.com/zonedb/zonedb"
)

func Fetch(query string) (*Response, error) {
	req, err := NewRequest(query)
	if err != nil {
		return nil, err
	}
	return DefaultClient.Fetch(req)
}

func Server(query string) (string, string, error) {
	z := zonedb.PublicZone(query)
	if z == nil {
		return "", "", fmt.Errorf("no public zone found for %s", query)
	}

	h := z.WhoisServer()
	if h != "" {
		return h, "", nil
	}

	return "", "", fmt.Errorf("no whois server found for %s", query)
}
