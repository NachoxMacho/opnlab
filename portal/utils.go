package portal

import (
	"fmt"
	"math"
	"strconv"
)

func HumanFileSize(size float64) string {
	fmt.Println(size)
	suffixes := []string{
		"B",
		"KB",
		"MB",
		"GB",
		"TB",
	}

	if size < 1 {
		return "0B"
	}

	base := math.Log(size) / math.Log(1024)
	getSize := Round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	fmt.Println(int(math.Floor(base)))
	getSuffix := suffixes[int(math.Floor(base))]
	return strconv.FormatFloat(getSize, 'f', -1, 64) + " " + string(getSuffix)
}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
