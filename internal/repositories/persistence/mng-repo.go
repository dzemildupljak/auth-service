package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"github.com/dzemildupljak/auth-service/internal/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MngRepo struct {
	db  *mongo.Database
	ctx context.Context
}

func NewMngRepo(ctx context.Context, dbConn *mongo.Database) *MngRepo {
	return &MngRepo{
		db:  dbConn,
		ctx: ctx,
	}
}

func (mngrepo *MngRepo) GetListusers() ([]domain.User, error) {
	return []domain.User{}, nil
}

func (mngrepo *MngRepo) GetUserById(id uuid.UUID) (domain.User, error) {
	return domain.User{}, nil
}

func (mngrepo *MngRepo) GetUserByMail(mail string) (domain.User, error) {
	collection := mngrepo.db.Collection("user")
	var usr domain.User

	query := bson.M{"email": strings.ToLower(mail)}
	err := collection.FindOne(mngrepo.ctx, query).Decode(&usr)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.ErrorLogger.Println(err)
			return domain.User{}, err
		}
		utils.ErrorLogger.Println(err)
		return domain.User{}, err
	}

	return usr, nil
}

func (mngrepo *MngRepo) CreateRegisterUser(usr domain.User) error {
	collection := mngrepo.db.Collection("user")

	usr.CreatedAt = time.Now()
	usr.UpdatedAt = time.Now()

	_, err := collection.InsertOne(mngrepo.ctx, &usr)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			utils.ErrorLogger.Println(err)
			return errors.New("user with that email already exist")
		}
		utils.ErrorLogger.Println(err)
		return err
	}
	return nil
}

func (mngrepo *MngRepo) DeleteUserById(id uuid.UUID) error {
	collection := mngrepo.db.Collection("user")

	query := bson.M{"_id": id}
	res, err := collection.DeleteOne(mngrepo.ctx, query)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.ErrorLogger.Println(err)
			return err
		}
		utils.ErrorLogger.Println(err)
		return err
	}
	fmt.Println(res, err)
	return nil
}