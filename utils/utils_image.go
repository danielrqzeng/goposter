package utils

import (
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"
)

//LoadFont 加载字体
func LoadFont(fontFile string) (font *truetype.Font, err error) {
	var fontBytes []byte
	fontBytes, err = ioutil.ReadFile(fontFile) // 读取字体文件
	if err != nil {
		err = fmt.Errorf("can not load from file=%s,err=%s", fontFile, err.Error())
		return
	}
	font, err = freetype.ParseFont(fontBytes) // 解析字体文件
	if err != nil {
		err = fmt.Errorf("font file parse err,err=%s", err.Error())
		return
	}
	return
}

// ToRGBA 将任意image.Image转换为image.RGBA类型
func ToRGBA(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, image.Point{}, draw.Src)
	return rgba
}

//Hex2RGB 将十六进制的颜色码转换为颜色对象,比如#123456=>(R:12,G:34,B:56), #12345678=>((R:12,G:34,B:56,A:78)
func Hex2RGB(hexStr string) (c color.Color, err error) {
	if len(hexStr) <= 0 {
		err = fmt.Errorf("hexStr is null")
		return
	}
	if hexStr[0] != '#' {
		err = fmt.Errorf("hexStr format must start with #")
		return
	}
	rStr, gStr, bStr, aStr := "", "", "", "FF"
	if len(hexStr) == 6+1 {
		rStr = hexStr[1:3]
		gStr = hexStr[3:5]
		bStr = hexStr[5:7]
	} else if len(hexStr) == 8+1 {
		rStr = hexStr[1:3]
		gStr = hexStr[3:5]
		bStr = hexStr[5:7]
		aStr = hexStr[7:9]
	} else {
		err = fmt.Errorf("hexStr format err")
		return
	}

	r, g, b, a := uint8(0), uint8(0), uint8(0), uint8(255)

	// red部分解析
	tmp, err := strconv.ParseInt(rStr, 16, 32)
	if err != nil {
		err = fmt.Errorf("hexStr format err for red part")
	}
	r = uint8(tmp)

	// green部分解析
	tmp, err = strconv.ParseInt(gStr, 16, 32)
	if err != nil {
		err = fmt.Errorf("hexStr format err for green part")
	}
	g = uint8(tmp)

	// blue部分解析
	tmp, err = strconv.ParseInt(bStr, 16, 32)
	if err != nil {
		err = fmt.Errorf("hexStr format err for blue part")
	}
	b = uint8(tmp)

	// alpha部分解析
	tmp, err = strconv.ParseInt(aStr, 16, 32)
	if err != nil {
		err = fmt.Errorf("hexStr format err for alpha part")
	}
	a = uint8(tmp)

	c = color.RGBA{R: r, G: g, B: b, A: a}
	return
}

//DrawHorizLine 画一条水平线
func DrawHorizLine(img *image.RGBA, c color.Color, fromX, toX, y int) {
	for x := fromX; x <= toX; x++ {
		img.Set(x, y, c)
	}
}

//DrawVertLine 画一条竖直线
func DrawVertLine(img *image.RGBA, c color.Color, x, fromY, toY int) {
	for y := fromY; y < toY; y++ {
		img.Set(x, y, c)
	}
}

//MeasureText 测量一行文字的宽高
func MeasureText(oneLineText string, tf *truetype.Font, dpi, fontSize float64) (width int, height int, err error) {
	face := truetype.NewFace(tf, &truetype.Options{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: 0,
	})
	bounds, _ := font.BoundString(face, oneLineText)
	width, height = bounds.Max.X.Ceil()-bounds.Min.X.Ceil(), face.Metrics().Descent.Ceil()+face.Metrics().Ascent.Ceil()

	return
}

//GetFontSizeByHeight 通过文字的高度得到文字的大小（fontsize)
func GetFontSizeByHeight(tf *truetype.Font, dpi float64, fontHeight int) (fontSize float64, err error) {
	//fromSize, toSize := 0.0, 1638.0 //字体在word中最大为1638
	fromSize, toSize := 0.0, 72.0 //字体在word中最大为1638
	const diff = 1.0
	const maxLoopTime = 100 //最大计算次数，如果超过这个值还算不出来，使用当前值
	idx := 0
	h := 0
	for true {
		idx++
		fontSize = (fromSize + toSize) / 2
		_, h, err = MeasureText("A", tf, dpi, fontSize)
		if err != nil {
			return
		}

		//超过了最大计算次数
		if idx > maxLoopTime {
			return
		}

		//在误差范围内，返回结果
		if math.Abs(float64(h-fontHeight)) < diff {
			break
		}

		//字号太大
		if h > fontHeight {
			toSize = fontSize
		}
		//字体太小
		if h < fontHeight {
			fromSize = fontSize
		}
	}
	return
}

//ParseNumPercentNumNone 解析字符串是整数型|百分比数|none,如果是非none，计算返回真正的尺寸
func ParseNumPercentNumNone(numStr string, baseSize int) (isNone bool, realSize int, err error) {
	isNone = false
	if strings.ToLower(strings.TrimSpace(numStr)) == "none" {
		isNone = true
		return
	}

	if len(numStr) == 0 {
		err = fmt.Errorf("unvalid format for sizeStr null")
		return
	}
	//如果是纯数字，则转换为int
	pattern := `^\d+(\.\d+)?$`
	matched, err := regexp.MatchString(pattern, numStr)
	if err != nil {
		return
	}
	// 纯数字类型的
	if matched {
		tmp, err1 := strconv.ParseFloat(numStr, 64)
		if err1 != nil {
			err = err1
			return
		}
		realSize = int(tmp)
		fmt.Println("pure number,numStr=", numStr, ",baseSize=", baseSize, ",size=", realSize)
		return
	}

	//如果是百分数，则通过父size转换得到最终大小，格式如45%
	pattern = `^\d+(\.\d+)?%$`
	matched, err = regexp.MatchString(pattern, numStr)
	if err != nil {
		return
	}
	// 纯数字类型的
	if matched {
		tmp, err1 := strconv.ParseFloat(numStr[:len(numStr)-1], 64)
		if err1 != nil {
			err = err1
			return
		}
		realSize = int(tmp / 100 * float64(baseSize))
		fmt.Println("percent number,numStr=", numStr, ",baseSize=", baseSize, ",numStr=", numStr, ",tmp=", tmp, ",size=", realSize)

		return
	}

	//格式不对
	err = fmt.Errorf("unvalid format for sizeStr=" + numStr)
	return
}

//ParseAbsoluteLocation 根据绝对定位配置，计算自身位置信息，ParseAbsoluteLocation("1% none none 2%",{100,200},{20,20})==>(20,20,40,40,nil)(如果没找到，则返回 math.MaxInt32）
func ParseAbsoluteLocation(absolutePositionStr string, parentSize [2]int, selfSize [2]int) (minX, minY, maxX, maxY int, err error) {
	tmp := strings.Split(absolutePositionStr, " ")
	if len(tmp) != 4 {
		err = fmt.Errorf("absolutePositionStr format err,absolutePositionStr=" + absolutePositionStr)
		return
	}

	const WidthIdx = 0
	const HeightIdx = 1

	minX, minY, maxX, maxY = math.MaxInt32, math.MaxInt32, math.MaxInt32, math.MaxInt32
	for idx, str := range tmp {
		//对于top进行计算
		if idx == 0 {
			parentHeight := parentSize[HeightIdx]
			isNone, realSize, err1 := ParseNumPercentNumNone(str, parentHeight)
			if err1 != nil {
				err = err1
				return
			}
			if isNone {
				continue
			}
			minY = realSize
			maxY = minY + selfSize[HeightIdx]
		}
		//对于right进行计算
		if idx == 1 {
			parentWeight := parentSize[WidthIdx]
			isNone, realSize, err1 := ParseNumPercentNumNone(str, parentWeight)
			if err1 != nil {
				err = err1
			}
			if isNone {
				continue
			}
			maxX = parentWeight - realSize
			minX = maxX - selfSize[WidthIdx]
		}

		//对于bottom进行计算
		if idx == 2 {
			parentHeight := parentSize[HeightIdx]
			isNone, realSize, err1 := ParseNumPercentNumNone(str, parentHeight)
			if err1 != nil {
				err = err1
			}
			if isNone {
				continue
			}
			maxY = parentHeight - realSize
			minY = maxY - selfSize[HeightIdx]
		}

		//对于left进行计算
		if idx == 3 {
			parentWeight := parentSize[WidthIdx]
			isNone, realSize, err1 := ParseNumPercentNumNone(str, parentWeight)
			if err1 != nil {
				err = err1
			}
			if isNone {
				continue
			}
			minX = realSize
			maxX = minX + selfSize[WidthIdx]
		}
	}
	return
}

//ParseRelativeLocation 根据相对定位配置，计算自身位置信息(如果没找到，则返回 math.MaxInt32）
func ParseRelativeLocation(relativePositionStr string, parentSize [2]int, selfSize [2]int, relativePosition [4]int) (minX, minY, maxX, maxY int, err error) {
	tmp := strings.Split(relativePositionStr, " ")
	if len(tmp) != 4 {
		err = fmt.Errorf("relativePositionStr format err,relativePositionStr=" + relativePositionStr)
		return
	}

	const WidthIdx = 0
	const HeightIdx = 1
	const MinXIdx = 0
	const MinYIdx = 1
	const MaxXIdx = 2
	const MaxYIdx = 3

	minX, minY, maxX, maxY = math.MaxInt32, math.MaxInt32, math.MaxInt32, math.MaxInt32
	/*
	                       +------------+
	                       |   relative |
	                       +------------+
	                       |    {top}   |
	   +----------+--------+------------+--------+------------+
	   | relative | {left} |    self    | {right}|  relative  |
	   +----------+--------+------------+--------+------------+
	                       |   {bottom} |
	                       +------------+
	                       |   relative |
	                       +------------+
	*/
	for idx, str := range tmp {
		//对于top进行计算(即自身会位置将会在相对定位体的下方）
		if idx == 0 {
			/*
				+------------+
				|	relative |
				+------------+
				|	{top}	 |
				+------------+
				|	self	 |
				+------------+
				自身的上边距离相对元素，有top的距离
			*/
			parentHeight := parentSize[HeightIdx]
			isNone, realSize, err1 := ParseNumPercentNumNone(str, parentHeight)
			if err1 != nil {
				err = err1
				return
			}
			if isNone {
				continue
			}
			minY = relativePosition[MaxYIdx] + realSize
			maxY = minY + selfSize[HeightIdx]
		}
		//对于right进行计算(即自身会位置将会在相对定位体的左方）
		if idx == 1 {
			/*
				+------------+--------+------------+
				|	self     | {right}|	 relative  |
				+------------+--------+------------+
				自身的右边距离相对元素，有right的距离
			*/
			parentWeight := parentSize[WidthIdx]
			isNone, realSize, err1 := ParseNumPercentNumNone(str, parentWeight)
			if err1 != nil {
				err = err1
			}
			if isNone {
				continue
			}

			maxX = relativePosition[MinXIdx] - realSize
			minX = maxX - selfSize[WidthIdx]
		}

		//对于bottom进行计算(即自身会位置将会在相对定位体的上方）
		if idx == 2 {
			/*
				+------------+
				|	 self    |
				+------------+
				|	{bottom} |
				+------------+
				|	relative |
				+------------+
				自身的下边距离相对元素，有bottom的距离
			*/
			parentHeight := parentSize[HeightIdx]
			isNone, realSize, err1 := ParseNumPercentNumNone(str, parentHeight)
			if err1 != nil {
				err = err1
			}
			if isNone {
				continue
			}
			maxY = relativePosition[MinYIdx] - realSize
			minY = maxY - selfSize[HeightIdx]
		}

		//对于left进行计算(即自身会位置将会在相对定位体的右方）
		if idx == 3 {
			/*
				+------------+--------+------------+
				|	relative | {left} |	self	   |
				+------------+--------+------------+
				自身的左边距离相对元素，有left的距离
			*/
			parentWeight := parentSize[WidthIdx]
			isNone, realSize, err1 := ParseNumPercentNumNone(str, parentWeight)
			if err1 != nil {
				err = err1
			}
			if isNone {
				continue
			}
			minX = relativePosition[MaxXIdx] + realSize
			maxX = minX + selfSize[WidthIdx]
		}
	}
	return
}

//CalcRealSize 判断和计算sizeStr是纯数字类型的指标，还是百分百类型的，都转换为数字类型的返回， CalcRealSize("30", 200)=>30, CalcRealSize("30%", 200)=>60,CalcRealSize("30.0%", 200)=>60,
func CalcRealSize(sizeStr string, parentSize int) (size int, err error) {
	if len(sizeStr) == 0 {
		//err = fmt.Errorf("unvalid format for sizeStr=" + sizeStr)
		size = parentSize
		return
	}
	//如果是纯数字，则转换为int
	pattern := `^\d+(\.\d+)?$`
	matched, err := regexp.MatchString(pattern, sizeStr)
	if err != nil {
		return
	}
	// 纯数字类型的
	if matched {
		tmp, err1 := strconv.ParseFloat(sizeStr, 64)
		if err1 != nil {
			err = err1
			return
		}
		size = int(tmp)
		return
	}

	//如果是百分数，则通过父size转换得到最终大小
	pattern = `^\d+(\.\d+)?%$`
	matched, err = regexp.MatchString(pattern, sizeStr)
	if err != nil {
		return
	}
	// 纯数字类型的
	if matched {
		tmp, err1 := strconv.ParseFloat(sizeStr[:len(sizeStr)-1], 64)
		if err1 != nil {
			err = err1
			return
		}
		size = int(tmp / 100 * float64(parentSize))

		return
	}

	//格式不对
	err = fmt.Errorf("unvalid format for sizeStr=" + sizeStr)
	return
}
