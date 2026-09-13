package main

import (
	"fgengine/language"
	"fmt"
)

func main() {
	text, err := language.LoadLang(language.English)
	if err != nil {
		panic(err)
	}
	fmt.Println(*text)
}
