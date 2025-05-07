package tokens

import (
	"context"
	"errors"
	"strconv"
	"time"

	"encore.dev/rlog"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type TokenPayload struct {
	Sub string `json:"sub"`
	Nbf *int64 `json:"nbf"`

	Type        TokenType   `json:"type"`
	ExtraClaims interface{} `json:"extra_claims"`
}

type TokensUsecase interface {
	Build(ctx context.Context, p *TokenPayload) (jwt.Token, error)
	Sign(ctx context.Context, token jwt.Token) (string, error)
	Parse(ctx context.Context, token string) (jwt.Token, error)
	NewTokenPayload(sub string, typ TokenType, extraClaims interface{}) *TokenPayload
}

type tokensService struct {
	issuer      string
	alg         jwa.KeyAlgorithm
	key         []byte
	access_exp  int64
	refresh_exp int64
	log         rlog.Ctx
}

var secrets struct {
	JWTKey        string
	JWTIssuer     string
	JWTAccessExp  string
	JWTRefreshExp string
}

func NewTokensUsecase() (TokensUsecase, error) {
	rctx := rlog.With("service", "tokensService")
	uc := &tokensService{
		log: rctx,
	}

	issuer := secrets.JWTIssuer
	if issuer == "" {
		issuer = "orka"
	}
	uc.issuer = issuer

	k := secrets.JWTKey
	if k == "" {
		k = "secret"
	}
	uc.key = []byte(k)
	uc.alg = jwa.HS256()

	aStr := secrets.JWTAccessExp
	rStr := secrets.JWTRefreshExp

	access_exp, err := strconv.ParseInt(aStr, 10, 64)
	if err != nil {
		uc.log.Warn("Invalid JWT Access Expiration, using default", "err", err.Error(), "given", aStr, "default", 3600)
		access_exp = 3600
	}
	refresh_exp, err := strconv.ParseInt(rStr, 10, 64)
	if err != nil {
		uc.log.Warn("Invalid JWT Refresh Expiration, using default", "err", err.Error(), "given", rStr, "default", 24*3600)
		refresh_exp = 24 * 3600
	}

	if access_exp == 0 {
		access_exp = 3600
	}
	uc.access_exp = access_exp
	if refresh_exp == 0 {
		refresh_exp = 24 * 3600
	}
	uc.refresh_exp = refresh_exp
	return uc, nil
}

func (uc *tokensService) newBuilder() *jwt.Builder {
	return jwt.NewBuilder().Issuer(uc.issuer).IssuedAt(time.Now())
}

func (uc *tokensService) Build(ctx context.Context, p *TokenPayload) (jwt.Token, error) {
	if p == nil {
		err := errors.New("Invalid Token Payload")
		uc.log.Error("Invalid Token Payload")
		return nil, err
	}
	builder := uc.newBuilder().Subject(p.Sub)
	var exp int64
	switch p.Type {
	case AccessToken:
		exp = uc.access_exp
	case RefreshToken:
		exp = uc.refresh_exp
	default:
		return nil, errors.New("Invalid Token Type")
	}

	builder.Expiration(time.Now().Add(time.Duration(exp) * time.Second))
	if p.Nbf != nil {
		builder.NotBefore(time.Unix(*p.Nbf, 0))
	}
	if p.ExtraClaims != nil {
		for k, v := range p.ExtraClaims.(map[string]interface{}) {
			builder.Claim(k, v)
		}
	}

	return builder.Build()
}

func (uc *tokensService) Sign(ctx context.Context, token jwt.Token) (string, error) {
	opts := []jwt.SignOption{
		jwt.WithKey(uc.alg, uc.key),
	}
	b, err := jwt.Sign(token, opts...)
	if err != nil {
		// TODO: match against jwt errors
		err := errors.New("Could not sign token: " + err.Error())
		return "", err
	}
	return string(b), nil
}

// only parses the token, does not validate
func (uc *tokensService) parse(ctx context.Context, token string) (jwt.Token, error) {
	opts := []jwt.ParseOption{
		jwt.WithKey(uc.alg, uc.key),
		jwt.WithValidate(false),
	}
	t, err := jwt.Parse([]byte(token), opts...)
	if err != nil {
		// TODO: match against jwt errors
		err := errors.New("Could not parse token: " + err.Error())
		return nil, err
	}
	return t, nil
}

func (uc *tokensService) Validate(ctx context.Context, token jwt.Token) error {
	opts := []jwt.ValidateOption{
		jwt.WithIssuer(uc.issuer),
	}
	err := jwt.Validate(token, opts...)
	return err
}

func (uc *tokensService) Parse(ctx context.Context, token string) (jwt.Token, error) {
	t, err := uc.parse(ctx, token)
	if err != nil {
		return nil, err
	}
	err = uc.Validate(ctx, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *tokensService) NewTokenPayload(sub string, typ TokenType, extraClaims interface{}) *TokenPayload {
	nbf := time.Now().Unix()
	p := &TokenPayload{
		Sub:         sub,
		Type:        typ,
		Nbf:         &nbf,
		ExtraClaims: extraClaims,
	}
	return p
}
