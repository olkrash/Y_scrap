package main

import (
	"Y_scrap/models"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const page = "https://www.yakaboo.ua/knigi/vospitanie-detej-knigi-dlja-roditelej.html?"

type scraper struct {
	cfg Config
	db  *gorm.DB
}

func (s *scraper) collectData() error {
	//check folder exists pictures
	if _, err := os.Stat(s.cfg.PictureFolder); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(s.cfg.PictureFolder, 0700)
			if err != nil {
				log.Println(err)
			}
		}
	}
	c := colly.NewCollector()

	//Visiting all links at each page
	c.OnHTML(".thumbnails_middle .content>a", func(e *colly.HTMLElement) {
		println(e.Attr("href"))
		err := e.Request.Visit(e.Attr("href"))
		if err != nil {
			return
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting the webpage", r.URL)
	})

	//Find and visit next page links
	c.OnHTML("#product-attribute-specs-table", func(e *colly.HTMLElement) {
		name := e.DOM.Closest("body").Find("#product-title h1").Text()
		name = strings.TrimPrefix(name, "Книга ")

		author := strings.TrimSpace(e.DOM.Find("tbody tr:first-child>td:last-child").Text())

		codeElement := e.DOM.Closest("body").Find(".product-sku span:nth-last-of-type(1)")
		code := strings.TrimSuffix(codeElement.Text(), "|")

		book := models.Book{
			BookName:   name,
			Author:     author,
			OriginalID: code,
		}

		s.db.Create(&book)
		//fmt.Printf("Visiting the name = %s\n author %s\n and the code = %s\n", name, author, code)

		//e.ForEach("tr:not(tr:first-child)", func(_ int, e *colly.HTMLElement) {
		//	charName := strings.TrimSpace(e.DOM.Find("td:first-child").Text())
		//	charValue := strings.TrimSpace(e.DOM.Find("td:nth-child(2)").Text())
		//
		//	s.db.Create(&models.BookChar{
		//		BookID:   book.ID,
		//		BookChar: charName,
		//		BookVal:  charValue,
		//	})
		//})

		//check if the record exists in DB
		//var oldBooks models.Book
		//if err := s.db.Where("name = ?", name).First(&oldBooks).Error; err == nil {
		//	return
		//}

		//check pictures folder exist, if not created, create
		path := s.cfg.PictureFolder + "/" + code
		err := os.Mkdir(path, 0700)
		if err != nil {
			log.Println(err)
		}

		e.DOM.Closest("body").Find("#media_popup_photos img").Each(func(_ int, s *goquery.Selection) {
			thumbSrc, _ := s.Attr("data-original")
			thumbSplit := strings.Split(thumbSrc, "/")
			thumbPic := thumbSplit[len(thumbSplit)-1]
			thumbPic = thumbPic[:len(thumbPic)-4]

			err := downloadFile(thumbSrc, path, thumbPic)
			if err != nil {
				log.Fatal(err)
			}
		})

	})

	//pagination
	c.OnHTML("li:nth-child(9) > a.next", func(e *colly.HTMLElement) {
		println("Next page link found:", e.Attr("href"))
		err := e.Request.Visit(e.Attr("href"))
		if err != nil {
			return
		}
	})

	return c.Visit("https://www.yakaboo.ua/knigi/vospitanie-detej-knigi-dlja-roditelej.html?")
}

func downloadFile(url, path, picID string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("received non 200 response code")
	}

	file, err := os.Create(path + "/" + picID + ".jpg")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
