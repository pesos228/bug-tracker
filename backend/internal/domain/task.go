package domain

import "time"

type CheckStatus string
type CheckResult string

const (
	NotChecked       CheckStatus = "not_checked"
	Checked          CheckStatus = "checked"
	PartiallyChecked CheckStatus = "partially_checked"
	Failed           CheckStatus = "failed_check"
)

const (
	Success CheckResult = "success"
	Failure CheckResult = "failure"
	Warning CheckResult = "warning"
)

type Task struct {
	BaseModel
	SoftName          string      `gorm:"type:varchar(255)"`
	RequestID         string      `gorm:"type:varchar(255)"`
	Description       string      `gorm:"type:text"`
	UserID            string      `gorm:"type:uuid;not null"`
	TestEnvDateUpdate time.Time   `gorm:"type:date"`
	CheckDate         *time.Time  `gorm:"type:date"`
	CheckStatus       CheckStatus `gorm:"type:varchar(20);not null;default:'not_checked'"`
	CheckResult       CheckResult `gorm:"type:varchar(20)"`
	Comment           string      `gorm:"type:text"`
}
