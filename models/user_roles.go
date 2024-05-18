package models

type User_role struct {
	ID     uint `gorm:"primary_key" json:"id"`
	UserId uint `json:"user_id"`
	RoleId uint `json:"role_id"`
}

func (User_role) TableName() string {
	return "user_role"
}

func FindRoleByUserId(userId interface{}) User_role {
	var uRole User_role
	OldDB.Where("user_id = ?", userId).First(&uRole)
	return uRole
}

func UpdateRole(user_id int, role_id int) {
	DB.Model(&User_role{}).Where(&User_role{UserId: uint(user_id)}).Update("role_id", role_id)
}

func CreateUserRole(userId uint, roleId uint) {
	uRole := &User_role{
		UserId: userId,
		RoleId: roleId,
	}
	OldDB.Create(uRole)
}
func DeleteRoleByUserId(userId interface{}) {
	OldDB.Where("user_id = ?", userId).Delete(User_role{})
}
