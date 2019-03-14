package orm

import (
	"context"

	"github.com/aymerick/raymond"
	"github.com/jmoiron/sqlx"
	"github.com/phogolabs/prana/sqlexec"
)

func namedSelectMany(ctx context.Context, preparer Preparer, provider *sqlexec.Provider, dest Entity, query NamedQuery) error {
	stmt, args, err := prepareNamedQuery(preparer, provider, query, dest)
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

func namedSelectOne(ctx context.Context, preparer Preparer, provider *sqlexec.Provider, dest Entity, query NamedQuery) error {
	stmt, args, err := prepareNamedQuery(preparer, provider, query, dest)
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

func namedQueryRows(ctx context.Context, preparer Preparer, provider *sqlexec.Provider, query NamedQuery) (*Rows, error) {
	stmt, args, err := prepareNamedQuery(preparer, provider, query, nil)
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

func namedQueryRow(ctx context.Context, preparer Preparer, provider *sqlexec.Provider, query NamedQuery) (*Row, error) {
	stmt, args, err := prepareNamedQuery(preparer, provider, query, nil)
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

func namedExec(ctx context.Context, preparer Preparer, provider *sqlexec.Provider, query NamedQuery) (Result, error) {
	stmt, args, err := prepareNamedQuery(preparer, provider, query, nil)
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

func prepareNamedQuery(preparer Preparer, provider *sqlexec.Provider, nquery NamedQuery, dest Entity) (*sqlx.NamedStmt, map[string]Param, error) {
	var err error

	switch q := nquery.(type) {
	case *routine:
		err = q.Prepare(provider)
	case *query:
		err = q.Prepare(dest)
	}

	if err != nil {
		return nil, nil, err
	}

	body, args := nquery.NamedQuery()

	template, err := raymond.Parse(body)
	if err != nil {
		return nil, nil, err
	}

	body, err = template.Exec(args)
	if err != nil {
		return nil, nil, err
	}

	namedStmt, err := preparer.PrepareNamed(body)
	if err != nil {
		return nil, nil, err
	}

	return namedStmt, args, nil
}
