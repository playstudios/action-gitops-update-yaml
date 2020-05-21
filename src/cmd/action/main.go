package main

import (
	actions "github.com/sethvargo/go-githubactions"
)

func main() {
	val := actions.GetInput("val")
	if val == "" {
		actions.Fatalf("missing 'val'")
	}
}
