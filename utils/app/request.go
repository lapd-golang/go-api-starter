package app

import (
	"github.com/astaxie/beego/validation"
	"go-admin-starter/utils"
)

func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		utils.Log.Warn(err.Key, err.Message)
	}

	return
}
