package gord

import (
	"fmt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"reflect"
)

type idType interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | string
}

type gormDatatype interface {
	datatypes.JSON | datatypes.Date | datatypes.Time | datatypes.JSONSlice[any] | datatypes.JSONType[any]
}

type updateType interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | string | gormDatatype
}

type UpdateMap map[string]any

func (m UpdateMap) valid() error {
	for key, value := range m {
		switch value.(type) {
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			string,
			datatypes.JSON, datatypes.Date, datatypes.Time, datatypes.JSONSlice[any], datatypes.JSONType[any]:
		default:
			return fmt.Errorf("invalid type for key '%s': %v (type: %s)", key, value, reflect.TypeOf(value).Name())
		}
	}
	return nil
}

type iCRUDRepository[T any, ID idType] interface {
	Count() (int64, error)
	Delete(model T) error
	DeleteAll(model T) error
	DeleteMany(models []T) error
	DeleteManyByID(model T, ids []ID) error
	DeleteByID(model T, id ID) error
	ExistByID(id ID) (bool, error)
	FindAll() ([]T, error)
	FindAllByID(ids []ID) ([]T, error)
	FindByID(id ID) (T, error)
	Save(t T) error
	SaveAll(ts []T) error
	Updates(t T, m UpdateMap) error
	Create(t T) error
}

type iRepository[T any, ID idType] interface {
	iCRUDRepository[T, ID]
}

type Repository[T any, ID idType] struct {
	tx *gorm.DB
}

func (r Repository[T, ID]) Count() (int64, error) {
	var count int64
	err := r.tx.Count(&count).Error
	return count, err
}

func (r Repository[T, ID]) Delete(model T) error {
	return r.tx.Delete(&model).Error
}

func (r Repository[T, ID]) DeleteAll(model T) error {
	return r.tx.Where("1 = 1").Delete(&model).Error
}

func (r Repository[T, ID]) DeleteMany(models []T) error {
	return r.tx.Delete(&models).Error
}

func (r Repository[T, ID]) DeleteManyByID(model T, ids []ID) error {
	return r.tx.Delete(&model, ids).Error
}

func (r Repository[T, ID]) DeleteByID(model T, id ID) error {
	return r.tx.Delete(&model, id).Error
}

func (r Repository[T, ID]) ExistByID(id ID) (bool, error) {
	var count int64
	err := r.tx.Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r Repository[T, ID]) FindAll() ([]T, error) {
	models := make([]T, 0)
	err := r.tx.Find(&models).Error
	return models, err
}

func (r Repository[T, ID]) FindAllByID(ids []ID) ([]T, error) {
	models := make([]T, 0)
	err := r.tx.Where("id in (?)", ids).Find(&models).Error
	return models, err
}

func (r Repository[T, ID]) FindByID(id ID) (T, error) {
	model := new(T)
	err := r.tx.Where("id = ?", id).Find(model).Error
	return *model, err
}

func (r Repository[T, ID]) Save(t T) error {
	return r.tx.Save(&t).Error
}

func (r Repository[T, ID]) SaveAll(ts []T) error {
	return r.tx.Save(ts).Error
}

func (r Repository[T, ID]) Updates(t T, m UpdateMap) error {

	if err := m.valid(); err != nil {
		return err
	}

	return r.tx.Model(&t).Updates(m).Error
}

func (r Repository[T, ID]) Create(t T) error {
	return r.tx.Create(&t).Error
}
