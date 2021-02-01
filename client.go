package whois

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"time"
)

const (
	DefaultTimeout   = 30 * time.Second
	DefaultReadLimit = 1 << 20 // 1 MB
)

type Client struct {
	Timeout time.Duration
}

var DefaultClient = NewClient(DefaultTimeout)

func NewClient(timeout time.Duration) *Client {
	return &Client{Timeout: timeout}
}

func (c *Client) dialContext(ctx context.Context, network, address string) (net.Conn, error) {
	conn, err := defaultDialer.DialContext(ctx, network, address)
	if err != nil {
		return nil, err
	}

	if deadline, ok := ctx.Deadline(); ok {
		err = conn.SetDeadline(deadline)
	}

	return conn, nil
}

var defaultDialer = &net.Dialer{}

func (c *Client) Fetch(req *Request) (*Response, error) {
	return c.FetchContext(context.Background(), req)
}

// FetchContext sends the Request to a whois server.
// If ctx cancels or times out before the request completes, it will return an error.
func (c *Client) FetchContext(ctx context.Context, req *Request) (*Response, error) {
	if c.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.Timeout)
		defer cancel()
	}

	return c.fetchWhois(ctx, req)
}

func (c *Client) fetchWhois(ctx context.Context, req *Request) (*Response, error) {
	if req.Host == "" {
		return nil, fmt.Errorf("no request host for %s", req.Query)
	}

	conn, err := c.dialContext(ctx, "tcp", req.Host+":43")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if _, err = conn.Write(req.Body); err != nil {
		logError(err)
		return nil, err
	}

	res := NewResponse(req.Query, req.Host)
	if res.Body, err = ioutil.ReadAll(io.LimitReader(conn, DefaultReadLimit)); err != nil {
		logError(err)
		return nil, err
	}

	return res, nil
}

func logError(err error) {
	switch t := err.(type) {
	case net.Error:
		fmt.Fprintf(os.Stderr, "net.Error timeout=%t, temp=%t: %s\n", t.Timeout(), t.Temporary(), err.Error())
	default:
		fmt.Fprintf(os.Stderr, "Unknown error %v: %s\n", t, err.Error())
	}
}
