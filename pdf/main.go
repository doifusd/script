package main

import (
	"flag"
	"fmt"
	"os/exec"
)

func main() {
	Pdf2Word()
}

func Pdf2Word() {
	var pdfPath, wordPath string
	flag.StringVar(&pdfPath, "pdfPath", "", "pdfPath")
	flag.StringVar(&wordPath, "wordPath", "", "wordPath")
	flag.Parse()

	// 使用Unoconv将PDF转换为Word文档
	cmd := exec.Command("unoconv", "-f", "docx", "-o", wordPath, pdfPath)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error converting PDF to Word:", err)
	} else {
		fmt.Println("PDF converted to Word successfully!")
	}
}

func Word2Pdf() {}
