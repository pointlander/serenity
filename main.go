// Copyright 2022 The Serenity Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	// MemorySize is the size of the working memory
	MemorySize = 1024 * 1024
)

// https://github.com/cvhariharan/goBrainFuck
type cpu struct {
	Program []rune
	Output  strings.Builder
}

// NewCPU produces a new cpu
func NewCPU(code string) *cpu {
	c := &cpu{}
	c.Program = []rune(code)
	return c
}

func (c *cpu) execute() {
	var (
		memory [MemorySize]int
		pc     int
		dc     int
	)
	program := c.Program
	length := len(program)

	for pc < length {
		opcode := program[pc]
		switch opcode {
		case '+':
			memory[dc] += 1
			pc++
		case '-':
			memory[dc] -= 1
			pc++
		case '>':
			dc++
			pc++
		case '<':
			if dc > 0 {
				dc--
			}
			pc++
		case '.':
			c.Output.WriteRune(rune(memory[dc]))
			pc++
		case ',':
			memory[dc] = c.input()
			pc++
		case '[':
			if memory[dc] == 0 {
				pc = c.findMatchingForward(pc) + 1
			} else {
				pc++
			}
		case ']':
			if memory[dc] != 0 {
				pc = c.findMatchingBackward(pc) + 1
			} else {
				pc++
			}
		default:
			pc++
		}
	}
}

func (c *cpu) findMatchingForward(position int) int {
	program, count := c.Program, 1
	length := len(program)
	for i := position + 1; i < length; i++ {
		if program[i] == ']' {
			count--
			if count == 0 {
				return i
			}
		} else if program[i] == '[' {
			count++
		}
	}

	return -1
}

func (c *cpu) findMatchingBackward(position int) int {
	program, count := c.Program, 1
	for i := position - 1; i >= 0; i-- {
		if program[i] == '[' {
			count--
			if count == 0 {
				return i
			}
		} else if program[i] == ']' {
			count++
		}
	}

	return -1
}

func (c *cpu) input() int {
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		log.Fatal(err)
	}
	return int(char)
}

func main() {
	machine := NewCPU(`>++++++++[-<+++++++++>]<.>>+>-[+]++>++>+++[>[->+++<<+++>]<<]>-----.>-> Comments can be added
+++..+++.>-.<<+[>[+>+]>>]<--------------.>>.+++.------.--------.>+.>+.`)
	machine.execute()
	fmt.Print(machine.Output.String())
}
