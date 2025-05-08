package common

import (
	"testing"
)

func TestInitSensitiveFilter(t *testing.T) {
	sensitiveFilter, err := InitSensitiveFilter()
	if err != nil {
		t.Error(err)
	}

	is := Validate(sensitiveFilter, "傻逼")
	if is {
		t.Log("敏感词出现")
	}

}
