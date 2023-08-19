package repository

import (
	"context"

	"github.com/de-bait/ent"
	"github.com/de-bait/ent/user"
)

type UserRepository struct {
	ctx    context.Context
	client *ent.Client
}

type UserInput struct {
	Nickname string `validate:"required"`
	Password string `validate:"required"`
}

func NewUserRepository(ctx context.Context, client *ent.Client) *UserRepository {
	return &UserRepository{
		ctx:    ctx,
		client: client,
	}
}

func (ur *UserRepository) Create(ui UserInput) (*ent.User, error) {
	u, err := ur.client.User.Create().SetNickname(ui.Nickname).SetPassword(ui.Password).Save(ur.ctx)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) Find() ([]*ent.User, error) {
	u, err := ur.client.User.Query().All(ur.ctx)
	if err != nil {
		return nil, err
	}
	return u, nil
}
func (ur *UserRepository) FindOne(id int) (*ent.User, error) {
	u, err := ur.client.User.Query().Where(user.ID(id)).Only(ur.ctx)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (ur *UserRepository) FindOneByNickname(nickname string) (*ent.User, error) {
	u, err := ur.client.User.Query().Where(user.NicknameEQ(nickname)).Only(ur.ctx)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (ur *UserRepository) Update(id int, ui UserInput) (*ent.User, error) {
	u, err := ur.client.User.UpdateOneID(id).SetNickname(ui.Nickname).Save(ur.ctx)
	if err != nil {
		return nil, err
	}
	return u, nil
}
