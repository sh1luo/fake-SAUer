package core

import (
	"errors"
	"fake-SAUer/conf"
	"fake-SAUer/notice"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	PostUrl   = "https://ucapp.sau.edu.cn/wap/login/invalid"             // log in post target
	HtmlUrl   = "https://app.sau.edu.cn/form/wap/default?formid=10"      // html url to get UUID
	submitURL = "https://app.sau.edu.cn/form/wap/default/save?formid=10" // submit address
)

type Faker struct {
	Cnt int
	
	Notifier notice.Notifier
	
	EnableHTTP bool
}

func NewFaker(enableHTTP bool) (f *Faker, err error) {
	if len(conf.GlobalConfig.StusInfo) == 0 {
		return nil, errors.New("has no valid students information")
	}
	
	f = &Faker{}
	if conf.GlobalConfig.NotifierInfo.Method != "" {
		inNotifier := reflect.ValueOf(conf.GlobalConfig.NotifierInfo).FieldByName(conf.GlobalConfig.NotifierInfo.Method).Interface()
		f.Notifier = notice.NewNotifier(conf.GlobalConfig.NotifierInfo.Method, inNotifier)
	}
	f.Cnt = len(conf.GlobalConfig.StusInfo)
	f.EnableHTTP = enableHTTP
	
	//reVal := reflect.ValueOf(conf.GlobalConfig.NotifierInfo)
	//reType := reflect.TypeOf(conf.GlobalConfig.NotifierInfo)
	// TODO: Now only handle first layer of nested struct
	//for i := 0; i < reVal.NumField(); i++ {
	//	if reVal.Field(i).Kind() == reflect.Struct || reVal.Field(i).Kind() == reflect.Ptr {
	//		reSubVal := reVal.Field(i)
	//		flag := -1
	//		for j := 0; j < reSubVal.NumField(); j++ {
	//			if reSubVal.Field(i).Interface().(string) == "" {
	//				flag = j
	//				break
	//			}
	//		}
	//		if flag != -1 {
	//			f.Notifier = notice.NewNotifier(reSubVal.Field(0).String(), "xx@qq.com", "xxxx", "smtp.qq.com", 465)
	//		}
	//	} else {
	//
	//	}
	//}
	
	return f, nil
}

// Do 返回打卡成功数量
func (f *Faker) Do() (done int) {
	wg := &sync.WaitGroup{}
	wg.Add(f.Cnt)
	for i := 0; i < f.Cnt; i++ {
		go func(i int) {
			defer func() {
				if e := recover(); e != nil && f.Notifier != nil && conf.GlobalConfig.StusInfo[i].To != "" {
					if e = f.Notifier.Notice(conf.GlobalConfig.StusInfo[i].To, "Failed Sign-in Message", fmt.Sprintf("Failed to perform the task today:\n\t%s", e)); e != nil {
						panic(e)
					}
				}
			}()
			defer wg.Done()
			acc, passwd := conf.GlobalConfig.StusInfo[i].Account, conf.GlobalConfig.StusInfo[i].Passwd
			cks, err := getCookies(acc, passwd)
			if err != nil {
				panic(err)
			}
			
			if conf.GlobalConfig.StusInfo[i].Uuid == "" {
				conf.GlobalConfig.StusInfo[i].Uuid = getUuid(cks)
			}
			body := structs2urlValues(conf.GlobalConfig.StusInfo[i])
			
			req, _ := http.NewRequest("POST", submitURL, strings.NewReader(body.Encode()))
			// TODO: reuse http header
			h := make(http.Header, 16)
			h.Set("Location", "辽宁省沈阳市")
			h.Set("Accept", "application/json, text/javascript, */*; q=0.01")
			h.Set("Accept-Encoding", "gzip,deflate,br")
			h.Set("Accept-Language", "zh-CN,zh;q=0.9,zh-TW;q=0.8,en-US;q=0.7,en;q=0.6")
			h.Set("Cache-Control", "no-cache")
			h.Set("Connection", "keep-alive")
			h.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
			h.Set("Host", "app.sau.edu.cn")
			h.Set("Origin", "https://app.sau.edu.cn")
			h.Set("Pragma", "no-cache")
			h.Set("sec-ch-ua-mobile", "?0")
			h.Set("Referer", "https://app.sau.edu.cn/form/wap/default/index?formid=10&nn=4669.797748311082")
			h.Set("Sec-Fetch-Dest", "empty")
			h.Set("Sec-Fetch-Mode", "cors")
			h.Set("Sec-Fetch-Site", "same-origin")
			h.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
			h.Set("X-Requested-With", "XMLHttpRequest")
			req.Header = h
			req.Header.Set("Content-Length", strconv.Itoa(len(body.Encode())))
			for _, c := range cks {
				req.AddCookie(c)
			}
			
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				panic(err)
			}
			
			decoder := mahonia.NewDecoder("utf-8")
			data, err := ioutil.ReadAll(decoder.NewReader(resp.Body))
			err = resp.Body.Close()
			if err != nil {
				panic(err)
			}
			
			_ = data
			
			done++
		}(i)
	}
	wg.Wait()
	return
}

// GetCookie 用账号密码获取cookies
func getCookies(account, passwd string) ([]*http.Cookie, error) {
	resp, err := http.PostForm(PostUrl, url.Values{
		"username": {account},
		"password": {passwd},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	cks := resp.Cookies()
	return cks, nil
}

// bindInfo 返回一个可以encode的url.Value结构
func structs2urlValues(f *conf.StuInfo) url.Values {
	u := url.Values{}
	u.Add("xingming", f.Name)
	u.Add("xuehao", f.Account)
	u.Add("shoujihao", f.Phone)
	u.Add("danweiyuanxi", f.College)
	u.Add("dangqiansuozaishengfen", f.Province)
	u.Add("dangqiansuozaichengshi", f.City)
	
	u.Add("shifouyuhubeiwuhanrenyuanmiqie", "否")
	u.Add("shifoujiankangqingkuang", "是")
	u.Add("shifoujiechuguohubeihuoqitayou", "否")
	u.Add("fanhuididian", "")
	u.Add("shifouweigelirenyuan", "否")
	u.Add("shentishifouyoubushizhengzhuan", "否")
	u.Add("shifouyoufare", "否")
	u.Add("qitaxinxi", "")
	u.Add("tiwen", "36.3")
	u.Add("tiwen1", "36.4")
	u.Add("tiwen2", "36.4")
	
	t := time.Now().Format("2006-01-02")
	u.Add("riqi", t)
	u.Add("id", f.Uuid)
	return u
}

func getUuid(cks []*http.Cookie) string {
	req, _ := http.NewRequest("GET", HtmlUrl, nil)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,zh-TW;q=0.8,en-US;q=0.7,en;q=0.6")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "app.sau.edu.cn")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("sec-ch-ua", "Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"91\", \"Chromium\";v=\"91")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36")
	
	for _, c := range cks {
		req.AddCookie(c)
	}
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return ""
	}
	
	var sid string
	doc.Find(".footReturn > input").Eq(0).Each(func(i int, selection *goquery.Selection) {
		sid, _ = selection.Attr("value")
	})
	return sid
}
