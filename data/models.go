package data

import (
	"context"
	"database/sql"

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
	PaymentID          int     `json:"-"`
	UserName           string  `json:"username"`
	PaymentAmount      float32 `json:"amount"`
	PaymentName        string  `json:"paymentName"`
	PaymentDescription string  `json:"paymentDescription"`
	PaymentDate        string  `json:"paymentDate"`
}

type PaymentHistory struct {
	PaymentHistoryID     int    `json:"-"`
	PaymentID            int    `json:"paymentID"`
	PaymentHistoryDate   string `json:"paymentHistoryDate"`
	PaymentHistoryStatus bool   `json:"paymentHistoryStatus"`
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
