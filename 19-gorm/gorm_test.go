package gorm

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

func OpenConnection() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/belajar_golang_gorm"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetConnMaxIdleTime(5 * time.Minute)
	sqlDb.SetConnMaxLifetime(60 * time.Minute)
	return db
}

var db = OpenConnection()

func TestExecuteSQL(t *testing.T) {
	err := db.Exec("insert into sample(id, name) values(?, ?)", 1, "dihanto").Error
	assert.Nil(t, err)

	err = db.Exec("insert into sample(id, name) values(?, ?)", 2, "budi").Error
	assert.Nil(t, err)

	err = db.Exec("insert into sample(id, name) values(?, ?)", 3, "joko").Error
	assert.Nil(t, err)

	err = db.Exec("insert into sample(id, name) values(?, ?)", 4, "siti").Error
	assert.Nil(t, err)
}

type Sample struct {
	Id   string
	Name string
}

func TestRawSQL(t *testing.T) {
	var sample Sample

	err := db.Raw("select id, name from sample where id = ?", 1).Scan(&sample).Error
	assert.Nil(t, err)
	assert.Equal(t, "dihanto", sample.Name)

	var samples []Sample

	err = db.Raw("select id, name from sample").Scan(&samples).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(samples))
}

func TestSQLRows(t *testing.T) {

	rows, err := db.Raw("select id, name from sample").Rows()
	assert.Nil(t, err)
	defer rows.Close()
	var samples []Sample
	for rows.Next() {
		var sample Sample
		err := rows.Scan(&sample.Id, &sample.Name)
		assert.Nil(t, err)
		samples = append(samples, sample)
	}

	assert.Equal(t, 4, len(samples))
}

func TestScanRows(t *testing.T) {
	rows, err := db.Raw("select id, name from sample").Rows()
	assert.Nil(t, err)
	defer rows.Close()
	var samples []Sample
	for rows.Next() {
		err := db.ScanRows(rows, &samples)
		assert.Nil(t, err)
	}

	assert.Equal(t, 4, len(samples))
}

func TestCrateUser(t *testing.T) {
	user := User{
		ID:       "1",
		Password: "12345",
		Name: Name{
			FirstName:  "Dihantoo",
			MiddleName: "To",
			LastName:   "P",
		},
		Information: "Information",
	}
	err := db.Create(&user).Error
	assert.Nil(t, err)

}

func TestBatchInster(t *testing.T) {
	var users []User

	for i := 2; i < 11; i++ {
		users = append(users, User{
			ID:       strconv.Itoa(i),
			Password: "123456",
			Name: Name{
				FirstName:  "User",
				MiddleName: "",
				LastName:   strconv.Itoa(i),
			},
			Information: "Information",
		})
	}
	response := db.Create(&users)

	assert.Nil(t, response.Error)
	assert.Equal(t, int64(9), response.RowsAffected)

}

func TestTransactionSuccess(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&User{
			ID:       "101",
			Password: "12345",
			Name: Name{
				FirstName:  "Dihanto",
				MiddleName: "To",
				LastName:   "P",
			},
			Information: "Information",
		}).Error

		if err != nil {
			return err
		}
		err = tx.Create(&User{
			ID:       "102",
			Password: "12345",
			Name: Name{
				FirstName:  "Dihanto",
				MiddleName: "To",
				LastName:   "P",
			},
			Information: "Information",
		}).Error

		if err != nil {
			return err
		}

		return nil
	})

	assert.Nil(t, err)
}

func TestTransactionError(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&User{
			ID:       "103",
			Password: "12345",
			Name: Name{
				FirstName:  "Dihanto",
				MiddleName: "To",
				LastName:   "P",
			},
			Information: "Information",
		}).Error

		if err != nil {
			return err
		}

		err = tx.Create(&User{
			ID:       "102",
			Password: "12345",
			Name: Name{
				FirstName:  "Dihanto",
				MiddleName: "To",
				LastName:   "P",
			},
			Information: "Information",
		}).Error

		if err != nil {
			return err
		}

		return nil
	})

	assert.NotNil(t, err)
}

func TestManualTransactionSuccess(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	err := tx.Create(&User{
		ID:       "103",
		Password: "12345",
		Name: Name{
			FirstName:  "Dihanto",
			MiddleName: "To",
			LastName:   "P",
		},
		Information: "Information",
	}).Error
	assert.Nil(t, err)

	err = tx.Create(&User{
		ID:       "104",
		Password: "12345",
		Name: Name{
			FirstName:  "Dihanto",
			MiddleName: "To",
			LastName:   "P",
		},
		Information: "Information",
	}).Error
	assert.Nil(t, err)
	if err == nil {
		tx.Commit()
	}
}

func TestManualTransactionError(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	err := tx.Create(&User{
		ID:       "105",
		Password: "12345",
		Name: Name{
			FirstName:  "Dihanto",
			MiddleName: "To",
			LastName:   "P",
		},
		Information: "Information",
	}).Error
	assert.Nil(t, err)

	err = tx.Create(&User{
		ID:       "104",
		Password: "12345",
		Name: Name{
			FirstName:  "Dihanto",
			MiddleName: "To",
			LastName:   "P",
		},
		Information: "Information",
	}).Error
	assert.NotNil(t, err)
	if err == nil {
		tx.Commit()
	}
}

func TestQuerySingleObject(t *testing.T) {
	user := User{}
	result := db.First(&user)
	assert.Nil(t, result.Error)
	assert.Equal(t, "10", user.ID)

	user = User{}
	result = db.Last(&user)
	assert.Nil(t, result.Error)
	assert.Equal(t, "99", user.ID)
}

func TestQuerySingleObjectInlinceCondition(t *testing.T) {

	var users []User
	result := db.First(&users, "id = ?", "10")
	assert.Nil(t, result.Error)
	assert.Equal(t, 1, len(users))

	result = db.Take(&users, "id = ?", "99")
	assert.Nil(t, result.Error)
	assert.Equal(t, 1, len(users))
}

func TestQueryAllObject(t *testing.T) {
	var users []User
	result := db.Find(&users, "id in ?", []string{"10", "101", "102"})
	assert.Nil(t, result.Error)
	assert.Equal(t, 3, len(users))
}

func TestQueryCondition(t *testing.T) {

	var users []User
	result := db.Where("first_name = ?", "Dihanto").
		Where("id = ?", "10").Find(&users)

	assert.Nil(t, result.Error)
	assert.Equal(t, 1, len(users))
}

func TestOrOperator(t *testing.T) {
	var users []User
	result := db.Where("first_name = ?", "Dihanto").
		Or("password = ?", "12345").Find(&users)

	assert.Nil(t, result.Error)
	assert.Equal(t, 102, len(users))
}

func TestNotOperator(t *testing.T) {
	var users []User
	result := db.Not("id like ?", "1%").Find(&users)
	assert.Nil(t, result.Error)
	assert.Equal(t, 88, len(users))
}

func TestSelectFields(t *testing.T) {

	var users []User

	result := db.Select("id, first_name").Find(&users)
	assert.Nil(t, result.Error)
	for _, user := range users {
		assert.NotNil(t, user.ID)
		assert.NotEqual(t, user.ID, "")
		assert.NotNil(t, user.Name.FirstName)
	}
}

func TestStructCondition(t *testing.T) {
	UserCondition := User{
		Name: Name{
			FirstName: "Dihanto",
		},
	}
	var users []User

	result := db.Where(&UserCondition).Find(&users)
	assert.Nil(t, result.Error)
	assert.Equal(t, 102, len(users))
}

func TestMapCondition(t *testing.T) {
	mapCondition := map[string]interface{}{
		"middle_name": "",
	}

	var users []User

	result := db.Where(mapCondition).Find(&users)

	assert.Nil(t, result.Error)
	assert.Equal(t, 0, len(users))
}

func TestOrderLimitOffset(t *testing.T) {

	var users []User
	err := db.Order("id asc, first_name desc").Limit(2).Offset(1).Find(&users).Error
	assert.Nil(t, err)

	assert.Equal(t, 2, len(users))
}

type UserResponse struct {
	ID        string
	FirstName string
	LastName  string
}

func TestQueryNonModel(t *testing.T) {
	var users []UserResponse

	result := db.Model(&User{}).Select("id, first_name, last_name").Find(&users)
	assert.Nil(t, result.Error)

	assert.Equal(t, 102, len(users))
}

func TestUpdate(t *testing.T) {

	var user User

	result := db.First(&user, "id = ?", "101")

	assert.Nil(t, result.Error)
	user.Name.FirstName = "budi"
	user.Name.MiddleName = "to"
	user.Name.LastName = "nugraha"
	user.Password = "1234567"

	result = db.Save(&user)
	assert.Nil(t, result.Error)
}

func TestSelectedColumn(t *testing.T) {
	result := db.Model(&User{}).Where("id = ?", "102").Updates(map[string]interface{}{
		"first_name": "reno",
		"last_name":  "saputra",
		"password":   "1234567",
	})

	assert.Nil(t, result.Error)

	result = db.Model(&User{}).Where("id = ?", "102").Update("password", "123456789")

	assert.Nil(t, result.Error)
	err := db.Where("id = ?", "102").Updates(&User{
		Name: Name{
			MiddleName: "rifki",
		},
	}).Error

	assert.Nil(t, err)
}

func TestAutoIncrement(t *testing.T) {
	for i := 0; i < 10; i++ {
		userLog := UserLog{
			UserId: "101",
			Action: "create",
		}
		result := db.Create(&userLog)
		assert.Nil(t, result.Error)

		assert.NotEqual(t, 0, userLog.ID)
	}

}

func TestSaveOrUpdate(t *testing.T) {
	userLog := UserLog{
		UserId: "101",
		Action: "create",
	}
	result := db.Save(&userLog)
	assert.Nil(t, result.Error)

	userLog.UserId = "102"
	result = db.Save(&userLog)
	assert.Nil(t, result.Error)
}

func TestSaveOrUpdateNotAutoIncrement(t *testing.T) {
	user := User{
		ID: "99",
		Name: Name{
			FirstName: "User 99",
		},
	}

	result := db.Save(&user)
	assert.Nil(t, result.Error)

	user.Name.FirstName = "User 99 Updated"
	result = db.Save(&user)
	assert.Nil(t, result.Error)
}

func TestConflict(t *testing.T) {

	user := User{
		ID: "88",
		Name: Name{
			FirstName: "User 88",
		},
	}
	err := db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&user).Error
	assert.Nil(t, err)
}

func TestDelete(t *testing.T) {
	var user User

	result := db.Take(&user, "id = ?", "88")
	assert.Nil(t, result.Error)

	result = db.Delete(&user)
	assert.Nil(t, result.Error)

	result = db.Delete(&user, "id = ?", "99")
	assert.Nil(t, result.Error)

	err := db.Where("id = ?", "10").Delete(&user).Error
	assert.Nil(t, err)

}

func TestSoftDelete(t *testing.T) {
	todo := Todo{
		UserId:      "1",
		Title:       "Todo 1",
		Description: "Todo 1 Description",
	}
	result := db.Create(&todo)
	assert.Nil(t, result.Error)

	result = db.Delete(&todo)
	assert.Nil(t, result.Error)
	assert.NotNil(t, todo.DeletedAt)

	var todos []Todo
	result = db.Find(&todos)
	assert.Nil(t, result.Error)
	assert.Equal(t, 0, len(todos))
}

func TestUnscoped(t *testing.T) {
	var todo Todo
	result := db.Unscoped().Where("id = ?", 3).First(&todo)
	assert.Nil(t, result.Error)
	assert.NotNil(t, todo.DeletedAt)
	err := db.Unscoped().Delete(&todo).Error
	assert.Nil(t, err)

	var todos []Todo
	result = db.Unscoped().Find(&todos)
	assert.Nil(t, result.Error)
	assert.Equal(t, 0, len(todos))
}

func TestLock(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		var user User
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&user, "id = ?", "1").Error
		if err != nil {
			return err
		}
		user.Name.FirstName = "Dihantoo"
		user.Name.LastName = "Saputra"
		err = tx.Save(&user).Error
		if err != nil {
			return err
		}
		return nil
	})

	assert.Nil(t, err)
}

func TestCreateWaller(t *testing.T) {
	wallet := Wallet{
		ID:      "1",
		UserId:  "1",
		Balance: 1000,
	}
	result := db.Create(&wallet)
	assert.Nil(t, result.Error)
}

func TestRetrieveRelation(t *testing.T) {
	var user User
	err := db.Model(&User{}).Preload("Wallet").Take(&user, "id = ?", "1").Error
	assert.Nil(t, err)

	assert.Equal(t, int64(1000), user.Wallet.Balance)

}

func TestRetrieveRelationJoin(t *testing.T) {
	var user User
	err := db.Model(&User{}).Joins("Wallet").Take(&user, "users.id = ?", "1").Error
	assert.Nil(t, err)

	assert.Equal(t, int64(1000), user.Wallet.Balance)

}

func TestAutoCreateUpdate(t *testing.T) {
	user := User{
		ID:       "20",
		Password: "rahasia",
		Name: Name{
			FirstName: "User 20",
		},
		Wallet: Wallet{
			ID:      "20",
			UserId:  "20",
			Balance: 1000,
		},
	}
	result := db.Create(&user)
	assert.Nil(t, result.Error)
}

func TestSkipAutoCreateUpdate(t *testing.T) {
	user := User{
		ID:       "21",
		Password: "rahasia",
		Name: Name{
			FirstName: "User 21",
		},
		Wallet: Wallet{
			ID:      "21",
			UserId:  "21",
			Balance: 1000,
		},
	}
	result := db.Omit(clause.Associations).Create(&user)
	assert.Nil(t, result.Error)
}

func TestUserAndAddresses(t *testing.T) {
	user := User{
		ID:       "23",
		Password: "rahasia",
		Name: Name{
			FirstName: "User 23",
		},
		Wallet: Wallet{
			ID:      "23",
			UserId:  "23",
			Balance: 1000,
		},
		Addresses: []Address{
			{
				ID:      "23",
				UserId:  "23",
				Address: "Jl. User 23 no 1",
			},
			{
				ID:      "23",
				UserId:  "23",
				Address: "Jl. User 23 no 2",
			},
		},
	}

	result := db.Create(&user)
	assert.Nil(t, result.Error)
}

func TestPreloadJoinOneToMany(t *testing.T) {
	var usersPreload []User

	err := db.Model(&User{}).Joins("Wallet").Preload("Addresses").Find(&usersPreload).Error
	assert.Nil(t, err)
}

func TestTakePreloadJoinOneToMany(t *testing.T) {
	var usersPreload User

	err := db.Model(&User{}).Joins("Wallet").Preload("Addresses").Take(&usersPreload, "users.id = ?", "22").Error
	assert.Nil(t, err)
}

func TestBelongTo(t *testing.T) {
	fmt.Println("preload")
	var addresses []Address
	err := db.Model(&Address{}).Preload("User").Find(&addresses).Error
	assert.Nil(t, err)
	assert.Equal(t, 2, len(addresses))

	fmt.Println("joins")
	var addressesJoin []Address
	err = db.Model(&Address{}).Joins("User").Find(&addressesJoin).Error
	assert.Nil(t, err)
	assert.Equal(t, 2, len(addressesJoin))

}

func TestBelongToWallet(t *testing.T) {
	fmt.Println("preload")
	var wallets []Wallet
	err := db.Model(&Wallet{}).Preload("User").Find(&wallets).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(wallets))

	fmt.Println("joins")
	var walletsJoin []Wallet
	err = db.Model(&Wallet{}).Joins("User").Find(&walletsJoin).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(walletsJoin))

}

func TestCreateManyToMany(t *testing.T) {
	product := Product{
		ID:    "P001",
		Name:  "product 1",
		Price: 1000,
	}
	result := db.Create(&product)
	assert.Nil(t, result.Error)

	err := db.Table("user_like_product").Create(map[string]interface{}{
		"user_id":    "1",
		"product_id": "P001",
	}).Error
	assert.Nil(t, err)

	err = db.Table("user_like_product").Create(map[string]interface{}{
		"user_id":    "2",
		"product_id": "P001",
	}).Error
	assert.Nil(t, err)
}

func TestPreloadManyToMany(t *testing.T) {
	var product Product
	err := db.Preload("LikedByUsers").Find(&product, "id = ?", "P001").Error
	assert.Nil(t, err)
	assert.Equal(t, 2, len(product.LikedByUsers))
}

func TestPreloadManyToManyUser(t *testing.T) {
	var user User
	err := db.Preload("LikeProducts").Find(&user, "id = ?", "1").Error
	assert.Nil(t, err)
	assert.Equal(t, 1, len(user.LikeProducts))
}

func TestAssociationFind(t *testing.T) {
	var product Product
	err := db.Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)

	var users []User
	err = db.Model(&product).Where("users.first_name LIKE ?", "User%").Association("LikedByUsers").Find(&users)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}

func TestAssociationAdd(t *testing.T) {
	var user User
	err := db.Take(&user, "id = ?", "3").Error
	assert.Nil(t, err)
	var product Product
	err = db.Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)
	err = db.Model(&product).Association("LikedByUsers").Append(&user)
	assert.Nil(t, err)
}

func TestAssociationReplace(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		var user User
		err := tx.Take(&user, "id = ?", "1").Error
		assert.Nil(t, err)

		wallet := Wallet{
			ID:      "W001",
			Balance: 1000,
			UserId:  user.ID,
		}

		err = tx.Model(&user).Association("Wallet").Replace(&wallet)
		assert.Nil(t, err)
		return nil
	})

	assert.Nil(t, err)
}

func TestAssociationDelete(t *testing.T) {
	var user User
	err := db.Take(&user, "id = ?", "3").Error
	assert.Nil(t, err)
	var product Product
	err = db.Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)
	err = db.Model(&product).Association("LikedByUsers").Delete(&user)
	assert.Nil(t, err)
}

func TestAssociationClear(t *testing.T) {
	var product Product
	err := db.Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)
	err = db.Model(&product).Association("LikedByUsers").Clear()
	assert.Nil(t, err)
}

func TestPreloadingWithCondition(t *testing.T) {
	var user User

	err := db.Preload("Wallet", "balance > ?", 100).Find(&user, "id = ?", "1").Error
	assert.Nil(t, err)
	fmt.Println(user)
}

func TestNestedPreloading(t *testing.T) {
	var wallet Wallet
	err := db.Preload("User.Addresses").Take(&wallet, "id = ?", "22").Error
	assert.Nil(t, err)
	fmt.Println(wallet)
	fmt.Println(wallet.User)
	fmt.Println(wallet.User.Addresses)
}

func TestPreloadAll(t *testing.T) {
	var user User
	err := db.Preload(clause.Associations).Find(&user, "id = ?", "22").Error
	assert.Nil(t, err)
}

func TestJoinQuery(t *testing.T) {
	var users []User

	err := db.Joins("join wallets on wallets.user_id = users.id").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 5, len(users))

	users = []User{}
	err = db.Joins("Wallet").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 16, len(users))

}

func TestJoinWithCondition(t *testing.T) {
	var users []User

	err := db.Joins("join wallets on wallets.user_id = users.id and wallets.balance > ?", 100).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 5, len(users))

	users = []User{}
	err = db.Joins("Wallet").Where("Wallet.balance > ?", 100).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 5, len(users))
}

func TestCount(t *testing.T) {
	var count int64
	err := db.Model(&User{}).Joins("Wallet").Where("Wallet.balance > ?", 100).Count(&count).Error
	assert.Nil(t, err)
	assert.Equal(t, int64(5), count)
}

type AggregationResult struct {
	TotalBalance   int64
	MinBalance     int64
	MaxBalance     int64
	AverageBalance float64
}

func TestAggreagtion(t *testing.T) {
	var result AggregationResult

	err := db.Model(&Wallet{}).Select("sum(balance) as total_balance, min(balance) as min_balance, max(balance) as max_balance, avg(balance) as average_balance").Take(&result).Error
	assert.Nil(t, err)
	assert.Equal(t, int64(5000), result.TotalBalance)
	assert.Equal(t, int64(1000), result.MinBalance)
	assert.Equal(t, int64(1000), result.MaxBalance)
	assert.Equal(t, float64(1000), result.AverageBalance)
}
func TestAggreagtionGroupByAndHaving(t *testing.T) {
	var result []AggregationResult

	err := db.Model(&Wallet{}).Select("sum(balance) as total_balance, min(balance) as min_balance, max(balance) as max_balance, avg(balance) as average_balance").
		Joins("User").Group("User.id").Having("sum(balance) > ?", 100).Find(&result).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(result))
}

func TestContext(t *testing.T) {
	ctx := context.Background()

	var users []User
	err := db.WithContext(ctx).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 15, len(users))
}

func BrokeWalletBalance(db *gorm.DB) *gorm.DB {
	return db.Where("balance = ?", 0)
}

func SultanWalletBalance(db *gorm.DB) *gorm.DB {
	return db.Where("balance > ?", 1000)
}

func TestScopes(t *testing.T) {
	var wallets []Wallet

	err := db.Scopes(BrokeWalletBalance).Find(&wallets).Error
	assert.Nil(t, err)
	assert.Equal(t, 0, len(wallets))

	wallets = []Wallet{}
	err = db.Scopes(SultanWalletBalance).Find(&wallets).Error
	assert.Nil(t, err)
	assert.Equal(t, 0, len(wallets))
}
