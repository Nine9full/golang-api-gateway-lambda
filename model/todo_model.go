package model

type Todo struct {
	ID        int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"` // Identity column
	Name      string `gorm:"column:name" json:"name"`
	Completed bool   `gorm:"column:active" json:"completed"`
}

func TableName() string {
	return "todos"
}
