package runtime

import "fmt"

type StructDef struct {
	Name   string
	Fields map[string]*HeapObject
	Impls  map[string]func(*HeapObject, ...*HeapObject) *HeapObject
}

var StructRegistry = map[string]*StructDef{}

func NewStruct(name string) *StructDef {
	s := &StructDef{
		Name:   name,
		Fields: map[string]*HeapObject{},
		Impls:  map[string]func(*HeapObject, ...*HeapObject) *HeapObject{},
	}
	StructRegistry[name] = s
	return s
}

func (s *StructDef) SetField(name string, val *HeapObject) {
	s.Fields[name] = val
}

func (s *StructDef) GetField(name string) *HeapObject {
	return s.Fields[name]
}

func (s *StructDef) AddMethod(name string, f func(*HeapObject, ...*HeapObject) *HeapObject) {
	s.Impls[name] = f
}

func (s *StructDef) CallMethod(name string, self *HeapObject, args ...*HeapObject) *HeapObject {
	if fn, ok := s.Impls[name]; ok {
		return fn(self, args...)
	}
	fmt.Printf("⚠️ Method not found: %s\n", name)
	return nil
}

// Trait 정의용
type Trait struct {
	Name   string
	Funcs  []string
}

var TraitRegistry = map[string]*Trait{}

func NewTrait(name string, funcs ...string) {
	TraitRegistry[name] = &Trait{Name: name, Funcs: funcs}
}
