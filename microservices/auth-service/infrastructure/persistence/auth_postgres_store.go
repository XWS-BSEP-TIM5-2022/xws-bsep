package persistence

import (
	"errors"
	"log"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/domain"
	"gorm.io/gorm"
)

type AuthPostgresStore struct {
	db *gorm.DB
}

func NewAuthPostgresStore(db *gorm.DB) (*AuthPostgresStore, error) {
	err := db.AutoMigrate(&domain.Authentication{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&domain.Permission{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&domain.Role{})
	if err != nil {
		return nil, err
	}

	return &AuthPostgresStore{
		db: db,
	}, nil
}

func (store *AuthPostgresStore) Create(auth *domain.Authentication) (*domain.Authentication, error) {
	err := store.db.Create(auth)
	return auth, err.Error
}

func (store *AuthPostgresStore) FindByUsername(username string) (*domain.Authentication, error) {
	var auth domain.Authentication
	err := store.db.First(&auth, "username = ?", username)
	return &auth, err.Error
}

func (store *AuthPostgresStore) FindAll() (*[]domain.Authentication, error) {
	var auths []domain.Authentication
	result := store.db.Preload("Roles").Find(&auths)
	if result.Error != nil {
		return nil, result.Error
	}
	return &auths, nil
}

func (store *AuthPostgresStore) DeleteAll() {
	store.db.Exec("DELETE FROM role_permissions CASCADE")
	store.db.Exec("DELETE FROM auth_roles CASCADE")
	store.db.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&domain.Authentication{})
	store.db.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&domain.Role{})
	store.db.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&domain.Permission{})
}

func (store *AuthPostgresStore) Insert(auth *domain.Authentication) error {
	result := store.db.Create(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (store *AuthPostgresStore) UpdateUsername(id, username string) (*domain.Authentication, error) {
	var auth domain.Authentication
	err := store.db.First(&auth, "id = ?", id)
	store.db.Model(&domain.Authentication{}).Where("Id = ?", id).Update("Username", username)
	log.Println(auth)
	if err != nil {
		return nil, err.Error
	}
	return &auth, nil
}

func updateUsernameById(tx *gorm.DB, auth *domain.Authentication, username string) error {
	tx = tx.Model(&domain.Authentication{}).
		Where("Id = ?", auth.Id).
		Update("Username", gorm.Expr(username))
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected != 1 {
		return errors.New("update error")
	}
	return nil
}

func (store *AuthPostgresStore) FindById(id string) (*domain.Authentication, error) {
	var auth domain.Authentication
	err := store.db.First(&auth, "id = ?", id)
	return &auth, err.Error
}

func (store *AuthPostgresStore) UpdatePassword(id, password string) error {
	var auth domain.Authentication
	err := store.db.First(&auth, "id = ?", id)
	store.db.Model(&domain.Authentication{}).Where("Id = ?", id).Update("Password", password)
	if err != nil {
		return err.Error
	}
	return nil
}

func (store *AuthPostgresStore) UpdateVerifactionCode(id, code string) error {
	var auth domain.Authentication
	err := store.db.First(&auth, "id = ?", id)
	store.db.Model(&domain.Authentication{}).Where("Id = ?", id).Update("VerificationCode", code)
	if err != nil {
		return err.Error
	}
	return nil
}

func (store *AuthPostgresStore) UpdateExpirationTime(id string, expTime int64) error {
	var auth domain.Authentication
	err := store.db.First(&auth, "id = ?", id)
	store.db.Model(&domain.Authentication{}).Where("Id = ?", id).Update("ExpirationTime", expTime)
	if err != nil {
		return err.Error
	}
	return nil
}

func (store *AuthPostgresStore) FindAllRolesAndPermissions() (*[]domain.Role, error) {
	var roles []domain.Role
	result := store.db.Preload("Permissions").Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return &roles, nil
}

func (store *AuthPostgresStore) InsertRole(role *domain.Role) error {
	result := store.db.Create(role)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (store *AuthPostgresStore) GetAllPermissionsByRole(role string) (*[]domain.Permission, error) {
	var permissions []domain.Permission
	roles, err := store.FindAllRolesAndPermissions()
	if err != nil {
		return nil, err
	}
	for _, role := range *roles {
		if role.Name == "User" {
			for _, permission := range role.Permissions {
				permissions = append(permissions, *permission)
			}
		}
	}
	return &permissions, nil
}

func (store *AuthPostgresStore) FindRoleByName(name string) (*[]domain.Role, error) {
	var roles []domain.Role
	rolesWithPermissions, err := store.FindAllRolesAndPermissions()
	for _, role := range *rolesWithPermissions {
		if role.Name == name {
			roles = append(roles, role)
		}
	}
	return &roles, err
}
