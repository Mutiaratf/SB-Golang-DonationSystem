package main

import (
	"log"

	"github.com/Mutiaratf/SB-Golang-DonationSystem/internal/campaign"
	campaignupdate "github.com/Mutiaratf/SB-Golang-DonationSystem/internal/campaign_update"
	"github.com/Mutiaratf/SB-Golang-DonationSystem/internal/category"
	"github.com/Mutiaratf/SB-Golang-DonationSystem/internal/config"
	"github.com/Mutiaratf/SB-Golang-DonationSystem/internal/donor"
	"github.com/Mutiaratf/SB-Golang-DonationSystem/internal/transaction"
	"github.com/gin-gonic/gin"
)

func main() {

	db := config.InitDatabase(config.Load())
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

	r := gin.Default()

	api := r.Group("/api")
	{
		categories := api.Group("/categories")
		{
			categories.POST("", categoryHandler.CreateCategory)
			categories.GET("", categoryHandler.GetAllCategories)
			categories.PUT("/:id", categoryHandler.UpdateCategory)
			categories.DELETE("/:id", categoryHandler.DeleteCategory)
		}

		campaigns := api.Group("/campaigns")
		{
			campaigns.POST("", campaignHandler.CreateCampaign)
			campaigns.GET("", campaignHandler.GetAllCampaigns)
			campaigns.GET("/:id/updates", campaignUpdateHandler.GetByCampaignID)
			campaigns.GET("/:id", campaignHandler.GetCampaignByID)
			campaigns.PUT("/:id", campaignHandler.UpdateCampaign)
			campaigns.DELETE("/:id", campaignHandler.DeleteCampaign)
		}

		donors := api.Group("/donors")
		{
			donors.POST("", donorHandler.CreateDonor)
			donors.GET("", donorHandler.GetAllDonors)
			donors.GET("/:id", donorHandler.GetDonorByID)
			donors.PUT("/:id", donorHandler.UpdateDonor)
			donors.DELETE("/:id", donorHandler.DeleteDonor)
		}

		transactions := api.Group("/transactions")
		{
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.GET("", transactionHandler.GetAllTransactions)
			transactions.PUT("/:id", transactionHandler.UpdateTransaction)
			transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
			transactions.GET("/history/:id_campaign", transactionHandler.GetTransactionHistory)
		}

	}

	log.Println("Server berjalan di http://localhost:8082")
	err := r.Run(":8082")
	if err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
