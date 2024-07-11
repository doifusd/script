package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	es "github.com/elastic/go-elasticsearch/v7"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type Result struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore interface{} `json:"max_score"`
		Hits     []struct {
			Index   string      `json:"_index"`
			Type    string      `json:"_type"`
			Id      string      `json:"_id"`
			Score   interface{} `json:"_score"`
			Ignored []string    `json:"_ignored"`
			Source  struct {
				RequestBody string `json:"request_body"`
			} `json:"_source"`
			Sort []int64 `json:"sort"`
		} `json:"hits"`
	} `json:"hits"`
}

type Result2 struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		MaxScore interface{} `json:"max_score"`
		Hits     []struct {
			Index   string      `json:"_index"`
			Type    string      `json:"_type"`
			Id      string      `json:"_id"`
			Version int         `json:"_version"`
			Score   interface{} `json:"_score"`
			Ignored []string    `json:"_ignored"`
			Source  struct {
				Message string `json:"message"`
			} `json:"_source"`
			Fields struct {
				FieldsHostname           []string    `json:"fields.hostname"`
				UAOsName                 []string    `json:"UA_os_name"`
				RemoteAddrKeyword        []string    `json:"remote_addr.keyword"`
				BodyBytesSent            []int       `json:"body_bytes_sent"`
				FieldsTypeKeyword        []string    `json:"fields.type.keyword"`
				UAOsKeyword              []string    `json:"UA_os.keyword"`
				UAOsNameKeyword          []string    `json:"UA_os_name.keyword"`
				HttpLocationKeyword      []string    `json:"http_location.keyword"`
				Hostname                 []string    `json:"hostname"`
				Protocol                 []string    `json:"protocol"`
				HttpHostKeyword          []string    `json:"http_host.keyword"`
				UAName                   []string    `json:"UA_name"`
				CookieUnameKeyword       []string    `json:"cookie_uname.keyword"`
				UAOsFull                 []string    `json:"UA_os_full"`
				ContentLength            []int       `json:"content_length"`
				HttpCookie               []string    `json:"http_cookie"`
				UrlPathKeyword           []string    `json:"url_path.keyword"`
				FieldsEnvKeyword         []string    `json:"fields.env.keyword"`
				Method                   []string    `json:"method"`
				CookieUname              []string    `json:"cookie_uname"`
				UpstreamAddrKeyword      []string    `json:"upstream_addr.keyword"`
				XRequestId               []string    `json:"x_request_id"`
				FieldsHostnameKeyword    []string    `json:"fields.hostname.keyword"`
				TimeLocal                []time.Time `json:"time_local"`
				UADevice                 []string    `json:"UA_device"`
				VersionKeyword           []string    `json:"@version.keyword"`
				LogOffset                []int64     `json:"log.offset"`
				Tags                     []string    `json:"tags"`
				UANameKeyword            []string    `json:"UA_name.keyword"`
				HttpCookieKeyword        []string    `json:"http_cookie.keyword"`
				HttpReferer              []string    `json:"http_referer"`
				ProtocolKeyword          []string    `json:"protocol.keyword"`
				UAOs                     []string    `json:"UA_os"`
				HttpLocation             []string    `json:"http_location"`
				HostnameKeyword          []string    `json:"hostname.keyword"`
				Status                   []string    `json:"status"`
				Request                  []string    `json:"request"`
				RealRemoteAddr           []string    `json:"real_remote_addr"`
				UpstreamAddr             []string    `json:"upstream_addr"`
				HttpRefererKeyword       []string    `json:"http_referer.keyword"`
				StatusKeyword            []string    `json:"status.keyword"`
				TagsKeyword              []string    `json:"tags.keyword"`
				HttpHost                 []string    `json:"http_host"`
				HttpUserAgent            []string    `json:"http_user_agent"`
				UriKeyword               []string    `json:"uri.keyword"`
				XRequestIdKeyword        []string    `json:"x_request_id.keyword"`
				RequestTime              []float64   `json:"request_time"`
				HttpUserAgentKeyword     []string    `json:"http_user_agent.keyword"`
				UADeviceKeyword          []string    `json:"UA_device.keyword"`
				Version                  []string    `json:"@version"`
				MethodKeyword            []string    `json:"method.keyword"`
				LogFilePathKeyword       []string    `json:"log.file.path.keyword"`
				UrlPath                  []string    `json:"url_path"`
				RemoteAddr               []string    `json:"remote_addr"`
				FieldsEnv                []string    `json:"fields.env"`
				FieldsType               []string    `json:"fields.type"`
				RequestKeyword           []string    `json:"request.keyword"`
				HttpXForwardedForKeyword []string    `json:"http_x_forwarded_for.keyword"`
				Message                  []string    `json:"message"`
				Uri                      []string    `json:"uri"`
				RealRemoteAddrKeyword    []string    `json:"real_remote_addr.keyword"`
				Timestamp                []time.Time `json:"@timestamp"`
				RequestBody              []string    `json:"request_body"`
				LogFilePath              []string    `json:"log.file.path"`
				UAOsFullKeyword          []string    `json:"UA_os_full.keyword"`
				HttpXForwardedFor        []string    `json:"http_x_forwarded_for"`
				UpstreamResponseTime     []float64   `json:"upstream_response_time"`
			} `json:"fields"`
			IgnoredFieldValues struct {
				MessageKeyword     []string `json:"message.keyword"`
				RequestBodyKeyword []string `json:"request_body.keyword"`
			} `json:"ignored_field_values"`
			Highlight struct {
				UriKeyword     []string `json:"uri.keyword"`
				Request        []string `json:"request"`
				RequestBody    []string `json:"request_body"`
				Message        []string `json:"message"`
				Uri            []string `json:"uri"`
				UrlPathKeyword []string `json:"url_path.keyword"`
				UrlPath        []string `json:"url_path"`
			} `json:"highlight"`
			Sort []int64 `json:"sort"`
		} `json:"hits"`
	} `json:"hits"`
}

var (
	client *es.Client
)

type InvoiceParam struct {
	InvoiceCode string
	InvoiceNo   string
}

//func init() {
//	var err error
//	client, err = es.NewClient(es.Config{
//		Addresses: []string{"10.1.10.185:9200"},
//		Username:  "lecoo",
//		Password:  "tL8hHQ6DB8pT",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//}

func main() {
	//doIt()
	doItV1()
}

func doItV1() {
	ReadTxtV2("./datasource/2023-08-09.txt")
}

func doIt() {
	//contentStr := make(chan string, 10)
	contentParam := make(chan *InvoiceParam, 10)
	//go ReadTxt("./datasource/daichuli.txt", contentStr)
	go ReadTxtV1("./datasource/daichuli.txt", contentParam)
	//dataArr := []string{
	//	"2023.05.14",
	//	"2023.05.15",
	//	"2023.05.16",
	//	"2023.05.17",
	//	"2023.05.18",
	//}
	waitNum := make(chan int, 10)
	//for {
	//	select {
	//	case data, ok := <-contentStr:
	//		if ok {
	//			for _, val := range dataArr {
	//				go searchData(waitNum, data, val)
	//			}
	//		}
	//	}
	//}

	for {
		select {
		case data, ok := <-contentParam:
			if ok {
				//fmt.Println(data.InvoiceNo, data.InvoiceCode, waitNum)
				go searchDataV1(waitNum, data.InvoiceNo, data.InvoiceCode)
			}
		}
	}
}

func ReadTxtV1(fileName string, contentParam chan *InvoiceParam) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("open file failed: %s \n", err.Error())
		return
	}
	defer file.Close()
	line := bufio.NewReader(file)
	for {
		content, _, err := line.ReadLine()
		if err == io.EOF {
			break
		}
		contentArr := strings.Fields(string(content))
		contentParam <- &InvoiceParam{
			InvoiceCode: contentArr[1],
			InvoiceNo:   contentArr[0],
		}
	}
}

func ReadTxtV2(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("open file failed: %s \n", err.Error())
		return
	}
	defer file.Close()
	line := bufio.NewReader(file)
	for {
		content, _, err := line.ReadLine()
		if err == io.EOF {
			break
		}
		contentArr := strings.Fields(string(content))
		searchDataV2(contentArr[0], contentArr[1])
	}
}

func ReadTxt(fileName string, contentStr chan string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("open file failed: %s \n", err.Error())
		return
	}
	defer file.Close()
	line := bufio.NewReader(file)
	for {
		content, _, err := line.ReadLine()
		if err == io.EOF {
			break
		}
		contentStr <- string(content)
	}
}

func searchDataV2(invoiceNo, invoiceCode string) {
	//defer func(waitNum chan int) {
	//	<-waitNum
	//}(waitNum)
	//dataTime := strings.ReplaceAll(dateStr, ".", "-")
	dataTimeStart := "2023-05-17"
	dataTimeStop := "2023-05-18"
	//url := "http://10.1.10.185:9200/logstash-prod-nginx-" + dateStr + "/_doc/_search"
	url := "http://10.1.10.185:9200/logstash-prod-nginx-*/_doc/_search"
	//payload := strings.NewReader("{\n  \"_source\": [\n    \"message\"\n  ],\n  \"query\": {\n    \"bool\": {\n      \"must\": [],\n      \"filter\": [\n        {\n          \"bool\": {\n            \"filter\": [\n              {\n                \"multi_match\": {\n                  \"type\": \"phrase\",\n                  \"query\": \"" + invoiceNo + "\",\n                  \"lenient\": true\n                }\n              },\n              {\n                \"multi_match\": {\n                  \"type\": \"phrase\",\n                  \"query\": \"/api/v1/invoice/notify\",\n                  \"lenient\": true\n                }\n              }\n            ]\n          }\n        },\n        {\n          \"range\": {\n            \"time_local\": {\n              \"format\": \"strict_date_optional_time\",\n              \"gte\": \"" + dataTime + "T00:00:00.000Z\",\n              \"lte\": \"" + dataTime + "T23:59:00.000Z\"\n            }\n          }\n        }\n      ],\n      \"should\": [],\n      \"must_not\": [\n        {\n          \"match_phrase\": {\n            \"fields.tag\": \"ngw\"\n          }\n        }\n      ]\n    }\n  }\n}")
	//payload := strings.NewReader("{\n  \"_source\": [\n    \"message\"\n  ],\n  \"query\": {\n    \"bool\": {\n      \"must\": [],\n      \"filter\": [\n        {\n          \"bool\": {\n            \"filter\": [\n              {\n                \"multi_match\": {\n                  \"type\": \"phrase\",\n                  \"query\": \"" + invoiceNo + "\",\n                  \"lenient\": true\n                }\n              },\n              {\n                      \"multi_match\": {\n                        \"type\": \"phrase\",\n                        \"query\": \"" + invoiceCode + "\",\n                        \"lenient\": true\n                      }\n                    },\n              {\n                \"multi_match\": {\n                  \"type\": \"phrase\",\n                  \"query\": \"/api/v1/invoice/notify\",\n                  \"lenient\": true\n                }\n              }\n            ]\n          }\n        },\n        {\n          \"range\": {\n            \"time_local\": {\n              \"format\": \"strict_date_optional_time\",\n              \"gte\": \"" + dataTimeStart + "T00:00:00.000Z\",\n              \"lte\": \"" + dataTimeStop + "T23:59:59.000Z\"\n            }\n          }\n        }\n      ],\n      \"should\": [],\n      \"must_not\": [\n        {\n          \"match_phrase\": {\n            \"fields.tag\": \"ngw\"\n          }\n        }\n      ]\n    }\n  }\n}")
	payload := strings.NewReader("{\n  \"_source\": [\n    \"message\"\n  ],\n  \"query\": {\n    \"bool\": {\n      \"must\": [],\n      \"filter\": [\n        {\n          \"bool\": {\n            \"filter\": [\n              {\n                \"multi_match\": {\n                  \"type\": \"phrase\",\n                  \"query\": \"" + invoiceNo + "\",\n                  \"lenient\": true\n                }\n              },\n              {\n                      \"multi_match\": {\n                        \"type\": \"phrase\",\n                        \"query\": \"" + invoiceCode + "\",\n                        \"lenient\": true\n                      }\n                    },\n        {\n          \"range\": {\n            \"time_local\": {\n              \"format\": \"strict_date_optional_time\",\n              \"gte\": \"" + dataTimeStart + "T00:00:00.000Z\",\n              \"lte\": \"" + dataTimeStop + "T23:59:59.000Z\"\n            }\n          }\n        }\n      ],\n      \"should\": [],\n      \"must_not\": [\n        {\n          \"match_phrase\": {\n            \"fields.tag\": \"ngw\"\n          }\n        }\n      ]\n    }\n  }\n}")
	req, _ := http.NewRequest("GET", url, payload)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	req.WithContext(ctx)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Basic bGVjb286dEw4aEhRNkRCOHBU")
	req.Header.Set("User-Agent", "自定义的浏览器")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	var messageWriter strings.Builder
	_, err := io.Copy(&messageWriter, res.Body)
	//fmt.Println("res:", string(body))
	if messageWriter.String() == "" {
		return
	}

	var resp Result2
	err = json.Unmarshal([]byte(messageWriter.String()), &resp)
	if err != nil {
		fmt.Println("json decoder err:", err)
	}
	if len(resp.Hits.Hits) > 0 {
		message := resp.Hits.Hits[0].Source.Message
		fmt.Println("msg:", message)
		if message != "" {
			//sampleRegexp := regexp.MustCompile(`[0-9a-zA-Z]+\|`)
			//message = sampleRegexp.ReplaceAllString(message, "$1")
			//respArr := strings.Split(message, "|")
			//dstStr, err := ReplaceStringByRegex(message, `[0-9a-zA-Z]+\\|`, `[0-9a-zA-Z]+`)
			//dstStr, err := ReplaceStringByRegex(message, "[0-9a-zA-Z]+\\|", `[0-9a-zA-Z]+`)
			//if err != nil {
			//	fmt.Println("err:", err)
			//}
			//respArr := strings.Split(dstStr, "|")
			//fmt.Println(respArr[6])
			//获取字符串长度
			//从151处开始拿
			//长度到}}}结束
			//message[128]
			//result := strings.LastIndex(message, "}}} |")
			//fmt.Println(message[128 : len(message)-128-233])
			//fmt.Println(message[134 : result-134+3])
		}
		//	fmt.Println(resp.Hits.Hits[0].Source.RequestBody)
		//	isJSON := json.Valid([]byte(resp.Hits.Hits[0].Source.RequestBody))
		//	if isJSON {
		//		tmp := strings.ReplaceAll(resp.Hits.Hits[0].Source.RequestBody, "\\\\u", "\\u")
		//		tmp = strings.ReplaceAll(tmp, "\\\"", "\"")
		//		fmt.Println(tmp)
		//	} else {
		//		//从message处获取
		//		getMessage(waitNum, dateStr, invoiceNo)
		//	}
	}
}

func searchDataV1(waitNum chan int, invoiceNo, invoiceCode string) {
	defer func(waitNum chan int) {
		<-waitNum
	}(waitNum)
	//dataTime := strings.ReplaceAll(dateStr, ".", "-")
	dataTimeStart := "2023-05-01"
	dataTimeStop := "2023-07-01"
	//url := "http://10.1.10.185:9200/logstash-prod-nginx-" + dateStr + "/_doc/_search"
	url := "http://10.1.10.185:9200/logstash-prod-nginx-*/_doc/_search"
	//payload := strings.NewReader("{\n  \"_source\": [\n    \"message\"\n  ],\n  \"query\": {\n    \"bool\": {\n      \"must\": [],\n      \"filter\": [\n        {\n          \"bool\": {\n            \"filter\": [\n              {\n                \"multi_match\": {\n                  \"type\": \"phrase\",\n                  \"query\": \"" + invoiceNo + "\",\n                  \"lenient\": true\n                }\n              },\n              {\n                \"multi_match\": {\n                  \"type\": \"phrase\",\n                  \"query\": \"/api/v1/invoice/notify\",\n                  \"lenient\": true\n                }\n              }\n            ]\n          }\n        },\n        {\n          \"range\": {\n            \"time_local\": {\n              \"format\": \"strict_date_optional_time\",\n              \"gte\": \"" + dataTime + "T00:00:00.000Z\",\n              \"lte\": \"" + dataTime + "T23:59:00.000Z\"\n            }\n          }\n        }\n      ],\n      \"should\": [],\n      \"must_not\": [\n        {\n          \"match_phrase\": {\n            \"fields.tag\": \"ngw\"\n          }\n        }\n      ]\n    }\n  }\n}")
	payload := strings.NewReader("{\n  \"_source\": [\n    \"message\"\n  ],\n  \"query\": {\n    \"bool\": {\n      \"must\": [],\n      \"filter\": [\n        {\n          \"bool\": {\n            \"filter\": [\n              {\n                \"multi_match\": {\n                  \"type\": \"phrase\",\n                  \"query\": \"" + invoiceNo + "\",\n                  \"lenient\": true\n                }\n              },\n              {\n                      \"multi_match\": {\n                        \"type\": \"phrase\",\n                        \"query\": \"" + invoiceCode + "\",\n                        \"lenient\": true\n                      }\n                    },\n              {\n                \"multi_match\": {\n                  \"type\": \"phrase\",\n                  \"query\": \"/api/v1/invoice/notify\",\n                  \"lenient\": true\n                }\n              }\n            ]\n          }\n        },\n        {\n          \"range\": {\n            \"time_local\": {\n              \"format\": \"strict_date_optional_time\",\n              \"gte\": \"" + dataTimeStart + "T00:00:00.000Z\",\n              \"lte\": \"" + dataTimeStop + "T23:59:59.000Z\"\n            }\n          }\n        }\n      ],\n      \"should\": [],\n      \"must_not\": [\n        {\n          \"match_phrase\": {\n            \"fields.tag\": \"ngw\"\n          }\n        }\n      ]\n    }\n  }\n}")

	req, _ := http.NewRequest("GET", url, payload)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	req.WithContext(ctx)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Basic bGVjb286dEw4aEhRNkRCOHBU")
	req.Header.Set("User-Agent", "自定义的浏览器")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	var messageWriter strings.Builder
	_, err := io.Copy(&messageWriter, res.Body)
	//fmt.Println("res:", string(body))
	if messageWriter.String() == "" {
		return
	}

	var resp Result2
	err = json.Unmarshal([]byte(messageWriter.String()), &resp)
	if err != nil {
		fmt.Println("json decoder err:", err)
	}
	if len(resp.Hits.Hits) > 0 {
		message := resp.Hits.Hits[0].Source.Message
		fmt.Println("msg:", message)
		if message != "" {
			//sampleRegexp := regexp.MustCompile(`[0-9a-zA-Z]+\|`)
			//message = sampleRegexp.ReplaceAllString(message, "$1")
			//respArr := strings.Split(message, "|")
			//dstStr, err := ReplaceStringByRegex(message, `[0-9a-zA-Z]+\\|`, `[0-9a-zA-Z]+`)
			//dstStr, err := ReplaceStringByRegex(message, "[0-9a-zA-Z]+\\|", `[0-9a-zA-Z]+`)
			//if err != nil {
			//	fmt.Println("err:", err)
			//}
			//respArr := strings.Split(dstStr, "|")
			//fmt.Println(respArr[6])
			//获取字符串长度
			//从151处开始拿
			//长度到}}}结束
			//message[128]
			//result := strings.LastIndex(message, "}}} |")
			//fmt.Println(message[128 : len(message)-128-233])
			//fmt.Println(message[134 : result-134+3])
		}
		//	fmt.Println(resp.Hits.Hits[0].Source.RequestBody)
		//	isJSON := json.Valid([]byte(resp.Hits.Hits[0].Source.RequestBody))
		//	if isJSON {
		//		tmp := strings.ReplaceAll(resp.Hits.Hits[0].Source.RequestBody, "\\\\u", "\\u")
		//		tmp = strings.ReplaceAll(tmp, "\\\"", "\"")
		//		fmt.Println(tmp)
		//	} else {
		//		//从message处获取
		//		getMessage(waitNum, dateStr, invoiceNo)
		//	}
	}
}

func searchData(waitNum chan int, invoiceNo, dateStr string) {
	defer func(waitNum chan int) {
		<-waitNum
	}(waitNum)
	dataTime := strings.ReplaceAll(dateStr, ".", "-")
	url := "http://10.1.10.185:9200/logstash-prod-nginx-" + dateStr + "/_doc/_search"
	payload := strings.NewReader("{\n  \"_source\": [\n    \"message\"\n  ],\n  \"query\": {\n    \"bool\": {\n      \"must\": [],\n      \"filter\": [\n        {\n          \"bool\": {\n            \"filter\": [\n              {\n                \"multi_match\": {\n                  \"type\": \"phrase\",\n                  \"query\": \"" + invoiceNo + "\",\n                  \"lenient\": true\n                }\n              },\n              {\n                \"multi_match\": {\n                  \"type\": \"phrase\",\n                  \"query\": \"/api/v1/invoice/notify\",\n                  \"lenient\": true\n                }\n              }\n            ]\n          }\n        },\n        {\n          \"range\": {\n            \"time_local\": {\n              \"format\": \"strict_date_optional_time\",\n              \"gte\": \"" + dataTime + "T00:00:00.000Z\",\n              \"lte\": \"" + dataTime + "T23:59:00.000Z\"\n            }\n          }\n        }\n      ],\n      \"should\": [],\n      \"must_not\": [\n        {\n          \"match_phrase\": {\n            \"fields.tag\": \"ngw\"\n          }\n        }\n      ]\n    }\n  }\n}")
	//payload := strings.NewReader("{\n  \"_source\": [\n    \"message\"\n  ],\n  \"query\": {\n    \"bool\": {\n      \"must\": [],\n      \"filter\": [\n        {\n          \"bool\": {\n            \"filter\": [\n              {\n                \"multi_match\": {\n                  \"type\": \"phrase\",\n                  \"query\": \"" + invoiceNo + "\",\n                  \"lenient\": true\n                }\n              },\n              {\n                      \"multi_match\": {\n                        \"type\": \"phrase\",\n                        \"query\": \"" + invoiceCode + "\",\n                        \"lenient\": true\n                      }\n                    },\n              {\n                \"multi_match\": {\n                  \"type\": \"phrase\",\n                  \"query\": \"/api/v1/invoice/notify\",\n                  \"lenient\": true\n                }\n              }\n            ]\n          }\n        },\n        {\n          \"range\": {\n            \"time_local\": {\n              \"format\": \"strict_date_optional_time\",\n              \"gte\": \"" + dataTime + "T00:00:00.000Z\",\n              \"lte\": \"" + dataTime + "T23:59:59.000Z\"\n            }\n          }\n        }\n      ],\n      \"should\": [],\n      \"must_not\": [\n        {\n          \"match_phrase\": {\n            \"fields.tag\": \"ngw\"\n          }\n        }\n      ]\n    }\n  }\n}")

	req, _ := http.NewRequest("GET", url, payload)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	req.WithContext(ctx)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Basic bGVjb286dEw4aEhRNkRCOHBU")
	req.Header.Set("User-Agent", "自定义的浏览器")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	var messageWriter strings.Builder
	_, err := io.Copy(&messageWriter, res.Body)
	//fmt.Println("res:", string(body))
	if messageWriter.String() == "" {
		return
	}

	var resp Result2
	err = json.Unmarshal([]byte(messageWriter.String()), &resp)
	if err != nil {
		fmt.Println("json decoder err:", err)
	}
	if len(resp.Hits.Hits) > 0 {
		message := resp.Hits.Hits[0].Source.Message
		fmt.Println("msg:", message)
		if message != "" {
			//sampleRegexp := regexp.MustCompile(`[0-9a-zA-Z]+\|`)
			//message = sampleRegexp.ReplaceAllString(message, "$1")
			//respArr := strings.Split(message, "|")
			//dstStr, err := ReplaceStringByRegex(message, `[0-9a-zA-Z]+\\|`, `[0-9a-zA-Z]+`)
			//dstStr, err := ReplaceStringByRegex(message, "[0-9a-zA-Z]+\\|", `[0-9a-zA-Z]+`)
			//if err != nil {
			//	fmt.Println("err:", err)
			//}
			//respArr := strings.Split(dstStr, "|")
			//fmt.Println(respArr[6])
			//获取字符串长度
			//从151处开始拿
			//长度到}}}结束
			//message[128]
			//result := strings.LastIndex(message, "}}} |")
			//fmt.Println(message[128 : len(message)-128-233])
			//fmt.Println(message[134 : result-134+3])
		}
		//	fmt.Println(resp.Hits.Hits[0].Source.RequestBody)
		//	isJSON := json.Valid([]byte(resp.Hits.Hits[0].Source.RequestBody))
		//	if isJSON {
		//		tmp := strings.ReplaceAll(resp.Hits.Hits[0].Source.RequestBody, "\\\\u", "\\u")
		//		tmp = strings.ReplaceAll(tmp, "\\\"", "\"")
		//		fmt.Println(tmp)
		//	} else {
		//		//从message处获取
		//		getMessage(waitNum, dateStr, invoiceNo)
		//	}
	}
}

func ReplaceStringByRegex(str, rule, replace string) (string, error) {
	reg, err := regexp.Compile(rule)
	if reg == nil || err != nil {
		return "", errors.New("正则MustCompile错误:" + err.Error())
	}
	return reg.ReplaceAllString(str, replace), nil
}
