package main

import (
	"encoding/base64"
	"flag"
	"fmt"
)

func main() {
	var input, style string
	flag.StringVar(&input, "input", "", "input")
	flag.StringVar(&style, "style", "", "input")
	flag.Parse()
	if style == "encode" {
		encode(input)
	} else {

		decode(input)
	}
}

func encode(data string) {
	str := base64.StdEncoding.EncodeToString([]byte(data))
	fmt.Println(str)
}

func decode(str string) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	// fmt.Printf("%q\n", data)
	fmt.Println(data)
}
