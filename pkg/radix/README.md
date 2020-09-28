# radix

Simple implementation of radix trees in Go.

A radix tree is a space-optimized/compressed version of a standard [trie](https://en.wikipedia.org/wiki/Trie), with every node that is a single child being merged with their parent. Unlike regular tries, the edges of a radix tree can hold strings, not only single characters.

### References
[Wikipedia article about Radix trees](https://en.wikipedia.org/wiki/Radix_tree)
[Compressing radix trees without too many tears](https://medium.com/basecs/compressing-radix-trees-without-too-many-tears-a2e658adb9a0)

### Related projects
[armon/go-radix](https://github.com/armon/go-radix)
[asergeyev/nradix/](https://github.com/asergeyev/nradix/)
[zmap/go-iptree](https://github.com/zmap/go-iptree)