package domain

type BaseModel struct {
	ID string `gorm:"type:uuid;primary_key"`
}
