package exchange

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kfukue/lyle-labs/libraries/utils"
)

const stepAssetsPath = "stepAssets"

func handleStepAssets(r *mux.Router, path string) {
	r.HandleFunc(fmt.Sprintf("/%s", path), getStepAssetsHandle).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s", path), createStepAssetHandle).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/{id}", path), getStepAssetHandle).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/{id}", path), putStepAssetHandle).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/{id}", path), deleteStepAsset).Methods("DELETE")

}

func getStepAssetsHandle(w http.ResponseWriter, r *http.Request) {
	stepAssets, err := getStepAssets()
	if err != nil {
		log.Fatal(err)
		utils.RespondWithJSON(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	results := make(map[string][]StepAsset)

	results["stepAssets"] = stepAssets
	utils.RespondWithJSON(w, http.StatusOK, results)
}

func createStepAssetHandle(w http.ResponseWriter, r *http.Request) {
	var stepAsset StepAsset
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&stepAsset); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	newInt, err := insertStepAsset(stepAsset)
	if err != nil {
		log.Fatal(err)
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	stepAsset.ID = &newInt
	results := make(map[string]StepAsset)
	results["stepAsset"] = stepAsset
	utils.RespondWithJSON(w, http.StatusCreated, results)
}

func getStepAssetHandle(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid stepAsset ID")
		return
	}
	stepAsset, err := getStepAsset(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	if stepAsset == nil {
		utils.RespondWithError(w, http.StatusNotFound, "StepAsset not found")
		return
	}
	results := make(map[string]StepAsset)
	results["stepAsset"] = *stepAsset
	utils.RespondWithJSON(w, http.StatusOK, results)
}

func putStepAssetHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid stepAsset ID")
		return
	}

	var stepAsset StepAsset
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&stepAsset); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	stepAsset.ID = &id

	err = updateStepAsset(stepAsset)
	if err != nil {
		log.Print(err)
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	results := make(map[string]StepAsset)
	results["stepAsset"] = stepAsset
	utils.RespondWithJSON(w, http.StatusOK, results)
	return
}

func deleteStepAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid StepAsset ID")
		return
	}
	deleteErr := removeStepAsset(id)
	if deleteErr != nil {
		log.Fatal(deleteErr)
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// SetupRoutes :
func SetupRoutes(r *mux.Router) {
	handleStepAssets(r, fmt.Sprintf("%s", stepAssetsPath))
}
