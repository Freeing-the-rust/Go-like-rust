// struct 정의 처리
if strings.HasPrefix(line, "struct ") {
    name := strings.TrimSpace(strings.TrimPrefix(line, "struct "))
    runtime.NewStruct(name)
    continue
}

// impl 정의
if strings.HasPrefix(line, "impl ") {
    parts := strings.Split(line, "{")
    structName := strings.TrimSpace(strings.TrimPrefix(parts[0], "impl "))
    currentStruct = runtime.StructRegistry[structName]
    continue
}
if currentStruct != nil {
    if line == "}" {
        currentStruct = nil
        continue
    }
    // 메서드 추가: 예) fn show(self) { println!(self.x); }
    if strings.HasPrefix(line, "fn ") {
        fnName := strings.TrimPrefix(line, "fn ")
        fnName = strings.Split(fnName, "(")[0]
        currentStruct.AddMethod(fnName, func(self *runtime.HeapObject, args ...*runtime.HeapObject) *runtime.HeapObject {
            fmt.Printf("[METHOD] %s.%s() called on %s\n", currentStruct.Name, fnName, self.Type)
            return runtime.Alloc("nil", nil)
        })
    }
}
