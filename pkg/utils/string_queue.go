package utils

import "errors"

type StringQueue interface {
	Add(element string) StringQueue
	RemoveFirst() (string, error)
	Length() int
}

func NewQueue() StringQueue {
	return &stringQueue{
		elements: make([]string, 0),
	}
}

type stringQueue struct {
	elements []string
}

func (q *stringQueue) Add(element string) StringQueue {
	q.elements = append(q.elements, element)
	return q
}

func (q *stringQueue) RemoveFirst() (string, error) {
	if len(q.elements) == 0 {
		return "", errors.New("empty-queue")
	}

	first := q.elements[0]
	q.elements = q.elements[1:]

	return first, nil
}

func (q *stringQueue) Length() int {
	return len(q.elements)
}
