package controllers

import (
	configs "blogappgolang/config"
	"blogappgolang/models"
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateCategory(c *fiber.Ctx) error {
	db := configs.DB.Db
	category := models.Category{}

	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	validate := validator.New()

	if err := validate.Struct(category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	result := db.Create(&category)

	if result.Error != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": false,
			"result": result.Error,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": true,
		"result": "Category is created",
	})

}

func AllCategory(c *fiber.Ctx) error {
	category := []models.Category{}

	result := configs.DB.Db.Find(&category)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"result": result.Error,
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"result": "Database is empty",
		})
	}

	// Convert this result to map
	categoryMap := []map[string]interface{}{}
	categoryJSON, err := json.Marshal(category)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"result": "Failed to marshal the category results",
		})
	}

	err = json.Unmarshal(categoryJSON, &categoryMap)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"result": "failed to unmarshal category results",
		})
	}

	// Now we will remove unwanted fields

	for _, cat := range categoryMap {
		delete(cat, "CreatedAt")
		delete(cat, "UpdatedAt")
		delete(cat, "DeletedAt")
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"status": true,
		"result": categoryMap,
	})

}

func SingleCategoryDetail(c *fiber.Ctx) error {
	id := c.Params("id")

	category := models.Category{}

	db := configs.DB.Db
	result := db.First(&category, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status": false,
				"result": fiber.ErrNotFound,
			})
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": false,
			"result": result.Error,
		})
	}

	categoryMap := make(map[string]interface{})
	categoryJSON, err := json.Marshal(category)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"result": "Failed to marshal category result",
		})
	}

	err = json.Unmarshal(categoryJSON, &categoryMap)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"result": "Failed to unmarshal category result",
		})
	}

	delete(categoryMap, "CreatedAt")
	delete(categoryMap, "UpdatedAt")
	delete(categoryMap, "DeletedAt")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"result": categoryMap,
	})

}

func UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	category := models.Category{}

	db := configs.DB.Db

	result := db.First(&category, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": false,
				"result": "Category not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"result": "Category failed to retrieve",
		})
	}

	var requestBody models.Category
	if err := c.BodyParser(&requestBody); err != nil {
		if err == fiber.ErrUnprocessableEntity {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": false,
				"result": "Invalid Request Body",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	if requestBody.Categoryname != "" {
		category.Categoryname = requestBody.Categoryname
	}

	result = db.Save(&category)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"result": "Failed to update category",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"result": category,
	})
}

func DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	db := configs.DB.Db

	category := models.Category{}

	if err := db.First(&category, id).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"result": err.Error(),
		})
	} else {
		db.Delete(&category)
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"status": true,
			"result": "Category Deleted Successfully",
		})
	}
}
