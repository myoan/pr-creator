package main

import (
	"context"
	"flag"
	"log"

	"github.com/BurntSushi/toml"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type TomlConfig struct {
	AccessToken string
	UserName    string
	Repository  string
	BaseBranch  string
}

type PullRequest struct {
	Branch string
	Title  string
	Body   string
}

func (pr PullRequest) Head(userName string) string {
	return userName + ":" + pr.Branch
}

func (pr PullRequest) Create(config TomlConfig) {
	client := githubClient(config)
	head := pr.Head(config.UserName)
	maintainerCanModify := true
	newPR := &github.NewPullRequest{
		Title:               &pr.Title,
		Head:                &head,
		Base:                &config.BaseBranch,
		Body:                &pr.Body,
		MaintainerCanModify: &maintainerCanModify,
	}
	_, _, err := client.PullRequests.Create(context.Background(), config.UserName, config.Repository, newPR)
	if err != nil {
		log.Fatal("PullRequests.Create returned error: %v", err)
	}
}

func githubClient(config TomlConfig) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.AccessToken})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func main() {
	configPath := flag.String("c", "config.toml", "toml config file path")
	branch := flag.String("b", "development", "head branch name")
	title := flag.String("t", "title", "PR title")
	body := flag.String("body", "body", "PR body")
	flag.Parse()
	var cfg TomlConfig
	_, err := toml.DecodeFile(*configPath, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	pr := PullRequest{Branch: *branch, Title: *title, Body: *body}
	pr.Create(cfg)
}
