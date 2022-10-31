package data

import (
	"context"
	"database/sql"
	"log"
	"time"
)

const dbTimeout = 3 * time.Second

var db *sql.DB

// Models contains all types we want to be available to our application
type Models struct {
	Exercise Exercise
}

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		Exercise: Exercise{},
	}
}

type Exercise struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
	Date  string `json:"date"`
}

func (e *Exercise) GetAll() ([]*Exercise, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, name, count, date from exercises order by date`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []*Exercise

	for rows.Next() {
		var exercise Exercise
		err = rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Count,
			&exercise.Date,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		exercises = append(exercises, &exercise)
	}

	return exercises, nil
}

func (e *Exercise) GetByDate(date string) ([]*Exercise, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, name, count, date from exercises where date = $1`

	rows, err := db.QueryContext(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []*Exercise

	for rows.Next() {
		var exercise Exercise
		err = rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Count,
			&exercise.Date,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		exercises = append(exercises, &exercise)
	}

	return exercises, nil
}

func (e *Exercise) GetByID(id string) ([]*Exercise, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, name, count, date from exercises where id = $1`

	rows, err := db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []*Exercise

	for rows.Next() {
		var exercise Exercise
		err = rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Count,
			&exercise.Date,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		exercises = append(exercises, &exercise)
	}

	return exercises, nil
}

func (e *Exercise) GetByName(name string) ([]*Exercise, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, name, count, date from exercises where name = $1`

	rows, err := db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []*Exercise

	for rows.Next() {
		var exercise Exercise
		err = rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Count,
			&exercise.Date,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		exercises = append(exercises, &exercise)
	}

	return exercises, nil
}

func (e *Exercise) GetByNameDate(name string, date string) (*Exercise, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, name, count, date from exercises where name = $1 and date = $2`

	var exercise Exercise
	row := db.QueryRowContext(ctx, query, name, date)
	err := row.Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Count,
		&exercise.Date,
	)

	if err != nil {
		return nil, err
	}

	return &exercise, nil
}

func (e *Exercise) Update(name, date string, count int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update exercises set count = $1 where name = $2 and date = $3`

	_, err := db.ExecContext(ctx, stmt,
		count,
		name,
		date,
	)

	if err != nil {
		return err
	}

	return nil
}

func (e *Exercise) Delete(name, date string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from exercises where name = $1 and date = $2`

	_, err := db.ExecContext(ctx, stmt, name, date)
	if err != nil {
		return err
	}

	return nil
}

func (e *Exercise) Insert() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `insert into exercises (name, count, date)
		values ($1, $2, $3) returning id`

	err := db.QueryRowContext(ctx, stmt,
		e.Name,
		e.Count,
		e.Date,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}
