package api

import (
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/interfaces"
	apiInterfaces "example-wikipedia-scraper/internal/interfaces/api"
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/internal/model/repository"
	"example-wikipedia-scraper/internal/service"
	"example-wikipedia-scraper/internal/service/mailer"
	paymentService "example-wikipedia-scraper/internal/service/payment"
	pkgRepo "example-wikipedia-scraper/pkg/repository"
)

type Container struct {
	cfg    config.ConfigInterface
	logger interfaces.LoggerInterface

	userRepo                     *repository.UserRepository
	pageRepo                     *repository.PageRepository
	filterRepo                   *repository.UserWantedPagesFilterRepository
	orderRepo                    *repository.OrderRepository
	orderItemRepo                pkgRepo.RepositoryInterface[*model.OrderItem]
	productRepo                  pkgRepo.RepositoryInterface[*model.Product]
	subscriptionLevelRepo        pkgRepo.RepositoryInterface[*model.SubscriptionLevel]
	subscriptionLevelProductRepo pkgRepo.RepositoryInterface[*model.SubscriptionLevelProduct]
}

func NewContainer(cfg config.ConfigInterface, logger interfaces.LoggerInterface) *Container {
	return &Container{
		cfg:                          cfg,
		logger:                       logger,
		userRepo:                     repository.NewUserRepository(),
		pageRepo:                     repository.NewPageRepository(),
		filterRepo:                   repository.NewUserWantedPagesFilterRepository(),
		orderRepo:                    repository.NewOrderRepository(),
		orderItemRepo:                pkgRepo.NewGenericRepository[*model.OrderItem](),
		productRepo:                  pkgRepo.NewGenericRepository[*model.Product](),
		subscriptionLevelRepo:        pkgRepo.NewGenericRepository[*model.SubscriptionLevel](),
		subscriptionLevelProductRepo: pkgRepo.NewGenericRepository[*model.SubscriptionLevelProduct](),
	}
}

func (c *Container) LoadModules(api apiInterfaces.ApiInterface) {
	mailerService := mailer.NewMailer(c.cfg.GetMailerConfig(), c.logger)
	userService := service.NewUserService(c.userRepo, mailerService, c.cfg)
	subscriptionService := service.NewSubscriptionService(
		c.filterRepo,
		c.subscriptionLevelRepo,
		c.subscriptionLevelProductRepo,
		c.userRepo,
	)
	pageFilterService := service.NewPageFilterService(c.filterRepo, c.pageRepo)
	paymentSvc := paymentService.NewPaymentService(c.cfg)
	orderService := service.NewOrderService(
		c.orderRepo,
		c.orderItemRepo,
		c.productRepo,
		paymentSvc,
		c.cfg,
	)

	api.LoadModule(NewUserApiModule(c.userRepo, userService, api, subscriptionService))
	api.LoadModule(NewPageApiModule(api, c.pageRepo))
	api.LoadModule(NewUserWantedFiltersApiModule(api, c.filterRepo, pageFilterService, subscriptionService))
	api.LoadModule(NewOrderApiModule(api, orderService, subscriptionService))
}
