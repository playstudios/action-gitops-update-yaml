package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
)

func CloneGithub(repository, token string) (*Repository, error) {
	tokenAuth := &http.TokenAuth{
		Token: token,
	}
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:  fmt.Sprintf("https://github.com/%s", repository),
		Auth: tokenAuth,
	})
	if err != nil {
		return nil, fmt.Errorf("error cloning repository %q: %w", repository, err)
	}
	w, err := r.Worktree()
	if err != nil {
		return nil, fmt.Errorf("error getting WorkTree: %w", err)
	}

	return &Repository{
		repo: r,
		tree: w,
		auth: tokenAuth,
		fs:   w.Filesystem,
	}, nil
}
