package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	f, err := os.Open("42d07f18-3a4e-4191-94e9-672e69f0e6d1.csv")
	if err != nil {
		fmt.Println("打开文件失败", err)
		os.Exit(1)
	}
	output := fmt.Sprintf("%d.csv", time.Now().UnixNano())
	wf, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 777)
	if err != nil {
		fmt.Println("打开文件失败", err)
		os.Exit(1)
	}
	defer f.Close()
	defer wf.Close()
	buf := bufio.NewReader(f)
	writeBuf := bufio.NewWriter(wf)

	i := 0

	for {
		i++
		lineBytes, err := buf.ReadBytes('\n')
		//fmt.Printf("原始数据:%v\n", string(lineBytes))
		newLine := handleLine(lineBytes)
		newLine = handleLine2(newLine)
		//fmt.Printf("处理后数据:%v\n", string(newLine))
		writeBuf.Write(newLine)
		fmt.Printf("写入第%d行\n", i)
		writeBuf.Flush()
		if err != nil {
			if err == io.EOF {
				os.Exit(0)
			}
		}
		//if i == 100 {
		//	os.Exit(0)
		//
		//}
	}
}

// 处理连续,,，转换为,"",
func handleLine2(lineBytes []byte) []byte {
	var newLineBytes []byte
	for i, lineByte := range lineBytes {
		// 连续出现,,
		if i+1 <= len(lineBytes) && lineBytes[i] == 44 && lineBytes[i+1] == 44 {
			newLineBytes = append(newLineBytes, lineByte)
			// 插入两个"
			newLineBytes = append(newLineBytes, 34)
			newLineBytes = append(newLineBytes, 34)
			continue
		}
		newLineBytes = append(newLineBytes, lineByte)
	}
	return newLineBytes
}

// 字节码：44 = ,
// 字节码：34 = "
func handleLine(lineBytes []byte) []byte {
	var positions []int
	arrayLength := len(lineBytes)
	for i, _ := range lineBytes {
		// 需要匹配"""
		if i+2 < arrayLength && lineBytes[i] == 34 && lineBytes[i+1] == 34 && lineBytes[i+2] == 34 {
			// 记录当前index
			if len(positions) == 0 {
				positions = append(positions, i)
			} else {
				positions = append(positions, i+3)
			}

		}
	}
	if len(positions) > 0 {
		// print sub string
		subBytes := lineBytes[positions[0]:positions[1]]
		prefixBytes := lineBytes[0:positions[0]]
		suffixBytes := lineBytes[positions[1]:]
		// 处理多引号包裹，结果为CSV标准模式
		var newSubBytes []byte
		for i, v := range subBytes {
			// 跳过第一个和最后一个"
			if i == 0 || i == len(subBytes)-1 {
				newSubBytes = append(newSubBytes, v)
				continue
			}
			if v == 34 {
				continue
			}
			newSubBytes = append(newSubBytes, v)
		}
		lineBytes = append(prefixBytes, newSubBytes...)
		lineBytes = append(lineBytes, suffixBytes...)
	}
	return lineBytes
}
