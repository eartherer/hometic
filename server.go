package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Pair struct {
	DeviceID int64
	UserID   int64
}

func main() {

	fmt.Println("This is hometic")

	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := zap.NewExample()
			l = l.With(zap.Namespace("hometic"), zap.String("I'm", "gopher"))
			l.Info("PairDevice")
			ctx := context.WithValue(r.Context(), "logger", l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	r.Handle("/pair-device", PairDeviceHandler(CreatePairDeviceFunc(createPairDeviceFunc))).Methods(http.MethodPost)

	addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	fmt.Println("addr: ", addr)

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	log.Println("Starting server @")
	log.Fatal(server.ListenAndServe())
}

func PairDeviceHandler(device Device) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Context().Value("logger").(*zap.Logger).Info("pair-device")
		var p Pair
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		defer r.Body.Close()
		fmt.Printf("pair %#v\n", p)
		err = device.Pair(p)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		w.Write([]byte(`{"status":"active"}`))
	}
}

type Device interface {
	Pair(p Pair) error
}

type CreatePairDeviceFunc func(p Pair) error

func (fn CreatePairDeviceFunc) Pair(p Pair) error {
	return fn(p)
}

func createPairDeviceFunc(p Pair) error {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO pairs VALUES($1,$2);", p.DeviceID, p.UserID)
	return err
}
