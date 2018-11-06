package models

import "admin-server/database"

type Tag struct {
	Base

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func (t *Tag) Get(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
	database.Eloquent.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)

	return
}

func (t *Tag) GetTotal(maps interface{}) (count int) {
	database.Eloquent.Model(&Tag{}).Where(maps).Count(&count)

	return
}

func (t *Tag) ExistByName(name string) bool {
	var tag Tag
	database.Eloquent.Select("id").Where("name = ?", name).First(&tag)
	if tag.ID > 0 {
		return true
	}

	return false
}

func (t *Tag) Add(name string, state int, createdBy string) int {
	tag := Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}

	database.Eloquent.Create(&tag)

	return tag.ID
}

func (t *Tag) ExistByID(id int) bool {
	var tag Tag
	database.Eloquent.Select("id").Where("id = ?", id).First(&tag)
	if tag.ID > 0 {
		return true
	}

	return false
}

func (t *Tag) Delete(id int) bool {
	database.Eloquent.Where("id = ?", id).Delete(&Tag{})

	return true
}

func (t *Tag) Edit(id int, data interface{}) bool {
	database.Eloquent.Model(&Tag{}).Where("id = ?", id).Updates(data)

	return true
}

func (t *Tag) CleanAll() bool {
	database.Eloquent.Unscoped().Where("deleted_on != ? ", 0).Delete(&Tag{})

	return true
}
