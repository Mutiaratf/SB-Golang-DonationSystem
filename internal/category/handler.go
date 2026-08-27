package category

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type categoryHandler struct {
	categoryService *CategoryService
}

func NewHandler(categoryService *CategoryService) *categoryHandler {
	return &categoryHandler{
		categoryService: categoryService,
	}
}

// @Summary Create category
// @Tags Categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CategoryRequest true "Category"
// @Success 201 {object} map[string]interface{}
// @Failure 400,401,500 {object} map[string]string
// @Router /categories [post]
func (h *categoryHandler) CreateCategory(c *gin.Context) {
	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request body",
		})
		return
	}

	category, err := h.categoryService.CreateCategory(&req)
	if err != nil {

		if errors.Is(err, ErrCategoryNameEmpty) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Category name is empty",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create category",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Category created successfully",
		"data":    category,
	})
}

// @Summary List categories
// @Tags Categories
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /categories [get]
func (h *categoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve categories",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Categories retrieved successfully",
		"data":    categories,
	})
}

// @Summary Update category
// @Tags Categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param request body CategoryRequest true "Category"
// @Success 200 {object} map[string]interface{}
// @Failure 400,401,404,500 {object} map[string]string
// @Router /categories/{id} [put]
func (h *categoryHandler) UpdateCategory(c *gin.Context) {
	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request body",
		})
		return
	}

	id := c.Param("id")
	categoryID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid category ID",
		})
		return
	}

	category, err := h.categoryService.UpdateCategory(categoryID, &req)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Category not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Category updated successfully",
		"data":    category,
	})
}

// @Summary Delete category
// @Tags Categories
// @Security BearerAuth
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]string
// @Failure 400,401,404,500 {object} map[string]string
// @Router /categories/{id} [delete]
func (h *categoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid category ID",
		})
		return
	}

	err = h.categoryService.DeleteCategory(int64(id))
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Category not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Category deleted successfully",
	})
}
