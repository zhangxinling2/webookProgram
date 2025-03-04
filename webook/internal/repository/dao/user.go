package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicate = errors.New("邮箱冲突")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	EditUser(ctx context.Context, u User) (User, error)
}
type GORMUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}
func (dao *GORMUserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint = 1062
		if mysqlErr.Number == 1062 {
			return ErrUserDuplicate
		}
	}
	return err
}
func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}
func (dao *GORMUserDAO) FindByPhone(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", email).First(&u).Error
	return u, err
}
func (dao *GORMUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return u, err
}
func (dao *GORMUserDAO) EditUser(ctx context.Context, u User) (User, error) {
	now := time.Now().UnixMilli()
	u.Utime = now
	err := dao.db.WithContext(ctx).Updates(&u).Error
	return u, err
}

type User struct {
	Id    int64          `gorm:"primaryKey,autoIncrement"`
	Email sql.NullString `gorm:"unique"`
	//唯一索引允许有空值但不能有多个""
	Phone        sql.NullString `gorm:"unique"`
	NickName     string
	Birth        int64
	Introduction string
	Password     string
	Ctime        int64
	Utime        int64
}
