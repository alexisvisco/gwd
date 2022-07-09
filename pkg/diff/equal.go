package diff

import (
	"bytes"

	"gopkg.in/src-d/go-git.v4/utils/merkletrie/noder"
)

var emptyNoderHash = make([]byte, 24)

func diffTreeIsEquals(a, b noder.Hasher) bool {
	hashA := a.Hash()
	hashB := b.Hash()

	if bytes.Equal(hashA, emptyNoderHash) || bytes.Equal(hashB, emptyNoderHash) {
		return false
	}

	return bytes.Equal(hashA, hashB)
}
