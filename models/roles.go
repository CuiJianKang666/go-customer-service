package models

import "fmt"

type Role struct {
	Id     int    `json:"role_id"`
	Name   string `json:"role_name"`
	Method string `json:"method"`
	Path   string `json:"path"`
}

func (Role) TableName() string {
	return "role"
}

func FindRoles() []Role {
	var roles []Role
	OldDB.Order("id desc").Find(&roles)
	return roles
}
func FindRole(id interface{}) Role {
	var role Role
	OldDB.Where("id = ?", id).First(&role)
	return role
}
func SaveRole(id string, name string, method string, path string) {
	role := &Role{
		Method: method,
		Name:   name,
		Path:   path,
	}
	OldDB.Model(role).Where("id=?", id).Update(role)
}

func CreateRole(name string, method string, path string) {
	role := Role{
		Method: method,
		Name:   name,
		Path:   path,
	}
	err := DB.Model(&Role{}).Create(&role).Error
	fmt.Println(err)
}

func DeleteRole(id int) {
	role := Role{
		Id: id,
	}
	DB.Model(&Role{}).Delete(&role)
}
