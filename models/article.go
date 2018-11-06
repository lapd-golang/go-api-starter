package models

import "admin-server/database"

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

func (a Article) ExistByID(id int) bool {
	var article Article
	database.Eloquent.Select("id").Where("id = ?", id).First(&article)

	if article.ID > 0 {
		return true
	}

	return false
}

func (a Article) GetTotal(maps interface {}) (count int){
	database.Eloquent.Model(&Article{}).Where(maps).Count(&count)

	return
}

func (a Article) Get(pageNum int, pageSize int, maps interface {}) (articles []Article) {
	database.Eloquent.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)

	return
}

func (a Article) GetById(id int) (article Article) {
	database.Eloquent.Where("id = ?", id).First(&article)
	database.Eloquent.Model(&article).Related(&article.Tag)

	return
}

func (a Article) Edit(id int, data interface {}) bool {
	database.Eloquent.Model(&Article{}).Where("id = ?", id).Updates(data)

	return true
}

func (a Article) Add(data map[string]interface {}) int {
	article := Article {
		TagID : data["tag_id"].(int),
		Title : data["title"].(string),
		Desc : data["desc"].(string),
		Content : data["content"].(string),
		CreatedBy : data["created_by"].(string),
		State : data["state"].(int),
		CoverImageUrl: data["cover_image_url"].(string),
	}

	database.Eloquent.Create(&article)

	return article.ID
}

func (a Article) Delete(id int) bool {
	database.Eloquent.Where("id = ?", id).Delete(Article{})

	return true
}

func (a Article) CleanAll() bool {
	database.Eloquent.Unscoped().Where("deleted_on != ? ", 0).Delete(&Article{})

	return true
}