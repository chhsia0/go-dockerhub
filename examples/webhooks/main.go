package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"

	"github.com/chhsia0/go-dockerhub/dockerhub"
)

var (
	repo     = flag.String("repository", "", "Docker Hub image repository")
	username = flag.String("username", "", "Docker Hub username")
	password = flag.String("password", "", "Docker Hub password")
)

const (
	alphanums = "bcdfghjklmnpqrstvwxz2456789"
)

func randSuffix(s string, n int) string {
	for i := 0; i < n; i++ {
		s += string(alphanums[rand.Intn(len(alphanums))])
	}
	return s
}

func printResponse(prefix string, res *dockerhub.Response, err error) {
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n%s: %s\n", prefix, res.Status)
}

func main() {
	flag.Parse()

	ctx := context.Background()
	c := dockerhub.NewDefaultClient()

	res, err := c.Auth.Login(ctx, dockerhub.BasicAuth{*username, *password})
	printResponse("Login", res, err)
	fmt.Println(c.Client.Transport.(*dockerhub.JWTAuthTransport).Token)

	created, res, err := c.Webhooks.Create(ctx, *repo, &dockerhub.WebhookInput{randSuffix("foo-", 5)})
	printResponse("Create", res, err)

	webhooks, res, err := c.Webhooks.List(ctx, *repo)
	printResponse("List", res, err)
	for _, webhook := range webhooks {
		fmt.Printf("%+v\n", webhook)
	}

	createdHook, res, err := c.Webhooks.CreateHook(ctx, *repo, created.ID, &dockerhub.HookInput{randSuffix("http://example.com/", 5)})
	printResponse("CreateHook", res, err)

	hooks, res, err := c.Webhooks.ListHooks(ctx, *repo, created.ID)
	printResponse("ListHooks", res, err)
	for _, hook := range hooks {
		fmt.Printf("%+v\n", hook)
	}

	_, res, err = c.Webhooks.Update(ctx, *repo, created.ID, &dockerhub.WebhookInput{randSuffix("bar-", 5)})
	printResponse("Update", res, err)

	created, res, err = c.Webhooks.Find(ctx, *repo, created.ID)
	printResponse("Find", res, err)
	fmt.Printf("%+v\n", created)

	_, res, err = c.Webhooks.UpdateHook(ctx, *repo, created.ID, createdHook.ID, &dockerhub.HookInput{randSuffix("http://example.com/", 5)})
	printResponse("UpdateHook", res, err)

	createdHook, res, err = c.Webhooks.FindHook(ctx, *repo, created.ID, createdHook.ID)
	printResponse("FindHook", res, err)
	fmt.Printf("%+v\n", createdHook)

	res, err = c.Webhooks.DeleteHook(ctx, *repo, created.ID, createdHook.ID)
	printResponse("DeleteHook", res, err)

	res, err = c.Webhooks.Delete(ctx, *repo, created.ID)
	printResponse("Delete", res, err)
}
