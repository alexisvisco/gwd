package diff

import (
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie/noder"
)

func getTree(repo *git.Repository, ref string) (noder.Noder, error) {
	var errToReturn error
	tree, err := getTreeByBranchOrTag(repo, ref)
	if err != nil {
		errToReturn = errors.Wrap(err, "failed to get tree by branch or tag")
		tree, err = getTreeByCommit(repo, ref)
		if err != nil {
			return nil, errors.Wrap(errors.Wrap(errToReturn, err.Error()), "failed to get tree by commit")
		}
	}
	if tree != nil {
		return object.NewTreeRootNode(tree), nil
	} else {
		return nil, errors.New("object tree is nil")
	}
}

func getTreeByBranchOrTag(repo *git.Repository, branchOrTag string) (*object.Tree, error) {
	ref, err := repo.Reference(plumbing.NewBranchReferenceName(branchOrTag), false)
	if err != nil {
		ref, err = repo.Reference(plumbing.NewTagReferenceName(branchOrTag), false)
	}

	if err != nil {
		return nil, err
	}

	encodedObject, err := repo.Storer.EncodedObject(plumbing.CommitObject, ref.Hash())
	if err != nil {
		return nil, err
	}

	commit, err := object.DecodeCommit(repo.Storer, encodedObject)
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	return tree, nil
}

func getTreeByCommit(repo *git.Repository, hash string) (*object.Tree, error) {
	encodedObject, err := repo.Storer.EncodedObject(plumbing.CommitObject, plumbing.NewHash(hash))
	if err != nil {
		return nil, err
	}

	commit, err := object.DecodeCommit(repo.Storer, encodedObject)
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	return tree, nil
}
