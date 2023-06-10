package pgsql

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

type Provider struct {
	conn *sql.DB
	tx   *sql.Tx
}

func NewProvider(conn *sql.DB) (*Provider, error) {
	p := &Provider{
		conn: conn,
	}

	return p, nil
}

func (p *Provider) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Provider, error) {
	tx, err := p.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	txProvider := &Provider{
		conn: p.conn,
		tx:   tx,
	}

	return txProvider, nil
}

func (p *Provider) DB() *sql.DB {
	return p.conn
}

func (p *Provider) Tx() *sql.Tx {
	return p.tx
}

func (p *Provider) Close() error {
	if p.conn == nil || p.tx != nil {
		return nil
	}

	err := p.conn.Close()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
