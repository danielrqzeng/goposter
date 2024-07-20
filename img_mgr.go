package goposter

import (
	"bytes"
	"fmt"
	"github.com/danielrqzeng/goposter/utils"
	"strconv"
	"strings"
	"sync"
)

//实例
var (
	ImageMgrInstance *ImageMgrType
	ImageMgrOnce     sync.Once
)

//ImageMgr 获取单例实例
func ImageMgr() *ImageMgrType {
	ImageMgrOnce.Do(func() {
		ImageMgrInstance = &ImageMgrType{}
		ImageMgrInstance.Init()
	})
	return ImageMgrInstance
}

//ImageMgrType 实例定义
type ImageMgrType struct {
	sync.Mutex
	debugColor [][2]string // 开启调试开关，调试背景色, [[rbga1,subImageID1],[rbga2,subImageID2],[rbga3,""]]
}

//Init 初始化
func (mgr *ImageMgrType) Init() {
	mgr.debugColor = make([][2]string, 0)
	// 生成调试背景色
	cl := []string{
		//"00",
		"11",
		"22",
		"33",
		"44",
		"55",
		"66",
		"77",
		"88",
		"99",
		"AA",
		"BB",
		"CC",
		"DD",
		"EE",
		//"FF",
	}
	rbgList := make([]string, 0)

	a := "02" //给1%的透明
	for _, c := range cl {
		rbgList = append(rbgList, "#"+c+"0000"+a)
		rbgList = append(rbgList, "#00"+c+"00"+a)
		rbgList = append(rbgList, "#0000"+c+a)
	}

	for i := 0; i < len(rbgList); i++ {
		mgr.debugColor = append(mgr.debugColor, [2]string{rbgList[i], ""})
	}

	go mgr.grLoop()
}

//grLoop 默认goroutine
func (mgr *ImageMgrType) grLoop() {
}

//grLoop 默认goroutine
func (mgr *ImageMgrType) GetDebugColor(subImageID string) (debugColor string) {
	//mgr.debugColor = [[rbga1,subImageID1],[rbga2,subImageID2],[rbga3,""]]
	for idx, info := range mgr.debugColor {
		//说明此颜色还没有被分配，分配之
		if info[1] == "" {
			mgr.debugColor[idx][1] = subImageID
			return info[0]
		}
		//此颜色已经被分配给这个id了
		if info[1] == subImageID {
			return info[0]
		}
	}
	//找不到更多的颜色分配了，默认返回透明黑色
	return "#000000CC"
}

//GenPhoneImage 生成海报
func (mgr *ImageMgrType) GenByImageConfig(imageConfigInfo *ImageConfigInfoType) (imgBuffer *bytes.Buffer, err error) {

	imageMap := make(map[string]IImage)

	//画布
	var canvas IImage
	isNone, canvasWidth, canvasHeight := true, 0, 0
	isNone, canvasWidth, err = utils.ParseNumPercentNumNone(imageConfigInfo.Width, 0)
	if err != nil {
		return
	}
	if isNone {
		err = fmt.Errorf("canvasWidth not allow none")
		return
	}
	isNone, canvasHeight, err = utils.ParseNumPercentNumNone(imageConfigInfo.Height, 0)
	if err != nil {
		return
	}
	if isNone {
		err = fmt.Errorf("canvasHeight not allow none")
		return
	}
	//检查图片类型
	if imageConfigInfo.OutputBufferType != ImageTypePNG &&
		imageConfigInfo.OutputBufferType != ImageTypeJPEG {
		err = fmt.Errorf("OutputBufferType err,not support image type=" + imageConfigInfo.OutputBufferType)
		return
	}
	switch imageConfigInfo.OutputBufferType {
	case ImageTypePNG:
		canvas = &ImagePNGType{}
	case ImageTypeJPEG:
		canvas = &ImageJPGType{}
	default:
		err = fmt.Errorf("OutputBufferType err,not support image type=" + imageConfigInfo.OutputBufferType)
		return
	}

	canvas = &ImagePNGType{}
	err = canvas.LoadFromNew(canvasWidth, canvasHeight, imageConfigInfo.CanvasBackgroundColor)
	if err != nil {
		return
	}
	canvas.SetID("canvas") // 画板id

	for _, s := range imageConfigInfo.SubImageInfoList {
		if strings.ToLower(s.Enable) == "false" {
			continue
		}
		if strings.ToLower(s.Enable) != "true" {
			err = fmt.Errorf("for subImageID=" + s.ID + " enable only can set true|false")
			return
		}
		for idx, a := range s.ActionList {
			//fmt.Println("++++++++++++draw image=", s.Name, " with action=", a.ActionType, "++++++++++++")
			var subImg IImage
			switch a.ActionType {
			case ActionTypeImage:
				switch a.ImageType {
				case ImageTypePNG:
					subImg = &ImagePNGType{}
				case ImageTypeJPEG:
					subImg = &ImageJPGType{}
				default:
					err = fmt.Errorf("for subImageID=" + s.ID + " not support image type=" + a.ImageType)
					return
				}
				// 从url加载
				if a.ImageURLFile != "" {
					err = subImg.LoadFromURL(a.ImageURLFile)
					if err != nil {
						return
					}
				}
				//从本地文件中加载
				if a.ImageLocalFile != "" {
					err = subImg.LoadFromFile(a.ImageLocalFile)
					if err != nil {
						return
					}
				}
				if a.ImageCustomFile != "" {
					err = fmt.Errorf("not support for ImageCustomFile")
					return
				}
			//类型为font的作图
			case ActionTypeFont:
				subImg = &ImageBaseType{}
				isNone, maxWidth := true, 0
				isNone, maxWidth, err = utils.ParseNumPercentNumNone(a.MaxWidth, canvas.GetImage().Bounds().Dx())
				if err != nil {
					return
				}
				if isNone {
					err = fmt.Errorf("MaxWidth not allow none")
					return
				}
				//dpi和fontsize的计算
				/*
					https://cloud.tencent.com/developer/article/2113226
					dip設置與分辯率無關,但写屏幕密度有关.在默认情况下,
					LDPI密度为120,系数为0.75,
					MDPI的密度为160,系数为1.0；
					HDPI的密码为240,系数为1.5；
					XHDPI的密度为320,系数为2.0；
					所谓密度即单位平方英寸中含像素的数量
					> 系数即为设备像素比
				*/
				pixelRatio, err1 := strconv.ParseFloat(imageConfigInfo.PixelRatio, 64)
				if err1 != nil {
					err = fmt.Errorf("PixelRatio=" + imageConfigInfo.PixelRatio + " unvalid")
					return
				}
				dpi := float64(160)
				if pixelRatio <= 0 {
					err = fmt.Errorf("PixelRatio=" + imageConfigInfo.PixelRatio + " unvalid")
					return
				} else if pixelRatio > 2.0 {
					dpi = 320
				} else if pixelRatio <= 2.0 {
					dpi = 320
				} else if pixelRatio <= 1.5 {
					dpi = 240
				} else if pixelRatio <= 1.0 {
					dpi = 160
				} else if pixelRatio <= 0.75 {
					dpi = 120
				} else {
					dpi = 72
				}

				//fontsize默认是在dpi为72时候的尺寸，但是dpi也是依据PixelRatio计算除了的，所以fontsize也需要做下变换
				parentSize := [2]int{canvas.GetImage().Bounds().Dx(), canvas.GetImage().Bounds().Dy()}
				isNone, fontHeight, err1 := utils.ParseNumPercentNumNone(a.FontSizeByHeightPercent, parentSize[1])
				if err1 != nil {
					err = err1
					return
				}
				if isNone {
					err = fmt.Errorf("FontSizeByHeightPercent=" + a.FontSizeByHeightPercent + " not valid")
					return
				}

				textFont, err1 := utils.LoadFont(a.FontFile) //文字库文件
				if err1 != nil {
					err = err1
					return
				}
				fontSize, err1 := utils.GetFontSizeByHeight(textFont, dpi, fontHeight)
				if err1 != nil {
					err = err1
					return
				}

				//fmt.Println("dpi=", dpi, ",fontSize=", fontSize)
				//fontSize := float64(a.FontSize) / pixelRatio
				err = subImg.LoadFromText(a.Text, maxWidth, a.FontFile, dpi, fontSize, a.FontColor, a.FontBackgroundColor)
				if err != nil {
					return
				}
			//类型为resize的重新调整大小
			case ActionTypeResize:
				subImg = imageMap[s.ID]
				isNone, ResizeWidth, ResizeHeight := true, 0, 0
				switch a.ResizeType {
				case ResizeTypeByWidthAndHeight:
					isNone, ResizeWidth, err = utils.ParseNumPercentNumNone(a.ResizeWidth, canvas.GetImage().Bounds().Dx())
					if err != nil {
						return
					}
					if isNone {
						err = fmt.Errorf("ResizeWidth not allow none")
						return
					}
					isNone, ResizeHeight, err = utils.ParseNumPercentNumNone(a.ResizeHeight, canvas.GetImage().Bounds().Dy())
					if err != nil {
						return
					}
					if isNone {
						err = fmt.Errorf("ResizeHeight not allow none")
						return
					}
				case ResizeTypeByWidth:
					isNone, ResizeWidth, err = utils.ParseNumPercentNumNone(a.ResizeWidth, canvas.GetImage().Bounds().Dx())
					if err != nil {
						return
					}
					if isNone {
						err = fmt.Errorf("ResizeWidth not allow none")
						return
					}
					ResizeHeight = ResizeWidth
				case ResizeTypeByHeight:
					isNone, ResizeHeight, err = utils.ParseNumPercentNumNone(a.ResizeHeight, canvas.GetImage().Bounds().Dy())
					if err != nil {
						return
					}
					if isNone {
						err = fmt.Errorf("ResizeHeight not allow none")
						return
					}
					ResizeWidth = ResizeHeight
				default:
					err = fmt.Errorf("not support resize type=%d", a.ResizeType)
					return
				}

				err = subImg.Resize(uint(ResizeWidth), uint(ResizeHeight))
				if err != nil {
					return
				}
			//裁剪矩形圆角
			case ActionTypeRound:
				subImg = imageMap[s.ID]
				isNone, RoundRadius := true, 0
				isNone, RoundRadius, err = utils.ParseNumPercentNumNone(a.RoundRadius, subImg.GetImage().Bounds().Dx())
				if err != nil {
					return
				}
				if isNone {
					err = fmt.Errorf("RoundRadius not allow none")
					return
				}
				err = subImg.DrawRoundImage(float64(RoundRadius))
				if err != nil {
					return
				}
			//裁剪圆
			case ActionTypeCircle:
				subImg = imageMap[s.ID]
				isNone, CircleOriginX, CircleOriginY, CircleRadius := true, 0, 0, 0
				isNone, CircleOriginX, err = utils.ParseNumPercentNumNone(a.CircleOriginX, subImg.GetImage().Bounds().Dx())
				if err != nil {
					return
				}
				if isNone {
					err = fmt.Errorf("CircleOriginX not allow none")
					return
				}
				isNone, CircleOriginY, err = utils.ParseNumPercentNumNone(a.CircleOriginY, subImg.GetImage().Bounds().Dx())
				if err != nil {
					return
				}
				if isNone {
					err = fmt.Errorf("CircleOriginY not allow none")
					return
				}
				isNone, CircleRadius, err = utils.ParseNumPercentNumNone(a.CircleRadius, subImg.GetImage().Bounds().Dx())
				if err != nil {
					return
				}
				if isNone {
					err = fmt.Errorf("CircleRadius not allow none")
					return
				}
				err = subImg.DrawCircleImage(CircleOriginX, CircleOriginY, CircleRadius)
				if err != nil {
					return
				}
			//定位且叠画子图到画布中
			case ActionTypeLocation:
				subImg = imageMap[s.ID]
				parentSize := [2]int{canvas.GetImage().Bounds().Dx(), canvas.GetImage().Bounds().Dy()}
				selfSize := [2]int{subImg.GetImage().Bounds().Dx(), subImg.GetImage().Bounds().Dy()}
				minX, minY, maxX, maxY := 0, 0, 0, 0
				debugColor := "" //如果为空，则代表不需要debug开启，debug背景色为空
				if imageConfigInfo.Debug {
					debugColor = mgr.GetDebugColor(s.ID)
				}

				switch a.LocationType {
				//绝对定位
				case LocationTypeAbsolute:
					minX, minY, maxX, maxY, err = utils.ParseAbsoluteLocation(a.AbsolutePosition, parentSize, selfSize)
					if err != nil {
						return
					}
					_, _ = maxX, maxY
					err = canvas.DrawSubImage(subImg, minX, minY, debugColor)
					if err != nil {
						return
					}
				//相对定位
				case LocationTypeRelative:
					var relativePosition [4]int
					relativePosition, err = canvas.GetSubImagePosition(a.RelativeToImageID)
					if err != nil {
						return
					}
					minX, minY, maxX, maxY, err = utils.ParseRelativeLocation(a.RelativePosition, parentSize, selfSize, relativePosition)
					if err != nil {
						return
					}
					_, _ = maxX, maxY
					err = canvas.DrawSubImage(subImg, minX, minY, debugColor)
					if err != nil {
						return
					}
				//混合定位
				case LocationTypeMixed:
					var relativePosition [4]int

					absoluteMinX, absoluteMinY, absoluteMaxX, absoluteMaxY := 0, 0, 0, 0
					relativeMinX, relativeMinY, relativeMaxX, relativeMaxY := 0, 0, 0, 0
					absoluteMinX, absoluteMinY, absoluteMaxX, absoluteMaxY, err = utils.ParseAbsoluteLocation(a.AbsolutePosition, parentSize, selfSize)
					if err != nil {
						return
					}
					_, _ = absoluteMaxX, absoluteMaxY

					relativePosition, err = canvas.GetSubImagePosition(a.RelativeToImageID)
					if err != nil {
						return
					}
					relativeMinX, relativeMinY, relativeMaxX, relativeMaxY, err = utils.ParseRelativeLocation(a.RelativePosition, parentSize, selfSize, relativePosition)
					if err != nil {
						return
					}
					_, _ = relativeMaxX, relativeMaxY
					minX, minY = utils.MinInt(absoluteMinX, relativeMinX), utils.MinInt(absoluteMinY, relativeMinY)

					err = canvas.DrawSubImage(subImg, minX, minY, debugColor)
					if err != nil {
						return
					}
				}

				//fmt.Println(fmt.Sprintf("[DRAW] %s(%d,%d) at %s(%d,%d) at (%d,%d)",
				//	subImg.ID(), subImg.GetImage().Bounds().Dx(), subImg.GetImage().Bounds().Dy(),
				//	canvas.ID(), canvas.GetImage().Bounds().Dx(), canvas.GetImage().Bounds().Dy(),
				//	minX, minY,
				//))
			default:
				err = fmt.Errorf("not support action type=" + a.ActionType)
				return
			} // end switch a.ActionType

			if idx == 0 {
				subImg.SetID(s.ID)
				imageMap[s.ID] = subImg
			}
		} // end for idx, a := range s.ActionList
		//configMap[s.ID] = s
	}

	imgBuffer, err = canvas.SaveToBuffer()
	if err != nil {
		return
	}

	return
}
