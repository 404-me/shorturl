package model

type URL struct {
	ID       uint   `gorm:"primaryKey"`
	Code     string `gorm:"uniqueIndex;default:null"`
	Original string
}
