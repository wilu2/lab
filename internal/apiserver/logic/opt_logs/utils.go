package opt_logs

import "time"

type optLogsJoinUsers struct {
	ID           int32     `gorm:"column:id;type:integer;primaryKey;autoIncrement:true" json:"id"`
	OptTime      time.Time `gorm:"column:opt_time;type:character varying" json:"opt_time"`
	Operation    string    `gorm:"column:operation;type:character varying" json:"operation"`
	Resource     string    `gorm:"column:resource;type:character varying" json:"resource"`
	ResourceType string    `gorm:"column:resource_type;type:character varying" json:"resource_type"`
	UserID       int32     `gorm:"column:user_id;type:integer" json:"user_id"`
	UserName     string    `gorm:"column:user_name;type:character varying" json:"user_name"`
}
