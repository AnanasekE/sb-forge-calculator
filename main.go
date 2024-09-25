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

type Recipe struct {
	ItemID          string         `json:"itemId"`
	TimeHours       int            `json:"timeHours"`
	HotmRequirement int            `json:"hotmRequirement"`
	Items           map[string]int `json:"items"`
}

// hypixel api url: https://api.hypixel.net/v2/skyblock/bazaar

func main() {
	productsFile, err := os.Open("products.json")
	if err != nil {
		return
	}
	productsJson, _ := io.ReadAll(productsFile)

	var marketData MarketData
	err = json.Unmarshal(productsJson, &marketData)
	if err != nil {
		return
	}

	forgeRecipesFile, err := os.Open("forge_recipes.json")
	if err != nil {
		return
	}
	forgeRecipesJson, _ := io.ReadAll(forgeRecipesFile)

	var recipes []Recipe
	err = json.Unmarshal(forgeRecipesJson, &recipes)
	if err != nil {
		return
	}

	fmt.Println(marketData.Products)
}
