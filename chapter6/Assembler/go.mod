module Assembler

go 1.20

require (
	Code v0.0.0
	Parser v0.0.0
	SymbolTable v0.0.0
)

replace (
	Code => ./Code
	Parser => ./Parser
	SymbolTable => ./SymbolTable
)
