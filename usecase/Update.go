package usecase

import (
	"context"
	"log"
	"time"

	"github.com/strwrd/scrapper/model"
)

// Creating channel buffer for pararelism
type channel struct {
	archieve *model.Archieve
	err      error
}

// Update data between JPTIIK dan DB
func (u *usecase) Update() error {
	// Time execution
	start := time.Now()

	// Creating context with timeout duration process
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()

	// Get all archieve from scrapper repository
	archieves, err := u.scrapperRepo.GetAllArchieve()
	if err != nil {
		return err
	}

	// Create array (slice) of archieve object
	var newArchieves []*model.Archieve

	// Check if archieve is exist in DB
	for _, archieve := range archieves {
		// Get archieve from DB by archieve code
		_, err := u.mysqlRepo.GetArchieveByCode(ctx, archieve.Code)

		// if archieve not exist then add to newArchieve array (slice)
		if err == model.ErrDataNotFound {
			// Add archieve
			newArchieves = append(newArchieves, archieve)
			log.Printf("New archieve: %v", archieve.Code)
		} else if err != nil && err != model.ErrDataNotFound {
			return err
		}
	}

	// Counter new journal from archieves
	var totalNewJournal int

	// Get new archieves journals
	for _, newArchieve := range newArchieves {
		// Get all journal from scrapper repository based on archieve
		journals, err := u.scrapperRepo.GetAllJournalByArchieveObject(newArchieve)
		if err != nil {
			return err
		}

		// Append Journals into archieve object
		newArchieve.Journals = journals
		totalNewJournal += len(newArchieve.Journals)
	}

	// Check if there's new archieve then saved new archieve into DB
	if len(newArchieves) > 0 {
		// Insert new archieves into DB
		if err := u.mysqlRepo.BatchArchieves(ctx, newArchieves); err != nil {
			return err
		}
	}

	log.Printf("Added %v archieve and %v journal (%v)m", len(newArchieves), totalNewJournal, time.Since(start).Minutes())

	// if there's no update then do nothing or finish pull data from archieve scrapper
	return nil
}
