package repository

import (
	"context"
	"todo-list/internal/entity"

	"gorm.io/gorm"
)

type TodoRepository interface {
	Create(ctx context.Context,todo *entity.Todo) error
	GetAll(ctx context.Context)	([]entity.Todo, error)
	GetByID(ctx context.Context,id uint) (*entity.Todo, error)
	GetByUserID(ctx context.Context,userID uint) ([]entity.Todo, error)
	Update(ctx context.Context,todo *entity.Todo) error
	Delete(ctx context.Context,id uint) error
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db}
}

func (r *todoRepository) Create(ctx context.Context,todo *entity.Todo) error {
	return r.db.WithContext(ctx).Create(todo).Error
}

func (r *todoRepository) GetAll(ctx context.Context) ([]entity.Todo, error) {
	var todos []entity.Todo
	if err := r.db.WithContext(ctx).
	Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}


func (r *todoRepository) GetByID(ctx context.Context,id uint) (*entity.Todo, error) {
	var todo entity.Todo
	if err := r.db.WithContext(ctx).First(&todo, id).Error; err != nil{
		return nil, err
	}
	return &todo, nil
}

func (r *todoRepository) GetByUserID(ctx context.Context,userID uint) ([]entity.Todo, error) {
	var todos []entity.Todo
	if err := r.db.WithContext(ctx).
	Where("user_id = ?", userID).
	Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *todoRepository) Update(ctx context.Context, todo *entity.Todo) error {
	return r.db.WithContext(ctx).Save(todo).Error
}

func (r *todoRepository) Delete(ctx context.Context,id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Todo{}, id).Error
}
