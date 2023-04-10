package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	var str string
	flag.StringVar(&str, "str", "", "str")
	flag.Parse()

	tmp := strings.ReplaceAll(str, "\\\"", "\"")
	tmp2 := strings.ReplaceAll(tmp, "\\\\u", "\\u")

	strArr := strings.Split(tmp2, ":")
	var strNew strings.Builder
	for _, val := range strArr {
		subOpt := strings.Index(val, "\\u")
		if subOpt != -1 {
			subOpts := strings.ReplaceAll(val, "\"", "")
			isBlock := strings.Index(val, "}")
			if isBlock != -1 {
				subOpts = strings.ReplaceAll(subOpts, "}", "")
			}

			msgStr, err := zhToUnicode([]byte(subOpts))
			if err != nil {
				fmt.Println("err:", err)
			}
			if isBlock != -1 {
				val = "\"" + string(msgStr) + "\"}"
			} else {
				val = "\"" + string(msgStr) + "\""
			}
		}
		strNew.WriteString(":" + val)
	}

	cleanStr := strings.TrimPrefix(strNew.String(), ":")
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(cleanStr), "", "    "); err != nil {
		fmt.Println("err:", err)
	}
	fmt.Printf("\033[32;42;40m%s\033[0m\n", prettyJSON.String())
}

func zhToUnicode(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}
