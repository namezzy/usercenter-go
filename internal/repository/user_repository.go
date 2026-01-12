package repository

import (
	"user-center/internal/model"

	"gorm.io/gorm"
)

// UserRepository 用户数据访问层
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByUserAccount 根据账号查询用户
func (r *UserRepository) FindByUserAccount(userAccount string) (*model.User, error) {
	var user model.User
	err := r.db.Where("userAccount = ? AND isDelete = 0", userAccount).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUserAccountAndPassword 根据账号和密码查询用户
func (r *UserRepository) FindByUserAccountAndPassword(userAccount, password string) (*model.User, error) {
	var user model.User
	err := r.db.Where("userAccount = ? AND userPassword = ? AND isDelete = 0", userAccount, password).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByPlanetCode 根据星球编号查询用户
func (r *UserRepository) FindByPlanetCode(planetCode string) (*model.User, error) {
	var user model.User
	err := r.db.Where("planetCode = ? AND isDelete = 0", planetCode).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 根据ID查询用户
func (r *UserRepository) FindByID(id int64) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ? AND isDelete = 0", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// SearchByUsername 根据用户名模糊搜索
func (r *UserRepository) SearchByUsername(username string) ([]model.User, error) {
	var users []model.User
	query := r.db.Where("isDelete = 0")
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	err := query.Find(&users).Error
	return users, err
}

// DeleteByID 根据ID删除用户（逻辑删除）
func (r *UserRepository) DeleteByID(id int64) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("isDelete", 1).Error
}
