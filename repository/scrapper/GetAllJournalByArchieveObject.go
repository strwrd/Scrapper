package scrapper

import (
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/strwrd/scrapper/model"
)

// GetAllJournalByArchieveObject : return journals where belongs to archieve object
func (r *repository) GetAllJournalByArchieveObject(archieve *model.Archieve) ([]*model.Journal, error) {
	// Do scrapping html
	doc, err := r.scrapping(archieve.Link)
	if err != nil {
		return nil, err
	}

	// Create array (slice) of journal object
	journals := make([]*model.Journal, 0)

	// Scrapping html then parse to archieve object with looping (IIFE style<wrapper function>)
	if err := func() error {
		var err error

		// Finding DOM element then looping each childrens element
		doc.Find(".media-list").Children().Each(func(i int, s *goquery.Selection) {
			// initialization variable of journal
			var link string
			var pdfLink string
			var title string
			var abstract string
			var authors string
			var published string

			// Scrapping data from html
			link, linkOk := s.Find(".item-recent-title.media-heading").Children().First().Attr("href")
			pdfLink, pdfLinkOk := s.Find(".media-right.media-top").Children().First().Attr("href")
			title = s.Find(".item-recent-title.media-heading").Children().First().Text()
			authors = s.Find(".authors").Text()

			// if each link are not empty then push to array of journal
			if linkOk && pdfLinkOk {
				var detailDoc *goquery.Document

				// Get detail journal object (abstract, published) Do scrapping again
				detailDoc, err = r.scrapping(link)

				// Scrapping data from html
				published = detailDoc.Find(".date-published").After("strong").Text()
				abstract = detailDoc.Find(".article-abstract").Children().Text()

				// Creating journal object
				journal := &model.Journal{
					ID:         r.uuid.CreateV4(),
					ArchieveID: archieve.ID,
					Abstract:   strings.TrimSpace(abstract),
					Authors:    strings.TrimSpace(authors),
					Link:       strings.TrimSpace(link),
					PDFLink:    strings.TrimSpace(pdfLink),
					Published:  strings.TrimSpace(published),
					Title:      strings.TrimSpace(title),
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}

				// Push archieve object to archieve array (slice)
				journals = append(journals, journal)

				// Print log journal
				log.Printf("✓ %v\n", journal.Title)
			}
		})
		return err
	}(); err != nil {
		return nil, err
	}

	return journals, nil
}
