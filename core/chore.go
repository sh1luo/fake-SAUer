package core

import (
	"fake-SAUer/conf"
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
func bindInfo(f *conf.StuInfo, sid string) url.Values {
	u := url.Values{}
	u.Add("xingming", f.Name)
	u.Add("xuehao", f.Account)
	u.Add("shoujihao", f.Phone)
	u.Add("danweiyuanxi", f.College)
	u.Add("dangqiansuozaishengfen", f.Province)
	u.Add("dangqiansuozaichengshi", f.City)

	t := time.Now().Format("2006-01-02 15:04:05")
	u.Add("riqi", t[:10])

	u.Add("shifouyuhubeiwuhanrenyuanmiqie", "否")
	u.Add("shifoujiankangqingkuang", "是")
	u.Add("shifoujiechuguohubeihuoqitayou", "否")
	u.Add("shifouweigelirenyuan", "否")
	u.Add("shentishifouyoubushizhengzhuan", "否")
	u.Add("shifouyoufare", "否")

	u.Add("tiwen", "36.5")
	u.Add("tiwen1", "36.5")
	u.Add("tiwen2", "36.5")

	u.Add("fanhuididian", "")
	u.Add("qitaxinxi", "")

	u.Add("id", sid)
	return u
}

func GetUuid(cks []*http.Cookie) string {
	req, _ := http.NewRequest("GET", htmlUrl, nil)
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

