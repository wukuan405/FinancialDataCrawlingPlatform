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
	//"fmt"
	// "math"
	// "time"
	//"log"
	//"log"
	"strings"
)

func init() {
	Jiangshuxt.Register()
}

var Jiangshuxt = &Spider{
	Name:        "江苏信托",
	Description: "江苏信托净值数据 [Auto Page] [http://www.jsitic.net/418.html]",
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
			ctx.Aid(map[string]interface{}{"loop": [2]int{1, 2}, "Rule": "生成请求"}, "生成请求")
		},

		Trunk: map[string]*Rule{

			"生成请求": {

				//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"基金ID",
					"名称",
					"净值",
					"累计净值",
					"估值日期",
				},

				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					page := 0

					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {

						page++

						ctx.AddQueue(&request.Request{
							Url:  "http://www.jsitic.net/418.html",
							Rule: aid["Rule"].(string),
							Temp: map[string]interface{}{
								"level1pages": page,
							},
						})
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					ss := query.Find(".float_left.text-left.m-l-30").Children()

					recordCount := 0

					var page1 int
					page1 = ctx.GetTemp("level1pages", &page1).(int)

					count := 0
					ss.Each(func(i int, goq *goquery.Selection) {

						style, exist := goq.Attr("style")

						if exist && style == "height:30px;" && count < 7 {
							count++
							titleLine := goq.Children().Eq(0).Text()
							if titleLine != "产品名称" {
								mingchen := strings.TrimSpace(goq.Children().Eq(0).Text())
								jingzhi := strings.TrimSpace(goq.Children().Eq(2).Text())
								leijijingzhi := strings.TrimSpace(goq.Children().Eq(2).Text())
								guzhiriqi := strings.TrimSpace(goq.Children().Eq(1).Text())

								fundID := "XTJIANGSHU" + "P1" + strconv.Itoa(page1) + "L" + strconv.Itoa(recordCount)

								ctx.Output(map[int]interface{}{
									0: fundID,
									1: mingchen,
									2: jingzhi,
									3: leijijingzhi,
									4: guzhiriqi,
								})
							}
						}

					})
				},
			},
		},
	},
}
