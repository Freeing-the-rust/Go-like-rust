#!/usr/bin/env python3
import sys
import re
import platform

"""
SpongeLang Meaning IR → NASM Assembly
Windows(PE64) + Linux ELF64 모두 지원.

지원 의미:
- print "text"
- wait-input
"""

def parse_ir(lines):
    actions = []
    for line in lines:
        line = line.strip()
        if line.startswith("print "):
            text = re.findall(r'print\s+"(.*)"', line)
            if text:
                actions.append(("print", text[0]))
        elif "wait-input" in line:
            actions.append(("wait-input", None))
    return actions


# PE64 (Windows) 생성
def gen_windows(actions):
    data = []
    code = []

    data.append("section .data")
    msg_count = 0

    for act, val in actions:
        if act == "print":
            txt = val.encode("utf-8")
            bytestr = ", ".join(str(b) for b in txt)
            data.append(f"msg{msg_count} db {bytestr}, 10")
            data.append(f"msg{msg_count}_len equ $ - msg{msg_count}")
            msg_count += 1

    code.append("default rel")
    code.append("extern GetStdHandle")
    code.append("extern WriteFile")
    code.append("extern ReadFile")
    code.append("extern ExitProcess")
    code.append("")
    code.append("section .text")
    code.append("global main")
    code.append("main:")

    msg_idx = 0

    for act, val in actions:
        if act == "print":
            code.append("    ; print")
            code.append("    mov rcx, -11")
            code.append("    call GetStdHandle")
            code.append("    mov rbx, rax")
            code.append("")
            code.append("    sub rsp, 32")
            code.append(f"    mov rcx, rbx")
            code.append(f"    lea rdx, [rel msg{msg_idx}]")
            code.append(f"    mov r8, msg{msg_idx}_len")
            code.append(f"    lea r9, [rsp]")
            code.append("    mov qword [rsp+16], 0")
            code.append("    call WriteFile")
            code.append("    add rsp, 32")
            code.append("")

            msg_idx += 1

        elif act == "wait-input":
            code.append("    ; wait-input")
            code.append("    mov rcx, -10")
            code.append("    call GetStdHandle")
            code.append("    mov rbx, rax")
            code.append("")
            code.append("    sub rsp, 32")
            code.append("    mov rcx, rbx")
            code.append("    mov rdx, rsp")
            code.append("    mov r8, 1")
            code.append("    lea r9, [rsp+8]")
            code.append("    mov qword [rsp+16], 0")
            code.append("    call ReadFile")
            code.append("    add rsp, 32")
            code.append("")

    code.append("    ; exit")
    code.append("    xor ecx, ecx")
    code.append("    call ExitProcess")

    return "\n".join(data) + "\n\n" + "\n".join(code)


# Linux ELF64 생성
def gen_linux(actions):
    data = []
    code = []

    data.append("section .data")
    msg_count = 0

    for act, val in actions:
        if act == "print":
            txt = val.encode("utf-8")
            bytestr = ", ".join(str(b) for b in txt)
            data.append(f"msg{msg_count} db {bytestr}, 10")
            data.append(f"msg{msg_count}_len equ $ - msg{msg_count}")
            msg_count += 1

    code.append("global _start")
    code.append("section .text")
    code.append("_start:")

    msg_idx = 0

    for act, val in actions:
        if act == "print":
            code.append("    mov rax, 1")
            code.append("    mov rdi, 1")
            code.append(f"    mov rsi, msg{msg_idx}")
            code.append(f"    mov rdx, msg{msg_idx}_len")
            code.append("    syscall")
            msg_idx += 1

        elif act == "wait-input":
            code.append("    mov rax, 0")
            code.append("    mov rdi, 0")
            code.append("    mov rsi, rsp")
            code.append("    mov rdx, 1")
            code.append("    syscall")

    code.append("    mov rax, 60")
    code.append("    mov rdi, 0")
    code.append("    syscall")

    return "\n".join(data) + "\n\n" + "\n".join(code)


def main():
    if len(sys.argv) < 2:
        print("Usage: sponge_to_nasm.py <ir_file>", file=sys.stderr)
        sys.exit(1)

    irpath = sys.argv[1]
    try:
        with open(irpath, "r", encoding="utf-8") as f:
            lines = f.readlines()
    except:
        print("Error: cannot read IR file", file=sys.stderr)
        sys.exit(1)

    actions = parse_ir(lines)

    # OS 감지
    is_windows = platform.system().lower().startswith("win")

    if is_windows:
        asm = gen_windows(actions)
    else:
        asm = gen_linux(actions)

    print(asm)


if __name__ == "__main__":
    main()
