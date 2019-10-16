package packages

import (
	"gopkg.in/src-d/go-git.v4/utils/merkletrie"
)

func FromChanges(changes merkletrie.Changes) Packages {
	packages := make(Packages)
	containModification := PresenceReason{Reason: ReasonModification}

	for _, ch := range changes {
		path := ch.To
		if path.String() == "" {
			path = ch.From
		}

		if path.IsDir() {
			packages[path.String()] = containModification
			continue
		}

		if len(path) == 1 {
			packages["."] = containModification
			continue
		}

		pathDir := path[:len(path)-1]

		if pathDir.IsDir() {
			packages[pathDir.String()] = containModification
		}
	}

	return packages
}
