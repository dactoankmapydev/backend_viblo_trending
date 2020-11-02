package repository

import (
	"context"
	"tech_posts_trending/model"
	"tech_posts_trending/model/req"
)

type UserRepo interface {
	CheckSignIn(context context.Context, signinReq req.ReqSignIn) (model.User, error)
	CheckEmailSignUp(context context.Context, emailReq req.ReqSignUp) (model.User, error)
	CheckEmail(context context.Context, emailReq req.ReqEmail) (model.User, error)
	UpdateUser(context context.Context, user model.User) (model.User, error)
	UpdatePassword(context context.Context, user model.User) (model.User, error)
	UpdateVerify(context context.Context, user model.User) (model.User, error)
	SaveUser(context context.Context, user model.User) (model.User, error)
	SelectUserByID(context context.Context, userID string) (model.User, error)
}
