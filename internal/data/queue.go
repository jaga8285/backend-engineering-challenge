package data

//Simple FIFO Queue data structure with a circular buffer. Inserts and unordered full queue reads are preformed in constant time

type FIFOQueue[T any] struct {
	queue []T
	size  int
	head  int
	count int
}

func NewFIFOQueue[T any](size int) *FIFOQueue[T] {
	return &FIFOQueue[T]{
		queue: make([]T, size),
		size:  size,
	}
}

func (q *FIFOQueue[T]) Enqueue(val T) {
	q.queue[q.head] = val
	q.head = (q.head + 1) % q.size
	if q.count < q.size {
		q.count++
	}
}

func (q FIFOQueue[T]) GetQueue() []T {
	return q.queue[:q.count]
}
