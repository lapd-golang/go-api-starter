package models

type Tag struct {
	Base

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func (t *Tag) Get(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
	Eloquent.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)

	return
}

func (t *Tag) GetTotal(maps interface{}) (count int) {
	Eloquent.Model(&Tag{}).Where(maps).Count(&count)

	return
}

func (t *Tag) ExistTagByName(name string) bool {
	var tag Tag
	Eloquent.Select("id").Where("name = ?", name).First(&tag)
	if tag.ID > 0 {
		return true
	}

	return false
}

func (t *Tag) AddTag(name string, state int, createdBy string) bool{
	Eloquent.Create(&Tag {
		Name : name,
		State : state,
		CreatedBy : createdBy,
	})

	return true
}

func (t *Tag) ExistTagByID(id int) bool {
	var tag Tag
	Eloquent.Select("id").Where("id = ?", id).First(&tag)
	if tag.ID > 0 {
		return true
	}

	return false
}

func (t *Tag) DeleteTag(id int) bool {
	Eloquent.Where("id = ?", id).Delete(&Tag{})

	return true
}

func (t *Tag) EditTag(id int, data interface {}) bool {
	Eloquent.Model(&Tag{}).Where("id = ?", id).Updates(data)

	return true
}

func (t *Tag) CleanAllTag() bool {
	Eloquent.Unscoped().Where("deleted_on != ? ", 0).Delete(&Tag{})

	return true
}