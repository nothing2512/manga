package controllers

import (
	"encoding/json"
	"main/models"
	"net/http"
	"strconv"
)

func (idb *InDb) Read(w http.ResponseWriter, r *http.Request) {
	var manga models.Manga
	idb.pg.First(&manga, "id = ?", r.URL.Query().Get("id"))

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	var chapter models.Chapter
	idb.pg.
		Order("id").
		Offset(page-1).
		First(&chapter, "manga_id = ?", manga.ID)

	manga.LastRead = page
	idb.pg.Save(&manga)

	var images []models.Image
	idb.pg.Find(&images, "chapter_id = ?", chapter.ID)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"manga":       manga,
		"images":      images,
		"canNext":     page < manga.Index,
		"canPrevious": page > 1,
		"page":        page,
	})
}
