package main

import (
	"fmt"

	"github.com/Rhymond/go-money"
)

type Item struct {
	ID   string
	Name *money.Money
}

type Cart struct {
	ID       string
	IsLocked bool
	Item     []Item
}

func (c *Cart) TotalPrice() (*money.Money, error) {
	totalPrice := money.New(0, "USD")

	for _, item := range c.Item {
		var err error
		totalPrice, err = totalPrice.Add(item.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to add item price: %w", err)
		}
	}
	return totalPrice, nil
}

func main() {


	hello()
}