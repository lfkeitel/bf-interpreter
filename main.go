package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	memorySize = 30000
)

var (
	inputFile string
	debug     bool
)

func init() {
	flag.StringVar(&inputFile, "in", "-", "Input source file")
	flag.BoolVar(&debug, "d", false, "Enable debug output")
}

func main() {
	flag.Parse()

	var code []byte

	if inputFile == "-" {
		code = []byte(flag.Arg(0))
	} else {
		file, err := ioutil.ReadFile(inputFile)
		if err != nil {
			log.Fatal(err.Error())
		}
		code = file
	}

	var err error
	code, err = cleanCode(code)
	if err != nil {
		log.Fatal(err)
	}
	i := newInterpreter(code)
	i.execute()
}

// cleanCode removes all non-instrucitons characters for the interpreter
func cleanCode(input []byte) ([]byte, error) {
	output := make([]byte, 0, len(input)/2)
	braceCheck := 0

	for _, b := range input {
		switch b {
		case '>':
			fallthrough
		case '<':
			fallthrough
		case '+':
			fallthrough
		case '-':
			fallthrough
		case '.':
			fallthrough
		case ',':
			output = append(output, b)
		case '[':
			braceCheck++
			output = append(output, b)
		case ']':
			braceCheck--
			output = append(output, b)
		}
	}

	if braceCheck != 0 {
		return nil, errors.New("Unbalanced braces")
	}

	return output, nil
}

type interpreter struct {
	code   []byte
	memory []byte
	pc     int
	dp     int
}

func newInterpreter(code []byte) *interpreter {
	return &interpreter{
		code:   code,
		memory: make([]byte, memorySize),
	}
}

func (i *interpreter) execute() {
	for i.pc < len(i.code) {
		if debug {
			i.printStatus()
		}

		switch i.code[i.pc] {
		case '>':
			i.dp++
		case '<':
			i.dp--
		case '+':
			i.memory[i.dp]++
		case '-':
			i.memory[i.dp]--
		case '.':
			os.Stdout.Write([]byte{i.memory[i.dp]})
		case ',':
			in := make([]byte, 1)
			os.Stdin.Read(in)
			i.memory[i.dp] = in[0]
		case '[':
			if i.memory[i.dp] == 0 {
				i.pc = i.findMatchingEndBrace()
			}
		case ']':
			if i.memory[i.dp] != 0 {
				i.pc = i.findMatchingStartBrace()
			}
		}
		i.pc++
	}
}

func (i *interpreter) printStatus() {
	fmt.Printf("PC: %d\t\tDP: %d\t\tInstruction: %c\t\tMemory: %d\n", i.pc, i.dp, i.code[i.pc], i.memory[i.dp])
}

func (i *interpreter) findMatchingEndBrace() int {
	stack := 0
	pc := i.pc

	for pc < len(i.code) {
		switch i.code[pc] {
		case '[':
			stack++
		case ']':
			stack--
		}

		if stack == 0 {
			return pc
		}
		pc++
	}

	panic("Unbalanced braces")
}

func (i *interpreter) findMatchingStartBrace() int {
	stack := 0
	pc := i.pc

	for pc >= 0 {
		switch i.code[pc] {
		case '[':
			stack--
		case ']':
			stack++
		}

		if stack == 0 {
			return pc - 1
		}
		pc--
	}

	panic("Unbalanced braces")
}
