package Parser

import (
	"strings"
)

/**
* 入力に対して構文解析を行う
 */
type Parser struct {
	line string
}

// TODO: advanceを実装する

func NewParser(line string) *Parser {
	return &Parser{line: line}
}

/** 入力にまだコマンドが存在するかを判定する */
func (p *Parser) HasMoreCommands() bool {
	return len(p.line) != 0
}

/**
* 現コマンドの種類を返す
* TODO: 引数をA, C, L_Commandに絞り込む
 */
func (p *Parser) CommandType() string {
	if p.line[0] == '@' {
		return "A_COMMAND"
	}
	if p.line[0] == '(' {
		return "L_COMMAND"
	}
	return "C_COMMAND"
}

/**
* 現A命令またはC命令のシンボルまたは数値を返す
* commandTypeがA_COMMANDまたはL_COMMANDの時だけ呼び出す
 */
func (p *Parser) Symbol() string {
	if p.CommandType() == "A_COMMAND" {
		return p.line[1:]
	}
	if p.CommandType() == "L_COMMAND" {
		return p.line[1 : len(p.line)-1]
	}
	return ""
}

/**
* 現C命令のdestニーモニックを返す
* commandTypeがC_COMMANDの時だけ呼び出す
 */
func (p *Parser) Dest() string {
	if p.CommandType() == "A_COMMAND" || p.CommandType() == "L_COMMAND" {
		return ""
	}
	var formula string = ""
	if strings.Contains(p.line, "=") {
		formula = strings.Split(p.line, "=")[0]
	}
	return formula
}

/**
* 現C命令のcompニーモニックを返す
* commandTypeがC_COMMANDの時だけ呼び出す
 */
func (p *Parser) Comp() string {
	if p.CommandType() == "A_COMMAND" || p.CommandType() == "L_COMMAND" {
		return ""
	}
	var formula string = ""
	if strings.Contains(p.line, "=") {
		formula = strings.Split(p.line, "=")[1]
	}
	if strings.Contains(p.line, ";") {
		formula = strings.Split(p.line, ";")[0]
	}
	return formula
}

/**
* 現C命令のjumpニーモニックを返す
* commandTypeがC_COMMANDの時だけ呼び出す
 */
func (p *Parser) Jump() string {
	if p.CommandType() == "A_COMMAND" || p.CommandType() == "L_COMMAND" {
		return ""
	}
	var formula string = ""
	if strings.Contains(p.line, ";") {
		formula = strings.Split(p.line, ";")[1]
	}
	return formula
}

// func Parser(line string) string {

// 	return line
// }
