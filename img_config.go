package goposter

const (
	ActionTypeNone     = "none"     // 作图类型-无
	ActionTypeCanvas   = "canvas"   // 作图类型-画布
	ActionTypeFont     = "font"     // 作图类型-画字
	ActionTypeImage    = "image"    // 作图类型-加载子图
	ActionTypeResize   = "resize"   // 作图类型-调整大小
	ActionTypeRound    = "round"    // 作图类型-裁剪矩形圆角
	ActionTypeCircle   = "circle"   // 作图类型-裁剪圆
	ActionTypeLocation = "location" // 作图类型-定位且叠画子图

	ImageTypeNone = "none" //图片类型-无
	ImageTypePNG  = "png"  //图片类型-png
	ImageTypeJPEG = "jpeg" //图片类型-jpg或者jpeg

	ResizeTypeByNone           = 0 //调整大小-无
	ResizeTypeByWidthAndHeight = 1 //调整大小-根据宽高的调整
	ResizeTypeByWidth          = 2 //调整大小-只根据宽度调整
	ResizeTypeByHeight         = 3 //调整大小-只根据高度调整

	LocationTypeNone     = "none"     //定位类型-无
	LocationTypeAbsolute = "absolute" //定位类型-绝对定位
	LocationTypeRelative = "relative" //定位类型-相对定位
	LocationTypeMixed    = "mixed"    //定位类型-混合定位
)

type ImageActionType struct {
	ActionType string `json:"action_type"` // 类型, font:文字,image:图片加载,resize:重新规划尺寸,round:圆角，circle:圆形，location：摆放在画布位置

	//ActionType==font时候的字段
	Text                    string `json:"text"`                        // 文本
	FontBackgroundColor     string `json:"font_background_color"`       // 文本背景颜色，（只支持#FFFFFF格式)
	FontFile                string `json:"font_file"`                   // 字体文件，必需字段
	FontSizeByHeightPercent string `json:"font_size_by_height_percent"` // 字体大小,画板高度的百分比
	FontColor               string `json:"font_color"`                  // 字体颜色（只支持#FFFFFF格式)，必需字段
	MaxWidth                string `json:"max_width"`                   // 字体所占的最大宽度（为0不限制)，必需字段

	//ActionType==image时候的字段(image只支持png,jpeg）
	ImageType       string `json:"image_type"`        //图片类型，image_type=png|jpeg
	ImageLocalFile  string `json:"image_local_file"`  //图片本地路径（为空代表不是本地图片),比如: /data/image/aa.jpeg
	ImageURLFile    string `json:"image_url_file"`    //图片url路径（为空代表不是本地图片),比如: https://baidu.com/data/image/aa.jpeg
	ImageCustomFile string `json:"image_custom_file"` //图片自定义来源，自定义方法获取，此类别让业务按照要求来使用，此处不启用

	//ActionType==resize时候的字段
	ResizeType   int    `json:"resize_type"`   // 调整大小的规则,0:none,1:根据指定宽高,2:只依据宽度,3:只依据高度
	ResizeWidth  string `json:"resize_width"`  // 图宽度,整数或者画布百分数
	ResizeHeight string `json:"resize_height"` // 图高度度,整数或者画布百分数

	//ActionType==round时候的字段
	RoundRadius string `json:"round_radius"` // 图的矩形圆角半径,整数或者百分数，如果是百分数，代表基于本图width的圆角半径

	//ActionType==circle时候的字段
	CircleOriginX string `json:"circle_origin_x"` // 图的原点x，整数或者百分数（百分数的话，是基于本图，而不是父图的），none代表不指定，由子图自计算
	CircleOriginY string `json:"circle_origin_y"` // 图的原点y，整数或者百分数
	CircleRadius  string `json:"circle_radius"`   // 图的半径，整数或者百分数

	//ActionType==location时候的字段
	LocationType      string `json:"location_type"`     // 定位类型，absolute|relative|mixed
	RelativeToImageID string `json:"sub_image_id"`      // 图id，location_type==relative(即相对定位）为同级图id
	AbsolutePosition  string `json:"absolute_position"` // 绝对定位时候的四部定位(顺时针方向：上右底左)，每部可为整数或者百分数(none代表忽略定位），至少需要两维信息，举例:10 10 0 0=> 离父元素顶部10，右部10摆放子图
	RelativePosition  string `json:"relative_position"` // 相对定位时候的四部定位(顺时针方向：上右底左)，每部可为整数或者百分数(none代表忽略定位），至少需要两维信息，举例:10 10 0 0=> 离定位元素顶部10，右部10摆放子图在画布中
}

// 子图图像
type SubImageConfigInfoType struct {
	ID         string            `json:"id"`          // id
	Name       string            `json:"name"`        // 名称
	Desc       string            `json:"desc"`        // 备注说明
	Enable     string            `json:"enable"`      // 是否启用,字符串类型，以方便做模板渲染，只能是true|false
	Other      string            `json:"other"`       // 其他信息，不参与到作图中，只是预留给业务使用
	ActionList []ImageActionType `json:"action_list"` // 动作说明

	//计算值，其是根据ActionList计算出来的坐标值
	//Top    int `json:"top"`    // 计算值-四角定位在父元素的左上角y坐标
	//Left   int `json:"top"`    // 计算值-四角定位在父元素的左上角x坐标
	//Bottom int `json:"bottom"` // 计算值-四角定位在父元素的右下角y坐标
	//Right  int `json:"right"`  // 计算值-四角定位在父元素的右下角x坐标
}

// 图像
type ImageConfigInfoType struct {
	ID                    string                   `json:"id"`                      // id
	Name                  string                   `json:"name"`                    // 名称
	Version               string                   `json:"version"`                 // 配置版本，比如版本1位v1.0.0后续可能添加了某个subimage，此时可以改动此版本以便重新生成
	Desc                  string                   `json:"desc"`                    // 备注说明
	Enable                string                   `json:"enable"`                  // 是否启用，注意此处为字符串类型,true为启用，其他为不启用
	Other                 string                   `json:"other"`                   // 其他信息，不参与到作图中，只是预留给业务使用
	Debug                 bool                     `json:"debug"`                   // 是否开启调试，如果开启了，则会给子图加入调式色块，以便于辨识面积和位置
	PixelRatio            string                   `json:"pixel_ratio"`             // 设备像素比，一个浮点数
	Width                 string                   `json:"width"`                   // 画布宽度,数字
	Height                string                   `json:"height"`                  // 画布高度,数字
	CanvasBackgroundColor string                   `json:"canvas_background_color"` //画布背景颜色，（只支持#FFFFFF格式)
	OutputBufferType      string                   `json:"output_buffer_type"`      // 输出图片类型，png|jpeg
	SubImageInfoList      []SubImageConfigInfoType `json:"sub_image_info_list"`     //子图列表
}
