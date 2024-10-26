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
	eqController := controller.NewEquipment(
		service.NewEquipment(&repo),
	)

	// Configure router
	router := corsrouter.CORSRouter{}
	eqRouter := router.PathPrefix("/equipment").Subrouter()
	eqRouter.HandleFunc("/", eqController.List).Methods("GET")
	eqRouter.HandleFunc("/", eqController.Create).Methods("POST")
	eqRouter.HandleFunc("/{id}", eqController.Update).Methods("PUT")
	eqRouter.HandleFunc("/{id}", eqController.Get).Methods("GET")
	eqRouter.HandleFunc("/{id}", eqController.Delete).Methods("DELETE")

	http.Handle("/", http.FileServer(http.Dir("./public")))

	server := rest.RestServer{}
	server.StartHTTP(":8080", &router)
}
