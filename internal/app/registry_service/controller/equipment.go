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

func (cont *Equipment) List(writer http.ResponseWriter, request *http.Request) {
	// TODO: Filter
	eql, err := cont.service.List()
	if err != nil {
		writeMessage(writer, http.StatusInternalServerError, "List error: %v", err)
	} else {
		writeJSON(writer, http.StatusOK, eql)
	}
}

func (cont *Equipment) Create(writer http.ResponseWriter, request *http.Request) {
	// TODO: batch create
	
	if eqc, err := dtos.FromRequestJSON[dtos.EquipmentCreate](request); err != nil {
		writeMessage(writer, http.StatusBadRequest, invalidJSONData, err)
	} else if id, err := cont.service.Create(eqc); err != nil {
		writeMessage(writer, http.StatusBadRequest, "Create equipment error: %v", err)
	} else {
		writeMessage(writer, http.StatusCreated, equipmentActionIsPerformed, id, "created")
	}
}

func (cont *Equipment) Update(writer http.ResponseWriter, request *http.Request) {
	// TODO: batch update
	if equ, err := dtos.FromRequestJSON[dtos.EquipmentUpdate](request); err != nil {
		writeMessage(writer, http.StatusBadRequest, invalidJSONData, err)
	} else if updated, err := cont.service.Update(equ); err != nil {
		writeMessage(writer, http.StatusBadRequest, equipmentIdError, "Update", equ.Id, err)
	} else if !updated {
		writeMessage(writer, http.StatusNotFound, unableToFindEquipment, equ.Id, "updating")
	} else {
		writeMessage(writer, http.StatusOK, equipmentActionIsPerformed, equ.Id, "updated")
	}
}

func (cont *Equipment) Get(writer http.ResponseWriter, request *http.Request) {
	if id, ok := mux.Vars(request)["id"]; !ok {
		writeMessage(writer, http.StatusBadRequest, parameterIsRequired, "id")
	} else if eqg, err := cont.service.Get(id); err != nil {
		writeMessage(writer, http.StatusNotFound, equipmentIdError, "Get", id, err)
	} else {
		writeJSON(writer, http.StatusOK, eqg)
	}
}

func (cont *Equipment) Delete(writer http.ResponseWriter, request *http.Request) {
	if id, ok := mux.Vars(request)["id"]; !ok {
		writeMessage(writer, http.StatusBadRequest, parameterIsRequired, "id")
	} else if deleted, err := cont.service.Delete(id); err != nil {
		writeMessage(writer, http.StatusBadRequest, equipmentIdError, "Delete", id, err)
	} else if !deleted {
		writeMessage(writer, http.StatusNotFound, unableToFindEquipment, id, "deleting")
	} else {
		writeMessage(writer, http.StatusOK, equipmentActionIsPerformed, id, "deleted")
	}
}
