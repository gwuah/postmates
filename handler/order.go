package handler

// func (h *Handler) GetOrder(c *gin.Context) {
// 	id := c.Param("id")
// 	if id == "" {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"message": "Order ID not found",
// 		})
// 	}

// 	order, err := h.Repo.FindOrder(uint(id))

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message": "Failed To Retrieve Order",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Success",
// 		"order":   order,
// 	})
// 	return

// }
