package faker

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"toolset/config"
	"toolset/email"
)

var (
	postUrl = "https://ucapp.sau.edu.cn/wap/login/invalid"        // 提交的目标地址
	htmlUrl = "https://app.sau.edu.cn/form/wap/default?formid=10" // html页面地址
	URL     = "https://app.sau.edu.cn/form/wap/default/save?formid=10"
)

type Faker struct {
	// 打卡人数
	Cnt int

	// 用户信息
	cf *config.Config

	// name->uuid
	us map[string]int

	e email.Email
}

func NewFaker() *Faker {
	f := Faker{
		us: make(map[string]int),
		cf: config.ReadConfig(),
	}
	for _, pStuInfo := range f.cf.StusInfos {
		fmt.Println(*pStuInfo)
	}

	f.GetUUID()
	f.Cnt = checkInfo(f.cf.StusInfos)
	return &f
}

// 将处理结果返回
func (f *Faker) Do() {
	var wg sync.WaitGroup
	var done int8
	for i := 0; i < f.Cnt; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			thisID := f.us[f.cf.StusInfos[i].Name]
			if thisID <= 0 {
				return
			}
			cks := GetCookie(f.cf.StusInfos[i].Account, f.cf.StusInfos[i].Passwd)
			u := bindInfo(f.cf.StusInfos[i], thisID)
			req, err := http.NewRequest("POST", URL, strings.NewReader(u.Encode()))
			if err != nil {
				panic("致命错误，POST提交表单失败！")
			}

			req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
			req.Header.Add("Accept-Encoding", "gzip,deflate,br")
			req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,zh-TW;q=0.8,en-US;q=0.7,en;q=0.6")
			req.Header.Add("Connection", "keep-alive")
			req.Header.Add("Content-Length", strconv.Itoa(len(u.Encode())))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
			req.Header.Add("Host", "app.sau.edu.cn")
			req.Header.Add("Origin", "https://app.sau.edu.cn")
			req.Header.Add("Referer", "https://app.sau.edu.cn/form/wap/default/index?formid=10&nn=4669.797748311082")
			req.Header.Add("Sec-Fetch-Dest", "empty")
			req.Header.Add("Sec-Fetch-Mode", "cors")
			req.Header.Add("Sec-Fetch-Site", "same-origin")
			req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36")
			req.Header.Add("X-Requested-With", "XMLHttpRequest")

			for _, c := range cks {
				req.AddCookie(c)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			// 读取返回信息并打印字符串
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return
			}

			// msg := gjson.Get(string(data), "m").String()
			if code := gjson.Get(string(data), "e").Int(); code == 0 {
				done++
				if f.cf.WithEmail.On {

				}
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("打卡完毕，一共%d个用户，拿到了%d个UUID\n", f.Cnt, done)
}
