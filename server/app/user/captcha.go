package user

import (
	"encoding/json"
	"fmt"
	"godocms/common"
	"godocms/libs"
	"godocms/middleware"
	"godocms/pkg/sms/alibaba"
	"godocms/utils"
	"log"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	images "github.com/wenlng/go-captcha-assets/resources/images_v2"
	"github.com/wenlng/go-captcha-assets/resources/tiles"
	"github.com/wenlng/go-captcha/v2/slide"
)

func init() {
	buildSlideCaptcha()
	middleware.RegisterRouter("POST", "user/smscode", sendSmsCodeHandler, 0, "发送手机验证码")
	middleware.RegisterRouter("POST", "user/emailcode", sendEmailHandler, 0, "发送邮件验证码")
	middleware.RegisterRouter("GET", "user/getcaptcha", slideCaptchaHandler, 0, "获取图形验证码")
	middleware.RegisterRouter("POST", "user/checkcaptcha", checkSlideDataHandler, 0, "校验验证码")
}

var slideCapt slide.Captcha

type SendSmsStatus struct {
	Number int       `json:"number"` // 发送次数
	Time   time.Time `json:"time"`
}

type SmsCodeRecord struct {
	Phone    string `json:"phone"`
	Sender   string `json:"sender"` // alibaba, tencent, baidu
	SendTime int64  `json:"send_time"`
	Err      string `json:"err"`
}

var (
	lastSenter = make(map[string]*SendSmsStatus) // 用于存储手机号或邮箱的上次发送时间
	mu         sync.Mutex                        // 保护 lastSenter map 的并发安全
)

func sendSmsCodeHandler(c *gin.Context) {
	type SendSmsReq struct {
		Phone    string `json:"phone" binding:"required"`
		ClientID string `json:"client_id" binding:"required"`
		// VaildStatus bool   `json:"valid_status" binding:"required"` // 验证状态, 验证成功才能发送短信
	}
	var req SendSmsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		libs.Error(c, "参数错误")
		return
	}

	// if !req.VaildStatus {
	// 	libs.Error(c, "请先验证")
	// 	return
	// }

	if !common.LoginConf.Phone.Enable {
		libs.Error(c, "短信服务未开启")
		return
	}

	// 获取当前时间
	now := time.Now()

	mu.Lock()
	defer mu.Unlock()
	// 检查该手机号是否在一分钟内发送过短信
	record, exists := lastSenter[req.Phone]
	if exists && now.Sub(record.Time) < time.Minute {
		libs.Error(c, "一分钟内只能发送一次短信")
		return
	}
	if exists {
		record.Number++
		record.Time = now
	} else {
		lastSenter[req.Phone] = &SendSmsStatus{
			Number: 1,
			Time:   now,
		}
	}

	if common.LoginConf.Phone.AliyunSms.AccessKeyId == "" || common.LoginConf.Phone.AliyunSms.AccessKeySecret == "" {
		libs.Error(c, "短信服务未配置")
		return
	}
	// TODO: 从阿里云 腾讯云 百度云 等云服务随机挑选发送短信(已配置)
	err := alibaba.SendSms(req.Phone)
	if err != nil {
		if record.Number >= 3 {
			// 关闭发送短信服务
			common.LoginConf.Phone.Enable = false
			record := SmsCodeRecord{
				Phone:    req.Phone,
				Sender:   "alibaba",
				SendTime: time.Now().Unix(),
				Err:      err.Error(),
			}
			common.SetCache("cloudsms", record, 60*time.Hour)
			lastSenter = nil
			libs.Error(c, "短信服务不可用")
			return
		}

		libs.Error(c, "发送短信失败")
		return
	}

	libs.Success(c, "发送成功", nil)
}

func sendEmailHandler(c *gin.Context) {
	type EmailReq struct {
		Email string `json:"email" binding:"required"`
	}
	var req EmailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		libs.Error(c, "邮箱信息缺失")
		return
	}
	if !common.LoginConf.Email.Enable {
		libs.Error(c, "邮箱功能未开启")
		return
	}
	if !common.EmailIsValid() {
		libs.Error(c, "邮箱配置错误")
		return
	}

	email := req.Email
	if !utils.IsEmail(email) {
		libs.Error(c, "邮箱格式不正确")
		return
	}

	now := time.Now()
	mu.Lock()
	record, exists := lastSenter[email]
	if exists && now.Sub(record.Time) < time.Minute {
		libs.Error(c, "请稍微等待再发送")
		return
	}
	if exists {
		record.Number++
		record.Time = now
	} else {
		lastSenter[email] = &SendSmsStatus{
			Number: 1,
			Time:   now,
		}
	}
	mu.Unlock()

	num := utils.GenerateRandomSixDigitNumber()
	emailKey := email + "_captcha"
	common.SetCache(emailKey, num, 10)
	log.Printf("captchaStr num:%v", num)
	body := fmt.Sprintf("您的验证码为：%s", num)
	if err := utils.Email(email, "GodoOS验证码", body); err != nil {
		libs.Error(c, "发送失败")
		return
	}

	libs.Success(c, "发送成功", nil)
}

func buildSlideCaptcha() {
	builder := slide.NewBuilder(
	// slide.WithGenGraphNumber(2),
	// slide.WithEnableGraphVerticalRandom(true),
	)

	// background images
	imgs, err := images.GetImages()
	if err != nil {
		log.Fatalln(err)
	}

	graphs, err := tiles.GetTiles()
	if err != nil {
		log.Fatalln(err)
	}

	var newGraphs = make([]*slide.GraphImage, 0, len(graphs))
	for i := 0; i < len(graphs); i++ {
		graph := graphs[i]
		newGraphs = append(newGraphs, &slide.GraphImage{
			OverlayImage: graph.OverlayImage,
			MaskImage:    graph.MaskImage,
			ShadowImage:  graph.ShadowImage,
		})
	}

	// set resources
	builder.SetResources(
		slide.WithGraphImages(newGraphs),
		slide.WithBackgrounds(imgs),
	)

	slideCapt = builder.Make()
}

func slideCaptchaHandler(c *gin.Context) {
	captData, err := slideCapt.Generate()
	if err != nil {
		slog.Error("generate captcha error", "err", err)
		libs.Error(c, "生成验证码失败")
		return
	}
	blockData := captData.GetData()
	if blockData == nil {
		libs.Error(c, "gen captcha data failed")
		return
	}

	var masterImageBase64, tileImageBase64 string
	masterImageBase64, err = captData.GetMasterImage().ToBase64()
	if err != nil {
		libs.Error(c, "生成验证码失败")
		return
	}
	tileImageBase64, err = captData.GetTileImage().ToBase64()
	if err != nil {
		libs.Error(c, "生成验证码失败")
		return
	}
	dotsByte, _ := json.Marshal(blockData)
	//log.Printf("====dotsByte: %v", dotsByte)
	key := utils.StringToMD5(string(dotsByte))
	//log.Printf("====key: %v", key)
	common.SetCache(key, dotsByte, 3*time.Minute)

	libs.Success(c, "生成验证码成功", gin.H{
		"captcha_key":  key,
		"tile_width":   blockData.Width,
		"tile_height":  blockData.Height,
		"tile_x":       blockData.TileX,
		"tile_y":       blockData.TileY,
		"image_base64": masterImageBase64,
		"tile_base64":  tileImageBase64,
	})
}

// CheckSlideData 处理滑动验证码的验证
func checkSlideDataHandler(c *gin.Context) {
	type SlideData struct {
		Point string `json:"point" binding:"required"`
		Key   string `json:"key" binding:"required"`
	}
	var req SlideData
	if err := c.ShouldBindJSON(&req); err != nil {
		libs.Error(c, "参数错误")
		return
	}
	//log.Printf("check slide data: %v", req)
	// 获取请求参数
	point := req.Point
	key := req.Key

	// 参数校验
	if point == "" || key == "" {
		libs.Error(c, "point or key param is empty")
		return
	}

	// 从缓存中读取数据
	cacheDataByte, _ := common.GetCache(key)
	if cacheDataByte == nil {
		libs.Error(c, "illegal key")
		return
	}
	//log.Printf("cacheDataByte: %v", cacheDataByte)
	dataBytes := common.GetCacheVal(cacheDataByte)
	//log.Printf("dataBytes: %v", dataBytes)
	// 反序列化缓存中的数据
	var dct *slide.Block
	if err := json.Unmarshal(dataBytes, &dct); err != nil {
		//slog.Error("json unmarshal error", "err", err)
		libs.Error(c, "illegal key")
		return
	}

	// 解析 point 参数
	src := strings.Split(point, ",")
	chkRet := false
	if len(src) == 2 {
		sx, err := strconv.ParseFloat(src[0], 64)
		if err != nil {
			//slog.Error("parse point error", "err", err)
			libs.Error(c, "invalid point formatr")
			return
		}
		sy, err := strconv.ParseFloat(src[1], 64)
		if err != nil {
			//slog.Error("parse point error", "err", err)
			libs.Error(c, "invalid point format")
			return
		}

		// 检查用户的滑动位置是否正确
		chkRet = slide.CheckPoint(int64(sx), int64(sy), int64(dct.X), int64(dct.Y), 4)
	}

	// 如果验证成功，返回 code 0
	if chkRet {
		libs.Success(c, "验证成功", nil)
		return
	}

	libs.Error(c, "验证失败")
}
