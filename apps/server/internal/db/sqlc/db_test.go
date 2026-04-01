package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

type mockDBTX struct{}

func (m *mockDBTX) Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (m *mockDBTX) Query(context.Context, string, ...interface{}) (pgx.Rows, error) {
	return nil, nil
}

func (m *mockDBTX) QueryRow(context.Context, string, ...interface{}) pgx.Row {
	return nil
}

func TestNew(t *testing.T) {
	mockDB := &mockDBTX{}

	queries := New(mockDB)

	assert.NotNil(t, queries, "New should return a non-nil Queries object")
	assert.Equal(t, mockDB, queries.db, "New should set the db field correctly")
}
