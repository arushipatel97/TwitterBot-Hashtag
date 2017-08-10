package main

import (
	"fmt"
	"strings"
)

//HashTag struct is block type for linked list
type HashTag struct {
	tag   string
	freq  int
	total int
	next  *HashTag
	prev  *HashTag
}

//PrintList goes through linked list printing the most popular hashtags/order of searching
//with frequency
func PrintList(first string) {
	count := 1
	for temp := startList; temp != nil; temp = temp.next {
		grammar1, grammar2 := "tweets", "tweets"
		if temp.prev != nil {
			prev := temp.prev.tag
			if temp.freq == 1 {
				grammar1 = "tweet"
			}
			if temp.total == 1 {
				grammar2 = "tweet"
			}
			fmt.Printf("%d.) %d %s had %s of the %d %s that had %s \n", count, temp.freq, grammar1, temp.tag, temp.total, grammar2, prev)
			count++
		}
	}
}

//AddToList adds next hashtag to be searched in linked list
func AddToList(text string, frequency int, total int) {
	block := &HashTag{
		tag:   text,
		freq:  frequency,
		total: total,
		next:  nil,
		prev:  nil,
	}
	var temp *HashTag
	for temp = startList; temp.next != nil; temp = temp.next {
	}
	temp.next = block
	block.prev = temp
}

//InList checks to see if the given string is already a tag for a block already
//in the list
func InList(tag string) bool {
	var temp *HashTag
	for temp = startList; temp != nil; temp = temp.next {
		if strings.EqualFold(tag, temp.tag) { //case insensitive check
			return true
		}
	}
	return false
}
