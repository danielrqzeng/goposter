package goposter

import (
	"bytes"
	"fmt"
	"github.com/danielrqzeng/goposter/utils"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
)

type ImageBaseType struct {
	img *image.RGBA
	id  string

	id2Position map[string][4]int //子图在父图中的四角定位,[id]=[minX,minY,maxX,maxY],此变量会在drawimage时候写入
}

//LoadFromFile 加载-从文件中加载
func (o *ImageBaseType) LoadFromFile(imgFile string) (err error) {
	// 部件自实现
	err = fmt.Errorf("ImageBaseType not support")
	return
}

//LoadFromBuffer 加载-从buffer中加载
func (o *ImageBaseType) LoadFromBuffer(imgBuffer *bytes.Buffer) (err error) {
	// 部件自实现
	err = fmt.Errorf("ImageBaseType not support")

	return
}

//LoadFromURL 加载-从url中加载
func (o *ImageBaseType) LoadFromURL(imgURL string) (err error) {
	// 部件自实现
	err = fmt.Errorf("ImageBaseType not support")
	return
}

//LoadFromImageRGBA 加载-设置当前实例为实例对象
func (o *ImageBaseType) LoadFromImageRGBA(srcImg *image.RGBA) (err error) {
	o.img = srcImg
	return
}

//LoadFromText 加载-根据需要的文字创建一个图像
func (o *ImageBaseType) LoadFromText(word string, maxWidth int, fontFile string, dpi, fontSize float64, fontColorInHex, bgColorInHex string) (err error) {
	//此处的一些概念
	// 字的宽度高度：有字体库决定
	// 字和字的间隔：是为了字体美观连续可读，可能为负或者正，意味着每个相邻的字间隔可能不同，比如字[a,b]可能间隔很小，字[一,二]可能间隔比较大
	//  行与行：maxWidth代表最大宽度，因为如果字太多一行摆不下，此时需要换行，0代表不限制
	//  fontSize单位为pt，一般用8pt小字，13pt大字
	// 坐标的分布，请请参考：像素坐标系
	//		.---------------------->x
	//		|
	//		|
	//		|
	//		|
	//		|
	//		y
	// 图片的左上角为坐标原点

	textFont, err := utils.LoadFont(fontFile) //文字库文件
	if err != nil {
		return
	}
	// scale converts truetype.FUnit to float64
	fontColor, err := utils.Hex2RGB(fontColorInHex)
	if err != nil {
		return
	}
	//计算排布的宽度
	lineTextList := make([]string, 0) // 字和位置信息,word="123456"==>lineTextList:["123","456"]，,word="1234\n56"==>lineTextList:["123","4","56"]
	canvasWidth := 0                  //底图的宽度度
	canvasHeight := 0                 //底图的高度
	lineHeight := 0                   //每行文字的高度

	lineTextList = append(lineTextList, "")

	for _, r := range word {
		currLineWord := lineTextList[len(lineTextList)-1] + string(r)
		w, h, err1 := utils.MeasureText(currLineWord, textFont, dpi, fontSize)
		//fmt.Println(fmt.Sprintf("draw '%s' at (%d,%d)", currLineWord, w, h))
		if err1 != nil {
			err = err1
			return
		}
		lineHeight = h
		//if string(r) == " " {
		//	lineTextList = append(lineTextList, "")
		//	continue
		//}
		// 换行-有换行符
		if string(r) == "\n" {
			lineTextList = append(lineTextList, "")
			continue
		}
		// 换行-达到最大宽度
		if maxWidth != 0 && w > maxWidth {
			lineTextList = append(lineTextList, "")
		}

		if canvasWidth < w {
			canvasWidth = w
		}

		lineTextList[len(lineTextList)-1] += string(r)

	}

	canvasHeight = lineHeight * len(lineTextList)

	//计算底图大小
	err = o.LoadFromNew(canvasWidth, canvasHeight, bgColorInHex)
	if err != nil {
		return
	}

	margin := 0
	x, y := 0, 0
	for i, lineText := range lineTextList {
		x = 0 + margin
		if i == 0 {
			y += margin
		} else {
			y += lineHeight
		}
		//fmt.Println(fmt.Sprintf("draw '%s' dpi=%f,fontSize=%f at (%d,%d) on (%d,%d)", lineText, dpi, fontSize, x, y, canvasWidth, canvasHeight))
		err = o.DrawOneLineText(lineText, textFont, dpi, fontSize, fontColor, x, y)
		if err != nil {
			return
		}
	}
	return
}

//LoadFromNew 加载-新建画板,fontColorInHex如果为空代表不需要底色
func (o *ImageBaseType) LoadFromNew(width, height int, fontColorInHex string) (err error) {

	rect := image.Rect(0, 0, width, height)
	o.img = image.NewRGBA(rect.Bounds())
	//fontColorInHex如果无底色，代表是透明底
	if fontColorInHex != "" {
		fontColor, err1 := utils.Hex2RGB(fontColorInHex)
		if err1 != nil {
			err = err1
			return
		}
		draw.Draw(o.img, o.img.Bounds(), &image.Uniform{C: fontColor}, image.Pt(0, 0), draw.Src)
	}
	return
}

//SetID 设置id
func (o *ImageBaseType) SetID(id string) {
	o.id = id
	if o.id2Position == nil {
		o.id2Position = make(map[string][4]int)
	}
}

//ID 获取id
func (o *ImageBaseType) ID() string {
	return o.id
}

//GetImage 获取-*image.RGBA格式的实例
func (o *ImageBaseType) GetImage() (srdImg *image.RGBA) {
	srdImg = o.img
	return
}

//Resize 操作-重新设置大小
func (o *ImageBaseType) Resize(newWidth, newHeight uint) (err error) {
	resizedImg := resize.Resize(newWidth, newHeight, o.img, resize.Lanczos3)
	o.img = utils.ToRGBA(resizedImg)
	return
}

//DrawOneLineText 操作-画文字（一行文字）
func (o *ImageBaseType) DrawOneLineText(oneLineText string, tf *truetype.Font, dpi, fontSize float64, fontColor color.Color, x int, y int) (err error) {
	face := truetype.NewFace(tf, &truetype.Options{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: 0,
		//GlyphCacheEntries: 0,
		//SubPixelsX:        0,
		//SubPixelsY:        0,
	})

	p := freetype.Pt(x, y+face.Metrics().Ascent.Ceil()) // 此处计算落点basepoint(加上Ascent，否则看不见文字）

	// 使用 DrawString 绘制文本
	d := &font.Drawer{
		Dst:  o.img,
		Src:  image.NewUniform(fontColor),
		Face: face,
		Dot:  p, // 设置起始点 p
	}

	d.DrawString(oneLineText)
	//utils.DrawHorizLine(o.img, image.NewUniform(color.RGBA{
	//	R: 0,
	//	G: 255,
	//	B: 0,
	//	A: 0,
	//}), o.img.Bounds().Min.X, o.img.Bounds().Max.X, face.Metrics().Ascent.Ceil())

	return
}

//DrawImage 操作-在当前对象的(x,y)位置，叠着画上sub
func (o *ImageBaseType) DrawImage(sub image.Image, x, y int) (err error) {
	r := sub.Bounds().Add(image.Pt(x, y)) //计算处理子图像，应该放置在底图的矩形坐标系的点
	draw.Draw(o.img, r, sub, image.Point{}, draw.Over)

	return
}

//DrawImage 操作-在当前对象的(x,y)位置，叠着画上sub
func (o *ImageBaseType) DrawSubImage(sub IImage, x, y int, debugColorInHex string) (err error) {
	//四角定位
	r := sub.GetImage().Bounds().Add(image.Pt(x, y)) //计算处理子图像，应该放置在底图的矩形坐标系的点
	o.id2Position[sub.ID()] = [4]int{r.Min.X, r.Min.Y, r.Max.X, r.Max.Y}

	err = o.DrawImage(sub.GetImage(), x, y)

	// 如果开启了debug模式，则做debug背景色块
	if debugColorInHex != "" {
		debugMask := image.NewRGBA(sub.GetImage().Bounds())
		debugColor, err1 := utils.Hex2RGB(debugColorInHex)
		if err1 != nil {
			err = err1
			return
		}
		draw.Draw(debugMask, debugMask.Bounds(), &image.Uniform{C: debugColor}, image.Point{}, draw.Over)
		err = o.DrawImage(debugMask, x, y)
	}

	return
}

//GetSubImagePosition 操作-获取子图的位置信息
func (o *ImageBaseType) GetSubImagePosition(subImageID string) (position [4]int, err error) {
	if _, ok := o.id2Position[subImageID]; !ok {
		err = fmt.Errorf("not found subImageID=%s", subImageID)
		return
	}
	return o.id2Position[subImageID], nil
}

//DrawRoundImage 操作-将sub裁剪为圆角
func (o *ImageBaseType) DrawRoundImage(radius float64) (err error) {
	maskCanvasImage := image.NewRGBA(o.img.Bounds())

	cm := &RoundMask{o.img, int(radius)}
	//draw.Draw(maskCanvasImage, maskCanvasImage.Bounds(), cm, image.Point{}, draw.Over)
	draw.Draw(maskCanvasImage, maskCanvasImage.Bounds(), cm, image.ZP, draw.Src) //如使用上面注释的，会出现锯齿
	o.img = maskCanvasImage
	return
}

//DrawRoundImage 操作-裁剪为圆
func (o *ImageBaseType) DrawCircleImage(originX, originY, radius int) (err error) {
	maskCanvasImage := image.NewRGBA(o.img.Bounds())

	cm := &CircleMask{o.img, originX, originY, radius}
	//draw.Draw(maskCanvasImage, maskCanvasImage.Bounds(), cm, image.Point{}, draw.Over)
	draw.Draw(maskCanvasImage, maskCanvasImage.Bounds(), cm, image.ZP, draw.Src)
	o.img = maskCanvasImage
	return
}

//SaveToBuffer 保存-图片为buffer
func (o *ImageBaseType) SaveToBuffer() (imgBuffer *bytes.Buffer, err error) {
	// 部件自实现

	//imgBuffer = bytes.NewBuffer(nil)
	//
	//switch imageType {
	//case ImageTypeJPG:
	//	err = jpeg.Encode(imgBuffer, o.img, &jpeg.Options{})
	//	if err != nil {
	//		return
	//	}
	//	return
	//case ImageTypePNG:
	//	err = png.Encode(imgBuffer, o.img)
	//	if err != nil {
	//		return
	//	}
	//	return
	//case ImageTypeWEBP:
	//	quality := float32(50.0)
	//	webpByte, err1 := webp.EncodeRGBA(o.img, quality)
	//	if err1 != nil {
	//		err = err1
	//		return
	//	}
	//	imgBuffer = bytes.NewBuffer(webpByte)
	//	return
	//default:
	//	err = fmt.Errorf("imageType=" + utils.Num2Str(imageType) + " not support yet")
	//	return
	//}
	return
}

//SaveToPNGFile 保存-图片为png格式
func (o *ImageBaseType) SaveToPNGFile(imgFile string) (err error) {
	dstFile, err := os.Create(imgFile)
	if err != nil {
		return
	}
	defer dstFile.Close()

	err = png.Encode(dstFile, o.img)
	if err != nil {
		return
	}
	return
}

//SaveToJPEGFile 保存-图片为jpeg或者jpg格式
func (o *ImageBaseType) SaveToJPEGFile(imgFile string) (err error) {
	dstFile, err := os.Create(imgFile)
	if err != nil {
		return
	}
	defer dstFile.Close()

	err = jpeg.Encode(dstFile, o.img, &jpeg.Options{})
	if err != nil {
		return
	}
	return
}
