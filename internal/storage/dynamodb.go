package storage

import (
	"context"
	"fmt"

	"github.com/Rhymond/go-money"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
	"github.com/raghavendrah25/golang-backend/internal/product"
	"strconv"
)

type DynamoDB struct {
	tableName string
	client    *dynamodb.Client
}

func NewDynamoDB(tableName string, c *gin.Context) (*DynamoDB, error) {
	cfg, err := config.LoadDefaultConfig(c.Request.Context(), config.WithRegion("us-east-2"), config.WithSharedConfigProfile("default"))
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration, %v", err)
	}
	client := dynamodb.NewFromConfig(cfg)

	return &DynamoDB{
		tableName: tableName,
		client:    client,
	}, nil
}

func (d *DynamoDB) CreateProduct(product product.Product, ctx context.Context) error {
	// Marshal product amounts (money)
	priceVatExcludedCents := product.PriceVATExcluded.Money.Amount()
	vatCents := product.VAT.Money.Amount()
	totalPriceCents := product.TotalPrice.Money.Amount()

	// Prepare DynamoDB item to be inserted
	item := make(map[string]types.AttributeValue)
	item["PK"] = &types.AttributeValueMemberS{Value: product.ID}   // Example of using product ID as PK
	item["SK"] = &types.AttributeValueMemberS{Value: product.Name} // Example of using product name as SK
	item["Name"] = &types.AttributeValueMemberS{Value: product.Name}
	item["Image"] = &types.AttributeValueMemberS{Value: product.Image}
	item["ShortDescription"] = &types.AttributeValueMemberS{Value: product.ShortDescription}
	item["Description"] = &types.AttributeValueMemberS{Value: product.Description}
	item["PriceVATExcluded"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", priceVatExcludedCents)} // Storing as integer (cents)
	item["VAT"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", vatCents)}                           // Storing VAT as integer (cents)
	item["TotalPrice"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", totalPriceCents)}             // Storing total price as integer (cents)

	// Perform DynamoDB PutItem operation
	_, err := d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &d.tableName,
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB, %v", err)
	}
	return nil
}

func (d *DynamoDB) GetProducts(ctx context.Context) ([]product.Product, error) {
	// Perform DynamoDB Scan operation
	res, err := d.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: &d.tableName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan DynamoDB, %v", err)
	}

	// Unmarshal DynamoDB items to product.Product
	var products []product.Product
	for _, item := range res.Items {
		// Create a new product object (renamed to productItem to avoid conflict)
		productItem := product.Product{
			ID:               item["PK"].(*types.AttributeValueMemberS).Value,
			Name:             item["Name"].(*types.AttributeValueMemberS).Value,
			Image:            item["Image"].(*types.AttributeValueMemberS).Value,
			ShortDescription: item["ShortDescription"].(*types.AttributeValueMemberS).Value,
			Description:      item["Description"].(*types.AttributeValueMemberS).Value,
		}

		// Handle PriceVATExcluded, VAT, and TotalPrice as Amount
		if priceVATExcluded, ok := item["PriceVATExcluded"].(*types.AttributeValueMemberN); ok {
			val, err := strconv.ParseInt(priceVATExcluded.Value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse PriceVATExcluded: %v", err)
			}
			productItem.PriceVATExcluded = product.Amount{
				Money:   money.New(val, "USD"),
				Display: fmt.Sprintf("$%d", val),
			}
		}

		if vat, ok := item["VAT"].(*types.AttributeValueMemberN); ok {
			val, err := strconv.ParseInt(vat.Value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse VAT: %v", err)
			}
			productItem.VAT = product.Amount{
				Money:   money.New(val, "USD"),
				Display: fmt.Sprintf("$%d", val),
			}
		}

		if totalPrice, ok := item["TotalPrice"].(*types.AttributeValueMemberN); ok {
			val, err := strconv.ParseInt(totalPrice.Value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse TotalPrice: %v", err)
			}
			productItem.TotalPrice = product.Amount{
				Money:   money.New(val, "USD"),
				Display: fmt.Sprintf("$%d", val),
			}
		}

		// Append the productItem to the products slice
		products = append(products, productItem)
	}

	return products, nil
}
