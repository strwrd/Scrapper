package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/strwrd/scrapper/model"
	"github.com/strwrd/scrapper/repository/mysql"
	"github.com/strwrd/scrapper/repository/scrapper"
	"github.com/strwrd/scrapper/tool/uuid"
)

type initializer struct {
	uuid         uuid.Tool
	mysqlRepo    mysql.Repository
	scrapperRepo scrapper.Repository
}

func (i *initializer) Execute() error {
	// Time execution
	start := time.Now()

	// Handling panic (unknown error)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			// find out exactly what the error was and set err
			var err error
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	log.Println("====== STARTING ======")

	// Set maximum timeout process
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Get all archieve from scrapper repository
	log.Println("====== SCRAPPING ALL ACHIEVE ======")
	archieves, err := i.scrapperRepo.GetAllArchieve()
	if err != nil {
		return err
	}

	log.Println("====== SCRAPPING ALL JOURNAL ======")

	// Creating syncronization for pararel process
	var wg sync.WaitGroup

	// Iterating each archieve to retrieve all journal
	for _, archieve := range archieves {
		// Adding worker process pararelism
		wg.Add(1)

		// Start pararel scrapping journals each archieve (IIFE style)
		go func(archieve *model.Archieve, wg *sync.WaitGroup) {
			journals, err := i.scrapperRepo.GetAllJournalByArchieveObject(archieve)
			if err != nil {
				log.Fatal(err)
			}
			archieve.Journals = journals
			wg.Done()
		}(archieve, &wg)
	}

	// Wait until all process of pararelism are done
	wg.Wait()

	var totalJournal int

	// Count available journal
	for _, archieve := range archieves {
		totalJournal += len(archieve.Journals)
	}

	log.Println("====== SAVING TO DB PLEASE WAIT ======")

	if err := i.mysqlRepo.BatchArchieves(ctx, archieves); err != nil {
		log.Fatal(err)
	}

	log.Printf("TOTAL ARCHIEVE:\t %v\n", len(archieves))
	log.Printf("TOTAL JOURNAL:\t %v\n", totalJournal)
	log.Printf("====== DONE ✓ (%v m)======", time.Since(start).Minutes())

	return nil
}

func main() {
	// Initial uuid tool object
	uuid := uuid.NewHandler()

	// Initial mysql repository object
	mysqlRepo, err := mysql.NewRepository()
	if err != nil {
		log.Fatal(err)
	}

	// Initial scrapper repository object
	scrapperRepo, err := scrapper.NewRepository(uuid)
	if err != nil {
		log.Fatal(err)
	}

	// Creating init object
	init := &initializer{
		uuid:         uuid,
		mysqlRepo:    mysqlRepo,
		scrapperRepo: scrapperRepo,
	}

	// Handling panic (unknown error)
	defer func() {
		if r := recover(); r != nil {
			var err error

			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
			if err != nil {
				mysqlRepo.Close()
				log.Fatal(err)
			}
		}
	}()

	// Do initialization execution
	if err := init.Execute(); err != nil {
		mysqlRepo.Close()
		log.Fatal(err)
	}
}
