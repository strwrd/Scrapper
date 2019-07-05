package scrapper

import (
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/strwrd/scrapper/model"
)

// GetAllArchieve : return all archieve from jptiik
func (r *repository) GetAllArchieve() ([]*model.Archieve, error) {
	// Do scrapping html
	doc, err := r.scrapping(_RepoLink)
	if err != nil {
		return nil, err
	}

	// Create array (slice) of archieve object
	archieves := make([]*model.Archieve, 0)

	// Scrapping html then parse to archieve object with looping
	if err := func() error {
		var err error

		// Finding DOM element then looping each childrens element
		doc.Find(".issues.media-list").Children().Each(func(i int, s *goquery.Selection) {
			var detailArchieveDoc *goquery.Document
			// get the band and title
			link, ok := s.Find(".title").Attr("href")
			code := s.Find(".series.lead").Text()

			// Scrapping detail archieve
			detailArchieveDoc, err = r.scrapping(link)

			// If archieve link is not empty then push archieve to archieve array (slice).
			// If archieve link is empty then skip
			if ok {
				// Creating archieve object
				archieve := &model.Archieve{
					ID:        r.uuid.CreateV4(),
					Code:      strings.TrimSpace(code),
					Link:      strings.TrimSpace(link),
					Published: strings.TrimSpace(detailArchieveDoc.Find(".published").Children().Children().Text()),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				// Push archieve object to archieve array (slice)
				archieves = append(archieves, archieve)

				// Print log archieve
				log.Printf("✓ %v\n", archieve.Code)
			}
		})

		return err
	}(); err != nil {
		return nil, err
	}

	return archieves, nil
}
