package main

import (
	"Code"
	"Parser"
	"SymbolTable"
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	assemblyFile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("input Error: ", err)
		panic(err)
	}
	defer assemblyFile.Close()

	binaryFile, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Println("output Error: ", err)
		panic(err)
	}
	defer binaryFile.Close()

	reader := bufio.NewReader(assemblyFile)

	symbolTable := SymbolTable.NewSymbolTable()

	lastSymbolAddress := 16

	preParseLineNumber := 0

	for {
		rawLine, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		line := strings.TrimSpace(rawLine)

		if len(line) == 0 {
			continue
		}
		if line[0] == '/' && line[1] == '/' {
			continue
		}

		parser := Parser.NewParser(line)

		if parser.CommandType() == "L_COMMAND" {
			if !symbolTable.Contains(parser.Symbol()) {
				symbolTable.AddEntry(parser.Symbol(), preParseLineNumber)
			}
			continue
		}

		preParseLineNumber++
	}

	_, err = assemblyFile.Seek(0, io.SeekStart)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	reader = bufio.NewReader(assemblyFile)

	for {
		rawLine, err := reader.ReadString('\n')

		line := strings.TrimSpace(rawLine)
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if len(line) == 0 {
			continue
		}

		if line[0] == '/' && line[1] == '/' {
			continue
		}

		var binaryCommand string = ""

		parser := Parser.NewParser(line)

		if parser.CommandType() == "L_COMMAND" {
			// if !symbolTable.Contains(parser.Symbol()) {
			// 	lastSymbolAddress++
			// 	symbolTable.AddEntry(parser.Symbol(), symbolTable.GetAddress(strconv.Itoa(lastSymbolAddress)))
			// }
			continue
		}
		// fmt.Println(line)

		if parser.CommandType() == "A_COMMAND" {
			// fmt.Println(parser.Symbol())
			isNumber := false
			_, err := strconv.Atoi(parser.Symbol())
			if err == nil {
				isNumber = true
			}
			if isNumber {
				binaryCommand = "0" + Code.SymbolBinary(parser.Symbol())
			} else {
				if !symbolTable.Contains(parser.Symbol()) && !isNumber {
					symbolTable.AddEntry(parser.Symbol(), lastSymbolAddress)
					lastSymbolAddress++
				}
				binaryCommand = "0" + Code.SymbolBinary(strconv.Itoa(symbolTable.GetAddress(parser.Symbol())))
			}

			// if !symbolTable.Contains(parser.Symbol()) {
			// 	lastSymbolAddress++
			// 	symbolTable.AddEntry(parser.Symbol(), symbolTable.GetAddress(strconv.Itoa(lastSymbolAddress)))
			// }

			// binaryCommand = "0" + Code.SymbolBinary(strconv.Itoa(symbolTable.GetAddress(parser.Symbol())))
		}

		if parser.CommandType() == "C_COMMAND" {
			binaryCommand = "111" + Code.Comp(parser.Comp()) + Code.Dest(parser.Dest()) + Code.Jump(parser.Jump())
		}

		// fmt.Println(binaryCommand)

		binaryFile.WriteString(binaryCommand + "\n")
	}
}
