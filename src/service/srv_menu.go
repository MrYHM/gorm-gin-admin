package service

import (
	"github.com/mitchellh/mapstructure"
	"github.com/olongfen/gorm-gin-admin/src/models"
	"github.com/olongfen/gorm-gin-admin/src/utils"
)

// AddMenu
func AddMenu(forms []*utils.FormAddMenu) (ret []*models.Menu, err error) {
	db := models.DB.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
			return
		}
		db.Commit()
	}()
	for _, v := range forms {
		var m = new(models.Menu)
		if v.ParentId != 0 {
			m1 := &models.Menu{}
			if err = m1.Get(int(v.ParentId)); err != nil {
				return
			}
		}
		m.ParentId = uint(v.ParentId)
		m.Name = v.Name
		m.Path = v.Path
		m.Component = v.Component
		m.Sort = v.Sort
		mapstructure.Decode(v.Meta, &m.Meta)
		if err = m.Insert(db); err != nil {
			return
		}
		ret = append(ret, m)
	}
	return
}

// GetMenu
func GetMenu(id int) (ret *models.Menu, err error) {
	var data = new(models.Menu)
	if err = data.Get(id); err != nil {
		return
	}
	ret = data
	return
}

// GetMenuList
func GetMenuList() (ret []*models.Menu, err error) {
	if ret, err = models.GetMenuList(); err != nil {
		return
	}
	for _, v := range ret {
		_ = v.Get(int(v.ID))
	}
	return
}

// DelMenu
func DelMenu(id int) (err error) {
	var data = new(models.Menu)
	return data.Delete(id)
}

func UpdateMenu(f *utils.FormUpdateMenu) (ret *models.Menu, err error) {
	data := new(models.Menu)
	if err = data.Update(f.Id, f.ToMap()); err != nil {
		return
	}
	ret = data
	return
}
