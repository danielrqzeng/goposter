package goposter

import (
	"bytes"
	"image"
)

type IImage interface {
	LoadFromFile(imgFile string) (err error)                                                                                         // 加载-从文件中加载
	LoadFromBuffer(imgBuffer *bytes.Buffer) (err error)                                                                              // 加载-从buffer中加载
	LoadFromURL(imgURL string) (err error)                                                                                           // 加载-从url中加载
	LoadFromImageRGBA(img *image.RGBA) (err error)                                                                                   // 加载-设置当前实例为实例对象
	LoadFromText(word string, maxWidth int, fontFile string, dpi, fontSize float64, fontColorInHex, bgColorInHex string) (err error) // 加载-根据需要的文字创建一个图像
	LoadFromNew(width, height int, fontColorInHex string) (err error)                                                                // 加载-新建画板
	SetID(id string)                                                                                                                 // 设置id
	ID() string                                                                                                                      // id
	GetImage() (img *image.RGBA)                                                                                                     // 获取-*image.RGBA格式的实例
	Resize(newWidth, newHeight uint) (err error)                                                                                     // 操作-重新设置大小
	DrawImage(sub image.Image, x, y int) (err error)                                                                                 // 操作-在当前对象的(x,y)位置，叠着画上sub
	DrawSubImage(sub IImage, x, y int, debugColorInHex string) (err error)                                                           // 操作-叠画子图
	GetSubImagePosition(subImageID string) (position [4]int, err error)                                                              // 获取子图的位置
	DrawRoundImage(radius float64) (err error)                                                                                       // 操作，给image加上矩形的圆角
	DrawCircleImage(originX, originY, radius int) (err error)                                                                        //操作，裁剪为圆
	SaveToBuffer() (imgBuffer *bytes.Buffer, err error)                                                                              // 保存-图片为buffer
	SaveToWEBPFile(imgFile string) (err error)                                                                                       // 保存-图片为webp格式
	SaveToPNGFile(imgFile string) (err error)                                                                                        // 保存-图片为png格式
	SaveToJPEGFile(imgFile string) (err error)                                                                                       // 保存-图片为jpeg或者jpg格式
}
