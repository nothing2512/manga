package controllers

import (
	"fmt"
	"io"
	"main/models"
	"net/http"

	"gorm.io/gorm"
)

type InDb struct {
	pg *gorm.DB
}

func NewInstance() *InDb {
	return &InDb{
		pg: models.GetConnection(),
	}
}

func (idb *InDb) write(w http.ResponseWriter, message string) {
	n, err := fmt.Fprintf(w, "data: %v\n\n", message)
	fmt.Println(n, message)
	if err != nil {
		fmt.Println(err.Error())
	}
	w.(http.Flusher).Flush()
}

func (idb *InDb) scrap(uri string) (string, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
