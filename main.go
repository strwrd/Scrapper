package main

import (
	"errors"
	"log"
	"time"

	"github.com/strwrd/scrapper/repository/mysql"
	"github.com/strwrd/scrapper/repository/scrapper"
	"github.com/strwrd/scrapper/tool/uuid"
	"github.com/strwrd/scrapper/usecase"
)

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

	usecaseObj := usecase.NewUsecase(mysqlRepo, scrapperRepo, uuid)

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

	// Looping every 20 minutes
	for true {
		if err := usecaseObj.Update(); err != nil {
			mysqlRepo.Close()
			log.Fatal(err)
		}
		time.Sleep(20 * time.Minute)
	}

}
