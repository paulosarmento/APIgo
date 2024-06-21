package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
)

// Event representa a estrutura de dados de um evento.
type Event struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Organization string `json:"organization"`
	Date         string `json:"date"`
	Price        int    `json:"price"`
	Rating       string `json:"rating"`
	ImageURL     string `json:"image_url"`
	CreatedAt    string `json:"created_at"`
	Location     string `json:"location"`
}

// Spot representa a estrutura de dados de um lugar (spot).
type Spot struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	EventID int    `json:"event_id"`
}

// Data encapsula os dados lidos do arquivo JSON.
type Data struct {
	Events []Event `json:"events"`
	Spots  []Spot  `json:"spots"`
}

var (
	data Data       // Variável global para armazenar os dados carregados do arquivo JSON.
	mu   sync.Mutex // Mutex para sincronização de acesso aos dados compartilhados.
)

// loadData carrega os dados do arquivo data.json para a variável global data.
func loadData() {
	// Abre o arquivo data.json
	file, err := os.Open("data.json")
	if err != nil {
		log.Fatalf("Failed to open data file: %v", err)
	}
	defer file.Close() // Garante que o arquivo seja fechado após a leitura

	// Lê o conteúdo completo do arquivo
	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read data file: %v", err)
	}

	// Deserializa os dados JSON para a estrutura Data
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		log.Fatalf("Failed to parse data file: %v", err)
	}
}
