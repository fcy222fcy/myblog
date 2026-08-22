package entity

import "time"

// ContentViewEvent records one daily unique view without storing raw visitor identity.
type ContentViewEvent struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	ContentType string    `gorm:"type:varchar(32);not null;uniqueIndex:uk_content_view,priority:1;index:idx_content_view_date,priority:1" json:"content_type"`
	ContentID   uint      `gorm:"not null;uniqueIndex:uk_content_view,priority:2" json:"content_id"`
	VisitorKey  string    `gorm:"type:char(64);not null;uniqueIndex:uk_content_view,priority:3" json:"-"`
	ViewDate    time.Time `gorm:"type:date;not null;uniqueIndex:uk_content_view,priority:4;index:idx_content_view_date,priority:2" json:"view_date"`
	CreatedAt   time.Time `json:"created_at"`
}

func (ContentViewEvent) TableName() string {
	return "content_view_events"
}
