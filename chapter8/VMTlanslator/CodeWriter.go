package main

import (
	"fmt"
	"os"

	"github.com/google/uuid"
)

type CodeWriter struct {
	file     *os.File
	count    int
	funcName string
	fileName string
}

func NewCodeWriter() *CodeWriter {
	return &CodeWriter{
		file:     nil,
		count:    0,
		funcName: "",
		fileName: "",
	}
}

func (cw *CodeWriter) SetVMFileName(fileName string) {
	cw.fileName = fileName
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

func (cw *CodeWriter) WriteIf(label string) {
	cw.writeToFile("@SP\n")
	cw.writeToFile("AM=M-1\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("@%s$%s\n", cw.funcName, label)
	cw.writeToFile("D;JNE\n")
}

func (cw *CodeWriter) WriteLabel(label string) {
	cw.writeToFile(fmt.Sprintf("(%s$%s)\n", cw.funcName, label))
}

func (cw *CodeWriter) WriteGoto(label string) {
	cw.writeToFile("@%s$%s\n", cw.funcName, label)
	cw.writeToFile("0;JMP\n")
}

func (cw *CodeWriter) WriteCall(funcName string, numberOfInput int) {
	count := cw.nextCount()
	// ?
	cw.writeToFile("@SP\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("@R13\n")
	cw.writeToFile("M=D\n")

	// push リターンアドレス
	cw.writeToFile("@ret.%s\n", count)
	cw.writeToFile("D=A\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M\n")
	cw.writeToFile("M=D\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("M=M+1\n")

	// push LCL
	cw.writeToFile("@LCL\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M\n")
	cw.writeToFile("M=D\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("M=M+1\n")

	// push ARG
	cw.writeToFile("@ARG\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M\n")
	cw.writeToFile("M=D\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("M=M+1\n")

	// push THIS
	cw.writeToFile("@THIS\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M\n")
	cw.writeToFile("M=D\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("M=M+1\n")

	// push THAT
	cw.writeToFile("@THAT\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M\n")
	cw.writeToFile("M=D\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("M=M+1\n")

	// ?
	cw.writeToFile("@R13\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("@%d\n", numberOfInput)
	cw.writeToFile("D=D-A\n")
	cw.writeToFile("@ARG\n")
	cw.writeToFile("M=D\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("@LCL\n")
	cw.writeToFile("M=D\n")
	cw.writeToFile("@%s\n", funcName)
	cw.writeToFile("0;JMP\n")
	cw.writeToFile("(ret.%s)\n", count)
}

func (cw *CodeWriter) WriteFunction(name string, numberOfLocal int) {
	cw.funcName = name
	// 関数の開始位置のラベルを宣言する
	cw.writeToFile("(%s)\n", name)

	// ローカル変数の数分スタックに0をpushする
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M\n")
	for i := 0; i < numberOfLocal; i++ {
		cw.writeToFile("M=0\n")
		cw.writeToFile("A=A+1\n")
	}
	cw.writeToFile("D=A\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("M=D\n")
}

func (cw *CodeWriter) WriteReturn() {
	// 関数呼び出し時に保存したリターンアドレス(LCL-5)を取得し、R13に格納する
	cw.writeToFile("@LCL\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("@5\n")
	cw.writeToFile("A=D-A\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("@R13\n")
	cw.writeToFile("M=D\n")

	// 関数の結果をARGに格納する
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M-1\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("@ARG\n")
	cw.writeToFile("A=M\n")
	cw.writeToFile("M=D\n")
	cw.writeToFile("D=A+1\n")

	// 呼び出し側のスタックポインタを戻す
	cw.writeToFile("@SP\n")
	cw.writeToFile("M=D\n")

	// 呼び出し側のTHAT, THIS, ARG, LCLを戻す
	cw.writeToFile("@LCL\n")
	cw.writeToFile("AM=M-1\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("@THAT\n")
	cw.writeToFile("M=D\n")

	cw.writeToFile("@LCL\n")
	cw.writeToFile("AM=M-1\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("@THIS\n")
	cw.writeToFile("M=D\n")

	cw.writeToFile("@LCL\n")
	cw.writeToFile("AM=M-1\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("@ARG\n")
	cw.writeToFile("M=D\n")

	cw.writeToFile("@LCL\n")
	cw.writeToFile("A=M-1\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("@LCL\n")
	cw.writeToFile("M=D\n")

	// 呼び出し側のリターンアドレスへジャンプする
	cw.writeToFile("@R13\n")
	cw.writeToFile("A=M\n")
	cw.writeToFile("0;JMP\n")
}

func (cw *CodeWriter) writeAdd() {
	cw.writeToFile("@SP\n")
	cw.writeToFile("AM=M-1\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M-1\n")
	cw.writeToFile("M=M+D\n")
}

func (cw *CodeWriter) writeSub() {
	cw.writeToFile("@SP\n")
	cw.writeToFile("AM=M-1\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("A=M-1\n")
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
	if index > 0 {
		// インデックスの値をDレジスタに格納し、シンボルの値に加算する
		cw.writeToFile("@%d\n", index)
		cw.writeToFile("D=A\n")
		cw.writeToFile("@%s\n", segment)
		cw.writeToFile("M=M+D\n")
		cw.writeToFile("A=M\n")
		cw.writeToFile("D=M\n")

		// スタックにDレジスタの値を追加し、スタックポインタをインクリメント
		cw.writeToFile("@SP\n")
		cw.writeToFile("A=M\n")
		cw.writeToFile("M=D\n")
		cw.writeToFile("@SP\n")
		cw.writeToFile("M=M+1\n")

		// インデックスの分加算したシンボルの値を減算する
		cw.writeToFile("@%d\n", index)
		cw.writeToFile("D=A\n")
		cw.writeToFile("@%s\n", segment)
		cw.writeToFile("M=M-D\n")
	} else {
		// セグメントが表すアドレスにおける値をDレジスタに格納
		cw.writeToFile("@%s\n", segment)
		cw.writeToFile("A=M\n")
		cw.writeToFile("D=M\n")

		// スタックにDレジスタの値を追加し、スタックポインタをインクリメント
		cw.writeToFile("@SP\n")
		cw.writeToFile("A=M\n")
		cw.writeToFile("M=D\n")
		cw.writeToFile("@SP\n")
		cw.writeToFile("M=M+1\n")
	}
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
	cw.writeToFile("@%s.%d\n", cw.fileName, index)
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
	if index > 0 {
		// インデックス分シンボルの値を加算する
		cw.writeToFile("@%d\n", index)
		cw.writeToFile("D=A\n")
		cw.writeToFile("@%s\n", segment)
		cw.writeToFile("M=M+D\n")

		// スタックの最下位の値をセグメントの値が示すアドレスの値に格納する
		cw.writeToFile("@SP\n")
		cw.writeToFile("A=M-1\n")
		cw.writeToFile("D=M\n")
		cw.writeToFile("@%s\n", segment)
		cw.writeToFile("A=M\n")
		cw.writeToFile("M=D\n")

		// インデックスの分加算したシンボルの値を減算する
		cw.writeToFile("@%d\n", index)
		cw.writeToFile("D=A\n")
		cw.writeToFile("@%s\n", segment)
		cw.writeToFile("M=M-D\n")
	} else {
		// スタックの最下位の値をDレジスタに格納する
		cw.writeToFile("@SP\n")
		cw.writeToFile("A=M-1\n")
		cw.writeToFile("D=M\n")

		cw.writeToFile("@%s\n", segment)
		cw.writeToFile("A=M\n")
		cw.writeToFile("M=D\n")
	}

	// スタックポインタの最下位の値を0に直す
	cw.writeToFile("@0\n")
	cw.writeToFile("D=A\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("M=M-1\n")
	cw.writeToFile("A=M\n")
	cw.writeToFile("M=D\n")
}

func (cw *CodeWriter) writePopTemp(index int) {
	cw.writeToFile("@R5\n")
	cw.writeToFile("D=A\n")
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
	cw.writeToFile("@%s.%d\n", cw.fileName, index)
	cw.writeToFile("D=A\n")
	cw.writeToFile("@R13\n")
	cw.writeToFile("M=D\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("AM=M-1\n")
	cw.writeToFile("D=M\n")
	cw.writeToFile("@R13\n")
	cw.writeToFile("A=M\n")
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

func (cw *CodeWriter) nextCount() string {
	cw.count++
	return fmt.Sprintf("%d", cw.count)

}

func (cw *CodeWriter) WriteBootstrap() {
	cw.writeToFile("@256\n")
	cw.writeToFile("D=A\n")
	cw.writeToFile("@SP\n")
	cw.writeToFile("M=D\n")
	cw.WriteCall("Sys.init", 0)
	cw.writeToFile("0;JMP\n")
}
