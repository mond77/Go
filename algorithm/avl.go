package algorithm

type Node struct {
	key int
	// 	val int
	l      *Node
	r *Node
	//
	height int
}

type AVLTree struct {
	root *Node
}

func NewAVLTree(nums []int) *AVLTree {
	avl := &AVLTree{}
	for _, v := range nums {
		avl.InsertV(v)
	}
	return avl
}

func (t *AVLTree) InsertV(v int) {
	insert(t.root, v)
	
}

func insert(node *Node, v int) *Node {
	if node == nil {
		return &Node{key: v}
	}
	if v == node.key {

	} else if v < node.key {
		node.l = insert(node.l, v)
	} else {
		node.r = insert(node.r, v)
	}
	return AVLBalance(node)
}

func AVLBalance(node *Node) *Node {
	//balanceFactor 平衡因子
	bf := hdiff(node)

	if bf > 1 {
		if hdiff(node.l) > 0{
			node = AVLTreeLL(node)
		}else{
			node = AVLTreeLR(node)
		}
	}
	if bf < -1{
		if hdiff(node.r) > 0{
			node = AVLTreeRL(node)
		}else{
			node = AVLTreeRR(node)
		}
	}
	//更新旋转后返回的“根节点”的height或者仅updateh
	updateh(node)
	return node
}

func AVLTreeLL(node *Node)*Node{
	return Rrotate(node)
}

func AVLTreeLR(node *Node)*Node{
	node.l = Lrotate(node.l)
	return Rrotate(node)
}

func AVLTreeRL(node *Node)*Node{
	node.r = Rrotate(node.r)
	return Lrotate(node)
}

func AVLTreeRR(node *Node)*Node{
	return Lrotate(node)
}


func Rrotate(node *Node) *Node{
	left := node.l
	node.l = left.r
	left.r = node
	updateh(node)//仅在旋转时更新被置为子节点的node的height
	//updateh(left)
	return left
}

func Lrotate(node *Node) *Node {
	right := node.r
	node.r = right.l
	right.l = node
	updateh(node)
	//updateh(right)
	return right
}




//helper

func hdiff(node *Node) int {
	if node == nil {
		return 0
	}
	return geth(node.l) - geth(node.r)
}

func geth(node *Node) int {
	if node == nil {
		return 0
	}
	return node.height
}

func updateh(node *Node) {
	if node != nil {
		node.height = 1 + max(geth(node.l), geth(node.r))
	}
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
