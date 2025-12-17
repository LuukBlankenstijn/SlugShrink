package models

import "github.com/google/uuid"

type RedirectModel struct {
	BaseModel
	DomainId  uuid.UUID
	Domain    DomainModel `gorm:"foreignKey:DomainId"`
	Path      string
	TargetUrl string
	Active    bool
}

func (RedirectModel) TableName() string {
	return "redirects"
}
