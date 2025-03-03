package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/database"
)

func main() {
	db, err := database.Connect(
		`postgres`, `localhost`, 54327, `equipment_api`, `postgres`, `postgres`, false,
	)
	if err == nil {
		defer db.Close()

		const basePath = `./internal/app/registry_service/recreate_db/`
		var content []byte
		if content, err = os.ReadFile(basePath + `script.list`); err == nil {
			for _, scriptName := range strings.Split(string(content), "\n") {
				if scriptName = strings.TrimSpace(scriptName); len(scriptName) > 0 {
					content, err = os.ReadFile(fmt.Sprintf(`%ssql/%s.sql`, basePath, scriptName))
					if err != nil {
						goto ERROR
					}
					_, err = db.Exec(string(content))
					if err != nil {
						goto ERROR
					}
				}
			}
			return // Success
		}
	}
ERROR:
	log.Fatalln(err)
	panic(err)
}
