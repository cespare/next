// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the THIRD-PARTY-LICENSE.txt file.

// This example demonstrates a priority queue built using a Heap.
package heap_test

import (
	"fmt"

	"github.com/cespare/next/container/heap"
)

// update modifies the priority and value of an Item in the queue.
func update(pq *heap.Heap[*Item], item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	pq.Fix(item.index)
}

// This example creates a priority queue with some items, adds and manipulates an item,
// and then removes the items in priority order.
func Example_heapPriorityQueue() {
	// Some items and their priorities.
	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := heap.New(func(item0, item1 *Item) bool {
		return item0.priority > item1.priority
	})
	pq.SetIndex = func(item *Item, i int) { item.index = i }
	var s []*Item
	i := 0
	for value, priority := range items {
		s = append(s, &Item{
			value:    value,
			priority: priority,
			index:    i,
		})
		i++
	}
	pq.Init(s)

	// Insert a new item and then modify its priority.
	item := &Item{
		value:    "orange",
		priority: 1,
	}
	pq.Push(item)
	update(pq, item, item.value, 5)

	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := pq.Pop()
		fmt.Printf("%.2d:%s ", item.priority, item.value)
	}
	// Output:
	// 05:orange 04:pear 03:banana 02:apple
}
