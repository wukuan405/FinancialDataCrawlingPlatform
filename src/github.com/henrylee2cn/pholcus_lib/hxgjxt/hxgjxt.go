package pholcus_lib

import (
	// 基础包
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析
	// "github.com/henrylee2cn/pholcus/logs"           //信息输出
	. "github.com/henrylee2cn/pholcus/app/spider" //必需
	// . "github.com/henrylee2cn/pholcus/app/spider/common" //选用

	// net包
	// "net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包
	// "regexp"
	"strconv"
	// "strings"
	// 其他包
	"fmt"
	// "math"
	// "time"
	//"log"
	//"log"
)

func init() {
	Hxgjxt.Register()
}

var Hxgjxt = &Spider{
	Name:        "华鑫国际信托",
	Description: "华鑫国际信托净值数据 [Auto Page] [http://www.cfitc.com/webfront/webpage/web/contentList/channelId/c8f5402dc3de432ab57257d2d652787b/pageNo/1]",
	// Pausetime: 300,
	// Keyin:   KEYIN,
	// Limit:        LIMIT,
	NotDefaultField: true,

	Namespace: func(*Spider) string {
		return "xintuo"
	},
	// 子命名空间相对于表名，可依赖具体数据内容，可选
	SubNamespace: func(self *Spider, dataCell map[string]interface{}) string {
		return "fund_src_nav"
	},

	EnableCookie: false,
	RuleTree: &RuleTree{

		Root: func(ctx *Context) {
			ctx.Aid(map[string]interface{}{"loop": [2]int{1, 10}, "Rule": "生成请求"}, "生成请求")
		},

		Trunk: map[string]*Rule{

			"生成请求": {

				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						ctx.AddQueue(&request.Request{
							Url:  "http://www.cfitc.com/nodejsService/articleListSearch/?callback=jQuery1830014371070777997375_1463928891321&website_id=c8f5402dc3de432ab57257d2d652787b&perPage=20&startNum=" + strconv.Itoa(loop[0]) + "&_=1463928891355",
							Rule: aid["Rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {

					query := ctx.GetText()
					fmt.Println(query)

					//query := ctx.GetDom()

					/*
						ss := query.Find("div#NewsList.List dl.FloatDiv.ff")

						ss.Each(func(i int, goq *goquery.Selection) {

							if url, ok := goq.Find("a").Attr("href"); ok {
								ctx.AddQueue(&request.Request{
									Url:  "http://www.cfitc.com/" + url,
									Rule: "获取结果",
								})
							}

						})
					*/
				},
			},

			"获取结果": {
				//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"名称",
					"净值",
					"累计净值",
					"估值日期",
				},

				ParseFunc: func(ctx *Context) {

					queryResult := ctx.GetDom()

					ssResult := queryResult.Find(".ContentTxt tbody").Find("tr")

					//var titleMingCheng string
					//ctx.GetTemp("title", &titleMingCheng)
					mingchen := queryResult.Find(".ContentTitle").Text()

					ssResult.Each(func(i int, goq *goquery.Selection) {

						titleLineResult := goq.Children().Eq(0).Text()
						if titleLineResult != "日期" && titleLineResult != "" {

							jingzhi := goq.Children().Eq(1).Text()
							leijijingzhi := goq.Children().Eq(3).Text()
							guzhiriqi := goq.Children().Eq(0).Text()

							ctx.Output(map[int]interface{}{
								0: mingchen,
								1: jingzhi,
								2: leijijingzhi,
								3: guzhiriqi,
							})
						}

					})
				},
			},
		},
	},
}