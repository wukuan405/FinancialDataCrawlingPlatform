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
	// "fmt"
	// "math"
	// "time"
	//"log"
	//"log"
)

func init() {
	Fzdyxt.Register()
}

var Fzdyxt = &Spider{
	Name:        "方正东亚信托",
	Description: "方正东亚信托净值数据 [Auto Page] [http://www.fd-trust.com/index/show/tid/47.html]",
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
							Url:  "http://www.fd-trust.com/home/index/show/tid/47/p/" + strconv.Itoa(loop[0]) + ".html",
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

					ss := query.Find(".pro_cont ul").Find("li")

					var page1 int
					ctx.GetTemp("level1pages", &page1)

					count := 0

					ss.Each(func(i int, goq *goquery.Selection) {

						goqueryResult := goq.Find("a")
						titleLine := goqueryResult.Children().Eq(1).Text()
						if titleLine != "产品名称" {
							mingchen := goqueryResult.Children().Eq(0).Text()
							jingzhi := goqueryResult.Children().Eq(1).Text()
							leijijingzhi := goqueryResult.Children().Eq(1).Text()
							guzhiriqi := goqueryResult.Children().Eq(2).Text()

							count++
							fundID := "XTFANGZHENGDONGYA" + "P1" + strconv.Itoa(page1) + "L" + strconv.Itoa(count)

							ctx.Output(map[int]interface{}{
								0: fundID,
								1: mingchen,
								2: jingzhi,
								3: leijijingzhi,
								4: guzhiriqi,
							})
						}

					})
				},
			},
		},
	},
}
