package models

import (
	"encoding/json"
	"fmt"

	"github.com/LuukBlankenstijn/slugshrink/internal/authconfig"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AuthConfigModel struct {
	BaseModel
	Type        authconfig.AuthStrategyType
	RawStrategy datatypes.JSON
	Strategy    authconfig.AuthStrategy `gorm:"-"`
}

func (AuthConfigModel) TableName() string {
	return "auth_config"
}

func (m *AuthConfigModel) BeforeSave(tx *gorm.DB) error {
	if m.Strategy != nil {
		m.Type = m.Strategy.Type()

		b, err := json.Marshal(m.Strategy)
		if err != nil {
			return err
		}
		m.RawStrategy = datatypes.JSON(b)
	}
	return nil
}

func (m *AuthConfigModel) AfterFind(tx *gorm.DB) error {
	factory := authconfig.AuthStrategyFactories[m.Type]
	if factory == nil {
		return fmt.Errorf("unknown auth strategy kind: %q", m.Type)
	}

	s := factory()
	if len(m.RawStrategy) > 0 {
		if err := json.Unmarshal(m.RawStrategy, s); err != nil {
			return err
		}
	}
	m.Strategy = s
	return nil
}
