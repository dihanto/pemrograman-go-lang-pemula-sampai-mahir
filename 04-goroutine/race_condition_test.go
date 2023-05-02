package goroutine

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRaceCondition(t *testing.T) {
	x := 0
	for i := 1; i <= 1000; i++ {
		go func() {
			for j := 1; j <= 100; j++ {
				x = x + 1
			}
		}()
	}

	time.Sleep(5 * time.Second)
	fmt.Println("Counter = ", x)
}
func TestMutex(t *testing.T) {
	x := 0
	var mutex sync.Mutex
	for i := 1; i <= 1000; i++ {
		go func() {
			for j := 1; j <= 1000; j++ {
				mutex.Lock()
				x = x + 1
				mutex.Unlock()
			}
		}()
	}

	time.Sleep(5 * time.Second)
	fmt.Println("Counter = ", x)
}

type BankAccount struct {
	RWMutex sync.RWMutex
	Balance int
}

func (account *BankAccount) AddBalance(amount int) {
	account.RWMutex.Lock()
	account.Balance = account.Balance + amount
	account.RWMutex.Unlock()
}

func (account *BankAccount) GetBalance() int {
	account.RWMutex.RLock()
	balance := account.Balance
	account.RWMutex.RUnlock()
	return balance
}

func TestRWMutex(t *testing.T) {
	account := BankAccount{}

	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				account.AddBalance(1)
				fmt.Println(account.GetBalance())
			}
		}()
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Total Balance", account.GetBalance())
}

type UserBanlace struct {
	sync.Mutex
	Name    string
	Balance int
}

func (user *UserBanlace) Lock() {
	user.Mutex.Lock()
}

func (user *UserBanlace) Unlock() {
	user.Mutex.Unlock()
}

func (user *UserBanlace) Change(amount int) {
	user.Balance = user.Balance + amount
}

func Transfer(user1 *UserBanlace, user2 *UserBanlace, amount int) {
	user1.Lock()
	fmt.Println("Lock user1", user1.Name)
	user1.Change(-amount)

	time.Sleep(1 * time.Second)

	user2.Lock()
	fmt.Println("Lock User2", user2.Name)
	user2.Change(amount)

	user1.Unlock()
	user2.Unlock()

}
func TestDeadlock(t *testing.T) {
	user1 := UserBanlace{
		Name:    "Kurniawan",
		Balance: 100000,
	}
	user2 := UserBanlace{
		Name:    "Kurniadi",
		Balance: 100000,
	}
	go Transfer(&user1, &user2, 100000)
	go Transfer(&user2, &user1, 100000)

	time.Sleep(2 * time.Second)

	fmt.Println("User ", user1.Name, ", Balance ", user1.Balance)
	fmt.Println("User ", user2.Name, ", Balance ", user2.Balance)
}
