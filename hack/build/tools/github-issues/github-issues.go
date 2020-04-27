/**
 * Copyright 2020 Appvia Ltd <info@appvia.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
 * Taken from https://github.com/bboreham/github-issues
 * - Changed format slightly
 */

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var iTemplate = template.Must(template.New("issues").Parse(
	"{{if .PullRequestLinks}}{{if .PRMerged}}- {{.Title}} [PR #{{.Number}}](https://github.com/appvia/kore/pull/{{.Number}})\n" +
		"{{end}}" +
		"{{else}}" +
		"- {{.Title}} [#{{.Number}}](https://github.com/appvia/kore/issues/{{.Number}})\n" +
		"{{end}}",
))

func main() {
	var tc *http.Client
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc = oauth2.NewClient(context.TODO(), ts)
	}
	client := github.NewClient(tc)
	ctx := context.Background()

	owner, repo := "appvia", "kore"
	milestone := ""
	if len(os.Args) > 1 {
		milestone = "Release " + os.Args[1]
	}
	milestoneNumber := ""
	listOptions := github.ListOptions{Page: 1}
	for listOptions.Page != 0 {
		milestones, response, err := client.Issues.ListMilestones(ctx, owner, repo, &github.MilestoneListOptions{State: "all", ListOptions: listOptions})
		if err != nil {
			log.Fatal(err)
		}
		for _, m := range milestones {
			if m.Title != nil && *m.Title == milestone && m.Number != nil {
				milestoneNumber = fmt.Sprintf("%d", *m.Number)
				break
			}
		}
		listOptions.Page = response.NextPage
	}
	if milestoneNumber == "" {
		log.Fatal("Unable to find milestone", milestone)
	}

	listOptions.Page = 1
	for listOptions.Page != 0 {
		issues, response, err := client.Issues.ListByRepo(ctx, owner, repo, &github.IssueListByRepoOptions{Milestone: milestoneNumber, State: "all", ListOptions: listOptions})
		if err != nil {
			log.Fatal(err)
		}

		for _, issue := range issues {
			wrapper := struct {
				*github.Issue
				PR       *github.PullRequest
				PRMerged bool
			}{Issue: issue}
			if issue.PullRequestLinks != nil {
				wrapper.PR, _, err = client.PullRequests.Get(ctx, owner, repo, *issue.Number)
				if err != nil {
					log.Fatal(err)
				}
				wrapper.PRMerged = *wrapper.PR.Merged
			}

			if err := iTemplate.Execute(os.Stdout, wrapper); err != nil {
				log.Fatalf("error executing template %s", err)
			}
		}
		listOptions.Page = response.NextPage
	}
}
