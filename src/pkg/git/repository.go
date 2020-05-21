package git

import (
	"context"
	"fmt"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"io"
	"time"
)

type Repository struct {
	repo *git.Repository
	tree *git.Worktree
	auth transport.AuthMethod
	fs   billy.Filesystem
}

func (r *Repository) Checkout(reference string) error {
	ref, err := r.repo.Reference(plumbing.ReferenceName(reference), true)
	if err != nil {
		return fmt.Errorf("error resolving reference: %w", err)
	}

	opts := &git.CheckoutOptions{}
	if ref.Name().IsBranch() {
		opts.Branch = ref.Name()
	} else {
		opts.Hash = ref.Hash()
	}
	if err := r.tree.Checkout(opts); err != nil {
		return fmt.Errorf("error checking out %s: %w", ref, err)
	}

	return nil
}

func (r *Repository) ReadFile(name string) (io.Reader, error) {
	f, err := r.tree.Filesystem.Open(name)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", name, err)
	}
	return f, nil
}

func (r *Repository) WriteFile(name string, b []byte) error {
	f, err := fileOpenOrCreate(r.fs, name)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", name, err)
	}
	if _, err := f.Write(b); err != nil {
		return fmt.Errorf("error writing to file %s: %w", name, err)
	}
	if _, err := r.tree.Add(name); err != nil {
		return fmt.Errorf("error staging file %s: %w", name, err)
	}
	return nil
}

func (r *Repository) Commit(msg string) error {
	if _, err := r.tree.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "",
			Email: "",
			When:  time.Now(),
		},
	}); err != nil {
		return fmt.Errorf("error creating commit: %w", err)
	}
	return nil
}

func (r *Repository) Push() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return r.repo.PushContext(ctx, &git.PushOptions{
		Auth: r.auth,
	})
}
