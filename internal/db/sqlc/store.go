package db

import (
	"context"
	"database/sql"
)

type Store interface {
	Querier
	ListProjectsWithTechnologies(ctx context.Context, arg ListProjectsWithTechnologiesParams) ([]ListProjectsWithTechnologiesRow, error)
}

// SQLStore provides all functions to execute db queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

type ListProjectsWithTechnologiesParams struct {
	CvProfileID int32
	Limit       int32
	Offset      int32
}

type ListProjectsWithTechnologiesRow struct {
	ID               int32                           `json:"id"`
	Title            string                          `json:"title"`
	ShortDescription string                          `json:"short_description"`
	Description      string                          `json:"description"`
	Image            string                          `json:"image"`
	HexThemeColor    string                          `json:"hex_theme_color"`
	ProjectUrl       string                          `json:"project_url"`
	Significance     int32                           `json:"significance"`
	TechnologiesUsed []ListTechnologiesForProjectRow `json:"technologies_used"`
}

// ListProjectsWithTechnologies returns a list of projects with technologies
func (store *SQLStore) ListProjectsWithTechnologies(ctx context.Context, arg ListProjectsWithTechnologiesParams) ([]ListProjectsWithTechnologiesRow, error) {
	params := ListProjectsParams{
		CvProfileID: arg.CvProfileID,
		Limit:       arg.Limit,
		Offset:      arg.Offset,
	}
	projects, err := store.ListProjects(ctx, params)
	if err != nil {
		return nil, err
	}

	var rows []ListProjectsWithTechnologiesRow
	for _, project := range projects {
		technologies, err := store.ListTechnologiesForProject(ctx, project.ID)
		if err != nil {
			return nil, err
		}

		rows = append(rows, ListProjectsWithTechnologiesRow{
			ID:               project.ID,
			Title:            project.Title,
			ShortDescription: project.ShortDescription,
			Description:      project.Description,
			Image:            project.Image,
			HexThemeColor:    project.HexThemeColor,
			ProjectUrl:       project.ProjectUrl,
			Significance:     project.Significance,
			TechnologiesUsed: technologies,
		})
	}

	return rows, nil
}
