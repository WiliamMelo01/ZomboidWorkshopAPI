package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

	// Initialize Redis connection with pool limits & timeouts
	rdb = redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     20,
		MinIdleConns: 5,
	})

	// Test Redis connection
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Falha ao conectar no Redis em %s: %v", redisAddr, err)
	}
	log.Printf("Conectado ao Redis (%s) com sucesso!", redisAddr)

	// Router
	mux := http.NewServeMux()
	mux.HandleFunc("/mappings", handleMappings)
	mux.HandleFunc("/mappings/bulk", handleBulkMappings)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Servidor Zomboid Workshop API iniciado na porta :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}

// Enable CORS for Tauri frontend & Web clients
func setupCORS(w *http.ResponseWriter, r *http.Request) bool {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
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

	// Retorna a totalidade do Hash "workshop_mappings"
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

	// Limita o tamanho do payload para até 10MB por requisição
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	var req BulkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON invalido ou payload muito grande", http.StatusBadRequest)
		return
	}

	if len(req.Mappings) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","processed":0}`))
		return
	}

	// Prepara inserção em pipeline para alta performance e zero estouro de memória
	pipe := rdb.Pipeline()
	count := 0

	for _, m := range req.Mappings {
		modID := strings.TrimSpace(m.ModID)
		workshopID := strings.TrimSpace(m.WorkshopID)

		if modID != "" && workshopID != "" {
			// Salva a versão original e a versão em minúsculo para busca insensível à caixa
			pipe.HSet(ctx, "workshop_mappings", modID, workshopID)
			pipe.HSet(ctx, "workshop_mappings", strings.ToLower(modID), workshopID)
			count++
		}
	}

	if count > 0 {
		_, err := pipe.Exec(ctx)
		if err != nil {
			log.Printf("Erro ao executar Pipeline no Redis: %v", err)
			http.Error(w, "Erro interno ao salvar no banco de dados", http.StatusInternalServerError)
			return
		}
	}

	log.Printf("Sucesso: processado lote com %d mapeamentos", count)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","processed":%d}`, count)
}
