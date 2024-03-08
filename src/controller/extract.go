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

func HandleExtract(w http.ResponseWriter, r *http.Request, s neo4j.SessionWithContext) {
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

	result, err := database.GetExtract(s, id)

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
