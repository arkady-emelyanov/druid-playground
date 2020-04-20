package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type TransactionStatus string

const (
	Created   TransactionStatus = "created"
	Confirmed                   = "confirmed"
	Cancelled                   = "cancelled"
)

type Transaction struct {
	Timestamp  time.Time
	MerchantId uint
	ClientId   uint
	Amount     uint
	Status     TransactionStatus
}

func main() {
	dayTransactionsChan := make(chan *Transaction, 1000)
	timeZone, _ := time.LoadLocation("Europe/Moscow")
	startDate := time.Date(2020, 4, 27, 0, 0, 0, 0, timeZone)
	numberOfDays := 30
	numberOfClients := 1500000
	maxTransactionsPerClientPerDay := 5

	// receive transactions and put all of them to file
	writeDoneChan := make(chan struct{}, 0)
	go func(dayChan <-chan *Transaction) {
		fp, _ := os.Create("./generated/events.json")
		defer fp.Close()

		i := 0
		for tr := range dayChan {
			if i%10000 == 0 {
				fmt.Println("Syncing file to disk, written:", i)
				fp.Sync()
			}

			b, err := json.Marshal(tr)
			if err != nil {
				panic(err)
			}

			if _, err := fp.Write(b); err != nil {
				panic(err)
			}
			if _, err := fp.Write([]byte("\n")); err != nil {
				panic(err)
			}
			i++
		}

		writeDoneChan <- struct{}{}
	}(dayTransactionsChan)

	// generate all transactions
	for d := 0; d < numberOfDays; d++ {
		for c := 0; c < numberOfClients; c++ {
			if c%10000 == 0 {
				fmt.Println("clients done:", c, "for day:", d)
			}

			numToday := rand.Intn(maxTransactionsPerClientPerDay) + 1
			startHour := rand.Intn(20-maxTransactionsPerClientPerDay) + 8

			for i := 0; i < numToday; i++ {
				trDate := startDate.
					AddDate(0, 0, d).
					Add(time.Hour * time.Duration(startHour)).
					Add(time.Minute * time.Duration(rand.Intn(59)))

				dayTransactionsChan <- &Transaction{
					Timestamp:  trDate,
					MerchantId: uint(rand.Intn(1000)),
					ClientId:   uint(c),
					Amount:     uint(rand.Intn(10000)),
					Status:     Created,
				}
				startHour++
			}
		}
	}

	close(dayTransactionsChan)
	fmt.Println("finished transaction generation")

	<-writeDoneChan
	fmt.Println("finished writing")
}
