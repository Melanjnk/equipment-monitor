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

func (controller *Equipment) Close() {
	controller.service.Close()
}

func (controller *Equipment) Create(writer http.ResponseWriter, request *http.Request) {
	const action = actionCreate
	if equipmentCreate, equipmentCreateMany, err := FromRequestJSON(request); err != nil {
		respondBadRequest(writer, fmt.Errorf(invalidJSONData, err), nil)
	} else if equipmentCreateMany != nil {
		if ids, err := controller.service.CreateMany(equipmentCreateMany); err != nil {
			respondBadRequest(writer, fmt.Errorf(actionFailed, action, err), nil)
		} else {
			respondCreated(writer, ids)
		}
	} else {
		if id, err := controller.service.CreateOne(equipmentCreate); err != nil {
			respondBadRequest(writer, fmt.Errorf(actionFailed, action, err), nil)
		} else {
			respondCreated(writer, id)
		}
	}
}

func (controller *Equipment) UpdateByIds(writer http.ResponseWriter, request *http.Request) {
	const action = actionUpdate
	if idSet, err := parseIds(mux.Vars(request)[`id`]); err != nil {
		respondBadRequest(writer, fmt.Errorf(invalidEquipmentId, action, err), nil)
	} else {
		var equipmentUpdate dtos.EquipmentUpdate
		if err = json.NewDecoder(request.Body).Decode(&equipmentUpdate); err != nil {
			respondBadRequest(writer, fmt.Errorf(invalidJSONData, err), nil)
		} else {
			switch ids := idSet.Slice(); len(ids) {
				case 0:
					respondBadRequest(writer, fmt.Errorf(parameterIsRequired, "id", action), nil)
				case 1: // Single id
					var found bool
					if found, err = controller.service.UpdateById(&equipmentUpdate, dtos.EquipmentFilterFromIds(ids)); err != nil {
						respondInternalError(writer, fmt.Errorf(actionFailed, action, err), nil)
					} else if !found {
						respondNotFound(writer, fmt.Errorf(equipmentNotFound, ids[0], action), nil) // TODO: check other reasons
					} else {
						respondNoContent(writer)
					}
				default:
					if updatedIds, err := controller.service.UpdateByIds(&equipmentUpdate, dtos.EquipmentFilterFromIds(ids)); err != nil {
						respondInternalError(writer, fmt.Errorf(actionFailed, action, err), ids)
					} else if len(updatedIds) == 0 {
						respondNotFound(writer, fmt.Errorf(equipmentNotFound, ids, action), nil)
					} else {
						idSet.ExcludeMultiply(updatedIds...)
						if idSet.IsEmpty() {
							respondNoContent(writer)
						} else {
							respondMulti(writer, map[string]any{"updated": updatedIds, "unfound": idSet.Slice()})
						}
					}
			}
		}
	}
}

func (controller *Equipment) UpdateByConditions(writer http.ResponseWriter, request *http.Request) {
	const action = actionUpdate
	if equipmentFilter, err := dtos.EquipmentFilterFromRequest(request); err != nil {
		respondBadRequest(writer, fmt.Errorf(invalidGETParameters, action, err), nil)
	} else {
		var equipmentUpdate dtos.EquipmentUpdate
		if err = json.NewDecoder(request.Body).Decode(&equipmentUpdate); err != nil {
			respondBadRequest(writer, fmt.Errorf(invalidJSONData, err), nil)
		} else if ids, err := controller.service.UpdateByConditions(&equipmentUpdate, equipmentFilter); err != nil {
			respondInternalError(writer, fmt.Errorf(actionFailed, action, err), nil)
		} else {
			respondOK(writer, ids)
		}
	}
}

func (controller *Equipment) DeleteByIds(writer http.ResponseWriter, request *http.Request) {
	const action = actionDelete
	if idSet, err := parseIds(mux.Vars(request)[`id`]); err != nil {
		respondBadRequest(writer, fmt.Errorf(invalidEquipmentId, action, err), nil)
	} else {
		switch ids := idSet.Slice(); len(ids) {
			case 0:
				respondBadRequest(writer, fmt.Errorf(parameterIsRequired, "id", action), nil)
			case 1: // Single id
				var found bool
				if found, err = controller.service.DeleteById(dtos.EquipmentFilterFromIds(ids)); err != nil {
					respondInternalError(writer, fmt.Errorf(actionFailed, action, err), nil)
				} else if !found {
					respondNotFound(writer, fmt.Errorf(equipmentNotFound, ids[0], action), nil) // TODO: check other reasons
				} else {
					respondNoContent(writer)
				}
			default: // Multiply ids
				if deletedIds, err := controller.service.DeleteByIds(dtos.EquipmentFilterFromIds(ids)); err != nil {
					respondInternalError(writer, fmt.Errorf(actionFailed, action, err), ids)
				} else if len(deletedIds) == 0 {
					respondNotFound(writer, fmt.Errorf(equipmentNotFound, ids, action), nil)
				} else {
					idSet.ExcludeMultiply(deletedIds...)
					if idSet.IsEmpty() {
						respondNoContent(writer)
					} else {
						respondMulti(writer, map[string]any{"deleted": deletedIds, "unfound": idSet.Slice()})
					}
				}
		}
	}
}

func (controller *Equipment) DeleteByConditions(writer http.ResponseWriter, request *http.Request) {
	const action = actionDelete
	if equipmentFilter, err := dtos.EquipmentFilterFromRequest(request); err != nil {
		respondBadRequest(writer, fmt.Errorf(invalidGETParameters, action, err), nil)
	} else if ids, err := controller.service.DeleteByConditions(equipmentFilter); err != nil {
		respondInternalError(writer, fmt.Errorf(actionFailed, action, err), nil)
	} else {
		respondOK(writer, ids)
	}
}


func (controller *Equipment) FindById(writer http.ResponseWriter, request *http.Request) {
	const action = actionSearch
	if idSet, err := parseIds(mux.Vars(request)[`id`]); err != nil {
		respondBadRequest(writer, fmt.Errorf(invalidEquipmentId, action, err), nil)
	} else {
		switch ids := idSet.Slice(); len(ids) {
			case 0:
				respondBadRequest(writer, fmt.Errorf(parameterIsRequired, "id", action), nil)
			case 1: // Single id
				if equipmentGet, err := controller.service.FindById(dtos.EquipmentFilterFromIds(ids), ); err != nil {
					respondNotFound(writer, fmt.Errorf(equipmentNotFound, ids[0], action), nil)
				} else {
					respondOK(writer, equipmentGet)
				}
			default: // Multiply ids
				if foundEquipment, err := controller.service.FindByIds(dtos.EquipmentFilterFromIds(ids)); err != nil {
					respondInternalError(writer, fmt.Errorf(actionFailed, action, err), ids)
				} else if len(foundEquipment) == 0 {
					respondNotFound(writer, fmt.Errorf(equipmentNotFound, ids, action), nil)
				} else {
					for _, equipment := range foundEquipment {
						idSet.Exclude(equipment.Id)
					}
					if idSet.IsEmpty() {
						respondOK(writer, foundEquipment)
					} else {
						respondMulti(writer, map[string]any{"found": foundEquipment, "unfound": idSet.Slice()})
					}
				}
		}
	}
}

func (controller *Equipment) FindByConditions(writer http.ResponseWriter, request *http.Request) {
	const action = actionSearch
	if equipmentFilter, err := dtos.EquipmentFilterFromRequest(request); err != nil {
		respondBadRequest(writer, fmt.Errorf(invalidGETParameters, action, err), nil)
	} else if equipments, err := controller.service.FindByConditions(equipmentFilter); err != nil {
		respondInternalError(writer, fmt.Errorf(actionFailed, action, err), nil)
	} else {
		respondOK(writer, equipments)
	}
}
