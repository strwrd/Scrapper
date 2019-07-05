package usecase

import (
	"github.com/strwrd/scrapper/repository/mysql"
	"github.com/strwrd/scrapper/repository/scrapper"
	"github.com/strwrd/scrapper/tool/uuid"
)

// Usecase : usecase interface contract
type Usecase interface {
	Update() error
}

// usecase object
type usecase struct {
	mysqlRepo    mysql.Repository
	scrapperRepo scrapper.Repository
	uuid         uuid.Tool
}

// NewUsecase : create usecase object
func NewUsecase(mysqlRepo mysql.Repository, scrapperRepo scrapper.Repository, uuid uuid.Tool) Usecase {
	return &usecase{
		mysqlRepo,
		scrapperRepo,
		uuid,
	}
}
