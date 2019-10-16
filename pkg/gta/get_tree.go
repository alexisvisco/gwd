package gta

import (
	"errors"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie/noder"
)

func GetTree(repo *git.Repository, ref string) (noder.Noder, error) {
	if tree := getTreeByBranchOrTag(repo, ref); tree != nil {
		return object.NewTreeRootNode(tree), nil
	}

	if tree := getTreeByCommit(repo, ref); tree != nil {
		return object.NewTreeRootNode(tree), nil
	}

	return nil, errors.New("reference is not a tag, branch or commit hash")
}

func getTreeByBranchOrTag(repo *git.Repository, branchOrTag string) *object.Tree {
	ref, err := repo.Reference(plumbing.NewBranchReferenceName(branchOrTag), false)
	if err != nil {
		ref, err = repo.Reference(plumbing.NewTagReferenceName(branchOrTag), false)
	}

	if err != nil {
		return nil
	}

	encodedObject, err := repo.Storer.EncodedObject(plumbing.CommitObject, ref.Hash())
	if err != nil {
		return nil
	}

	commit, err := object.DecodeCommit(repo.Storer, encodedObject)
	if err != nil {
		return nil
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil
	}

	return tree
}

func getTreeByCommit(repo *git.Repository, hash string) *object.Tree {
	encodedObject, err := repo.Storer.EncodedObject(plumbing.CommitObject, plumbing.NewHash(hash))
	if err != nil {
		return nil
	}

	commit, err := object.DecodeCommit(repo.Storer, encodedObject)
	if err != nil {
		return nil
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil
	}

	return tree
}
