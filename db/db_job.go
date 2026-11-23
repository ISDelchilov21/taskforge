package db

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"taskforge/models"
)

const jobColumns = `id, type, status, description, payload, priority_weight`

func scanJob(row pgx.Row) (models.Job, error) {
	var j models.Job
	err := row.Scan(&j.ID, &j.Type, &j.Status, &j.Description, &j.Payload, &j.PriorityWeight)
	return j, err
}

func CreateJob(ctx context.Context, pool *pgxpool.Pool, job *models.Job) error {
	if len(job.Payload) == 0 {
		job.Payload = json.RawMessage(`{}`)
	}
	return pool.QueryRow(ctx,
		`INSERT INTO jobs (type, description, payload, priority_weight)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, status`,
		job.Type, job.Description, job.Payload, job.PriorityWeight,
	).Scan(&job.ID, &job.Status)
}

func GetJob(ctx context.Context, pool *pgxpool.Pool, id int) (models.Job, error) {
	j, err := scanJob(pool.QueryRow(ctx,
		`SELECT `+jobColumns+` FROM jobs WHERE id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return j, ErrNotFound
	}
	return j, err
}

func ListJobs(ctx context.Context, pool *pgxpool.Pool) ([]models.Job, error) {
	rows, err := pool.Query(ctx, `SELECT `+jobColumns+` FROM jobs ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := []models.Job{}
	for rows.Next() {
		j, err := scanJob(rows)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, rows.Err()
}

func UpdateJobStatus(ctx context.Context, pool *pgxpool.Pool, id int, status string) error {
	return execOne(ctx, pool,
		`UPDATE jobs SET status = $1, updated_at = now() WHERE id = $2`,
		status, id)
}

func DeleteJob(ctx context.Context, pool *pgxpool.Pool, id int) error {
	return execOne(ctx, pool, `DELETE FROM jobs WHERE id = $1`, id)
}
