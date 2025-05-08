package utils

import (
	"github.com/spf13/cast"
)

func GetFirstChar(str string) int64 {
	if len(str) == 0 {
		return 0
	}
	level := cast.ToInt64(str[0:1])
	if level >= 1 && level <= 5 {
		return level
	}
	return 0
}
