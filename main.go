package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx = context.Background()

type MappingItem struct {
	ModID      string `json:"modId"`
	WorkshopID string `json:"workshopId"`
}

type BulkRequest struct {
	Mappings []MappingItem `json:"mappings"`
}

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	// Initialize Redis connection
	rdb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Test Redis connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Falha ao conectar no Redis em %s: %v", redisAddr, err)
	}
	log.Printf("Conectado ao Redis (%s) com sucesso!", redisAddr)

	// Routes
	http.HandleFunc("/mappings", handleMappings)
	http.HandleFunc("/mappings/bulk", handleBulkMappings)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor Zomboid Workshop API iniciado na porta :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}

// Enable CORS for Tauri frontend & Web clients
func setupCORS(w *http.ResponseWriter, r *http.Request) bool {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		(*w).WriteHeader(http.StatusNoContent)
		return true
	}
	return false
}

func handleMappings(w http.ResponseWriter, r *http.Request) {
	if setupCORS(&w, r) {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Metodo nao permitido", http.StatusMethodNotAllowed)
		return
	}

	// Fetch all mappings from Redis Hash
	mappings, err := rdb.HGetAll(ctx, "workshop_mappings").Result()
	if err != nil {
		log.Printf("Erro ao buscar no Redis: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mappings); err != nil {
		log.Printf("Erro ao encodar JSON: %v", err)
	}
}

func handleBulkMappings(w http.ResponseWriter, r *http.Request) {
	if setupCORS(&w, r) {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Metodo nao permitido", http.StatusMethodNotAllowed)
		return
	}

	var req BulkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON invalido", http.StatusBadRequest)
		return
	}

	if len(req.Mappings) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","message":"no items to process"}`))
		return
	}

	// Prepare fields for bulk HSet
	var pairs []string
	for _, m := range req.Mappings {
		modID := strings.TrimSpace(m.ModID)
		workshopID := strings.TrimSpace(m.WorkshopID)
		if modID != "" && workshopID != "" {
			pairs = append(pairs, modID, workshopID)
		}
	}

	if len(pairs) > 0 {
		err := rdb.HSet(ctx, "workshop_mappings", pairs).Err()
		if err != nil {
			log.Printf("Erro ao salvar no Redis: %v", err)
			http.Error(w, "Erro interno ao salvar os dados", http.StatusInternalServerError)
			return
		}
	}

	log.Printf("Sucesso: processado lote com %d mapeamentos", len(pairs)/2)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","processed":%d}`, len(pairs)/2)
}
