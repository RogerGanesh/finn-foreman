package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/RogerGanesh/finn-foreman/data"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Config struct {
	DB     *sql.DB
	Models data.Models
}

const (
	webPort = "9002"
)

type RecurringPayment struct {
	PaymentID          int     `json:"-"`
	UserName           string  `json:"username"`
	PaymentAmount      float32 `json:"amount"`
	PaymentName        string  `json:"paymentName"`
	PaymentDescription string  `json:"paymentDescription"`
	PaymentDate        string  `json:"paymentDate"`
}

var counts int64

func main() {
	if os.Getenv("DSN") == "" {
		os.Setenv("DSN", "host=localhost port=5433 user=postgres password=password dbname=finn sslmode=disable timezone=UTC connect_timeout=5")
	}
	log.Println("Starting transactions service")

	conn := connectToDB()

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	ticker := time.NewTicker(5 * time.Second)
	for t := range ticker.C {
		app.checkRecurringPayments(t)
	}

	// srv := &http.Server{
	// 	Addr:    fmt.Sprintf(":%s", webPort),
	// 	Handler: app.routes(),
	// }

	// err := srv.ListenAndServe()

	// if err != nil {
	// 	log.Panic(err)
	// }
}

func (app *Config) checkRecurringPayments(t time.Time) error {
	recurring, err := app.Models.RecurringPayment.GetAllReccurringPayments()

	for _, reccurance := range recurring {
		jsonbody, err := json.Marshal(reccurance)
		if err != nil {
			// do error check
			fmt.Println(err)
			return err
		}
		recurr := RecurringPayment{}
		if err := json.Unmarshal(jsonbody, &recurr); err != nil {
			// do error check
			fmt.Println(err)
			return err
		}

		var username = recurr.UserName
		recurrance_date, _ := time.Parse(time.RFC3339, recurr.PaymentDate)
		currentTime := time.Now()

		if recurrance_date.Truncate(24 * time.Hour).Equal(currentTime.Truncate(24 * time.Hour)) {
			balance, err := app.Models.RecurringPayment.GetUserBalance(username)

			res, err := app.Models.RecurringPayment.UpdateBalance(username, -recurr.PaymentAmount, recurr.PaymentName, recurr.PaymentDescription)
			if err != nil {
				// do error check
				fmt.Println(err)
				app.Models.PaymentHistory.InsertPaymentHistory(recurr.PaymentID, false)
				return err
			}

			app.Models.PaymentHistory.InsertPaymentHistory(recurr.PaymentID, true)
			log.Println(balance, res)
		}
	}

	if err != nil {
		log.Panic(err)
		return err
	}
	return nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			counts++
		} else {
			log.Println("Connected to Postgres")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}
