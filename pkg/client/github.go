package client

import (
	"context"
	"github.com/elek/klepif/pkg"
	"github.com/google/go-github/v22/github"
	"golang.org/x/oauth2"
	"strings"
)

type GithubClient struct {
	ctx          context.Context
	githubClient *github.Client
	Org          string
	Repo         string
}

func CreateGithubClient(config *pkg.GithubConfig) GithubClient {
	ctx := context.Background();
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return GithubClient{
		githubClient: client,
		ctx:          ctx,
		Org:          config.Org,
		Repo:         config.Repo,
	}
}

func (gh *GithubClient) ListCommitsOfPr(org string, repo string, number int) ([]*github.RepositoryCommit, error) {
	commits, _, err := gh.githubClient.PullRequests.ListCommits(gh.ctx, org, repo, number, &github.ListOptions{});
	return commits, err
}

func (gh *GithubClient) ListCommentsOfPr(org string, repo string, number int) ([]*github.IssueComment, error) {
	comments, _, err := gh.githubClient.Issues.ListComments(gh.ctx, org, repo, number, &github.IssueListCommentsOptions{
		Sort:      "created",
		Direction: "desc",
	});
	return comments, err
}

func (gh *GithubClient) GetPullRequest(org string, repo string, number int) (*github.PullRequest, error) {
	pr, _, err := gh.githubClient.PullRequests.Get(gh.ctx, org, repo, number)
	return pr, err
}

func (gh *GithubClient) ListOpenPullRequests(org string, repo string) ([]*github.PullRequest, error) {
	listResponse, _, err := gh.githubClient.PullRequests.List(gh.ctx, org, repo, &github.PullRequestListOptions{
		State:     "open",
		Sort:      "updated",
		Direction: "desc",
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	})
	return listResponse, err
}

func (gh *GithubClient) AddLabel(id int, s string) error {
	_, _, err := gh.githubClient.Issues.AddLabelsToIssue(gh.ctx, gh.Org, gh.Repo, id, []string{s})
	return err
}

func GetCommand(command string, comments []*github.IssueComment) (bool, string) {
	for _, comment := range comments {
		message := strings.TrimSpace(*comment.Body)
		if strings.HasPrefix(message, command) {
			label := strings.TrimSpace(strings.TrimPrefix(message, command))
			return true, label
		}
	}
	return false, ""
}
