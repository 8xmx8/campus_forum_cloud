package common

import (
	"github.com/importcjj/sensitive"
)

func InitSensitiveFilter() (*sensitive.Filter, error) {
	snf := sensitive.New()
	err := snf.LoadWordDict("./common/sensitive_word_dic.txt")
	if err != nil {
		return nil, err
	}
	return snf, nil
}

func Validate(sn *sensitive.Filter, word string) bool {
	is, _ := sn.Validate(word)
	return !is
}
