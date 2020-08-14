package dockerhub

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type Response struct {
	*http.Response
	Page Page
}

type Page struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
}

type Error struct {
	json.RawMessage
}

func (e Error) Error() string {
	return string(e.RawMessage)
}

type Client struct {
	mu           sync.Mutex
	Client       *http.Client
	BaseURL      *url.URL
	Auth         *AuthService
	Webhooks     *WebhookService
	DumpResponse func(*http.Response, bool) ([]byte, error)
}

type service struct {
	client *Client
}

func NewClient(uri string) (*Client, error) {
	base, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %w", uri, err)
	}
	if !strings.HasSuffix(base.Path, "/") {
		base.Path += "/"
	}

	s := &service{new(Client)}
	s.client.BaseURL = base
	s.client.Auth = (*AuthService)(s)
	s.client.Webhooks = (*WebhookService)(s)

	return s.client, nil
}

func NewDefaultClient() *Client {
	client, _ := NewClient("https://hub.docker.com/")
	return client
}

func (c *Client) Do(ctx context.Context, method, path string, in, out interface{}) (*Response, error) {
	uri, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %w", path, err)
	}

	var body io.ReadWriter
	if in != nil {
		body = &bytes.Buffer{}
		enc := json.NewEncoder(body)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(in); err != nil {
			return nil, fmt.Errorf("encoding request body: %w", err)
		}
	}

	req, err := http.NewRequest(method, uri.String(), body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req = req.WithContext(ctx)

	client := c.Client
	if client == nil {
		client = http.DefaultClient
	}

	r, err := client.Do(req)
	res := &Response{Response: r}
	if err != nil {
		return res, fmt.Errorf("sending request: %w", err)
	}

	if c.DumpResponse != nil {
		if _, err := c.DumpResponse(res.Response, true); err != nil {
			return res, fmt.Errorf("dumping response: %w", err)
		}
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		err := new(Error)
		json.NewDecoder(res.Body).Decode(err)
		return res, fmt.Errorf("%s: %w", res.Status, err)
	}

	if out != nil {
		if err := json.NewDecoder(res.Body).Decode(out); err != nil && !errors.Is(err, io.EOF) {
			return res, fmt.Errorf("decoding response body: %w", err)
		}
	}

	return res, nil
}
