package models

import (
	"fmt"

	"github.com/LuukBlankenstijn/slugshrink/internal/authconfig"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RedirectModel struct {
	BaseModel
	DomainId  uuid.UUID   `gorm:"uniqueIndex:idx_unique_target"`
	Domain    DomainModel `gorm:"foreignKey:DomainId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Path      string      `gorm:"uniqueIndex:idx_unique_target"`
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
