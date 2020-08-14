package dockerhub

import (
	"context"
	"fmt"
	"time"
)

type Webhook struct {
	ID                  int       `json:"id"`
	Name                string    `json:"name"`
	Active              bool      `json:"active"`
	ExpectFinalCallback bool      `json:"expect_final_callback"`
	Creator             string    `json:"creator"`
	LastUpdated         time.Time `json:"last_updated"`
	LastUpdater         string    `json:"last_updater"`
	Hooks               []Hook    `json:"hooks"`
}

type WebhookInput struct {
	Name string `json:"name"`
}

type Hook struct {
	ID          int       `json:"id"`
	Creator     string    `json:"creator"`
	LastUpdater string    `json:"last_updater"`
	HookURL     string    `json:"hook_url"`
	DateAdded   time.Time `json:"date_added"`
	LastUpdated time.Time `json:"last_updated"`
}

type HookInput struct {
	HookURL string `json:"hook_url"`
}

type WebhookService service

func (s *WebhookService) Find(ctx context.Context, repo string, id int) (*Webhook, *Response, error) {
	path := fmt.Sprintf("v2/repositories/%s/webhooks/%d/", repo, id)
	out := new(Webhook)
	res, err := s.client.Do(ctx, "GET", path, nil, out)
	return out, res, err
}

func (s *WebhookService) FindHook(ctx context.Context, repo string, id, hook int) (*Hook, *Response, error) {
	path := fmt.Sprintf("v2/repositories/%s/webhooks/%d/hooks/%d/", repo, id, hook)
	out := new(Hook)
	res, err := s.client.Do(ctx, "GET", path, nil, out)
	return out, res, err
}

func (s *WebhookService) List(ctx context.Context, repo string) ([]*Webhook, *Response, error) {
	path := fmt.Sprintf("v2/repositories/%s/webhooks/", repo)
	out := new(struct {
		Page
		Results []*Webhook `json:"results"`
	})
	res, err := s.client.Do(ctx, "GET", path, nil, out)
	res.Page = out.Page
	return out.Results, res, err
}

func (s *WebhookService) ListHooks(ctx context.Context, repo string, id int) ([]*Hook, *Response, error) {
	path := fmt.Sprintf("v2/repositories/%s/webhooks/%d/hooks/", repo, id)
	out := new(struct {
		Page
		Results []*Hook `json:"results"`
	})
	res, err := s.client.Do(ctx, "GET", path, nil, out)
	res.Page = out.Page
	return out.Results, res, err
}

func (s *WebhookService) Create(ctx context.Context, repo string, in *WebhookInput) (*Webhook, *Response, error) {
	path := fmt.Sprintf("v2/repositories/%s/webhooks/", repo)
	out := new(Webhook)
	res, err := s.client.Do(ctx, "POST", path, in, out)
	return out, res, err
}

func (s *WebhookService) CreateHook(ctx context.Context, repo string, id int, in *HookInput) (*Hook, *Response, error) {
	path := fmt.Sprintf("v2/repositories/%s/webhooks/%d/hooks/", repo, id)
	out := new(Hook)
	res, err := s.client.Do(ctx, "POST", path, in, out)
	return out, res, err
}

func (s *WebhookService) Update(ctx context.Context, repo string, id int, in *WebhookInput) (*Webhook, *Response, error) {
	path := fmt.Sprintf("v2/repositories/%s/webhooks/%d/", repo, id)
	out := new(Webhook)
	res, err := s.client.Do(ctx, "PATCH", path, in, out)
	return out, res, err
}

func (s *WebhookService) UpdateHook(ctx context.Context, repo string, id, hook int, in *HookInput) (*Hook, *Response, error) {
	path := fmt.Sprintf("v2/repositories/%s/webhooks/%d/hooks/%d/", repo, id, hook)
	out := new(Hook)
	res, err := s.client.Do(ctx, "PATCH", path, in, out)
	return out, res, err
}

func (s *WebhookService) Delete(ctx context.Context, repo string, id int) (*Response, error) {
	path := fmt.Sprintf("v2/repositories/%s/webhooks/%d/", repo, id)
	return s.client.Do(ctx, "DELETE", path, nil, nil)
}

func (s *WebhookService) DeleteHook(ctx context.Context, repo string, id, hook int) (*Response, error) {
	path := fmt.Sprintf("v2/repositories/%s/webhooks/%d/hooks/%d/", repo, id, hook)
	return s.client.Do(ctx, "DELETE", path, nil, nil)
}
