package signup

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	postUrl = "https://ucapp.sau.edu.cn/wap/login/invalid"        // 提交的目标地址
	htmlUrl = "https://app.sau.edu.cn/form/wap/default?formid=10" // html页面地址
	URL     = "https://app.sau.edu.cn/form/wap/default/save?formid=10"
)

type postForm struct {
	// base info
	StuName     string
	StuNumber   string
	PhoneNumber string
	College     string
	Province    string
	City        string

	id   string //区分用户的唯一标识码
	riqi string
}

func NewPostForm() *postForm {
	t := time.Now().Format("2006-01-02 15:04:05")
	return &postForm{
		StuName:     "",
		StuNumber:   "",
		PhoneNumber: "",
		College:     "",
		Province:    "",
		City:        "",
		riqi:        t[:10],
	}
}

// 将处理结果返回，发到邮箱
func (p *postForm) Signup(account, passwd string) (retStr string, err error) {
	// 获取用户标识的唯一ID
	cks := getCookie(account, passwd)
	fmt.Println("logincookie",cks)
	spID := getSpID(cks)
	p.id = spID
	u := handleFuckInfo(p)

	fmt.Println("spid:", spID) //debug info

	// fmt.Println("strings.NewReader:",strings.NewReader(u.Encode()))

	req, err := http.NewRequest("POST", URL, strings.NewReader(u.Encode()))
	if err != nil {
		return "", err
	}

	// 根据浏览器，模拟请求头
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Accept-Encoding", "gzip,deflate,br")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,zh-TW;q=0.8,en-US;q=0.7,en;q=0.6")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Length", strconv.Itoa(len(u.Encode())))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	//req.Header.Add("Cookie", "vjuid=18297; vjvd=152d0065583078311c9f0ea3360c8c24; vt=114490343")
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
		fmt.Println(c)
	}

	// 新建默认参数的客户端，执行预先给定的 request 请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取返回信息并打印字符串
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 校验返回值是否正确
	code := gjson.Get(string(data), "e").Int()
	// msg := gjson.Get(string(data), "m").String()

	fmt.Println("用户唯一id为：",spID)


	if code == 0 {
<<<<<<< HEAD
		retStr = fmt.Sprintf("%s打卡信息为：\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
			p.riqi, p.StuName, p.StuNumber, p.PhoneNumber, p.College, p.Province, p.City, p.id)
=======
		retStr = fmt.Sprintf("%s 打卡信息为：\n%s\n%s\n%s\n%s\n%s\n%s\n",
			p.riqi,p.StuName, p.StuNumber, p.PhoneNumber, p.College, p.Province, p.City)
>>>>>>> 73f7a7a4b3134b68b63cd66b224bccbcfa969cc1
		return retStr, nil
	}
	return "", errors.New("未知错误")
}

// 整合信息，返回一个可以encode的url.Value格式结构
func handleFuckInfo(p *postForm) url.Values {
	u := url.Values{}
	u.Add("xingming", p.StuName)
	u.Add("xuehao", p.StuNumber)
	u.Add("shoujihao", p.PhoneNumber)
	u.Add("danweiyuanxi", p.College)
	u.Add("dangqiansuozaishengfen", p.Province)
	u.Add("dangqiansuozaichengshi", p.City)
	u.Add("riqi", p.riqi)

	u.Add("shifouyuhubeiwuhanrenyuanmiqie", "否")
	u.Add("shifoujiankangqingkuang", "是")
	u.Add("shifoujiechuguohubeihuoqitayou", "否")
	u.Add("shifouweigelirenyuan", "否")
	u.Add("shentishifouyoubushizhengzhuan", "否")
	u.Add("shifouyoufare", "否")

	u.Add("tiwen", "36.6")
<<<<<<< HEAD
	u.Add("tiwen1", "36.2")
	u.Add("tiwen2", "36.2")
=======
	u.Add("tiwen1", "36.6")
	u.Add("tiwen2", "36.6")
>>>>>>> 73f7a7a4b3134b68b63cd66b224bccbcfa969cc1

	u.Add("fanhuididian", "")
	u.Add("qitaxinxi", "")

	u.Add("id", p.id)
	return u
}

// 根据账号密码模拟登陆获取Cookie
func getCookie(account, passwd string) []*http.Cookie {
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

// 获取用户唯一标识符
func getSpID(cookies []*http.Cookie) (spID string) {
	req, err := http.NewRequest("POST", htmlUrl, nil)
	if err != nil {
		log.Fatalf("获取表单页出错！")
		return
	}
	// 将所需的Cookie加入http头
	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("客户端请求失败！")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".footReturn > input").Eq(0).Each(func(i int, selection *goquery.Selection) {
		spID, _ = selection.Attr("value")
	})
	return
}
