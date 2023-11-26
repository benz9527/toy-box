# Go ASM

## 常见错误
```bash
./xxx.s:<line>: unexpected EOF
asm: assembly of ./xxx.s failed
```
Root cause： RET 指令后面没有换行符（newline）