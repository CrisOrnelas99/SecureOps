package model

import "time"

type WafEvent struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Method    string    `gorm:"not null" json:"method"`
	Path      string    `gorm:"not null" json:"path"`
	Reason    string    `gorm:"not null" json:"reason"`
	SourceIP  string    `gorm:"column:source_ip" json:"sourceIp"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (WafEvent) TableName() string {
	return "waf_events"
}
