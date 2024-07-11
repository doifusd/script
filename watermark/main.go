package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/golang/freetype"
)

func main() {
	imageSrc := "./123.jpeg"
	imgb, err := os.Open(imageSrc)
	if err != nil {
		panic(err)
	}
	defer imgb.Close()

	var img image.Image
	img, _, err = image.Decode(imgb)
	if err != nil {
		fmt.Println("get_mi_picture err", err)
		return
	}
	//读取图片
	offset := image.Pt(img.Bounds().Dx()-210, img.Bounds().Dy()-10)
	b := img.Bounds()
	m := image.NewNRGBA(b) //按原图生成新图

	//文字水印
	//fontBytes, err1 := ioutil.ReadFile("Arial Unicode.ttf") //读取字体文件
	fontBytes, err1 := ioutil.ReadFile("./Arial.ttf") //读取字体文件
	if err1 != nil {
		log.Println(err1)
	}

	font, err2 := freetype.ParseFont(fontBytes)
	if err2 != nil {
		log.Println(err2)
	}

	f := freetype.NewContext()
	f.SetDPI(72)      //设置DPI
	f.SetFont(font)   //设置字体
	f.SetFontSize(24) //设置字号
	f.SetClip(img.Bounds())
	f.SetDst(m)
	f.SetSrc(image.NewUniform(color.RGBA{R: 255, G: 0, B: 0, A: 255})) //设置颜色

	//新图写入原图和背景图
	draw.Draw(m, b, img, image.ZP, draw.Src)
	fmt.Println("Dx:", img.Bounds().Dx())
	fmt.Println("Dy:", img.Bounds().Dy())
	fmt.Println("x:", offset.X)
	fmt.Println("y:", offset.Y)
	content := time.Now().Format("2006-01-02 15:04:05")
	pt := freetype.Pt(offset.X-len(content), offset.Y)
	_, err = f.DrawString(content, pt)

	//输出图像
	imgw, _ := os.Create("new.jpg")
	jpeg.Encode(imgw, m, &jpeg.Options{100})

}
