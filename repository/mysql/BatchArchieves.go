package mysql

import (
	"context"

	"github.com/strwrd/scrapper/model"
)

// BatchArchieves : Insert bulk data (archieves & journals)
func (r *repository) BatchArchieves(ctx context.Context, archieves []*model.Archieve) error {
	// Creating context mysql transaction
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Query insert for archieve object
	archieveQuery := "INSERT INTO archieve(archieve_id, code, link, published, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?)"

	// Query insert for journal object
	journalQuery := "INSERT INTO journal(journal_id, archieve_id, title, authors, abstract, link, pdf_link, published, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	// Preparing archieve transaction
	archieveStmt, err := tx.PrepareContext(ctx, archieveQuery)
	if err != nil {
		return err
	}
	defer archieveStmt.Close()

	// Preparing journal transaction
	journalStmt, err := tx.PrepareContext(ctx, journalQuery)
	if err != nil {
		return err
	}
	defer journalStmt.Close()

	// Insert into transaction buffer before inserted into DB
	for _, archieve := range archieves {
		// Insert archieve
		if _, err := archieveStmt.ExecContext(
			ctx,
			archieve.ID,
			archieve.Code,
			archieve.Link,
			archieve.Published,
			archieve.CreatedAt,
			archieve.UpdatedAt,
		); err != nil {
			return err
		}

		// Insert every journal in archieve
		for _, journal := range archieve.Journals {
			if _, err := journalStmt.ExecContext(
				ctx,
				journal.ID,
				journal.ArchieveID,
				journal.Title,
				journal.Authors,
				journal.Abstract,
				journal.Link,
				journal.PDFLink,
				journal.Published,
				journal.CreatedAt,
				journal.UpdatedAt,
			); err != nil {
				return err
			}
		}
	}

	// Commited transaction into DB
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
