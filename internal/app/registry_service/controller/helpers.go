package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"github.com/gofrs/uuid"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/dtos"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/stringset"
)

const(
	actionCreate = "create"
	actionUpdate = "update"
	actionDelete = "delete"
	actionSearch = "search"
)

const(
	parameterIsRequired		= "Parameter `%s` is required for %s."
	invalidJSONData			= "Invalid JSON data: `%v`"
	equipmentNotFound		= "Unable to find equipment #%v for %s"
	invalidEquipmentId		= "Invalid id of equipment to %s: %v"
	invalidGETParameters	= "Invalid GET parameters on %s: %v"
	actionFailed			= "Failed to %s: %v"
)

func respond(writer http.ResponseWriter, status int, data any) {
	if data == nil {
		writer.WriteHeader(status)
	} else {
		writer.Header().Set(`Content-Type`, `application/json`)
		writer.WriteHeader(status)
		_ = json.NewEncoder(writer).Encode(data)
	}
}

func respondNoContent(writer http.ResponseWriter) {
	respond(writer, http.StatusNoContent, nil)
}

func respondCreated(writer http.ResponseWriter, data any) {
	respond(writer, http.StatusCreated, data)
}

func respondOK(writer http.ResponseWriter, data any) {
	respond(writer, http.StatusOK, data)
}

func respondMulti(writer http.ResponseWriter, data any) {
	respond(writer, http.StatusMultiStatus, data)
}

func respondError(writer http.ResponseWriter, status int, err error, extra any) {
	response := map[string]any{"error": err}
	if extra != nil {
		response["details"] = err
	}
	respond(writer, status, extra)
}

func respondBadRequest(writer http.ResponseWriter, err error, extra any) {
	respondError(writer, http.StatusBadRequest, err, extra)
}

func respondNotFound(writer http.ResponseWriter, err error, extra any) {
	respondError(writer, http.StatusNotFound, err, extra)
}

func respondInternalError(writer http.ResponseWriter, err error, extra any) {
	respondError(writer, http.StatusInternalServerError, err, extra)
}

func isBreaker(b byte) bool {
	switch b {
		case ',', ' ', '\t', '\r', '\n':
			return true
	}
	return false
}

func isDecimal(str string) bool {
	for _, c := range str {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func parseIds(idString string) (stringset.StringSet, error) {
	var buffer []byte
	var b byte
	ids := stringset.New()
	for i, l := 0, len(idString); i < l; i++ {
		b = idString[i]
		if isBreaker(b) {
			continue
		}
		for {
			buffer = append(buffer, b)
			i++
			if i == l {
				break
			}
			b = idString[i]
			if isBreaker(b) {
				break
			}
		}
		id := string(buffer)
		if _, err := uuid.FromString(id); err != nil { // Validation
			return ids, err
		}
		ids.Include(id)
		buffer = buffer[:0]
	}
	return ids, nil
}

func FromRequestJSON(request *http.Request) (*dtos.EquipmentCreate, []dtos.EquipmentCreate, error) {
	buffer, err := io.ReadAll(request.Body)
	if err == nil {
		for i, b := range buffer {
			switch b {
				case ' ', '\t', '\r', '\n': // Skip whitespace
					continue
				case '{': // Single object
					var equipmentCreate dtos.EquipmentCreate
					err = json.Unmarshal(buffer[i:], &equipmentCreate)
					if err == nil {
						err = equipmentCreate.Validate()
						if err == nil {
							return &equipmentCreate, nil, nil
						}
					}
					goto ERROR
				case '[': // Array of objects
					var equipmentCreateMany []dtos.EquipmentCreate
					err = json.Unmarshal(buffer[i:], &equipmentCreateMany)
					if err == nil {
						for _, equipmentCreate := range equipmentCreateMany {
							err = equipmentCreate.Validate()
							if err != nil {
								goto ERROR
							}
						}
						return nil, equipmentCreateMany, nil
					}
				default: // Invalid character
					err = fmt.Errorf("Unexpected character `%r`", rune(b))
					goto ERROR
			}
		}
		err = errors.New("Request's body is empty")
	}
ERROR:
	return nil, nil, err
}
