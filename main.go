package main

import (
	"net/http"

	"github.com/Rhymond/go-money"
	"github.com/gin-gonic/gin"
	"github.com/raghavendrah25/golang-backend/internal/category"
	"github.com/raghavendrah25/golang-backend/internal/product"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/products", func(c *gin.Context) {
		products := []product.Product{
			{
				ID:               "1",
				Name:             "Apple",
				Description:      "Apple iPhone 12",
				PriceVATExcluded: money.New(1000, "USD"),
				VAT:              money.New(200, "USD"),
			},
			{
				ID:               "2",
				Name:             "Samsung",
				Description:      "Samsung Galaxy S21",
				PriceVATExcluded: money.New(900, "USD"),
				VAT:              money.New(180, "USD"),
			},
			{
				ID:               "3",
				Name:             "OnePlus",
				Description:      "OnePlus 9 Pro",
				PriceVATExcluded: money.New(800, "USD"),
				VAT:              money.New(160, "USD"),
			},
		}
		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusOK, products)
	})

	r.GET("/categories", func(c *gin.Context) {
		categories := []category.Category{
			{
				ID:          "1",
				Name:        "Electronics",
				Description: "Electronic Items",
			},
			{
				ID:          "2",
				Name:        "Clothing",
				Description: "Clothing Items",
			},
			{
				ID:          "3",
				Name:        "Books",
				Description: "Books Items",
			},
		}
		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusOK, categories)
	})
	r.Run()
}
