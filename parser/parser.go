package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bamboV/kinopoisk"
	"strconv"
	"strings"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding"
)

const urlFormat = "https://www.kinopoisk.ru/user/%v/movies/list/type/%v/page/%v"

type Parser struct {
	User kinopoisk.User
}

func (p *Parser) ParseFolder(folderId int) ([]kinopoisk.Movie, error) {
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
	decoder := charmap.Windows1251.NewDecoder()

	document.Find("#itemList").Find("li").Each(func(i int, s *goquery.Selection) {
		attrId, _ := s.Attr("id")
		id, _ := strconv.Atoi(strings.Trim(attrId, "film_"))
		infoDiv := s.Find("div.info")
		translatedName := decodeWindows1251(fmt.Sprint(infoDiv.Find("a.name").Text()), *decoder)
		name := decodeWindows1251(infoDiv.Find("span").First().Text(), *decoder)
		movies = append(movies, kinopoisk.Movie{Id: id, Name: name, TranslatedName: translatedName})

	})

	return movies, nil
}

func decodeWindows1251(str string, decoder encoding.Decoder) (string) {
	res, _ := decoder.String(str)
	return res
}
