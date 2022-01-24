// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hashbidimap implements a bidirectional map backed by two hashmaps.
//
// A bidirectional map, or hash bag, is an associative data structure in which the (key,value) pairs form a one-to-one correspondence.
// Thus the binary relation is functional in each direction: value can also act as a key to key.
// A pair (a,b) thus provides a unique coupling between 'a' and 'b' so that 'b' can be found when 'a' is used as a key and 'a' can be found when 'b' is used as a key.
//
// Elements are unordered in the map.
//
// Structure is not thread safe.
//
// Reference: https://en.wikipedia.org/wiki/Bidirectional_map
package hashbidimap

import (
	"fmt"
	"github.com/dairongpeng/gds/maps"
	"github.com/dairongpeng/gds/maps/hashmap"
)

func assertMapImplementation() {
	var _ maps.BidiMap = (*Map)(nil)
}

// Map holds the elements in two hashmaps.
// Map 双向索引的Map，既可以通过key查询value也可以通过value查询key，时间复杂度都是O(1)。基于两个HashMap实现
type Map struct {
	// kay -> value的Map
	forwardMap hashmap.Map
	// value -> key的Map
	inverseMap hashmap.Map
}

// New instantiates a bidirectional map.
func New() *Map {
	return &Map{*hashmap.New(), *hashmap.New()}
}

// Put inserts element into the map.
func (m *Map) Put(key interface{}, value interface{}) {
	// 先检查key value Map中是否存在该key，存在则要移除另外一个value key map中以该value为键的元素
	if valueByKey, ok := m.forwardMap.Get(key); ok {
		m.inverseMap.Remove(valueByKey)
	}
	// 同理，检查value key map 中是否存在以value为key的元素，存在则移除key value map中对应的与元素
	if keyByValue, ok := m.inverseMap.Get(value); ok {
		m.forwardMap.Remove(keyByValue)
	}
	// 都清理干净后，安全的插入key-value map value-key map中
	m.forwardMap.Put(key, value)
	m.inverseMap.Put(value, key)
}

// Get searches the element in the map by key and returns its value or nil if key is not found in map.
// Second return parameter is true if key was found, otherwise false.
// Get 根据键查值，从forwardMap中寻找
func (m *Map) Get(key interface{}) (value interface{}, found bool) {
	return m.forwardMap.Get(key)
}

// GetKey searches the element in the map by value and returns its key or nil if value is not found in map.
// Second return parameter is true if value was found, otherwise false.
// GetKey 以value当键，查看value对应的key是否存在，从inverseMap中查找
func (m *Map) GetKey(value interface{}) (key interface{}, found bool) {
	return m.inverseMap.Get(value)
}

// Remove removes the element from the map by key.
// Remove 以key来移除元素，先查key-value的map，找到的话，需要把inverseMap中的对应元素也移除
func (m *Map) Remove(key interface{}) {
	if value, found := m.forwardMap.Get(key); found {
		m.forwardMap.Remove(key)
		m.inverseMap.Remove(value)
	}
}

// Empty returns true if map does not contain any elements
func (m *Map) Empty() bool {
	return m.Size() == 0
}

// Size returns number of elements in the map.
func (m *Map) Size() int {
	return m.forwardMap.Size()
}

// Keys returns all keys (random order).
// Keys 寻找key的列表， 即为forwardMap的key列表
func (m *Map) Keys() []interface{} {
	return m.forwardMap.Keys()
}

// Values returns all values (random order).
// Values 寻找value的列表，即为inverseMap中的key列表
func (m *Map) Values() []interface{} {
	return m.inverseMap.Keys()
}

// Clear removes all elements from the map.
// 清理hashbidmap
func (m *Map) Clear() {
	m.forwardMap.Clear()
	m.inverseMap.Clear()
}

// String returns a string representation of container
// 输出hashbidmap
func (m *Map) String() string {
	str := "HashBidiMap\n"
	str += fmt.Sprintf("%v", m.forwardMap)
	return str
}
