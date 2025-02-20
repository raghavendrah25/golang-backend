package server

import (
	"net/http"

	"github.com/Rhymond/go-money"
	"github.com/gin-gonic/gin"
	"github.com/raghavendrah25/golang-backend/internal/category"
	"github.com/raghavendrah25/golang-backend/internal/product"
)

type Server struct {
	Engine *gin.Engine
	Config Config
}

type Config struct {
	Port string
}

func NewServer(config Config) (*Server, error) {
	engine := gin.Default()
	s := &Server{
		Engine: engine,
		Config: config,
	}
	engine.Use(s.CORSMiddleware, s.CheckRequest)
	engine.GET("/products", s.Products)
	engine.GET("/categories", s.Categories)

	return s, nil
}

func (s *Server) Run() error {
	return s.Engine.Run(":" + s.Config.Port)
}

func (s *Server) CORSMiddleware(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "http://localhost:3001")
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
	twoUSD := money.New(200, "USD")
	fourUSD := money.New(400, "USD")
	products := []product.Product{
		{
			ID:               "1",
			Name:             "Test",
			Image:            "https://www.practical-go-lessons.com/img/practical-go-lessons-book10.a8a05387.jpg",
			ShortDescription: "New",
			Description:      "This is my product",
			PriceVATExcluded: product.Amount{
				Money:   twoUSD,
				Display: twoUSD.Display(),
			},
			VAT: product.Amount{
				Money:   twoUSD,
				Display: twoUSD.Display(),
			},
			TotalPrice: product.Amount{
				Money:   fourUSD,
				Display: fourUSD.Display(),
			},
		},
		{
			ID:               "2",
			Name:             "Test",
			Image:            "https://www.practical-go-lessons.com/img/practical-go-lessons-book10.a8a05387.jpg",
			ShortDescription: "New",
			Description:      "This is my product",
			PriceVATExcluded: product.Amount{
				Money:   twoUSD,
				Display: twoUSD.Display(),
			},
			VAT: product.Amount{
				Money:   twoUSD,
				Display: twoUSD.Display(),
			},
			TotalPrice: product.Amount{
				Money:   fourUSD,
				Display: fourUSD.Display(),
			},
		},
		{
			ID:               "3",
			Name:             "Test",
			Image:            "https://www.practical-go-lessons.com/img/practical-go-lessons-book10.a8a05387.jpg",
			ShortDescription: "New",
			Description:      "This is my product",
			PriceVATExcluded: product.Amount{
				Money:   twoUSD,
				Display: twoUSD.Display(),
			},
			VAT: product.Amount{
				Money:   twoUSD,
				Display: twoUSD.Display(),
			},
			TotalPrice: product.Amount{
				Money:   fourUSD,
				Display: fourUSD.Display(),
			},
		},
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