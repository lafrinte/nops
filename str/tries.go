package str

import (
	"github.com/yihleego/trie"
)

func New(keyword ...string) *trie.Trie {
	return trie.New(keyword...)
}
