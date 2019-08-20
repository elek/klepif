package plugins

import "github.com/google/go-github/v22/github"

func HasLabel(pr *github.PullRequest, labeltoFind string) bool {
	for _, label := range pr.Labels {
		if label.GetName() == labeltoFind {
			return true
		}
	}
	return false
}
