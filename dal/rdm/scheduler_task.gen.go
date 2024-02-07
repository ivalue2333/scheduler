// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package rdm

import (
	"time"
)

const TableNameSchedulerTask = "scheduler_task"

// SchedulerTask 任务表
type SchedulerTask struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement:true;comment:主键" json:"id"`                               // 主键
	Name      string    `gorm:"column:name;not null;comment:任务名,绑定执行器" json:"name"`                                         // 任务名,绑定执行器
	Namespace string    `gorm:"column:namespace;not null;comment:实例空间" json:"namespace"`                                    // 实例空间
	Category  string    `gorm:"column:category;not null;comment:实例空间分类" json:"category"`                                    // 实例空间分类
	Priority  int32     `gorm:"column:priority;not null;default:100;comment:任务优先级, 1 - 100, 1:数值越小,优先级越高" json:"priority"`  // 任务优先级, 1 - 100, 1:数值越小,优先级越高
	Status    int32     `gorm:"column:status;not null;default:1;comment:状态; 1:等待中,2:执行中,3:暂停,4:停止,5:成功,6:过期" json:"status"` // 状态; 1:等待中,2:执行中,3:暂停,4:停止,5:成功,6:过期
	Extra     string    `gorm:"column:extra;comment:业务数据" json:"extra"`                                                     // 业务数据
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`        // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`        // 更新时间
}

// TableName SchedulerTask's table name
func (*SchedulerTask) TableName() string {
	return TableNameSchedulerTask
}
