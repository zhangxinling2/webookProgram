package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}
func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint = 1062
		if mysqlErr.Number == 1062 {
			return ErrUserDuplicateEmail
		}
	}
	return err
}
func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}
func (dao *UserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return u, err
}
func (dao *UserDAO) EditUser(ctx context.Context, u User) (User, error) {
	now := time.Now().UnixMilli()
	u.Utime = now
	err := dao.db.WithContext(ctx).Updates(&u).Error
	return u, err
}

type User struct {
	Id           int64  `gorm:"primaryKey,autoIncrement"`
	Email        string `gorm:"unique"`
	NickName     string
	Birth        int64
	Introduction string
	Password     string
	Ctime        int64
	Utime        int64
}
