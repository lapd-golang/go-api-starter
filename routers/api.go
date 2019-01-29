package routers

import (
	"github.com/gin-gonic/gin"
	"go-admin-starter/controllers/api"
	"go-admin-starter/controllers/api/v1"
	"go-admin-starter/middleware/authrole"
	"go-admin-starter/middleware/jwt"
)

func initApiRouter(r *gin.Engine) *gin.Engine {
	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/auth", api.GetAuth)
		apiGroup.POST("/register", api.Register)
		apiGroup.POST("/refreshToken", api.RefreshToken)//前后台共用

		apiv1 := apiGroup.Group("/v1")
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

		//后台管理api
		admin := apiGroup.Group("/admin")
		admin.POST("/auth", api.AdminGetAuth)

		admin.Use(jwt.JWTAuth())
		admin.Use(authrole.AuthCheckRole())
		{
			admin.POST("/addrole", v1.AddCasbin)//添加角色权限
		}
	}

	return r
}
