package main

import (
	"context"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/buger/goreplay/proto"
	elastigo "github.com/mattbaird/elastigo/lib"

	"github.com/olivere/elastic"
)

type ESUriErorr struct{}

func (e *ESUriErorr) Error() string {
	return "Wrong ElasticSearch URL format. Expected to be: scheme://host/index_name"
}

type ESPlugin struct {
	Url     string
	Active  bool
	ApiPort string
	eConn   *elastigo.Conn
	Host    string
	Index   string
	indexor *elastigo.BulkIndexer
	done    chan bool
	client  *elastic.Client
}

func (p *ESPlugin) RttDurationToMs(d time.Duration) int64 {
	sec := d / time.Second
	nsec := d % time.Second
	fl := float64(sec) + float64(nsec)*1e-6
	return int64(fl)
}

type ESRequestResponse struct {
	ReqHost           string `json:"Req_Host"`
	ReqMethod         string `json:"Req_Method"`
	ReqURL            string `json:"Req_URL"`
	ReqBody           string `json:"Req_Body"`
	ReqUserAgent      string `json:"Req_User-Agent"`
	ReqXRealIP        string `json:"Req_X-Real-IP"`
	ReqXForwardedFor  string `json:"Req_X-Forwarded-For"`
	ReqConnection     string `json:"Req_Connection,omitempty"`
	ReqCookies        string `json:"Req_Cookies,omitempty"`
	RespStatusCode    string `json:"Resp_Status-Code"`
	RespBody          string `json:"Resp_Body"`
	RespProto         string `json:"Resp_Proto,omitempty"`
	RespContentLength string `json:"Resp_Content-Length,omitempty"`
	RespContentType   string `json:"Resp_Content-Type,omitempty"`
	RespSetCookie     string `json:"Resp_Set-Cookie,omitempty"`
	Rtt               int64  `json:"RTT"`
	Timestamp         time.Time
}

func parseURI(URI string) (err error, host, index string) {

	parsedUrl, parseErr := url.Parse(URI)

	if parseErr != nil {
		err = new(ESUriErorr)
		return
	}

	//  check URL validity by extracting host and index values.
	host = parsedUrl.Host
	urlPathParts := strings.Split(parsedUrl.Path, "/")
	index = urlPathParts[len(urlPathParts)-1]

	// force index specification in uri : ie no implicit index
	if host == "" || index == "" {
		err = new(ESUriErorr)
	}

	return
}

var initOnce sync.Once

func (p *ESPlugin) Init(URI string) {
	p.Url = URI
	var err error

	err, p.Host, p.Index = parseURI(URI)
	log.Println("Initializing Elasticsearch Plugin", p.Index, p.Host)
	t := time.Now()
	if p.Index == "" {
		p.Index = "gor-" + t.Format("2006-01-02")
	}

	if err != nil {
		log.Fatal("Can't initialize ElasticSearch plugin.", err)
	}

	initOnce.Do(func() {
		p.client, err = elastic.NewSimpleClient(
			elastic.SetURL("http://"+p.Host),
			// 设置错误日志
			elastic.SetErrorLog(log.New(os.Stderr, "ES-ERROR ", log.LstdFlags)),
			elastic.SetBasicAuth("elastic", "elastic"), // 账号密码
			// 设置info日志
			elastic.SetInfoLog(log.New(os.Stdout, "ES-INFO ", log.LstdFlags)),
		)
		if err != nil {
			log.Println(err)
		}
	})

	exists, err := p.client.IndexExists(p.Index).Do(context.Background())
	if err != nil {
		log.Println(err)
	}

	if !exists {
		_, err := p.client.CreateIndex(p.Index).Do(context.Background())
		if err != nil {
			log.Println(err)
		}
	}
	log.Println("Initialized Elasticsearch Plugin")
	return
}

func (p *ESPlugin) ResponseAnalyze(req, resp []byte, start, stop time.Time) {
	if len(resp) == 0 && len(req) == 0 {
		// nil http response - skipped elasticsearch export for this request
		log.Println("ResponseAnalyze ", resp, req)
		return
	}

	t := time.Now()
	rtt := p.RttDurationToMs(stop.Sub(start))
	// req = payloadBody(req)

	host := ESRequestResponse{
		ReqHost:           string(proto.Header(req, []byte("Host"))),
		ReqMethod:         string(proto.Method(req)),
		ReqURL:            string(proto.Path(req)),
		ReqBody:           string(proto.Body(req)),
		ReqUserAgent:      string(proto.Header(req, []byte("User-Agent"))),
		ReqXRealIP:        string(proto.Header(req, []byte("X-Real-IP"))),
		ReqXForwardedFor:  string(proto.Header(req, []byte("X-Forwarded-For"))),
		ReqConnection:     string(proto.Header(req, []byte("Connection"))),
		ReqCookies:        string(proto.Header(req, []byte("Cookie"))),
		RespStatusCode:    string(proto.Status(resp)),
		RespProto:         string(proto.Method(resp)),
		RespBody:          string(proto.Body(resp)),
		RespContentLength: string(proto.Header(resp, []byte("Content-Length"))),
		RespContentType:   string(proto.Header(resp, []byte("Content-Type"))),
		RespSetCookie:     string(proto.Header(resp, []byte("Set-Cookie"))),
		Timestamp:         t,
		Rtt:               rtt,
	}

	h, err := p.client.Index().Index(p.Index).BodyJson(host).Type("_doc").Do(context.Background()) //Type("ESRequestResponse").
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Indexed data with ID %s to index %s, type %s\n", h.Id, h.Index, h.Type)
	return
}
