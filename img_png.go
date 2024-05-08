package goposter

import (
	"bytes"
	"goposter/utils"
	"image/png"
	"os"
)

type ImagePNGType struct {
	ImageBaseType
}

//LoadFromFile 加载-从文件中加载
func (o *ImagePNGType) LoadFromFile(imgFile string) (err error) {
	// 根据路径打开模板文件
	imgFileHandler, err := os.Open(imgFile)
	if err != nil {
		return
	}
	defer imgFileHandler.Close()
	tmp, err := png.Decode(imgFileHandler)
	if err != nil {
		return
	}
	o.img = utils.ToRGBA(tmp)
	return
}

//LoadFromBuffer 加载-从buffer中加载
func (o *ImagePNGType) LoadFromBuffer(imgBuffer *bytes.Buffer) (err error) {
	tmp, err := png.Decode(imgBuffer)
	if err != nil {
		return
	}
	o.img = utils.ToRGBA(tmp)
	return
}

//SaveToBuffer 保存-图片为buffer
func (o *ImagePNGType) SaveToBuffer() (imgBuffer *bytes.Buffer, err error) {
	imgBuffer = bytes.NewBuffer(nil)
	err = png.Encode(imgBuffer, o.img)
	if err != nil {
		return
	}
	return
}

//NewImagePNGType 新建一个webp文件加载的图像
func NewImagePNGType() (ImagePNG IImage) {
	ImagePNG = &ImagePNGType{}
	return
}
