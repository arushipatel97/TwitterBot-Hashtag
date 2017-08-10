package main

import (
	"fmt"
	"sync"
)

type Tree struct {
	root *HashTagTree
	sync.RWMutex
}

//HashTag struct is block type for linked list
type HashTagTree struct {
	tag    string
	freq   int
	total  int
	depth  int
	parent *HashTagTree
	left   *HashTagTree
	right  *HashTagTree
}

//adds first and second most popular hashtags (corresponding nodes to BinaryTree)
func AddToTree(bestTag string, bestFreq int, secTag string, secFreq int, total int, parent string) {
	parentNode := Find(parent, BTree.root)
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
func Find(parent string, root *HashTagTree) *HashTagTree {
	BTree.RLock()
	defer BTree.RUnlock()
	var found *HashTagTree
	if root != nil {
		if root.tag == parent {
			return root
		}
		found = Find(parent, root.left)
		if found == nil {
			found = Find(parent, root.right)
		}
		return found
	}
	return nil
}

//Prints all nodes of list, depth-wise
func PrintTreeB() {
	// BTree.RLock()
	// defer BTree.RUnlock()
	queue := New()
	Enq(queue, BTree.root)
	for !IsEmpty(queue) {
		grammar1, grammar2 := "tweets", "tweets"
		curr := Deq(queue)
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
		fmt.Printf("%d.) %d %s had %s of the %d %s that had %s \n", curr.depth, curr.freq, grammar1, curr.tag, curr.total, grammar2, parent)
		if curr.left != nil {
			Enq(queue, curr.left)
		}
		if curr.right != nil {
			Enq(queue, curr.right)
		}
	}
}
