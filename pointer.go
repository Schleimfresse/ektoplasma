package main

import (
	"fmt"
	"log"
)

type Memory struct {
	data map[string]*Value
	next int
}

func NewMemory() *Memory {
	return &Memory{
		data: make(map[string]*Value),
		next: 1,
	}
}

func (m *Memory) Allocate(value *Value) string {
	addr := fmt.Sprintf("0x%08X", m.next)
	m.data[addr] = value
	m.next++
	return addr
}

func (m *Memory) Get(addr string) (*Value, bool) {
	val, ok := m.data[addr]
	return val, ok
}

func (m *Memory) Set(addr string, value *Value) {
	m.data[addr] = value
}

func addressOf(value *Value, memory *Memory) *Value {
	pointer := &Pointer{Addr: memory.Allocate(value)}
	return &Value{Pointer: pointer}
}

func dereference(ptr *Pointer, memory *Memory) *Value {
	val, ok := memory.Get(ptr.Addr)
	if !ok {
		panic(fmt.Sprintf("Invalid memory access at address %s", ptr.Addr))
	}
	return val
}

func assignToPointer(ptr *Pointer, value *Value, memory *Memory) {
	log.Println("Assigning value to pointer", ptr.Addr)
	memory.Set(ptr.Addr, value)
}
