package dockerhub

import (
	"context"
	"fmt"
	"net/http"
)

type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JWTAuth struct {
	Token string `json:"token"`
}

type AuthTransport interface {
	Base() http.RoundTripper
}

type JWTAuthTransport struct {
	base http.RoundTripper
	JWTAuth
}

func (t *JWTAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header.Get("Authorization") != "" {
		return t.Base().RoundTrip(req)
	}
	// Shallow copy of the struct.
	r := *req
	// Deep copy of the Header.
	r.Header = req.Header.Clone()
	r.Header.Set("Authorization", fmt.Sprintf("JWT %s", t.Token))
	base := t.Base()
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(&r)
}

func (t *JWTAuthTransport) Base() http.RoundTripper {
	return t.base
}

type AuthService service

func (s *AuthService) Login(ctx context.Context, in BasicAuth) (*Response, error) {
	s.client.mu.Lock()
	defer s.client.mu.Unlock()

	path := "v2/users/login/"
	out := JWTAuth{}
	res, err := s.client.Do(ctx, "POST", path, in, &out)
	if err != nil {
		return res, err
	}
	if s.client.Client == nil {
		s.client.Client = &http.Client{}
	}
	base := s.client.Client.Transport
	if t, ok := base.(AuthTransport); ok {
		base = t.Base()
	}
	s.client.Client.Transport = &JWTAuthTransport{
		base:    base,
		JWTAuth: out,
	}
	return res, nil
}
