package routers

import (
	"admin-server/controllers/api"
	"admin-server/controllers/api/v1"
	"admin-server/middleware/jwt"
	"github.com/gin-gonic/gin"
)

func initApiRouter(r *gin.Engine) *gin.Engine {
	r.POST("/auth", api.GetAuth)
	r.POST("/refreshToken", api.RefreshToken)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWTAuth())
	{
		//获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		//新建标签
		apiv1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)

		//获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		//获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		//新建文章
		apiv1.POST("/articles", v1.AddArticle)
		//更新指定文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		//删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return r
}
