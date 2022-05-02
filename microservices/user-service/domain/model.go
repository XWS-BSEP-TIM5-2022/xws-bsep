package domain

type User struct {
	// Id   int `bson:"_id"`
	// Name string             `bson:"name"`
	// Id   int    `json:"id"`
	// Name string `json:"name"`
	Id   string `gorm:"index:idx_name,unique"`
	Name string `gorm:"index:idx_name,unique"`
}
