// Copyright 2022 The Serenity Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"

	"github.com/texttheater/golang-levenshtein/levenshtein"
)

const (
	// MemorySize is the size of the working memory
	MemorySize = 1024 * 1024
	// CyclesLimit is the limit on cycles
	CyclesLimit = 1024 * 1024
	// PopulationSize population size
	PopulationSize = 1024
)

var (
	// Genes are the genes
	Genes = [...]rune{'+', '-', '>', '<', '.', '[', ']'}
)

// Program is a program
// https://github.com/cvhariharan/goBrainFuck
type Program []rune

// Execute executes a program
func (p Program) Execute(size int) *strings.Builder {
	var (
		memory [MemorySize]int
		pc     int
		dc     int
		i      int
		output strings.Builder
	)
	length := len(p)

	for pc < length && i < CyclesLimit {
		opcode := p[pc]
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
			output.WriteRune(rune(memory[dc]))
			if len([]rune(output.String())) == size {
				return &output
			}
			pc++
		case ',':
			memory[dc] = p.input()
			pc++
		case '[':
			if memory[dc] == 0 {
				pc = p.findMatchingForward(pc) + 1
			} else {
				pc++
			}
		case ']':
			if memory[dc] != 0 {
				pc = p.findMatchingBackward(pc) + 1
			} else {
				pc++
			}
		default:
			pc++
		}
		i++
	}
	return &output
}

func (p Program) findMatchingForward(position int) int {
	count, length := 1, len(p)
	for i := position + 1; i < length; i++ {
		if p[i] == ']' {
			count--
			if count == 0 {
				return i
			}
		} else if p[i] == '[' {
			count++
		}
	}

	return length - 1
}

func (p Program) findMatchingBackward(position int) int {
	count := 1
	for i := position - 1; i >= 0; i-- {
		if p[i] == '[' {
			count--
			if count == 0 {
				return i
			}
		} else if p[i] == ']' {
			count++
		}
	}

	return -1
}

func (p Program) input() int {
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		log.Fatal(err)
	}
	return int(char)
}

func Generate(rnd *rand.Rand, program *strings.Builder) {
	count := rnd.Intn(16) + 1
	for i := 0; i < count; i++ {
		switch rnd.Intn(16) {
		case 0, 1, 2, 3:
			count := rnd.Intn(255) + 1
			for j := 0; j < count; j++ {
				program.WriteRune('+')
			}
		case 4, 5, 6, 7:
			count := rnd.Intn(255) + 1
			for j := 0; j < count; j++ {
				program.WriteRune('-')
			}
		case 8:
			program.WriteRune('>')
		case 9:
			program.WriteRune('<')
		case 10:
			program.WriteRune('.')
		case 11:
			program.WriteRune('[')
			Generate(rnd, program)
			program.WriteRune(']')
		}
	}
}

// Genome is a genome
type Genome struct {
	Program Program
	Output  string
	Fitness float64
}

// InsertGene inserts a gene into a genome
func (p Program) InsertGene(rnd *rand.Rand, gene rune, i, index int, child *strings.Builder) {
	for i < index {
		child.WriteRune(p[i])
		i++
	}
	if i == index {
		child.WriteRune(gene)
		length := len(p)
		for i < length {
			child.WriteRune(p[i])
			i++
		}
	}
}

// UpdateGene updates a gene in a genome
func (p Program) UpdateGene(rnd *rand.Rand, gene rune, i, index int, child *strings.Builder) {
	for i < index {
		child.WriteRune(p[i])
		i++
	}
	if i == index {
		child.WriteRune(gene)
		length := len(p)
		i++
		for i < length {
			child.WriteRune(p[i])
			i++
		}
	}
}

// DeleteGene deletes a gene from a genome
func (p Program) DeleteGene(rnd *rand.Rand, i, index int, child *strings.Builder) {
	for i < index {
		child.WriteRune(p[i])
		i++
	}
	if i == index {
		length := len(p)
		i++
		for i < length {
			child.WriteRune(p[i])
			i++
		}
	}
}

// Breed breeds two programs
func Breed(rnd *rand.Rand, a, b Program) (x, y Program) {
	lengtha, lengthb := len(a), len(b)
	a1, a2 := rnd.Intn(lengtha), rnd.Intn(lengtha)
	if a1 > a2 {
		a1, a2 = a2, a1
	}
	b1, b2 := rnd.Intn(lengthb), rnd.Intn(lengthb)
	if b1 > b2 {
		b1, b2 = b2, b1
	}

	x = a[:a1]
	x = append(x, b[b1:b2]...)
	x = append(x, a[a2:]...)

	y = b[:b1]
	y = append(y, a[a1:a2]...)
	y = append(y, b[b2:]...)

	return
}

func main() {
	rnd, target := rand.New(rand.NewSource(1)), []rune("ab")
	length := len(target)
	genomes := make([]Genome, 0, 8)
	for i := 0; i < PopulationSize; i++ {
		program := strings.Builder{}
		Generate(rnd, &program)
		code := Program(program.String())
		output := code.Execute(length)
		distance := levenshtein.DistanceForStrings([]rune(output.String()), target, levenshtein.DefaultOptions)
		genomes = append(genomes, Genome{
			Program: code,
			Output:  output.String(),
			Fitness: float64(distance),
		})
	}
	sort.Slice(genomes, func(i, j int) bool {
		return genomes[i].Fitness < genomes[j].Fitness
	})

	for i := 0; i < 128; i++ {
		size := len(genomes)
		for j := 0; j < size; j++ {
			// insert
			for _, gene := range Genes {
				index, child := rnd.Intn(len(genomes[j].Program)+1), strings.Builder{}
				genomes[j].Program.InsertGene(rnd, gene, 0, index, &child)
				code := Program([]rune(child.String()))
				output := code.Execute(length)
				distance := levenshtein.DistanceForStrings([]rune(output.String()), target, levenshtein.DefaultOptions)
				genomes = append(genomes, Genome{
					Program: code,
					Output:  output.String(),
					Fitness: float64(distance),
				})
			}

			// update
			for _, gene := range Genes {
				index, child := rnd.Intn(len(genomes[j].Program)+1), strings.Builder{}
				genomes[j].Program.UpdateGene(rnd, gene, 0, index, &child)
				code := Program([]rune(child.String()))
				output := code.Execute(length)
				distance := levenshtein.DistanceForStrings([]rune(output.String()), target, levenshtein.DefaultOptions)
				genomes = append(genomes, Genome{
					Program: code,
					Output:  output.String(),
					Fitness: float64(distance),
				})
			}

			// delete
			if len(genomes[j].Program) > 0 {
				index, child := rnd.Intn(len(genomes[j].Program)), strings.Builder{}
				genomes[j].Program.DeleteGene(rnd, 0, index, &child)
				code := Program([]rune(child.String()))
				output := code.Execute(length)
				distance := levenshtein.DistanceForStrings([]rune(output.String()), target, levenshtein.DefaultOptions)
				genomes = append(genomes, Genome{
					Program: code,
					Output:  output.String(),
					Fitness: float64(distance),
				})
			}
		}

		for j := 0; j < 10; j++ {
			a, b := rnd.Intn(10), rnd.Intn(10)
			x, y := Breed(rnd, genomes[a].Program, genomes[b].Program)

			output := x.Execute(length)
			distance := levenshtein.DistanceForStrings([]rune(output.String()), target, levenshtein.DefaultOptions)
			genomes = append(genomes, Genome{
				Program: x,
				Output:  output.String(),
				Fitness: float64(distance),
			})

			output = y.Execute(length)
			distance = levenshtein.DistanceForStrings([]rune(output.String()), target, levenshtein.DefaultOptions)
			genomes = append(genomes, Genome{
				Program: y,
				Output:  output.String(),
				Fitness: float64(distance),
			})
		}

		sort.Slice(genomes, func(i, j int) bool {
			return genomes[i].Fitness < genomes[j].Fitness
		})
		fmt.Println(i, genomes[0].Fitness)
		if genomes[0].Fitness == 0 {
			fmt.Println(genomes[0].Output)
			fmt.Println(string(genomes[0].Program))
			break
		}
		genomes = genomes[:PopulationSize]
	}
}
