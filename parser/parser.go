package parser

import (
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"strconv"
	"strings"
)

const urlFormat = "https://www.kinopoisk.ru/user/%v/movies/list/type/%v/page/%v"

func (k *Kinopoisk)ParseFolder(folderId int) ([]Movie, error) {
	allMovies := []Movie{}
	page := 1
	for {
		movies, err := k.parseFolderPage(folderId, page)

		if err != nil {
			return nil, err
		}

		if len(movies) <= 0 {
			return allMovies, nil
		}

		allMovies = append(allMovies, movies...)
		page++
	}

}

func (k *Kinopoisk) parseFolderPage(folderId int, page int) ([]Movie, error) {
	url := fmt.Sprintf(urlFormat, k.UserId, folderId, page)

	movies := []Movie{}
	document, err := goquery.NewDocument(url)

	if err != nil {
		return nil, err
	}

	document.Find("#itemList").Find("li").Each(func(i int, s *goquery.Selection) {
		attrId, _ := s.Attr("id")

		id, _ := strconv.Atoi(strings.Trim(attrId, "film_"))
		infoDiv := s.Find("div.info")
		translatedName := infoDiv.Find("a.name").Text()
		name := infoDiv.Find("span").First().Text()
		movies = append(movies, Movie{Id:id, Name:name, TranslatedName:translatedName})
	})

	return movies, nil
}