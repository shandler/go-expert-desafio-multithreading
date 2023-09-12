package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
)

type ViaCEP struct {
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
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

func BuscarAPI(url string, ch chan<- struct {
	URL     string
	Result  interface{}
	Elapsed float64
}) {
	start := time.Now()
	resp, err := http.Get(url)
	elapsed := time.Since(start).Seconds()

	if err != nil {
		ch <- struct {
			URL     string
			Result  interface{}
			Elapsed float64
		}{URL: url, Result: fmt.Sprintf("Erro ao fazer a requisição para %s: %s", url, err), Elapsed: elapsed}
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ch <- struct {
			URL     string
			Result  interface{}
			Elapsed float64
		}{URL: url, Result: fmt.Sprintf("API %s retornou código de status %d", url, resp.StatusCode), Elapsed: elapsed}
		return
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		ch <- struct {
			URL     string
			Result  interface{}
			Elapsed float64
		}{URL: url, Result: fmt.Sprintf("Erro ao decodificar JSON da API %s: %s", url, err), Elapsed: elapsed}
		return
	}

	ch <- struct {
		URL     string
		Result  interface{}
		Elapsed float64
	}{URL: url, Result: result, Elapsed: elapsed}
}

func main() {
	cep := os.Args[1]

	apicepURL := "https://cdn.apicep.com/file/apicep/" + cep + ".json"
	viacepURL := "https://cdn.apicep.com/file/apicep/" + cep + ".json"
	//viacepURL := "http://viacep.com.br/ws/" + cep + "/json"

	ch := make(chan struct {
		URL     string
		Result  interface{}
		Elapsed float64
	}, 2)

	go BuscarAPI(apicepURL, ch)
	go BuscarAPI(viacepURL, ch)

	var result struct {
		URL     string
		Result  interface{}
		Elapsed float64
	}

	fmt.Println(" ")

	select {
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout - nenhuma API respondeu dentro do tempo limite")
	case result = <-ch:
		fmt.Printf("Resultado da API %s (Tempo: %.2f segundos):\n", result.URL, result.Elapsed)
		switch result.URL {
		case apicepURL:
			//RetornoCEP := ApiCep{}
			// if data, ok := result.Result.(map[string]interface{}); ok {
			// 	if err := mapstructure.Decode(data, &RetornoCEP); err == nil {
			// 		imprimirStruct(RetornoCEP)
			// 	}
			// }
			converterEImprimir(result.Result, ApiCep{})
		case viacepURL:
			//RetornoCEP := ViaCEP{}
			// if data, ok := result.Result.(map[string]interface{}); ok {
			// 	if err := mapstructure.Decode(data, &RetornoCEP); err == nil {
			// 		imprimirStruct(RetornoCEP)
			// 	}
			// }
			converterEImprimir(result.Result, ViaCEP{})
		}

	}
}

func converterEImprimir(data interface{}, estrutura interface{}) {
	if dataMap, ok := data.(map[string]interface{}); ok {
		println("Aqui 01")
		if err := mapstructure.Decode(dataMap, estrutura); err != nil {
			println(estrutura)
			imprimirStruct(estrutura)
		}
	}
}

func imprimirStruct(data interface{}) {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Struct {
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			value := val.Field(i)
			fmt.Printf("%s: %v\n", field.Name, value.Interface())
		}
	}
}
