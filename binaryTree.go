package main

import (
	"fmt"
	"sync"
)

//Tree root of entire Tree
type Tree struct {
	root *HashTagTree
	sync.RWMutex
}

//HashTagTree struct is node type for Binary Tree
type HashTagTree struct {
	tag    string
	freq   int
	total  int
	depth  int
	parent *HashTagTree
	left   *HashTagTree
	right  *HashTagTree
}

//AddToTree adds first and second most popular hashtags (corresponding nodes to BinaryTree)
func AddToTree(bestTag string, bestFreq int, secTag string, secFreq int, total int, parent string) {
	parentNode := Find(parent, BTree.root)
	if parentNode == nil {
		fmt.Println("ERROR:", parent)
	}
	BTree.Lock()
	defer BTree.Unlock()
	depth := parentNode.depth + 1
	blockRight := &HashTagTree{
		tag:    bestTag,
		freq:   bestFreq,
		total:  total,
		depth:  depth,
		parent: parentNode,
	}
	blockLeft := &HashTagTree{
		tag:    secTag,
		freq:   secFreq,
		total:  total,
		depth:  depth,
		parent: parentNode,
	}
	parentNode.left = blockLeft
	parentNode.right = blockRight
}

//Find looks for specific node in BinaryTree
func Find(wantNode string, root *HashTagTree) *HashTagTree {
	BTree.RLock()
	defer BTree.RUnlock()
	if BTree.root.tag == wantNode {
		return BTree.root
	}
	var found *HashTagTree
	if root != nil {
		if root.tag == wantNode {
			return root
		}
		found = Find(wantNode, root.left)
		if found == nil {
			found = Find(wantNode, root.right)
		}
		return found
	}
	return nil
}

//PrintTreeB prints all nodes of list, depth-wise
func PrintTreeB() {
	BTree.RLock()
	defer BTree.RUnlock()
	queue := New()
	Enq(queue, BTree.root)
	for !IsEmpty(queue) {
		curr := Deq(queue)
		fmt.Println(formatPrint(curr))
		if curr.left != nil {
			Enq(queue, curr.left)
		}
		if curr.right != nil {
			Enq(queue, curr.right)
		}
	}
}

//based on HashTagTree node deq-ed, formats what will be printed
func formatPrint(curr *HashTagTree) string {
	grammar1, grammar2 := "tweets", "tweets"
	if curr.freq == 1 {
		grammar1 = "tweet"
	}
	if curr.total == 1 {
		grammar1 = "tweet"
	}
	parent := first + " (initial search)"
	if curr.parent != nil {
		parent = curr.parent.tag
	}
	format := fmt.Sprintf("%d.) %d %s had %s of the %d %s that had %s", curr.depth, curr.freq, grammar1, curr.tag, curr.total, grammar2, parent)
	return format
}
