package mysql

import (
	"context"
	"database/sql"

	"github.com/strwrd/scrapper/model"
)

// GetArchieveByCode : return archieve based on code
func (r *repository) GetArchieveByCode(ctx context.Context, code string) (*model.Archieve, error) {
	// Mysql query
	query := "SELECT a.archieve_id, a.code, a.link, a.published, a.created_at, a.updated_at FROM archieve AS a WHERE a.code = ?"

	// Create archieve object
	archieve := new(model.Archieve)

	// Query execution and copy result into archieve object
	err := r.conn.QueryRowContext(ctx, query, code).Scan(
		&archieve.ID,
		&archieve.Code,
		&archieve.Link,
		&archieve.Published,
		&archieve.CreatedAt,
		&archieve.UpdatedAt,
	)

	// Check if any error during execution
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrDataNotFound
		}
		return nil, err
	}

	return archieve, nil
}
