package models

type Role struct {
	Id     string `json:"role_id"`
	Name   string `json:"role_name"`
	Method string `json:"method"`
	Path   string `json:"path"`
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
