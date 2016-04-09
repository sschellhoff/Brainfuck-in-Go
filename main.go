package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type Cell struct {
	prev  *Cell
	next  *Cell
	value int
}

type VirtualMachine struct {
	currentCell *Cell
	pc          int
	program     string
}

func NewVM() *VirtualMachine {
	v := new(VirtualMachine)
	v.pc = 0
	v.program = ""
	v.currentCell = &Cell{prev: nil, next: nil, value: 0}
	return v
}

func (v *VirtualMachine) run(program string) error {
	v.program = program
	v.pc = 0

	if !v.validateBraces() {
		return errors.New("invalid program")
	}

	for v.pc < len(v.program) {
		switch v.program[v.pc] {
		case '<':
			v.previousCell()
			break
		case '>':
			v.nextCell()
			break
		case '.':
			v.output()
			break
		case ',':
			v.input()
			break
		case '+':
			v.incMemory()
			break
		case '-':
			v.decMemory()
			break
		case '[':
			v.conditionalFwdJump()
			break
		case ']':
			v.conditionalBwdJump()
			break
		}
		v.pc++
	}
	return nil
}

func (v *VirtualMachine) validateBraces() bool {
	openBraces := 0
	for _, c := range v.program {
		if c == '[' {
			openBraces++
		} else if c == ']' {
			openBraces--
			if openBraces < 0 {
				return false
			}
		}
	}
	return openBraces == 0
}

func (v *VirtualMachine) incMemory() {
	v.currentCell.value++
}

func (v *VirtualMachine) decMemory() {
	v.currentCell.value--
}

func (v *VirtualMachine) nextCell() {
	if v.currentCell.next == nil {
		v.currentCell.next = &Cell{prev: v.currentCell, next: nil, value: 0}
		v.currentCell.next.prev = v.currentCell
	}
	v.currentCell = v.currentCell.next
}

func (v *VirtualMachine) previousCell() {
	if v.currentCell.prev == nil {
		v.currentCell.prev = &Cell{prev: nil, next: v.currentCell, value: 0}
		v.currentCell.prev.next = v.currentCell
	}
	v.currentCell = v.currentCell.prev
}

func (v *VirtualMachine) input() {
	v.currentCell.value = int(getch())
}

func (v *VirtualMachine) output() {
	fmt.Print(string(v.currentCell.value))
}

func (v *VirtualMachine) conditionalFwdJump() {
	if v.currentCell.value == 0 {
		counter := 1
		v.pc++
		for v.pc < len(v.program) && counter > 0 {
			if v.program[v.pc] == '[' {
				counter++
			} else if v.program[v.pc] == ']' {
				counter--
			}
			v.pc++
		}
	}
}

func (v *VirtualMachine) conditionalBwdJump() {
	if v.currentCell.value != 0 {
		counter := 1
		v.pc--
		for v.pc >= 0 && counter > 0 {
			if v.program[v.pc] == ']' {
				counter++
			} else if v.program[v.pc] == '[' {
				counter--
			}
			v.pc--
		}
	}
}

func main() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()

	program := "++++++++++[>+++++++>++++++++++>+++>+<<<<-]>++.>+.+++++++..+++.>++.<<+++++++++++++++.>.+++.------.--------.>+.>.+++."

	vm := NewVM()
	err := vm.run(program)
	if err != nil {
		fmt.Println(err)
	}
}

func getch() byte {
	var b []byte = make([]byte, 1)
	os.Stdin.Read(b)
	return b[0]
}
