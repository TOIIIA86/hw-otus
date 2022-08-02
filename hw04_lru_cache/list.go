package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	firstItem *ListItem
	lastItem  *ListItem
	length    int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.firstItem
}

func (l *list) Back() *ListItem {
	return l.lastItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := new(ListItem)
	item.Value = v
	if l.firstItem != nil {
		l.firstItem.Prev = item
	}
	item.Next = l.firstItem
	l.firstItem = item
	if l.lastItem == nil {
		l.lastItem = item
	}
	l.length++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := new(ListItem)
	item.Value = v
	if l.lastItem != nil {
		l.lastItem.Next = item
	}
	item.Prev = l.lastItem
	l.lastItem = item
	if l.firstItem == nil {
		l.firstItem = item
	}
	l.length++
	return item
}

func (l *list) Remove(i *ListItem) {
	if l.length == 0 {
		return
	}

	if i.Prev != nil && i.Next != nil { // элемент в середине
		i.Prev.Next = i.Next
	}

	if i.Prev == nil && i.Next != nil { // элемент в начале
		i.Next.Prev = nil
		l.firstItem = i.Next
	}

	if i.Prev != nil && i.Next == nil { // элемент в конце
		i.Prev.Next = nil
		l.lastItem = i.Prev
	}

	if i.Prev == nil && i.Next == nil { // один элемент в списке
		l.firstItem = nil
		l.lastItem = nil
	}
	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil { // уже вначале
		return
	}
	i.Prev.Next = i.Next
	i.Prev = nil
	i.Next = l.firstItem
	l.firstItem.Prev = i
	l.firstItem = i
}
