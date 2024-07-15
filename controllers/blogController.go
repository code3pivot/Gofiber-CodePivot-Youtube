package controllers

import (
	configs "blogappgolang/config"
	"blogappgolang/models"
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateBlog(c *fiber.Ctx) error {
	db := configs.DB.Db
	blog := models.Blog{}

	if err := c.BodyParser(&blog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	validate := validator.New()
	err := validate.StructExcept(blog, "User", "Category")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	// User Association
	user := models.User{}
	result := db.Find(&user, blog.UserID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": false,
				"result": fiber.ErrNotFound,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"result": result.Error,
		})
	}

	blog.User = user

	// Category Association
	category := models.Category{}
	result = db.Find(&category, blog.CategoryID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": false,
				"result": fiber.ErrNotFound,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"result": result.Error,
		})
	}

	blog.Category = category

	result = db.Create(&blog)

	if result.Error != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": false,
			"result": result.Error,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": true,
		"result": "Blog is created",
	})
}

func GetAllBlogs(c *fiber.Ctx) error {
	db := configs.DB.Db
	blogs := []models.Blog{}

	err := db.Preload("Category").Preload("User").Find(&blogs).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"result": err.Error(),
		})
	}
	blogMaps := make([]map[string]interface{}, len(blogs))

	for i, blogs := range blogs {
		blogBytes, err := json.Marshal(blogs)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to marshal blog data",
				"result":  err.Error(),
			})
		}
		var blogMap map[string]interface{}
		err = json.Unmarshal(blogBytes, &blogMap)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to unmarshal blog data",
				"result":  err.Error(),
			})
		}

		// To remove unwanted fields from the JSON results
		delete(blogMap, "CreatedAt")
		delete(blogMap, "UpdatedAt")
		delete(blogMap, "DeletedAt")

		if category, ok := blogMap["Category"].(map[string]interface{}); ok {
			delete(category, "CreatedAt")
			delete(category, "UpdatedAt")
			delete(category, "DeletedAt")

			blogMap["Category"] = category
		}

		if user, ok := blogMap["User"].(map[string]interface{}); ok {
			delete(user, "CreatedAt")
			delete(user, "UpdatedAt")
			delete(user, "DeletedAt")
			delete(user, "password")

			blogMap["User"] = user
		}

		blogMaps[i] = blogMap

	}

	if len(blogMaps) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"result": "Database is empty",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"result": blogMaps,
	})

}

func GetSingleBlogByID(c *fiber.Ctx) error {
	id := c.Params("id")

	db := configs.DB.Db

	blog := models.Blog{}

	result := db.Preload("Category").Preload("User").First(&blog, id)

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

	blogBytes, err := json.Marshal(blog)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to marshal blog data",
			"result":  err.Error(),
		})
	}

	var blogMap map[string]interface{}
	err = json.Unmarshal(blogBytes, &blogMap)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to unmarshal blog data",
			"result":  err.Error(),
		})
	}

	delete(blogMap, "CreatedAt")
	delete(blogMap, "UpdatedAt")
	delete(blogMap, "DeletedAt")
	delete(blogMap, "userID")
	delete(blogMap, "categoryID")

	if category, ok := blogMap["Category"].(map[string]interface{}); ok {
		delete(category, "CreatedAt")
		delete(category, "UpdatedAt")
		delete(category, "DeletedAt")

		blogMap["Category"] = category
	}

	if user, ok := blogMap["User"].(map[string]interface{}); ok {
		delete(user, "CreatedAt")
		delete(user, "UpdatedAt")
		delete(user, "DeletedAt")
		delete(user, "password")

		blogMap["User"] = user
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"result": blogMap,
	})

}

// Updating the blog we require token from the author in order to update blog
func UpdateBlog(c *fiber.Ctx) error {
	id := c.Params("id")
	userId := c.Locals("user").(models.User)

	blog := models.Blog{}

	db := configs.DB.Db

	result := db.Where("id=? AND user_id =?", id, userId.ID).First(&blog)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": false,
				"result": "blog not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"result": "Failed to retrieve blog",
		})
	}

	var requestBody models.Blog

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

	if requestBody.Blogtitle != "" {
		blog.Blogtitle = requestBody.Blogtitle
	}

	if requestBody.Blogsubtitle != "" {
		blog.Blogsubtitle = requestBody.Blogsubtitle
	}

	if requestBody.Blogimage != "" {
		blog.Blogimage = requestBody.Blogimage
	}

	if requestBody.Blogdescription != "" {
		blog.Blogdescription = requestBody.Blogdescription
	}

	result = db.Save(&blog)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"result": "Failed to update blog",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Blog Updated Successfully",
		"result":  blog,
	})

}

func DeleteBlog(c *fiber.Ctx) error {
	id := c.Params("id")
	userId := c.Locals("user").(models.User)

	db := configs.DB.Db
	blog := models.Blog{}

	if err := db.Where("id=? AND user_id=?", id, userId.ID).First(&blog).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"result": err.Error(),
		})
	} else {
		db.Delete(&blog)
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"status": true,
			"result": "Blog Deleted Successful",
		})
	}

}
