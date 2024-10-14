package main

import (
	"fmt"
	"sqlx"
	"github.com/gofrs/uuid"

	eqm "github.com/Melanjnk/equipment-monitor/internal/app/registry-service"
)

func main() {
	// PostgreSQL connection string
	dsn := "host=localhost user=docker password=docker dbname=equipment_api port=54327 sslmode=disable"

	// Open connection to PostgreSQL using sqlx
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// Create table with params as JSONB NOT NULL
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS equipment (
			id UUID PRIMARY KEY,
			type SMALLINT,
			status SMALLINT,
			params JSONB NOT NULL,
			CONSTRAINT type_check CHECK (type IN (0, 1, 2, 3)),
			CONSTRAINT status_check CHECK (status IN (0, 1, 2))
		);
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Table created successfully!")

	// Generate a UUID version 6 (using a library)
	newUUID, err := uuid.NewV6()
	if err != nil {
		log.Fatalln("Failed to generate UUID v6:", err)
	}

	// Insert data with Params as JSONB (not null)
	newEquipment := Equipment{
		Id: newUUID,
		Type: eqm.DrillMachine,
		Status: eqm.Operational,
		Params: map[string]interface{}{
			"power":  "500W",
			"weight": "1.5kg",
			"voltage": 220,
		},
	}

	insertSQL := `INSERT INTO equipment (id, type, status, params) VALUES (:id, :type, :status, :params)`

	_, err = db.NamedExec(insertSQL, newEquipment)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Inserted equipment with JSON params successfully!")
}
