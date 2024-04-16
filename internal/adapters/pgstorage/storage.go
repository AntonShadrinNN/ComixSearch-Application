package pgstorage

import (
	"comixsearch/internal/models"
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StoragePGX struct {
	db *pgxpool.Pool
}

var (
	pgInstance *StoragePGX
	//make sure that the database connection will only be established once per our application lifetime
	pgOnce sync.Once
)

func NewStorage(ctx context.Context, connString string) (*StoragePGX, error) {
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			// TODO: error
			fmt.Errorf("unable to create connection pool: %w", err)
		}

		pgInstance = &StoragePGX{db}
	})

	return pgInstance, nil
}

func (s *StoragePGX) Ping(ctx context.Context) error {
	return s.db.Ping(ctx)
}

func (s *StoragePGX) Close() {
	s.db.Close()
}

func (s *StoragePGX) Write(ctx context.Context, comix models.Comic) error {
	// query := `INSERT INTO comixes (id, title, content, altcontent, link) VALUES (@Id, @Title, @Content, @AltContent, @Link)`

	// batch := &pgx.Batch{}
	// for _, comix := range comixes {
	// 	args := pgx.NamedArgs{
	// 		"Id":         comix.Id,
	// 		"Title":      comix.Title,
	// 		"Content":    comix.Content,
	// 		"AltContent": comix.AltContent,
	// 		"Link":       comix.Link,
	// 	}

	// 	batch.Queue(query, args)
	// }

	// results := s.db.SendBatch(ctx, batch)
	// defer results.Close()

	// for _, comix := range comixes {
	// 	_, err := results.Exec()
	// 	if err != nil {
	// 		var pgErr *pgconn.PgError
	// 		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
	// 			log.Printf("comix %s already exists", comix.Title)
	// 			continue
	// 		}

	// 		return fmt.Errorf("unable to insert row: %w", err)
	// 	}
	// }

	// return results.Close()
	query := `INSERT INTO comixes (id, title, content, altcontent, link) VALUES (@Id, @Title, @Content, @Link)`
	args := pgx.NamedArgs{
		"Id":      comix.Id,
		"Title":   comix.Title,
		"Content": comix.Content,
		// "AltContent": comix.AltContent,
		"Link": comix.Link,
	}

	_, err := s.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	return nil
}

func (s *StoragePGX) Get(ctx context.Context, keywords []string, isContentSearch bool) ([]string, error) {
	var results []string
	var query string

	if isContentSearch {
		query = `
            SELECT title, link
            FROM comixes
            WHERE to_tsvector('russian', content) @@ to_tsquery('russian', $1)
            ORDER BY ts_rank(to_tsvector('russian', content), to_tsquery('russian', $1)) DESC
            LIMIT 10
        `
	} else {
		query = `
            SELECT title, link
            FROM comixes
            WHERE to_tsvector('russian', title) @@ to_tsquery('russian', $1)
            ORDER BY ts_rank(to_tsvector('russian', title), to_tsquery('russian', $1)) DESC
            LIMIT 10
        `
	}

	args := []interface{}{strings.Join(keywords, " & ")}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var title, link string
		if err := rows.Scan(&title, &link); err != nil {
			return nil, err
		}
		results = append(results, title+" - "+link)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
	// query := `SELECT * FROM comixes`

	// rows, err := s.db.Query(ctx, query)
	// if err != nil {
	// 	return nil, fmt.Errorf("unable to query comixes: %w", err)
	// }

	// defer rows.Close()
	// return pgx.CollectRows(rows, pgx.RowToStructByName[models.Comix])
}
