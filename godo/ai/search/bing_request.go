package search

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func (bing *Bing) Search() []Entity {
	bing.Req.url = bing.urlWrap()
	log.Printf("req.url: %s\n", bing.Req.url)
	resp := &Resp{}
	resp, _ = bing.send()
	bing.resp = *resp
	return bing.toEntityList()
}

func (bing *Bing) urlWrap() (url string) {
	return fmt.Sprintf(BingSearch, bing.Req.Q)
}

func (bing *Bing) toEntityList() []Entity {
	var resList []Entity
	if bing.resp.doc != nil {
		var wg sync.WaitGroup
		results := make(chan Entity, SearchNum)
		done := make(chan bool)

		bing.resp.doc.Find("ol#b_results>li[class=b_algo]").EachWithBreak(func(i int, s *goquery.Selection) bool {
			if len(resList) >= SearchNum { // 如果已经处理了5篇，停止遍历
				return false
			}

			// For each item found, get the Title
			title := s.Find("div[class=b_title]>h2>a").Text()
			url := s.Find("div[class=b_title]>h2>a").AttrOr("href", "")
			subTitle := s.Find("div[class=b_caption]>p").Text()
			entity := Entity{From: BingFrom}
			entity.Title = title
			entity.SubTitle = subTitle
			entity.Url = url
			host := strings.ReplaceAll(url, "http://", "")
			host = strings.ReplaceAll(host, "https://", "")
			entity.Host = strings.Split(host, "/")[0]

			wg.Add(1)
			go func(entity Entity) {
				defer wg.Done()
				// 获取详细页面的纯文本内容
				content, err := FetchContent(entity.Url)
				if err != nil {
					log.Printf("Error fetching content for URL %s: %v\n", entity.Url, err)
					results <- entity
					return
				}
				entity.Content = content
				results <- entity
			}(entity)

			return true
		})

		go func() {
			wg.Wait()
			close(results)
			done <- true
		}()

		for entity := range results {
			resList = append(resList, entity)
		}

		<-done
	}
	return resList
}

func (bing *Bing) send() (resp *Resp, err error) {

	client := &http.Client{
		Transport: tr,
	}
	//提交请求
	request, err := http.NewRequest("GET", bing.urlWrap(), nil)
	if err != nil {
		log.Println(err)
	}

	//增加header选项
	request.Header.Add("User-Agent", UserAgent)
	request.Header.Add("Host", BingDomain)
	request.Header.Add("Cookie", BingCoolkie)
	request.Header.Add("Accept", BingAccept)

	return SendDo(client, request)
}
