package main

import (
	"fmt"
	"gcrust/runtime"
)

func main() {
	fmt.Println("🦀 Go-like-Rust Runtime Start")

	// 힙에 값 생성
	a := runtime.Alloc("int", 10)
	b := runtime.Alloc("int", 20)

	// Rust식 add 함수 호출
	result := runtime.Call(runtime.Add, a, b)
	runtime.Print(result)

	fmt.Println("✅ Done (heap-only execution)")
}
