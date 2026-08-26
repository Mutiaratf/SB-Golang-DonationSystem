package main

import (
	"log"

	"github.com/Mutiaratf/SB-Golang-DonationSystem/internal/campaign"
	campaignupdate "github.com/Mutiaratf/SB-Golang-DonationSystem/internal/campaign_update"
	"github.com/Mutiaratf/SB-Golang-DonationSystem/internal/category"
	"github.com/Mutiaratf/SB-Golang-DonationSystem/internal/config"
	"github.com/Mutiaratf/SB-Golang-DonationSystem/internal/donor"
	"github.com/Mutiaratf/SB-Golang-DonationSystem/internal/middleware"
	"github.com/Mutiaratf/SB-Golang-DonationSystem/internal/notification"
	"github.com/Mutiaratf/SB-Golang-DonationSystem/internal/transaction"
	"github.com/Mutiaratf/SB-Golang-DonationSystem/internal/user"
	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.Load()
	if err := cfg.ValidateJWT(); err != nil {
		log.Fatalf("Invalid JWT configuration: %v", err)
	}
	db := config.InitDatabase(cfg)
	defer db.Close()

	categoryRepo := category.NewRepository(db)
	categoryService := category.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categoryService)

	campaignRepo := campaign.NewRepository(db)
	campaignService := campaign.NewService(campaignRepo)
	campaignHandler := campaign.NewHandler(campaignService)
	campaignUpdateRepo := campaignupdate.NewRepository(db)
	campaignUpdateService := campaignupdate.NewService(campaignUpdateRepo)
	campaignUpdateHandler := campaignupdate.NewHandler(campaignUpdateService)

	donorRepo := donor.NewRepository(db)
	donorService := donor.NewService(donorRepo)
	donorHandler := donor.NewHandler(donorService)

	transactionRepo := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepo, donorService, campaignService)
	transactionHandler := transaction.NewHandler(transactionService)

	notificationRepo := notification.NewRepository(db)
	notificationService := notification.NewService(notificationRepo)
	notificationHandler := notification.NewHandler(notificationService)

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService, cfg)

	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/auth/login", userHandler.Login)
		api.POST("/auth/register", userHandler.Register)

		categories := api.Group("/categories")
		{
			categories.GET("", categoryHandler.GetAllCategories)
		}

		campaigns := api.Group("/campaigns")
		{
			campaigns.GET("", campaignHandler.GetAllCampaigns)
			campaigns.GET("/:id/updates", campaignUpdateHandler.GetByCampaignID)
			campaigns.GET("/:id", campaignHandler.GetCampaignByID)
		}

		transactions := api.Group("/transactions")
		{
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.GET("/history/:id_campaign", transactionHandler.GetTransactionHistory)
		}
		api.GET("/print-transaction/:transaction_id", transactionHandler.PrintTransaction)

		protected := api.Group("")
		protected.Use(middleware.JWT(cfg.JWTSecret, cfg.JWTIssuer))
		{
			protectedCategories := protected.Group("/categories")
			protectedCategories.POST("", categoryHandler.CreateCategory)
			protectedCategories.PUT("/:id", categoryHandler.UpdateCategory)
			protectedCategories.DELETE("/:id", categoryHandler.DeleteCategory)

			protectedCampaigns := protected.Group("/campaigns")
			protectedCampaigns.POST("", campaignHandler.CreateCampaign)
			protectedCampaigns.PUT("/:id", campaignHandler.UpdateCampaign)
			protectedCampaigns.DELETE("/:id", campaignHandler.DeleteCampaign)
			protectedCampaigns.POST("/:id/updates", campaignUpdateHandler.Create)
			protectedCampaigns.PUT("/:id/updates/:update_id", campaignUpdateHandler.Update)
			protectedCampaigns.DELETE("/:id/updates/:update_id", campaignUpdateHandler.Delete)

			protectedDonors := protected.Group("/donors")
			protectedDonors.POST("", donorHandler.CreateDonor)
			protectedDonors.GET("", donorHandler.GetAllDonors)
			protectedDonors.GET("/:id", donorHandler.GetDonorByID)
			protectedDonors.PUT("/:id", donorHandler.UpdateDonor)
			protectedDonors.DELETE("/:id", donorHandler.DeleteDonor)

			protectedTransactions := protected.Group("/transactions")
			protectedTransactions.GET("", transactionHandler.GetAllTransactions)
			protectedTransactions.PUT("/:id", transactionHandler.UpdateTransaction)
			protectedTransactions.DELETE("/:id", transactionHandler.DeleteTransaction)
			protectedTransactions.POST("/notifications/:id_transaksi", notificationHandler.SendTransactionNotification)
		}

	}

	err := r.Run(":8082")
	if err != nil {
		log.Fatalf("Failed run server: %v", err)
	}
}
