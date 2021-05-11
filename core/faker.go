package core

import (
	"fake-SAUer/global"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	postUrl = "https://ucapp.sau.edu.cn/wap/login/invalid" // log in post target
	// URL htmlUrl = "https://app.sau.edu.cn/form/wap/default?formid=10"
	URL = "https://app.sau.edu.cn/form/wap/default/save?formid=10" // submit address
)

type Faker struct {
	Cnt int        // punch counts
	mu  sync.Mutex // protect global.StuInfo
}

func NewFaker() (*Faker, error) {
	err := global.ReadConfig()
	if err != nil {
		return nil, err
	}

	return &Faker{Cnt: checkInfo()}, nil
}

// Do 执行任务，返回成功数量
func (f *Faker) Do() (done int8) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(f.Cnt)

	// 复用一份header，只需要修改Content-Length即可
	h := make(http.Header, 16)
	h.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	h.Set("Accept-Encoding", "gzip,deflate,br")
	h.Set("Accept-Language", "zh-CN,zh;q=0.9,zh-TW;q=0.8,en-US;q=0.7,en;q=0.6")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	h.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	h.Set("Host", "app.sau.edu.cn")
	h.Set("Origin", "https://app.sau.edu.cn")
	h.Set("Pragma", "no-cache")
	h.Set("Referer", "https://app.sau.edu.cn/form/wap/default/index?formid=10&nn=4669.797748311082")
	h.Set("Sec-Fetch-Dest", "empty")
	h.Set("Sec-Fetch-Mode", "cors")
	h.Set("Sec-Fetch-Site", "same-origin")
	h.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
	h.Set("X-Requested-With", "XMLHttpRequest")

	for i := 0; i < f.Cnt; i++ {
		go func(i int) {

			defer wg.Done()
			thisID := global.G_CONF.StusInfos[i].UUID
			if thisID == "" {
				fmt.Println(global.G_CONF.StusInfos[i].Name, "的UUID为空，执行打卡失败")
				return
			}

			cks := GetCookie(global.G_CONF.StusInfos[i].Account, global.G_CONF.StusInfos[i].Passwd)
			u := bindInfo(global.G_CONF.StusInfos[i], thisID)
			req, err := http.NewRequest("POST", URL, strings.NewReader(u.Encode()))
			if err != nil {
				panic("致命错误，构造POST表单失败！")
			}

			req.Header = h
			req.Header.Set("Content-Length", strconv.Itoa(len(u.Encode())))

			for _, c := range cks {
				req.AddCookie(c)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println("c.Do() err: ", err)
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("data, err := ioutil.ReadAll(resp.Body): ", err)
				return
			}

			e := gjson.Get(string(data), "e").Int()
			if e == 0 {
				mu.Lock()
				done++
				mu.Unlock()
			} else {
				if global.G_CONF.E.Enabled == true {
					if err := global.G_CONF.E.SendMail(global.G_CONF.StusInfos[i].Email, "打卡通告", "今日打卡失败，请手动打卡"); err != nil {
						log.Printf("发送邮件失败%s\n", err)
					}
				}
			}
		}(i)
	}

	wg.Wait()
	return done
}
