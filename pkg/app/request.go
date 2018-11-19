package app

import (
	"admin-server/pkg/util"
	"github.com/astaxie/beego/validation"
)

func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		util.Log.Warn(err.Key, err.Message)
	}

	return
}
