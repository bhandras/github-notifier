package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

const GITHUB_TOKEN = ""
const USERNAME = ""

type repo struct {
	owner string
	name  string
}

var repositories []repo = []repo{
	{"lightninglabs", "loop"},
	{"lightningnetwork", "lnd"},
}

func listImportantStuff(ctx context.Context, client *github.Client) {
	notifications, _, err := client.Activity.ListNotifications(ctx, nil)
	if err != nil {
		panic(err)
	}

	var myMentions []*github.Notification

	for _, notificaton := range notifications {
		if *notificaton.Reason == "mention" {
			myMentions = append(myMentions, notificaton)
		}
	}

	var myPRs []*github.PullRequest

	for _, repo := range repositories {
		prs, _, err := client.PullRequests.List(
			ctx, repo.owner, repo.name, nil,
		)
		if err != nil {
			panic(err)
		}

		for _, pr := range prs {
			reviewers, _, err := client.PullRequests.ListReviewers(
				ctx, repo.owner, repo.name, pr.GetNumber(), nil,
			)
			if err != nil {
				panic(err)
			}

			for _, reviewer := range reviewers.Users {
				if *reviewer.Login != USERNAME {
					continue
				}

				myPRs = append(myPRs, pr)
			}
		}
	}

	if len(myMentions) == 0 && len(myPRs) == 0 {
		fmt.Printf("No mentions or reviews!")
	}

	for _, notificaton := range myMentions {
		subject := notificaton.GetSubject()
		fmt.Printf("** Mentioned: '%v' (%v)\n",
			subject.GetTitle(),
			subject.GetURL())
	}

	for _, pr := range myPRs {
		fmt.Printf("++ Review requested: '%v' (%v)\n",
			pr.GetTitle(), pr.GetHTMLURL())
	}

}

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GITHUB_TOKEN},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	for {
		fmt.Printf("\n=== %v - Getting mentions and review requests\n",
			time.Now().Format(time.UnixDate))
		listImportantStuff(ctx, client)
		time.Sleep(time.Minute)
	}
}
