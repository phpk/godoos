package search

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func (baidu *Baidu) Search() []Entity {
	baidu.Req.url = baidu.urlWrap()
	log.Printf("req.url: %s\n", baidu.Req.url)
	resp := &Resp{}
	resp, _ = baidu.send()
	baidu.resp = *resp
	return baidu.toEntityList()
}

func (baidu *Baidu) urlWrap() (url string) {
	return fmt.Sprintf(BaiduSearch, baidu.Req.Q)
}

func (baidu *Baidu) toEntityList() []Entity {
	var resList []Entity

	if baidu.resp.doc != nil {
		var wg sync.WaitGroup
		results := make(chan Entity, SearchNum)
		done := make(chan bool)

		baidu.resp.doc.Find("div[srcid]").EachWithBreak(func(i int, s *goquery.Selection) bool {
			if i >= SearchNum { // 如果已经处理了5篇，停止遍历
				return false
			}
			// For each item found, get the Title
			title := s.Find("h3").Find("a").Text()
			url := s.AttrOr("mu", "")
			tpl := s.AttrOr("tpl", "")
			if tpl != "se_com_default" {
				return true
			}
			subTitle := s.Find(".c-gap-top-small").Find("span").Text()
			entity := Entity{From: BaiduFrom}
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

func (baidu *Baidu) send() (resp *Resp, err error) {

	client := &http.Client{
		Transport: tr,
	}
	//提交请求
	request, err := http.NewRequest("GET", baidu.urlWrap(), nil)
	if err != nil {
		log.Println(err)
	}

	//增加header选项
	request.Header.Add("User-Agent", UserAgent)
	request.Header.Add("Host", BaiduDomain)
	request.Header.Add("Cookie", BaiduCookie)
	request.Header.Add("Accept", BaiduAccept)

	return SendDo(client, request)

}
