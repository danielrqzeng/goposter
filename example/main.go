package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/danielrqzeng/goposter"

	"github.com/spf13/viper"
	"text/template"
)

func main() {
	/*
		手机海报样式
		|-----------------------|
		| +	| name				|
		|-----------------------|
		|						|
		|						|
		|						|
		|						|
		|						|
		|						|
		|-----------------------|
		| Title			|	|	|
		| desc			|--	。--|
		|				|	|	|
		|-----------------------|
	*/

	viper.SetConfigFile("./app.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	var posterList []goposter.ImageConfigInfoType
	err = viper.UnmarshalKey("posterList", &posterList)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(fmt.Sprintf("%+v", posterList))

	jsonStr, err := json.Marshal(posterList[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(string(jsonStr))

	tmpl := template.New("poster")
	tmpl, err = tmpl.Parse(string(jsonStr))
	if err != nil {
		fmt.Println(err)
		return
	}

	pixelRatio := 2.0 // 逻辑像素和物理像素比例
	width := 320      // iphone5的尺寸
	height := 548
	avatarUrl := "./avatar.jpeg" // 用户头像
	name := "唐*李涉"               // 用户名称
	mainUrl := "./main.jpg"      // 主图片
	highlight := "闲"             //高亮字
	title := "题鹤林寺僧舍"            //标题
	desc := "因过竹院逢僧话，偷得浮生半日闲"    // 描述
	qrCodeURL := "./qrcode.png"  // 二维码

	params := map[string]string{
		"pixel_ratio": fmt.Sprintf("%f", pixelRatio),
		"width":       fmt.Sprintf("%d", int(float64(width)*pixelRatio)),
		"height":      fmt.Sprintf("%d", int(float64(height)*pixelRatio)),
		"avatar_url":  avatarUrl,
		"name":        name,
		"main_url":    mainUrl,
		"highlight":   highlight,
		"title":       title,
		"desc":        desc,
		"qr_code_url": qrCodeURL,
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, params)
	if err != nil {
		return
	}

	imageConfigInfo := &goposter.ImageConfigInfoType{}
	err = json.Unmarshal(buf.Bytes(), imageConfigInfo)
	if err != nil {
		return
	}

	err = goposter.GenByImageConfig(imageConfigInfo)
	if err != nil {
		return
	}

}