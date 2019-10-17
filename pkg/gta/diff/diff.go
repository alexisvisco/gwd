package diff

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/format/gitignore"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie/filesystem"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie/noder"

	"github.com/alexisvisco/gta/pkg/gta/packages"
)

func LocalDiff(repo *git.Repository, previous noder.Noder) (packages.Packages, error) {
	wt, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	submodules, err := getSubmodulesStatus(wt)
	if err != nil {
		return nil, err
	}
	current := filesystem.NewRootNode(wt.Filesystem, submodules)

	return Diff(repo, previous, current)
}

func getSubmodulesStatus(w *git.Worktree) (map[string]plumbing.Hash, error) {
	o := map[string]plumbing.Hash{}

	sub, err := w.Submodules()
	if err != nil {
		return nil, err
	}

	status, err := sub.Status()
	if err != nil {
		return nil, err
	}

	for _, s := range status {
		if s.Current.IsZero() {
			o[s.Path] = s.Expected
			continue
		}

		o[s.Path] = s.Current
	}

	return o, nil
}

func Diff(repo *git.Repository, previous noder.Noder, current noder.Noder) (packages.Packages, error) {
	wt, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	changes, err := merkletrie.DiffTree(previous, current, diffTreeIsEquals)
	if err != nil {
		return nil, err
	}

	return packages.FromChanges(excludeIgnoredChanges(wt, changes), wt), nil
}

func excludeIgnoredChanges(w *git.Worktree, changes merkletrie.Changes) merkletrie.Changes {
	patterns, err := gitignore.ReadPatterns(w.Filesystem, nil)
	if err != nil {
		return changes
	}

	patterns = append(patterns, w.Excludes...)

	if len(patterns) == 0 {
		return changes
	}

	m := gitignore.NewMatcher(patterns)

	var res merkletrie.Changes
	for _, ch := range changes {
		var path []string
		for _, n := range ch.To {
			path = append(path, n.Name())
		}
		if len(path) == 0 {
			for _, n := range ch.From {
				path = append(path, n.Name())
			}
		}
		if len(path) != 0 {
			isDir := (len(ch.To) > 0 && ch.To.IsDir()) || (len(ch.From) > 0 && ch.From.IsDir())
			if m.Match(path, isDir) {
				continue
			}
		}
		res = append(res, ch)
	}
	return res
}