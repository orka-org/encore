package accounts

import (
	"context"

	"encore.app/accounts/lib"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"golang.org/x/crypto/bcrypt"
)

type AuthLogic struct {
	log    rlog.Ctx
	repo   UsersRepo
	tokens lib.TokensUsecase
}

func NewAuthLogic(repo UsersRepo) *AuthLogic {
	logger := rlog.With("scope", "authLogic")
	tokens, err := lib.NewTokensUsecase()
	if err != nil {
		logger.Error("Could not create tokens service", "err", err.Error())
		return nil
	}
	return &AuthLogic{
		log:    logger,
		repo:   repo,
		tokens: tokens,
	}
}

func (l *AuthLogic) Register(ctx context.Context, params *RegisterParams) (*RegisterResponse, error) {
	pwHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		l.log.Error("Could not hash password", "err", err.Error())
		err := &errs.Error{
			Code:    errs.Internal,
			Message: "Could not hash password" + err.Error(),
		}
		return nil, err
	}
	l.log.Debug("Hashed password", "password", params.Password, "hash", string(pwHash))
	u, err := l.repo.Create(ctx, params.Username, params.Email, string(pwHash))
	if err != nil {
		l.log.Error("Could not create user", "err", err.Error())
		err := &errs.Error{
			Code:    errs.Internal,
			Message: "Could not create user" + err.Error(),
		}
		return nil, err
	}
	l.log.Info("User created", "user", u)

	access, err := l.buildToken(ctx, lib.AccessToken, u)
	if err != nil {
		return nil, err
	}
	refresh, err := l.buildToken(ctx, lib.RefreshToken, u)
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (l *AuthLogic) Login(ctx context.Context, params *LoginParams) (*LoginResponse, error) {
	u, err := l.repo.GetByUsername(ctx, params.Username)
	if err != nil {
		err := &errs.Error{
			Code:    errs.Internal,
			Message: "Could not get user" + err.Error(),
		}
		l.log.Error("Could not get user", "err", err.Error())
		return nil, err
	}
	l.log.Debug("User found", "user", u)
	if u == nil {
		err := &errs.Error{
			Code:    errs.NotFound,
			Message: "User not found",
		}
		l.log.Error(err.Message, "username", params.Username)
		return nil, err
	}
	l.log.Info("User found", "user", u)

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(params.Password))
	if err != nil {
		err := &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "Invalid password",
		}
		l.log.Error(err.Message, "username", params.Username)
		return nil, err
	}

	access, err := l.buildToken(ctx, lib.AccessToken, u)
	if err != nil {
		return nil, err
	}
	refresh, err := l.buildToken(ctx, lib.RefreshToken, u)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (l *AuthLogic) Refresh(ctx context.Context, params *RefreshParams) (*RefreshResponse, error) {
	t, err := l.tokens.Parse(ctx, params.RefreshToken)
	if err != nil {
		l.log.Error("Could not parse refresh token", "err", err.Error())
		return nil, err
	}
	id, ok := t.Subject()
	if !ok || id == "" {
		err := &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "Invalid refresh token",
		}
		l.log.Error(err.Message, "err", err.Error())
		return nil, err
	}
	u, err := l.repo.GetByID(ctx, id)
	if err != nil {
		l.log.Error("Could not get user", "err", err.Error())
		return nil, err
	}
	access, err := l.buildToken(ctx, lib.AccessToken, u)
	if err != nil {
		return nil, err
	}
	return &RefreshResponse{
		AccessToken: access,
	}, nil
}

func (l *AuthLogic) Validate(ctx context.Context, params *ValidateParams) (*ValidateResponse, error) {
	t, err := l.tokens.Parse(ctx, params.AccessToken)
	if err != nil {
		l.log.Error("Could not parse access token", "err", err.Error())
		return nil, err
	}
	id, ok := t.Subject()
	if !ok || id == "" {
		err := &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "Invalid access token",
		}
		l.log.Error(err.Message, "err", err.Error())
		return nil, err
	}
	u, err := l.repo.GetByID(ctx, id)
	if err != nil {
		l.log.Error("Could not get user", "err", err.Error())
		return nil, err
	}
	exp, ok := t.Expiration()
	if !ok {
		err := &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "Invalid access token",
		}
		l.log.Error(err.Message, "err", err.Error())
		return nil, err
	}
	iat, ok := t.IssuedAt()
	if !ok {
		err := &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "Invalid access token",
		}
		l.log.Error(err.Message, "err", err.Error())
		return nil, err
	}
	iss, ok := t.Issuer()
	if !ok {
		err := &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "Invalid access token",
		}
		l.log.Error(err.Message, "err", err.Error())
		return nil, err
	}

	return &ValidateResponse{
		Subject:  u.ID,
		Expires:  exp.Unix(),
		IssuedAt: iat.Unix(),
		Issuer:   iss,
		Username: u.Username,
		Email:    u.Email,
		Role:     u.Role,
	}, nil
}

func (l *AuthLogic) buildToken(ctx context.Context, tokenType lib.TokenType, user *User) (string, error) {
	extraClaims := map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	}

	p := l.tokens.NewTokenPayload(user.ID, tokenType, extraClaims)
	token, err := l.tokens.Build(ctx, p)
	if err != nil {
		l.log.Error("Could not build token", "type", tokenType, "err", err.Error())
		return "", err
	}
	b, err := l.tokens.Sign(ctx, token)
	if err != nil {
		l.log.Error("Could not sign token", "type", tokenType, "err", err.Error())
		return "", err
	}
	return b, nil
}
