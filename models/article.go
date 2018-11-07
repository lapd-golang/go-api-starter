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

func (a *Article) ExistByID(id int) bool {
	var article Article
	database.Eloquent.Select("id").Where("id = ?", id).First(&article)

	if article.ID > 0 {
		return true
	}

	return false
}

func (a *Article) GetTotal() (count int){
	database.Eloquent.Model(&Article{}).Where(a).Count(&count)

	return
}

func (a *Article) Get(pageNum int, pageSize int) (articles []Article) {
	database.Eloquent.Preload("Tag").Where(a).Offset(pageNum).Limit(pageSize).Find(&articles)

	return
}

func (a *Article) GetById(id int) (article Article) {
	database.Eloquent.Where("id = ?", id).First(&article)
	database.Eloquent.Model(&article).Related(&article.Tag)

	return
}

func (a *Article) Update(id int) (updateArticle Article, err error) {
	if err = database.Eloquent.Select([]string{"id"}).First(&updateArticle, id).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = database.Eloquent.Model(&updateArticle).Updates(&a).Error; err != nil {
		return
	}
	return
}

func (a *Article) Insert() (id int, err error) {
	result := database.Eloquent.Create(&a)
	id = a.ID
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}

func (a *Article) Delete(id int) (article Article, err error) {
	if err = database.Eloquent.Select([]string{"id"}).First(&a, id).Error; err != nil {
		return
	}

	if err = database.Eloquent.Delete(&a).Error; err != nil {
		return
	}
	article = *a
	return
}

func (a *Article) CleanAll() bool {
	database.Eloquent.Unscoped().Where("deleted_on != ? ", 0).Delete(&Article{})

	return true
}