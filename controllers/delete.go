package controllers

import (
	"main/models"
	"net/http"
)

func (idb *InDb) Delete(w http.ResponseWriter, r *http.Request) {
	var manga models.Manga
	idb.pg.Find(&manga, "id = ?", r.URL.Query().Get("id"))

	var chapters []models.Chapter
	idb.pg.Find(&chapters, "manga_id = ?", manga.ID)

	var chIds []int
	for _, c := range chapters {
		chIds = append(chIds, c.ID)
	}

	idb.pg.Delete(&models.Image{}, "chapter_id in (?)", chIds)
	idb.pg.Delete(&models.Chapter{}, "manga_id = ?", manga.ID)
	idb.pg.Delete(&models.Manga{}, "id = ?", manga.ID)

	http.Redirect(w, r, "/list", http.StatusTemporaryRedirect)
}
