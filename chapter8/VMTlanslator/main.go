package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	assemblyFileName := getAssemblyFileName(os.Args[1])

	codeWriter := NewCodeWriter()
	codeWriter.SetFileName(assemblyFileName)
	codeWriter.WriteBootstrap()

	vmFiles := getVMFiles(os.Args[1])

	for _, file := range vmFiles {
		parser := NewParser(file)
		if parser == nil {
			panic("Failed to create parser")
		}
		codeWriter.SetVMFileName(strings.Split(strings.Split(file.Name(), "/")[len(strings.Split(file.Name(), "/"))-1], ".")[0])

		switch parser.CommandType() {
		case "C_ARITHMETIC":
			codeWriter.WriteArithmetic(parser.Arg1())
		case "C_PUSH":
			codeWriter.WritePushPop("push", parser.Arg1(), parser.Arg2())
		case "C_POP":
			codeWriter.WritePushPop("pop", parser.Arg1(), parser.Arg2())
		case "C_IF":
			codeWriter.WriteIf(parser.Arg1())
		case "C_LABEL":
			codeWriter.WriteLabel(parser.Arg1())
		case "C_GOTO":
			codeWriter.WriteGoto(parser.Arg1())
		case "C_FUNCTION":
			codeWriter.WriteFunction(parser.Arg1(), parser.Arg2())
		case "C_RETURN":
			codeWriter.WriteReturn()
		case "C_CALL":
			codeWriter.WriteCall(parser.Arg1(), parser.Arg2())
		}

		for parser.HasMoreCommands() {
			fmt.Println(parser.currentCommand)
			parser.Advance()
			switch parser.CommandType() {
			case "C_ARITHMETIC":
				codeWriter.WriteArithmetic(parser.Arg1())
			case "C_PUSH":
				codeWriter.WritePushPop("push", parser.Arg1(), parser.Arg2())
			case "C_POP":
				codeWriter.WritePushPop("pop", parser.Arg1(), parser.Arg2())
			case "C_IF":
				codeWriter.WriteIf(parser.Arg1())
			case "C_LABEL":
				codeWriter.WriteLabel(parser.Arg1())
			case "C_GOTO":
				codeWriter.WriteGoto(parser.Arg1())
			case "C_FUNCTION":
				codeWriter.WriteFunction(parser.Arg1(), parser.Arg2())
			case "C_RETURN":
				codeWriter.WriteReturn()
			case "C_CALL":
				codeWriter.WriteCall(parser.Arg1(), parser.Arg2())
			}
		}
	}

}

func getVMFiles(dirPath string) []*os.File {
	// 引数に渡されたパスがファイルだった場合、そのファイルを含む長さ1の配列を返す
	if strings.Contains(dirPath, ".vm") {
		file, err := os.Open(dirPath)
		if err != nil {
			panic(err)
		}
		return []*os.File{file}
	}

	// 引数に渡されたパスがディレクトリだった場合、そのディレクトリ内の.vmファイルを含む配列を返す
	dir, err := os.Open(dirPath)
	if err != nil {
		panic(err)
	}
	defer dir.Close()

	files, err := dir.Readdir(0)
	if err != nil {
		panic(err)
	}

	vmFiles := []*os.File{}

	for _, file := range files {
		if strings.Contains(file.Name(), ".vm") {
			fmt.Println(file.Name())
			vmFile, err := os.Open(dirPath + "/" + file.Name())
			if err != nil {
				panic(err)
			}
			vmFiles = append(vmFiles, vmFile)
		}
	}

	// Sys.vmが存在する場合、Sys.vmを配列の先頭に移動する
	for i, file := range vmFiles {
		if strings.Contains(file.Name(), "Sys.vm") {
			vmFiles[0], vmFiles[i] = vmFiles[i], vmFiles[0]
		}
	}
	return vmFiles

}

func getAssemblyFileName(input string) string {

	// コマンドの引数にvmファイルが指定された場合、そのファイル名を.asmに変換して返す
	if strings.Contains(input, ".vm") {
		return strings.Split(input, ".")[0] + ".asm"
	}

	// コマンドの引数にディレクトリが指定された場合、そのディレクトリ内の.vmファイルを含む.asmファイル名を返す
	return input + "/" + strings.Split(input, "/")[len(strings.Split(input, "/"))-1] + ".asm"
}
