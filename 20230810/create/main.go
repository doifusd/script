package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	readLog()
}

func readLog() {
	//读取目录下文件
	//判断获取的文件日期
	startTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("run start: ", startTime)

	fileName := "../data/tmp1.txt"

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("open file failed: %s \n", err.Error())
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	i := 1
	for {
		conent, errs := reader.ReadString('\n')
		if errs != nil {
			fmt.Println("err:", errs)
			break
		}
		httpClient(conent, i)
		i++
	}
	stopTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("run complete: ", stopTime)
}

func httpClient(content string, count int) {
	var netTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		}).DialContext,
		DisableKeepAlives:     true,
		TLSHandshakeTimeout:   10 * time.Second, // 限制TLS握手使用的时间
		MaxIdleConns:          10,
		MaxIdleConnsPerHost:   10,
		MaxConnsPerHost:       10,
		ResponseHeaderTimeout: 10 * time.Second, // 限制读取响应报文头使用的时间
		IdleConnTimeout:       90 * time.Second, // 连接最大空闲时间，超过这个时间就会被关闭
		ExpectContinueTimeout: 0,                // 等待服务器的第一个响应headers的时间，0表示没有超时
	}
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: netTransport,
	}
	//url := "http://127.0.0.1:8002/api/v1/invoice/create"
	//url := "https://service.test.lecoosys.com/api/v1/invoice/fix"
	url := "https://service.lecoosys.com/api/v1/invoice/create"
	//fmt.Println("content:", content)
	body := strings.NewReader(content)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatal(err)
	}

	//ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	//req.WithContext(ctx)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Accept-Charset", "utf-8")
	//req.Header.Add("Accept-Encoding","br, gzip, deflate")
	//req.Header.Add("Accept-Language", "zh-cn")
	//req.Header.Add("Connection", "keep-alive")
	//req.Header.Add("Cookie","xxxxxxxxxxxxxxx")
	//req.Header.Add("Content-Lenght",xxx)
	//req.Header.Add("Host", "www.baidu.com")
	//req.Header.Add("User-Agent", "http client 1.1.0")
	rep, err := client.Do(req)
	if err != nil {
		fmt.Println("http client err:", err)
		log.Fatal(err)
	}
	var data strings.Builder
	_, err = io.Copy(&data, rep.Body)
	rep.Body.Close()
	if err != nil {
		fmt.Println("http client resp err:", err)
		log.Fatal(err)
	}

	fmt.Printf("line-%d: resp:%s \n", count, data.String())
	// return string(data)
}
