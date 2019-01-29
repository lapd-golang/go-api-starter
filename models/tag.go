package models

type Tag struct {
	Base
	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func (t *Tag) Get(pageNum int, pageSize int) (tags []Tag) {
	db.Where(t).Offset(pageNum).Limit(pageSize).Find(&tags)

	return
}

func (t *Tag) GetTotal() (count int) {
	db.Model(&Tag{}).Where(t).Count(&count)

	return
}

func (t *Tag) ExistByName(name string) bool {
	var tag Tag
	db.Select("id").Where("name = ?", name).First(&tag)
	if tag.ID > 0 {
		return true
	}

	return false
}

func (t *Tag) Insert() (id int, err error) {
	result := db.Create(&t)
	id = t.ID
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}

func (t *Tag) ExistByID(id int) bool {
	var tag Tag
	db.Select("id").Where("id = ?", id).First(&tag)
	if tag.ID > 0 {
		return true
	}

	return false
}

func (t *Tag) Delete(id int) (tag Tag, err error) {
	if err = db.Select([]string{"id"}).First(&t, id).Error; err != nil {
		return
	}

	if err = db.Delete(&t).Error; err != nil {
		return
	}
	tag = *t
	return
}

func (t *Tag) Update(id int) (updateTag Tag, err error) {
	if err = db.Select([]string{"id"}).First(&updateTag, id).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = db.Model(&updateTag).Updates(&t).Error; err != nil {
		return
	}
	return
}

func (t *Tag) CleanAll() bool {
	db.Unscoped().Where("deleted_on != ? ", 0).Delete(&Tag{})

	return true
}
