package trie

import (
	"sort"
	"strings"
)

// Single node of the Trie
type TrieNode struct {
	children map[rune]*TrieNode
	isWord   bool
}

// Trie data structure
type Trie struct {
	_root *TrieNode
}

// Create a new empty Trie
func NewTrie() *Trie {
	return &Trie{_root: &TrieNode{children: map[rune]*TrieNode{}}}
}

// Insert a word
func (t *Trie) Insert(word string) {
	if word == "" {
		return
	}

	node := t._root
	for _, r := range word {
		if node.children[r] == nil {
			node.children[r] = &TrieNode{children: map[rune]*TrieNode{}}
		}
		node = node.children[r]
	}
	node.isWord = true
}

// Insert several words
func (t *Trie) InsertMany(words []string) {
	for _, w := range words {
		t.Insert(w)
	}
}

// Shows if the word already exist
func (t *Trie) Search(word string) bool {
	node := t._root
	for _, r := range word {
		if node.children[r] == nil {
			return false
		}
		node = node.children[r]
	}
	return node.isWord
}

// Shows if there is any word that starts with the given prefix
func (t *Trie) StartsWith(prefix string) bool {
	node := t._root
	for _, r := range prefix {
		if node.children[r] == nil {
			return false
		}
		node = node.children[r]
	}
	return true
}

// Do DFS to pick words from a node
func (t *Trie) Collect(node *TrieNode, prefix []rune, out *[]string, limit int) {
	if node == nil || len(*out) >= limit {
		return
	}

	if node.isWord {
		*out = append(*out, string(prefix))
		if len(*out) >= limit {
			return
		}
	}

	// Iterate children in order to get deterministics results
	keys := make([]rune, 0, len(node.children))
	for k := range node.children {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	for _, r := range keys {
		child := node.children[r]
		t.Collect(child, append(prefix, r), out, limit)
		if len(*out) >= limit {
			return
		}
	}
}

// Suggest returns up to 'limit' suggestions (completions) for the given prefix
// Returns all the complete words (prefix + sufix)
func (t *Trie) Suggest(prefix string, limit int) []string {
	res := make([]string, 0, limit)
	node := t._root
	prefixRunes := []rune(prefix)
	for _, r := range prefixRunes {
		if node.children[r] == nil {
			return res
		}
		node = node.children[r]
	}
	t.Collect(node, prefixRunes, &res, limit)
	return res
}

// Words returns all words stored in the trie
func (t *Trie) Words() []string {
	out := []string{}
	t.Collect(t._root, []rune{}, &out, int(^uint(0)>>1)) // large limit
	return out
}

// BuildTrieFromText helps to create an unique vocabulary (split between spaces) and fill the trie
func BuildTrieFromText(text string, normalizeFunc func(string) string) *Trie {
	t := NewTrie()
	words := strings.Fields(text)
	uniq := map[string]struct{}{}
	for _, w := range words {
		n := strings.TrimSpace(w)
		if normalizeFunc != nil {
			n = normalizeFunc(n)
		}
		if n == "" {
			continue
		}
		uniq[n] = struct{}{}
	}
	list := make([]string, 0, len(uniq))
	for w := range uniq {
		list = append(list, w)
	}
	return t
}
