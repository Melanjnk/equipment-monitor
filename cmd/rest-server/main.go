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

	equipmentRepository := repository.NewEquipment(db)
	equipmentController := controller.NewEquipment(
		service.NewEquipment(equipmentRepository),
	)

	// Configure router
	router := corsrouter.CORSRouter{}
	equipmentRouter := router.PathPrefix("/equipment").Subrouter()
	equipmentRouter.HandleFunc("/", equipmentController.Create).Methods(http.MethodPost)
	equipmentRouter.HandleFunc("/{id}", equipmentController.UpdateByIds).Methods(http.MethodPatch)
	equipmentRouter.HandleFunc("/", equipmentController.UpdateByConditions).Methods(http.MethodPatch)
	equipmentRouter.HandleFunc("/{id}", equipmentController.DeleteByIds).Methods(http.MethodDelete)
	equipmentRouter.HandleFunc("/", equipmentController.DeleteByConditions).Methods(http.MethodDelete)
	equipmentRouter.HandleFunc("/{id}", equipmentController.FindById).Methods(http.MethodGet)
	equipmentRouter.HandleFunc("/", equipmentController.FindByConditions).Methods(http.MethodGet)

	http.Handle("/", http.FileServer(http.Dir("./public")))

	server := rest.RestServer{}
	server.StartHTTP(":8080", &router)
}
