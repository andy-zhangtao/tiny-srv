package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	ModuleName = "I am tiny-srv"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {

	logrus.WithFields(logrus.Fields{
		"version": "v0.0.2",
		"auther":  "andy-zhangtao",
		"email":   "ztao8607@gmail.com",
	})

	logrus.WithFields(logrus.Fields{
		"apis": []string{
			"/_ping",
			"/v1/test",
			"/v1/echo",
		},
	}).Info("API LIST")

	r := mux.NewRouter()

	r.HandleFunc("/_ping", func(writer http.ResponseWriter, request *http.Request) {
		logrus.WithFields(logrus.Fields{"url": "/_ping", "Request RemoteAddr": request.RemoteAddr, "Header": request.Header, "Host": request.Host})
		d, _ := ioutil.ReadAll(request.Body)
		logrus.WithFields(logrus.Fields{"body": string(d)}).Info(ModuleName)
		writer.Write([]byte(ModuleName))
	})

	r.HandleFunc("/v1/test", func(writer http.ResponseWriter, request *http.Request) {
		logrus.WithFields(logrus.Fields{"url": "/v1/test", "Request RemoteAddr": request.RemoteAddr, "Header": request.Header, "Host": request.Host})
		d, _ := ioutil.ReadAll(request.Body)
		logrus.WithFields(logrus.Fields{"header": request.Header, "body": string(d)}).Info(ModuleName)
		writer.Write([]byte("U sucd!"))
	})

	r.HandleFunc("/v1/echo", func(writer http.ResponseWriter, request *http.Request) {
		logrus.WithFields(logrus.Fields{"url": "/v1/echo", "Request RemoteAddr": request.RemoteAddr, "Header": request.Header, "Host": request.Host})
		data, err := json.Marshal(request.Header)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":  err.Error(),
				"origin": request.Header,
			}).Error("JSON Marshal Error")
		}

		qd, err := json.Marshal(request.URL.Query())
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":  err.Error(),
				"origin": request.URL.Query(),
			}).Error("JSON Marshal Error")
		}

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Read Body Error")
		}

		type respon struct {
			Header string `json:"header"`
			Query  string `json:"query"`
			Body   string `json:"body"`
		}

		r := &respon{
			Header: string(data),
			Query:  string(qd),
			Body:   string(body),
		}

		logrus.WithFields(logrus.Fields{"header": request.Header, "body": string(body)}).Info(ModuleName)
		json.NewEncoder(writer).Encode(r)
	})

	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	logrus.Println(http.ListenAndServe(":8000", handlers.CORS(headersOk, originsOk, methodsOk)(r)))
}
