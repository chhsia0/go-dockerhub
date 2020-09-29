package main

import (
	"context"
	"flag"
	"fmt"

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

	repository, res, err := c.Repositories.Find(ctx, *repo)
	printResponse("Find", res, err)
	fmt.Printf("%+v\n", repository)

	tag, res, err := c.Repositories.FindTag(ctx, *repo, "latest")
	printResponse("FindTag", res, err)
	fmt.Printf("%+v\n", tag)
}
