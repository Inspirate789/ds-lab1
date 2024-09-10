package repository

import (
	"context"
	"database/sql"
	"github.com/Inspirate789/ds-lab1/internal/models"
	"github.com/Inspirate789/ds-lab1/internal/person/usecase"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"log/slog"
)

type sqlxRepository struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewSqlxRepository(db *sqlx.DB, logger *slog.Logger) usecase.Repository {
	return &sqlxRepository{
		db:     db,
		logger: logger,
	}
}

func (r *sqlxRepository) HealthCheck(_ context.Context) error {
	return r.db.Ping()
}

func (r *sqlxRepository) GetPersons(ctx context.Context, offset, limit int64) ([]models.Person, error) {
	var persons Persons

	err := r.db.SelectContext(ctx, &persons, selectPersonsQuery, offset, limit)
	if errors.Is(err, sql.ErrNoRows) {
		return make([]models.Person, 0), nil
	}

	return persons.ToModel(), err
}

func (r *sqlxRepository) CreatePerson(ctx context.Context, person models.PersonProperties) (models.Person, error) {
	properties := NewPersonProperties(person)

	namedQuery, args, err := sqlx.Named(insertPersonQuery, &properties)
	if err != nil {
		return models.Person{}, err
	}

	var identifiedPerson Person

	err = r.db.GetContext(ctx, &identifiedPerson, r.db.Rebind(namedQuery), args...)

	return identifiedPerson.ToModel(), err
}

func (r *sqlxRepository) GetPerson(ctx context.Context, personID int) (models.Person, bool, error) {
	var identifiedPerson Person

	err := r.db.GetContext(ctx, &identifiedPerson, selectPersonQuery, personID)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Person{}, false, nil
	}

	return identifiedPerson.ToModel(), true, err
}

type txFunc func(tx *sqlx.Tx) error

type txRunner interface {
	BeginTxx(context.Context, *sql.TxOptions) (*sqlx.Tx, error)
}

func runTx(ctx context.Context, db txRunner, f txFunc) (err error) {
	var tx *sqlx.Tx

	opts := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}

	tx, err = db.BeginTxx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}

	defer func() {
		if err != nil {
			err = multierr.Combine(err, tx.Rollback())
		} else {
			err = tx.Commit()
		}
	}()

	return f(tx)
}

func (r *sqlxRepository) updatePersonTx(ctx context.Context, tx sqlx.ExtContext, person Person) (Person, error) {
	var res Person

	err := sqlx.GetContext(ctx, tx, &res, selectPersonQuery, person.ID)
	if err != nil {
		return Person{}, err
	}

	res = res.UpdateBy(person.PersonProperties)

	namedQuery, args, err := sqlx.Named(updatePersonQuery, &res)
	if err != nil {
		return Person{}, err
	}

	err = sqlx.GetContext(ctx, tx, &res, r.db.Rebind(namedQuery), args...)
	if err != nil {
		return Person{}, err
	}

	return res, nil
}

func (r *sqlxRepository) UpdatePerson(ctx context.Context, person models.Person) (models.Person, bool, error) {
	var res Person
	var err error

	err = runTx(ctx, r.db, func(tx *sqlx.Tx) error {
		res, err = r.updatePersonTx(ctx, tx, NewPerson(person))
		return err
	})

	if err != nil {
		return models.Person{}, false, err
	}

	if errors.Is(err, sql.ErrNoRows) {
		return models.Person{}, false, nil
	}

	return res.ToModel(), true, nil
}

func (r *sqlxRepository) DeletePerson(ctx context.Context, personID int) (bool, error) {
	_, err := r.db.ExecContext(ctx, deletePersonQuery, personID)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
