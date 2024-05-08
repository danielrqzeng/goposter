package goposter

import (
	"goposter/utils"
	"math"
	"testing"
)

func TestParseNumPercentNumNone(t *testing.T) {
	for idx, unit := range []struct {
		//params
		numStr   string
		baseSize int
		//return
		isNone   bool
		realSize int
	}{
		{"100", 100, false, 100},
		{"10%", 100, false, 10},
		{"none", 100, true, 0},
		{"23%", 100, false, 23},
		//{"", 100, false, 23},
	} {
		isNone, realSize, err := utils.ParseNumPercentNumNone(unit.numStr, unit.baseSize)
		if err != nil {
			t.Fatal(err)
		}
		if isNone != unit.isNone {
			t.Errorf("idx=%d expected isNone:[%v] , actually: [%v]", idx, unit.isNone, isNone)
		}
		if realSize != unit.realSize {
			t.Errorf("idx=%d expected realSize:[%v] , actually: [%v]", idx, unit.realSize, realSize)
		}
	}
}

func TestParseAbsoluteLocation(t *testing.T) {
	for idx, unit := range []struct {
		//params
		positionStr string
		parentSize  [2]int
		selfSize    [2]int
		//return
		minX, minY, maxX, maxY int
	}{
		//定位左上角
		{"10% none none 1%", [2]int{320, 640}, [2]int{40, 10}, 3, 64, 43, 74},
		//定位右上角
		{"10% 1% none none", [2]int{320, 640}, [2]int{40, 10}, 320 - 3 - 40, 64, 320 - 3, 74},
		//定位右下角
		{"none 1% 10% none", [2]int{320, 640}, [2]int{40, 10}, 320 - 3 - 40, 640 - 64 - 10, 320 - 3, 640 - 64},
		//定位左下角
		{"none none 10% 1%", [2]int{320, 640}, [2]int{40, 10}, 3, 640 - 64 - 10, 40 + 3, 640 - 64},
		{"none none 64 3", [2]int{320, 640}, [2]int{40, 10}, 3, 640 - 64 - 10, 40 + 3, 640 - 64},
		//只有x维度
		{"none none none 3", [2]int{320, 640}, [2]int{40, 10}, 3, math.MaxInt32, 40 + 3, math.MaxInt32},
		//只有y维度
		{"none none 10% none", [2]int{320, 640}, [2]int{40, 10}, math.MaxInt32, 640 - 64 - 10, math.MaxInt32, 640 - 64},
		//格式错误
		//{"none none ", [2]int{320, 640}, [2]int{40, 10}, -1, 640 - 64 - 10, -1, 640 - 64},
	} {
		minX, minY, maxX, maxY, err := utils.ParseAbsoluteLocation(unit.positionStr, unit.parentSize, unit.selfSize)
		if err != nil {
			t.Fatal(err)
		}
		if minX != unit.minX {
			t.Errorf("idx=%d expected minX:[%v] , actually: [%v]", idx, unit.minX, minX)
		}
		if minY != unit.minY {
			t.Errorf("idx=%d expected minY:[%v] , actually: [%v]", idx, unit.minY, minY)
		}
		if maxX != unit.maxX {
			t.Errorf("idx=%d expected maxX:[%v] , actually: [%v]", idx, unit.maxX, maxX)
		}
		if maxY != unit.maxY {
			t.Errorf("idx=%d expected maxY:[%v] , actually: [%v]", idx, unit.maxY, maxY)
		}

	}
}

func TestParseRelativeLocation(t *testing.T) {
	for idx, unit := range []struct {
		//params
		positionStr      string
		parentSize       [2]int
		selfSize         [2]int
		relativePosition [4]int
		//return
		minX, minY, maxX, maxY int
	}{
		//相对于相对体，本体往下10%,往左1%，即本体在相对体的左下方,以相对体的左下角(11,24)做偏移
		{"10% none none 1%", [2]int{320, 640}, [2]int{40, 10}, [4]int{11, 12, 23, 24}, 11 - 3 - 40, 24 + 64, 11 - 3, 24 + 64 + 10},
		//相对于相对体，本体往下10%,往右1%，即本体在相对体的右下方,以相对体的左下角(23,24)做偏移
		{"10% 1% none none", [2]int{320, 640}, [2]int{40, 10}, [4]int{11, 12, 23, 24}, 23 + 3, 24 + 64, 23 + 3 + 40, 24 + 64 + 10},
		//相对于相对体，本体往上10%,往右1%，即本体在相对体的右上方,以相对体的右上角(23,12)做偏移
		{"none 1% 10% none", [2]int{320, 640}, [2]int{40, 10}, [4]int{11, 12, 23, 24}, 23 + 3, 12 - 64 - 10, 23 + 3 + 40, 12 - 64},
		//相对于相对体，本体往上10%,往左1%，即本体在相对体的左上方,以相对体的左上角(11,12)做偏移
		{"none none 10% 1%", [2]int{320, 640}, [2]int{40, 10}, [4]int{11, 12, 23, 24}, 11 - 3 - 40, 12 - 64 - 10, 11 - 3, 12 - 64},
		{"none none 64 3", [2]int{320, 640}, [2]int{40, 10}, [4]int{11, 12, 23, 24}, 11 - 3 - 40, 12 - 64 - 10, 11 - 3, 12 - 64},
		//只有x维度
		{"none none none 3", [2]int{320, 640}, [2]int{40, 10}, [4]int{11, 12, 23, 24}, 11 - 3 - 40, math.MaxInt32, 11 - 3, math.MaxInt32},
		////只有y维度
		//{"none none 10% none", [2]int{320, 640}, [2]int{40, 10}, [4]int{11, 12, 23, 24}, 23 + 3, 24 + 64, 23 + 3 + 40, 24 + 64 + 10},
		//格式错误
		//{"none none ", [2]int{320, 640}, [2]int{40, 10}, [4]int{11, 12, 23, 24}, 23 + 3, 24 + 64, 23 + 3 + 40, 24 + 64 + 10},
	} {
		minX, minY, maxX, maxY, err := utils.ParseRelativeLocation(unit.positionStr, unit.parentSize, unit.selfSize, unit.relativePosition)
		if err != nil {
			t.Fatal(err)
		}
		if minX != unit.minX {
			t.Errorf("idx=%d expected minX:[%v] , actually: [%v]", idx, unit.minX, minX)
		}
		if minY != unit.minY {
			t.Errorf("idx=%d expected minY:[%v] , actually: [%v]", idx, unit.minY, minY)
		}
		if maxX != unit.maxX {
			t.Errorf("idx=%d expected maxX:[%v] , actually: [%v]", idx, unit.maxX, maxX)
		}
		if maxY != unit.maxY {
			t.Errorf("idx=%d expected maxY:[%v] , actually: [%v]", idx, unit.maxY, maxY)
		}

	}
}
