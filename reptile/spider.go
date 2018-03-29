package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

// Spider 爬虫数据类型
type Spider struct {
	url     string
	headers map[string]string
}

// getHTML 抓取 html
func (sp Spider) getHTML() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", sp.url, nil)
	if err != nil {
		log.Panic(err)
	}

	for key, value := range sp.headers {
		req.Header.Add(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	return string(body)
}

func parse() {
	header := map[string]string{
		"Host":                      "movie.douban.com",
		"Connection":                "keep-alive",
		"Cache-Control":             "max-age=0",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		"Referer":                   "https://movie.douban.com/top250",
	}
	//创建excel文件
	f, err := os.Create("./test.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	//写入标题
	f.WriteString("电影名称,评分评价人数" + "\t" + "\r\n")

	//循环每页解析并把结果写入excel
	for i := 0; i < 10; i++ {
		fmt.Println("正在抓取第" + strconv.Itoa(i) + "页......")
		url := "https://movie.douban.com/top250?start=" + strconv.Itoa(i*25) + "&filter="
		spider := &Spider{url, header}
		html := spider.getHTML()

		//评价人数
		pattern2 := `<span>(.*?)评价</span>`
		rp2 := regexp.MustCompile(pattern2)
		findTxt2 := rp2.FindAllStringSubmatch(html, -1)

		//评分
		pattern3 := `property="v:average">(.*?)</span>`
		rp3 := regexp.MustCompile(pattern3)
		findTxt3 := rp3.FindAllStringSubmatch(html, -1)

		//电影名称
		pattern4 := `img width="100" alt="(.*?)" src=`
		rp4 := regexp.MustCompile(pattern4)
		findTxt4 := rp4.FindAllStringSubmatch(html, -1)

		// f.WriteString("\xEF\xBB\xBF")
		//  打印全部数据和写入excel文件

		print(findTxt3)
		for i := 0; i < len(findTxt2); i++ {
			fmt.Printf("%s %s %s\n", findTxt4[i][1], findTxt3[i][1], findTxt2[i][1])
			f.WriteString(findTxt4[i][1] + "," + findTxt3[i][1] + "," + findTxt2[i][1] + "\t" + "\r\n")

		}
	}

}
