package data

type Permissions []string

type PermissionsRepository interface {
	GetAllForUser(userId int64) (Permissions, error)
	AddForUser(userId int64, codes ...string) error
}

func (p *Permissions) Included(code string) bool {
	for _, p := range *p {
		if p == code {
			return true
		}
	}
	return false
}
