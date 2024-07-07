package goposter

import (
	"bytes"
	"github.com/chai2010/webp"
	"github.com/danielrqzeng/goposter/utils"
	"os"
)

type ImageWEBPType struct {
	ImageBaseType
}

//LoadFromFile 加载-从文件中加载
func (o *ImageWEBPType) LoadFromFile(imgFile string) (err error) {
	// 根据路径打开模板文件
	imgFileHandler, err := os.Open(imgFile)
	if err != nil {
		return
	}
	defer imgFileHandler.Close()

	tmp, err := webp.Decode(imgFileHandler)
	if err != nil {
		return
	}
	o.img = utils.ToRGBA(tmp)
	return
}

//LoadFromURL 加载-从url中加载
func (o *ImageWEBPType) LoadFromURL(imgURL string) (err error) {
	resReader, _, err := utils.RequestResource(imgURL)
	if err != nil {
		return
	}
	defer resReader.Close()
	var buff bytes.Buffer
	_, err = buff.ReadFrom(resReader)
	if err != nil {
		return
	}
	err = o.LoadFromBuffer(&buff)
	if err != nil {
		return
	}
	return
}

//LoadFromBuffer 加载-从buffer中加载
func (o *ImageWEBPType) LoadFromBuffer(imgBuffer *bytes.Buffer) (err error) {
	tmp, err := webp.Decode(imgBuffer)
	if err != nil {
		return
	}
	o.img = utils.ToRGBA(tmp)
	return
}

//SaveToBuffer 保存-图片为buffer
func (o *ImageWEBPType) SaveToBuffer() (imgBuffer *bytes.Buffer, err error) {
	imgBuffer = bytes.NewBuffer(nil)
	quality := float32(50.0)
	webpByte, err1 := webp.EncodeRGBA(o.img, quality)
	if err1 != nil {
		err = err1
		return
	}
	imgBuffer = bytes.NewBuffer(webpByte)
	return
}

//NewImageWEBPType 新建一个webp文件加载的图像
func NewImageWEBPType() (ImageWEBP IImage) {
	ImageWEBP = &ImageWEBPType{}
	return
}
