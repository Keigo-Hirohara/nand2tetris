@1
D=A
@ARG
M=M+D
A=M
D=M
@SP
A=M
M=D
@SP
M=M+1
@1
D=A
@ARG
M=M-D
@SP
AM=M-1
D=M
@THAT
M=D
@0
D=A
@SP
A=M
M=D
@0
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
A=M-1
D=M
@THAT
A=M
M=D
@0
D=A
@SP
M=M-1
A=M
M=D
@1
D=A
@SP
A=M
M=D
@SP
M=M+1
@1
D=A
@THAT
M=M+D
@SP
A=M-1
D=M
@THAT
A=M
M=D
@1
D=A
@THAT
M=M-D
@0
D=A
@SP
M=M-1
A=M
M=D
@ARG
A=M
D=M
@SP
A=M
M=D
@SP
M=M+1
@2
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
@SP
A=M-1
M=M-D
@0
D=A
@SP
A=M
M=D
@SP
A=M-1
D=M
@ARG
A=M
M=D
@0
D=A
@SP
M=M-1
A=M
M=D
(LOOP)
@ARG
A=M
D=M
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
@COMPUTE_ELEMENT
D;JNE
@END
0;JEQ
(COMPUTE_ELEMENT)
@THAT
A=M
D=M
@SP
A=M
M=D
@SP
M=M+1
@1
D=A
@THAT
M=M+D
A=M
D=M
@SP
A=M
M=D
@SP
M=M+1
@1
D=A
@THAT
M=M-D
@SP
AM=M-1
D=M
@SP
A=M-1
M=M+D
@0
D=A
@SP
A=M
M=D
@2
D=A
@THAT
M=M+D
@SP
A=M-1
D=M
@THAT
A=M
M=D
@2
D=A
@THAT
M=M-D
@0
D=A
@SP
M=M-1
A=M
M=D
@THAT
D=M
@SP
A=M
M=D
@SP
M=M+1
@1
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
@SP
A=M-1
M=M+D
@0
D=A
@SP
A=M
M=D
@SP
AM=M-1
D=M
@THAT
M=D
@0
D=A
@SP
A=M
M=D
@ARG
A=M
D=M
@SP
A=M
M=D
@SP
M=M+1
@1
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
@SP
A=M-1
M=M-D
@0
D=A
@SP
A=M
M=D
@SP
A=M-1
D=M
@ARG
A=M
M=D
@0
D=A
@SP
M=M-1
A=M
M=D
@LOOP
0;JEQ
(END)
