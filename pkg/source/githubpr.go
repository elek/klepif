package source

import (
	"fmt"
	"github.com/elek/klepif/pkg/client"
	"github.com/elek/klepif/pkg/persistence"
	"github.com/google/go-github/v22/github"
	"github.com/sirupsen/logrus"
	"time"
)

type GithubPrChange struct {
	Pr         *github.PullRequest
	Commits    []*github.RepositoryCommit
	Comments   []*github.IssueComment
	LastCommit *github.RepositoryCommit
}

func createGithubPrChange(pr *github.PullRequest) *GithubPrChange {
	return &GithubPrChange{
		Pr:       pr,
		Commits:  make([]*github.RepositoryCommit, 0),
		Comments: make([]*github.IssueComment, 0),
	}
}

type GithubPr struct {
	Client      *client.GithubClient
	Persistence persistence.Persistence
	Org         string
	Repo        string
}

func (self *GithubPr) GetEventsOfPr(prNumber int) ([]*GithubPrChange, error) {
	results := make([]*GithubPrChange, 0);
	pr, err := self.Client.GetPullRequest(self.Org, self.Repo, prNumber)
	if err != nil {
		return nil, err
	}
	lastTime := time.Unix(0, 0)
	prChange, err := self.GenerateEvents(pr, lastTime)
	if err != nil {
		return results, err
	}

	results = append(results, prChange);
	return results, nil
}

func (self *GithubPr) GetEventsSinceLastCheck() ([]*GithubPrChange, error) {
	results := make([]*GithubPrChange, 0);
	prs, err := self.Client.ListOpenPullRequests(self.Org, self.Repo)
	if err != nil {
		return nil, err
	}

	for _, pr := range prs {

		persistence_key := fmt.Sprintf("github/%s/%s/PR-%d/updated", self.Org, self.Repo, *pr.Number);
		lastTime, err := persistence.ReadTime(self.Persistence, persistence_key)
		if err != nil {
			panic(err)
		}

		if ! lastTime.Before(*pr.UpdatedAt) {
			//not changed
			continue
		}

		prChange, err := self.GenerateEvents(pr, lastTime)
		if err != nil {
			return results, err
		}

		results = append(results, prChange);
		err = persistence.WriteTime(self.Persistence, persistence_key, *pr.UpdatedAt)

		if err != nil {
			panic(err)
		}

	}
	return results, nil

}

func (self *GithubPr) GenerateEvents(pr *github.PullRequest, lastTime time.Time) (*GithubPrChange, error) {
	logrus.Infof("Found PR %d", pr.GetNumber())
	prChange := createGithubPrChange(pr)
	commits, err := self.Client.ListCommitsOfPr(self.Org, self.Repo, pr.GetNumber())
	if err != nil {
		return nil, err
	}
	for _, commit := range commits {
		if prChange.LastCommit == nil ||
			prChange.LastCommit.Commit.Committer.Date.Before(*commit.Commit.Committer.Date) {
			prChange.LastCommit = commit
		}
		if lastTime.Before(commit.GetCommit().GetCommitter().GetDate()) {
			prChange.Commits = append(prChange.Commits, commit)
		}
	}
	comments, err := self.Client.ListCommentsOfPr(self.Org, self.Repo, pr.GetNumber())
	if err != nil {
		return nil, err
	}
	for _, comment := range comments {
		if lastTime.Before(comment.GetUpdatedAt()) {
			prChange.Comments = append(prChange.Comments, comment)
		}
	}
	return prChange, nil
}

func labeled(labels []*github.Label) bool {
	for _, label := range labels {
		if *label.Name == "ozone" {
			return true
		}
	}
	return false

}
