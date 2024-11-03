package controller

import (
	"fmt"
	"encoding/json"
	"net/http"
)

const(
	parameterIsRequired string = "Parameter `%s` is required."
	invalidJSONData = "Invalid JSON data: `%v`"
	unableToFindEquipment = "Unable to find equipment `%v` for %s"
	equipmentIdError = "%s equipment `%v` error: %v"
	equipmentActionIsPerformed = "Equipment `%v` is %s"
)

func writeMessage(writer http.ResponseWriter, status int, message string, parameters ...interface{}) {
	writer.WriteHeader(status)
	fmt.Fprintf(writer, message, parameters...)
}

func writeJSON(writer http.ResponseWriter, status int, data interface{}) {
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(data)
}
