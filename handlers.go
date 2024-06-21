package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// getEvents retorna todos os eventos em formato JSON.
func getEvents(w http.ResponseWriter, r *http.Request) {
	// Define o tipo de conteúdo da resposta como JSON
	w.Header().Set("Content-Type", "application/json")

	// Codifica e escreve a lista de eventos como resposta
	if err := json.NewEncoder(w).Encode(data.Events); err != nil {
		http.Error(w, "Failed to encode events", http.StatusInternalServerError)
	}
}

// eventHandler manipula as requisições para eventos específicos e seus sub-recursos.
func eventHandler(w http.ResponseWriter, r *http.Request) {
	// Divide o caminho da URL para determinar o tipo de requisição
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/events/"), "/")

	if len(pathParts) == 1 {
		// Se houver apenas um elemento no pathParts, chama getEvent para obter os detalhes do evento
		getEvent(w, r, pathParts[0])
	} else if len(pathParts) == 2 && pathParts[1] == "spots" {
		// Se houver dois elementos e o segundo for "spots", chama getEventSpots para obter os lugares do evento
		getEventSpots(w, r, pathParts[0])
	} else {
		// Se o caminho não corresponder aos padrões esperados, retorna um erro de requisição inválida
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}
}

// getEvent retorna os detalhes de um evento específico.
func getEvent(w http.ResponseWriter, r *http.Request, eventIDStr string) {
	// Define o tipo de conteúdo da resposta como JSON
	w.Header().Set("Content-Type", "application/json")

	// Converte o ID do evento de string para inteiro
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	// Busca o evento pelo ID
	var foundEvent *Event
	for _, event := range data.Events {
		if event.ID == eventID {
			foundEvent = &event
			break
		}
	}

	// Verifica se o evento foi encontrado
	if foundEvent == nil {
		http.NotFound(w, r)
		return
	}

	// Codifica e envia os detalhes do evento encontrado como resposta
	if err := json.NewEncoder(w).Encode(foundEvent); err != nil {
		http.Error(w, "Failed to encode event", http.StatusInternalServerError)
	}
}

// getEventSpots retorna todos os lugares (spots) de um evento específico.
func getEventSpots(w http.ResponseWriter, r *http.Request, eventIDStr string) {
	// Define o tipo de conteúdo da resposta como JSON
	w.Header().Set("Content-Type", "application/json")

	// Converte o ID do evento de string para inteiro
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	// Lista para armazenar os lugares do evento
	var eventSpots []Spot
	for _, spot := range data.Spots {
		if spot.EventID == eventID {
			eventSpots = append(eventSpots, spot)
		}
	}

	// Verifica se foram encontrados lugares para o evento
	if len(eventSpots) == 0 {
		http.NotFound(w, r)
		return
	}

	// Codifica e envia os lugares do evento como resposta
	if err := json.NewEncoder(w).Encode(eventSpots); err != nil {
		http.Error(w, "Failed to encode spots", http.StatusInternalServerError)
	}
}

// reserveSpot reserva um lugar (spot) para um evento.
func reserveSpot(w http.ResponseWriter, r *http.Request) {
	// Define o tipo de conteúdo da resposta como JSON
	w.Header().Set("Content-Type", "application/json")

	// Divide o caminho da URL para extrair o ID do evento e verificar se é uma requisição de reserva
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/event/"), "/")
	if len(pathParts) < 2 || pathParts[1] != "reserve" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Converte o ID do evento de string para inteiro
	eventID, err := strconv.Atoi(pathParts[0])
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	// Decodifica o payload da requisição para obter o lugar solicitado
	var requestedSpot Spot
	if err := json.NewDecoder(r.Body).Decode(&requestedSpot); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Mutex para garantir que a reserva seja feita de forma segura (sem concorrência)
	mu.Lock()
	defer mu.Unlock()

	// Busca o lugar solicitado pelo nome e verifica se está disponível para reserva
	for i, spot := range data.Spots {
		if spot.EventID == eventID && spot.Name == requestedSpot.Name {
			if spot.Status == "available" {
				// Marca o lugar como reservado
				data.Spots[i].Status = "reserved"
				// Codifica e envia os detalhes do lugar reservado como resposta
				if err := json.NewEncoder(w).Encode(data.Spots[i]); err != nil {
					http.Error(w, "Failed to encode spot", http.StatusInternalServerError)
				}
				return
			}
			// Se o lugar não estiver disponível, retorna um erro
			http.Error(w, "Spot already reserved", http.StatusBadRequest)
			return
		}
	}

	// Se o lugar não for encontrado, retorna um erro de lugar não encontrado
	http.NotFound(w, r)
}
