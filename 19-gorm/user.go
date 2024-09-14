package gorm

type User struct {
	ID           string    `gorm:"primaryKey;column:id"`
	Password     string    `gorm:"column:password"`
	Name         Name      `gorm:"embedded"`
	CreatedAt    int64     `gorm:"column:created_at;autoCreateTime:mili"`
	UpdatedAt    int64     `gorm:"column:updated_at;autoCreateTime:mili;autoUpdateTime:mili"`
	Information  string    `gorm:"-"`
	Wallet       Wallet    `gorm:"foreignKey:user_id;references:id"`
	Addresses    []Address `gorm:"foreignKey:user_id;references:id"`
	LikeProducts []Product `gorm:"many2many:user_like_product;foreignKey:id;joinForeignKey:user_id;references:id;joinReferences:product_id"`
}

func (u *User) TableName() string {
	return "users"
}

type Name struct {
	FirstName  string `gorm:"column:first_name"`
	MiddleName string `gorm:"column:middle_name"`
	LastName   string `gorm:"column:last_name"`
}

type UserLog struct {
	ID        int    `gorm:"primaryKey;column:id;autoIncrement"`
	UserId    string `gorm:"column:user_id"`
	Action    string `gorm:"column:action"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime:mili"`
	UpdatedAt int64  `gorm:"column:updated_at;autoCreateTime:mili;autoUpdateTime:mili"`
}
