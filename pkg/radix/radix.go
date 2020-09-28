package radix

import (
	"fmt"
	"strings"
)

// Tree implements a compressed radix tree, a
// space-optimized/compressed version of a standard trie,
// with every node that is a single child being merged
// with their parent. Unlike regular tries, the edges
// of a radix tree can hold strings, not only single
// characters.
type Tree struct {
	root *node
	size int
}

// WalkFcn ...
type WalkFcn func(string, interface{}) bool

// NewTree creates and return a new radix tree
func NewTree() *Tree {
	return &Tree{
		size: 0,
		root: &node{
			edges: []*edge{},
		},
	}
}

type edge struct {
	label string
	node  *node
}

// Set adds or updates a leaf node prefixed by 'k'.
func (t *Tree) Set(k string, v interface{}) error {

	n := t.root
	search := k

	for {

		if len(search) == 0 {
			// If the current node is a leaf, update its value
			if n.isLeaf() {
				n.leaf.value = v
				return nil
			}
			// Otherwise, insert a leaf
			n.leaf, t.size = &leaf{key: k, value: v}, t.size+1
			return nil
		}

		// Find longest common prefix between the search string and all edges.
		if p, e := n.getEdgeWithLongestCommonPrefix(search); p != "" {

			// The common prefix corresponds to the matching
			// edge label, and we simply walk down this edge.
			if p == e.label {
				n, search = e.node, strings.TrimPrefix(search, p)
				continue
			}

			// Otherwise, we found a common prefix shorter than the
			// matching edge label, and the node must be split.
			nn := &node{}
			nn.addEdge(&edge{
				label: strings.TrimPrefix(e.label, p),
				node:  e.node,
			})
			e.node, e.label = nn, p
			n, search = e.node, strings.TrimPrefix(search, p)
		}

		// If no prefix match was found, crete a new leaf node
		// and add an edge to it.
		ln := &node{
			leaf: &leaf{key: k, value: v},
		}

		n.addEdge(&edge{
			label: search,
			node:  ln,
		})

		t.size++
		break
	}

	return nil
}

// Get searches the tree and returns the value for the
// leaf node whose prefix exactly matches the key passed
// as argument. Additionaly, a bool is returned to indicate
// if a match was found or not.
func (t *Tree) Get(k string) (interface{}, bool) {

	n := t.root
	search := k

	for {

		if len(search) == 0 {
			if n.isLeaf() {
				return n.leaf.value, true
			}
			return nil, false
		}

		if p, e := n.getEdgeWithLongestCommonPrefix(search); p != "" {
			if p == e.label {
				n, search = e.node, strings.TrimPrefix(search, p)
				continue
			}
		}

		break
	}

	return nil, false
}

// GetClosest searches the tree and returns the value for the
// leaf node corresponding to the longest prefix match to the
// key passed as argument. Besides the value itself, the longest
// prefix and a bool indicating whether or not a match was found
// are also returned.
func (t *Tree) GetClosest(k string) (string, interface{}, bool) {

	n := t.root
	search := k

	prefix := ""
	var prevLeaf *leaf

	for {

		if len(search) == 0 {
			return prefix, prevLeaf.value, true
		}

		if n.isLeaf() {
			prevLeaf = n.leaf
		}

		if p, e := n.getEdgeWithLongestCommonPrefix(search); p != "" {
			if p == e.label {
				prefix += e.label
				n, search = e.node, strings.TrimPrefix(search, p)
				continue
			}
		}

		break
	}

	if prevLeaf != nil {
		return prefix, prevLeaf.value, true
	}

	return "", nil, false
}

// Delete removes the leaf node prefixed by 'k' from the tree.
func (t *Tree) Delete(k string) error {

	n := t.root
	search := k

	var prevEdge *edge

	for {

		if p, e := n.getEdgeWithLongestCommonPrefix(search); p != "" {

			if p == e.label {

				// Check if the next node is our target prior
				// to walking down to it.
				if len(strings.TrimPrefix(search, p)) == 0 {

					if e.node.isLeaf() {

						e.node.leaf, t.size = nil, t.size-1

						// If the next node has no edges, remove it.
						if len(e.node.edges) == 0 {
							n.deleteEdge(e.label)
						}

						// If the next node has a single edge, it is not needed,
						// and we can merge the edges.
						if len(e.node.edges) == 1 {
							e.label = e.label + e.node.edges[0].label
							e.node = e.node.edges[0].node
						}

						// If the current node is neither a leaf node nor the
						// root, and contains a single edge, it can too be merged
						// with its parent.
						if !n.isLeaf() && n != t.root && len(n.edges) == 1 {
							prevEdge.label = prevEdge.label + n.edges[0].label
							prevEdge.node = n.edges[0].node
						}

					} else {
						return nil
					}
				}
				n, prevEdge, search = e.node, e, strings.TrimPrefix(search, p)
				continue
			}
		}

		break
	}
	return nil
}

// Walk ...
func (t *Tree) Walk(f WalkFcn) {
	walkNode(t.root, f)
}

func walkNode(n *node, f WalkFcn) bool {

	if n.isLeaf() && f(n.leaf.key, n.leaf.value) {
		return true
	}

	for _, e := range n.edges {
		if walkNode(e.node, f) {
			return true
		}
	}

	return false
}

// Size returns the number of leaves in the radix tree.
func (t *Tree) Size() int {
	return t.size
}

// String returns a string representation of the tree
// which is useful for debugging.
func (t *Tree) String() string {
	return treeString(t.root, 0)
}

// TODO: make padding adaptive, and based on the prefix length.
func treeString(n *node, lvl int) string {

	s := ""
	p := 14

	if len(n.edges) == 0 {
		return s
	}

	for _, e := range n.edges {
		padding := strings.Repeat(" ", lvl*p)
		s += fmt.Sprintf("\n%s|", padding)
		s += fmt.Sprintf("\n%s|----%s----(%v)", padding, e.label, e.node.leaf)
		s += treeString(e.node, lvl+1)
	}

	return s
}

// longestCommonPrefix finds the longest common prefix of the input strings.
// It compares by bytes instead of runes (Unicode code points).
// It's up to the caller to do Unicode normalization if desired
// (e.g. see golang.org/x/text/unicode/norm).
// Credits: https://rosettacode.org/wiki/Longest_common_prefix
func longestCommonPrefix(l ...string) string {
	switch len(l) {
	case 0:
		return ""
	case 1:
		return l[0]
	}
	min, max := l[0], l[0]
	for _, s := range l[1:] {
		switch {
		case s < min:
			min = s
		case s > max:
			max = s
		}
	}
	for i := 0; i < len(min) && i < len(max); i++ {
		if min[i] != max[i] {
			return min[:i]
		}
	}
	return min
}
