package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Exchange struct {
	USDBRL struct {
		Code        string `json:"code"`
		Codein      string `json:"codein"`
		Name        string `json:"name"`
		High        string `json:"high"`
		Low         string `json:"low"`
		VarBid      string `json:"varBid"`
		PctChange   string `json:"pctChange"`
		Bid         string `json:"bid"`
		Ask         string `json:"ask"`
		Timestamp   string `json:"timestamp"`
		Create_date string `json:"create_date"`
	}
}

type ExchangeDB struct {
	ID          int `gorm:"primaryKey"`
	Code        string
	Codein      string
	Name        string
	High        string
	Low         string
	VarBid      string
	PctChange   string
	Bid         string
	Ask         string
	Timestamp   string
	Create_date string
}

func main() {
	http.HandleFunc("/cotacao", SearchExchangeHandler)
	http.ListenAndServe(":8080", nil)
}

func SearchExchangeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	ctx := r.Context()
	log.Println("Request Initiate")
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()
	defer log.Println("Request Finished")

	res, err := SearchExchange(ctx)
	if err != nil {
		log.Println("Error consuming API")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res.USDBRL.Bid))

	ctx, cancel = context.WithTimeout(ctx, 10*time.Nanosecond)
	// register in database here
	err = CreateExchange(ctx, res)
	if err != nil {
		log.Println("Error writing record")
	}
}

func CreateExchange(ctx context.Context, exchange *Exchange) error {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	db.AutoMigrate(&ExchangeDB{})

	if err := db.WithContext(ctx).Create(&ExchangeDB{
		Code:        exchange.USDBRL.Code,
		Codein:      exchange.USDBRL.Codein,
		Name:        exchange.USDBRL.Name,
		High:        exchange.USDBRL.High,
		Low:         exchange.USDBRL.Low,
		VarBid:      exchange.USDBRL.VarBid,
		PctChange:   exchange.USDBRL.PctChange,
		Bid:         exchange.USDBRL.Bid,
		Ask:         exchange.USDBRL.Ask,
		Timestamp:   exchange.USDBRL.Timestamp,
		Create_date: exchange.USDBRL.Create_date,
	}).Error; err != nil {
		return err
	}

	// select all
	var exchangesDB []ExchangeDB
	if err := db.Find(&exchangesDB).Error; err != nil {
		return err
	}
	for _, ext := range exchangesDB {
		fmt.Println(ext)
	}
	return nil
}

func SearchExchange(ctx context.Context) (*Exchange, error) {
	req, error := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if error != nil {
		return nil, error
	}
	resp, error := http.DefaultClient.Do(req)
	if error != nil {
		return nil, error
	}
	defer resp.Body.Close()

	body, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}

	var e Exchange
	error = json.Unmarshal(body, &e)
	if error != nil {
		return nil, error
	}
	return &e, nil
}
