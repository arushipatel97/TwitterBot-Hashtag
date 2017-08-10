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
	level  int
	parent *HashTagTree
	left   *HashTagTree
	right  *HashTagTree
}

func AddToTree(bestTag string, bestFreq int, secTag string, secFreq int, total int, parent string) {
	parentNode := Find(parent, BTree.root)
	BTree.Lock()
	defer BTree.Unlock()
	level := parentNode.level + 1
	blockRight := &HashTagTree{
		tag:    bestTag,
		freq:   bestFreq,
		total:  total,
		level:  level,
		left:   nil,
		right:  nil,
		parent: parentNode,
	}
	blockLeft := &HashTagTree{
		tag:    secTag,
		freq:   secFreq,
		total:  total,
		level:  level,
		left:   nil,
		right:  nil,
		parent: parentNode,
	}
	parentNode.left = blockLeft
	parentNode.right = blockRight
}

//adds first and second most popular hashtags (corresponding nodes to BinaryTree)
func Find(parent string, root *HashTagTree) *HashTagTree {
	BTree.RLock()
	defer BTree.RUnlock()
	found := &HashTagTree{}
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
		fmt.Printf("%d.) %d %s had %s of the %d %s that had %s \n", curr.level, curr.freq, grammar1, curr.tag, curr.total, grammar2, parent)
		if curr.left != nil {
			Enq(queue, curr.left)
		}
		if curr.right != nil {
			Enq(queue, curr.right)
		}
	}
}
