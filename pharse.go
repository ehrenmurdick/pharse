package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	tag      string
	text     []byte
	parent   *Node
	children []*Node
}

type Stack []*Node

func (stack Stack) push(node *Node) Stack {
	return append(stack, node)
}

func (stack Stack) pop() (node *Node, res Stack) {
	if stack.isEmpty() {
		return
	}
	node = stack[len(stack)-1]
	res = stack[:len(stack)-1]
	return
}

func (stack Stack) peek() (node *Node) {
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
		b   byte
		err error

		state string
		tag   []byte

		dom Stack

		cur   *Node
		child *Node
	)

	cur = new(Node)
	cur.tag = "document"
	child = new(Node)

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
				child = new(Node)
				child.tag = string(tag)
				child.parent = cur
				cur.children = append(cur.children, child)
				dom = dom.push(cur)
				cur = child
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
	return *cur
}

func ident(n int) (s string) {
	for i := 1; i < n; i++ {
		s += "  "
	}
	return
}

func (doc Node) walk(i int) {
	for _, n := range doc.children {
		fmt.Println(ident(i), n.tag, len(n.children))
		n.walk(i + 1)
	}
}

func main() {
	var (
		file   *os.File
		reader *bufio.Reader
		doc    Node
	)

	if len(os.Args) > 1 {
		file, _ = os.Open(os.Args[1])
		reader = bufio.NewReader(file)
	} else {
		reader = bufio.NewReader(os.Stdin)
	}

	doc = parse(reader)
	doc.walk(1)
}
