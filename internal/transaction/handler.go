package transaction

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	receiptpdf "github.com/Mutiaratf/SB-Golang-DonationSystem/internal/pdf"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) CreateTransaction(c *gin.Context) {
	var request TransactionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		transactionError(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	item, err := h.service.Create(&request)
	if err != nil {
		handleError(c, err, "Failed to create transaction")
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Transaction created successfully",
		"data":    item,
	})
}

func (h *Handler) GetAllTransactions(c *gin.Context) {
	items, err := h.service.GetAll()
	if err != nil {
		transactionError(c, http.StatusInternalServerError, "Failed to retrieve transactions")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Transactions retrieved successfully",
		"data":    items,
	})
}

func (h *Handler) UpdateTransaction(c *gin.Context) {
	id, ok := transactionID(c)
	if !ok {
		return
	}
	var request TransactionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		transactionError(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	item, err := h.service.Update(id, &request)
	if err != nil {
		handleError(c, err, "Failed to update transaction")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Transaction updated successfully",
		"data":    item,
	})
}

func (h *Handler) DeleteTransaction(c *gin.Context) {
	id, ok := transactionID(c)
	if !ok {
		return
	}
	if err := h.service.Delete(id); err != nil {
		handleError(c, err, "Failed to delete transaction")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Transaction deleted successfully",
	})
}

func (h *Handler) GetTransactionHistory(c *gin.Context) {
	campaignID, err := strconv.ParseInt(c.Param("id_campaign"), 10, 64)
	if err != nil || campaignID <= 0 {
		transactionError(c, http.StatusBadRequest, "Invalid campaign ID")
		return
	}
	items, err := h.service.GetHistory(campaignID)
	if err != nil {
		transactionError(c, http.StatusInternalServerError, "Failed to retrieve transaction history")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Transaction history retrieved successfully",
		"data":    items,
	})
}

func (h *Handler) PrintTransaction(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("transaction_id"), 10, 64)
	if err != nil || id <= 0 {
		transactionError(c, http.StatusBadRequest, "Invalid transaction ID")
		return
	}
	receipt, err := h.service.GetReceipt(id)
	if err != nil {
		handleError(c, err, "Failed to retrieve transaction")
		return
	}
	pdfBytes, err := receiptpdf.GenerateReceipt(receipt, time.Now())
	if err != nil {
		transactionError(c, http.StatusInternalServerError, "Failed to generate donation receipt")
		return
	}
	c.Header("Content-Disposition", `inline; filename="bukti-donasi-`+strconv.FormatInt(receipt.ID, 10)+`.pdf"`)
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

func transactionID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		transactionError(c, http.StatusBadRequest, "Invalid transaction ID")
		return 0, false
	}
	return id, true
}

func handleError(c *gin.Context, err error, fallback string) {
	switch {
	case errors.Is(err, ErrTransactionNotFound):
		transactionError(c, http.StatusNotFound, "Transaction not found")
	case errors.Is(err, ErrCampaignInactive):
		transactionError(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrAmountTooLow), errors.Is(err, ErrTransactionAmount), errors.Is(err, ErrDonorDataRequired),
		errors.Is(err, ErrDonorGenderInvalid), errors.Is(err, ErrDonorPhoneInvalid), errors.Is(err, ErrDonorEmailInvalid):
		transactionError(c, http.StatusBadRequest, err.Error())
	default:
		transactionError(c, http.StatusInternalServerError, fallback)
	}
}

func transactionError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"status":  "error",
		"message": message,
	})
}
