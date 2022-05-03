package domain

type Authentication struct {
	Id      string `gorm:"index:idx_name,unique"`
	Name    string `gorm:"index:idx_name,unique"`
	Pasword string `gorm:"index:idx_name,unique"`
	Role    string `gorm:"index:idx_name,unique"` // TODO S: ispraviti
}
