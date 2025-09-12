package main

import "fmt"

type testAction struct {
}

func (a *testAction) Execute(args ...any) error {
	fmt.Printf("LOADED\n")
	return nil
}

var Action = testAction{}
