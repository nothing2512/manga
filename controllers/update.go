package controllers

import (
	"fmt"
	"main/models"
	"net/http"
	"slices"
	"strings"
	"time"
)

func (idb *InDb) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	idb.write(w, "Getting All Mangas")
	var mangas []models.Manga
	idb.pg.Find(&mangas)

	idb.write(w, "Updating Mangas")
	total := len(mangas)
	for k, m := range mangas {
		data, err := idb.scrap(m.Link)
		if err != nil {
			idb.write(w, err.Error())
			idb.write(w, "EOF")
			return
		}
		idb.updateChapter(w, data, m, fmt.Sprintf("[%v/%v]", k+1, total))
	}
	idb.write(w, "Done")
	idb.write(w, "EOF")
	return
}

func (idb *InDb) updateChapter(w http.ResponseWriter, data string, manga models.Manga, counter string) error {
	idb.write(w, counter+" Updating "+manga.Name)
	// Getting list chapters
	data = strings.Split(strings.Split(data, "id=\"chapter_list\"")[1], "</div>")[0]
	chapters := strings.Split(data, "<span class=\"lchx\">")[1:]
	slices.Reverse(chapters)

	if len(chapters) <= manga.Index {
		return nil
	}

	total := len(chapters) - manga.Index

	for k, ch := range chapters[manga.Index:] {
		idb.write(w, fmt.Sprintf("%v Updating %v [%v/%v]", counter, manga.Name, k+1, total))
		chapter := models.Chapter{
			MangaId: manga.ID,
			Link:    strings.Split(strings.Split(ch, "href=\"")[1], "\"")[0],
		}
		chData, err := idb.scrap(chapter.Link)
		if err != nil {
			return err
		}
		imagesData := strings.Split(strings.Split(strings.Split(chData, "<div id=\"chimg-auh\">")[1], "</div>")[0], "<img src=\"")[1:]
		var images []models.Image
		if len(imagesData) == 0 {
			return nil
		}
		idb.pg.Create(&chapter)
		for _, im := range imagesData {
			images = append(images, models.Image{
				ChapterId: chapter.ID,
				Link:      strings.Split(im, "\"")[0],
			})
		}
		idb.pg.Create(&images)
		manga.Index += 1
		manga.LastUpdated = time.Now()
		idb.pg.Save(&manga)
	}
	return nil
}
