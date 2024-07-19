package goposter

import (
	"bytes"
	"github.com/danielrqzeng/goposter/utils"
	"image/jpeg"
	"os"
)

type ImageJPGType struct {
	ImageBaseType
}

//LoadFromFile 加载-从文件中加载
func (o *ImageJPGType) LoadFromFile(imgFile string) (err error) {
	// 根据路径打开模板文件
	imgFileHandler, err := os.Open(imgFile)
	if err != nil {
		return
	}
	defer imgFileHandler.Close()
	tmp, err := jpeg.Decode(imgFileHandler)
	if err != nil {
		return
	}
	o.img = utils.ToRGBA(tmp)
	return
}

//LoadFromURL 加载-从url中加载
func (o *ImageJPGType) LoadFromURL(imgURL string) (err error) {
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
func (o *ImageJPGType) LoadFromBuffer(imgBuffer *bytes.Buffer) (err error) {
	tmp, err := jpeg.Decode(imgBuffer)
	if err != nil {
		return
	}
	o.img = utils.ToRGBA(tmp)
	return
}

//SaveToBuffer 保存-图片为buffer
func (o *ImageJPGType) SaveToBuffer() (imgBuffer *bytes.Buffer, err error) {
	imgBuffer = bytes.NewBuffer(nil)
	err = jpeg.Encode(imgBuffer, o.img, &jpeg.Options{})
	if err != nil {
		return
	}
	return
}

//NewImageJPGType 新建一个jpeg文件加载的图像
func NewImageJPGType() (ImageJPG IImage) {
	ImageJPG = &ImageJPGType{}
	return
}
