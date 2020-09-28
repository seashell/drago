package radix

import "sort"

type node struct {
	leaf  *leaf
	edges []*edge
}

type leaf struct {
	key   string
	value interface{}
}

func (n *node) isLeaf() bool {
	return n.leaf != nil
}

func (n *node) addEdge(e *edge) {
	n.edges = append(n.edges, e)
	n.sortEdges()
}

func (n *node) deleteEdge(s string) {
	// Apply binary search, since our edges sorted
	idx := sort.Search(len(n.edges), func(i int) bool {
		return n.edges[i].label >= s
	})
	n.edges = append(n.edges[:idx], n.edges[idx+1:]...)
}

func (n *node) sortEdges() {
	sort.Slice(n.edges, func(i, j int) bool {
		return n.edges[i].label < n.edges[j].label
	})
}

// getEdgeWithLongestCommonPrefix takes a query string 's' and finds
// amongst all edges in the node the longest prefix they have in
// common with 's', returning both the prefix and the edge.
//
// Example:
//
// query = "abcdef", edges = [{"aaa"},{"abcde"},{"de"},{"bde"}]
//
// output = "ab", {"abcde"}
func (n *node) getEdgeWithLongestCommonPrefix(s string) (string, *edge) {
	for _, e := range n.edges {
		if p := longestCommonPrefix(s, e.label); p != "" {
			return p, e
		}
	}
	return "", nil
}
