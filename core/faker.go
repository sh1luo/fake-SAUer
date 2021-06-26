package core

import (
	"encoding/json"
	"errors"
	"fake-SAUer/conf"
	_ "fake-SAUer/conf"
	"fake-SAUer/notice"
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
	postUrl   = "https://ucapp.sau.edu.cn/wap/login/invalid"             // log in post target
	htmlUrl   = "https://app.sau.edu.cn/form/wap/default?formid=10"      // html url to get UUID
	submitURL = "https://app.sau.edu.cn/form/wap/default/save?formid=10" // submit address
	
	Username = "YOUR USERNAME"
	Password = "YOUR PASSWORD"
)

type Faker struct {
	Cnt        int // punch counts
	Notifier   notice.Notifier
	EnableHTTP bool
	mu         sync.Mutex // protect ...
}

func NewFaker(enableHTTP bool) (*Faker, error) {
	if len(conf.G_Conf.StusInfos) == 0 {
		return nil, errors.New("没有有效的学生信息")
	}
	f := &Faker{}
	f.Notifier = notice.NewNotifier("email", "3450047248@qq.com", "lgmjasugxqefchjj", "smtp.qq.com", 465)
	f.Cnt = len(conf.G_Conf.StusInfos)
	f.EnableHTTP = enableHTTP
	return f, nil
}

// Do 执行任务,返回打卡成功数量
func (f *Faker) Do() (done int8) {
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
	h.Set("sec-ch-ua-mobile", "?0")
	h.Set("Referer", "https://app.sau.edu.cn/form/wap/default/index?formid=10&nn=4669.797748311082")
	h.Set("Sec-Fetch-Dest", "empty")
	h.Set("Sec-Fetch-Mode", "cors")
	h.Set("Sec-Fetch-Site", "same-origin")
	h.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
	h.Set("X-Requested-With", "XMLHttpRequest")
	
	for i := 0; i < f.Cnt; i++ {
		cks := GetCookie(conf.G_Conf.StusInfos[i].Account, conf.G_Conf.StusInfos[i].Passwd)
		sid := ""
		if conf.G_Conf.StusInfos[i].Uuid != "" {
			sid = conf.G_Conf.StusInfos[i].Uuid
		} else {
			sid = GetUuid(cks)
		}
		v := bindInfo(conf.G_Conf.StusInfos[i], sid)
		req, err := http.NewRequest("POST", submitURL, strings.NewReader(v.Encode()))
		if err != nil {
			panic("致命错误，构造POST表单失败！")
		}
		
		req.Header = h
		req.Header.Set("Content-Length", strconv.Itoa(len(v.Encode())))
		fmt.Println(sid, strconv.Itoa(len(v.Encode())))
		for _, c := range cks {
			req.AddCookie(c)
		}
		
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("POST - %s err: %s\n", submitURL, err.Error())
			return
		}
		
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("read submit.Resp.Body err: %s\n", err)
			return
		}
		
		if e := gjson.Get(string(data), "e").Int(); e == 0 {
			done++
		} else {
			if f.Notifier != nil {
				if err := f.Notifier.Notice(conf.G_Conf.StusInfos[i].To, "Punch Message", "今日打卡失败"); err != nil {
					log.Printf("通知失败: %s\n", err)
				}
			}
		}
	}
	
	return done
}

func (f *Faker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		uuid := r.URL.Query().Get("uuid")
		username := r.URL.Query().Get("username")
		passwd := r.URL.Query().Get("passwd")
		for _, s := range conf.G_Conf.StusInfos {
			if s.Account == username && s.Passwd == passwd {
				s.Uuid = uuid
				if _, err := w.Write([]byte("uuid设置成功")); err != nil {
					log.Printf("write msg err:%s\n", err.Error())
					return
				}
			}
		}
		w.WriteHeader(http.StatusBadRequest)
	case http.MethodPost:
		if f.EnableHTTP {
			bs, err := ioutil.ReadAll(r.Body)
			if err != nil || len(bs) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			defer r.Body.Close()
			s := &conf.StuInfo{}
			if err = json.Unmarshal(bs, s); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("unmarshal err:%s\n", err.Error())
				return
			}
			conf.G_Conf.StusInfos = append(conf.G_Conf.StusInfos, s)
			f.Cnt = len(conf.G_Conf.StusInfos)
			w.WriteHeader(http.StatusOK)
			if _, err = w.Write([]byte("Add successfully!")); err != nil {
				log.Printf("add new student err:%s\n", err.Error())
				return
			}
		} else {
			w.WriteHeader(http.StatusForbidden)
			if _, err := w.Write([]byte("This Service had prohibited registration!")); err != nil {
				log.Printf("write msg err:%s\n", err.Error())
				return
			}
		}
	case http.MethodPatch:
		if u, p, ok := r.BasicAuth(); ok && u == Username && p == Password {
			if r.URL.Path == "/switch" {
				f.EnableHTTP = !f.EnableHTTP
				w.WriteHeader(http.StatusOK)
				if _, err := w.Write([]byte("Switch successfully")); err != nil {
					log.Printf("write msg err:%s\n", err.Error())
					return
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := w.Write([]byte("Method Not Allowed")); err != nil {
			log.Printf("add new student err:%s\n", err.Error())
			return
		}
	}
}
