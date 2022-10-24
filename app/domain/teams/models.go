package teams

type Team struct {
	ID          int64  `gorm:"AUTO_INCREMENT;column:id;primary_key"`
	Name        string `gorm:"column:name;type:varchar;not null;unique"`
	Identifier  string `gorm:"column:identifier;type:uuid; not null"`
	Description string `gorm:"column:description;type:varchar;not null"`
}

func (Team) TableName() string {
	return "team_tb"
}
