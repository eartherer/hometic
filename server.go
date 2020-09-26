package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"server/hometic/zaplogger"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Pair struct {
	DeviceID int64
	UserID   int64
}

func main() {

	fmt.Println("This is hometic")

	r := mux.NewRouter()
	r.Use(zaplogger.Middleware)
	r.Handle("/pair-device", CustomHandlerFunc(PairDeviceHandler(CreatePairDeviceFunc(createPairDeviceFunc)))).Methods(http.MethodPost)

	addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	fmt.Println("addr: ", addr)

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	log.Println("Starting server @")
	log.Fatal(server.ListenAndServe())
}

type CustomHandlerFunc func(w CustomResponseWriter, r *http.Request)

func (handler CustomHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler(&JSONResponseWriter{w}, r)
}

type CustomResponseWriter interface {
	http.ResponseWriter
	JSON(statusCode int, data interface{})
}
type JSONResponseWriter struct {
	http.ResponseWriter
}

func (w *JSONResponseWriter) JSON(statusCode int, data interface{}) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func PairDeviceHandler(device Device) func(w CustomResponseWriter, r *http.Request) {
	return func(w CustomResponseWriter, r *http.Request) {
		zaplogger.L(r.Context()).Info("pair-device")
		//r.Context().Value("logger").(*zap.Logger).Info("pair-device")

		var p Pair
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			w.JSON(http.StatusBadRequest, err.Error())
			return
		}

		defer r.Body.Close()
		fmt.Printf("pair %#v\n", p)
		err = device.Pair(p)
		if err != nil {
			w.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		w.JSON(http.StatusOK, map[string]interface{}{"status": "active"})
		return
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
