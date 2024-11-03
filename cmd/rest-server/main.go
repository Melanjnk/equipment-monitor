package main

import (
	"log"
	"net/http"
	"github.com/Melanjnk/equipment-monitor/cmd/rest-server/corsrouter"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/controller"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/database"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/repository"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/server/rest"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/service"
)

func main() {
	db, err := database.Connect(
		"postgres", "localhost", 54327, "equipment_api", "postgres", "postgres", false,
	)
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
	defer db.Close()

	repo := repository.NewEquipment(db)
	cont := controller.NewEquipment(
		service.NewEquipment(&repo),
	)

	// Configure router
	router := corsrouter.CORSRouter{}
	eqRouter := router.PathPrefix("/equipment").Subrouter()
	eqRouter.HandleFunc("/", cont.Create).Methods(http.MethodPost)
	eqRouter.HandleFunc("/", cont.Update).Methods(http.MethodPut)
	eqRouter.HandleFunc("/", cont.List).Methods(http.MethodGet)
	eqRouter.HandleFunc("/{id}", cont.Get).Methods(http.MethodGet)
	eqRouter.HandleFunc("/{id}", cont.Delete).Methods(http.MethodDelete)

	http.Handle("/", http.FileServer(http.Dir("./public")))

	server := rest.RestServer{}
	server.StartHTTP(":8080", &router)
}
