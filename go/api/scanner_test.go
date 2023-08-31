package api

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestBufioScannerInput(t *testing.T) {
	scanner := bufio.NewScanner(os.Stdin)
	// scanner.Scan() 调用一次就读一行，连续调用就会忽略前面的行
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func TestBufioScannerInput_BufferSize(t *testing.T) {
	scanner := bufio.NewScanner(os.Stdin)
	bufSize := make([]byte, 1024*4) // MaxScanTokenSize 64*1024
	scanner.Buffer(bufSize, len(bufSize))

	for {
		line := scanner.Text()
		if len(line) == 0 || line == "\r\n" {
			break
		}
		if strings.HasPrefix(line, "hi") {
			fmt.Println("hello")
		}
	}
}

func TestBufioScannerInput_MultipleLines(t *testing.T) {
	scanner := bufio.NewScanner(os.Stdin)
	// 每次读取一行
	list := make([]int, 0)
	for scanner.Scan() {
		// 结束基本是靠空字符串
		input := scanner.Text()
		if len(input) == 0 || input == "\r\n" || input == "\n" {
			break
		}
		// 分隔符
		data := strings.Split(input, " ")
		var sum int
		for i := range data {
			val, _ := strconv.Atoi(data[i])
			sum += val
		}
		list = append(list, sum)
	}
	fmt.Println(list)
}

func TestBufioScannerInput_MultipleLines2(t *testing.T) {
	var n int
	// 不能使用 fmt.Scan("%d\n", n) 的形式
	fmt.Scanln(&n)

	scanner := bufio.NewScanner(os.Stdin)
	// 每次读取一行
	list := make([]string, 0)
	for scanner.Scan() {
		// 结束基本是靠空字符串
		input := scanner.Text()
		if len(input) == 0 || input == "\r\n" || input == "\n" {
			break
		}
		// 分隔符
		data := strings.Split(input, " ")
		list = append(list, strings.Join(data, ""))
	}
	fmt.Println(n, list)
}

func TestBufioReader(t *testing.T) {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || len(line) == 0 || line == "\r\n" {
			break
		}
		if strings.HasPrefix(line, "hi") {
			fmt.Println("hello")
		}
	}
}

func TestBufioReaderRune(t *testing.T) {
	reader := bufio.NewReader(os.Stdin)
	ch, size, err := reader.ReadRune()
	if err != nil {
		return
	}
	fmt.Println(ch, size)
}

func TestACMInputMode(t *testing.T) {
	var n int
	// 不能捕获回车，如果后续仍旧还有输入的话，这个回车就会被其他输入函数捕获
	fmt.Scan(&n)
	fmt.Scanf("%d", n)
	// 可以捕获回车，不会影响后续输入
	fmt.Scanln(&n)
	fmt.Scanf("%d\n", n)
}

func TestACMInputMode_MultipleLines(t *testing.T) {
	/*
		第一行输入一个 N，第二行输入N个值
	*/
	var n int
	fmt.Scanln(&n)

	list := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Scanf("%d", &list[i])
		// 或者是下面这个
		// fmt.Scan(&list[i])
	}
}

func TestACMInputMode_LoopInput(t *testing.T) {
	/*
			先输入一个数N,回车
			再输入一行数据（但是个数未知）,回车
			再输入一个数据M回车

			1
			1,2,3,4,5
			2

			1
			1 2 3 4 5
		    2

			1
			1 2 3 4,5
		    2
	*/
	var n int
	fmt.Scanln(&n)

	list := make([]int, 0)
	for {
		tmp := 0
		// 这里就是扫描需要的数据类型，否则这里会报错 unexpected newline，结束输入
		_, err := fmt.Scanf("%d", &tmp)
		if err != nil {
			// Scanf当检查到最后一个数据的时候，err不是 nil，直接跳出就行
			break
		}
		list = append(list, tmp)
	}
	var m int
	fmt.Scanln(&m)
	fmt.Println(list, n, m)
}
