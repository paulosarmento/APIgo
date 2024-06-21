package main

import (
	"log"
	"net/http"
)

func main() {
	// Carregar dados do arquivo JSON
	loadData()

	// Registrar manipuladores de rotas
	http.HandleFunc("/events", getEvents)
	http.HandleFunc("/events/", eventHandler)
	http.HandleFunc("/event/", reserveSpot)

	// Iniciar servidor HTTP na porta 8080
	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
