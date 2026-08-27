package schema

type HoldReason struct {
	ID     int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Status string  `json:"status" gorm:"type:enum('Active','Inactive') NOT NULL;default:'Active';column:status"`
	Reason *string `json:"reason" gorm:"type:varchar(255);column:reason"`
}

func (HoldReason) TableName() string {
	return "tbl_hold_reason"
}
