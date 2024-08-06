package main

import (
	"os"
	"strings"
)

func main() {
	assemblyFileName := getAssemblyFileName(os.Args[1])

	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	parser := NewParser(file)
	if parser == nil {
		panic("Failed to create parser")
	}
	codeWriter := NewCodeWriter()
	codeWriter.SetFileName(assemblyFileName)

	switch parser.CommandType() {
	case "C_ARITHMETIC":
		codeWriter.WriteArithmetic(parser.Arg1())
	case "C_PUSH":
		codeWriter.WritePushPop("push", parser.Arg1(), parser.Arg2())
	}

	for parser.HasMoreCommands() {
		parser.Advance()
		switch parser.CommandType() {
		case "C_ARITHMETIC":
			codeWriter.WriteArithmetic(parser.Arg1())
		case "C_PUSH":
			codeWriter.WritePushPop("push", parser.Arg1(), parser.Arg2())
		case "C_POP":
			codeWriter.WritePushPop("pop", parser.Arg1(), parser.Arg2())
		}
	}

}

func getAssemblyFileName(input string) string {
	filePathList := strings.Split(input, "/")
	vmFileName := filePathList[len(filePathList)-1]
	assemblyFileName := strings.Split(vmFileName, ".")[0] + ".asm"
	if len(filePathList) > 1 {
		assemblyFileName = strings.Join(filePathList[:len(filePathList)-1], "/") + "/" + assemblyFileName
	}
	return assemblyFileName
}

// func getVMFileName(input string) string {
// 	filePathList := strings.Split(input, "/")
// 	vmFileName := filePathList[len(filePathList)-1]
// 	return vmFileName
// }
