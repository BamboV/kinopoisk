package parser

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"github.com/bamboV/kinopoisk"
)

const urlFormat = "https://www.kinopoisk.ru/user/%v/movies/list/type/%v/page/%v"

type Parser struct {
	User kinopoisk.User
}

func (p *Parser)ParseFolder(folderId int) ([]kinopoisk.Movie, error) {
	allMovies := []kinopoisk.Movie{}
	page := 1
	for {
		movies, err := p.parseFolderPage(folderId, page)

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

func (p *Parser) parseFolderPage(folderId int, page int) ([]kinopoisk.Movie, error) {
	url := fmt.Sprintf(urlFormat, p.User.Id, folderId, page)

	movies := []kinopoisk.Movie{}
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
		movies = append(movies, kinopoisk.Movie{Id:id, Name:name, TranslatedName:translatedName})

	})

	return movies, nil
}