package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Parser struct {
	file *os.File
	// lineCount         int
	currentLineNumber int
	currentCommand    string
}

func NewParser(
	file *os.File,
) *Parser {
	lineNumber := 1
	// コメントアウト、または空行をスキップしていき、最終的に実行可能なコマンド、行数を取得する
	command, line, err := getExecutableNextCommand(file, lineNumber)
	if err != nil {
		fmt.Println("failed to generate parser", err)
		return nil
	}

	return &Parser{
		file: file,
		// lineCount:         getLineLength(file),
		currentLineNumber: line,
		currentCommand:    command,
	}
}

/**
* 入力において、さらにコマンドが存在するかどうかを判定する
 */
func (p *Parser) HasMoreCommands() bool {
	_, lineNumber, err := getExecutableNextCommand(p.file, p.currentLineNumber)
	if err != nil {
		fmt.Println("failed to check if there are more commands", err)
	}
	return err == nil && lineNumber != 0
}

func getExecutableNextCommand(file *os.File, currentLineNumber int) (command string, lineNumber int, error error) {
	_, err := file.Seek(0, 0)
	if err != nil {
		return "", 0, err
	}
	scanner := bufio.NewScanner(file)
	lineNumber = 1
	for scanner.Scan() {
		if lineNumber <= currentLineNumber {
			lineNumber++
			continue
		}
		lineText := strings.TrimSpace(scanner.Text())
		if idx := strings.Index(lineText, "//"); idx != -1 {
			lineText = strings.TrimSpace(lineText[:idx])
		}

		if !isCommandExecutable(lineText) {
			lineNumber++
			continue
		}
		fmt.Println(lineText)
		return lineText, lineNumber, nil
	}
	return "", 0, nil
}

func isCommandExecutable(command string) bool {
	// trimmedCommand := strings.Trim(command, " ")
	if command == "" {
		return false
	}
	if command[0] == '/' && command[1] == '/' {
		return false
	}
	return true
}

/**
* 入力から次のコマンドを読み、それを現在のコマンドにする
* このルーチンは、HasMoreCommands()がtrueの場合のみ呼び出される
 */
func (p *Parser) Advance() {
	command, lineNumber, err := getExecutableNextCommand(p.file, p.currentLineNumber)
	if err != nil {
		fmt.Println(err)
		return
	}
	p.currentLineNumber = lineNumber
	p.currentCommand = command
}

/**
* 現コマンドの種類を返す
 */
func (p *Parser) CommandType() string {
	order := strings.Split(p.currentCommand, " ")[0]
	switch order {
	case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
		return "C_ARITHMETIC"
	case "push":
		return "C_PUSH"
	case "pop":
		return "C_POP"
	case "label":
		return "C_LABEL"
	case "goto":
		return "C_GOTO"
	case "if-goto":
		return "C_IF"
	case "function":
		return "C_FUNCTION"
	case "return":
		return "C_RETURN"
	case "call":
		return "C_CALL"
	default:
		return "C_UNKNOWN"
	}
}

/**
* 現コマンドの最初の引数を返す
 */
func (p *Parser) Arg1() string {
	if p.CommandType() == "C_RETURN" {
		return ""
	}
	commandList := strings.Split(p.currentCommand, " ")
	result := commandList[0]
	if len(commandList) > 1 {
		result = commandList[1]
	}
	return result
}

/**
* 現コマンドの2番目の引数を返す
 */
func (p *Parser) Arg2() int {
	if p.CommandType() != "C_PUSH" && p.CommandType() != "C_POP" && p.CommandType() != "C_FUNCTION" && p.CommandType() != "C_CALL" {
		return 0
	}
	commandList := strings.Split(p.currentCommand, " ")
	result, error := strconv.Atoi(commandList[2])
	if error != nil {
		fmt.Println("failed to convert string to int", error)
		return 0
	}

	return result
}
