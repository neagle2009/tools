package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "io"
	"io/ioutil"
	"log"
	_ "net/http"
	"net/mail"
	"net/smtp"
	"net/url"
	"os"
	_ "path/filepath"
	"strconv"
	_ "strings"
	"time"
)

const (
	INIT_LIST_NUM    = 4
	GRUB_URL         = "https://sjz.zu.ke.com/zufang/rco11rs%E9%87%91%E7%9F%B3%E5%B0%8F%E5%8C%BA/?unique_id=d38d1a2c-fd9b-491b-b809-d0f7e820fcfazufangrco11rs%E9%87%91%E7%9F%B3%E5%B0%8F%E5%8C%BA1556893750938"
	SHELL_HOUSE_HOST = "https://sjz.zu.ke.com"

	DATA_FILE = "/home/neagle/house.json"
	SMTP_USER = "itisasecretbaby@163.com"
	SMTP_PASS = "130910*******" //注意修改密码  @TODO
	SMTP_HOST = "smtp.163.com"
	SMTP_PORT = "25"

	MAIL_MIME_HTML = "text/html"
	MAIL_MIME_TEXT = "text/plain"

	BASE64_ENCODE_STRING = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

	NOTICE_EMAIL = "gu217@126.com"

	ERROR_MIME = MAIL_MIME_TEXT

	SUCCESS_MIME  = MAIL_MIME_HTML
	SUCCESS_TITLE = "有新房了,抓紧!!!! %d条 [BY ID CHECK]"
	SUCCESS_BODY  = `<!doctype html>
<html>
<head>
	<meta charset="utf-8" />
	<title>新房通知</title>
</head>
<body style="margin: 0;padding: 0;">
<h3><a href="https://sjz.zu.ke.com/zufang/rs%E9%87%91%E7%9F%B3%E5%B0%8F%E5%8C%BA/?unique_id=d38d1a2c-fd9b-491b-b809-d0f7e820fcfazufangrs%E9%87%91%E7%9F%B3%E5%B0%8F%E5%8C%BA1556672325357">有房了, 抓紧啦! 点击查看</a></h3>
</body>
</html>
`
)

type House struct {
	Name string `json:"name"`
	Href string `json:"href"`
}

func main() {

	doc, err := goquery.NewDocument(GRUB_URL)
	logOut("goQuery init")
	if err != nil {
		SendToMail(NOTICE_EMAIL, "query init error", "", ERROR_MIME)
		log.Fatal(err)
	}

	textNum := doc.Find("span.content__title--hl").Text()
	grubNum, err := strconv.Atoi(textNum)
	if nil != err {
		SendToMail(NOTICE_EMAIL, "grub get list item num erro", "", MAIL_MIME_TEXT)
		log.Fatal(err)
	}

	if NeedNotice(doc) {
		SendToMail(NOTICE_EMAIL, fmt.Sprintf(SUCCESS_TITLE, grubNum), SUCCESS_BODY, SUCCESS_MIME)
		return
	}

	return
}

func NeedNotice(doc *goquery.Document) bool {
	s := make(map[string]*House)
	doc.Find("div.content__article").Find("a.content__list--item--aside").Each(func(i int, sel *goquery.Selection) {
		href := sel.AttrOr("href", "null")
		id := ""
		u, err := url.Parse(href)
		if nil == err {
			id = u.Query().Get("h")
		}
		title := sel.Find("img.lazyload").AttrOr("alt", "empty-alt")
		s[id] = &House{title, SHELL_HOUSE_HOST + href}
	})
	return !checkHouseExist(s)
}

func storeData(hl map[string]*House) {
	if j, e := json.Marshal(hl); e == nil {
		ioutil.WriteFile(DATA_FILE, j, 0644)
	}
}

func checkHouseExist(houseList map[string]*House) bool {
	if _, e := os.Stat(DATA_FILE); os.IsNotExist(e) {
		logOut(DATA_FILE + " NOT EXIST")
		storeData(houseList)
		return false
	}

	data, err := ioutil.ReadFile(DATA_FILE)
	if err != nil {
		logOut(DATA_FILE + " READ ERROR")
		storeData(houseList)
		return false
	}

	var oldList map[string]*House
	if err = json.Unmarshal(data, &oldList); err != nil {
		logOut(DATA_FILE + " json.Unmarshal error")
		storeData(houseList)
		log.Fatal(err)
	}

	//检查新房ID是否在旧列表中,不存在则重新存储新列表
	for k, _ := range houseList {
		if _, ok := oldList[k]; !ok {
			logOut("new list itme add for:" + k)
			storeData(houseList)
			return false
		}
	}

	return true
}

func SendToMail(toEmail, subject, body, mailType string) {
	b64 := base64.NewEncoding(BASE64_ENCODE_STRING)

	from := mail.Address{"[租房中心]", SMTP_USER}
	to := mail.Address{"", toEmail}

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", b64.EncodeToString([]byte(subject)))
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = mailType + "; charset=UTF-8"
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + b64.EncodeToString([]byte(body))

	auth := smtp.PlainAuth("", SMTP_USER, SMTP_PASS, SMTP_HOST)
	logOut("mail auth done")
	err := smtp.SendMail(SMTP_HOST+":"+SMTP_PORT, auth, SMTP_USER, []string{to.Address}, []byte(message))
	logOut("mail send done")
	if err != nil {
		log.Fatal(err)
	}
}

func logOut(s string) {
	fmt.Printf("%+v\t%s\n", time.Now().Format("2006-01-02 15:04:05"), s)
}
