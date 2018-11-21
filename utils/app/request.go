package app

import (
	"admin-server/utils"
	"github.com/astaxie/beego/validation"
)

func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		utils.Log.Warn(err.Key, err.Message)
	}

	return
}
