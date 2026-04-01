package db

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

type mockTx struct {
	pgx.Tx
}

func TestWithTx(t *testing.T) {
	q := &Queries{}
	tx := &mockTx{}

	qWithTx := q.WithTx(tx)

	assert.NotNil(t, qWithTx)
	assert.Equal(t, tx, qWithTx.db)
}
