package main

import (
	"dojer/cmd"
	"dojer/store"
	"fmt"
)

func main() {
	data, _ := store.Search("in", 1)
	fmt.Printf("data: %v\n", data)
	cmd.Execute()
}
