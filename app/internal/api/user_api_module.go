package api

import (
	"errors"
	apiMiddleware "example-wikipedia-scraper/internal/api/middleware"
	"example-wikipedia-scraper/internal/dto"
	apiInterfaces "example-wikipedia-scraper/internal/interfaces/api"
	"example-wikipedia-scraper/internal/interfaces/auth"
	"example-wikipedia-scraper/internal/interfaces/repository"
	serviceInterfaces "example-wikipedia-scraper/internal/interfaces/service"
	"example-wikipedia-scraper/internal/model"
	types "example-wikipedia-scraper/internal/types/api"
	"net/http"
	"slices"
	"strconv"
)

type UserApiModule struct {
	userRepo            repository.UserRepositoryInterface
	userService         serviceInterfaces.UserServiceInterface
	api                 apiInterfaces.ApiInterface
	subscriptionService serviceInterfaces.SubscriptionServiceInterface
}

func NewUserApiModule(
	userRepo repository.UserRepositoryInterface,
	userService serviceInterfaces.UserServiceInterface,
	api apiInterfaces.ApiInterface,
	subscriptionService serviceInterfaces.SubscriptionServiceInterface,
) *UserApiModule {
	return &UserApiModule{
		userRepo:            userRepo,
		userService:         userService,
		api:                 api,
		subscriptionService: subscriptionService,
	}
}

func (a *UserApiModule) GetRoutes() []*types.Route {
	return []*types.Route{
		{
			Method:           "POST",
			Path:             "/",
			Handler:          a.createUser,
			Middlewares:      []types.Middleware{apiMiddleware.IdempotentMiddleware},
			AfterMiddlewares: []types.AfterMiddleware{apiMiddleware.CacheIdempotentResponse},
		},
		{
			Method:      "GET",
			Path:        "/:id",
			Handler:     a.getUser,
			Middlewares: []types.Middleware{a.api.AuthenticateMiddleware, a.AuthorizeUserMiddleware},
		},
		{
			Method:      "DELETE",
			Path:        "/:id",
			Handler:     a.deleteUser,
			Middlewares: []types.Middleware{a.api.AuthenticateMiddleware, a.AuthorizeUserMiddleware},
		},
		{
			Method:           "POST",
			Path:             "/login",
			Handler:          a.loginUser,
			Middlewares:      []types.Middleware{apiMiddleware.IdempotentMiddleware},
			AfterMiddlewares: []types.AfterMiddleware{apiMiddleware.CacheIdempotentResponse},
		},
		{
			Method:      "POST",
			Path:        "/change_email",
			Handler:     a.changeEmail,
			Middlewares: []types.Middleware{a.api.AuthenticateMiddleware},
		},
		{
			Method:      "POST",
			Path:        "/change_password",
			Handler:     a.changePassword,
			Middlewares: []types.Middleware{a.api.AuthenticateMiddleware},
		},
		{
			Method:  "POST",
			Path:    "/verify_email",
			Handler: a.verifyEmail,
		},
		{
			Method:  "POST",
			Path:    "/forgot_password",
			Handler: a.forgotPassword,
		},
		{
			Method:  "POST",
			Path:    "/reset_password",
			Handler: a.resetPassword,
		},
		{
			Method:      "POST",
			Path:        "/logout_all",
			Handler:     a.logoutAll,
			Middlewares: []types.Middleware{a.api.AuthenticateMiddleware},
		},
		{
			Method:      "GET",
			Path:        "/subscription_levels",
			Handler:     a.getSubscriptionLevelProducts,
			Middlewares: []types.Middleware{a.api.AuthenticateMiddleware},
		},
		{
			Method:      "GET",
			Path:        "/subscription_levels/:level",
			Handler:     a.getSubscriptionLevel,
			Middlewares: []types.Middleware{a.api.AuthenticateMiddleware},
		},
	}
}

func (a *UserApiModule) GetRoutePrefix() string {
	return "user"
}

func (a *UserApiModule) createUser(request *types.ApiRequest) *types.ApiResponse {
	errApiResp, dto := ValidateAndUnmarshalDTO[dto.CreateUserDTO](request.Body)
	if errApiResp != nil {
		return errApiResp
	}
	if user, _ := a.userRepo.GetByEmail(dto.Email); user.ID != 0 {
		return &types.ApiResponse{
			StatusCode: http.StatusConflict,
			Body: &types.ApiResponseBody{
				Success: false,
				Errors:  []types.ErrorDetail{{Message: "Email already in use", Field: "email"}},
			},
		}
	}
	if user, _ := a.userRepo.GetByUsername(dto.Username); user.ID != 0 {
		return &types.ApiResponse{
			StatusCode: http.StatusConflict,
			Body: &types.ApiResponseBody{
				Success: false,
				Errors:  []types.ErrorDetail{{Message: "Username already in use", Field: "username"}},
			},
		}
	}
	user, err := a.userService.CreateUser(*dto)
	if err != nil {
		return InternalErrorResponse()
	}
	return &types.ApiResponse{
		StatusCode: http.StatusCreated,
		Body: &types.ApiResponseBody{
			Success: true,
			Data:    user,
		},
	}
}

func (a *UserApiModule) loginUser(request *types.ApiRequest) *types.ApiResponse {
	errApiResp, dto := ValidateAndUnmarshalDTO[dto.LoginUserDTO](request.Body)
	if errApiResp != nil {
		return errApiResp
	}
	req := *dto
	token, err := a.api.GetAuthManager().Login(req.Login, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return &types.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Body: &types.ApiResponseBody{
					Success: false,
					Errors:  []types.ErrorDetail{{Message: "Invalid login or password"}},
				},
			}
		}
		if errors.Is(err, auth.ErrUserNotVerified) {
			return ForbiddenResponseWithMsg("Email not verified")
		}
		return InternalErrorResponse()
	}
	return &types.ApiResponse{
		StatusCode: http.StatusOK,
		Body: &types.ApiResponseBody{
			Success: true,
			Data:    map[string]string{"token": token},
		},
	}
}

func (a *UserApiModule) getUser(request *types.ApiRequest) *types.ApiResponse {
	a.api.GetLogger().Info("getUser called")
	idParam := request.PathParams["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return &types.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Body: &types.ApiResponseBody{
				Success: false,
				Errors:  []types.ErrorDetail{{Message: "Invalid user ID"}},
			},
		}
	}
	user, err := a.userRepo.GetByID(uint(id))
	if err != nil {
		return InternalErrorResponse()
	}
	if user.ID == 0 {
		return &types.ApiResponse{
			StatusCode: http.StatusNotFound,
			Body: &types.ApiResponseBody{
				Success: false,
				Errors:  []types.ErrorDetail{{Message: "User not found"}},
			},
		}
	}
	return &types.ApiResponse{
		StatusCode: http.StatusOK,
		Body: &types.ApiResponseBody{
			Success: true,
			Data:    user,
		},
	}
}

func (a *UserApiModule) deleteUser(request *types.ApiRequest) *types.ApiResponse {
	idParam := request.PathParams["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return &types.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Body: &types.ApiResponseBody{
				Success: false,
				Errors:  []types.ErrorDetail{{Message: "Invalid user ID"}},
			},
		}
	}
	err = a.userRepo.Disable(uint(id))
	if err != nil {
		return InternalErrorResponse()
	}
	return &types.ApiResponse{
		StatusCode: http.StatusNoContent,
		Body: &types.ApiResponseBody{
			Success: true,
		},
	}
}

func (a *UserApiModule) AuthorizeUserMiddleware(request *types.ApiRequest) *types.ApiResponse {
	userObj := request.User
	if userObj == nil {
		a.api.GetLogger().Debug("User not found in context")
		return &types.ApiResponse{
			StatusCode: http.StatusUnauthorized,
			Body: &types.ApiResponseBody{
				Success: false,
				Errors:  []types.ErrorDetail{{Message: "Unauthorized"}},
			},
		}
	}
	if userObj.Role == model.RoleAdmin {
		return nil
	}
	userIDStr := request.PathParams["id"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		return &types.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Body: &types.ApiResponseBody{
				Success: false,
				Errors:  []types.ErrorDetail{{Message: "Invalid user ID"}},
			},
		}
	}
	if userObj.ID == uint(userID) {
		return nil
	}
	return ForbiddenResponse()
}

func (a *UserApiModule) changeEmail(request *types.ApiRequest) *types.ApiResponse {
	user := request.User
	if user == nil {
		return &types.ApiResponse{
			StatusCode: http.StatusUnauthorized,
			Body: &types.ApiResponseBody{
				Success: false,
				Errors:  []types.ErrorDetail{{Message: "Unauthorized"}},
			},
		}
	}
	errApiResp, dto := ValidateAndUnmarshalDTO[dto.ChangeEmailDTO](request.Body)
	if errApiResp != nil {
		return errApiResp
	}
	_, err := a.userService.ChangeEmail(user, *dto)
	if err != nil {
		return InternalErrorResponse()
	}
	return &types.ApiResponse{
		StatusCode: http.StatusOK,
		Body: &types.ApiResponseBody{
			Success: true,
		},
	}
}

func (a *UserApiModule) verifyEmail(request *types.ApiRequest) *types.ApiResponse {
	errApiResp, dto := ValidateAndUnmarshalDTO[dto.VerifyEmailDTO](request.Body)
	if errApiResp != nil {
		return errApiResp
	}
	err := a.userService.VerifyEmail(dto.Token)
	if err != nil {
		if errors.Is(err, serviceInterfaces.ErrInvalidVerificationToken) {
			return &types.ApiResponse{
				StatusCode: http.StatusNotFound,
				Body: &types.ApiResponseBody{
					Success: false,
					Errors:  []types.ErrorDetail{{Message: "Invalid verification token"}},
				},
			}
		}
		return InternalErrorResponse()
	}
	return &types.ApiResponse{
		StatusCode: http.StatusOK,
		Body: &types.ApiResponseBody{
			Success: true,
		},
	}
}

func (a *UserApiModule) changePassword(request *types.ApiRequest) *types.ApiResponse {
	user := request.User
	errApiResp, dto := ValidateAndUnmarshalDTO[dto.ChangePasswordDTO](request.Body)
	if errApiResp != nil {
		return errApiResp
	}
	err := a.userService.ChangePassword(user, *dto)
	if err != nil {
		if errors.Is(err, serviceInterfaces.ErrInvalidPassword) {
			return &types.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Body: &types.ApiResponseBody{
					Success: false,
					Errors:  []types.ErrorDetail{{Message: "Invalid current password"}},
				},
			}
		}
		return InternalErrorResponse()
	}
	return &types.ApiResponse{
		StatusCode: http.StatusOK,
		Body: &types.ApiResponseBody{
			Success: true,
		},
	}
}

func (a *UserApiModule) forgotPassword(request *types.ApiRequest) *types.ApiResponse {
	errApiResp, dto := ValidateAndUnmarshalDTO[dto.ForgotPasswordDTO](request.Body)
	if errApiResp != nil {
		return errApiResp
	}
	err := a.userService.ForgotPassword(*dto)
	if err != nil {
		return InternalErrorResponse()
	}
	return &types.ApiResponse{
		StatusCode: http.StatusOK,
		Body: &types.ApiResponseBody{
			Success: true,
		},
	}
}

func (a *UserApiModule) resetPassword(request *types.ApiRequest) *types.ApiResponse {
	errApiResp, dto := ValidateAndUnmarshalDTO[dto.ResetPasswordDTO](request.Body)
	if errApiResp != nil {
		return errApiResp
	}
	err := a.userService.ResetPassword(*dto)
	if err != nil {
		if errors.Is(err, serviceInterfaces.ErrInvalidResetToken) {
			return &types.ApiResponse{
				StatusCode: http.StatusNotFound,
				Body: &types.ApiResponseBody{
					Success: false,
					Errors:  []types.ErrorDetail{{Message: "Invalid reset token"}},
				},
			}
		}
		if errors.Is(err, serviceInterfaces.ErrExpiredResetToken) {
			return &types.ApiResponse{
				StatusCode: http.StatusBadRequest,
				Body: &types.ApiResponseBody{
					Success: false,
					Errors:  []types.ErrorDetail{{Message: "Expired reset token"}},
				},
			}
		}
		return InternalErrorResponse()
	}
	return &types.ApiResponse{
		StatusCode: http.StatusOK,
		Body: &types.ApiResponseBody{
			Success: true,
		},
	}
}

func (a *UserApiModule) logoutAll(request *types.ApiRequest) *types.ApiResponse {
	err := a.userService.Logout(request.User)
	if err != nil {
		return InternalErrorResponse()
	}
	return &types.ApiResponse{
		StatusCode: http.StatusOK,
		Body: &types.ApiResponseBody{
			Success: true,
		},
	}
}

func (a *UserApiModule) getSubscriptionLevelProducts(request *types.ApiRequest) *types.ApiResponse {
	subscriptionLevelProducts, err := a.subscriptionService.GetSubscriptionLevelProducts()
	if err != nil {
		return InternalErrorResponse()
	}
	result := make(map[uint]map[uint]*model.SubscriptionLevelProduct)
	for _, product := range subscriptionLevelProducts {
		if product.SubscriptionLevel.Level == 0 {
			continue
		}
		id := uint(product.SubscriptionLevel.Level)
		resultMap := result[id]
		if resultMap == nil {
			resultMap = make(map[uint]*model.SubscriptionLevelProduct)
			result[id] = resultMap
		}
		intDsc, err := strconv.Atoi(product.Product.Description)
		if err != nil {
			return InternalErrorResponse()
		}
		resultMap[uint(intDsc)] = product
	}
	sortedResult := make(map[uint][]*model.SubscriptionLevelProduct)
	for levelID, productsMap := range result {
		var keys []uint
		for k := range productsMap {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		var sortedProducts []*model.SubscriptionLevelProduct
		for _, k := range keys {
			sortedProducts = append(sortedProducts, productsMap[k])
		}
		sortedResult[levelID] = sortedProducts
	}
	return &types.ApiResponse{
		StatusCode: http.StatusOK,
		Body: &types.ApiResponseBody{
			Success: true,
			Data:    sortedResult,
		},
	}
}

func (a *UserApiModule) getSubscriptionLevel(request *types.ApiRequest) *types.ApiResponse {
	subscriptionLevelStr := request.PathParams["level"]
	subscriptionLevel, err := strconv.Atoi(subscriptionLevelStr)
	if err != nil {
		return BadRequestResponseWithMsg("Invalid subscription level")
	}
	subscriptionLevelData, err := a.subscriptionService.GetRepository().FindOneBy("level = ?", subscriptionLevel)
	if err != nil {
		return InternalErrorResponse()
	}
	if subscriptionLevelData == nil {
		return NotFoundResponse("Subscription level not found")
	}
	return &types.ApiResponse{
		StatusCode: http.StatusOK,
		Body: &types.ApiResponseBody{
			Success: true,
			Data:    subscriptionLevelData,
		},
	}
}
