package main

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/guotie/gogb2312"
)

func main() {
	r := Request{
		URL: "https://www.x23us.com/class/1_1.html",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36",
		},
	}

	res := r.GetHTML()

	output, _, _, _ := gogb2312.ConvertGB2312(res)

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(output))
	Panic(err)

	lists := doc.Find("table").Find("tr")

	lists.Each(func(i int, el *goquery.Selection) {
		bookName := el.Find("td").First().Find("a").Last().Text()
		bookID, _ := el.Find("td").First().Find("a").First().Attr("href")

		reBookID := regexp.MustCompile(`/(\d+)$`)

		arr := reBookID.FindStringSubmatch(bookID)

		if len(arr) == 2 {
			bookID = arr[1]
		}

		fmt.Println(bookName, bookID)

	})

}
