package controller

import (
	"net/http"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/model"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/service"
)

type Equipment struct {
	service service.Equipment
}

func NewEquipment(service service.Equipment) Equipment {
	return Equipment{service: service}
}

func (eqc *Equipment) List(w http.ResponseWriter, r *http.Request) {
	eql, err := eqc.service.List()
	if err != nil {
		writeMessage(w, http.StatusInternalServerError, "List error: %v", err)
	} else {
		writeJSON(w, http.StatusOK, eql)
	}
}

func (eqc *Equipment) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: batch create
	var j map[string]interface{}
	if err := readJSON(r, j); err != nil {
		writeMessage(w, http.StatusBadRequest, "Invalid JSON data: %v", err)
		return
	}
	eqType, ok := j["type"]
	if !ok {
		writeMessage(w, http.StatusBadRequest, parameterIsRequired, "type")
		return
	}
	eqt := model.ParseEquipmentType(eqType.(string))
	if eqt == nil {
		writeMessage(w, http.StatusBadRequest, "Invalid equipment type: `%s`", eqType)
		return
	}
	eqParams, ok := j["parameters"]
	if !ok {
		writeMessage(w, http.StatusBadRequest, parameterIsRequired, "parameters")
		return
	}
	id, err := eqc.service.Create(*eqt, eqParams.(model.Params))
	if err != nil {
		writeMessage(w, http.StatusBadRequest, "Create equipment error: %v", err)
		return
	}

	writeMessage(w, http.StatusCreated, "Equipment %s is created", id)
}

func (eqc *Equipment) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: batch update
	vars := mux.Vars(r)
	eqId, ok := vars["id"]
	if !ok {
		writeMessage(w, http.StatusBadRequest, parameterIsRequired, "id")
		return
	}
	id, err := uuid.FromString(eqId)
	if err != nil {
		writeMessage(w, http.StatusBadRequest, "Invalid UUID: `%s`", eqId)
	}

	var j map[string]interface{}
	if err := readJSON(r, j); err != nil {
		writeMessage(w, http.StatusBadRequest, "Invalid JSON data: %v", err)
		return
	}
	
	eqStatus, ok := j["status"]
	var eqos *model.OperationalStatus
	if ok {
		eqos = model.ParseOperationalStatus(eqStatus.(string))
		if eqos == nil {
			writeMessage(w, http.StatusBadRequest, "Invalid equipment operational status: `%s`", eqStatus)
			return
		}
	} else {
		eqos = nil
	}

	eqParams, ok := j["parameters"]
	var eqp model.Params
	if ok {
		eqp = eqParams.(model.Params)
	} else {
		eqp = nil
	}

	err = eqc.service.Update(id, eqos, eqp)
	if err != nil {
		writeMessage(w, http.StatusBadRequest, "Update equipment `%s` error: %v", id, err)
		return
	}

	writeMessage(w, http.StatusOK, "Equipment %s is updated", id)
}

func (eqc *Equipment) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		writeMessage(w, http.StatusBadRequest, parameterIsRequired, "id")
		return
	}
	eq, err := eqc.service.Get(id)
	if err != nil {
		writeMessage(w, http.StatusNotFound, "Get equipment `%s` error: %v", id, err)
		return
	}

	writeJSON(w, http.StatusOK, eq)
}

func (eqc *Equipment) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		writeMessage(w, http.StatusBadRequest, parameterIsRequired, "id")
		return
	}

	deleted, err := eqc.service.Delete(id)
	if err != nil {
		writeMessage(w, http.StatusBadRequest, "Delete equipment `%s` error: %v", id, err)
		return
	}

	if !deleted {
		writeMessage(w, http.StatusNotFound, "Unable to find equipment `%s` for deleting", id)
		return
	}

	writeMessage(w, http.StatusOK, "Equipment %s was deleted", id)
}
