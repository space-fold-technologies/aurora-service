package environments

type Entry struct {
	Key   string `gorm:"column:key;type:varchar;not null"`
	Value string `gorm:"column:value;type:varchar;not null"`
}

type EnvVar struct {
	ID     uint   `gorm:"AUTO_INCREMENT;column:id;primary_key"`
	Scope  string `gorm:"column:scope;type:varchar;not null"`
	Target string `gorm:"column:target;type:varchar;not null"`
	Key    string `gorm:"column:key;type:varchar;not null;unique"`
	Value  string `gorm:"column:value;type:varchar;not null"`
}

func (EnvVar) TableName() string {
	return "environment_variable_tb"
}
