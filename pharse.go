package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	tag string
	text []byte
	children []Node
}

type Stack []Node

func (stack Stack) push(node Node) Stack {
	return append(stack, node)
}

func (stack Stack) pop() (node Node, res Stack) {
	if stack.isEmpty() {
		return
	}
	node = stack[len(stack)-1]
	res = stack[:len(stack)-1]
	return
}

func (stack Stack) peek() (node Node) {
	if stack.isEmpty() {
		return
	}
	node = stack[len(stack)-1]
	return
}

func (stack Stack) isEmpty() bool {
	return len(stack) < 1
}

func parse(reader *bufio.Reader) Node {
	var (
		b byte
		err error

		state string
		tag []byte

		dom Stack

		cur Node
		par Node
	)

	for err == nil {
		b, err = reader.ReadByte()
		switch state {
		default:
			switch b {
			default:
			case '<':
				state = "tag.open.start"
			}
		case "tag.open":
			switch b {
			default:
				tag = append(tag, b)
			case '>':
				cur.tag = string(tag)
				par.children = append(par.children, cur)
				dom = dom.push(cur)
				par = cur
				fmt.Println(ident(len(dom)), cur.tag, string(cur.text))
				tag = make([]byte, 10)
				state = ""
			}
		case "tag.close":
			switch b {
			default:
				tag = append(tag, b)
			case '>':
				state = ""
				cur, dom = dom.pop()
				tag = make([]byte, 10)
			}
		case "tag.open.start":
			switch b {
			default:
				tag = append(tag, b)
				state = "tag.open"
			case '/':
				state = "tag.close"
			}
		}

	}
	return cur
}

func ident(n int) (s string) {
	for i:=1; i < n; i++ {
		s += "  "
	}
	return
}

func main() {
	var (
		file   *os.File
		reader *bufio.Reader
		doc Node
	)

	if len(os.Args) > 1 {
		file, _ = os.Open(os.Args[1])
		reader = bufio.NewReader(file)
	} else {
		reader = bufio.NewReader(os.Stdin)
	}

	doc = parse(reader)
	fmt.Println(doc.tag)
}
