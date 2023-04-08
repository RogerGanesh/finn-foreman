package data

import (
	"context"
	"database/sql"
	"errors"

	// "errors"

	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		RecurringPayment: RecurringPayment{},
		PaymentHistory:   PaymentHistory{},
	}
}

type Models struct {
	RecurringPayment RecurringPayment
	PaymentHistory   PaymentHistory
}

type RecurringPayment struct {
	PaymentID          int     `json:"paymentID"`
	UserName           string  `json:"username"`
	AccountName        string  `json:"accountname"`
	PaymentAmount      float32 `json:"amount"`
	PaymentName        string  `json:"paymentName"`
	PaymentDescription string  `json:"paymentDescription"`
	PaymentDate        string  `json:"paymentDate"`
	PaymentType        string  `json:"paymentType"`
}

type PaymentHistory struct {
	PaymentHistoryID     int    `json:"-"`
	PaymentID            int    `json:"paymentID"`
	PaymentHistoryDate   string `json:"paymentHistoryDate"`
	PaymentHistoryStatus bool   `json:"paymentHistoryStatus"`
}

func (t *RecurringPayment) GetUserBalance(email string, account string) (float32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select SUM(TransactionAmount) From mrkrabs.Transactions
	where Username = $1 and accountname = $2
	group by Username`

	var totalBalance float32

	row := db.QueryRowContext(ctx, query, email, account)
	err := row.Scan(&totalBalance)

	if err != nil {
		return 0, err
	}

	return totalBalance, nil
}

func (t *RecurringPayment) UpdateBalance(username string, account string, transactionAmount float32, transactionName string, transactionDescription string, transactionCategory string) (float32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `insert into mrkrabs.Transactions (Username, AccountName, TransactionAmount, TransactionName, TransactionDescription, Category) values
	($1,$2,$3,$4,$5,$6)`

	balance, err := t.GetUserBalance(username, account)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return 0, err
	}
	if (balance + transactionAmount) < 0 {
		return 0, errors.New("error. can not decrement balance below zero")
	}

	_, err = db.ExecContext(ctx, query, username, account, transactionAmount, transactionName, transactionDescription, transactionCategory)
	if err != nil {
		return 0, err
	}

	return balance + transactionAmount, nil
}

func (t *RecurringPayment) GetAllReccurringPayments() ([]RecurringPayment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `SELECT paymentid, username, accountname, paymentamount, paymentname, paymentdescription, paymentdate, paymenttype
	FROM foreman.recurring_payment rp
	WHERE rp.paymentid NOT IN 
	(SELECT paymentid FROM foreman.payment_history WHERE paymenthistorydate::date = CURRENT_DATE::date OR paymenthistorystatus != true)
	AND EXTRACT(day from paymentdate) = EXTRACT(day from CURRENT_DATE)`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recurring_payments []RecurringPayment

	for rows.Next() {
		var recurring RecurringPayment
		if err := rows.Scan(&recurring.PaymentID, &recurring.UserName, &recurring.AccountName, &recurring.PaymentAmount, &recurring.PaymentName, &recurring.PaymentDescription, &recurring.PaymentDate, &recurring.PaymentType); err != nil {
			return recurring_payments, err
		}
		recurring_payments = append(recurring_payments, recurring)
	}
	if err = rows.Err(); err != nil {
		return recurring_payments, err
	}
	return recurring_payments, nil
}

func (t *RecurringPayment) GetReccurringPayments(username string) ([]RecurringPayment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `SELECT paymentid, username, paymentamount, paymentname, paymentdescription, paymentdate
	FROM foreman.recurring_payment WHERE username = $1`

	rows, err := db.QueryContext(ctx, query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recurring_payments []RecurringPayment

	for rows.Next() {
		var recurring RecurringPayment
		if err := rows.Scan(&recurring.PaymentID, &recurring.UserName, &recurring.PaymentAmount, &recurring.PaymentName, &recurring.PaymentDescription, &recurring.PaymentDate); err != nil {
			return recurring_payments, err
		}
		recurring_payments = append(recurring_payments, recurring)
	}
	if err = rows.Err(); err != nil {
		return recurring_payments, err
	}
	return recurring_payments, nil
}

func (t *RecurringPayment) AddReccurringPayment(username string, paymentAmount float32, paymentName string, paymentDescription string, paymentDate string) (float32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `INSERT INTO foreman.recurring_payment(
		username, paymentamount, paymentname, paymentdescription, paymentdate)
		VALUES ($1, $2, $3, $4, $5);`

	_, err := db.ExecContext(ctx, query, username, paymentAmount, paymentName, paymentDescription, paymentDate)

	if err != nil {
		return 0, err
	}

	return 1, err
}

func (*PaymentHistory) InsertPaymentHistory(paymentID int, PaymentHistoryStatus bool) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `INSERT INTO foreman.payment_history(
		paymentid, paymenthistorydate, paymenthistorystatus)
		VALUES ($1, $2, $3);`

	currentTime := time.Now()

	_, err := db.ExecContext(ctx, query, paymentID, currentTime.UTC(), PaymentHistoryStatus)

	if err != nil {
		return 0, err
	}

	return 1, err
}

func (t *PaymentHistory) GetPaymentHistory(paymentID int) ([]PaymentHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `SELECT paymenthistoryid, paymentid, paymenthistorydate, paymenthistorystatus
				FROM foreman.payment_history WHERE paymentid = $1`

	rows, err := db.QueryContext(ctx, query, paymentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []PaymentHistory

	for rows.Next() {
		var payment PaymentHistory
		if err := rows.Scan(&payment.PaymentHistoryID, &payment.PaymentID, &payment.PaymentHistoryDate, &payment.PaymentHistoryStatus); err != nil {
			return payments, err
		}
		payments = append(payments, payment)
	}
	if err = rows.Err(); err != nil {
		return payments, err
	}
	return payments, nil
}
