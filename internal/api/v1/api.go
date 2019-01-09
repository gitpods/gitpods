package v1

import (
	"net/http"

	"github.com/gitpods/gitpods/internal/api/v1/models"
	"github.com/gitpods/gitpods/internal/api/v1/restapi"
	"github.com/gitpods/gitpods/internal/api/v1/restapi/operations"
	"github.com/gitpods/gitpods/internal/api/v1/restapi/operations/repositories"
	"github.com/gitpods/gitpods/internal/api/v1/restapi/operations/users"
	"github.com/gitpods/gitpods/repository"
	"github.com/gitpods/gitpods/session"
	"github.com/gitpods/gitpods/user"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

// API has the http.Handler for the OpenAPI implementation
type API struct {
	Handler http.Handler
}

// New creates a new API that adds our own Handler implementations
func New(rs repository.Service, us user.Service) (*API, error) {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, err
	}

	gitpodsAPI := operations.NewGitpodsAPI(swaggerSpec)

	gitpodsAPI.RepositoriesGetOwnerRepositoriesHandler = GetOwnerRepositoriesHandler(rs)
	gitpodsAPI.RepositoriesGetRepositoryHandler = GetRepositoryHandler(rs)
	gitpodsAPI.UsersListUsersHandler = ListUsersHandler(us)
	gitpodsAPI.UsersGetUserMeHandler = GetUserMeHandler(us)
	gitpodsAPI.UsersGetUserHandler = GetUserHandler(us)
	gitpodsAPI.UsersUpdateUserHandler = UpdateUserHandler(us)

	return &API{
		Handler: gitpodsAPI.Serve(nil),
	}, nil
}

func convertRepository(r *repository.Repository) *models.Repository {
	return &models.Repository{
		ID:            strfmt.UUID(r.ID),
		Name:          &r.Name,
		Description:   r.Description,
		DefaultBranch: r.DefaultBranch,
		Website:       r.Website,
		CreatedAt:     strfmt.DateTime(r.Created),
		UpdatedAt:     strfmt.DateTime(r.Updated),
		//Owner: nil,
	}
}

func GetOwnerRepositoriesHandler(rs repository.Service) repositories.GetOwnerRepositoriesHandlerFunc {
	return func(params repositories.GetOwnerRepositoriesParams) middleware.Responder {
		list, _, err := rs.List(params.HTTPRequest.Context(), params.Owner)
		if err != nil {
			if err == repository.ErrOwnerNotFound {
				message := "owner not found"
				return repositories.NewGetOwnerRepositoriesNotFound().WithPayload(&models.Error{
					Message: &message,
				})
			}
			return repositories.NewGetOwnerRepositoriesDefault(http.StatusInternalServerError)
		}

		var payload []*models.Repository
		for _, r := range list {
			payload = append(payload, convertRepository(r))
		}

		return repositories.NewGetOwnerRepositoriesOK().WithPayload(payload)
	}
}

func GetRepositoryHandler(rs repository.Service) repositories.GetRepositoryHandlerFunc {
	return func(params repositories.GetRepositoryParams) middleware.Responder {
		r, _, err := rs.Find(params.HTTPRequest.Context(), params.Owner, params.Name)
		if err != nil {
			if err == repository.ErrRepositoryNotFound {
				message := "repository not found"
				return repositories.NewGetRepositoryNotFound().WithPayload(&models.Error{
					Message: &message,
				})
			}
			return repositories.NewGetRepositoryDefault(http.StatusInternalServerError)
		}

		return repositories.NewGetRepositoryOK().WithPayload(convertRepository(r))
	}
}

func convertUser(u *user.User) *models.User {
	return &models.User{
		ID:        strfmt.UUID(u.ID),
		Email:     strfmt.Email(u.Email),
		Username:  &u.Username,
		Name:      u.Name,
		CreatedAt: strfmt.DateTime(u.Created),
		UpdatedAt: strfmt.DateTime(u.Updated),
	}
}

// ListUsersHandler gets a list of users from the user.Service and returns a API response
func ListUsersHandler(us user.Service) users.ListUsersHandlerFunc {
	return func(params users.ListUsersParams) middleware.Responder {
		list, err := us.FindAll(params.HTTPRequest.Context())
		if err != nil {
			return users.NewListUsersDefault(http.StatusInternalServerError)
		}

		var payload []*models.User

		for _, u := range list {
			payload = append(payload, convertUser(u))
		}

		return users.NewListUsersOK().WithPayload(payload)
	}
}

// GetUserMeHandler gets the currently authenticated user
func GetUserMeHandler(us user.Service) users.GetUserMeHandlerFunc {
	return func(params users.GetUserMeParams) middleware.Responder {
		sessUser := session.GetSessionUser(params.HTTPRequest.Context())

		u, err := us.FindByUsername(params.HTTPRequest.Context(), sessUser.Username)
		if err != nil {
			return users.NewGetUserMeDefault(http.StatusInternalServerError)
		}

		return users.NewGetUserMeOK().WithPayload(convertUser(u))
	}
}

// GetUserHandler gets a user from the user.Service and returns a API response
func GetUserHandler(us user.Service) users.GetUserHandlerFunc {
	return func(params users.GetUserParams) middleware.Responder {
		u, err := us.FindByUsername(params.HTTPRequest.Context(), params.Username)
		if err != nil {
			if err == user.NotFoundError {
				message := "user not found"
				return users.NewGetUserNotFound().WithPayload(&models.Error{
					Message: &message,
				})
			}
			return users.NewGetUserDefault(http.StatusInternalServerError)
		}

		return users.NewGetUserOK().WithPayload(convertUser(u))
	}
}

// UpdateUserHandler receives a updated user and returns a API response after updating
func UpdateUserHandler(us user.Service) users.UpdateUserHandlerFunc {
	return func(params users.UpdateUserParams) middleware.Responder {
		old, err := us.FindByUsername(params.HTTPRequest.Context(), params.Username)
		if err != nil {
			if err == user.NotFoundError {
				message := "user not found"
				return users.NewUpdateUserNotFound().WithPayload(&models.Error{
					Message: &message,
				})
			}
			return users.NewUpdateUserDefault(http.StatusInternalServerError)
		}

		old.Name = *params.UpdatedUser.Name

		updated, err := us.Update(params.HTTPRequest.Context(), old)
		if err != nil {
			message := "updated user is invalid"
			return users.NewUpdateUserUnprocessableEntity().WithPayload(&models.ValidationError{
				Message: &message,
			})
		}

		return users.NewUpdateUserOK().WithPayload(convertUser(updated))
	}
}
