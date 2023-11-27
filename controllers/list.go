package controllers

import (
	"encoding/json"
	"main/models"
	"net/http"
)

func (idb *InDb) List(w http.ResponseWriter, r *http.Request) {
	var mangas []models.Manga
	idb.pg.
		Order("last_updated DESC").
		Find(&mangas)
	_ = json.NewEncoder(w).Encode(mangas)
}
