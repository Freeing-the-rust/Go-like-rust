package runtime

import "fmt"

type BuiltinFunc func(...*HeapObject) *HeapObject

var builtins = map[string]BuiltinFunc{
	"add": Add,
	"sub": Sub,
}

func GetBuiltin(name string) BuiltinFunc {
	if f, ok := builtins[name]; ok {
		return f
	}
	return nil
}

func Call(fn func(...*HeapObject) *HeapObject, args ...*HeapObject) *HeapObject {
	return fn(args...)
}

func Add(args ...*HeapObject) *HeapObject {
	a := args[0].Data.(int)
	b := args[1].Data.(int)
	return Alloc("int", a+b)
}

func Sub(args ...*HeapObject) *HeapObject {
	a := args[0].Data.(int)
	b := args[1].Data.(int)
	return Alloc("int", a-b)
}

func Print(o *HeapObject) {
	fmt.Printf("[PRINT] %v (%s)\n", o.Data, o.Type)
}
