package main

import (
	"fgengine/language"
	"fmt"
)

func main() {
	text, err := language.ImportYAML("./ptbr.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Println(*text)
}
