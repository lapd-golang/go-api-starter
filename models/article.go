package models

type Article struct {
	Base

	TagID int `json:"tag_id" gorm:"index"`
	Tag Tag `json:"tag"`

	Title string `json:"title"`
	Desc string `json:"desc"`
	Content string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	CreatedBy string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State int `json:"state"`
}

func ExistArticleByID(id int) bool {
	var article Article
	Eloquent.Select("id").Where("id = ?", id).First(&article)

	if article.ID > 0 {
		return true
	}

	return false
}

func GetArticleTotal(maps interface {}) (count int){
	Eloquent.Model(&Article{}).Where(maps).Count(&count)

	return
}

func GetArticles(pageNum int, pageSize int, maps interface {}) (articles []Article) {
	Eloquent.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)

	return
}

func GetArticle(id int) (article Article) {
	Eloquent.Where("id = ?", id).First(&article)
	Eloquent.Model(&article).Related(&article.Tag)

	return
}

func EditArticle(id int, data interface {}) bool {
	Eloquent.Model(&Article{}).Where("id = ?", id).Updates(data)

	return true
}

func AddArticle(data map[string]interface {}) bool {
	Eloquent.Create(&Article {
		TagID : data["tag_id"].(int),
		Title : data["title"].(string),
		Desc : data["desc"].(string),
		Content : data["content"].(string),
		CreatedBy : data["created_by"].(string),
		State : data["state"].(int),
		CoverImageUrl: data["cover_image_url"].(string),
	})

	return true
}

func DeleteArticle(id int) bool {
	Eloquent.Where("id = ?", id).Delete(Article{})

	return true
}

func CleanAllArticle() bool {
	Eloquent.Unscoped().Where("deleted_on != ? ", 0).Delete(&Article{})

	return true
}