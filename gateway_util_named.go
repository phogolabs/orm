package oak

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func namedSelectMany(ctx context.Context, preparer Preparer, dest Entity, query NamedQuery) error {
	stmt, args, err := prepareNamedQuery(preparer, query)
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

func namedSelectOne(ctx context.Context, preparer Preparer, dest Entity, query NamedQuery) error {
	stmt, args, err := prepareNamedQuery(preparer, query)
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

func namedQueryRows(ctx context.Context, preparer Preparer, query NamedQuery) (*Rows, error) {
	stmt, args, err := prepareNamedQuery(preparer, query)
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

func namedQueryRow(ctx context.Context, preparer Preparer, query NamedQuery) (*Row, error) {
	stmt, args, err := prepareNamedQuery(preparer, query)
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

func namedExec(ctx context.Context, preparer Preparer, query NamedQuery) (Result, error) {
	stmt, args, err := prepareNamedQuery(preparer, query)
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

func prepareNamedQuery(preparer Preparer, query NamedQuery) (*sqlx.NamedStmt, map[string]interface{}, error) {
	body, args := query.NamedQuery()

	stmt, err := preparer.PrepareNamed(body)
	if err != nil {
		return nil, nil, err
	}

	return stmt, args, nil
}
