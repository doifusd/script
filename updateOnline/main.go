package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var invoiceKindType = map[string]string{
	"纸质发票(增值税普通发票)": "04",
	"电子发票(增值税普通发票)": "10",
	"电子发票(增值税专用发票)": "08",
	"纸质发票(增值税专用发票)": "01",
}

var invoiceType = map[string]string{
	"纸质发票(增值税普通发票)": "2",
	"电子发票(增值税普通发票)": "1",
	"电子发票(增值税专用发票)": "1",
	"纸质发票(增值税专用发票)": "2",
}

type InvoiceRequestContent struct {
	NoticeType    string                   `json:"notice_type"`
	InvoiceSource string                   `json:"invoice_source"`
	Code          string                   `json:"code"`
	Msg           string                   `json:"msg"`
	Content       InvoiceRequestSubContent `json:"content"`
}

type InvoiceRequestSubContent struct {
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
	ItemList                 []*ItemGoods      `json:"itemList"`
	InvPreviewQrcodePath     string            `json:"invPreviewQrcodePath"`
	InvPreviewQrcode         string            `json:"invPreviewQrcode"`
	RedSource                string            `json:"redSource"`
	TakerTel                 string            `json:"takerTel"`
	Unit                     string            `json:"unit"`
	InvoiceExtend            InvoiceExtendData `json:"invoiceExtend"`
}

type InvoiceExtendData struct {
	CipherText          string `json:"cipherText"`
	RedFlag             string `json:"redFlag"`
	SmsStatus           string `json:"smsStatus"`
	EmailStatus         string `json:"emailStatus"`
	ElectronicInvoiceNo string `json:"electronicInvoiceNo"`
}
type ItemGoods struct {
	GoodsName              string  `json:"goodsName"`
	TaxClassificationCode  string  `json:"taxClassificationCode"`
	SpecificationModel     string  `json:"specificationModel"`
	MeteringUnit           string  `json:"meteringUnit"`
	Quantity               float64 `json:"quantity"`
	UnitPrice              string  `json:"unitPrice"`
	TaxIncludeFlag         int     `json:"taxIncludeFlag"`
	ItemAmount             float64 `json:"itemAmount"`
	TaxRateValue           float64 `json:"taxRateValue"`
	TaxRateAmount          float64 `json:"taxRateAmount"`
	PreferentialPolicyFlag string  `json:"preferentialPolicyFlag"`
	ZeroTaxFlag            string  `json:"zeroTaxFlag"`
	VatSpecialManage       string  `json:"vatSpecialManage"`
	ItemProperty           int     `json:"itemProperty"`
	ItemNo                 string  `json:"itemNo"`
	DiscountAmount         float64 `json:"discountAmount"`
	DiscountTaxRateAmount  float64 `json:"discountTaxRateAmount"`
}

type InvoiceParam struct {
	InvoiceCode string
	InvoiceNo   string
}

type PiaoyitongAll struct {
	InvoiceType           string `gorm:"invoice_type" json:"invoice_type"` // 发票类型
	InvoiceNo             string `gorm:"invoice_no" json:"invoice_no"`
	InvoiceCode           string `gorm:"invoice_code" json:"invoice_code"`
	Qiyebianhao           string `gorm:"qiyebianhao" json:"qiyebianhao"`                         // 销方编号(企业)
	SellerTaxNo           string `gorm:"seller_tax_no" json:"seller_tax_no"`                     // 销方纳税人识别号
	SellerName            string `gorm:"seller_name" json:"seller_name"`                         // 销方名称
	SellerAddrTel         string `gorm:"seller_addr_tel" json:"seller_addr_tel"`                 // 销方地址电话
	SellerBanknameAccount string `gorm:"seller_bankname_account" json:"seller_bankname_account"` // 销方银行名称账号
	BuyetCompany          string `gorm:"buyet_company" json:"buyet_company"`                     // 购方编号(企业)
	BuyerTaxNo            string `gorm:"buyer_tax_no" json:"buyer_tax_no"`                       // 购方纳税人识别号
	BuyerName             string `gorm:"buyer_name" json:"buyer_name"`                           // 购方名称
	BuyerAddressTel       string `gorm:"buyer_address_tel" json:"buyer_address_tel"`             // 购方地址电话
	BuyerNameAccount      string `gorm:"buyer_name_account" json:"buyer_name_account"`           // 购方银行名称账号
	TaxRate               string `gorm:"tax_rate" json:"tax_rate"`                               // 税率
	InvoiceNoAmount       string `gorm:"invoice_no_amount" json:"invoice_no_amount"`             // 不含税金额
	TaxAmount             string `gorm:"tax_amount" json:"tax_amount"`                           // 税额
	InvoiceAmount         string `gorm:"invoice_amount" json:"invoice_amount"`                   // 含税金额
	InvoiceDate           string `gorm:"invoice_date" json:"invoice_date"`                       // 发票日期
	InvoiceOrigin         string `gorm:"invoice_origin" json:"invoice_origin"`                   // 发票来源
	Remark                string `gorm:"remark" json:"remark"`                                   // 备注
	InvoiceState          string `gorm:"invoice_state" json:"invoice_state"`                     // 发票状态
	OldInvoiceNo          string `gorm:"old_invoice_no" json:"old_invoice_no"`                   // 原发票号码
	OldInvoiceCdoe        string `gorm:"old_invoice_cdoe" json:"old_invoice_cdoe"`               // 原发票代码
	Ext1                  string `gorm:"ext_1" json:"ext_1"`                                     // 扩展字段1
	Ext2                  string `gorm:"ext_2" json:"ext_2"`                                     // 扩展字段2
	Ext3                  string `gorm:"ext_3" json:"ext_3"`                                     // 扩展字段3
	ElectronicInvoiceNo   string `gorm:"electronic_invoice_no" json:"electronic_invoice_no"`     // 全电发票号码
	CreatedAt             string `gorm:"created_at" json:"created_at"`                           // 创建时间
	BuyerBankAccount      string `gorm:"buyer_bank_account" json:"buyer_bank_account"`           // 购方银行账号
	BuyerBankName         string `gorm:"buyer_bank_name" json:"buyer_bank_name"`                 // 购方银行名称
	BuyerTel              string `gorm:"buyer_tel" json:"buyer_tel"`                             // 购方电话
	BuyerAddr             string `gorm:"buyer_addr" json:"buyer_addr"`                           // 购方地址
	RedTime               string `gorm:"red_time" json:"red_time"`                               // 红冲时间
	RedState              string `gorm:"red_state" json:"red_state"`                             // 红冲状态
	McCode                string `gorm:"mc_code" json:"mc_code"`                                 // 机器编码
	InvoiceStoreName      string `gorm:"invoice_store_name" json:"invoice_store_name"`           // 开票点名称
	InvoiceStoreCode      string `gorm:"invoice_store_code" json:"invoice_store_code"`           // 开票点代码
	InvoiceMcCode         string `gorm:"invoice_mc_code" json:"invoice_mc_code"`                 // 开机票号
	InvoiceAgentCode      string `gorm:"invoice_agent_code" json:"invoice_agent_code"`           // 税控终端码
	RefundDate            string `gorm:"refund_date" json:"refund_date"`                         // 退回时间
	RefundRemark          string `gorm:"refund_remark" json:"refund_remark"`                     // 退回备注
	RefundSate            string `gorm:"refund_sate" json:"refund_sate"`                         // 退回状态
	SpecillModel          string `gorm:"specill_model" json:"specill_model"`                     // 特殊发票标识
	SellerBankName        string `gorm:"seller_bank_name" json:"seller_bank_name"`               // 销方银行名称
	SellerTel             string `gorm:"seller_tel" json:"seller_tel"`                           // 销方电话
	SellerAddr            string `gorm:"seller_addr" json:"seller_addr"`                         // 销方地址
	SellerBankAccount     string `gorm:"seller_bank_account" json:"seller_bank_account"`         // 销方银行账号
	CheckCode             string `gorm:"check_code" json:"check_code"`                           // 校验码
	RedReason             string `gorm:"red_reason" json:"red_reason"`                           // 作废红冲原因
	DisableTime           string `gorm:"disable_time" json:"disable_time"`                       // 作废时间
	ScreatData            string `gorm:"screat_data" json:"screat_data"`                         // 密文
	DisableState          string `gorm:"disable_state" json:"disable_state"`                     // 预作废状态
	CasherName            string `gorm:"casher_name" json:"casher_name"`                         // 收款人姓名
	ReviewerName          string `gorm:"reviewer_name" json:"reviewer_name"`                     // 复合人姓名
	DrawerName            string `gorm:"drawer_name" json:"drawer_name"`                         // 开票人姓名
}

type PiaoyitongDetail struct {
	Mingxihao                 string  `gorm:"mingxihao" json:"mingxihao"`                                 // 发票明细号
	InvoiceNo                 string  `gorm:"invoice_no" json:"invoice_no"`                               // 发票号码
	InvoiceCode               string  `gorm:"invoice_code" json:"invoice_code"`                           // 发票代码
	Huoyingshoufu             string  `gorm:"huoyingshoufu" json:"huoyingshoufu"`                         // 货物或应税劳务代码
	Yingshuilaowu             string  `gorm:"yingshuilaowu" json:"yingshuilaowu"`                         // 货物或应税劳务名称
	Guige                     string  `gorm:"guige" json:"guige"`                                         // 规格型号
	Num                       float64 `gorm:"num" json:"num"`                                             // 数量
	Unit                      string  `gorm:"unit" json:"unit"`                                           // 数量单位
	Danjia                    string  `gorm:"danjia" json:"danjia"`                                       // 不含税单价
	Jine                      string  `gorm:"jine" json:"jine"`                                           // 不含税金额
	Rate                      string  `gorm:"rate" json:"rate"`                                           // 税率
	Tax                       float64 `gorm:"tax" json:"tax"`                                             // 税额
	InvoiceAmount             float64 `gorm:"invoice_amount" json:"invoice_amount"`                       // 含税金额
	Buhanshuizhekou           string  `gorm:"buhanshuizhekou" json:"buhanshuizhekou"`                     // 不含税折扣金额
	Zhekoushuie               float64 `gorm:"zhekoushuie" json:"zhekoushuie"`                             // 折扣税额
	Hanshuizhekoujine         float64 `gorm:"hanshuizhekoujine" json:"hanshuizhekoujine"`                 // 含税折扣金额
	Shangpinshuimu            string  `gorm:"shangpinshuimu" json:"shangpinshuimu"`                       // 商品税目
	Shuishoufenleibianma      string  `gorm:"shuishoufenleibianma" json:"shuishoufenleibianma"`           // 税收分类编码
	OrderNo                   string  `gorm:"order_no" json:"order_no"`                                   // 订单号
	CreatedAt                 string  `gorm:"created_at" json:"created_at"`                               // 创建时间
	Version                   string  `gorm:"version" json:"version"`                                     // 编码版本号
	Dingdanmingxi             string  `gorm:"dingdanmingxi" json:"dingdanmingxi"`                         // 订单明细号
	Jiagefangshi              string  `gorm:"jiagefangshi" json:"jiagefangshi"`                           // 价格方式
	Lingshuilvbiaozhi         string  `gorm:"lingshuilvbiaozhi" json:"lingshuilvbiaozhi"`                 // 零税率标志
	Xianshouyouhui            string  `gorm:"xianshouyouhui" json:"xianshouyouhui"`                       // 享受税收优惠政策内容
	Mianjidanwei              string  `gorm:"mianjidanwei" json:"mianjidanwei"`                           // 面积单位
	Kuadishibiaozhi           string  `gorm:"kuadishibiaozhi" json:"kuadishibiaozhi"`                     // 跨地市标志
	Jianzhuxiangmu            string  `gorm:"jianzhuxiangmu" json:"jianzhuxiangmu"`                       // 建筑项目名称
	Tudizengzhishui           string  `gorm:"tudizengzhishui" json:"tudizengzhishui"`                     // 土地增值税项目编号
	Zulinqi                   string  `gorm:"zulinqi" json:"zulinqi"`                                     // 租赁期起
	Zulinqijiezhi             string  `gorm:"zulinqijiezhi" json:"zulinqijiezhi"`                         // 租赁期止
	Fashengxiangxidizhi       string  `gorm:"fashengxiangxidizhi" json:"fashengxiangxidizhi"`             // 发生详细地址
	Hedingjishuijine          string  `gorm:"hedingjishuijine" json:"hedingjishuijine"`                   // 核定计税价格
	Shijichengjiaohanshuijine string  `gorm:"shijichengjiaohanshuijine" json:"shijichengjiaohanshuijine"` // 实际成交含税金额
	Shifouxiangshouyouhui     string  `gorm:"shifouxiangshouyouhui" json:"shifouxiangshouyouhui"`         // 是否享受税收优惠政策
	Zhekoujine                string  `gorm:"zhekoujine" json:"zhekoujine"`                               // 扣除额
}

var db *gorm.DB

func init() {
	var err error
	dsn := "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Silent),
		QueryFields: true,
	})

	//db.Logger = logger.Default.LogMode(logger.Silent)
	if err != nil {
		fmt.Println("err:", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(64)
}

func main() {
	getData()
}

// 查询本地的库
func getData() {
	fileName := "./dataSource.txt"

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
		contentArr := strings.Fields(conent)
		//fmt.Println("content:", contentArr)
		resp := GetInvoiceMain(contentArr[0], contentArr[1])
		//fmt.Printf("resp:%#v", resp)
		res, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		httpClient(string(res), i)
		time.Sleep(2 * time.Second)
		i++
	}
	//循环读取文件
	//从数据库中获取数据
	//组装
}

// 组装参数
func GetInvoiceMain(invoiceNo, invoiceCode string) (reqContent InvoiceRequestContent) {
	var invoiceMain PiaoyitongAll
	var invoiceItem []*PiaoyitongDetail
	err := db.Table("piaoyitong_all").Where("invoice_no=?", invoiceNo).Where("invoice_code=?", invoiceCode).First(&invoiceMain).Error
	if err != nil {
		fmt.Println("get invoice main failed")
	}
	err = db.Table("piaoyitong_detail").Where("invoice_no=?", invoiceNo).Where("invoice_code=?", invoiceCode).Find(&invoiceItem).Error
	if err != nil {
		fmt.Println("get invoice main failed")
	}

	reqContent = InvoiceRequestContent{
		NoticeType:    "invoice",
		InvoiceSource: "piaotong",
		Code:          "0000",
		Msg:           "提交成功",
		Content: InvoiceRequestSubContent{
			InvoiceIssueSource: "FPXF", //发票类型
			InvoiceReqSerialNo: fmt.Sprintf("fPXF-%d", time.Now().UnixNano()),
			InvoiceType:        invoiceType[invoiceMain.InvoiceType], //发票类型
			InvoiceKindCode:    invoiceKindType[invoiceMain.InvoiceType],
			InvoiceNo:          invoiceMain.InvoiceNo,   //发票号码
			InvoiceCode:        invoiceMain.InvoiceCode, //发票代码
			SellerTaxpayerNum:  invoiceMain.SellerTaxNo, //销方纳税人识别号
			SellerName:         invoiceMain.SellerName,  //销方名称
			//SellerAddress:      invoiceMain.SellerAddr,        //销方地址电话
			//SellerBankName:     invoiceMain.SellerBankName,    //销方银行名称账号
			//SellerBankAccount:  invoiceMain.SellerBanknameAccount, //销方银行名称账号
			BuyerTaxpayerNum: invoiceMain.BuyerTaxNo, //购方纳税人识别号
			BuyerName:        invoiceMain.BuyerName,  //购方名称
			//BuyerAddress:       invoiceMain.BuyerAddr,         //购方地址电话
			//BuyerBankName : row.Col(13)    //购方银行名称账号
			BuyerBankAccount: invoiceMain.BuyerBankAccount, //购方银行名称账号
			AmountWithTax:    invoiceMain.InvoiceAmount,    //含税金额
			InvoiceDate:      invoiceMain.InvoiceDate,      //发票开票日期
			Remark:           invoiceMain.Remark,           //备注
			OldInvoiceNo:     invoiceMain.OldInvoiceNo,     //原发票号码
			OldInvoiceCode:   invoiceMain.OldInvoiceCdoe,   //原发票代码
			ExtendData:       invoiceMain.Ext1 + invoiceMain.Ext2 + invoiceMain.Ext3,
			//	//InvoiceType : row.Col(29) //创建时间
			//	//BuyerBankAccount : row.Col(30) //购方银行账号
			BuyerBankName: invoiceMain.BuyerBankName, //购方银行名称
			BuyerTel:      invoiceMain.BuyerTel,      //购方电话
			BuyerAddress:  invoiceMain.BuyerAddr,     //购方地址
			//InvoiceType : row.Col(34)   //红冲时间
			//	//InvoiceType : row.Col(35)   //红冲状态
			MachineCode: invoiceMain.McCode, //机器编码

			SpecialInvoiceKind: invoiceMain.SpecillModel,      //特殊发票标识
			SellerBankName:     invoiceMain.SellerBankName,    //销方银行名称
			SellerTel:          invoiceMain.SellerTel,         //销方电话
			SellerAddress:      invoiceMain.SellerAddr,        //销方地址
			SellerBankAccount:  invoiceMain.SellerBankAccount, //销方银行账号
			CheckCode:          invoiceMain.CheckCode,         //校验码
			//	//InvoiceType : row.Col(51)        //作废时间
			//	//InvoiceType : row.Col(53)        //预作废状态
			//TakerEmail : invoiceMain. //发票开具后自动发送邮件的邮箱
			CasherName: invoiceMain.CasherName, //收款人姓名
			//InvoiceType : row.Col(56) //系统来源
			ReviewerName: invoiceMain.ReviewerName, //复核人姓名
			DrawerName:   invoiceMain.DrawerName,   //开票人姓名
			InvoiceExtend: InvoiceExtendData{
				CipherText: invoiceMain.ScreatData, //密文
				//RedFlag:,
				ElectronicInvoiceNo: invoiceMain.ElectronicInvoiceNo, //全电发票号码,
				SmsStatus:           fmt.Sprintf("%d", time.Now().UnixNano()),
			},
		},
	}

	var itemGoods []*ItemGoods
	for _, val := range invoiceItem {
		ZeroTaxFlag := "999"
		if val.Lingshuilvbiaozhi != "" {
			ZeroTaxFlag = val.Lingshuilvbiaozhi
		}
		TaxRateValue := 0.0
		if val.Rate != "" {
			rateStr := strings.ReplaceAll(val.Rate, "%", "")
			rateNum, _ := strconv.ParseFloat(rateStr, 64)
			TaxRateValue = rateNum / 100
		}
		TaxClassificationCode := ""
		if IsNum(val.Shuishoufenleibianma) {
			TaxClassificationCode = val.Shuishoufenleibianma
		}
		TaxIncludeFlag := 0
		if val.Tax != 0.00 {
			TaxIncludeFlag = 1
		}
		tmp := &ItemGoods{
			GoodsName:              val.Yingshuilaowu,
			TaxClassificationCode:  TaxClassificationCode,
			SpecificationModel:     val.Guige,
			MeteringUnit:           val.Unit,
			Quantity:               val.Num,
			UnitPrice:              val.Danjia,
			TaxIncludeFlag:         TaxIncludeFlag,
			ItemAmount:             val.InvoiceAmount,
			TaxRateValue:           TaxRateValue,
			TaxRateAmount:          val.Tax,
			PreferentialPolicyFlag: val.Shifouxiangshouyouhui,
			ZeroTaxFlag:            ZeroTaxFlag,
			DiscountAmount:         val.Hanshuizhekoujine,
			DiscountTaxRateAmount:  val.Zhekoushuie,
			ItemProperty:           1,
		}
		itemGoods = append(itemGoods, tmp)
	}
	reqContent.Content.ItemList = itemGoods
	return
}

// 同步线上接口
func httpClient(content string, count int) {
	//defer func(waitNum chan int) {
	//	<-waitNum
	//}(waitNum)

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
	url := "http://127.0.0.1:8002/api/v1/invoice/fix"
	//url := "https://service.test.lecoosys.com/api/v1/invoice/fix"
	//url := "https://service.lecoosys.com/api/v1/invoice/fix"
	fmt.Println("content:", content)
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

func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
