package tree 

const(
	RED bool = true
	BLACK bool = false
)

type RBTree struct{
	root *Node
	height int
}

func NewRBTree()*RBTree{
	return &RBTree{}
}

func (rb *RBTree)Put(key,val int){
	rb.root = PutRBNode(rb.root,key,val)
	if rb.root.isred == RED{
		rb.height++
		rb.root.isred = BLACK
	}
}

func PutRBNode(node *Node,key,val int)*Node{
	if node == nil{
		return &Node{key: key,val: val,isred: RED}
	}
	if node.key>key{
		node.l = PutRBNode(node.l,key,val)
	}else if node.key<key{
		node.r = PutRBNode(node.r,key,val)
	}else{
		node.val = val
	}

	if isRed(node.r) && isRed(node.l){
		Lrotate_rb(node)
	}
	if isRed(node.l) && isRed(node.l.l){
		Rrotate_rb(node)
	}
	if isRed(node.l) && isRed(node.r){
		FilpColors(node)
	}
	return node
}

func isRed(node *Node)bool{
	if node == nil{
		return false
	}
	return node.isred == RED
}


func Lrotate_rb(node *Node)*Node{
	r := node.r
	node.r = r.l
	r.l = node
	//左旋的情况，右子节点一定isRed
	r.isred = node.isred
	node.isred = RED
	
	return r
}

func Rrotate_rb(node *Node)*Node{
	l := node.l
	node.l = l.r
	l.r = node
	l.isred = node.isred
	node.isred = RED

	return l
}

func FilpColors(node *Node){
	node.isred = RED
	node.l.isred = BLACK
	node.r.isred = BLACK
}