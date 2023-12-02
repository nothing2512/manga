package controllers

import (
	"errors"
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
	totalUpdated := 0
	var errs []error
	for k, m := range mangas {
		data, err := idb.scrap(m.Link)
		if err != nil {
			idb.write(w, err.Error())
			idb.write(w, "EOF")
			return
		}
		err, updated := idb.updateChapter(w, data, m, fmt.Sprintf("[%v/%v]", k+1, total), false)
		totalUpdated += updated
		if err != nil {
			errs = append(errs, err)
		}
	}
	idb.write(w, "")
	idb.write(w, "Error : ")
	for _, e := range errs {
		idb.write(w, e.Error())
	}
	idb.write(w, fmt.Sprintf("Updated : %v", totalUpdated))
	idb.write(w, "Done")
	idb.write(w, "EOF")
	return
}

func (idb *InDb) updateChapter(w http.ResponseWriter, data string, manga models.Manga, counter string, cached bool) (error, int) {
	idb.write(w, counter+" Updating "+manga.Name)
	// Getting list chapters
	_data := strings.Split(data, "id=\"chapter_list\"")
	if len(_data) <= 1 {
		if cached {
			return errors.New(manga.Name), 0
		}
		d, err := idb.scrap(fmt.Sprintf("https://webcache.googleusercontent.com/search?q=cache:%v", manga.Link))
		if err != nil {
			return errors.New(manga.Name), 0
		}
		return idb.updateChapter(w, d, manga, counter, true)
	}
	data = strings.Split(_data[1], "</div>")[0]
	chapters := strings.Split(data, "<span class=\"lchx\">")[1:]
	slices.Reverse(chapters)

	if len(chapters) <= manga.Index {
		return nil, 0
	}

	total := len(chapters) - manga.Index
	updated := 0

	for k, ch := range chapters[manga.Index:] {
		idb.write(w, fmt.Sprintf("%v Updating %v [%v/%v]", counter, manga.Name, k+1, total))
		chapter := models.Chapter{
			MangaId: manga.ID,
			Link:    strings.Split(strings.Split(ch, "href=\"")[1], "\"")[0],
		}
		chData, err := idb.scrap(chapter.Link)
		if err != nil {
			return errors.New(manga.Name), 0
		}
		if strings.Contains(chData, "Just a moment...") {
			chData, err = idb.scrap(fmt.Sprintf("https://webcache.googleusercontent.com/search?q=cache:%v", chapter.Link))
			if err != nil {
				return errors.New(manga.Name), 0
			}
		}
		if strings.Contains(chData, "Just a moment...") {
			return errors.New(manga.Name), 0
		}
		imagesData := strings.Split(strings.Split(strings.Split(chData, "<div id=\"chimg-auh\">")[1], "</div>")[0], "<img src=\"")[1:]
		var images []models.Image
		if len(imagesData) == 0 {
			return errors.New(manga.Name), 0
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
		updated += 1
		manga.LastUpdated = time.Now()
		idb.pg.Save(&manga)
	}
	return nil, updated
}
