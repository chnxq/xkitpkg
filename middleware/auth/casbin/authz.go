package casbin

import (
	"context"
	"errors"
	"fmt"

	jwtV5 "github.com/golang-jwt/jwt/v5"

	"github.com/chnxq/xkitpkg/middleware/auth/jwt"
	"github.com/chnxq/xkitpkg/transport"
	transportHttp "github.com/chnxq/xkitpkg/transport/http"
)

type ISecurityUser interface {
	// ParseFromContext parses the user from the context.
	ParseFromContext(ctx context.Context) error
	// GetSubject returns the subject of the token.
	GetSubject() string
	// GetObject returns the object of the token.
	GetObject() string
	// GetAction returns the action of the token.
	GetAction() string
	// GetDomain returns the domain of the token.
	GetDomain() string
}

type SecurityUserCreator func() ISecurityUser

const (
	ClaimAuthorityId = "roleId"
	Domain           = "domain"
)

type SecurityUser struct {
	Path        string
	Domain      string
	Method      string
	AuthorityId string
}

func NewSecurityUser() ISecurityUser {
	return &SecurityUser{}
}

func (su *SecurityUser) ParseFromContext(ctx context.Context) error {
	err := su.ParseAccessJwtTokenFromContext(ctx)
	if err != nil {
		return err
	}

	if header, ok := transport.FromServerContext(ctx); ok {
		http, httpOk := header.(transportHttp.Transporter)
		if !httpOk {
			return errors.New("no http transporter")
		}
		httpRequest := http.Request()
		su.Path = httpRequest.URL.Path
		su.Method = httpRequest.Method
	} else {
		return errors.New("jwt claim missing")
	}

	return nil
}

func (su *SecurityUser) GetSubject() string {
	return su.AuthorityId
}

func (su *SecurityUser) GetObject() string {
	return su.Path
}

func (su *SecurityUser) GetAction() string {
	return su.Method
}

func (su *SecurityUser) GetDomain() string {
	return su.Domain
}

func (su *SecurityUser) CreateAccessJwtToken(secretKey []byte) string {
	claims := jwtV5.NewWithClaims(jwtV5.SigningMethodHS256,
		jwtV5.MapClaims{
			ClaimAuthorityId: su.AuthorityId,
		})

	signedToken, err := claims.SignedString(secretKey)
	if err != nil {
		return ""
	}

	return signedToken
}

func (su *SecurityUser) ParseAccessJwtTokenFromContext(ctx context.Context) error {
	claims, ok := jwt.FromContext(ctx)
	if !ok {
		fmt.Println("ParseAccessJwtTokenFromContext 1")
		return errors.New("no jwt token in context")
	}
	if err := su.ParseAccessJwtToken(claims); err != nil {
		fmt.Println("ParseAccessJwtTokenFromContext 2")
		return err
	}
	return nil
}

func (su *SecurityUser) ParseAccessJwtTokenFromString(token string, secretKey []byte) error {
	parseAuth, err := jwtV5.Parse(token, func(*jwtV5.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return err
	}

	claims, ok := parseAuth.Claims.(jwtV5.MapClaims)
	if !ok {
		return errors.New("no jwt token in context")
	}

	if err = su.ParseAccessJwtToken(claims); err != nil {
		return err
	}

	return nil
}

func (su *SecurityUser) ParseAccessJwtToken(claims jwtV5.Claims) error {
	if claims == nil {
		return errors.New("claims is nil")
	}

	mc, ok := claims.(jwtV5.MapClaims)
	if !ok {
		return errors.New("claims is not map claims")
	}

	authorityIdStr, authorityIdOk := mc[ClaimAuthorityId]
	if !authorityIdOk {
		return errors.New("authorityId is missing")
	}
	su.AuthorityId = authorityIdStr.(string)
	domainStr, domainOk := mc[Domain]
	if domainOk {
		su.Domain = domainStr.(string)
	}

	return nil
}
