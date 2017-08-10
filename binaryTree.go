package main

import "fmt"

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
	parentNode := Find(parent, root)
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

func Find(parent string, root *HashTagTree) *HashTagTree {
	found := &HashTagTree{}
	if root != nil {
		if root.tag == parent {
			return root
		} else {
			found = Find(parent, root.left)
			if found == nil {
				found = Find(parent, root.right)
			}
			return found
		}
	} else {
		return nil
	}
}

func PrintTreeB() {
	queue := New()

	if root.left != nil {
		Enq(queue, root.left)
	}
	if root.right != nil {
		Enq(queue, root.right)
	}
	for !IsEmpty(queue) {
		grammar1, grammar2 := "tweets", "tweets"
		curr := Deq(queue)
		if curr.freq == 1 {
			grammar1 = "tweet"
		}
		if curr.total == 1 {
			grammar1 = "tweet"
		}
		fmt.Printf("%d.) %d %s had %s of the %d %s that had %s \n", curr.level, curr.freq, grammar1, curr.tag, curr.total, grammar2, curr.parent.tag)

	}
}
