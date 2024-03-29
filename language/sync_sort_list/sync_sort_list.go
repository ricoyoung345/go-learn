package sync_sort_list

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	NodeExist   = uint32(0)
	NodeDeleted = uint32(1)
)

type IntList struct {
	head   *IntNode
	length int64
}

type IntNode struct {
	value   *int32
	next    *IntNode
	deleted uint32
	mu      sync.Mutex
}

func (n *IntNode) atomicLoadNext() *IntNode {
	return (*IntNode)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&n.next))))
}

func (n *IntNode) atomicStoreNext(next *IntNode) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&n.next)), unsafe.Pointer(next))
}

func newIntNode(value int32) *IntNode {
	return &IntNode{value: &value, deleted: NodeExist}
}

func NewInt() *IntList {
	return &IntList{head: newIntNode(0)}
}

func (l *IntList) Insert(value int) bool {
	// if lock a fail, or check fail after lock
	// insert need do a loop to ensure value insert
	for {
		a := l.head
		b := a.atomicLoadNext()

		// if list is not empty keep looking
		// must use atomic.Load, to support one write and multi read
		for b != nil && *b.value < int32(value) {
			a = b
			b = b.atomicLoadNext()
		}

		// if node aleady exist, just return false
		if b != nil && *b.value == int32(value) {
			return false
		}

		// lock node a and check the following conditions
		// if check fail, just restart looking for new node 'a' from head
		// 1: a.next == b, to ensure that no other goroutine add a node after a
		// 2: a.marked == true, to ensure that a is not delete by other goroutine
		a.mu.Lock()
		if a.atomicLoadNext() != b || a.deleted == NodeDeleted {
			a.mu.Unlock()
			continue
		}

		x := newIntNode(int32(value))
		x.atomicStoreNext(b)
		a.atomicStoreNext(x)
		a.mu.Unlock()

		// increase l.length atomicly
		atomic.AddInt64(&l.length, 1)

		return true
	}
}

func (l *IntList) Delete(value int) bool {
	// if lock a fail, or check fail after lock
	// insert need do a loop to ensure value insert
	for {
		a := l.head
		b := a.atomicLoadNext()

		// if list is not empty keep looking
		// must use atomic.Load, to support one write and multi read
		for b != nil && *b.value < int32(value) {
			a = b
			b = b.atomicLoadNext()
		}

		// if node not exist, just return false
		if b == nil || *b.value != int32(value) {
			return false
		}

		// lock node b, and check whether node b is delete
		b.mu.Lock()
		if b.deleted == NodeDeleted {
			b.mu.Unlock()
			continue
		}

		// lock node a and check the following conditions
		// if check fail, just restart looking for new node 'a' from head
		// before restart, unlock a and b successively
		// 1: a.next == b, to ensure that no other goroutine add a node after a
		// 2: a.marked == true, to ensure that a is not delete by other goroutine
		a.mu.Lock()
		if a.atomicLoadNext() != b || a.deleted == NodeDeleted {
			a.mu.Unlock()
			b.mu.Unlock()
			continue
		}

		atomic.StoreUint32(&b.deleted, NodeDeleted)
		a.atomicStoreNext(b.atomicLoadNext())

		a.mu.Unlock()
		b.mu.Unlock()

		// decrease l.length atomicly
		atomic.AddInt64(&l.length, -1)

		return true
	}
}

func (l *IntList) Contains(value int) bool {
	x := l.head.atomicLoadNext()
	for x != nil && *x.value < int32(value) {
		x = x.atomicLoadNext()
	}

	if x == nil {
		return false
	}

	return *x.value == int32(value)
}

func (l *IntList) Range(f func(value int) bool) {
	x := l.head.atomicLoadNext()
	for x != nil {
		if !f(int(*x.value)) {
			break
		}
		x = x.atomicLoadNext()
	}
}

func (l *IntList) Len() int {
	return int(atomic.LoadInt64(&l.length))
}
