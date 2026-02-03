package handlers

import (
	"net/http"
	// "strings"

	"subscription-service/internal/models"
	"subscription-service/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Handler struct {
	repo   *repository.SubscriptionRepository
	logger *zap.SugaredLogger
}


func NewHandler(repo *repository.SubscriptionRepository, logger *zap.SugaredLogger) *Handler {
	return &Handler{repo: repo, logger: logger}
}









func (h *Handler) CreateSubscription(c *gin.Context) {
	var input models.Subscription
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Warnf("Invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Create(&input); err != nil {
		h.logger.Errorf("Failed to create subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}

	h.logger.Infof("Subscription created: %s", input.ID)
	c.JSON(http.StatusCreated, input)
}







func (h *Handler) GetAllSubscriptions(c *gin.Context) {
	subs, err := h.repo.GetAll()
	if err != nil {
		h.logger.Errorf("Failed to fetch subscriptions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	c.JSON(http.StatusOK, subs)
}








func (h *Handler) GetSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	sub, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sub)
}










func (h *Handler) UpdateSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var input models.Subscription
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.ID = id 

	if err := h.repo.Update(&input); err != nil {
		h.logger.Warnf("Failed to update %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	h.logger.Infof("Subscription updated: %s", id)
	c.JSON(http.StatusOK, input)
}







func (h *Handler) DeleteSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	if err := h.repo.Delete(id); err != nil {
		h.logger.Warnf("Failed to delete %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Infof("Subscription deleted: %s", id)
	c.Status(http.StatusNoContent)
}


func (h *Handler) GetTotalCost(c *gin.Context) {
	userIDStr := c.Query("user_id")
	serviceName := c.Query("service_name")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if userIDStr == "" || startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required params"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	
	
	
	var start, end models.CustomDate
	
	if err := start.UnmarshalJSON([]byte("\"" + startDateStr + "\"")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date"})
		return
	}
	if err := end.UnmarshalJSON([]byte("\"" + endDateStr + "\"")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date"})
		return
	}

	filter := repository.CostFilter{
		UserID:      userID,
		ServiceName: serviceName,
		StartDate:   start,
		EndDate:     end,
	}

	total, err := h.repo.GetTotalCost(filter)
	if err != nil {
		h.logger.Errorf("Cost calc error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_cost": total})
}
