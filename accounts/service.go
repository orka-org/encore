package accounts

import (
	"context"

	"encore.dev/rlog"
)

//encore:service
type Service struct {
	l   *AuthLogic
	log rlog.Ctx
}

func initService() (*Service, error) {
	db := usersDB
	repo := NewUsersRepo(db)
	l := NewAuthLogic(repo)
	logger := rlog.With("scope", "accountsService")
	return &Service{
		log: logger,
		l:   l,
	}, nil
}

func (s *Service) Shutdown(force context.Context) {}

//encore:api public path=/auth/register
func (s *Service) Register(ctx context.Context, params *RegisterParams) (*RegisterResponse, error) {
	return s.l.Register(ctx, params)
}

//encore:api public path=/auth/login
func (s *Service) Login(ctx context.Context, params *LoginParams) (*LoginResponse, error) {
	return s.l.Login(ctx, params)
}

//encore:api public path=/auth/refresh
func (s *Service) Refresh(ctx context.Context, params *RefreshParams) (*RefreshResponse, error) {
	return s.l.Refresh(ctx, params)
}

//encore:api public path=/auth/validate
func (s *Service) Validate(ctx context.Context, params *ValidateParams) (*ValidateResponse, error) {
	return s.l.Validate(ctx, params)
}
