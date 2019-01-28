package v1

import (
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"go-admin-starter/models"
	"go-admin-starter/utils"
	"go-admin-starter/utils/app"
	"go-admin-starter/utils/config"
	"go-admin-starter/utils/e"
)

//获取多个文章标签
func GetTags(c *gin.Context) {
	name := c.Query("name")

	var user models.Tag
	if name != "" {
		user.Name = name
	}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		user.State = state
	}

	conf := config.New()
	data := make(map[string]interface{})
	data["lists"] = user.Get(utils.GetPage(c), conf.App.PageSize)
	data["total"] = user.GetTotal()

	app.Response(c, e.SUCCESS, "ok", data)
}

//新增文章标签
func AddTag(c *gin.Context) {
	name := c.PostForm("name")
	state := com.StrTo(c.DefaultPostForm("state", "0")).MustInt()
	createdBy := c.PostForm("created_by")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		app.Response(c, e.INVALID_PARAMS, valid.Errors[0].Message, nil)
		return
	}

	var tag models.Tag

	if ! tag.ExistByName(name) {
		tag.Name = name
		tag.CreatedBy = createdBy
		tag.State = state
		id, err := tag.Insert()

		if err != nil {
			app.Response(c, e.ERROR, "添加失败", nil)
			return
		}

		app.Response(c, e.SUCCESS, "ok", id)
		return
	} else {
		app.Response(c, e.ERROR_EXIST_TAG, "已存在该标签名称", nil)
		return
	}
}

//修改文章标签
func EditTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.PostForm("name")
	modifiedBy := c.PostForm("modified_by")

	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		app.Response(c, e.INVALID_PARAMS, valid.Errors[0].Message, nil)
		return
	}

	var tag models.Tag

	if tag.ExistByID(id) {
		tag.ModifiedBy = modifiedBy
		if name != "" {
			tag.Name = name
		}
		if state != -1 {
			tag.State = state
		}

		result, err := tag.Update(id)
		if err != nil || result.ID == 0 {
			app.Response(c, e.ERROR, "修改失败", nil)
			return
		}
		app.Response(c, e.SUCCESS, "ok", nil)
		return
	} else {
		app.Response(c, e.ERROR_NOT_EXIST_TAG, "该标签不存在", nil)
		return
	}
}

//删除文章标签
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("请求参数错误")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		app.Response(c, e.INVALID_PARAMS, valid.Errors[0].Message, nil)
		return
	}

	var tag models.Tag
	if tag.ExistByID(id) {
		result, err := tag.Delete(id)
		if err != nil || result.ID == 0 {
			app.Response(c, e.ERROR, "删除失败", nil)
			return
		}

		app.Response(c, e.SUCCESS, "ok", nil)
		return
	} else {
		app.Response(c, e.ERROR_NOT_EXIST_TAG, "该标签不存在", nil)
		return
	}
}
