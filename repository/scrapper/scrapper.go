package scrapper

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/strwrd/scrapper/model"
	"github.com/strwrd/scrapper/tool/uuid"
)

// Configuration loader
var (
	_RepoLink = "http://j-ptiik.ub.ac.id/index.php/j-ptiik/issue/archive"
)

// Repository : repository contract interface
type Repository interface {
	GetAllArchieve() ([]*model.Archieve, error)
	GetAllJournalByArchieveObject(archieve *model.Archieve) ([]*model.Journal, error)
}

// repository object
type repository struct {
	uuid uuid.Tool
}

// NewRepository : create repository scrapper object
func NewRepository(uuid uuid.Tool) (Repository, error) {
	return &repository{uuid}, nil
}

// Scrapping : translate link into document html (helper)
func (r *repository) scrapping(link string) (*goquery.Document, error) {
	// Issue new request to http link
	res, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Check status response code
	if res.StatusCode != 200 {
		return nil, err
	}

	// Reading html DOM element from http body then parse to html document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
