package main

import (
	"fmt"
	"net/http"
	"time"
)

func fetchFromAPI(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	elapsed := time.Since(start).Seconds()

	if err != nil {
		ch <- fmt.Sprintf("Erro ao fazer a requisição para %s: %s", url, err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ch <- fmt.Sprintf("API %s retornou código de status %d", url, resp.StatusCode)
		return
	}

	ch <- fmt.Sprintf("Resultado da API %s: %s", url, elapsed)
}

func main() {
	cep := "01001-000" // Substitua pelo CEP desejado

	apicepURL := "https://cdn.apicep.com/file/apicep/" + cep + ".json"
	viacepURL := "http://viacep.com.br/ws/" + cep + "/json"

	ch := make(chan string, 2)

	go fetchFromAPI(apicepURL, ch)
	go fetchFromAPI(viacepURL, ch)

	select {
	case result := <-ch:
		fmt.Println(result)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout - nenhuma API respondeu dentro do tempo limite")
	}
}
