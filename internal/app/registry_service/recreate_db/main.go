package main

import (
	"log"
	"github.com/jmoiron/sqlx"
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
	for _, function := range []func(*sqlx.DB) error {
		migrations.DropTableEquipment,
		migrations.CreateGeneratorUUIDv6,
		migrations.CreateAutoUpdate,
		migrations.CreateTableEquipment,
	} {
		if err := function(db); err != nil {
			log.Fatalln(err)
			panic(err)
		}
	}
}
