// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package arraylist

import "github.com/dairongpeng/gds/containers"

func assertEnumerableImplementation() {
	var _ containers.EnumerableWithIndex = (*List)(nil)
}

// Each calls the given function once for each element, passing that element's index and value.
// Each 对列表中的每一个位置的元素进行一遍处理，处理函数通过回调的方式传入
func (list *List) Each(f func(index int, value interface{})) {
	iterator := list.Iterator()
	for iterator.Next() {
		f(iterator.Index(), iterator.Value())
	}
}

// Map invokes the given function once for each element and returns a
// container containing the values returned by the given function.
// Map对list列表中的每一个元素进行一遍处理，值处理结果保存到一个新的list中返回给Map的调用者
func (list *List) Map(f func(index int, value interface{}) interface{}) *List {
	newList := &List{}
	iterator := list.Iterator()
	for iterator.Next() {
		newList.Add(f(iterator.Index(), iterator.Value()))
	}
	return newList
}

// Select returns a new container containing all elements for which the given function returns a true value.
// Select 对list进行一遍遍历并通过传入的回调函数f进行筛选，返回筛选后的结果list
func (list *List) Select(f func(index int, value interface{}) bool) *List {
	newList := &List{}
	iterator := list.Iterator()
	for iterator.Next() {
		if f(iterator.Index(), iterator.Value()) {
			newList.Add(iterator.Value())
		}
	}
	return newList
}

// Any passes each element of the collection to the given function and
// returns true if the function ever returns true for any element.
// Any 对每个list中的元素进行判断，判断函数f通过回调的方式传入，list只有有任意元素判断为true时，Any返回true给调用者，否则返回false
func (list *List) Any(f func(index int, value interface{}) bool) bool {
	iterator := list.Iterator()
	for iterator.Next() {
		if f(iterator.Index(), iterator.Value()) {
			return true
		}
	}
	return false
}

// All passes each element of the collection to the given function and
// returns true if the function returns true for all elements.
// All 对每个list中的元素进行判断，判断函数f通过回调的方式传入，list只有当所有元素都判断为true时，All返回true给调用者，否则返回false
func (list *List) All(f func(index int, value interface{}) bool) bool {
	iterator := list.Iterator()
	for iterator.Next() {
		if !f(iterator.Index(), iterator.Value()) {
			return false
		}
	}
	return true
}

// Find passes each element of the container to the given function and returns
// the first (index,value) for which the function is true or -1,nil otherwise
// if no element matches the criteria.
// Find 对list中的元素进行查找，如果满足判断条件，返回第一次找到的元素位置和值，否则返回-1和nil表示没找到
func (list *List) Find(f func(index int, value interface{}) bool) (int, interface{}) {
	iterator := list.Iterator()
	for iterator.Next() {
		if f(iterator.Index(), iterator.Value()) {
			return iterator.Index(), iterator.Value()
		}
	}
	return -1, nil
}
