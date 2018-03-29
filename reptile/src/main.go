package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/djimenez/iconv-go"
)

func main() {

	// r := Request{
	// 	URL: "http://www.23us.so/files/article/html/10/10839/3652687.html",
	// 	Headers: map[string]string{
	// 		"Host":       "www.23us.so",
	// 		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36",
	// 		"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
	// 	},
	// 	proxy: "http://127.0.0.1:80",
	// }
	data := make(url.Values)
	data.Add("searchkey", "剑来")
	data.Add("searchtype", "articlename")
	data.Add("action", "login")
	data.Add("submit", "&#160;搜&#160;&#160;索&#160;")

	res, err := http.PostForm("http://www.biquge.vip/modules/article/search.php", data)

	if err != nil {
		log.Panic(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Panic(err)
	}

	body, err = iconv.ConvertString(string(body), "utf-8", "GB2312")
	fmt.Printf(string(body))
}
