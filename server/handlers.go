package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	vbox "github.com/blacktop/vm-proxy/drivers/virtualbox"
	"github.com/gorilla/mux"
	"github.com/riobard/go-virtualbox"
)

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

// Index root route
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

// VBoxList route lists all VMs
func VBoxList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	machines, err := virtualbox.ListMachines()
	assert(err)
	for _, machine := range machines {
		fmt.Println(machine.Name)
	}

	if err := json.NewEncoder(w).Encode(machines); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}

// VBoxStatus displays the machine readable status of a VM
func VBoxStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	nameOrID := vars["nameOrID"]

	machine, err := virtualbox.GetMachine(nameOrID)
	assert(err)
	if err := json.NewEncoder(w).Encode(machine.State); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}

// VBoxStart router starts a VM
func VBoxStart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	vars := mux.Vars(r)
	nameOrID := vars["nameOrID"]

	machine, err := virtualbox.GetMachine(nameOrID)
	assert(err)
	assert(machine.Start())

	w.WriteHeader(http.StatusOK)
}

// VBoxStop router stops a VM
func VBoxStop(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	vars := mux.Vars(r)
	nameOrID := vars["nameOrID"]

	d := vbox.NewDriver(nameOrID, "")
	err := d.Stop()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
}

// VBoxSnapshotRestore restores a certain snapshot
func VBoxSnapshotRestore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	vars := mux.Vars(r)
	nameOrID := vars["nameOrID"]
	snapShot := vars["snapShot"]

	d := vbox.NewDriver(nameOrID, "")
	outPut, err := d.RestoreSnapshot(nameOrID, snapShot)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(outPut))
}

// VBoxSnapshotRestoreCurrent restores the most resent snapshot
func VBoxSnapshotRestoreCurrent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	vars := mux.Vars(r)
	nameOrID := vars["nameOrID"]

	d := vbox.NewDriver(nameOrID, "")
	outPut, err := d.RestoreCurrentSnapshot(nameOrID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(outPut))
}
