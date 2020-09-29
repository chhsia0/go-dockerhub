package dockerhub

import (
	"context"
	"fmt"
	"time"
)

type Repository struct {
	User            string    `json:"user"`
	Name            string    `json:"name"`
	Namespace       string    `json:"namespace"`
	RepositoryType  string    `json:"repository_type"`
	Status          int       `json:"status"`
	Description     string    `json:"description"`
	IsPrivate       bool      `json:"is_private"`
	IsAutomated     bool      `json:"is_automated"`
	CanEdit         bool      `json:"can_edit"`
	StarCount       int       `json:"star_count"`
	PullCount       int       `json:"pull_count"`
	LastUpdated     time.Time `json:"last_updated"`
	IsMigrated      bool      `json:"is_migrated"`
	HasStarred      bool      `json:"has_starred"`
	FullDescription *string   `json:"full_description"`
	Affiliation     string    `json:"affiliation"`
	Permissions     struct {
		Read  bool `json:"read"`
		Write bool `json:"write"`
		Admin bool `json:"admin"`
	} `json:"permissions"`
}

type Tag struct {
	Creator int     `json:"creator"`
	ID      int     `json:"id"`
	ImageID *string `json:"image_id"`
	Images  []struct {
		Architecture string  `json:"architecture"`
		Features     string  `json:"features"`
		Variant      *string `json:"variant"`
		Digest       string  `json:"digest"`
		OS           string  `json:"os"`
		OSFeatures   string  `json:"os_features"`
		OSVersion    *string `json:"os_version"`
		Size         int64   `json:"size"`
	} `json:"images"`
	LastUpdated         time.Time `json:"last_updated"`
	LastUpdater         int       `json:"last_updater"`
	LastUpdaterUsername string    `json:"last_updater_username"`
	Name                string    `json:"name"`
	Repository          int       `json:"repository"`
	FullSize            int64     `json:"full_size"`
	V2                  bool      `json:"v2"`
}

type RepositoryService service

func (s *RepositoryService) Find(ctx context.Context, repo string) (*Repository, *Response, error) {
	path := fmt.Sprintf("v2/repositories/%s/", repo)
	out := new(Repository)
	res, err := s.client.Do(ctx, "GET", path, nil, out)
	return out, res, err
}

func (s *RepositoryService) FindTag(ctx context.Context, repo, tag string) (*Tag, *Response, error) {
	path := fmt.Sprintf("v2/repositories/%s/tags/%s/", repo, tag)
	out := new(Tag)
	res, err := s.client.Do(ctx, "GET", path, nil, out)
	return out, res, err
}
