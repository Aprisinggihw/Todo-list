package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"todo-list/internal/entity"
	"todo-list/internal/repository"
	"todo-list/pkg/cache"
	"todo-list/pkg/token"
)

type TodoService interface {
	CreateTodo(ctx context.Context,userID uint, title string) (*entity.Todo, error)
	GetTodos(ctx context.Context) ([]entity.Todo, error)
	GetTodosByUserID(ctx context.Context,userID uint) ([]entity.Todo, error)
	UpdateTodo(ctx context.Context,userID, todoID uint, title string, done bool) error
	DeleteTodo(ctx context.Context,userID, todoID uint) error
}

type todoService struct {
	repo repository.TodoRepository
	tokenUseCase   token.TokenUseCase
	cacheable      cache.Cacheable
}

func NewTodoService(
	repo repository.TodoRepository, 
	tokenUseCase token.TokenUseCase,
	cacheable cache.Cacheable,
	) TodoService {
	return &todoService{repo, tokenUseCase, cacheable}
}

func (s *todoService) CreateTodo(ctx context.Context,userID uint, title string) (*entity.Todo, error) {
	todo := &entity.Todo{UserID: userID, Title: title}
	err := s.repo.Create(ctx, todo)
	if err != nil {
		return nil, err
	}
	keyGetTodos := "todo-list:todos:get-todos"
	err = s.cacheable.Delete(keyGetTodos) // Menghapus cache lama
	if err != nil {
		return nil, errors.New("falied deleting key cache")
	}
	return todo, err
}

func (s *todoService) GetTodos(ctx context.Context) (result []entity.Todo, err error) {
	//tambahkan cache redis di service nya
	keyGetTodos := "todo-list:todos:get-todos"
	data := s.cacheable.Get(keyGetTodos)
	if data == ""{
		result, err = s.repo.GetAll(ctx)
		if err != nil{
			return nil, err
		}
		marshalledData, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}

		err = s.cacheable.Set(keyGetTodos, marshalledData, 5*time.Minute)
		if err != nil {
			return nil, err
		}
	}else{
		err = json.Unmarshal([]byte(data), &result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (s *todoService) GetTodosByUserID(ctx context.Context,userID uint) (result []entity.Todo, err error) {
	keyGetTodos := "todo-list:todos:get-todos"
	data := s.cacheable.Get(keyGetTodos)
	if data == ""{
		result, err = s.repo.GetByUserID(ctx, userID)
		if err != nil{
			return nil, err
		}
		marshalledData, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}

		err = s.cacheable.Set(keyGetTodos, marshalledData, 5*time.Minute)
		if err != nil {
			return nil, err
		}
	}else{
		err = json.Unmarshal([]byte(data), &result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}


func (s *todoService) UpdateTodo(ctx context.Context, userID, todoID uint, title string, done bool) error {
	todo, err := s.repo.GetByID(ctx, todoID)
	if err != nil || (todo.UserID != userID) {
		return errors.New("unauthorized or not found")
	}
	todo.Title = title
	todo.Done = done
	
	keyGetTodos := "todo-list:todos:get-todos"
	err = s.cacheable.Delete(keyGetTodos) // Menghapus cache lama
	if err != nil {
		return errors.New("falied deleting key cache")
	}
	return s.repo.Update(ctx, todo)
}

func (s *todoService) DeleteTodo(ctx context.Context,userID, todoID uint) error {
	todo, err := s.repo.GetByID(ctx, todoID)
	if err != nil || ( todo.UserID != userID) {
		return errors.New("unauthorized or not found")
	}

	keyGetTodos := "todo-list:todos:get-todos"
	err = s.cacheable.Delete(keyGetTodos) // Menghapus cache lama
	if err != nil {
		return errors.New("falied deleting key cache")
	}
	return s.repo.Delete(ctx, todoID)
}
