package main

type Queue struct {
	buffer []HashTagTree
	head   int
	tail   int
	maxLen int
}

func New() *Queue {
	buffer := make([]HashTagTree, 64, 64)
	new := &Queue{
		buffer: buffer,
		maxLen: 64,
	}
	return new
}

func Enq(queue *Queue, node *HashTagTree) {
	if (queue.tail - queue.head) == queue.maxLen {
		resize(queue)
	}
	queue.buffer[queue.tail] = *node
	queue.tail++
}

func Deq(queue *Queue) *HashTagTree {
	node := &HashTagTree{}
	if (queue.tail - queue.head) > 0 {
		node = &queue.buffer[queue.head]
		queue.head++
	}
	for i := queue.head; i < queue.tail; i++ {
		queue.buffer[i-queue.head] = queue.buffer[i]
	}

	return node
}

func resize(queue *Queue) {
	var temp []HashTagTree
	if queue.tail-queue.head == queue.maxLen {
		temp = make([]HashTagTree, (queue.tail - queue.head), 2*queue.maxLen)
		for i := queue.head; i < queue.tail; i++ {
			temp[i-queue.head] = queue.buffer[i]
		}
	} else {
		temp = make([]HashTagTree, (queue.tail - queue.head), queue.maxLen/2)
		for i := queue.head; i < queue.tail; i++ {
			temp[i-queue.head] = queue.buffer[i]
		}
	}
	queue.buffer = temp
}

func IsEmpty(queue *Queue) bool {
	len := queue.tail - queue.head
	if len == 0 {
		return true
	}
	return false
}
