package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// Root struct
type MarketData struct {
	Success     bool               `json:"success"`
	LastUpdated int64              `json:"lastUpdated"`
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
		log.Fatalf("Error loading json: %s", err.Error())
	}
	productsJson, _ := io.ReadAll(productsFile)

	var marketData MarketData
	err = json.Unmarshal(productsJson, &marketData)
	if err != nil {
		log.Fatalf("Error parsing json: %s", err.Error())
	}

	forgeRecipesFile, err := os.Open("forge_recipes.json")
	if err != nil {
		log.Fatalf("Error loading json: %s", err.Error())
	}
	forgeRecipesJson, _ := io.ReadAll(forgeRecipesFile)

	var recipes []Recipe
	err = json.Unmarshal(forgeRecipesJson, &recipes)
	if err != nil {
		log.Fatalf("Error parsing json: %s", err.Error())
	}

	var found bool
	for _, recipe := range recipes {
		found = false
		for _, product := range marketData.Products {
			if recipe.ItemID == product.ProductID {
				//log.Printf("Item Found %s", recipe.ItemID)
				found = true
				break
			}
		}
		if found == false {
			log.Printf("Item Not Found %s", recipe.ItemID)
		}
	}
}
