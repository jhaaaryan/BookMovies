package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jhaaaryan/BookMovies/internal/adapters/redis"
	"github.com/jhaaaryan/BookMovies/internal/booking"
	"github.com/jhaaaryan/BookMovies/internal/utils"
)

func main() {
	mux := http.NewServeMux()

	// -------------------------
	// ENV CONFIG
	// -------------------------
	redisAddr := mustGetEnv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASS")

	log.Println("Connecting to Redis at:", redisAddr)

	// -------------------------
	// INIT SERVICE ONCE (IMPORTANT FIX)
	// -------------------------
	store := booking.NewRedisStore(redis.NewClient(redisAddr, redisPassword))
	svc := booking.NewService(store)
	handler := booking.NewHandler(svc)

	// -------------------------
	// API ROUTES
	// -------------------------
	mux.HandleFunc("GET /movies", listMovies)

	mux.HandleFunc("GET /movies/{movieID}/seats", handler.ListSeats)
	mux.HandleFunc("POST /movies/{movieID}/seats/{seatID}/hold", handler.HoldSeat)
	mux.HandleFunc("PUT /sessions/{sessionID}/confirm", handler.ConfirmSession)
	mux.HandleFunc("DELETE /sessions/{sessionID}", handler.ReleaseSession)

	// -------------------------
	// STATIC + HOME ROUTING FIX
	// -------------------------
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "static/index.html")
			return
		}

		http.FileServer(http.Dir("static")).ServeHTTP(w, r)
	})

	// -------------------------
	// SERVER SETUP (Render safe)
	// -------------------------
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

// -------------------------
// HELPERS
// -------------------------
func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("%s is required", key)
	}
	return val
}

// -------------------------
// MOCK DATA
// -------------------------
var movies = []movieResponse{
	{ID: "inception", Title: "Inception", Rows: 5, SeatsPerRow: 8},
	{ID: "the-prestige", Title: "The Prestige", Rows: 5, SeatsPerRow: 8},
	{ID: "the-dark-knight", Title: "The Dark Knight", Rows: 5, SeatsPerRow: 8},
	{ID: "dunkirk", Title: "Dunkirk", Rows: 4, SeatsPerRow: 6},
	{ID: "tenet", Title: "Tenet", Rows: 4, SeatsPerRow: 6},
	{ID: "oppenheimer", Title: "Oppenheimer", Rows: 5, SeatsPerRow: 8},
	{ID: "the-odyssey", Title: "The Odyssey", Rows: 5, SeatsPerRow: 10},
}

func listMovies(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, movies)
}

type movieResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Rows        int    `json:"rows"`
	SeatsPerRow int    `json:"seats_per_row"`
}
