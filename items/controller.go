package items

import (
	"go-api/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Service Service
}

func NewController(service Service) Controller {
	return Controller{Service: service}
}

func generateMessage(action string, success bool) string {
	if success {
		return action + " successfully"
	}
	return "Failed to " + action
}

func sendSuccess(c *gin.Context, action string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"message": generateMessage(action, true), "data": data})
}

func sendError(c *gin.Context, status int, action string) {
	c.JSON(status, gin.H{"message": generateMessage(action, false)})
}

func (ctrl *Controller) GetItems(c *gin.Context) {
	items, err := ctrl.Service.GetItems()
	if err != nil {
		sendError(c, http.StatusInternalServerError, "retrieve items")
		return
	}
	sendSuccess(c, "retrieve items", items)
}

func (ctrl *Controller) GetItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(c, http.StatusBadRequest, "parse ID")
		return
	}
	item, err := ctrl.Service.GetItemByID(uint(id))
	if err != nil {
		sendError(c, http.StatusNotFound, "find item")
		return
	}
	sendSuccess(c, "retrieve item", item)
}

func (ctrl *Controller) CreateItem(c *gin.Context) {
	var item model.Item
	if err := c.BindJSON(&item); err != nil {
		sendError(c, http.StatusBadRequest, "parse request")
		return
	}
	item.Status = "PENDING"
	if err := ctrl.Service.CreateItem(&item); err != nil {
		sendError(c, http.StatusInternalServerError, "create item")
		return
	}
	sendSuccess(c, "create item", item)
}

func (ctrl *Controller) UpdateItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(c, http.StatusBadRequest, "parse ID")
		return
	}
	var updatedItem model.Item
	if err := c.BindJSON(&updatedItem); err != nil {
		sendError(c, http.StatusBadRequest, "parse request")
		return
	}
	updatedItem.ID = uint(id)
	updatedItem.Status = "PENDING"
	if err := ctrl.Service.UpdateItem(updatedItem); err != nil {
		sendError(c, http.StatusInternalServerError, "update item")
		return
	}
	sendSuccess(c, "update item", updatedItem)
}

func (ctrl *Controller) PatchItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(c, http.StatusBadRequest, "parse ID")
		return
	}
	var updatedFields map[string]interface{}
	if err := c.BindJSON(&updatedFields); err != nil {
		sendError(c, http.StatusBadRequest, "parse request")
		return
	}
	item, err := ctrl.Service.GetItemByID(uint(id))
	if err != nil {
		sendError(c, http.StatusNotFound, "find item")
		return
	}
	if title, ok := updatedFields["title"].(string); ok {
		item.Title = title
	}
	if amount, ok := updatedFields["amount"].(float64); ok {
		item.Amount = int(amount)
	}
	if quantity, ok := updatedFields["quantity"].(float64); ok {
		item.Quantity = int(quantity)
	}
	if status, ok := updatedFields["status"].(string); ok {
		item.Status = status
	}
	if ownerID, ok := updatedFields["owner_id"].(float64); ok {
		item.OwnerID = uint(ownerID)
	}
	if err := ctrl.Service.UpdateItem(item); err != nil {
		sendError(c, http.StatusInternalServerError, "update item")
		return
	}
	sendSuccess(c, "update item", item)
}

func (ctrl *Controller) DeleteItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(c, http.StatusBadRequest, "parse ID")
		return
	}
	item := model.Item{ID: uint(id)}
	if err := ctrl.Service.DeleteItem(item); err != nil {
		sendError(c, http.StatusInternalServerError, "delete item")
		return
	}
	sendSuccess(c, "delete item", nil)
}
