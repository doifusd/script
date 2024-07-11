package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func main() {

	res := strings.ReplaceAll(conent, "\\\\u", "\\u")
	str := strings.ReplaceAll(res, "\\\"", "\"")

	str2, errs2 := zh2Unicode([]byte(str))
	if errs2 != nil {
		fmt.Println("to chinses err:", errs2)
	}
	isJSON := json.Valid([]byte(str2))
	fmt.Println("isJSON:", isJSON)
	fmt.Println("str2:", str2)
	//if isJSON {
	//	waitNum <- 1
	//	//todo go routine
	//	// res := httpClient(str2)
	//	// wg.Add(1)
	//	go httpClient(waitNum, str2, i)
	//	// wg.Wait()
	//	time.Sleep(time.Microsecond * 5000)
	//}
}

func zh2Unicode(raw []byte) (str string, err error) {
	str, err = strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return
	}
	return
}
