package server

import (
	"github.com/Applessr/hello-sekai-shop-tutorial/modules/payment/paymentHandler"
	"github.com/Applessr/hello-sekai-shop-tutorial/modules/payment/paymentRepository"
	"github.com/Applessr/hello-sekai-shop-tutorial/modules/payment/paymentUsecase"
)

func (s *server) paymentService() {
	repo := paymentRepository.NewPaymentRepository(s.db)
	usecase := paymentUsecase.NewPaymentUsecase(repo)
	httpHandler := paymentHandler.NewPaymentHttpHandler(s.cfg, usecase)
	// queueHandler := paymentHandler.NewPaymentQueueHandler(s.cfg, usecase)

	_ = httpHandler

	payment := s.app.Group("/payment_v1")

	payment.GET("", s.healthCheckService)

	payment.POST("/payment/buy", httpHandler.BuyItem, s.middleware.JwtAuthorization)
	payment.POST("/payment/sell", httpHandler.SellItem, s.middleware.JwtAuthorization)
}
