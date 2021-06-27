package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"time"
)

// GetCookie 用账号密码获取cookies
func GetCookie(account, passwd string) []*http.Cookie {
	resp, err := http.PostForm(postUrl, url.Values{
		"username": {account},
		"password": {passwd},
	})
	if err != nil {
		log.Fatalf("致命错误！获取Cookie失败：%v", err)
	}
	defer resp.Body.Close()

	cks := resp.Cookies()
	return cks
}

// bindInfo 返回一个可以encode的url.Value结构
func bindInfo(f *StuInfo) url.Values {
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
	u.Add("tiwen2", "36.5")

	t := time.Now().Format("2006-01-02")
	u.Add("riqi", t)
	u.Add("id", f.Uuid)
	return u
}

func GetUuid(cks []*http.Cookie) string {
	req, _ := http.NewRequest("GET", htmlUrl, nil)
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
