// Package app contain multiple Go source files that are related to each other. It helps organize
// and group related code together for better maintainability and reusability.
package app

import "context"

// App defines a method for searching with keywords and a limit.
type App interface {
	Search(ctx context.Context, keywords []string, limit int) (map[string]string, error)
}
