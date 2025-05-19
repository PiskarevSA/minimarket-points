package repo

type Error struct {
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

var (
	ErrNoBalanceFound      = &Error{"no balance found"}
	ErrNoTransactionsFound = &Error{"no transactions found"}
	ErrNotEnoughtBalance   = &Error{"not enought balance"}
)
