# Plan9 ASM for Go

## 缩写
- SP: Stack Pointer： The highest address within the local stack frame
- SB: Static Base Pointer: Global symbols
- PC: Program Counter: Jumps and branches
- FP: Frame Pointer: The arguments and local variables of the current function

## 基本指令
### 栈调整
```asm 
SUBQ $0x18, SP // 对 SP 做减法，为函数分配函数栈帧
ADDQ $0x18, SP // 对 SP 做加法，清除函数栈帧
```

### 数据移动
```asm
MOVB $1, DI     // 移动 1 byte move byte
MOVW $0x10, BX  // 移动 2 bytes move word
MOVD $1, DX     // 移动 4 bytes move dword
MOVQ $0x0, AX   // 移动 8 bytes move qword
```
plan9 的汇编的操作数的方向是和 intel 汇编相反的，与 AT&T 类似
```plaintext
MOVQ $0x10, AX ===== mov rax, 0x10
       |    |------------|      |
       |------------------------|
```

### 计算指令
```asm
ADDQ  AX, BX   // BX += AX 8 bytes
SUBQ  AX, BX   // BX -= AX 8 bytes
IMULQ AX, BX   // BX *= AX 8 bytes
```
长度不同的操作数，通过不同的后缀来调用 ADDQ/ADDW/ADDL/ADDB

### 条件跳转
```asm
// 无条件跳转
JMP addr   // 跳转到地址，地址可为代码中的地址，不过实际上手写不会出现这种东西
JMP label  // 跳转到标签，可以跳转到同一函数内的标签位置
JMP 2(PC)  // 以当前指令为基础，向前/后跳转 x 行
JMP -2(PC) // 同上

// 有条件跳转
JZ target // 如果 zero flag 被 set 过，则跳转
```

### 通用寄存器
- amd64
  ```bash
  (lldb) reg read
    General Purpose Registers:
    rax = 0x0000000000000005
    rbx = 0x000000c420088000
    rcx = 0x0000000000000000
    rdx = 0x0000000000000000
    rdi = 0x000000c420088008
    rsi = 0x0000000000000000
    rbp = 0x000000c420047f78
    rsp = 0x000000c420047ed8
    r8 = 0x0000000000000004
    r9 = 0x0000000000000000
    r10 = 0x000000c420020001
    r11 = 0x0000000000000202
    r12 = 0x0000000000000000
    r13 = 0x00000000000000f1
    r14 = 0x0000000000000011
    r15 = 0x0000000000000001
    rip = 0x000000000108ef85  int`main.main + 213 at int.go:19
    rflags = 0x0000000000000212
    cs = 0x000000000000002b
    fs = 0x0000000000000000
    gs = 0x0000000000000000
  ```
  在 plan9 汇编里都是可以使用的，应用代码层面会用到的通用寄存器主要是: rax, rbx, rcx, rdx, rdi, rsi, r8~r15 这 14 个寄存器，虽然 rbp 和 rsp 也可以用，不过 bp 和 sp 会被用来管理栈顶和栈底，最好不要拿来进行运算。
  plan9 中使用寄存器不需要带 r 或 e 的前缀，例如 rax，只要写 AX 即可
  下面是通用通用寄存器的名字在 X64 和 plan9 中的对应关系:
  X64	rax	rbx	rcx	rdx	rdi	rsi	rbp	rsp	r8	r9	r10	r11	r12	r13	r14	rip
  Plan9	AX	BX	CX	DX	DI	SI	BP	SP	R8	R9	R10	R11	R12	R13	R14	PC
- 伪寄存器
  - FP
    - 使用形如 symbol+offset(FP) 的方式，引用函数的输入参数。例如 arg0+0(FP)，arg1+8(FP)，使用 FP 不加 symbol 时，无法通过编译，在汇编层面来讲，symbol 并没有什么用，加 symbol 主要是为了提升代码可读性。另外，官方文档虽然将伪寄存器 FP 称之为 frame pointer，实际上它根本不是 frame pointer，按照传统的 x86 的习惯来讲，frame pointer 是指向整个 stack frame 底部的 BP 寄存器。假如当前的 callee 函数是 add，在 add 的代码中引用 FP，该 FP 指向的位置不在 callee 的 stack frame 之内，而是在 caller 的 stack frame 上。
  - PC
    - 实际上就是在体系结构的知识中常见的 pc 寄存器，在 x86 平台下对应 ip 寄存器，amd64 上则是 rip。除了个别跳转之外，手写 plan9 代码与 PC 寄存器打交道的情况较少。
  - SB
    - 全局静态基指针，一般用来声明函数或全局变量
  - SP
    - plan9 的这个 SP 寄存器指向当前栈帧的局部变量的开始位置，使用形如 symbol+offset(SP) 的方式，引用函数的局部变量。offset 的合法取值是 [-framesize, 0)，注意是个左闭右开的区间。假如局部变量都是 8 字节，那么第一个局部变量就可以用 localvar0-8(SP) 来表示。这也是一个词不表意的寄存器。与硬件寄存器 SP 是两个不同的东西，在栈帧 size 为 0 的情况下，伪寄存器 SP 和硬件寄存器 SP 指向同一位置。手写汇编代码时，如果是 symbol+offset(SP) 形式，则表示伪寄存器 SP。如果是 offset(SP) 则表示硬件寄存器 SP。务必注意。对于编译输出(go tool compile -S / go tool objdump)的代码来讲，目前所有的 SP 都是硬件寄存器 SP，无论是否带 symbol。
  ```plaintext
  1. 伪 SP 和硬件 SP 不是一回事，在手写代码时，伪 SP 和硬件 SP 的区分方法是看该 SP 前是否有 symbol。如果有 symbol，那么即为伪寄存器，如果没有，那么说明是硬件 SP 寄存器。
  2. SP 和 FP 的相对位置是会变的，所以不应该尝试用伪 SP 寄存器去找那些用 FP + offset 来引用的值，例如函数的入参和返回值。
  3. 官方文档中说的伪 SP 指向 stack 的 top，是有问题的。其指向的局部变量位置实际上是整个栈的栈底(除 caller BP 之外)，所以说 bottom 更合适一些。
  4. 在 go tool objdump/go tool compile -S 输出的代码中，是没有伪 SP 和 FP 寄存器的，我们上面说的区分伪 SP 和硬件 SP 寄存器的方法，对于上述两个命令的输出结果是没法使用的。在编译和反汇编的结果中，只有真实的 SP 寄存器。
  5. FP 和 Go 的官方源代码里的 framepointer 不是一回事，源代码里的 framepointer 指的是 caller BP 寄存器的值，在这里和 caller 的伪 SP 是值是相等的。
  ```
  
### 变量声明
在汇编里所谓的变量，一般是存储在 .rodata 或者 .data 段中的只读值。对应到应用层的话，就是已初始化过的全局的 const、var、static 变量/常量。

使用 DATA 结合 GLOBL 来定义一个变量。DATA 的用法为:
```asm
DATA    symbol+offset(SB)/width, value
```
大多数参数都是字面意思，不过这个 offset 需要稍微注意。其含义是该值相对于符号 symbol 的偏移，而不是相对于全局某个地址的偏移。

使用 GLOBL 指令将变量声明为 global，额外接收两个参数，一个是 flag，另一个是变量的总大小。
```asm
GLOBL divtab(SB), RODATA, $64
GLOBL 必须跟在 DATA 指令之后，下面是一个定义了多个 readonly 的全局变量的完整例子:

DATA age+0x00(SB)/4, $18  // forever 18
GLOBL age(SB), RODATA, $4

DATA pi+0(SB)/8, $3.1415926
GLOBL pi(SB), RODATA, $8

DATA birthYear+0(SB)/4, $1988
GLOBL birthYear(SB), RODATA, $4
```

正如之前所说，所有符号在声明时，其 offset 一般都是 0。

有时也可能会想在全局变量中定义数组，或字符串，这时候就需要用上非 0 的 offset 了，例如:
```asm 
DATA bio<>+0(SB)/8, $"oh yes i"
DATA bio<>+8(SB)/8, $"am here "
GLOBL bio<>(SB), RODATA, $16
```
大部分都比较好理解，不过这里我们又引入了新的标记 <>，这个跟在符号名之后，表示该全局变量只在当前文件中生效，类似于 C 语言中的 static。如果在另外文件中引用该变量的话，会报 relocation target not found 的错误。

本小节中提到的 flag，还可以有其它的取值:

- NOPROF = 1
(For TEXT items.) Don't profile the marked function. This flag is deprecated.
- DUPOK = 2
It is legal to have multiple instances of this symbol in a single binary. The linker will choose one of the duplicates to use.
- NOSPLIT = 4
(For TEXT items.) Don't insert the preamble to check if the stack must be split. The frame for the routine, plus anything it calls, must fit in the spare space at the top of the stack segment. Used to protect routines such as the stack splitting code itself.
- RODATA = 8
(For DATA and GLOBL items.) Put this data in a read-only section.
- NOPTR = 16
(For DATA and GLOBL items.) This data contains no pointers and therefore does not need to be scanned by the garbage collector.
- WRAPPER = 32
(For TEXT items.) This is a wrapper function and should not count as disabling recover.
  - NEEDCTXT = 64
  (For TEXT items.) This function is a closure so it uses its incoming context register.
  当使用这些 flag 的字面量时，需要在汇编文件中添加 #include "textflag.h"，其在 go SDK `src/runtime/textflag.h` 中。

#### .s 和 .go 文件中的变量互用
refer.go:
```go
package main

var a = 999
func get() int

func main() {
  println(get())
}
```
refer.s:
```asm
#include "textflag.h"

TEXT ·get(SB), NOSPLIT, $0-8
MOVQ ·a(SB), AX
MOVQ AX, ret+0(FP)
RET
```

·a(SB)，表示该符号需要链接器来帮我们进行重定向(relocation)，如果找不到该符号，会输出 relocation target not found 的错误。

### 函数声明
```asm
#include "textflag.h"

// 只能在同一个包下面的任意 .go 文件中声明只有函数头，没有函数体的函数
// 这里使用 TEXT 是因为代码在二进制问中是存储子啊 .text 段中的。
// plan9 TEXT 就是专门用来定义函数的。
// 这里的 packagename 是应该省略的，表示当前包，否则包的重命名又得发生改变
// 中点 · 比较特殊，是一个 unicode 的中点，该点在 mac 下的输入方法是 option+shift+9。在程序被链接之后，
// 所有的中点· 都会被替换为句号.，比如你的方法是 runtime·main，在编译之后的程序里的符号则是 runtime.main
// func add(a, b int) int
TEXT packagename·add(SB),NODPLIT,$0-8
  MOVQ a+0(FP), AX
  MOVQ b+8(FP), BX
  ADDQ AX, BX
  MOVQ BX, ret+16(FP)
  RET
```
```plaintext

                              参数及返回值大小
                                      | 
 TEXT packagename·add(SB),NOSPLIT,$32-32
       |           |               |
      包名       函数名         栈帧大小(局部变量+可能需要的额外调用函数的参数空间的总大小，但不包括调用其它函数时的 ret address 的大小)

```

### 函数的栈结构
```plaintext
                       -----------------                                           
                       current func arg0                                           
                       ----------------- <----------- FP(pseudo FP)                
                        caller ret addr                                            
                       +---------------+                                           
                       | caller BP(*)  | <----------- Caller 的 BP 寄存器，在编译期有编译器插入                                          
                       ----------------- <----------- SP(pseudo SP，实际上是当前栈帧的 BP 位置)
                       |   Local Var0  |                                           
                       -----------------                                           
                       |   Local Var1  |                                           
                       -----------------                                           
                       |   Local Var2  |                                           
                       -----------------                -                          
                       |   ........    |                                           
                       -----------------                                           
                       |   Local VarN  |                                           
                       -----------------                                           
                       |               |                                           
                       |               |                                           
                       |  temporarily  |                                           
                       |  unused space |                                           
                       |               |                                           
                       |               |                                           
                       -----------------                                           
                       |  call retn    |                                           
                       -----------------                                           
                       |  call ret(n-1)|                                           
                       -----------------                                           
                       |  ..........   |                                           
                       -----------------                                           
                       |  call ret1    |                                           
                       -----------------                                           
                       |  call argn    |                                           
                       -----------------                                           
                       |   .....       |                                           
                       -----------------                                           
                       |  call arg3    |                                           
                       -----------------                                           
                       |  call arg2    |                                           
                       |---------------|                                           
                       |  call arg1    |                                           
                       -----------------   <------------  hardware SP 位置           
                         return addr                                               
                       +---------------+
```
FP 伪寄存器指向函数的传入参数的开始位置，因为栈是朝低地址方向增长，为了通过寄存器引用参数时方便，所以参数的摆放方向和栈的增长方向是相反的
```plaintext
                              FP
high ----------------------> low
argN, ... arg3, arg2, arg1, arg0
```
假设所有参数均为 8 字节，这样我们就可以用 symname+0(FP) 访问第一个 参数，symname+8(FP) 访问第二个参数，以此类推。用伪 SP 来引用局部变量，原理上来讲差不多，不过因为伪 SP 指向的是局部变量的底部，所以 symname-8(SP) 表示的是第一个局部变量，symname-16(SP)表示第二个，以此类推。当然，这里假设局部变量都占用 8 个字节。

```plaintext
                                                                                                                              
                                       caller                                                                                 
                                 +------------------+                                                                         
                                 |                  |                                                                         
       +---------------------->  --------------------                                                                         
       |                         |                  |                                                                         
       |                         | caller parent BP |                                                                         
       |           BP(pseudo SP) --------------------                                                                         
       |                         |                  |                                                                         
       |                         |   Local Var0     |                                                                         
       |                         --------------------                                                                         
       |                         |                  |                                                                         
       |                         |   .......        |                                                                         
       |                         --------------------                                                                         
       |                         |                  |                                                                         
       |                         |   Local VarN     |                                                                         
                                 --------------------                                                                         
 caller stack frame              |                  |                                                                         
                                 |   callee arg2    |                                                                         
       |                         |------------------|                                                                         
       |                         |                  |                                                                         
       |                         |   callee arg1    |                                                                         
       |                         |------------------|                                                                         
       |                         |                  |                                                                         
       |                         |   callee arg0    |                                                                         
       |                         ----------------------------------------------+   FP(virtual register)                       
       |                         |                  |                          |                                              
       |                         |   return addr    |  parent return address   |                                              
       +---------------------->  +------------------+---------------------------    <-------------------------------+         
                                                    |  caller BP               |                                    |         
                                                    |  (caller frame pointer)  |                                    |         
                                     BP(pseudo SP)  ----------------------------                                    |         
                                                    |                          |                                    |         
                                                    |     Local Var0           |                                    |         
                                                    ----------------------------                                    |         
                                                    |                          |                                              
                                                    |     Local Var1           |                                              
                                                    ----------------------------                            callee stack frame
                                                    |                          |                                              
                                                    |       .....              |                                              
                                                    ----------------------------                                    |         
                                                    |                          |                                    |         
                                                    |     Local VarN           |                                    |         
                                  SP(Real Register) ----------------------------                                    |         
                                                    |                          |                                    |         
                                                    |                          |                                    |         
                                                    |                          |                                    |         
                                                    |                          |                                    |         
                                                    |                          |                                    |         
                                                    +--------------------------+    <-------------------------------+         
                                                                                                                              
                                                              callee
```

### 地址运算
load effective address
```asm
LEAQ (BX)(AX*8), CX
// 上面代码中的 8 代表 scale
// scale 只能是 0、2、4、8
// 如果写成其它值:
// LEAQ (BX)(AX*3), CX
// ./a.s:6: bad scale: 3

// 用 LEAQ 的话，即使是两个寄存器值直接相加，也必须提供 scale
// 下面这样是不行的
// LEAQ (BX)(AX), CX
// asm: asmidx: bad address 0/2064/2067
// 正确的写法是
LEAQ (BX)(AX*1), CX


// 在寄存器运算的基础上，可以加上额外的 offset
LEAQ 16(BX)(AX*1), CX

// 三个寄存器做运算，还是别想了
// LEAQ DX(BX)(AX*8), CX
// ./a.s:13: expected end of operand, found (
```



