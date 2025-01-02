package repository_test

import (
	"context"
	"testing"

	"github.com/AugustineAurelius/eos/example/common"
	"github.com/AugustineAurelius/eos/example/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_WithSQLite(t *testing.T) {
	db, err := common.NewSqliteInMemory(context.Background())
	assert.NoError(t, err)

	_, err = db.Exec(context.Background(), `CREATE TABLE users (
    id TEXT PRIMARY KEY,
    name TEXT,
	email TEXT);`)
	assert.NoError(t, err)

	userRepo := repository.New(&db)

	id := uuid.New()
	testUser := &repository.User{ID: id, Name: "name", Email: "email"}
	err = userRepo.CreateUser(context.Background(), testUser)
	assert.NoError(t, err)

	user, err := userRepo.GetUser(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, testUser, user)

	f := repository.NewFilter().AddOneToIDs(id)
	users, err := userRepo.GetManyUsers(context.Background(), *f)
	assert.NoError(t, err)
	assert.Equal(t, []repository.User{*testUser}, users)
}
