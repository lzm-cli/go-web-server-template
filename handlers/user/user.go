package user

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/<%= organization %>/<%= repo %>/durables"
	"github.com/<%= organization %>/<%= repo %>/models"
	"github.com/<%= organization %>/<%= repo %>/session"
	"github.com/<%= organization %>/<%= repo %>/tools"
)

func AuthenticateUserByToken(ctx context.Context, authenticationToken string) (*models.User, error) {
	var user *models.User
	var queryErr error
	token, err := jwt.Parse(authenticationToken, func(token *jwt.Token) (interface{}, error) {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, session.BadDataError(ctx)
		}

		_, ok = token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, session.BadDataError(ctx)
		}
		user, queryErr = FindUserById(ctx, fmt.Sprint(claims["jti"]))
		if queryErr != nil {
			return nil, queryErr
		}
		if user == nil {
			return nil, session.BadDataError(ctx)
		}
		sum := sha256.Sum256([]byte(user.AccessToken))
		return sum[:], nil
	})

	if queryErr != nil {
		return nil, queryErr
	}
	if err != nil || !token.Valid {
		return nil, nil
	}
	return user, nil
}

func FindUserById(ctx context.Context, userId string) (*models.User, error) {
	var user models.User
	err := session.DB(ctx).First(&user, "user_id=?", userId).Error
	if durables.CheckEmptyError(err) != nil {
		return nil, err
	}
	if user.UserId == "" {
		u, err := session.MixinClient(ctx).ReadUser(ctx, userId)
		if err != nil {
			return nil, err
		}
		user = models.User{
			UserId:         u.UserID,
			IdentityNumber: u.IdentityNumber,
			FullName:       u.FullName,
			AvatarURL:      u.AvatarURL,
			CreatedAt:      u.CreatedAt,
		}
		if err := session.DB(ctx).Create(&user).Error; err != nil {
			tools.Log(err)
		}
		return nil, nil
	}
	return &user, nil
}