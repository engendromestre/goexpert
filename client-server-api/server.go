package main

import (
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
	ID          int    `gorm:"primaryKey"`
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
	defer log.Println("Request Finished")
	select {
	case <-time.After(200 * time.Millisecond):
		// print in command line of the server (stdout)
		log.Println("Request successfully processed")
		// print in browser
		exchange, error := SearchExchange()
		if error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(exchange.USDBRL.Bid))
		select {
		case <-time.After(10 * time.Millisecond):
			// register in database here
			CreateExchange(exchange)
		}
	case <-ctx.Done():
		log.Println("Request Cancelled by Client")
	}
}

func CreateExchange(exchange *Exchange) error {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&ExchangeDB{})

	db.Create(&ExchangeDB {     
		Code: exchange.USDBRL.Code,       
		Codein: exchange.USDBRL.Codein,     
		Name: exchange.USDBRL.Name,       
		High: exchange.USDBRL.High,       
		Low: exchange.USDBRL.Low,        
		VarBid: exchange.USDBRL.VarBid,    
		PctChange: exchange.USDBRL.PctChange,  
		Bid: exchange.USDBRL.Bid,        
		Ask: exchange.USDBRL.Ask,        
		Timestamp: exchange.USDBRL.Timestamp,
		Create_date: exchange.USDBRL.Create_date,
	})
	// select all
	var exchangesDB []ExchangeDB
	db.Find(&exchangesDB)
	for _, ext := range exchangesDB {
		fmt.Println(ext)
	}

	return nil
}

func SearchExchange() (*Exchange, error) {
	resp, error := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
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
