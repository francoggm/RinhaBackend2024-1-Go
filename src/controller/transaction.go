package controller

import (
	"crebito/database"
	"crebito/models"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func HandleTransaction(w http.ResponseWriter, r *http.Request, s neo4j.SessionWithContext) {
	defer r.Body.Close()

	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if id < 1 || id > 5 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var req models.TransactionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if req.Value < 1 || (req.Type != "d" && req.Type != "c") || (len(req.Description) < 1 || len(req.Description) > 10) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if req.Type == "d" {
		req.Value = -1 * req.Value
	}

	result, err := database.ExecuteTransaction(s, id, req)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrUserNotFound):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusUnprocessableEntity)
		}

		return
	}

	res, _ := json.Marshal(result)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
