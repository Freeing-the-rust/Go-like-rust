package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gcrust/runtime"
)

// 간단한 인터프리터: hello.rs 같은 파일을 직접 실행
func RunScript(path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	vars := map[string]*runtime.HeapObject{}
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		switch {
		case strings.HasPrefix(line, "let "):
			// 예: let a = 10;
			parseLet(line, vars)

		case strings.HasPrefix(line, "println!"):
			parsePrint(line, vars)

		case strings.Contains(line, "=") && strings.Contains(line, "(") && strings.Contains(line, ")"):
			// 예: c = add(a, b);
			parseCall(line, vars)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func parseLet(line string, vars map[string]*runtime.HeapObject) {
	// let a = 10;
	parts := strings.Split(line, "=")
	left := strings.TrimSpace(strings.TrimPrefix(parts[0], "let"))
	right := strings.TrimSpace(strings.TrimSuffix(parts[1], ";"))

	val := 0
	fmt.Sscanf(right, "%d", &val)
	vars[left] = runtime.Alloc("int", val)
}

func parsePrint(line string, vars map[string]*runtime.HeapObject) {
	// println!(a);
	start := strings.Index(line, "(")
	end := strings.Index(line, ")")
	name := strings.TrimSpace(line[start+1 : end])
	obj := vars[name]
	runtime.Print(obj)
}

func parseCall(line string, vars map[string]*runtime.HeapObject) {
	// c = add(a, b);
	parts := strings.Split(line, "=")
	left := strings.TrimSpace(parts[0])
	right := strings.TrimSpace(strings.TrimSuffix(parts[1], ";"))

	fnName := right[:strings.Index(right, "(")]
	argsRaw := right[strings.Index(right, "(")+1 : strings.Index(right, ")")]
	argNames := strings.Split(argsRaw, ",")

	var args []*runtime.HeapObject
	for _, name := range argNames {
		name = strings.TrimSpace(name)
		args = append(args, vars[name])
	}

	switch fnName {
	case "add":
		vars[left] = runtime.Call(runtime.Add, args...)
	default:
		fmt.Printf("⚠️ Unknown function: %s\n", fnName)
	}
}
