package controller

import (
	"context"
	"crebito/database"
	"crebito/models"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func HandleExtract(w http.ResponseWriter, r *http.Request, s neo4j.SessionWithContext) {
	defer r.Body.Close()

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if id < 1 || id > 5 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	c := context.Background()

	result, err := s.ExecuteRead(c,
		func(tx neo4j.ManagedTransaction) (any, error) {
			result, err := tx.Run(c, database.ExtractQuery,
				map[string]any{
					"id": id,
				})
			if err != nil {
				return nil, err
			}

			record, err := result.Single(c)
			if err != nil {
				return nil, models.ErrUserNotFound
			}

			var res models.ExtractResponse

			res.UserInfo.Balance = record.AsMap()["saldo"].(int64)
			res.UserInfo.Limit = record.AsMap()["limite"].(int64)
			res.UserInfo.Date = time.Now()

			values := record.AsMap()["transacoes"].([]any)
			for _, v := range values {
				var t models.Transaction

				transaction := v.(map[string]any)

				t.Value = transaction["valor"].(int64)
				t.Type = transaction["tipo"].(string)
				t.Description = transaction["descricao"].(string)

				date := transaction["data"].(int64)
				t.Date = time.UnixMilli(date)

				res.Transactions = append(res.Transactions, t)
			}

			return res, nil
		})

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
