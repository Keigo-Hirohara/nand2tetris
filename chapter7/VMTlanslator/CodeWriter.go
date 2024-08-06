package main

import (
	"fmt"
	"os"

	"github.com/google/uuid"
)

type CodeWriter struct {
	file *os.File
}

func NewCodeWriter() *CodeWriter {
	return &CodeWriter{
		file: nil,
	}
}

func (cw *CodeWriter) SetFileName(fileName string) {
	existingAsmFile, err := os.Open(fileName)
	if err == nil && existingAsmFile != nil {
		existingAsmFile.Close()
		err := os.Remove(fileName)
		if err != nil {
			fmt.Println("Failed to remove existing file")
			panic(err)
		}
	}
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	cw.file = file
}

func (cw *CodeWriter) WriteArithmetic(command string) {
	switch command {
	case "add":
		cw.writeAdd()
	case "sub":
		cw.writeSub()
	case "neg":
		cw.writeNeg()
	case "eq":
		cw.writeEq()
	case "gt":
		cw.writeGt()
	case "lt":
		cw.writeLt()
	case "and":
		cw.writeAnd()
	case "or":
		cw.writeOr()
	case "not":
		cw.writeNot()
	}
}

func (cw *CodeWriter) writeAdd() {
	cw.writeToFile("@SP\n")
	cw.writeToFile("AM=M-1\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("A=A-1\n")
	cw.writeToFile("M=M+D\n")
}

func (cw *CodeWriter) writeSub() {
	cw.writeToFile("@SP\n")
	cw.writeToFile("AM=M-1\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("A=A-1\n")
	cw.writeToFile("M=M-D\n")
}

func (cw *CodeWriter) writeNeg() {
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M-1\n")
	cw.writeToFile("M=-M\n")
}

func (cw *CodeWriter) writeEq() {
	trueLabel := fmt.Sprintf("EQ_TRUE_%s", GenerateUniqueID())
	endLabel := fmt.Sprintf("EQ_END_%s", GenerateUniqueID())
	cw.writeToFile("@SP\n")
	cw.writeToFile("AM=M-1\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("A=A-1\n")
	cw.writeToFile("D=M-D\n")
	cw.writeToFile("@%s\n", trueLabel)
	cw.writeToFile("D;JEQ\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M-1\n")
	cw.writeToFile("M=0\n")
	cw.writeToFile("@%s\n", endLabel)
	cw.writeToFile("0;JMP\n")
	cw.writeToFile("(%s)\n", trueLabel)
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M-1\n")
	cw.writeToFile("M=-1\n")
	cw.writeToFile("(%s)\n", endLabel)
}

func (cw *CodeWriter) writeGt() {
	trueLabel := fmt.Sprintf("GT_TRUE_%s", GenerateUniqueID())
	endLabel := fmt.Sprintf("GT_END_%s", GenerateUniqueID())
	cw.writeToFile("@SP\n")
	cw.writeToFile("AM=M-1\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("A=A-1\n")
	cw.writeToFile("D=M-D\n")
	cw.writeToFile("@%s\n", trueLabel)
	cw.writeToFile("D;JGT\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M-1\n")
	cw.writeToFile("M=0\n")
	cw.writeToFile("@%s\n", endLabel)
	cw.writeToFile("0;JMP\n")
	cw.writeToFile("(%s)\n", trueLabel)
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M-1\n")
	cw.writeToFile("M=-1\n")
	cw.writeToFile("(%s)\n", endLabel)
}

func (cw *CodeWriter) writeLt() {
	trueLabel := fmt.Sprintf("LT_TRUE_%s", GenerateUniqueID())
	endLabel := fmt.Sprintf("LT_END_%s", GenerateUniqueID())
	cw.writeToFile("@SP\n")
	cw.writeToFile("AM=M-1\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("A=A-1\n")
	cw.writeToFile("D=M-D\n")
	cw.writeToFile("@%s\n", trueLabel)
	cw.writeToFile("D;JLT\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M-1\n")
	cw.writeToFile("M=0\n")
	cw.writeToFile("@%s\n", endLabel)
	cw.writeToFile("0;JMP\n")
	cw.writeToFile("(%s)\n", trueLabel)
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M-1\n")
	cw.writeToFile("M=-1\n")
	cw.writeToFile("(%s)\n", endLabel)
}

func (cw *CodeWriter) writeAnd() {
	cw.writeToFile("@SP\n")
	cw.writeToFile("AM=M-1\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("A=A-1\n")
	cw.writeToFile("M=M&D\n")
}

func (cw *CodeWriter) writeOr() {
	cw.writeToFile("@SP\n")
	cw.writeToFile("AM=M-1\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("A=A-1\n")
	cw.writeToFile("M=M|D\n")
}

func (cw *CodeWriter) writeNot() {
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M-1\n")
	cw.writeToFile("M=!M\n")
}

func (cw *CodeWriter) WritePushPop(command string, segment string, index int) {
	switch command {
	case "push":
		cw.writePush(segment, index)
	case "pop":
		cw.writePop(segment, index)
	}
}

func (cw *CodeWriter) writePush(segment string, index int) {
	switch segment {
	case "constant":
		cw.writePushConstant(index)
	case "local":
		cw.writePushSegment("LCL", index)
	case "argument":
		cw.writePushSegment("ARG", index)
	case "this":
		cw.writePushSegment("THIS", index)
	case "that":
		cw.writePushSegment("THAT", index)
	case "temp":
		cw.writePushTemp(index)
	case "pointer":
		cw.writePushPointer(index)
	case "static":
		cw.writePushStatic(index)
	}
}

func (cw *CodeWriter) writePushConstant(index int) {
	cw.writeToFile("@%d\n", index)
	cw.writeToFile("D=A\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M\n")
	cw.writeToFile("M=D\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("M=M+1\n")
}

func (cw *CodeWriter) writePushSegment(segment string, index int) {
	cw.writeToFile("@%s\n", segment)
	cw.writeToFile("D=M\n")
	cw.writeToFile("@%d\n", index)
	cw.writeToFile("A=D+A\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M\n")
	cw.writeToFile("M=D\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("M=M+1\n")
}

func (cw *CodeWriter) writePushTemp(index int) {
	cw.writeToFile("@%d\n", 5+index)
	cw.writeToFile("D=M\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M\n")
	cw.writeToFile("M=D\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("M=M+1\n")
}

func (cw *CodeWriter) writePushPointer(index int) {
	if index == 0 {
		cw.writeToFile("@THIS\n")
	} else if index == 1 {
		cw.writeToFile("@THAT\n")
	} else {
		panic("ポインターセグメントには0または1のインデックスしか指定できません")
	}
	cw.writeToFile("D=M\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M\n")
	cw.writeToFile("M=D\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("M=M+1\n")
}

func (cw *CodeWriter) writePushStatic(index int) {
	cw.writeToFile("@%s.%d\n", "STATIC", index)
	cw.writeToFile("D=M\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M\n")
	cw.writeToFile("M=D\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("M=M+1\n")
}

func (cw *CodeWriter) writePop(segment string, index int) {
	switch segment {
	case "constant":
		panic("constantセグメントに対してpopコマンドは実行できません")
	case "local":
		cw.writePopSegment("LCL", index)
	case "argument":
		cw.writePopSegment("ARG", index)
	case "this":
		cw.writePopSegment("THIS", index)
	case "that":
		cw.writePopSegment("THAT", index)
	case "temp":
		cw.writePopTemp(index)
	case "pointer":
		cw.writePopPointer(index)
	case "static":
		cw.writePopStatic(index)
	default:
		panic("未知のセグメントです")
	}
}

func (cw *CodeWriter) writePopSegment(segment string, index int) {
	cw.writeToFile("@%s\n", segment)
	cw.writeToFile("D=M\n")
	cw.writeToFile("@%d\n", index)
	cw.writeToFile("D=D+A\n")
	cw.writeToFile("@R13\n")
	cw.writeToFile("M=D\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("AM=M-1\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("@R13\n")
	cw.writeToFile("A=M\n")
	cw.writeToFile("M=D\n")
}

func (cw *CodeWriter) writePopTemp(index int) {
	cw.writeToFile("@SP\n")
	cw.writeToFile("AM=M-1\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("@%d\n", 5+index)
	cw.writeToFile("M=D\n")
}

func (cw *CodeWriter) writePopPointer(index int) {
	cw.writeToFile("@SP\n")
	cw.writeToFile("AM=M-1\n")
	cw.writeToFile("D=M\n")
	if index == 0 {
		cw.writeToFile("@THIS\n")
	} else if index == 1 {
		cw.writeToFile("@THAT\n")
	} else {
		panic("ポインターセグメントには0または1のインデックスしか指定できません")
	}
	cw.writeToFile("M=D\n")
	cw.writeToFile("@0\n")
	cw.writeToFile("D=A\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M\n")
	cw.writeToFile("M=D\n")
}

func (cw *CodeWriter) writePopStatic(index int) {
	cw.writeToFile("@SP\n")
	cw.writeToFile("AM=M-1\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("@%s.%d\n", "STATIC", index)
	cw.writeToFile("M=D\n")
}

func (cw *CodeWriter) writeToFile(format string, args ...interface{}) {
	if cw.file == nil {
		panic("ファイルがオープンされていません")
	}
	_, err := cw.file.WriteString(fmt.Sprintf(format, args...))
	if err != nil {
		panic(err)
	}
}

func GenerateUniqueID() string {
	id := uuid.New()
	return id.String()
}
