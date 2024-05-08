package goposter

import (
	"image"
	"image/color"
	"math"
)

//圆形效果
type CircleMask struct {
	image   image.Image //在这个image的基础上圆形（要想得到最好的裁剪结果，这个image需要是正方形才最好）
	originX int         // 圆形的圆心x轴（如果为0，则宽度一半坐标）
	originY int         // 圆形的圆心y轴（如果为0，则高度一半坐标）
	radius  int         // 圆形半径（如果为0，则使用image最短边为半径）
}

func (m *CircleMask) ColorModel() color.Model {
	return m.image.ColorModel()
}

func (m *CircleMask) Bounds() image.Rectangle {
	return m.image.Bounds()
}

func (m *CircleMask) At(x, y int) color.Color {
	/*

		.-------------------+->x
		|					|
		|					|
		|		  +(origin) |
		|					|
		|				 	|
		+-------------------+
		|		  |  radius |
		|
		y
	*/

	if m.originX == 0 || m.originY == 0 || m.radius == 0 {
		dx := m.image.Bounds().Dx()
		dy := m.image.Bounds().Dy()
		m.radius = int(dx / 2)
		if dy < dx {
			m.radius = int(dy / 2)
		}
		m.originX = int(dx / 2)
		m.originY = int(dy / 2)
	}

	dis := math.Sqrt(math.Pow(float64(x-m.originX), 2) + math.Pow(float64(y-m.originY), 2))
	if dis > float64(m.radius) {
		return color.Transparent
	}
	return m.image.At(x, y)
}
