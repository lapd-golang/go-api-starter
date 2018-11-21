package v1

import (
	"admin-server/models"
	"admin-server/pkg/app"
	"admin-server/pkg/config"
	"admin-server/pkg/e"
	"admin-server/pkg/upload"
	"admin-server/pkg/util"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @Summary 获取单个文章
// @Tags articles
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path string true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"message":"ok"}"
// @Router /api/v1/articles/{id} [get]
func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID错误")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
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
// @Tags articles
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param state query int false "State"
// @Param tag_id query int false "tag_id"
// @Success 200 {string} json "{"code":200,"data":{"lists": [], "total": 0},"message":"ok"}"
// @Router /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	valid := validation.Validation{}
	var article models.Article

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		article.State = state

		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		article.TagID = tagId

		valid.Min(tagId, 1, "tag_id").Message("标签ID错误")
	}

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		app.Response(c, e.INVALID_PARAMS, valid.Errors[0].Message, nil)
		return
	}

	data := make(map[string]interface{})
	data["lists"] = article.Get(util.GetPage(c), config.Conf.App.PageSize)
	data["total"] = article.GetTotal()

	app.Response(c, e.SUCCESS, "ok", data)
	return
}

// @Summary 新增文章
// @Tags articles
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param title formData string true "Title"
// @Param desc formData string true "Desc"
// @Param content formData string true "Content"
// @Param created_by formData string true "Created_by"
// @Param cover_image_url formData file true "cover_image_url"
// @Param created_by formData string true "CreatedBy"
// @Param state formData int false "State"
// @Success 200 {string} json "{"code":200,"data":ID,"message":"ok"}"
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
		app.MarkErrors(valid.Errors)
		app.Response(c, e.INVALID_PARAMS, valid.Errors[0].Message, nil)
		return
	}

	var tag models.Tag
	var article models.Article
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

		article.TagID = tagId
		article.Title = title
		article.Desc = desc
		article.Content = content
		article.CoverImageUrl = fullPath
		article.CreatedBy = createdBy
		article.State = state

		id, err := article.Insert()
		if err != nil {
			app.Response(c, e.ERROR, "添加失败", nil)
			return
		}

		app.Response(c, e.SUCCESS, "ok", id)
		return
	} else {
		app.Response(c, e.ERROR_NOT_EXIST_TAG, "该标签不存在", nil)
		return
	}
}

// @Summary 修改文章
// @Tags articles
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "ID"
// @Param tag_id formData int true "tag_id"
// @Param title formData string false "title"
// @Param desc formData string false "desc"
// @Param content formData string false "content"
// @Param modified_by formData string false "ModifiedBy"
// @Param state formData int false "State"
// @Success 200 {string} json "{"code":200,"data":null,"message":"ok"}"
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
		app.MarkErrors(valid.Errors)
		app.Response(c, e.INVALID_PARAMS, valid.Errors[0].Message, nil)
		return
	}

	var article models.Article
	var tag models.Tag
	if article.ExistByID(id) {
		if tag.ExistByID(tagId) {
			article.ModifiedBy = modifiedBy
			if tagId > 0 {
				article.TagID = tagId
			}
			if title != "" {
				article.Title = title
			}
			if desc != "" {
				article.Desc = desc
			}
			if content != "" {
				article.Content = content
			}

			result, err := article.Update(id)
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
	} else {
		app.Response(c, e.ERROR_NOT_EXIST_ARTICLE, "该文章不存在", nil)
		return
	}
}

// @Summary 删除文章
// @Tags articles
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "ID"
// @Success 200 {string} json "{"code":200,"data":null,"message":"ok"}"
// @Router /api/v1/articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID错误")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		app.Response(c, e.INVALID_PARAMS, valid.Errors[0].Message, nil)
		return
	}

	var article models.Article
	if article.ExistByID(id) {
		result, err := article.Delete(id)
		if err != nil || result.ID == 0 {
			app.Response(c, e.ERROR, "删除失败", nil)
			return
		}
		app.Response(c, e.SUCCESS, "ok", nil)
		return
	} else {
		app.Response(c, e.ERROR_NOT_EXIST_ARTICLE, "该文章不存在", nil)
		return
	}
}
