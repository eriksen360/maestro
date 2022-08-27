package storage

import "gorm.io/gorm"

type HostModel struct {
	gorm.Model
	Fqdn   string      `gorm:"type:text"`
	IpAddr string      `gorm:"type:text"`
	PubKey string      `gorm:"type:text"`
	Groups string      `gorm:"type:text"`
	Tasks  []TaskModel `gorm:"foreignKey:HostRefer"`
}

type TaskModel struct {
	gorm.Model
	HostRefer      uint
	Name           string `gorm:"type:text"`
	InstructionSet string `gorm:"type:text"` // TODO: should be InstructionModel
}
