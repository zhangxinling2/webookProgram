//go:build wireinject

// 让wire来注入这里的代码
package wire

import (
	"github.com/google/wire"
	"webookProgram/wire/repository"
	"webookProgram/wire/repository/dao"
)

func InitRepository() *repository.UserRepository {
	wire.Build(repository.NewUserRepository, dao.NewUserDAO, InitDB)
	return new(repository.UserRepository)
}
