package structures

// Queue implementation without thread-safety
type Queue[T any] struct {
	items []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{items: make([]T, 0)}
}

func (q *Queue[T]) Push(item T) {
	q.items = append(q.items, item)
}

func (q *Queue[T]) Pop() (T, bool) {
	var zero T
	if len(q.items) == 0 {
		return zero, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

func (q *Queue[T]) Len() int {
	return len(q.items)
}
