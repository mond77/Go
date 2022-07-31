package main

import (
	"crypto/sha256"
	"fmt"
)

type MerkleTree struct {
	root  *MerkleNode
	hight int
}

func (t *MerkleTree) GetRoot() *MerkleNode {
	return t.root
}

func NewMerkleTree(list []Content) MerkleTree {
	n := len(list)
	if n == 1 {
		return MerkleTree{root: &MerkleNode{sha256.Sum256([]byte(list[0].x)), nil, nil}}
	}
	if n%2 != 0 {
		list = append(list, list[len(list)-1])
		n++
	}
	values := make([][32]byte, n)
	for i := range values {
		values[i] = sha256.Sum256([]byte(list[i].x))
	}
	nodes := make([]*MerkleNode, n)
	for i := range nodes {
		nodes[i] = &MerkleNode{values[i], nil, nil}
	}
	hight := 1
	for n != 1 {
		if n%2 != 0 {
			nodes[n] = nodes[n-1]
			n++
		}
		for i := 0; i < n/2; i++ {
			newHvalue := make([]byte, 64)
			//构建父节点哈希值时简单顺序连接再哈希，在SPV时就需要判断相连时的顺序
			for j := 0; j < 32; j++ {
				newHvalue[j] = nodes[2*i].hvalue[j]
				newHvalue[32+j] = nodes[2*i+1].hvalue[j]
			}
			NewHvalue := sha256.Sum256(newHvalue)
			nodes[i] = &MerkleNode{NewHvalue, nodes[2*i], nodes[2*i+1]}
		}
		n /= 2
		hight++
	}
	return MerkleTree{root: nodes[0], hight: hight}
}

func (t *MerkleTree) getPath(c Content) ([]*MerkleNode, []bool) {
	targetV := sha256.Sum256([]byte(c.x))
	return t.root.getPathTo(targetV)
}

func (t *MerkleTree) SPV(c Content) bool {
	path, onleft := t.getPath(c)
	if len(path) == 0 {
		return false
	}
	verifyingV := sha256.Sum256([]byte(c.x))
	for i := range path {
		newHvalue := make([]byte, 64)
		if onleft[i] {
			for j := 0; j < 32; j++ {
				newHvalue[j] = verifyingV[j]
				newHvalue[32+j] = path[i].hvalue[j]
			}
		} else {
			for j := 0; j < 32; j++ {
				newHvalue[32+j] = verifyingV[j]
				newHvalue[j] = path[i].hvalue[j]
			}
		}
		NewHvalue := sha256.Sum256(newHvalue)
		verifyingV = NewHvalue
	}
	if verifyingV == t.root.hvalue {
		return true
	}
	return false
}

type MerkleNode struct {
	hvalue [32]byte
	left   *MerkleNode
	right  *MerkleNode
}

func (n *MerkleNode) getPathTo(targetV [32]byte) (path []*MerkleNode, onleft []bool) {

	var dfs func(*MerkleNode) *MerkleNode
	dfs = func(n *MerkleNode) *MerkleNode {
		if n == nil {
			return nil
		} else if n.hvalue == targetV {
			return n
		}
		left := dfs(n.left)
		if left != nil && left.hvalue == targetV {
			path = append(path, n.right)
			onleft = append(onleft, true)
			return left
		}
		right := dfs(n.right)
		if right != nil && right.hvalue == targetV {
			path = append(path, n.left)
			onleft = append(onleft, false)
			return right
		}
		return nil
	}
	dfs(n)
	return
}

type Content struct {
	x string
}

func main() {
	t1 := Content{"这是第一个交易"}
	t2 := Content{"这是第二个交易"}
	t3 := Content{"这是第三个交易"}
	tree1 := NewMerkleTree([]Content{t1, t2, t3})
	fmt.Println(tree1.root.hvalue)
	fmt.Println(tree1.SPV(t3))
}
