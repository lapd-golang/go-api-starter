package v1

import (
	"admin-server/models"
	"admin-server/pkg/app"
	"admin-server/pkg/config"
	"admin-server/pkg/e"
	"admin-server/pkg/logging"
	"admin-server/pkg/upload"
	"admin-server/pkg/util"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @Summary 获取单个文章
// @Produce json
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles/{id} [get]
func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID错误")

	if valid.HasErrors() {
		logging.Info(valid.Errors)
		app.Response(c, e.INVALID_PARAMS, valid.Errors[0].Message, nil)
		return
	}

	article := models.Article{}
	if article.ExistByID(id) {
		data := article.GetById(id)
		app.Response(c, e.SUCCESS, "ok", data)
		return
	} else {
		app.Response(c, e.ERROR_NOT_EXIST_ARTICLE, "该文章不存在", nil)
		return
	}
}

// @Summary 获取多个文章
// @Produce json
// @Success 200 {string} json "{"code":200,"data":{"lists": [], "total": 0},"msg":"ok"}"
// @Router /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state

		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId

		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	}

	if valid.HasErrors() {
		logging.Info(valid.Errors)
		app.Response(c, e.INVALID_PARAMS, valid.Errors[0].Message, nil)
		return
	}

	article := models.Article{}

	data["lists"] = article.Get(util.GetPage(c), config.AppSetting.PageSize, maps)
	data["total"] = article.GetTotal(maps)

	app.Response(c, e.SUCCESS, "ok", data)
	return
}

// @Summary 新增文章
// @Produce  json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":ID,"msg":"ok"}"
// @Router /api/v1/articles [post]
func AddArticle(c *gin.Context) {
	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()
	title := c.PostForm("title")
	desc := c.PostForm("desc")
	content := c.PostForm("content")
	createdBy := c.PostForm("created_by")
	state := com.StrTo(c.DefaultPostForm("state", "0")).MustInt()
	coverImageUrl, header, _ := c.Request.FormFile("cover_image_url")

	valid := validation.Validation{}
	valid.Min(tagId, 1, "tag_id").Message("标签ID错误")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(desc, "desc").Message("简述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(coverImageUrl, "cover_image_url").Message("封面地址不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	if valid.HasErrors() {
		logging.Info(valid.Errors)

		app.Response(c, e.INVALID_PARAMS, valid.Errors[0].Message, nil)
		return
	}

	tag := models.Tag{}
	article := models.Article{}
	if tag.ExistByID(tagId) {
		//save image
		imageName := upload.GetImageName(header.Filename)
		path := upload.GetImagePath()
		fullPath := path + imageName
		err := upload.SaveImage(coverImageUrl, fullPath)
		if err != nil {
			app.Response(c, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, err.Error(), nil)
			return
		}

		data := make(map[string]interface{})
		data["tag_id"] = tagId
		data["title"] = title
		data["desc"] = desc
		data["content"] = content
		data["cover_image_url"] = fullPath
		data["created_by"] = createdBy
		data["state"] = state

		id := article.Add(data)

		app.Response(c, e.SUCCESS, "ok", id)
		return
	} else {
		app.Response(c, e.ERROR_NOT_EXIST_TAG, "该标签不存在", nil)
		return
	}
}

// @Summary 修改文章
// @Produce  json
// @Param id param int true "ID"
// @Param name query string true "ID"
// @Param state query int false "State"
// @Param modified_by query string true "ModifiedBy"
// @Success 200 {string} json "{"code":200,"data":null,"msg":"ok"}"
// @Router /api/v1/articles/{id} [put]
func EditArticle(c *gin.Context) {
	valid := validation.Validation{}

	id := com.StrTo(c.Param("id")).MustInt()
	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()
	title := c.PostForm("title")
	desc := c.PostForm("desc")
	content := c.PostForm("content")
	modifiedBy := c.PostForm("modified_by")

	var state int = -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Min(id, 1, "id").Message("ID错误")
	valid.MaxSize(title, 100, "title").Message("标题最长为100字符")
	valid.MaxSize(desc, 255, "desc").Message("简述最长为255字符")
	valid.MaxSize(content, 65535, "content").Message("内容最长为65535字符")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")

	if valid.HasErrors() {
		logging.Info(valid.Errors)
		app.Response(c, e.INVALID_PARAMS, valid.Errors[0].Message, nil)
		return
	}

	article := models.Article{}
	tag := models.Tag{}
	if article.ExistByID(id) {
		if tag.ExistByID(tagId) {
			data := make(map[string]interface{})
			if tagId > 0 {
				data["tag_id"] = tagId
			}
			if title != "" {
				data["title"] = title
			}
			if desc != "" {
				data["desc"] = desc
			}
			if content != "" {
				data["content"] = content
			}

			data["modified_by"] = modifiedBy

			article.Edit(id, data)
			app.Response(c, e.SUCCESS, "ok", nil)
			return
		} else {
			app.Response(c, e.ERROR_NOT_EXIST_TAG, "该标签不存在", nil)
			return
		}
	} else {
		app.Response(c, e.ERROR_NOT_EXIST_ARTICLE, "该文章不存在", nil)
		return
	}
}

// @Summary 删除文章
// @Produce json
// @Param id param int true "ID"
// @Success 200 {string} json "{"code":200,"data":null,"msg":"ok"}"
// @Router /api/v1/articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID错误")

	if valid.HasErrors() {
		logging.Info(valid.Errors)
		app.Response(c, e.INVALID_PARAMS, valid.Errors[0].Message, nil)
		return
	}

	article := models.Article{}
	if article.ExistByID(id) {
		article.Delete(id)
		app.Response(c, e.SUCCESS, "ok", nil)
		return
	} else {
		app.Response(c, e.ERROR_NOT_EXIST_ARTICLE, "该文章不存在", nil)
		return
	}
}
