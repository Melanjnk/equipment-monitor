package main

import (
	"log"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/database"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/migrations"
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
	if err := migrations.DropTableEquipment(db); err != nil {
		log.Fatalln(err)
		panic(err)
	}
	if err := migrations.CreateTableEquipment(db); err != nil {
		log.Fatalln(err)
		panic(err)
	}
}
