package main

import "testing"

func TestBasicPaths(t *testing.T) {
	segments := getSegments("/foo/bar")
	if len(segments) != 2 {
		t.Fatalf("Expected 2 segments, got %d", len(segments))
	}
	if segments[0] != "foo" {
		t.Fatalf("Unexpected segment value %s", segments[0])
	}
	if segments[1] != "bar" {
		t.Fatalf("Unexpected segment value %s", segments[1])
	}
}

func TestMkdir(t *testing.T) {
	f := NewFilesystem()
	node, err := f.Mkdir("/foo/bar")
	if err != nil {
		t.Fatal(err)
	}
	if node.Name != "bar" {
		t.Fatal(err)
	}
	if node.Parent == nil || node.Parent.Name != "foo" {
		t.Fatal(err)
	}
	if node.Parent.Parent != f.Root {
		t.Fatal(err)
	}
	if node.Type != Folder {
		t.Fatalf("Unexpected type: %s", node.Type)
	}
	// errors to handle:
	// NotTraversable, ChildAlreadyExists
	_, err = f.Mkdir("/foo/bar")
	if err != AlreadyExists {
		t.Fatal(err)
	}

	// no empty names
	_, err = f.Touch("//")
	if err != InvalidName {
		t.Fatal(err)
	}
}

func TestTouch(t *testing.T) {
	f := NewFilesystem()
	node, err := f.Touch("/foo/bar")
	if err != nil {
		t.Fatal(err)
	}
	if node.Name != "bar" {
		t.Fatal(err)
	}
	if node.Parent == nil || node.Parent.Name != "foo" {
		t.Fatal(err)
	}
	if node.Parent.Parent != f.Root {
		t.Fatal(err)
	}
	if node.Type != File {
		t.Fatalf("Unexpected type for the touched file: %s", node.Type)
	}
	if node.Parent.Type != Folder {
		t.Fatalf("Unexpected type for parent: %s", node.Parent.Type)
	}
	// errors to handle:
	// NotTraversable, ChildAlreadyExists
	_, err = f.Touch("/foo/bar")
	if err != AlreadyExists {
		t.Fatal(err)
	}

	nodes, err := f.List("/foo")
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 {
		t.Fatal("Unexpectedly wrong number of nodes: %d", len(nodes))
	}
	if nodes[0].Type != File || nodes[0].Name != "bar" {
		t.Fatalf("Unexpected result listing foo: %+v", nodes[0])
	}

	_, err = f.Touch("/foo/bar/baz")
	if err != NotTraversableError {
		t.Fatal(err)
	}
}

func TestList(t *testing.T) {
	f := NewFilesystem()
	_, err := f.Touch("/foo/bar")
	if err != nil {
		t.Fatal(err)
	}
	listedFiles, err := f.List("")
	if err != nil {
		t.Fatal(err)
	}
	if len(listedFiles) != 1 {
		t.Fatalf("Unexpected file list length %d", len(listedFiles))
	}
	if listedFiles[0].Name != "foo" {
		t.Fatalf("Unexpected file %+v", listedFiles[0])
	}
}

func TestSegments(t *testing.T) {
	segments := getSegments("//")
	if len(segments) != 2 {
		t.Fatal(len(segments))
	}
	if segments[0] != "" || segments[1] != "" {
		t.Fatal("expected empty strings")
	}
}
