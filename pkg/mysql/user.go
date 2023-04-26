package mysql

import (
	"context"
	"database/sql"

	"github.com/markhaur/messapp-backend/pkg"
	"github.com/markhaur/messapp-backend/pkg/mysql/gen"
)

type userRepository struct {
	queries *gen.Queries
}

func NewUserRepository(db *sql.DB) pkg.UserRepository {
	return &userRepository{queries: gen.New(db)}
}

func (u *userRepository) Insert(ctx context.Context, user *pkg.User) error {
	inserted, err := u.queries.CreateUser(ctx, gen.CreateUserParams{Name: user.Name, Password: user.Password, Designation: user.Designation, EmployeeID: user.EmployeeID, IsActive: user.IsActive, IsAdmin: user.IsAdmin, CreatedAt: user.CreatedAt})
	if err != nil {
		return err
	}
	user.ID, _ = inserted.LastInsertId()
	return nil
}

func (u *userRepository) FindAll(ctx context.Context) ([]pkg.User, error) {
	users, err := u.queries.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	var list []pkg.User
	for _, user := range users {
		list = append(list, pkg.User{ID: user.ID, Name: user.Name, Password: user.Password, Designation: user.Designation, IsActive: user.IsActive, IsAdmin: user.IsAdmin, EmployeeID: user.EmployeeID, CreatedAt: user.CreatedAt})
	}
	return list, nil
}

func (u *userRepository) FindByID(ctx context.Context, id int64) (*pkg.User, error) {
	user, err := u.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &pkg.User{ID: user.ID, Name: user.Name, Password: user.Password, Designation: user.Designation, IsActive: user.IsActive, IsAdmin: user.IsAdmin, EmployeeID: user.EmployeeID, CreatedAt: user.CreatedAt}, nil
}

func (u *userRepository) FindByEmployeeID(ctx context.Context, employee_id string) (*pkg.User, error) {
	user, err := u.queries.GetUserByEmployeeID(ctx, employee_id)
	if err != nil {
		return nil, err
	}
	return &pkg.User{ID: user.ID, Name: user.Name, Password: user.Password, Designation: user.Designation, IsActive: user.IsActive, IsAdmin: user.IsAdmin, EmployeeID: user.EmployeeID, CreatedAt: user.CreatedAt}, nil
}

func (u *userRepository) Update(ctx context.Context, user *pkg.User) error {
	return u.queries.UpdateUser(ctx, gen.UpdateUserParams{ID: user.ID, Password: user.Password})
}

func (u *userRepository) DeleteByID(ctx context.Context, id int64) error {
	return u.queries.DeleteUser(ctx, id)
}
