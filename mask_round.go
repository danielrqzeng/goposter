package goposter

import (
	"image"
	"image/color"
	"math"
)

//圆角效果
type RoundMask struct {
	image  image.Image //在这个image的基础上做圆角
	radius int         //圆角半径
}

func (m *RoundMask) ColorModel() color.Model {
	return m.image.ColorModel()
}

func (m *RoundMask) Bounds() image.Rectangle {
	return m.image.Bounds()
}

func (m *RoundMask) At(x, y int) color.Color {
	/*

		.-------------------+->x
		|					|
		|					|
		|					|
		|					|
		|					|
		|				 	|
		+----------------+--+ ==>四个角的radius
		|			    |radius|
		|
		y
	*/

	maxX := m.image.Bounds().Max.X - 1
	maxY := m.image.Bounds().Max.Y - 1

	//左上角计算
	if x < m.radius && y < m.radius {

		//移动到圆心为原点的地方做计算，其距离圆心，是否在半径范围内
		rx := -(m.radius - x)
		ry := -(m.radius - y)
		dis := math.Sqrt(math.Pow(float64(rx), 2) + math.Pow(float64(ry), 2))
		if dis > float64(m.radius) {
			//fmt.Println("up-left m.radius=", m.radius, "(x,y)=(", x, ",", y, ")/(", maxX, ",", maxY, "),", " a=0")
			return color.Transparent // 全透明

		} else {
			//fmt.Println("up-left m.radius=", m.radius, "(x,y)=(", x, ",", y, ")/(", maxX, ",", maxY, "),", " a=raw")
			return m.image.At(x, y) // 用原来的颜色
		}
	}

	//右上角
	if x > maxX-m.radius && y < m.radius {
		//移动到圆心为原点的地方做计算，其距离圆心，是否在半径范围内
		rx := x - (maxX - m.radius)
		ry := -(m.radius - y)
		dis := math.Sqrt(math.Pow(float64(rx), 2) + math.Pow(float64(ry), 2))
		if dis > float64(m.radius) {
			//fmt.Println("up-right m.radius=", m.radius, "(x,y)=(", x, ",", y, ")/(", maxX, ",", maxY, "),", " a=0")
			return color.Transparent // 全透明

		} else {
			//fmt.Println("up-right m.radius=", m.radius, "(x,y)=(", x, ",", y, ")/(", maxX, ",", maxY, "),", "raw")

			return m.image.At(x, y) // 用原来的颜色
		}
	}

	//左下角
	if x < m.radius && y > maxY-m.radius {
		//fmt.Println("down-left (x,y)=(", x, ",", y, ")")

		//移动到圆心为原点的地方做计算，其距离圆心，是否在半径范围内
		rx := -(m.radius - x)
		ry := y - (maxY - m.radius)
		dis := math.Sqrt(math.Pow(float64(rx), 2) + math.Pow(float64(ry), 2))
		if dis > float64(m.radius) {
			//fmt.Println("down-left m.radius=", m.radius, "(x,y)=(", x, ",", y, ")/(", maxX, ",", maxY, "),", " a=0")
			return color.Transparent // 全透明
		} else {
			//fmt.Println("down-left m.radius=", m.radius, " (x,y)=(", x, ",", y, ")/(", maxX, ",", maxY, "),", " raw")
			return m.image.At(x, y) // 用原来的颜色
		}
	}

	//右下角
	if x > maxX-m.radius && y > maxY-m.radius {

		//移动到圆心为原点的地方做计算，其距离圆心，是否在半径范围内
		rx := x - (maxX - m.radius)
		ry := y - (maxY - m.radius)
		dis := math.Sqrt(math.Pow(float64(rx), 2) + math.Pow(float64(ry), 2))
		if dis > float64(m.radius) {
			//fmt.Println("down-right m.radius=", m.radius, "(x,y)=(", x, ",", y, ")/(", maxX, ",", maxY, "),", "a=0")
			return color.Transparent // 全透明

		} else {
			//fmt.Println("down-right m.radius=", m.radius, "(x,y)=(", x, ",", y, ")/(", maxX, ",", maxY, "),", " raw")
			return m.image.At(x, y) // 用原来的颜色
		}
	}
	return m.image.At(x, y) // 用原来的颜色

}
