// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package arraylist implements the array list.
//
// Structure is not thread safe.
//
// Reference: https://en.wikipedia.org/wiki/List_%28abstract_data_type%29
package arraylist

import (
	"fmt"
	"strings"

	"github.com/dairongpeng/gds/lists"
	"github.com/dairongpeng/gds/utils"
)

// 无需显示调用，主要是提供给编译器做检查。用来保证当前List结构实现了lists.List接口
func assertListImplementation() {
	var _ lists.List = (*List)(nil)
}

// List holds the elements in a slice
// List 列表结构，基于切片实现
type List struct {
	elements []interface{}
	size     int
}

const (
	growthFactor = float32(2.0)  // growth by 100%
	shrinkFactor = float32(0.25) // shrink when size is 25% of capacity (0 means never shrink)
)

// New instantiates a new list and adds the passed values, if any, to the list
// New 实例化一个列表，如果有初始化传入的values则添加这些values
func New(values ...interface{}) *List {
	list := &List{}
	if len(values) > 0 {
		list.Add(values...)
	}
	return list
}

// Add appends a value at the end of the list
// 从列表的尾部添加值
func (list *List) Add(values ...interface{}) {
	// 判断是否触发了扩容条件。触发则进行切片2倍扩容
	list.growBy(len(values))
	for _, value := range values {
		list.elements[list.size] = value
		list.size++
	}
}

// Get returns the element at index.
// Second return parameter is true if index is within bounds of the array and array is not empty, otherwise false.
// Get 根据切片下标返回列表的元素值
// 如果下标未越界，且切片数组不为空的时候，第二个布尔的返回值为false，否则第二个返回值为false
func (list *List) Get(index int) (interface{}, bool) {

	// 检查下标是否越界，如果越界，则第二个返回值为false
	if !list.withinRange(index) {
		return nil, false
	}

	return list.elements[index], true
}

// Remove removes the element at the given index from the list.
// Remove 按照给定的下标移除列表的元素
func (list *List) Remove(index int) {

	// 检查越界
	if !list.withinRange(index) {
		return
	}

	// 解除Index位置的元素引用
	list.elements[index] = nil // cleanup reference
	// 后续元素向前移动一位。copy函数会把第二个参数范围赋值到第一个参数范围上
	// list.elements[index+1:list.size]范围填充list.elements[index:]范围，从而填补index位置被移除的位置
	copy(list.elements[index:], list.elements[index+1:list.size]) // shift to the left by one (slow operation, need ways to optimize this)
	list.size--

	// 根据缩容因子，判断是否要进行缩容，这里缩容因子默认为0.25
	list.shrink()
}

// Contains checks if elements (one or more) are present in the list.
// All elements have to be present in the list for the method to return true.
// Performance time complexity of n^2.
// Returns true if no arguments are passed at all, i.e. list is always super-list of empty list.
// Contains 检查一个或多个元素的值，在不在当前列表中。所有需要检查的元素values都在列表List中，则返回true
// 时间复杂度为n^2
// 如果没有任何参数被传递，则返回true，即List总是空集的超集。
func (list *List) Contains(values ...interface{}) bool {

	for _, searchValue := range values {
		found := false
		for _, element := range list.elements {
			if element == searchValue {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// Values returns all elements in the list.
// Values 返回所有的元素，通过一个一个数组切片返回
func (list *List) Values() []interface{} {
	newElements := make([]interface{}, list.size, list.size)
	copy(newElements, list.elements[:list.size])
	return newElements
}

// IndexOf returns index of provided element
// IndexOf 返回值等于传入value的第一次出现的下标，找不到value则返回-1
func (list *List) IndexOf(value interface{}) int {
	if list.size == 0 {
		return -1
	}
	for index, element := range list.elements {
		if element == value {
			return index
		}
	}
	return -1
}

// Empty returns true if list does not contain any elements.
// Empty 返回true当列表不包含任务元素的时候
func (list *List) Empty() bool {
	return list.size == 0
}

// Size returns number of elements within the list.
// Size 返回列表的长度
func (list *List) Size() int {
	return list.size
}

// Clear removes all elements from the list.
// Clear 移除列表内所有的元素
func (list *List) Clear() {
	list.size = 0
	list.elements = []interface{}{}
}

// Sort sorts values (in-place) using.
// Sort 通过传入的比较器来排序列表list的中的元素
func (list *List) Sort(comparator utils.Comparator) {
	if len(list.elements) < 2 {
		return
	}
	utils.Sort(list.elements[:list.size], comparator)
}

// Swap swaps the two values at the specified positions.
// Swap 交换列表两个位置元素的值
func (list *List) Swap(i, j int) {
	// 数组越界检查
	if list.withinRange(i) && list.withinRange(j) {
		list.elements[i], list.elements[j] = list.elements[j], list.elements[i]
	}
}

// Insert inserts values at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Does not do anything if position is negative or bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
// Insert 往列表List的指定位置开始追加一组元素值
func (list *List) Insert(index int, values ...interface{}) {

	// 检查指定的位置是否越界
	if !list.withinRange(index) {
		// Append
		// 如果指定的位置刚好是列表的最后一个元素的下一个位置，直接追加
		if index == list.size {
			list.Add(values...)
		}
		return
	}

	l := len(values)
	// 检查追加之后元素个数是否会触发扩容
	list.growBy(l)
	// 扩大size
	list.size += l
	copy(list.elements[index+l:], list.elements[index:list.size-l])
	copy(list.elements[index:], values)
}

// Set the value at specified index
// Does not do anything if position is negative or bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
// Set 更改列表list的指定位置上的值，如果越界则改为往list后追加
func (list *List) Set(index int, value interface{}) {

	if !list.withinRange(index) {
		// Append
		if index == list.size {
			list.Add(value)
		}
		return
	}

	list.elements[index] = value
}

// String returns a string representation of container
func (list *List) String() string {
	str := "ArrayList\n"
	values := []string{}
	for _, value := range list.elements[:list.size] {
		values = append(values, fmt.Sprintf("%v", value))
	}
	str += strings.Join(values, ", ")
	return str
}

// Check that the index is within bounds of the list
// 检查index位置在list上是否越界
func (list *List) withinRange(index int) bool {
	return index >= 0 && index < list.size
}

// 调整list的容量扩容到cap大小
func (list *List) resize(cap int) {
	newElements := make([]interface{}, cap, cap)
	copy(newElements, list.elements)
	list.elements = newElements
}

// Expand the array if necessary, i.e. capacity will be reached if we add n elements
// 检查往列表追加n长度的元素是否会触发列表的扩容动作
func (list *List) growBy(n int) {
	// When capacity is reached, grow by a factor of growthFactor and add number of elements
	currentCapacity := cap(list.elements)
	if list.size+n >= currentCapacity {
		newCapacity := int(growthFactor * float32(currentCapacity+n))
		list.resize(newCapacity)
	}
}

// Shrink the array if necessary, i.e. when size is shrinkFactor percent of current capacity
// 检查列表是否要进行缩容
func (list *List) shrink() {
	if shrinkFactor == 0.0 {
		return
	}
	// Shrink when size is at shrinkFactor * capacity
	// 当前列表的容量
	currentCapacity := cap(list.elements)
	// 缩容因子运算出来的目标容量仍大于切片的size，则切片的cap缩容到自身的size
	if list.size <= int(float32(currentCapacity)*shrinkFactor) {
		list.resize(list.size)
	}
}
