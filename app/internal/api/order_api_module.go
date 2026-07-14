package api

import (
	apiMiddleware "example-wikipedia-scraper/internal/api/middleware"
	"example-wikipedia-scraper/internal/dto"
	apiInterfaces "example-wikipedia-scraper/internal/interfaces/api"
	serviceInterfaces "example-wikipedia-scraper/internal/interfaces/service"
	"example-wikipedia-scraper/internal/queue"
	orderHandlers "example-wikipedia-scraper/internal/queue/handlers/order"
	types "example-wikipedia-scraper/internal/types/api"
	orderTypes "example-wikipedia-scraper/internal/types/order"
)

type OrderApiModule struct {
	api          apiInterfaces.ApiInterface
	orderService serviceInterfaces.OrderServiceInterface
	subService   serviceInterfaces.SubscriptionServiceInterface
}

func NewOrderApiModule(api apiInterfaces.ApiInterface, orderService serviceInterfaces.OrderServiceInterface, subService serviceInterfaces.SubscriptionServiceInterface) *OrderApiModule {
	module := &OrderApiModule{
		api:          api,
		orderService: orderService,
		subService:   subService,
	}
	module.registerOrderQueueHandlers()
	return module
}

func (a *OrderApiModule) registerOrderQueueHandlers() {
	orderQueueHandler := orderHandlers.NewOrderQueueHandler(a.orderService, a.subService, a.api.GetLogger())
	queue.GetMessageQueueService().RegisterHandler("order_payment_notfied", orderQueueHandler.HandleOrderPaymentNotfied)
}

func (a *OrderApiModule) GetRoutes() []*types.Route {
	return []*types.Route{
		{
			Method:           "POST",
			Path:             "/",
			Handler:          a.processOrder,
			Middlewares:      []types.Middleware{a.api.AuthenticateMiddleware, apiMiddleware.IdempotentMiddleware},
			AfterMiddlewares: []types.AfterMiddleware{apiMiddleware.CacheIdempotentResponse},
		},
		{
			Method:      "GET",
			Path:        "/payment_methods",
			Handler:     a.getPaymentMethods,
			Middlewares: []types.Middleware{a.api.AuthenticateMiddleware},
		},
		{
			Method:  "PATCH",
			Path:    "/payment/:method/notify",
			Handler: a.paymentNotify,
		},
		{
			Method:  "GET",
			Path:    "/example_payment/:paymentId",
			Handler: a.payExamplePayment,
		},
	}
}

func (a *OrderApiModule) GetRoutePrefix() string {
	return "order"
}

func (a *OrderApiModule) processOrder(request *types.ApiRequest) *types.ApiResponse {
	errApiResp, dto := ValidateAndUnmarshalDTO[dto.OrderCreateDto](request.Body)
	dto.ClientIP = request.ClientIP
	dto.UserID = request.User.ID
	if errApiResp != nil {
		return errApiResp
	}
	// tu wywalic te errory
	orderWithItems, err := a.orderService.CreateOrder(dto)
	if err != nil {
		return BadRequestResponseWithMsg(err.Error())
	}
	paymentProcessorResponse, err := a.orderService.ProcessPayment(orderWithItems)
	if err != nil {
		return BadRequestResponseWithMsg(err.Error())
	}
	response := &orderTypes.OrderCreated{
		ID:          orderWithItems.Order.ID,
		Status:      orderWithItems.Order.Status,
		PaymentId:   paymentProcessorResponse.PaymentID,
		RedirectURL: paymentProcessorResponse.RedirectURL,
	}
	return &types.ApiResponse{
		StatusCode: 201,
		Body: &types.ApiResponseBody{
			Success: true,
			Data:    response,
		},
	}
}

func (a *OrderApiModule) getPaymentMethods(request *types.ApiRequest) *types.ApiResponse {
	paymentMethods := a.orderService.GetAvailablePaymentMethods()
	return &types.ApiResponse{
		StatusCode: 200,
		Body: &types.ApiResponseBody{
			Success: true,
			Data:    paymentMethods,
		},
	}
}

func (a *OrderApiModule) paymentNotify(request *types.ApiRequest) *types.ApiResponse {
	orderPaymentNotification := &orderTypes.OrderPaymentNotification{
		ServiceIp:     request.ClientIP,
		Body:          request.Body,
		PaymentMethod: request.PathParams["method"],
		Headers:       request.Headers,
	}
	err := a.orderService.HandlePaymentNotification(orderPaymentNotification)
	if err != nil {
		return BadRequestResponseWithMsg(err.Error())
	}
	return &types.ApiResponse{
		StatusCode: 200,
		Body: &types.ApiResponseBody{
			Success: true,
		},
	}
}

func (a *OrderApiModule) payExamplePayment(request *types.ApiRequest) *types.ApiResponse {
	requestCpy := *request
	requestCpy.Body = `{"payment_id":"` + request.PathParams["paymentId"] + `", "action":"accept"}`
	requestCpy.PathParams = map[string]string{"method": "example"}
	requestCpy.Method = "PATCH"
	requestCpy.URL = "/order/payment/example/notify"
	result := a.paymentNotify(&requestCpy)
	if result.StatusCode != 200 {
		return BadRequestResponseWithMsg("Failed to process example payment")
	}
	return &types.ApiResponse{
		StatusCode: 302,
		Headers: map[string]string{
			"Location": "/dashboard/panel/user",
		},
	}
}
