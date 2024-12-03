package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/dtos"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/service"
)

type Equipment struct {
	service service.Equipment
}

func NewEquipment(service service.Equipment) Equipment {
	return Equipment{service: service}
}

func (controller *Equipment) Create(writer http.ResponseWriter, request *http.Request) {
	if equipmentCreate, equipmentCreateMany, err := FromRequestJSON(request); err != nil {
		respondBadRequest(writer, fmt.Errorf(invalidJSONData, err), nil)
	} else if equipmentCreateMany != nil {
		if ids, err := controller.service.CreateMany(equipmentCreateMany); err != nil {
			respondBadRequest(writer, fmt.Errorf(equipmentNotCreated, err), nil)
		} else {
			respondCreated(writer, ids)
		}
	} else {
		if id, err := controller.service.CreateOne(equipmentCreate); err != nil {
			respondBadRequest(writer, fmt.Errorf(equipmentNotCreated, err), nil)
		} else {
			respondCreated(writer, id)
		}
	}
}

func (controller *Equipment) UpdateByIds(writer http.ResponseWriter, request *http.Request) {
	if idSet, err := parseIds(mux.Vars(request)["id"]); err != nil {
		respondBadRequest(writer, fmt.Errorf(invalidEquipmentId, "update", err), nil)
	} else {
		var equipmentUpdate dtos.EquipmentUpdate
		if err := json.NewDecoder(request.Body).Decode(&equipmentUpdate); err != nil {
			respondBadRequest(writer, fmt.Errorf(invalidJSONData, err), nil)
		} else {
			switch ids := idSet.Slice(); len(ids) {
				case 0:
					respondBadRequest(writer, fmt.Errorf(parameterIsRequired, "id", "update"), nil)
				case 1: // Single id
					if err = controller.service.UpdateById(&equipmentUpdate, ids[0]); err != nil {
						respondNotFound(writer, fmt.Errorf(equipmentNotFound, "update", err), nil) // TODO: check other reasons
					} else {
						respond(writer, http.StatusNoContent, nil)
					}
				default:
					if updatedIds, err := controller.service.UpdateByIds(&equipmentUpdate, ids); err != nil {
						respondBadRequest(writer, fmt.Errorf(actionFailed, "update", err), ids)
					} else {
						respondOK(writer, splitSucceededAndFailed("updated", "unfound", idSet, updatedIds...))
					}
			}
		}
	}
}

func (controller *Equipment) UpdateByConditions(writer http.ResponseWriter, request *http.Request) {
	if equipmentFilter, err := dtos.EquipmentFilterFromRequest(request); err != nil {
		respondBadRequest(writer, fmt.Errorf(invalidGETParameters, "update", err), nil)
	} else {
		var equipmentUpdate dtos.EquipmentUpdate
		if err = json.NewDecoder(request.Body).Decode(&equipmentUpdate); err != nil {
			respondBadRequest(writer, fmt.Errorf(invalidJSONData, err), nil)
		} else if ids, err := controller.service.UpdateByConditions(&equipmentUpdate, equipmentFilter); err != nil {
			respondBadRequest(writer, fmt.Errorf(actionFailed, "update", err), nil) // TODO: BadRequest?
		} else {
			respondOK(writer, ids)
		}
	}
}

func (controller *Equipment) DeleteByIds(writer http.ResponseWriter, request *http.Request) {
	if idSet, err := parseIds(mux.Vars(request)["id"]); err != nil {
		respondBadRequest(writer, fmt.Errorf(invalidEquipmentId, "delete", err), nil)
	} else {
		switch ids := idSet.Slice(); len(ids) {
			case 0:
				respondBadRequest(writer, fmt.Errorf(parameterIsRequired, "id", "delete"), nil)
			case 1: // Single id
				if err = controller.service.DeleteById(ids[0]); err != nil {
					respondNotFound(writer, fmt.Errorf(equipmentNotFound, "delete", err), nil) // TODO: check other reasons
				} else {
					respond(writer, http.StatusNoContent, nil)
				}
			default: // Multiply ids
				if deletedIds, err := controller.service.DeleteByIds(ids); err != nil {
					respondBadRequest(writer, fmt.Errorf(actionFailed, "delete", err), ids)
				} else {
					respondOK(writer, splitSucceededAndFailed("deleted", "unfound", idSet, deletedIds...))
				}
		}
	}
}

func (controller *Equipment) DeleteByConditions(writer http.ResponseWriter, request *http.Request) {
	if equipmentFilter, err := dtos.EquipmentFilterFromRequest(request); err != nil {
		respondBadRequest(writer, fmt.Errorf(invalidGETParameters, "delete", err), nil)
	} else if ids, err := controller.service.DeleteByConditions(equipmentFilter); err != nil {
		respondBadRequest(writer, fmt.Errorf(actionFailed, "delete", err), nil) // TODO: BadRequest?
	} else {
		respondOK(writer, ids)
	}
}


func (controller *Equipment) FindById(writer http.ResponseWriter, request *http.Request) {
	if idSet, err := parseIds(mux.Vars(request)["id"]); err != nil {
		respondBadRequest(writer, fmt.Errorf(invalidEquipmentId, "delete", err), nil)
	} else {
		switch ids := idSet.Slice(); len(ids) {
			case 0:
				respondBadRequest(writer, fmt.Errorf(parameterIsRequired, "id", "search"), nil)
			case 1: // Single id
				if equipmentGet, err := controller.service.FindById(ids[0]); err != nil {
					respondNotFound(writer, fmt.Errorf(equipmentNotFound, "search", err), ids[0])
				} else {
					respondOK(writer, equipmentGet)
				}
			default: // Multiply ids
				if foundEquipment, err := controller.service.FindByIds(ids); err != nil {
					respondBadRequest(writer, fmt.Errorf(actionFailed, "search", err), ids)
				} else {
					response := make(map[string]interface{})
					if len(foundEquipment) > 0 {
						response["found"] = foundEquipment
						for _, equipment := range foundEquipment {
							idSet.Exclude(equipment.Id)
						}
					}
					if !idSet.IsEmpty() {
						response["unfound"] = idSet.Slice()
					}
					respondOK(writer, response)
				}
		}
	}
}

func (controller *Equipment) FindByConditions(writer http.ResponseWriter, request *http.Request) {
	if equipmentFilter, err := dtos.EquipmentFilterFromRequest(request); err != nil {
		respondBadRequest(writer, fmt.Errorf(invalidGETParameters, "search", err), nil)
	} else if equipments, err := controller.service.FindByConditions(equipmentFilter); err != nil {
		respondBadRequest(writer, fmt.Errorf(actionFailed, "search", err), nil) // TODO: BadRequest?
	} else {
		respondOK(writer, equipments)
	}
}
