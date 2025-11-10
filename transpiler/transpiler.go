package transpiler

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TranspileRustToGo translates extended Rust-like syntax to Go code
func TranspileRustToGo(inputFile string) (string, error) {
	in, err := os.Open(inputFile)
	if err != nil {
		return "", err
	}
	defer in.Close()

	outputFile := strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "_gen.go"
	out, err := os.Create(outputFile)
	if err != nil {
		return "", err
	}
	defer out.Close()

	writer := bufio.NewWriter(out)

	fmt.Fprintln(writer, "package main")
	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, "import (")
	fmt.Fprintln(writer, `    "fmt"`)
	fmt.Fprintln(writer, `    "gcrust/runtime"`)
	fmt.Fprintln(writer, ")")
	fmt.Fprintln(writer, "")

	structMode := false
	traitMode := false
	implMode := false
	var structName, traitName string

	scanner := bufio.NewScanner(in)
	fmt.Fprintln(writer, "func main() {")
	fmt.Fprintln(writer, `    fmt.Println("🦀 Transpiled Go code start")`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		// struct
		if strings.HasPrefix(line, "struct ") {
			structMode = true
			structName = strings.TrimSpace(strings.TrimPrefix(line, "struct "))
			fmt.Fprintf(writer, "type %s struct {\n", structName)
			continue
		}
		if structMode {
			if line == "}" {
				fmt.Fprintln(writer, "}")
				structMode = false
				continue
			}
			// 예: x: i32;
			field := strings.Split(line, ":")[0]
			field = strings.TrimSpace(strings.TrimSuffix(field, ";"))
			fmt.Fprintf(writer, "    %s *runtime.HeapObject\n", field)
			continue
		}

		// trait
		if strings.HasPrefix(line, "trait ") {
			traitMode = true
			traitName = strings.TrimSpace(strings.TrimPrefix(line, "trait "))
			fmt.Fprintf(writer, "type %s interface {\n", traitName)
			continue
		}
		if traitMode {
			if line == "}" {
				fmt.Fprintln(writer, "}")
				traitMode = false
				continue
			}
			if strings.HasPrefix(line, "fn ") {
				fnName := strings.Split(strings.TrimPrefix(line, "fn "), "(")[0]
				fmt.Fprintf(writer, "    %s()\n", strings.TrimSpace(fnName))
			}
			continue
		}

		// impl
		if strings.HasPrefix(line, "impl ") {
			implMode = true
			if strings.Contains(line, "for") {
				parts := strings.Split(line, "for")
				traitName = strings.TrimSpace(strings.TrimPrefix(parts[0], "impl "))
				structName = strings.TrimSpace(strings.TrimSuffix(parts[1], "{"))
			} else {
				structName = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "impl "), "{"))
			}
			continue
		}
		if implMode {
			if line == "}" {
				implMode = false
				continue
			}
			if strings.HasPrefix(line, "fn ") {
				fnName := strings.Split(strings.TrimPrefix(line, "fn "), "(")[0]
				fmt.Fprintf(writer, "func (p *%s) %s() {\n", structName, fnName)
				fmt.Fprintln(writer, "    fmt.Println(\"[METHOD CALL]\")")
				fmt.Fprintln(writer, "}")
			}
			continue
		}

		// 일반 문장
		goLine := translateStatement(line)
		if goLine != "" {
			fmt.Fprintf(writer, "    %s\n", goLine)
		}
	}

	fmt.Fprintln(writer, `    fmt.Println("✅ Transpiled Done")`)
	fmt.Fprintln(writer, "}")
	writer.Flush()
	fmt.Printf("✅ Transpiled → %s\n", outputFile)
	return outputFile, nil
}

func translateStatement(line string) string {
	switch {
	case strings.HasPrefix(line, "let "):
		parts := strings.Split(line, "=")
		left := strings.TrimSpace(strings.TrimPrefix(parts[0], "let"))
		right := strings.TrimSpace(strings.TrimSuffix(parts[1], ";"))
		return fmt.Sprintf(`%s := runtime.Alloc("int", %s)`, left, right)

	case strings.HasPrefix(line, "println!"):
		start := strings.Index(line, "(")
		end := strings.Index(line, ")")
		arg := strings.TrimSpace(line[start+1 : end])
		return fmt.Sprintf("runtime.Print(%s)", arg)

	case strings.HasPrefix(line, "if "):
		cond := strings.TrimPrefix(line, "if ")
		cond = strings.TrimSuffix(cond, "{")
		return fmt.Sprintf("if %s.Data.(int) > 0 {", strings.TrimSpace(cond))

	case strings.HasPrefix(line, "loop "):
		count := 0
		fmt.Sscanf(line, "loop %d", &count)
		return fmt.Sprintf("for i := 0; i < %d; i++ {", count)

	case line == "}":
		return "}"

	case strings.Contains(line, "=") && strings.Contains(line, "("):
		parts := strings.Split(line, "=")
		left := strings.TrimSpace(parts[0])
		right := strings.TrimSpace(strings.TrimSuffix(parts[1], ";"))
		fn := right[:strings.Index(right, "(")]
		args := right[strings.Index(right, "(")+1 : strings.Index(right, ")")]
		return fmt.Sprintf(`%s = runtime.Call(runtime.%s, %s)`, left, capitalize(fn), args)

	default:
		return fmt.Sprintf("// unhandled: %s", line)
	}
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
