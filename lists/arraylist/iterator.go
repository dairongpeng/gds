// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package arraylist

import "github.com/dairongpeng/gods/containers"

func assertIteratorImplementation() {
	var _ containers.ReverseIteratorWithIndex = (*Iterator)(nil)
}

// Iterator holding the iterator's state
// Iterator 是List结构的一个迭代器，保存List迭代器的状态
type Iterator struct {
	list  *List
	index int
}

// Iterator returns a stateful iterator whose values can be fetched by an index.
// Iterator 初始化一个List结构的迭代器返回给调用者
func (list *List) Iterator() Iterator {
	return Iterator{list: list, index: -1}
}

// Next moves the iterator to the next element and returns true if there was a next element in the container.
// If Next() returns true, then next element's index and value can be retrieved by Index() and Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
// Next 如果下一个元素存在，且在列表list的合理下标范围内，迭代器移动到下一个元素位置并且返回true
func (iterator *Iterator) Next() bool {
	if iterator.index < iterator.list.size {
		iterator.index++
	}
	return iterator.list.withinRange(iterator.index)
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
// Prev 如果上一个元素存在，且在列表list的合理下标范围内，迭代器移动到上一个元素位置并返回为true
func (iterator *Iterator) Prev() bool {
	if iterator.index >= 0 {
		iterator.index--
	}
	return iterator.list.withinRange(iterator.index)
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
// Value 返回迭代器所处位置上，列表在相对应的值
func (iterator *Iterator) Value() interface{} {
	return iterator.list.elements[iterator.index]
}

// Index returns the current element's index.
// Does not modify the state of the iterator.
// Index 返回迭代器当前所处的位置
func (iterator *Iterator) Index() int {
	return iterator.index
}

// Begin resets the iterator to its initial state (one-before-first)
// Call Next() to fetch the first element if any.
// Begin 重置列表的迭代器
func (iterator *Iterator) Begin() {
	iterator.index = -1
}

// End moves the iterator past the last element (one-past-the-end).
// Call Prev() to fetch the last element if any.
// End 移动列表List迭代器到最后位置，也就是列表的最后一个元素位置的下一个位置
func (iterator *Iterator) End() {
	iterator.index = iterator.list.size
}

// First moves the iterator to the first element and returns true if there was a first element in the container.
// If First() returns true, then first element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
// First 移动迭代器到列表的第一个位置上
func (iterator *Iterator) First() bool {
	iterator.Begin()
	return iterator.Next()
}

// Last moves the iterator to the last element and returns true if there was a last element in the container.
// If Last() returns true, then last element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
// Last 移动迭代器到列表的最后一个元素的位置
func (iterator *Iterator) Last() bool {
	iterator.End()
	return iterator.Prev()
}
