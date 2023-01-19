package repositories

import (
	"context"

	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type MngRepo struct {
	db *mongo.Client
}

func NewMngRepo(dbConn *mongo.Client) *MngRepo {
	return &MngRepo{
		db: dbConn,
	}
}

func (mngrepo *MngRepo) GetListusers(ctx context.Context) ([]domain.User, error) {
	return []domain.User{}, nil
}
func (mngrepo *MngRepo) GetUserById(ctx context.Context, id int64) (domain.User, error) {
	return domain.User{}, nil
}
func (mngrepo *MngRepo) GetUserByMail(ctx context.Context, mail string) (domain.User, error) {
	return domain.User{}, nil
}
func (mngrepo *MngRepo) CreateRegisterUser(ctx context.Context, usr domain.User) error {
	return nil
}
