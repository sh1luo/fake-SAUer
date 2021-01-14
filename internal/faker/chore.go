package faker

import (
	"fake-SAUer/config"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

// GetUUID ...
func (f *Faker) GetUUID() {
	var wg sync.WaitGroup
	var done int8
	req, _ := http.NewRequest("POST", htmlUrl, nil)

	for i := 0; i < f.Cnt; i++ {

		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			if f.Cf.StusInfos[i].UUID != 0 {
				fmt.Printf("用户%s无需再获取UUID\n", f.Cf.StusInfos[i].Name)
				return
			}

			cks := GetCookie(f.Cf.StusInfos[i].Account, f.Cf.StusInfos[i].Passwd)
			for _, c := range cks {
				req.AddCookie(c)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return
			}

			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				return
			}
			doc.Find(".footReturn > input").Eq(0).Each(func(i int, selection *goquery.Selection) {
				sid, _ := selection.Attr("value")
				iid, err := strconv.Atoi(sid)
				if err != nil || iid == 0 {
					fmt.Println(err)
					return
				}
				f.mu.Lock()
				f.Cf.StusInfos[i].UUID = iid

				fmt.Println("map操作：", f.Cf.StusInfos[i].Name, "=", iid)

				f.mu.Unlock()

				fmt.Printf("用户%s获取UUID成功! UUID=%d\n", f.Cf.StusInfos[i].Name, iid)
				done++
			})
		}(i)
	}

	wg.Wait()
	fmt.Printf("GetUUID完毕，一共%d个用户，本次重新获取了%d个用户的UUID\n", f.Cnt, done)
	return
}

// 根据账号密码模拟登陆获取Cookie
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
func bindInfo(f *config.StuInfo, iid int) url.Values {
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

	u.Add("tiwen", "36.6")
	u.Add("tiwen1", "36.6")
	u.Add("tiwen2", "36.6")

	u.Add("fanhuididian", "")
	u.Add("qitaxinxi", "")

	sid := strconv.Itoa(iid)
	u.Add("id", sid)
	return u
}

func checkInfo(infos []*config.StuInfo) int {
	var nums int
	for _, stu := range infos {
		if stu.Name != "" && stu.Phone != "" && stu.City != "" &&
			stu.Province != "" && stu.Account != "" && stu.Passwd != "" && stu.College != "" {
			nums++
		}
	}
	return nums
}
