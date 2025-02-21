package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raghavendrah25/golang-backend/internal/category"
	"github.com/raghavendrah25/golang-backend/internal/product"
	"github.com/raghavendrah25/golang-backend/internal/storage"
)

type Server struct {
	Engine  *gin.Engine
	Config  Config
	storage storage.Storage
}

type Config struct {
	Port string
	storage storage.Storage
}

func NewServer(config Config) (*Server, error) {
	engine := gin.Default()
	s := &Server{
		Engine: engine,
		Config: config,
		storage: config.storage,
	}
	engine.Use(s.CORSMiddleware)
	engine.POST("/products", s.Products)
	engine.GET("/getProducts", s.GetProducts)
	engine.GET("/categories", s.Categories)
	engine.GET("/ping", s.PingTest)

	return s, nil
}

func (s *Server) Run() error {
	return s.Engine.Run(":" + s.Config.Port)
}

func (s *Server) CORSMiddleware(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
}

func (s *Server) PingTest(c *gin.Context) {
	(*c).JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (s *Server) Categories(c *gin.Context) {
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
	c.JSON(http.StatusOK, categories)
}

func (s *Server) Products(c *gin.Context) {
	dynamoDB, err := storage.NewDynamoDB("ecommerce-dev", c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// New product from request body (should be dynamically populated, not hardcoded)
	var newProduct product.Product
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product data"})
		return
	}

	// Save the product to DynamoDB
	err = dynamoDB.CreateProduct(newProduct, c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product created successfully", "product": newProduct})
}

func (s *Server) GetProducts(c *gin.Context) {
	dynamoDB, err := storage.NewDynamoDB("ecommerce-dev", c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get products from DynamoDB
	products, err := dynamoDB.GetProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (s *Server) CheckRequest(c *gin.Context) {
	authVault := c.GetHeader("Authorization")
	if authVault != "1234567890" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}
