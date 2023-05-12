package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Exchange struct {
	ID     int `gorm:"primaryKey"`
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

func main() {
	http.HandleFunc("/cotacao", SearchExchangeHandler)
	http.ListenAndServe(":8080", nil)

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Exchange{})
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
			json.NewEncoder(w).Encode(exchange)
		}
	case <-ctx.Done():
		log.Println("Request Cancelled by Client")
	}
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
