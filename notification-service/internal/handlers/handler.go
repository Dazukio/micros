package handlers

//
//import (
//	"encoding/json"
//	"io"
//	"log"
//	"net/http"
//)
//
//func Notify(w http.ResponseWriter, r *http.Request) {
//	body, err := io.ReadAll(r.Body)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//	}
//	defer r.Body.Close()
//	notific := struct {
//		Message string `json:"title"`
//	}{}
//	if err := json.Unmarshal(body, &notific); err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//	}
//	log.Println("Sendind notification:", notific.Message)
//	log.Printf("%+v", notific)
//
//}
