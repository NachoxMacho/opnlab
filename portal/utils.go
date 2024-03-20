package portal

import (
	"context"
	"encoding/json"
	"math"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func HumanFileSize(size float64) string {
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

func getObjectsFromCache[T any](client *redis.Client, key string) ([]T, error) {
	var objs []T
	result, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(result), &objs)
	if err != nil {
		return nil, err
	}

	return objs, nil
}
