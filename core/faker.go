package core

import (
	"fake-SAUer/global"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	postUrl = "https://ucapp.sau.edu.cn/wap/login/invalid"             // log in post target
	htmlUrl = "https://app.sau.edu.cn/form/wap/default?formid=10"      // html address
	URL     = "https://app.sau.edu.cn/form/wap/default/save?formid=10" // submit address
)

type Faker struct {
	Cnt  int 			// punch counts
	mu sync.Mutex 		// protect global.StuInfo
}

func NewFaker() (*Faker, error) {
	err := global.ReadConfig()
	if err != nil {
		return nil, err
	}

	for _, pStu := range global.G_CONF.StusInfos {
		fmt.Println(pStu, pStu.College, pStu.City, pStu.Account, pStu.Email, pStu.Province, pStu.Phone)
	}

	f := Faker{
		Cnt: checkInfo(),
	}

	return &f, nil
}

// Do 执行任务，返回成功数量
func (f *Faker) Do() (done int8) {
	var wg sync.WaitGroup
	for i := 0; i < f.Cnt; i++ {
		wg.Add(1)
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

			req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
			req.Header.Add("Accept-Encoding", "gzip,deflate,br")
			req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,zh-TW;q=0.8,en-US;q=0.7,en;q=0.6")
			req.Header.Add("Cache-Control", "no-cache")
			req.Header.Add("Connection", "keep-alive")

			req.Header.Add("Content-Length", strconv.Itoa(len(u.Encode())))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
			req.Header.Add("Host", "app.sau.edu.cn")
			req.Header.Add("Origin", "https://app.sau.edu.cn")
			req.Header.Add("Pragma", "no-cache")
			req.Header.Add("Referer", "https://app.sau.edu.cn/form/wap/default/index?formid=10&nn=4669.797748311082")

			req.Header.Add("Sec-Fetch-Dest", "empty")
			req.Header.Add("Sec-Fetch-Mode", "cors")
			req.Header.Add("Sec-Fetch-Site", "same-origin")
			req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
			req.Header.Add("X-Requested-With", "XMLHttpRequest")

			for _, c := range cks {
				req.AddCookie(c)
			}

			resp, err := http.DefaultClient.Do(req)
			//c := http.Client{
			//	CheckRedirect: func(req *http.Request, via []*http.Request) error {
			//		fmt.Println("redirect!")
			//		return http.ErrUseLastResponse
			//	},
			//}

			//resp, err := c.Do(req)
			if err != nil {
				fmt.Println("c.Do() err: ", err)
				return
			}

			//if err == http.ErrUseLastResponse {
			//	fmt.Println("最后一次重定向")
			//	defer resp.Body.Close()
			//	data, err := ioutil.ReadAll(resp.Body)
			//	if err != nil {
			//		fmt.Println("data, err := ioutil.ReadAll(resp.Body): ",err)
			//		return
			//	}
			//	fmt.Println(string(data))
			//	return
			//}
			//else if err != nil {
			//	fmt.Println("resp, err := http.DefaultClient.Do(req) err: ",err)
			//	return
			//}

			// 读取返回信息并打印字符串
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("data, err := ioutil.ReadAll(resp.Body): ", err)
				return
			}

			// msg := gjson.Get(string(data), "m").String()
			e := gjson.Get(string(data), "e").Int()
			if e == 0 {
				done++
			} else {
				if global.G_CONF.WithEmail.Account != "" {
					//TODO: email result.
					// gomail: could not send email 1: 550 Mail content denied. http://service.mail.qq.com/cgi-bin/help?subtype=1&&id=20022&&no=1000726 [MFzFTLSV4lzOGwIfv+UqxoSSC6s1Cw9zqHAGgKkhM21V12ZU/zcxWo5jtQFePQGG4w== IP: 223.88.165.204]

					//if err := f.E.SendMail(f.Conf.StusInfos[i].Email, "打卡通告", "今日打卡失败，请手动打卡"); err != nil {
					//	log.Printf("发送邮件失败%s\n", err)
					//}
				}
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("打卡完毕，一共%d个用户，成功了%d个\n", f.Cnt, done)
	return nil, done
}
