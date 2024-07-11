package main

//
//import (
//	"encoding/csv"
//	"fmt"
//	"github.com/extrame/xls"
//	"github.com/xuri/excelize/v2"
//	"gorm.io/driver/mysql"
//	"gorm.io/gorm"
//	"gorm.io/gorm/logger"
//	"io/ioutil"
//	"log"
//	"net/http"
//	"os"
//	"strconv"
//	"strings"
//	"sync"
//	"time"
//)
//
//type queue struct {
//	sync.WaitGroup
//	count chan int
//}
//
//type InvoiceData struct {
//	SellerTaxNo string `json:"seller_tax_no"`
//	InvoiceCode string `json:"invoice_code"`
//	InvoiceNo   string `json:"invoice_no"`
//}
//
//var db *gorm.DB
//
//func init() {
//	var err error
//	//dsn := "zhangzj28:Y0puLj7qUT3gI6FW@tcp(rm-2ze50gp4067jmfh1r.mysql.rds.aliyuncs.com:3306)/lecoo_service?charset=utf8mb4&parseTime=True&loc=Local"
//	//dsn_dev := "p_lecoo_service:BkWsDQ72Fr@tcp(rm-2ze6qp1n47r2vrx0v.mysql.rds.aliyuncs.com:3306)/lecoo_service?charset=utf8mb4&parseTime=True&loc=Local"
//	dsn := "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
//	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
//		Logger: logger.Default.LogMode(logger.Warn),
//	})
//	if err != nil {
//		fmt.Println("err:", err)
//	}
//	sqlDB, _ := db.DB()
//	sqlDB.SetMaxIdleConns(10)
//	sqlDB.SetMaxOpenConns(64)
//}
//
//func main() {
//	readLog()
//}
//
//// 读取excel文件
//func readLog() {
//	//读取目录下文件
//	//判断获取的文件日期
//	startTime := time.Now().Format("2006-01-02 15:04:05")
//	fmt.Println("run start: ", startTime)
//	//读取文件夹下excel
//	//将数据库中不存在的数据存储到日志中
//	//files, err := os.ReadDir("./data")
//	//if err != nil {
//	//	fmt.Println(err)
//	//	return
//	//}
//	fileName := "/Users/sky/Documents/script/checkdata/data1/invoice_check.csv"
//	//fileName := "/Users/sky/Documents/script/checkdata/f.log"
//
//	ReadCsv(fileName)
//
//	//for _, file := range files {
//	//	fileName := pathName + "" + file.Name()
//	//	fmt.Println("fileName:", fileName)
//	//	//file, err := os.Open(fileName)
//	//	//if err != nil {
//	//	//	log.Fatalf("open file failed: %s \n", err.Error())
//	//	//}
//	//	//defer file.Close()
//	//	//ReadXls(fileName)cd ..
//	//	ReadXlsx(fileName)
//	//}
//	////将数据库中不存在的数据存储到日志中，将原始日志中不存在单独存储
//
//	//reader := bufio.NewReader(file)
//	//i := 0
//	//// var wg sync.WaitGroup
//	//
//	//waitNum := make(chan int, 10)
//	//for {
//	//	conent, errs := reader.ReadString('\n')
//	//	if errs != nil {
//	//		fmt.Println("err:", errs)
//	//		//log.Fatal("read err:", errs)
//	//		//等于eof 退出
//	//		//errs != io.EOF
//	//		break
//	//	}
//	//	res := strings.ReplaceAll(conent, "\\\\u", "\\u")
//	//	str := strings.ReplaceAll(res, "\\\"", "\"")
//	//
//	//	str2, errs2 := zh2Unicode([]byte(str))
//	//	if errs2 != nil {
//	//		fmt.Println("to chinses err:", errs2)
//	//	}
//	//	isJSON := json.Valid([]byte(str2))
//	//	if isJSON {
//	//		waitNum <- 1
//	//		//todo go routine
//	//		// res := httpClient(str2)
//	//		// wg.Add(1)
//	//		go httpClient(waitNum, str2, i)
//	//		// wg.Wait()
//	//		time.Sleep(time.Microsecond * 5000)
//	//	}
//	//	i++
//	//}
//	//stopTime := time.Now().Format("2006-01-02 15:04:05")
//	//fmt.Println("run complete: ", stopTime)
//}
//
//func ReadCsv(file_path string) (res [][]string) {
//	file, err := os.Open(file_path)
//	if err != nil {
//		fmt.Errorf("open_err:", err)
//		return
//	}
//	defer file.Close()
//	// 初始化csv-reader
//	reader := csv.NewReader(file)
//	// 设置返回记录中每行数据期望的字段数，-1 表示返回所有字段
//	reader.FieldsPerRecord = -1
//	// 允许懒引号（忘记遇到哪个问题才加的这行）
//	reader.LazyQuotes = true
//	// 返回csv中的所有内容
//	record, read_err := reader.ReadAll()
//	if read_err != nil {
//		fmt.Errorf("read_err:", read_err)
//		return
//	}
//	waitNum := make(chan int, 10)
//	for _, value := range record {
//		//strArr := strings.Split(value[0], " ")
//
//		waitNum <- 1
//		//go checkData(value[2], value[1], waitNum)
//		go updateData(value[2], value[1], value[3], waitNum)
//		//go checkData(strArr[0], strArr[1], waitNum)
//	}
//	return record
//}
//
//func ReadXls(file_path string) (res [][]string) {
//	if xlFile, err := xls.Open(file_path, "utf-8"); err == nil {
//		fmt.Println(xlFile.Author)
//		//第一个sheet
//		sheet := xlFile.GetSheet(0)
//		if sheet.MaxRow != 0 {
//			//temp := make([][]string, sheet.MaxRow)
//			for i := 0; i < int(sheet.MaxRow); i++ {
//				row := sheet.Row(i)
//				fmt.Println("row:", row)
//				//data := make([]string, 0)
//				if row.LastCol() > 0 {
//					//for j := 0; j < row.LastCol(); j++ {
//					//	col := row.Col(j)
//					//	data = append(data, col)
//					//}
//					//temp[i] = data
//					fmt.Println("data:", row.Col(3), row.Col(4))
//				}
//			}
//			//res = append(res, temp...)
//		}
//	} else {
//		fmt.Println("open_err:", err)
//	}
//	return res
//}
//
//func ReadXlsx(file_path string) (res [][]string) {
//	xlFile, err := excelize.OpenFile(file_path)
//	rows, err := xlFile.GetRows("Sheet1")
//	if err != nil {
//		return
//	}
//	fmt.Println("rows:", rows)
//	//if xlFile, err := excelize.OpenFile(file_path); err == nil {
//	//	for index, sheet := range xlFile.Sheets {
//	//		//第一个sheet
//	//		if index == 0 {
//	//			temp := make([][]string, len(sheet.Rows))
//	//			for k, row := range sheet.Rows {
//	//				var data []string
//	//				for _, cell := range row.Cells {
//	//					data = append(data, cell.Value)
//	//				}
//	//				temp[k] = data
//	//			}
//	//			res = append(res, temp...)
//	//		}
//	//	}
//	//} else {
//	//	fmt.Println("open_err:", err)
//	//}
//	return res
//}
//
//func zh2Unicode(raw []byte) (str string, err error) {
//	str, err = strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
//	if err != nil {
//		return
//	}
//	return
//}
//
////func DataSource() {
////	var err error
////	dsn := "zhangzj28:Y0puLj7qUT3gI6FW@tcp(rm-2ze50gp4067jmfh1r.mysql.rds.aliyuncs.com:3306)/lecoo_service?charset=utf8mb4&parseTime=True&loc=Local"
////	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
////	if err != nil {
////		fmt.Println("err:", err)
////	}
////
//// 查找 code 字段值为 D42 的记录
////}
//
//func checkData(invoiceCode, invoiceNo string, waitNum chan int) {
//	defer func(waitNum chan int) {
//		<-waitNum
//	}(waitNum)
//	var invoiceRes InvoiceData
//	result := db.Table("server_invoice").Where("invoice_code=?", invoiceCode).Where("invoice_no=?", invoiceNo).Order("id desc").First(&invoiceRes)
//	if result.Error != nil {
//		fmt.Println("err:", result.Error)
//	}
//	if result.RowsAffected == 0 {
//		fmt.Printf("%s %s", invoiceCode, invoiceNo)
//	}
//}
//
//func updateData(invoiceCode, invoiceNo, SellerTaxNo string, waitNum chan int) {
//	defer func(waitNum chan int) {
//		<-waitNum
//	}(waitNum)
//	var cnt int64
//	//result := db.Table("server_invoice").Where("invoice_code=?", invoiceCode).Where("invoice_no=?", invoiceNo).Order("id desc").First(&invoiceRes)
//	result := db.Table("invoice_result").Where("invoice_code=?", invoiceCode).Where("invoice_no=?", invoiceNo).Count(&cnt)
//	if result.Error != nil {
//		fmt.Println("err:", result.Error)
//	}
//	if cnt > 0 {
//		res := db.Table("invoice_result").Where("invoice_code=?", invoiceCode).Where("invoice_no=?", invoiceNo).Update("seller_tax_no", SellerTaxNo)
//		if res.Error != nil {
//			fmt.Println("update err:", res.Error)
//		}
//	}
//}
//
//func httpClient(waitNum chan int, content string, count int) {
//	defer func(waitNum chan int) {
//		<-waitNum
//	}(waitNum)
//	client := &http.Client{}
//	// url := "http://127.0.0.1:8002/api/v1/invoice/fix"
//	url := "https://service.test.lecoosys.com/api/v1/invoice/fix"
//
//	body := strings.NewReader(content)
//	req, err := http.NewRequest("POST", url, body)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	req.Header.Add("content-type", "application/json")
//	req.Header.Add("Accept-Charset", "utf-8")
//	//req.Header.Add("Accept-Encoding","br, gzip, deflate")
//	//req.Header.Add("Accept-Language", "zh-cn")
//	//req.Header.Add("Connection", "keep-alive")
//	//req.Header.Add("Cookie","xxxxxxxxxxxxxxx")
//	//req.Header.Add("Content-Lenght",xxx)
//	//req.Header.Add("Host", "www.baidu.com")
//	//req.Header.Add("User-Agent", "http client 1.1.0")
//	rep, err := client.Do(req)
//	if err != nil {
//		fmt.Println("http client err:", err)
//		log.Fatal(err)
//	}
//	data, err := ioutil.ReadAll(rep.Body)
//	rep.Body.Close()
//	if err != nil {
//		fmt.Println("http client resp err:", err)
//		log.Fatal(err)
//	}
//
//	fmt.Printf("line%d: http content:%s resp:%s \n", count, content, string(data))
//	// return string(data)
//}
//
////func GetAllFile(pathname string) error {
////	rd, err := ioutil.ReadDir(pathname)
////	if err != nil {
////		fmt.Println("read dir fail:", err)
////		return err
////	}
////
////	for _, fi := range rd {
////		if !fi.IsDir() {
////			//fullName := pathname + "/" + fi.Name()
////			//s = append(s, fullName)
////
////		}
////	}
////	return nil
////}
