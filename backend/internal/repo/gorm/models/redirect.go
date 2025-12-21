package models

import (
	"fmt"

	"github.com/LuukBlankenstijn/gewish/internal/authconfig"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RedirectModel struct {
	BaseModel
	DomainId  uuid.UUID
	Domain    DomainModel `gorm:"foreignKey:DomainId"`
	Path      string
	TargetUrl string
	Active    bool
	Creator   *string `gorm:"->;<-:create"`
}

func (RedirectModel) TableName() string {
	return "redirects"
}

func (r *RedirectModel) BeforeCreate(tx *gorm.DB) error {
	authState, _ := authconfig.GetAuthState(tx.Statement.Context)
	if authState == nil {
		return fmt.Errorf("user not authenticated")
	}
	r.Creator = authState.UserId
	return nil
}
