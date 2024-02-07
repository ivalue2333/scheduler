// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package rdm

import (
	"time"
)

const TableNameSchedulerInstance = "scheduler_instance"

// SchedulerInstance 调度实例表
type SchedulerInstance struct {
	ID          int64     `gorm:"column:id;primaryKey;comment:主键" json:"id"`                                                 // 主键
	HeartBeatAt time.Time `gorm:"column:heart_beat_at;not null;default:CURRENT_TIMESTAMP;comment:心跳时间" json:"heart_beat_at"` // 心跳时间
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`       // 创建时间
}

// TableName SchedulerInstance's table name
func (*SchedulerInstance) TableName() string {
	return TableNameSchedulerInstance
}