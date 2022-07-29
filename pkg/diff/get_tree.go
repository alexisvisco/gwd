package diff

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/merkletrie/noder"
	"github.com/pkg/errors"
)

func getTree(repo *git.Repository, ref string) (noder.Noder, error) {
	getter, err := getTreeGetter(repo, ref)
	if err != nil {
		return nil, err
	}

	tree, err := getter.Tree()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tree")
	}

	if tree != nil {
		return object.NewTreeRootNode(tree), nil
	} else {
		return nil, errors.New("object tree is nil")
	}
}

func getTreeGetter(repo *git.Repository, reference string) (TreeGetter, error) {
	// by branch name
	ref, err := repo.Reference(plumbing.NewBranchReferenceName(reference), false)
	if err == nil {
		encodedObject, err := repo.Storer.EncodedObject(plumbing.CommitObject, ref.Hash())
		if err == nil {
			return object.DecodeCommit(repo.Storer, encodedObject)
		}
	}

	// by tag name
	ref, err = repo.Reference(plumbing.NewTagReferenceName(reference), false)
	if err == nil {
		encodedObject, err := repo.Storer.EncodedObject(plumbing.TagObject, ref.Hash())
		if err == nil {
			return object.DecodeTag(repo.Storer, encodedObject)
		}
	}

	// by hash reference
	encodedObject, err := repo.Storer.EncodedObject(plumbing.CommitObject, plumbing.NewHash(reference))
	if err == nil {
		return object.DecodeCommit(repo.Storer, encodedObject)
	}

	return nil, errors.Wrap(err, "failed to get tree getter: commit, branch or tag reference not found")
}

type TreeGetter interface {
	Tree() (*object.Tree, error)
}
