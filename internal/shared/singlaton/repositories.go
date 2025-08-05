package singlaton

import (
	"sync"

	"github.com/gabrielmatsan/teste-api/cmd/db"
	"github.com/gabrielmatsan/teste-api/internal/user/repository"
)

var (
	userRepositoryInstance *repository.UserRepository
	userRepositoryOnce     sync.Once
)

func GetUserRepository() *repository.UserRepository {
	userRepositoryOnce.Do(func() {

		database := db.GetDB()
		userRepositoryInstance = repository.NewUserRepository(database)
	})
	return userRepositoryInstance
}
