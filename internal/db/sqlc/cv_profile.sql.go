// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: cv_profile.sql

package db

import (
	"context"
	"database/sql"
)

const createCvProfile = `-- name: CreateCvProfile :one
INSERT INTO cv_profiles (name, email, phone, address, linkedin_url, github_url, bio, profile_picture)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, name, email, phone, address, linkedin_url, github_url, bio, created_at, profile_picture
`

type CreateCvProfileParams struct {
	Name           string         `json:"name"`
	Email          string         `json:"email"`
	Phone          string         `json:"phone"`
	Address        string         `json:"address"`
	LinkedinUrl    sql.NullString `json:"linkedin_url"`
	GithubUrl      string         `json:"github_url"`
	Bio            string         `json:"bio"`
	ProfilePicture string         `json:"profile_picture"`
}

func (q *Queries) CreateCvProfile(ctx context.Context, arg CreateCvProfileParams) (CvProfile, error) {
	row := q.db.QueryRowContext(ctx, createCvProfile,
		arg.Name,
		arg.Email,
		arg.Phone,
		arg.Address,
		arg.LinkedinUrl,
		arg.GithubUrl,
		arg.Bio,
		arg.ProfilePicture,
	)
	var i CvProfile
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Phone,
		&i.Address,
		&i.LinkedinUrl,
		&i.GithubUrl,
		&i.Bio,
		&i.CreatedAt,
		&i.ProfilePicture,
	)
	return i, err
}

const getCvProfile = `-- name: GetCvProfile :one
SELECT id, name, email, phone, address, linkedin_url, github_url, bio, created_at, profile_picture
FROM cv_profiles
WHERE id = $1
`

func (q *Queries) GetCvProfile(ctx context.Context, id int32) (CvProfile, error) {
	row := q.db.QueryRowContext(ctx, getCvProfile, id)
	var i CvProfile
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Phone,
		&i.Address,
		&i.LinkedinUrl,
		&i.GithubUrl,
		&i.Bio,
		&i.CreatedAt,
		&i.ProfilePicture,
	)
	return i, err
}
