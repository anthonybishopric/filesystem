package main

import (
	"errors"
	"strings"
)

type NodeType int

const (
	File NodeType = iota
	Folder
	Symlink
)

type Segment string

type Node struct {
	Name     string
	Children map[Segment]*Node
	Parent   *Node
	Type     NodeType
}

func (n *Node) Get(segment Segment) (*Node, bool) {
	node, ok := n.Children[segment]
	return node, ok
}

type Filesystem struct {
	Root *Node
}

func NewFilesystem() *Filesystem {
	return &Filesystem{
		Root: &Node{
			Name:     "__",
			Children: make(map[Segment]*Node),
			Type:     Folder,
		},
	}
}

var NotTraversableError = errors.New("Not traversable")
var AlreadyExists = errors.New("File already exists")
var NoPathGiven = errors.New("No path given")
var InvalidName = errors.New("Invalid name given")

func (n *Node) addChild(segment Segment, theType NodeType) (*Node, error) {
	child, exists := n.Get(segment)
	if exists {
		return child, AlreadyExists
	}
	if len(string(segment)) == 0 {
		return nil, InvalidName
	}
	newNode := &Node{
		Type:     theType,
		Parent:   n,
		Name:     string(segment),
		Children: make(map[Segment]*Node),
	}
	n.Children[segment] = newNode
	return newNode, nil
}

func (n *Node) Mkdir(segment Segment) (*Node, error) {
	if n.Type != Folder {
		return nil, NotTraversableError
	}
	return n.addChild(segment, Folder)
}

func (n *Node) Touch(segment Segment) (*Node, error) {
	if n.Type != Folder {
		return nil, NotTraversableError
	}
	return n.addChild(segment, File)
}

func (f *Filesystem) Mkdir(path string) (*Node, error) {
	segments := getSegments(path) // [foo, bar]
	return f.mkdirRecursive(segments)
}

func (f *Filesystem) mkdirRecursive(segments []Segment) (*Node, error) {
	active := f.Root
	for index, segment := range segments {
		child, err := active.Mkdir(segment)
		if err == AlreadyExists && index == len(segments)-1 {
			return child, err
		}
		// ok, continue
		active = child
	}
	return active, nil
}

func (f *Filesystem) List(path string) ([]*Node, error) {
	segments := getSegments(path)
	active := f.Root
	for _, segment := range segments {
		child, ok := active.Get(segment)
		if !ok {
			return nil, NotTraversableError
		}
		active = child
	}
	nodes := []*Node{}
	for _, v := range active.Children {
		nodes = append(nodes, v)
	}
	return nodes, nil
}

func (f *Filesystem) Touch(path string) (*Node, error) {
	touchSegments := getSegments(path)
	if len(touchSegments) == 0 {
		return nil, NoPathGiven
	}
	parentDir, err := f.mkdirRecursive(touchSegments[:len(touchSegments)-1])
	if err != nil && err != AlreadyExists {
		return nil, err
	}
	// verify that type == Folder here
	return parentDir.Touch(touchSegments[len(touchSegments)-1])
}

func getSegments(path string) []Segment {
	path = strings.TrimSpace(path)
	segments := []Segment{}
	for _, str := range strings.Split(path, "/")[1:] {
		segments = append(segments, Segment(str))
	}
	return segments
}
