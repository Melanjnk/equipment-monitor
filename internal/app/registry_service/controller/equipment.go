package controller

import (
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

func (controller *Equipment) List(writer http.ResponseWriter, request *http.Request) {
	// TODO: Filter
	equipmentList, err := controller.service.List()
	if err != nil {
		writeMessage(writer, http.StatusInternalServerError, "List error: %v", err)
	} else {
		writeJSON(writer, http.StatusOK, equipmentList)
	}
}

func (controller *Equipment) Create(writer http.ResponseWriter, request *http.Request) {
	if equipmentCreate, err := dtos.FromRequestJSON[dtos.EquipmentCreate](request); err != nil {
		writeMessage(writer, http.StatusBadRequest, invalidJSONData, err)
	} else if id, err := controller.service.Create(equipmentCreate); err != nil {
		writeMessage(writer, http.StatusBadRequest, "Create equipment error: %v", err)
	} else {
		writeMessage(writer, http.StatusCreated, equipmentActionIsPerformed, id, "created")
	}
}

func (controller *Equipment) Update(writer http.ResponseWriter, request *http.Request) {
	if equipmentUpdate, err := dtos.FromRequestJSON[dtos.EquipmentUpdate](request); err != nil {
		writeMessage(writer, http.StatusBadRequest, invalidJSONData, err)
	} else if updated, err := controller.service.Update(equipmentUpdate); err != nil {
		writeMessage(writer, http.StatusBadRequest, equipmentIdError, "Update", equipmentUpdate.Id, err)
	} else if !updated {
		writeMessage(writer, http.StatusNotFound, unableToFindEquipment, equipmentUpdate.Id, "updating")
	} else {
		writeMessage(writer, http.StatusOK, equipmentActionIsPerformed, equipmentUpdate.Id, "updated")
	}
}

func (controller *Equipment) Get(writer http.ResponseWriter, request *http.Request) {
	if id, ok := mux.Vars(request)["id"]; !ok {
		writeMessage(writer, http.StatusBadRequest, parameterIsRequired, "id")
	} else if eqg, err := controller.service.Get(id); err != nil {
		writeMessage(writer, http.StatusNotFound, equipmentIdError, "Get", id, err)
	} else {
		writeJSON(writer, http.StatusOK, eqg)
	}
}

func (controller *Equipment) Delete(writer http.ResponseWriter, request *http.Request) {
	if id, ok := mux.Vars(request)["id"]; !ok {
		writeMessage(writer, http.StatusBadRequest, parameterIsRequired, "id")
	} else if deleted, err := controller.service.Delete(id); err != nil {
		writeMessage(writer, http.StatusBadRequest, equipmentIdError, "Delete", id, err)
	} else if !deleted {
		writeMessage(writer, http.StatusNotFound, unableToFindEquipment, id, "deleting")
	} else {
		writeMessage(writer, http.StatusOK, equipmentActionIsPerformed, id, "deleted")
	}
}
