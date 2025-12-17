package models

type DomainModel struct {
	BaseModel
	Name   string
	Domain string
}

func (DomainModel) TableName() string {
	return "domains"
}
