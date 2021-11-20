package core

import (
	"fake-SAUer/conf"
	"fmt"
	"log"
	"net/http"
)

func StartHTTPServer(faker *Faker) {
	log.Printf("Starting HTTP Server at 'localhost:9000'...\n")
	http.Handle("/", faker)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatalf("HTTP Server Start error: %s\n", err.Error())
	}
}

func (f *Faker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := ""
	for i := 0; i < f.Cnt; i++ {
		resp += fmt.Sprintf("%d\t%s\t%s\t%s\n", i, conf.GlobalConfig.StusInfo[i].Name, conf.GlobalConfig.StusInfo[i].Account, conf.GlobalConfig.StusInfo[i].College)
	}
	w.Write([]byte(resp))
}

//func (f *Faker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	switch r.Method {
//	// update uuid forcibly, like `?uuid=xxx&username=xxx&password=xxx`
//	case http.MethodPut:
//		uuid := r.URL.Query().Get("uuid")
//		username := r.URL.Query().Get("username")
//		passwd := r.URL.Query().Get("passwd")
//		for _, s := range main.G_Conf.StusInfo {
//			if s.Account == username && s.Passwd == passwd {
//				s.Uuid = uuid
//				if _, err := w.Write([]byte("uuid设置成功")); err != nil {
//					log.Printf("write msg err:%s\n", err.Error())
//					return
//				}
//			}
//		}
//		w.WriteHeader(http.StatusBadRequest)
//	// add a user
//	case http.MethodPost:
//		if f.EnableHTTP {
//			bs, err := ioutil.ReadAll(r.Body)
//			if err != nil || len(bs) == 0 {
//				w.WriteHeader(http.StatusBadRequest)
//				return
//			}
//			defer r.Body.Close()
//			s := &main.StuInfo{}
//			if err = json.Unmarshal(bs, s); err != nil {
//				w.WriteHeader(http.StatusInternalServerError)
//				log.Printf("unmarshal err:%s\n", err.Error())
//				return
//			}
//			main.G_Conf.StusInfo = append(main.G_Conf.StusInfo, s)
//			f.Cnt = len(main.G_Conf.StusInfo)
//			w.WriteHeader(http.StatusOK)
//			if _, err = w.Write([]byte("Add successfully!")); err != nil {
//				log.Printf("add new student err:%s\n", err.Error())
//				return
//			}
//		} else {
//			w.WriteHeader(http.StatusForbidden)
//			if _, err := w.Write([]byte("This Service had prohibited registration!")); err != nil {
//				log.Printf("write msg err:%s\n", err.Error())
//				return
//			}
//		}
//	// switch service status
//	case http.MethodPatch:
//		if u, p, ok := r.BasicAuth(); ok && u == Username && p == Password {
//			if r.URL.Path == "/switch" {
//				f.EnableHTTP = !f.EnableHTTP
//				w.WriteHeader(http.StatusOK)
//				if _, err := w.Write([]byte("Switch successfully")); err != nil {
//					log.Printf("write msg err:%s\n", err.Error())
//					return
//				}
//			} else {
//				w.WriteHeader(http.StatusNotFound)
//			}
//		} else {
//			w.WriteHeader(http.StatusForbidden)
//		}
//	default:
//		w.WriteHeader(http.StatusMethodNotAllowed)
//		if _, err := w.Write([]byte("Method Not Allowed")); err != nil {
//			log.Printf("add new student err:%s\n", err.Error())
//			return
//		}
//	}
//}