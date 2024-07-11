package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/extrame/xls"
)

type InvoiceRequestContent struct {
	NoticeType    string `json:"notice_type"`
	InvoiceSource string `json:"invoice_source"`
	Code          string `json:"code"`
	Msg           string `json:"msg"`
	Content       struct {
		InvoiceReqSerialNo       string            `json:"invoiceReqSerialNo"`
		InvoiceIssueSource       string            `json:"invoiceIssueSource"`
		InvoiceIssuePlatformName string            `json:"invoiceIssuePlatformName"`
		InvoiceIssueWay          string            `json:"invoiceIssueWay"`
		InvoiceOperationCode     string            `json:"invoiceOperationCode"`
		SellerTaxpayerNum        string            `json:"sellerTaxpayerNum"`
		SellerName               string            `json:"sellerName"`
		SellerAddress            string            `json:"sellerAddress"`
		SellerTel                string            `json:"sellerTel"`
		SellerBankName           string            `json:"sellerBankName"`
		SellerBankAccount        string            `json:"sellerBankAccount"`
		BuyerTaxpayerNum         string            `json:"buyerTaxpayerNum"`
		BuyerName                string            `json:"buyerName"`
		BuyerAddress             string            `json:"buyerAddress"`
		BuyerBankName            string            `json:"buyerBankName"`
		BuyerBankAccount         string            `json:"buyerBankAccount"`
		BuyerTel                 string            `json:"buyerTel"`
		TakerEmail               string            `json:"takerEmail"`
		DrawerName               string            `json:"drawerName"`
		CasherName               string            `json:"casherName"`
		ReviewerName             string            `json:"reviewerName"`
		NoTaxAmount              string            `json:"noTaxAmount"`
		TaxAmount                string            `json:"taxAmount"`
		AmountWithTax            string            `json:"amountWithTax"`
		Remark                   string            `json:"remark"`
		TradeNo                  string            `json:"tradeNo"`
		TaxCategoryCodeVersion   string            `json:"taxCategoryCodeVersion"`
		ItemName                 string            `json:"itemName"`
		InvoiceKindCode          string            `json:"invoiceKindCode"`
		SpecialInvoiceKind       string            `json:"specialInvoiceKind"`
		InvoiceType              string            `json:"invoiceType"`
		TaxRateFlag              string            `json:"taxRateFlag"`
		SpecialInvoiceRedFlag    string            `json:"specialInvoiceRedFlag"`
		DetailedListFlag         string            `json:"detailedListFlag"`
		DetailedListItemName     string            `json:"detailedListItemName"`
		InvoiceDate              string            `json:"invoiceDate"`
		InvoiceStatus            string            `json:"invoiceStatus"`
		InvoiceCode              string            `json:"invoiceCode"`
		InvoiceNo                string            `json:"invoiceNo"`
		CheckCode                string            `json:"checkCode"`
		QrCode                   string            `json:"qrCode"`
		CipherText               string            `json:"cipherText"`
		InvoicePdf               string            `json:"invoicePdf"`
		InvoiceLayoutFileType    string            `json:"invoiceLayoutFileType"`
		DownloadUrl              string            `json:"downloadUrl"`
		RedFlag                  string            `json:"redFlag"`
		OldInvoiceCode           string            `json:"oldInvoiceCode"`
		OldInvoiceNo             string            `json:"oldInvoiceNo"`
		OldInvCanRedNoTaxAmount  string            `json:"oldInvCanRedNoTaxAmount"`
		OldInvCanRedTaxAmount    string            `json:"oldInvCanRedTaxAmount"`
		OldInvRedFlag            string            `json:"oldInvRedFlag"`
		DestroyFlag              string            `json:"destroyFlag"`
		RequestTime              string            `json:"requestTime"`
		ReceiveTime              string            `json:"receiveTime"`
		SmsStatus                string            `json:"smsStatus"`
		EmailStatus              string            `json:"emailStatus"`
		ExtensionNum             string            `json:"extensionNum"`
		MachineCode              string            `json:"machineCode"`
		AgentInvoiceFlag         string            `json:"agentInvoiceFlag"`
		ExtendData               string            `json:"extendData"`
		ItemList                 []ItemGoods       `json:"itemList"`
		InvPreviewQrcodePath     string            `json:"invPreviewQrcodePath"`
		InvPreviewQrcode         string            `json:"invPreviewQrcode"`
		RedSource                string            `json:"redSource"`
		TakerTel                 string            `json:"takerTel"`
		Unit                     string            `json:"unit"`
		InvoiceExtend            InvoiceExtendData `json:"invoiceExtend"`
	} `json:"content"`
}

type InvoiceExtendData struct {
	CipherText          string `json:"cipherText"`
	RedFlag             string `json:"redFlag"`
	SmsStatus           string `json:"smsStatus"`
	EmailStatus         string `json:"emailStatus"`
	ElectronicInvoiceNo string `json:"electronicInvoiceNo"`
}
type ItemGoods struct {
	GoodsName              string `json:"goodsName"`
	TaxClassificationCode  string `json:"taxClassificationCode"`
	SpecificationModel     string `json:"specificationModel"`
	MeteringUnit           string `json:"meteringUnit"`
	Quantity               string `json:"quantity"`
	UnitPrice              string `json:"unitPrice"`
	TaxIncludeFlag         string `json:"taxIncludeFlag"`
	ItemAmount             string `json:"itemAmount"`
	TaxRateValue           string `json:"taxRateValue"`
	TaxRateAmount          string `json:"taxRateAmount"`
	PreferentialPolicyFlag string `json:"preferentialPolicyFlag"`
	ZeroTaxFlag            string `json:"zeroTaxFlag"`
	VatSpecialManage       string `json:"vatSpecialManage"`
	ItemProperty           string `json:"itemProperty"`
	ItemNo                 string `json:"itemNo"`
}

type InvoiceParam struct {
	InvoiceCode string
	InvoiceNo   string
}

//1，从配置中找到数据
//2，去excel找到数据
//3，组装成请求需要的数据
//4，将数据发送到fix接口

func main() {
	readExcel("./data/data1")
}

var invoiceKindType = map[string]string{
	"纸质发票(增值税普通发票)": "04",
	"电子发票(增值税普通发票)": "10",
	"电子发票(增值税专用发票)": "08",
	"纸质发票(增值税专用发票)": "01",
}

var invoiceType = map[string]string{
	"纸质发票(增值税普通发票)": "1",
	"电子发票(增值税普通发票)": "1",
	"电子发票(增值税专用发票)": "1",
	"纸质发票(增值税专用发票)": "1",
}

func GetSourceFile() {
	contentParam := make(chan *InvoiceParam, 10)
	go ReadDataSource("./data/diffSource.txt", contentParam)

}

func ReadDataSource(fileName string, contentParam chan *InvoiceParam) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("open file failed: %s \n", err.Error())
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Fatalf("close file failed: %s \n", err.Error())
		}
	}(file)
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

func readExcel(filePath string) {
	//从文件夹中查找
	files, err := os.ReadDir(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, file := range files {
		//fmt.Println(file.Name())
		fileName := filePath + "/" + file.Name()
		//读取excel文件
		ReadXls(fileName)
	}
	//遍历文件
}

func ReadXls(fileName string) (res [][]string) {
	fmt.Println(fileName)
	//发票类型

	if xlFile, err := xls.Open(fileName, "utf-8"); err == nil {
		//第一个sheet
		sheetMain := xlFile.GetSheet(0)
		data := make([]InvoiceRequestContent, 0)
		for i := 0; i < int(sheetMain.MaxRow); i++ {
			row := sheetMain.Row(i)
			if row.LastCol() > 0 {
				fmt.Println(row.Col(1))
				var reqContent InvoiceRequestContent
				reqContent.NoticeType = "invoice"
				reqContent.InvoiceSource = "piaotong"
				reqContent.Content.InvoiceIssueSource = "FPXF"           //发票类型
				reqContent.Content.InvoiceType = invoiceType[row.Col(1)] //发票类型
				reqContent.Content.InvoiceKindCode = invoiceKindType[row.Col(1)]
				reqContent.Content.InvoiceNo = row.Col(2)   //发票号码
				reqContent.Content.InvoiceCode = row.Col(3) //发票代码
				//reqContent.Content.InvoiceType = row.Col(4)  //销方编号(企业)
				reqContent.Content.SellerTaxpayerNum = row.Col(5) //销方纳税人识别号
				reqContent.Content.SellerName = row.Col(6)        //销方名称
				//reqContent.Content.SellerAddress = row.Col(7)     //销方地址电话
				//reqContent.Content.SellerBankName = row.Col(8) //销方银行名称账号
				reqContent.Content.SellerBankAccount = row.Col(8) //销方银行名称账号
				//reqContent.Content.InvoiceType = row.Col(9)       //销方编号(企业)
				reqContent.Content.BuyerTaxpayerNum = row.Col(10) //购方纳税人识别号
				reqContent.Content.BuyerName = row.Col(11)        //购方名称
				reqContent.Content.BuyerAddress = row.Col(12)     //购方地址电话
				//reqContent.Content.BuyerBankName = row.Col(13)    //购方银行名称账号
				reqContent.Content.BuyerBankAccount = row.Col(13) //购方银行名称账号
				//reqContent.Content.InvoiceType = row.Col(14)      //税率
				//reqContent.Content.InvoiceType = row.Col(15)      //不含税金额
				//reqContent.Content.InvoiceType = row.Col(16)      //税额
				reqContent.Content.AmountWithTax = row.Col(17) //含税金额
				reqContent.Content.InvoiceDate = row.Col(18)   //发票开票日期
				//reqContent.Content.InvoiceType = row.Col(19)   //发票来源
				reqContent.Content.Remark = row.Col(20)         //备注
				reqContent.Code = "0000"                        //发票状态
				reqContent.Content.OldInvoiceNo = row.Col(22)   //原发票号码
				reqContent.Content.OldInvoiceCode = row.Col(23) //原发票代码
				//reqContent.Content.InvoiceType = row.Col(25) //扩展字段1
				//reqContent.Content.InvoiceType = row.Col(26) //扩展字段2
				//reqContent.Content.InvoiceType = row.Col(27) //扩展字段3
				reqContent.Content.ExtendData = row.Col(25) + row.Col(26) + row.Col(27)

				//reqContent.Content.InvoiceType = row.Col(24) //红字信息表编号

				//reqContent.Content.InvoiceType = row.Col(29) //创建时间
				//reqContent.Content.BuyerBankAccount = row.Col(30) //购方银行账号
				reqContent.Content.BuyerBankName = row.Col(31) //购方银行名称
				reqContent.Content.BuyerTel = row.Col(32)      //购方电话
				reqContent.Content.BuyerAddress = row.Col(33)  //购方地址
				//reqContent.Content.InvoiceType = row.Col(34)   //红冲时间
				//todo
				//reqContent.Content.InvoiceType = row.Col(35)   //红冲状态
				reqContent.Content.MachineCode = row.Col(36) //机器编码
				//reqContent.Content.InvoiceType = row.Col(37)   //开票点名称
				//reqContent.Content.InvoiceType = row.Col(38)   //开票点代码
				//reqContent.Content.InvoiceType = row.Col(39)   //开机票号
				//reqContent.Content.InvoiceType = row.Col(40)   //税控终端码
				//reqContent.Content.InvoiceType = row.Col(41)   //退回时间
				//reqContent.Content.InvoiceType = row.Col(42)   //退回备注
				//reqContent.Content.InvoiceType = row.Col(43)   //退回状态
				reqContent.Content.SpecialInvoiceKind = row.Col(44) //特殊发票标识
				reqContent.Content.SellerBankName = row.Col(45)     //销方银行名称
				reqContent.Content.SellerTel = row.Col(46)          //销方电话
				reqContent.Content.SellerAddress = row.Col(47)      //销方地址
				reqContent.Content.InvoiceType = row.Col(48)        //销方银行账号
				reqContent.Content.CheckCode = row.Col(49)          //校验码
				reqContent.Msg = row.Col(50)                        //作废红冲原因
				//reqContent.Content.InvoiceType = row.Col(51)        //作废时间
				//reqContent.Content.InvoiceType = row.Col(53)        //预作废状态
				reqContent.Content.TakerEmail = row.Col(54) //发票开具后自动发送邮件的邮箱
				reqContent.Content.CasherName = row.Col(55) //收款人姓名
				//reqContent.Content.InvoiceType = row.Col(56) //系统来源
				reqContent.Content.ReviewerName = row.Col(57) //复核人姓名
				reqContent.Content.DrawerName = row.Col(58)   //开票人姓名
				reqContent.Content.InvoiceExtend = InvoiceExtendData{
					CipherText: row.Col(52), //密文
					//RedFlag:,
					ElectronicInvoiceNo: row.Col(28), //全电发票号码,
				}
				data = append(data, reqContent)
			}
		}

		//sheetDetail := xlFile.GetSheet(1)
		//for i := 0; i < int(sheetDetail.MaxRow); i++ {
		//	row := sheetDetail.Row(i)
		//	//data := make([]string, 0)
		//	if row.LastCol() > 0 {
		//		//for j := 0; j < row.LastCol(); j++ {
		//		//	col := row.Col(j)
		//		//	data = append(data, col)
		//		//}
		//		//temp[i] = data
		//		fmt.Println("data:", row.Col(2), row.Col(3))
		//		//发票明细号
		//		//发票号码
		//		//发票代码
		//		//货物或应税劳务代码
		//		//货物或应税劳务名称
		//		//规格型号
		//		//数量
		//		//数量单位
		//		//不含税单价
		//		//不含税金额
		//		//税率
		//		//税额
		//		//含税金额
		//		//不含税折扣金额
		//		//折扣税额
		//		//含税折扣金额
		//		//商品税目
		//		//税收分类编码
		//		//订单号
		//		//创建时间
		//		//编码版本号
		//		//订单明细号
		//		//价格方式
		//		//零税率标志
		//		//享受税收优惠政策内容
		//		//面积单位
		//		//跨地市标志
		//		//建筑项目名称
		//		//土地增值税项目编号
		//		//租赁期起
		//		//租赁期止
		//		//发生详细地址
		//		//建筑服务发生地
		//		//不动产详细地址
		//		//不动产单元代码/网签合同备案编码
		//		//不动产地址（省市区）
		//		//房屋产权证书/不动产权证号码
		//		//核定计税价格
		//		//实际成交含税金额
		//		//是否享受税收优惠政策
		//		//扣除额
		//	}
		//}
	} else {
		fmt.Println("open_err:", err)
	}
	return res
}

// 组装数据
//func formatData() {}
//
//func httoClient() {
//
//}
