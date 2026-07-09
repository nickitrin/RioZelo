package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Ocorrencia struct {
	ID             int    `json:"id"`
	MacroCategoria string `json:"macroCategoria"`
	Subcategoria   string `json:"subcategoria"`
	DetalheOutros   string `json:"detalheOutros"`
	Bairro         string `json:"bairro text"`
	Rua            string `json:"rua"`
}

var (
	ocorrenciasChan = make(chan Ocorrencia)
	mu              sync.Mutex
	idCounter       = 0
)

func main() {
	// Middleware simples de CORS para permitir que o Next.js acesse o Go localmente
	cors := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			if r.Method == "OPTIONS" {
				return
			}
			next(w, r)
		}
	}

	// RF03 - Endpoint para receber a ocorrência do formulário
	http.HandleFunc("/ocorrencias", cors(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		var o Ocorrencia
		if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		mu.Lock()
		idCounter++
		o.ID = idCounter
		mu.Unlock()

		// RNF01 - Dispara uma goroutine para empurrar o dado de forma assíncrona para o canal SSE
		go func() {
			ocorrenciasChan <- o
		}()

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"status": "sucesso"}`))
	}))

	// RF04 - Endpoint do Streaming SSE para o Dashboard
	http.HandleFunc("/dashboard/stream", cors(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		for {
			select {
			case o := <-ocorrenciasChan:
				data, _ := json.Marshal(o)
				// O formato do SSE exige o prefixo "data: " seguido de duas quebras de linha (\n\n)
				fmt.Fprintf(w, "data: %s\n\n", data)
				w.(http.Flusher).Flush() 
			case <-r.Context().Done():
				return // Conexão fechada pelo navegador
			}
		}
	}))

	println("🚀 Servidor Go rodando na porta :8080")
	http.ListenAndServe(":8080", nil)
}