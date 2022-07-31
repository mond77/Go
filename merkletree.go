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

type MerkleNode struct {
	hvalue [32]byte
	left   *MerkleNode
	right  *MerkleNode
}

type Content struct {
	x string
}

func main() {
	t1 := Content{"这是第一个交易"}
	t2 := Content{"这是第二个交易"}
	t3 := Content{"这是第三个交易"}
	fmt.Println(NewMerkleTree([]Content{t1, t2, t3}).root.hvalue)
	fmt.Println(NewMerkleTree([]Content{t1, t2, t3}).root.hvalue)
}
