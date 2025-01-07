package search

import (
	"encoding/json"
	"godo/libs"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"jaytaylor.com/html2text"
)

const (
	SearchNum = 5
	Debug     = false
	UserAgent = "user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"
	ProxyURL  = "http://127.0.0.1:7890"
	ProxyOpen = false

	// =================BAIDU==================
	//BaiduCookie = ""
	BaiduCookie = "BIDUPSID=3535F2B8A915447A4839A7DD194BA7B3; PSTM=1617718622; __yjs_duid=1_76a86308fbc1b29f8b2c852e8c9e24fd1620886035244; BD_UPN=123253; MCITY=-%3A; BDUSS=m5TRVBaQjk0NXF0SFpmMHZtN2ZWSTlyZVo2d0ZJTVhlU1BMTUpqbXRrZEl0WFZqSUFBQUFBJCQAAAAAAAAAAAEAAAC8QhUx1MbUqtauwfpoaQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEgoTmNIKE5jc; BDUSS_BFESS=m5TRVBaQjk0NXF0SFpmMHZtN2ZWSTlyZVo2d0ZJTVhlU1BMTUpqbXRrZEl0WFZqSUFBQUFBJCQAAAAAAAAAAAEAAAC8QhUx1MbUqtauwfpoaQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEgoTmNIKE5jc; BDSFRCVID=rY_OJeC62lCrte6jotU8bVRNE2SdnBRTH6aotxm4whxuChbLecyMEG0Pyf8g0KubiKd_ogKK0eOTHktF_2uxOjjg8UtVJeC6EG0Ptf8g0f5; H_BDCLCKID_SF=JJkO_D_atKvDqTrP-trf5DCShUFsyMvlB2Q-XPoO3KJnMU74M4Qtj4Fp3-ON--QiW5cpoMbgylRp8P3y0bb2DUA1y4vp5MnqQeTxoUJ2fnRJEUcGqj5Ah--ebPRiJPQ9QgbWLpQ7tt5W8ncFbT7l5hKpbt-q0x-jLTnhVn0MBCK0MC09j6KhDTPVKgTa54cbb4o2WbCQQq3O8pcN2b5oQT81jnnHatoh32jZ-pb4bb5vOPQKDpOUWfAkXpJvQnJjt2JxaqRCKhv-Sl5jDh3Me-AsLn6te6jzaIvy0hvctn5cShnc5MjrDRLbXU6BK5vPbNcZ0l8K3l02V-bIe-t2XjQhDNKtt5_jJJIsBP_8aJ7bHn7gbJK_-P4DeP6wexRZ5mAqoDQ6tbbnHRcH3noJ-lKw3MTA-55lba6naIQqa-cKDR38Kp58hPT00tb8Qnb43bRT2PPy5KJvfJo1BTrYhP-UyN3-Wh3725nlMKoaMp78jR093JO4y4Ldj4oxJpOJ5JbMonLafDD3fb7kbP6Eq4D_MfOtetJyaR3OoM7vWJ5WqR7jD5DbMP4qQfcrQxTn3Nrfb-OIJPTKShbXKxonKlLObbOGLM6EfNcZ0l8K3l02V-bRDDcfQJQDQt7JKPRMW20j0l7mWnvMsxA45J7cM4IseboJLfT-0bc4KKJxbnLWeIJEjjCKejbQDGADq6nfb5kXLn7J-J5HfbTkbITjhPrMKRrdWMT-0bFHWpO-bnK-jR7m3j7D3lDNQH5gKP6GLHn7_JjObPnVM-3d2boTM-DTbfb8WxQxtNR80DnjtpvhHRTobhnobUPUDMo9LUvWbgcdot5yBbc8eIna5hjkbfJBQttjQn3hfIkj0DKLK-oj-D8lej7P; BAIDUID=9A9CC424161261B757A36C59CFED76E5:FG=1; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; H_PS_PSSID=36542_37557_37519_36884_37627_36786_37539_37500_26350_37343_37461; BAIDUID_BFESS=9A9CC424161261B757A36C59CFED76E5:FG=1; BDSFRCVID_BFESS=rY_OJeC62lCrte6jotU8bVRNE2SdnBRTH6aotxm4whxuChbLecyMEG0Pyf8g0KubiKd_ogKK0eOTHktF_2uxOjjg8UtVJeC6EG0Ptf8g0f5; H_BDCLCKID_SF_BFESS=JJkO_D_atKvDqTrP-trf5DCShUFsyMvlB2Q-XPoO3KJnMU74M4Qtj4Fp3-ON--QiW5cpoMbgylRp8P3y0bb2DUA1y4vp5MnqQeTxoUJ2fnRJEUcGqj5Ah--ebPRiJPQ9QgbWLpQ7tt5W8ncFbT7l5hKpbt-q0x-jLTnhVn0MBCK0MC09j6KhDTPVKgTa54cbb4o2WbCQQq3O8pcN2b5oQT81jnnHatoh32jZ-pb4bb5vOPQKDpOUWfAkXpJvQnJjt2JxaqRCKhv-Sl5jDh3Me-AsLn6te6jzaIvy0hvctn5cShnc5MjrDRLbXU6BK5vPbNcZ0l8K3l02V-bIe-t2XjQhDNKtt5_jJJIsBP_8aJ7bHn7gbJK_-P4DeP6wexRZ5mAqoDQ6tbbnHRcH3noJ-lKw3MTA-55lba6naIQqa-cKDR38Kp58hPT00tb8Qnb43bRT2PPy5KJvfJo1BTrYhP-UyN3-Wh3725nlMKoaMp78jR093JO4y4Ldj4oxJpOJ5JbMonLafDD3fb7kbP6Eq4D_MfOtetJyaR3OoM7vWJ5WqR7jD5DbMP4qQfcrQxTn3Nrfb-OIJPTKShbXKxonKlLObbOGLM6EfNcZ0l8K3l02V-bRDDcfQJQDQt7JKPRMW20j0l7mWnvMsxA45J7cM4IseboJLfT-0bc4KKJxbnLWeIJEjjCKejbQDGADq6nfb5kXLn7J-J5HfbTkbITjhPrMKRrdWMT-0bFHWpO-bnK-jR7m3j7D3lDNQH5gKP6GLHn7_JjObPnVM-3d2boTM-DTbfb8WxQxtNR80DnjtpvhHRTobhnobUPUDMo9LUvWbgcdot5yBbc8eIna5hjkbfJBQttjQn3hfIkj0DKLK-oj-D8lej7P; delPer=0; BD_CK_SAM=1; PSINO=3; BA_HECTOR=8581ala48k0la181aga50fm81hm8prs1e; ZFY=MW5tgLZfxnXk8HLzWL5vKMTfHTBgdsBmZudM7p1:BE:Aw:C; sugstore=1; H_PS_645EC=6985pYm8uvzjRIBPLtrwALaCau4%2Fxj7uefaYB7sGfR5qYLEGfRZ8Eyb%2F1oc"
	BaiduUrl    = "https://www.baidu.com"
	BaiduDomain = "www.baidu.com"
	BaiduAccept = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	BaiduSearch = BaiduUrl + "/s?wd=%s" + "&usm=3&rsv_idx=2&rsv_page=1"
	BaiduFrom   = "百度"

	// =================BING==================
	BingCoolkie = ""
	//BingCoolkie = "MUID=16AE9B2FE9DC6568044F8B3EE8F2640C; MUIDB=16AE9B2FE9DC6568044F8B3EE8F2640C; SRCHD=AF=BDVEHC; SRCHUID=V=2&GUID=7321ED22410B41459C503CB6D2628196&dmnchg=1; _UR=QS=0&TQS=0; imgv=flts=20220704; _tarLang=default=zh-Hans; _TTSS_IN=hist=WyJlbiIsImF1dG8tZGV0ZWN0Il0=; _TTSS_OUT=hist=WyJ6aC1IYW5zIl0=; _HPVN=CS=eyJQbiI6eyJDbiI6MjQsIlN0IjoyLCJRcyI6MCwiUHJvZCI6IlAifSwiU2MiOnsiQ24iOjI0LCJTdCI6MCwiUXMiOjAsIlByb2QiOiJIIn0sIlF6Ijp7IkNuIjoyNCwiU3QiOjEsIlFzIjowLCJQcm9kIjoiVCJ9LCJBcCI6dHJ1ZSwiTXV0ZSI6dHJ1ZSwiTGFkIjoiMjAyMi0wOS0wNlQwMDowMDowMFoiLCJJb3RkIjowLCJHd2IiOjAsIkRmdCI6bnVsbCwiTXZzIjowLCJGbHQiOjAsIkltcCI6Nzd9; ANIMIA=FRE=1; MMCASM=ID=3FEB6F3855CC49E584F2DA61F6E5E44C; ZHCHATSTRONGATTRACT=TRUE; _SS=SID=11649FAAD94F6D9716388DF6D8296CB4&PC=U316; SRCHS=PC=U316; ABDEF=V=13&ABDV=13&MRB=0&MRNB=1668441959593; SUID=M; _EDGE_S=SID=11649FAAD94F6D9716388DF6D8296CB4&ui=zh-cn; SRCHUSR=DOB=20210406&T=1668472757000&TPC=1668472758000; ZHCHATWEAKATTRACT=TRUE; ipv6=hit=1668476360196&t=4; ZHLASTACTIVECHAT=0; ZHSEARCHCHATSTATUS=STATUS=0; SNRHOP=I=&TS=; RECSEARCH=SQs=[{\"q\":\"giac%20%E4%B8%8A%E6%B5%B7\"%2C\"c\":1%2C\"ad\":true}%2C{\"q\":\"sessioncachesize\"%2C\"c\":1%2C\"ad\":false}%2C{\"q\":\"rset\"%2C\"c\":1%2C\"ad\":false}%2C{\"q\":\"xx%3Ainitiatingheapoccupancypercent\"%2C\"c\":1%2C\"ad\":false}%2C{\"q\":\"%E7%AE%A1%E7%90%86%E7%9A%84%E5%B8%B8%E8%AF%86\"%2C\"c\":1%2C\"ad\":true}%2C{\"q\":\"%E7%AE%A1%E7%90%86%E7%9A%84%E5%B8%B8%E8%AF%86%20%E5%BE%B7%E9%B2%81%E5%85%8B\"%2C\"c\":1%2C\"ad\":true}%2C{\"q\":\"%E7%AE%A1%E7%90%86\"%2C\"c\":1%2C\"ad\":true}%2C{\"q\":\"goquery%20%E5%BE%AA%E7%8E%AF\"%2C\"c\":1%2C\"ad\":false}%2C{\"q\":\"yuanbiguo\"%2C\"c\":1%2C\"ad\":false}]; SRCHHPGUSR=SRCHLANGV2=zh-Hans&BRW=W&BRH=S&CW=1396&CH=435&DPR=2&UTC=480&DM=0&WTS=63804069557&HV=1668474094&BZA=0&SRCHLANG=zh-Hans&SW=1440&SH=900&PV=11.2.3&EXLTT=6&SCW=1381&SCH=1408&PRVCW=1396&PRVCH=764"
	BingUrl    = "https://cn.bing.com"
	BingDomain = "cn.bing.com"
	BingAccept = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	BingSearch = BingUrl + "/search?q=%s" + "&PC=U316&FORM=CHROMN"
	BingFrom   = "Bing"

	// =================Google==================
	//GoogleCookie = ""
	GoogleCookie = "CONSENT=YES+srp.gws-20220523-0-RC1.zh-CN+FX; HSID=AIDFOIfZRXMjhRznJ; SSID=ATYvcUlmXr3mPFerl; APISID=fTLh3HcYiZa0Ch2l/AiukySr7MDg_GhSRo; SAPISID=-JpKVkIJZXDgucyp/AjPiMHGTZYbtqrQbt; __Secure-1PAPISID=-JpKVkIJZXDgucyp/AjPiMHGTZYbtqrQbt; __Secure-3PAPISID=-JpKVkIJZXDgucyp/AjPiMHGTZYbtqrQbt; SEARCH_SAMESITE=CgQIk5YB; SID=Qggrdf-HuR3Im3foqYJWVZDPEay8b5U0-O_E6L489Ppxqy2bY358yp5YaaIs5NjEgyW23g.; __Secure-1PSID=Qggrdf-HuR3Im3foqYJWVZDPEay8b5U0-O_E6L489Ppxqy2brNnvCnxg8UCPQbHqZClMmA.; __Secure-3PSID=Qggrdf-HuR3Im3foqYJWVZDPEay8b5U0-O_E6L489Ppxqy2bJ0s1kXYhOOAgur0ax0qp6w.; OTZ=6777442_24_24__24_; AEC=AakniGNW8IBpZgW_-IGOh45Otu8fCK-Ty1n0eSZhgs05euV7s2hhmWLhCQ; NID=511=k96f3xX7KMB1Wvhg478KURvWSAd53Y3rTkVG3bMb1FhtdJwQUbUiHJOnVTgFgW0Mad_1X5gKnWLMt0lPHl33nQdVCTzEiitFOC2dIicYusLP1zl_L6Wh9l-XO6x9VRh4ZxGlcu1bCmlEbBYvz2eL2ioCJ3RMGZtGcr6_1tdpGZH8DRcj3c8X6FxRJqcX5peACa9pGELYmZk4TfOYiUON7p4ht5MTkfA4-hmLUD_JQOT_bK4Aub5SzhXukCFBt5qx2UMkKKtheWlK6cRlR-EIqbC0ppccy2pUdT0YyfE; DV=c77sl5f1j0hXACm1QWzNqJ3oPLjLSRge9TXKlwa5UAEAAMBLUjyomhJKegAAABD63jYjJ8URJwAAAKJBqsVeeJWmDgAAAA; 1P_JAR=2022-11-22-00; SIDCC=AIKkIs0Jk8Fs8wt-yHi9yEM8vFzc6y1bmzVm7mYb3eELJ0t_5yixpCQToLaIDwlOnk3yTvFwKw; __Secure-1PSIDCC=AIKkIs2Sk9zwJs1YfCaHpcFU4ZAKHzDO7Atrzk3b-uXf384Xuc1nQ0U41MiX_VQcEBrJ0TSR1ZY; __Secure-3PSIDCC=AIKkIs1iF4en_Ln8Jonhw8Q79QB_dwbKwcw0khXXIX0vsIwOjBgxmO8dd6kfi10evURBz7YdDg"
	GoogleUrl    = "https://www.google.com"
	GoogleDomain = "www.google.com"
	GoogleAccept = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	GoogleSearch = GoogleUrl + "/search?q=%s" + "&ie=UTF-8"
	GoogleFrom   = "Google"

	// =================Wx==================
	WxCookie = ""
	//WxCookie = "SUID=F84DCC781539960A000000006235D4E8; SUV=1647695080857391; ssuid=6792259580; weixinIndexVisited=1; IPLOC=CN3100; ABTEST=0|1668775030|v1; JSESSIONID=aaa2PBWFjHNS8hQ8_tfpy; cd=1668913493&0f942166ea05ede01cfe88195d36508d; rd=tyllllllll20WBOSYTuBqQ2iuqV0WBOqAfJbLZllll9llllxVllll5@@@@@@@@@@; ld=6Zllllllll20WBOSYTuBqQ2DdNH0WBOqAfJbLZllll9lllllVklll5@@@@@@@@@@; LSTMV=217%2C66; LCLKINT=1482; PHPSESSID=udut3pen5cml0b9849o68jch40; SNUID=9D60E0542D28C24D1776D6A42DF58CE7; ariaDefaultTheme=undefined"
	WxUrl    = "https://weixin.sogou.com"
	WxDomain = "weixin.sogou.com"
	WxAccept = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	WxSearch = WxUrl + "/weixin?type=2&s_from=input&query=%s&ie=utf8&_sug_=n&_sug_type_="
	WxFrom   = "微信公众号"
)

type SearchEngine interface {
	Search() (result []Entity)
	urlWrap() (url string)
	toEntityList() (entityList []Entity)
	send() (resp *Resp, err error)
}

type Baidu struct {
	Req  Req
	resp Resp
}

type Bing struct {
	Req  Req
	resp Resp
}

type Req struct {
	Q         string
	url       string
	userAgent string
	http.Cookie
}

type Resp struct {
	code int
	body string
	doc  *goquery.Document
}

type Entity struct {
	Title    string `json:"title"`
	Host     string `json:"host"`
	Url      string `json:"url"`
	SubTitle string `json:"subTitle"`
	From     string `json:"from"`
	Content  string `json:"content"`
}

func SearchWebhandler(w http.ResponseWriter, r *http.Request) {
	//log.Println(r.URL)
	_ = r.ParseForm()
	q := r.Form.Get("q")
	if q == "" {
		libs.ErrorMsg(w, "查询内容不能为空")
		return
	}
	q = url.QueryEscape(q)
	//log.Printf("查询内容:%s\n", q)
	jsonResult := SearchWeb(q)

	// 构造返回
	body, err := json.Marshal(jsonResult)
	if err != nil {
		libs.ErrorMsg(w, "解析错误")
		return
	}
	libs.SuccessMsg(w, body, "success")
}
func SearchWeb(q string) []Entity {

	q = url.QueryEscape(q)
	//log.Printf("查询内容:%s\n", q)
	jsonResult := []Entity{}

	array := [...]SearchEngine{
		&Bing{Req: Req{Q: q}},
		&Baidu{Req: Req{Q: q}},
	}

	var wg sync.WaitGroup
	c := make(chan []Entity, len(array))

	for _, engine := range array {
		wg.Add(1)
		go func(engine SearchEngine) {
			defer wg.Done()
			result := engine.Search()
			c <- result
		}(engine)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for result := range c {
		if result != nil {
			jsonResult = append(jsonResult, result...)
		}
	}

	return jsonResult
}
func SendDo(client *http.Client, request *http.Request) (*Resp, error) {
	resp := &Resp{code: 200}

	//处理返回结果
	response, e := client.Do(request)
	if response == nil {
		resp.code = 500
		log.Printf("response nil: %v\n", e)
		return resp, nil
	}
	if response.StatusCode != 200 {
		resp.code = response.StatusCode
		log.Printf("status code error: %d %s\n", response.StatusCode, response.Status)
		return resp, nil
	}
	defer response.Body.Close()
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Println(err)
	}

	resp.code = 200
	resp.doc = doc

	return resp, nil
}
func FetchContent(url string) (content string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return "", err
	}
	text, err := html2text.FromString(string(body), html2text.Options{PrettyTables: false})
	if err != nil {
		return "", err
	}
	// if len(text) > 1000 {
	// 	text = text[:1000]
	// }
	return text, nil
}
