package account

import (
	"runtime"
	"sync/atomic"
)

type Account struct {
	balance int64
	closed  int32
}

// If initial deposit is a negative, it return nil.
func Open(initialDeposit int64) *Account {
	if initialDeposit < 0 {
		return nil
	}

	return &Account{
		balance: initialDeposit,
		closed:  0,
	}
}

// Get account balance, if account closed, always return ok false
// If account closed, it not modify the account and return ok = false.
func (a *Account) Balance() (balance int64, ok bool) {
	if a.isClosed() {
		return
	}
	return atomic.LoadInt64(&a.balance), true
}

// handle a negative amount as a withdrawal.
// Withdrawals not succeed if they result in a negative balance and return ok = false
// If account closed, it not modify the account and return ok = false.
func (a *Account) Deposit(amount int64) (newBalance int64, ok bool) {
	if a.isClosed() {
		return
	}

	for {
		balance := atomic.LoadInt64(&a.balance)
		if balance < -amount {
			break
		}
		if atomic.CompareAndSwapInt64(&a.balance, balance, balance+amount) {
			newBalance = balance + amount
			ok = true
			break
		} else {
			runtime.Gosched() // Effectively in concurrency
		}
	}

	return
}

// Close account, If account closed, it return ok = false
func (a *Account) Close() (payout int64, ok bool) {
	if !a.isClosed() && atomic.CompareAndSwapInt32(&a.closed, 0, 1) {
		ok = true
		payout = atomic.LoadInt64(&a.balance)
		atomic.StoreInt64(&a.balance, 0)
		return
	}
	return
}

// Check account to closed
func (a *Account) isClosed() bool {
	if atomic.LoadInt32(&a.closed) == 1 {
		return true
	}
	return false
}

type Response struct {
	Balance int64 `json:"amount"`
}

type Request struct {
	Balance int64 `json:"initialAmount"`
}
