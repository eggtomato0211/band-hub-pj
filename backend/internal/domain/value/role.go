package value

import "fmt"

type Role string

// 有効なロールの定数定義
const (
	RoleMember Role = "member"
	RoleAdmin  Role = "admin"
)

var validRoles = map[Role]bool{
	RoleMember: true,
	RoleAdmin:  true,
}

func NewRole(s string) (Role, error) {
	r := Role(s)
	if validRoles[r] {
		return r, nil
	}
	return "", fmt.Errorf("無効なロールです: %s", s)
}

func (r Role) IsAdmin() bool {
	return r == RoleAdmin
}