package singlaton

import (
	"sync"

	"github.com/gabrielmatsan/teste-api/cmd/db"
	"github.com/gabrielmatsan/teste-api/internal/user/repository"
)

var (
	userRepository     *repository.UserRepository
	userRepositoryOnce sync.Once
)

func GetUserRepository() *repository.UserRepository {
	userRepositoryOnce.Do(func() {

		userRepository = repository.NewUserRepository(db.GetDB())
	})
	return userRepository
}
