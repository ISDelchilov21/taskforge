package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"taskforge/models"
)

func ClaimNextJob(ctx context.Context, pool *pgxpool.Pool) (models.Job, error) {
	j, err := scanJob(pool.QueryRow(ctx,
		`UPDATE jobs SET status = 'running', updated_at = now()
		 WHERE id = (
		     SELECT id FROM jobs
		     WHERE status = 'pending'
		     ORDER BY priority_weight DESC, id
		     LIMIT 1
		     FOR UPDATE SKIP LOCKED
		 )
		 RETURNING `+jobColumns))
	if errors.Is(err, pgx.ErrNoRows) {
		return j, ErrNoJobs
	}
	return j, err
}

func MarkJobDone(ctx context.Context, pool *pgxpool.Pool, id int) error {
	return finishJob(ctx, pool, id, "done")
}

func MarkJobFailed(ctx context.Context, pool *pgxpool.Pool, id int) error {
	return finishJob(ctx, pool, id, "failed")
}

func finishJob(ctx context.Context, pool *pgxpool.Pool, id int, status string) error {
	return execOne(ctx, pool,
		`UPDATE jobs SET status = $1, updated_at = now()
		 WHERE id = $2 AND status = 'running'`,
		status, id)
}
