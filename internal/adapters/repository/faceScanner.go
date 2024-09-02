package repository

import (
	"context"
	"database/sql"
	"errors"
	"faceScanner/internal/constants"
	scannerErrors "faceScanner/internal/errors"
	"faceScanner/internal/models"
	"fmt"
	"github.com/jackc/pgconn"
	"golang.org/x/sync/errgroup"
)

const queryExtendFaceScannerTask = `insert into public.image_tasks (task_id, image_data, image_id, file_name) values ($1, $2, $3, $4)`

func (r *Repository) ExtendFaceScannerTask(ctx context.Context, task models.ExtendFaceScannerTaskParamsRepository) (err error) {
	_, err = r.db.ExecContext(ctx, queryExtendFaceScannerTask, task.TaskUUID, task.Image, task.ImageUUID, task.FileName)
	if err != nil {
		err = fmt.Errorf("r.db.ExecContext(...): %w", err)
		return err
	}

	return nil
}

const queryGetFaceScannerTaskData = `select image_data, api_response, image_id, file_name from public.image_tasks where task_id = $1`
const queryGetFaceScannerTask = `select id, status from public.tasks where id = $1`

func (r *Repository) GetFaceScannerTask(ctx context.Context, taskUUID string) (task models.GetFaceScannerTaskResponseRepository, err error) {
	err = r.db.SelectContext(ctx, &task.ImagesData, queryGetFaceScannerTaskData, taskUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return task, scannerErrors.ErrTaskNotFound
		}
		err = fmt.Errorf("r.db.GetContext(...): %w", err)
		return task, err
	}

	err = r.db.GetContext(ctx, &task, queryGetFaceScannerTask, taskUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return task, scannerErrors.ErrTaskNotFound
		}
		err = fmt.Errorf("r.db.GetContext(...): %w", err)
		return task, err
	}

	return task, nil
}

const queryChangeFaceScannerTaskStatus = `update public.tasks set status = $1 where id = $2`

func (r *Repository) ChangeFaceScannerTaskStatus(ctx context.Context, taskUUID string, statusID int) (err error) {

	_, err = r.db.ExecContext(ctx, queryChangeFaceScannerTaskStatus, statusID, taskUUID)
	if err != nil {
		err = fmt.Errorf("r.db.ExecContext(...): %w", err)
		return err
	}

	return nil
}

const queryDeleteFaceScannerTask = `delete from public.tasks where id = $1`
const queryDeleteFaceScannerTaskImages = `delete from public.image_tasks where task_id = $1`

func (r *Repository) DeleteFaceScannerTask(ctx context.Context, taskUUID string) (err error) {
	eg, ctx := errgroup.WithContext(ctx)
	ctx = context.WithoutCancel(ctx)
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		err = fmt.Errorf("r.db.BeginTx(...): %w", err)
		return err
	}

	eg.Go(func() error {
		_, err = tx.ExecContext(ctx, queryDeleteFaceScannerTask, taskUUID)
		if err != nil {
			err = fmt.Errorf("tx.ExecContext(...): %w", err)
			return err
		}
		return nil
	})

	eg.Go(func() error {
		_, err = tx.ExecContext(ctx, queryDeleteFaceScannerTaskImages, taskUUID)
		if err != nil {
			err = fmt.Errorf("tx.ExecContext(...): %w", err)
			return err
		}
		return nil
	})
	if err = eg.Wait(); err != nil {
		tx.Rollback()
		err = fmt.Errorf("eg.Wait(): %w", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("tx.Commit(): %w", err)
		return err
	}

	return nil
}

const queryCreateFaceScannerTask = `insert into public.tasks (id, status) values ($1, $2)`
const queryCreateFaceScannerTaskImage = `insert into public.image_tasks (task_id, image_data, image_id, file_name) values ($1, $2, $3, $4)`

func (r *Repository) CreateFaceScannerTask(ctx context.Context, task models.CreateFaceScannerTaskParamsRepository) (err error) {
	eg, ctx := errgroup.WithContext(ctx)
	ctx = context.WithoutCancel(ctx)
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		err = fmt.Errorf("r.db.BeginTx(...): %w", err)
		return err
	}
	eg.Go(func() error {
		_, err = tx.ExecContext(ctx, queryCreateFaceScannerTask, task.TaskUUID, constants.StatusNew)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				if pgErr.Code == "23505" {
					return scannerErrors.ErrDuplicateTask
				}
			}
			err = fmt.Errorf("r.db.ExecContext(...): %w", err)
			return err
		}
		return nil
	})
	eg.Go(func() error {
		_, err = tx.ExecContext(ctx, queryCreateFaceScannerTaskImage, task.TaskUUID, task.Image, task.ImageUUID, task.FileName)
		if err != nil {
			err = fmt.Errorf("r.db.ExecContext(...): %w", err)
			return err
		}
		return nil
	})

	if err = eg.Wait(); err != nil {
		tx.Rollback()
		err = fmt.Errorf("eg.Wait(): %w", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("tx.Commit(): %w", err)
		return err
	}

	return nil
}

const queryUpdateTaskImageInfo = `update public.image_tasks set api_response = $1 where task_id = $2`

func (r *Repository) UpdateTaskImageInfo(ctx context.Context, task models.UpdateTaskImageInfoParamsRepository) (err error) {
	_, err = r.db.ExecContext(ctx, queryUpdateTaskImageInfo, task.ApiResponse, task.TaskUUID)
	if err != nil {
		err = fmt.Errorf("r.db.SelectContext(...): %w", err)
		return err
	}

	return nil
}
