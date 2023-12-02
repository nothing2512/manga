package controllers

import (
	"main/models"
	"net/http"
	"strings"
	"time"
)

func (idb *InDb) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var manga models.Manga

	manga.Link = r.URL.Query().Get("uri")

	idb.pg.First(&manga, "link = ?", manga.Link)
	if manga.ID != 0 {
		idb.write(w, "manga exist")
		idb.write(w, "EOF")
		return
	}

	manga.LastUpdated = time.Now()

	idb.write(w, "Getting manga")

	data, err := idb.scrap(manga.Link)
	if err != nil {
		idb.write(w, err.Error())
		idb.write(w, "EOF")
		return
	}

	idb.write(w, "Getting Image")
	manga.Image = strings.Split(strings.Split(strings.Split(data, "<div class=\"thumb\"")[1], "src=\"")[1], "\"")[0]
	idb.write(w, "image : "+manga.Image)

	idb.write(w, "Getting Name")
	manga.Name = strings.ReplaceAll(strings.Split(strings.Split(strings.Split(data, "<h1 class=\"entry-title\"")[1], ">")[1], "<")[0], "\n", "")
	if manga.Name[:5] == "Komik" {
		manga.Name = manga.Name[5:]
	}
	idb.write(w, "name : "+manga.Name)

	idb.write(w, "Saving Manga")
	idb.pg.Create(&manga)
	idb.write(w, "Manga Saved")

	err, _ = idb.updateChapter(w, data, manga, "[1/1]", false)
	if err != nil {
		idb.write(w, err.Error())
		idb.write(w, "EOF")
		return
	}
	idb.write(w, "Done")
	idb.write(w, "EOF")
	return
}
