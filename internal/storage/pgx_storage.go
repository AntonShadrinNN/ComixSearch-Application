package storage

import (
	"comixsearch/internal/models"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// StoragePGX provides facilities for malipulate postgres storage.
type StoragePGX struct {
	db *pgxpool.Pool
}

var (
	pgInstance *StoragePGX
	//make sure that the database connection will only be established once per our application lifetime
	pgOnce sync.Once
)

// The NewStorage function initializes an StoragePGX instance.
func NewStorage(ctx context.Context, connString string) (*StoragePGX, error) {
	var (
		err error
		db  *pgxpool.Pool
	)
	pgOnce.Do(func() {
		db, err = pgxpool.New(ctx, connString)
		if err != nil {
			return
		}

		pgInstance = &StoragePGX{db}
	})

	if err != nil {
		log.Printf("err %s\n", err)
		return nil, err
	}
	return pgInstance, nil
}

func (s *StoragePGX) Ping(ctx context.Context) error {
	return s.db.Ping(ctx)
}

// Close close the storage connection.
func (s *StoragePGX) Close() {
	s.db.Close()
}

// Write is used for writing comices data to db.
func (s *StoragePGX) Write(ctx context.Context, data []models.Comic) error {
	query := `INSERT INTO comices (comic_id, title, content, link) VALUES (@Id, @Title, @Content, @Link)`

	batch := &pgx.Batch{}
	for _, comic := range data {
		args := pgx.NamedArgs{
			"Id":      comic.Id,
			"Title":   comic.Title,
			"Content": comic.Content,
			"Link":    comic.Link,
		}

		batch.Queue(query, args)
	}

	results := s.db.SendBatch(ctx, batch)
	defer results.Close()

	for _, comix := range data {
		_, err := results.Exec()
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				log.Printf("comic %s already exists", comix.Title)
				continue
			}

			return err
		}
	}

	return results.Close()
}

// GetComices is used for retrieving comics based on keywords and limit.
func (s *StoragePGX) GetComices(ctx context.Context, keywords []string, limit int) (map[string]string, error) {
	query := fmt.Sprintf(`
	SELECT title, link
	FROM (
		SELECT title, link, ts_rank(to_tsvector(title), to_tsquery($1)) as ranke
		FROM comices
		WHERE to_tsvector(title) @@ to_tsquery($1)
		UNION
		(SELECT title, link, ts_rank(to_tsvector(content), to_tsquery($1)) as ranke
		FROM comices
		WHERE to_tsvector(content) @@ to_tsquery($1) 
		)
	) as subquery
	ORDER BY ranke DESC
	LIMIT %d;
	`, limit)

	args := []any{strings.Join(keywords, " | ")}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make(map[string]string)
	for rows.Next() {
		var title, link string
		if err := rows.Scan(&title, &link); err != nil {
			return nil, err
		}
		data[title] = link
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

// GetLastId is used for getting the last comic ID from db.
func (s *StoragePGX) GetLastId(ctx context.Context) (int64, error) {
	query := `SELECT MAX(comic_id) AS id FROM comices;`

	row := s.db.QueryRow(ctx, query)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, err
}
