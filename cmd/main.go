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

	// API routes
	mux.HandleFunc("GET /movies", listMovies)

	mux.HandleFunc("GET /movies/{movieID}/seats", bookingHandler().ListSeats)
	mux.HandleFunc("POST /movies/{movieID}/seats/{seatID}/hold", bookingHandler().HoldSeat)
	mux.HandleFunc("PUT /sessions/{sessionID}/confirm", bookingHandler().ConfirmSession)
	mux.HandleFunc("DELETE /sessions/{sessionID}", bookingHandler().ReleaseSession)

	// ✅ FIXED ROOT ROUTING (homepage + static files)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If root → serve homepage
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "static/index.html")
			return
		}

		// otherwise serve static files
		http.FileServer(http.Dir("static")).ServeHTTP(w, r)
	})

	redisAddr := mustGetEnv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASS")

	log.Println("Connecting to Redis at:", redisAddr)

	store := booking.NewRedisStore(redis.NewClient(redisAddr, redisPassword))
	svc := booking.NewService(store)
	_ = svc // already used in handler

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

// helper to avoid redeclaring handler multiple times
func bookingHandler() *booking.Handler {
	redisAddr := mustGetEnv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASS")

	store := booking.NewRedisStore(redis.NewClient(redisAddr, redisPassword))
	svc := booking.NewService(store)

	return booking.NewHandler(svc)
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("%s is required", key)
	}
	return val
}

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
