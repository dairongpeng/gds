// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package treemap implements a map backed by red-black tree.
//
// Elements are ordered by key in the map.
//
// Structure is not thread safe.
//
// Reference: http://en.wikipedia.org/wiki/Associative_array
package treemap

import (
	"fmt"
	"github.com/dairongpeng/gds/maps"
	rbt "github.com/dairongpeng/gds/trees/redblacktree"
	"github.com/dairongpeng/gds/utils"
	"strings"
)

func assertMapImplementation() {
	var _ maps.Map = (*Map)(nil)
}

// Map holds the elements in a red-black tree
// Map 有序的map，基于红黑树实现。
// 红黑树是二叉搜索树，对于任意节点，左子树的节点都比自己小，右子树的节点都比自己大
type Map struct {
	tree *rbt.Tree
}

// NewWith instantiates a tree map with the custom comparator.
// NewWith 实例化一个有序表treemap，传入自定义的比较器
func NewWith(comparator utils.Comparator) *Map {
	return &Map{tree: rbt.NewWith(comparator)}
}

// NewWithIntComparator instantiates a tree map with the IntComparator, i.e. keys are of type int.
// NewWithIntComparator 创建一个int类型的有序表treemap
func NewWithIntComparator() *Map {
	return &Map{tree: rbt.NewWithIntComparator()}
}

// NewWithStringComparator instantiates a tree map with the StringComparator, i.e. keys are of type string.
// NewWithStringComparator 创建一个string类型的有序表treemap
func NewWithStringComparator() *Map {
	return &Map{tree: rbt.NewWithStringComparator()}
}

// Put inserts key-value pair into the map.
// Key should adhere to the comparator's type assertion, otherwise method panics.
// Put 往有序表treemap中put一组k-v, key是可比较类型，基于传入的比较器
func (m *Map) Put(key interface{}, value interface{}) {
	m.tree.Put(key, value)
}

// Get searches the element in the map by key and returns its value or nil if key is not found in tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
// Get 是从有序表中查找key对应的元素是否存在，存在则返回该key对应的value且用ok标识是否查找成功
func (m *Map) Get(key interface{}) (value interface{}, found bool) {
	return m.tree.Get(key)
}

// Remove removes the element from the map by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
// Remove 从有序表中按照key移除一个元素
func (m *Map) Remove(key interface{}) {
	m.tree.Remove(key)
}

// Empty returns true if map does not contain any elements
// Empty 检查treemap是否是空的
func (m *Map) Empty() bool {
	return m.tree.Empty()
}

// Size returns number of elements in the map.
func (m *Map) Size() int {
	return m.tree.Size()
}

// Keys returns all keys in-order
// Keys 会遍历输出有序表treemap的所有key，且是有序的数据，排序基于传入的比较器
func (m *Map) Keys() []interface{} {
	return m.tree.Keys()
}

// Values returns all values in-order based on the key.
// Values 会遍历输出有序表treemap的所有value，且是基于key排序后对应的value顺序，排序基于传入的比较器
func (m *Map) Values() []interface{} {
	return m.tree.Values()
}

// Clear removes all elements from the map.
func (m *Map) Clear() {
	m.tree.Clear()
}

// Min returns the minimum key and its value from the tree map.
// Returns nil, nil if map is empty.
// Min 找到有序表treemap中最小的key对应的k-v,如果treemap为空，则返回nil
func (m *Map) Min() (key interface{}, value interface{}) {
	if node := m.tree.Left(); node != nil {
		return node.Key, node.Value
	}
	return nil, nil
}

// Max returns the maximum key and its value from the tree map.
// Returns nil, nil if map is empty.
// Max 找到有序表treemap中最大key对应的key-value，如果treemap为空，则返回nil
func (m *Map) Max() (key interface{}, value interface{}) {
	if node := m.tree.Right(); node != nil {
		return node.Key, node.Value
	}
	return nil, nil
}

// Floor finds the floor key-value pair for the input key.
// In case that no floor is found, then both returned values will be nil.
// It's generally enough to check the first value (key) for nil, which determines if floor was found.
//
// Floor key is defined as the largest key that is smaller than or equal to the given key.
// A floor key may not be found, either because the map is empty, or because
// all keys in the map are larger than the given key.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
// Floor 查找给定key在treemap中对应的前置节点。
func (m *Map) Floor(key interface{}) (foundKey interface{}, foundValue interface{}) {
	node, found := m.tree.Floor(key)
	if found {
		return node.Key, node.Value
	}
	return nil, nil
}

// Ceiling finds the ceiling key-value pair for the input key.
// In case that no ceiling is found, then both returned values will be nil.
// It's generally enough to check the first value (key) for nil, which determines if ceiling was found.
//
// Ceiling key is defined as the smallest key that is larger than or equal to the given key.
// A ceiling key may not be found, either because the map is empty, or because
// all keys in the map are smaller than the given key.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
// Ceiling  查找给定key在treemap中对应的后置节点。
func (m *Map) Ceiling(key interface{}) (foundKey interface{}, foundValue interface{}) {
	node, found := m.tree.Ceiling(key)
	if found {
		return node.Key, node.Value
	}
	return nil, nil
}

// String returns a string representation of container
func (m *Map) String() string {
	str := "TreeMap\nmap["
	it := m.Iterator()
	for it.Next() {
		str += fmt.Sprintf("%v:%v ", it.Key(), it.Value())
	}
	return strings.TrimRight(str, " ") + "]"

}
