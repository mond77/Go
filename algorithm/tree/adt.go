package tree


type Node struct {
	key int
	val int
	l      *Node
	r *Node
	//
	height int
	//rbtree
	isred 	bool
}