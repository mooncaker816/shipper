package go_micro_srv_user

import (
	"github.com/jinzhu/gorm"
	gouuid "github.com/satori/go.uuid"
)

func (model *User) BeforeCreate(scope *gorm.Scope) error {
	uuid, _ := gouuid.NewV4()
	return scope.SetColumn("Id", uuid.String())
}
