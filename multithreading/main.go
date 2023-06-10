package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ApiCep struct {
	Status   int    `json:"status"`
	Code     string `json:"code"`
	State    string `json:"state"`
	City     string `json:"city"`
	District string `json:"district"`
	Address  string `json:"address"`
}

func main() {
	fmt.Println("Enter Your CEP (numbers only): ")
	var cep string
	fmt.Scanln(&cep)
	c1 := make(chan string)
	c2 := make(chan string)
	go func() {
		cepParam := cep[0:5] + "-" + cep[len(cep)-3:]
		result, error := SearchApiCep(cepParam)
		if error != nil {
			panic(error)
		}
		c1 <- result.Address + ", " + result.District + ", " + result.City + ", " + result.State + " - " + result.Code
	}()
	go func() {
		result, error := SearchViaCep(cep)
		if error != nil {
			panic(error)
		}
		c2 <- result.Logradouro + ", " + result.Bairro + ", " + result.Localidade + ", " + result.Uf + " - " + result.Cep
	}()

	select {
	case msg1 := <-c1:
		println("API CEP -> ", msg1)
	case msg2 := <-c2:
		println("VIA CEP -> ", msg2)
	case <-time.After(time.Second):
		println("timeout")
	}
}

func SearchApiCep(cep string) (*ApiCep, error) {
	resp, error := http.Get("https://cdn.apicep.com/file/apicep/" + cep + ".json")
	if error != nil {
		return nil, error
	}
	defer resp.Body.Close()
	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}
	var c ApiCep
	error = json.Unmarshal(body, &c)
	if error != nil {
		return nil, error
	}
	return &c, nil
}

func SearchViaCep(cep string) (*ViaCep, error) {
	resp, error := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if error != nil {
		return nil, error
	}
	defer resp.Body.Close()
	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}
	var c ViaCep
	error = json.Unmarshal(body, &c)
	if error != nil {
		return nil, error
	}
	return &c, nil
}
