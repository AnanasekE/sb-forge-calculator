package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Root struct
type MarketData struct {
	Success     bool               `json:"success"`
	LastUpdated time.Time          `json:"lastUpdated"`
	Products    map[string]Product `json:"products"`
}

// Product struct
type Product struct {
	ProductID   string         `json:"product_id"`
	SellSummary []TradeSummary `json:"sell_summary"`
	BuySummary  []TradeSummary `json:"buy_summary"`
	QuickStatus QuickStatus    `json:"quick_status"`
}

// TradeSummary struct
type TradeSummary struct {
	Amount       int     `json:"amount"`
	PricePerUnit float64 `json:"pricePerUnit"`
	Orders       int     `json:"orders"`
}

// QuickStatus struct
type QuickStatus struct {
	ProductID      string  `json:"productId"`
	SellPrice      float64 `json:"sellPrice"`
	SellVolume     int     `json:"sellVolume"`
	SellMovingWeek int     `json:"sellMovingWeek"`
	SellOrders     int     `json:"sellOrders"`
	BuyPrice       float64 `json:"buyPrice"`
	BuyVolume      int     `json:"buyVolume"`
	BuyMovingWeek  int     `json:"buyMovingWeek"`
	BuyOrders      int     `json:"buyOrders"`
}

func main() {
	productsFile, err := os.Open("products.json")
	if err != nil {
		fmt.Println("BORKED LOADING JSON")
	}
	productsJson, _ := io.ReadAll(productsFile)

	var marketData MarketData
	json.Unmarshal(productsJson, &marketData)

	forgeRecipesFile, err := os.Open("forge_recipes.json")
	if err != nil {
		fmt.Println("BORKED LOADING FORGE RECIPES")
	}
	forgeRecipesJson, _ := io.ReadAll(forgeRecipesFile)
	// TODO add recipe struct
}
