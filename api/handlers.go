package api

import(
	"verdmell/check"
	"verdmell/service"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Welcome to Verdmell API!")
}

func GetAllChecks(w http.ResponseWriter, r *http.Request) {
	checks := box.GetObject(CHECKS).(*check.CheckSystem)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	
  if err := json.NewEncoder(w).Encode(checks.GetChecks()); err != nil {
  	panic(err)
  }
}

func GetCheck(w http.ResponseWriter, r *http.Request) {
	checks := box.GetObject(CHECKS).(*check.CheckSystem)
	vars := mux.Vars(r)
	check := vars["check"]
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

  if err := json.NewEncoder(w).Encode(checks.GetCheck(check)); err != nil {
  	panic(err)
  }
}



func GetAllServices(w http.ResponseWriter, r *http.Request) {
	services := box.GetObject(SERVICES).(*service.ServiceSystem)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	
  if err := json.NewEncoder(w).Encode(services.GetServices()); err != nil {
  	panic(err)
  }
}

func GetService(w http.ResponseWriter, r *http.Request) {
	services := box.GetObject(SERVICES).(*service.ServiceSystem)
	vars := mux.Vars(r)
	service := vars["service"]
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

  if err := json.NewEncoder(w).Encode(services.GetService(service)); err != nil {
  	panic(err)
  }
}