package oak

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func selectMany(ctx context.Context, preparer Preparer, dest Entity, query Query) error {
	stmt, args, err := prepareQuery(preparer, query)
	if err != nil {
		return err
	}

	defer func() {
		if stmtErr := stmt.Close(); err == nil {
			err = stmtErr
		}
	}()

	err = stmt.SelectContext(ctx, dest, args)
	return err
}

func selectOne(ctx context.Context, preparer Preparer, dest Entity, query Query) error {
	stmt, args, err := prepareQuery(preparer, query)
	if err != nil {
		return err
	}

	defer func() {
		if stmtErr := stmt.Close(); err == nil {
			err = stmtErr
		}
	}()

	err = stmt.GetContext(ctx, dest, args)
	return err
}

func queryRows(ctx context.Context, preparer Preparer, query Query) (*Rows, error) {
	stmt, args, err := prepareQuery(preparer, query)
	if err != nil {
		return nil, err
	}

	defer func() {
		if stmtErr := stmt.Close(); err == nil {
			err = stmtErr
		}
	}()

	var rows *Rows
	rows, err = stmt.QueryxContext(ctx, args)
	return rows, err
}

func queryRow(ctx context.Context, preparer Preparer, query Query) (*Row, error) {
	stmt, args, err := prepareQuery(preparer, query)
	if err != nil {
		return nil, err
	}

	defer func() {
		if stmtErr := stmt.Close(); err == nil {
			err = stmtErr
		}
	}()

	return stmt.QueryRowxContext(ctx, args), nil
}

func exec(ctx context.Context, preparer Preparer, query Query) (Result, error) {
	stmt, args, err := prepareQuery(preparer, query)
	if err != nil {
		return nil, err
	}

	defer func() {
		if stmtErr := stmt.Close(); err == nil {
			err = stmtErr
		}
	}()

	var result Result
	result, err = stmt.ExecContext(ctx, args)
	return result, err
}

func prepareQuery(preparer Preparer, query Query) (*sqlx.NamedStmt, map[string]interface{}, error) {
	body, args := query.NamedQuery()

	stmt, err := preparer.PrepareNamed(body)
	if err != nil {
		return nil, nil, err
	}

	return stmt, args, nil
}
