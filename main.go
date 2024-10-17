package main

import (
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/text/language"
	_ "golang.org/x/text/language"
	"golang.org/x/text/message"
	_ "golang.org/x/text/message"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// MarketData Root struct
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
	ProfitPerHour   int
	ProfitTotal     int
	Cost            int
}

// hypixel api url: https://api.hypixel.net/v2/skyblock/bazaar

func main() {
	// time.Now().Unix() - UNIX TIMESTAMP
	marketData := loadMarketData()
	recipes := loadRecipes()

	if time.Now().Unix()-(marketData.LastUpdated/1000) > 600 {
		downloadBazaarPrices()
		marketData = loadMarketData()
	}

	var found bool
	var newRecipes []Recipe
	for _, recipe := range recipes {
		found = false
		for _, product := range marketData.Products {
			if recipe.ItemID == product.ProductID {
				found = true

				cost := 0
				for itemName, itemAmount := range recipe.Items {
					for _, prod := range marketData.Products {
						if itemName == prod.ProductID {
							cost += int(prod.QuickStatus.BuyPrice) * itemAmount
						}
					}
				}
				recipe.Cost = cost

				recipe.ProfitPerHour = (int(product.QuickStatus.BuyPrice) / recipe.TimeHours) - cost/recipe.TimeHours
				recipe.ProfitTotal = int(product.QuickStatus.BuyPrice) - cost
				newRecipes = append(
					newRecipes,
					recipe,
				)
				break
			}
		}
		if found == false {
			log.Printf("Item Not Found %s", recipe.ItemID)
		}
	}

	fmt.Println("Enter your HOTM Level: ")
	var levelStr string
	_, err := fmt.Scan(&levelStr)
	if err != nil {
		log.Fatalf("Error while scanning for HOTM Level: %s", err)
	}
	hotmLevel, err := strconv.Atoi(levelStr)
	if err != nil {
		log.Fatalf("Error while parsing HOTM Level: %s", err)
	}

	var slotsStr string
	fmt.Println("Enter how many forge slots you have: ")
	_, err = fmt.Scan(&slotsStr)
	if err != nil {
		log.Fatalf("Error while scanning for forge slots: %s", err)
	}
	slots, err := strconv.Atoi(slotsStr)
	if err != nil {
		log.Fatalf("Error while parsing HOTM Level: %s", err)
	}

	p := message.NewPrinter(language.English)
	writer := table.NewWriter()
	writer.AppendHeader(table.Row{"ItemID",
		"Cost",
		"PROFIT PER HOUR",
		"Profit",
		"Time",
		fmt.Sprintf("Profit for %s slots per hour", fmt.Sprint(slots)),
		"Total Cost",
		"Total Profit",
		"HOTM Req",
		"sorting"})
	for _, recipe := range newRecipes {
		if recipe.HotmRequirement > hotmLevel {
			continue
		}
		row := table.Row{recipe.ItemID,
			p.Sprint(recipe.Cost),
			p.Sprint(recipe.ProfitPerHour),
			p.Sprint(recipe.ProfitTotal),
			recipe.TimeHours,
			p.Sprint(recipe.ProfitPerHour * slots),
			p.Sprint(recipe.Cost * slots),
			p.Sprint(recipe.ProfitTotal * slots),
			recipe.HotmRequirement,
			recipe.ProfitPerHour,
		}

		writer.AppendRow(row)
	}
	writer.SetStyle(table.StyleDefault)
	writer.SortBy([]table.SortBy{{Number: 10, Mode: table.DscNumericAlpha}})
	fmt.Println(writer.Render())

	fmt.Println("Press ENTER to end program")
	_, _ = fmt.Scan()
}

func loadRecipes() []Recipe {
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
	return recipes
}

func loadMarketData() MarketData {

	productsFile, err := os.Open("products.json")
	if err != nil {
		//log.Fatalf("Error loading json: %s", err.Error())
		log.Println("Error loading json, downloading a new one")
		downloadBazaarPrices()
	}
	productsJson, _ := io.ReadAll(productsFile)

	var marketData MarketData
	err = json.Unmarshal(productsJson, &marketData)
	if err != nil {
		log.Fatalf("Error parsing json: %s", err.Error())
	}
	return marketData
}

func downloadBazaarPrices() {
	response, err := http.Get("https://api.hypixel.net/v2/skyblock/bazaar")
	if err != nil {
		log.Fatalf("Error while fetching bz prices: %s", err)
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error while reading bz prices: %s", err)
	}

	err = os.WriteFile("products.json", data, os.ModePerm)
	if err != nil {
		log.Fatalf("Error while saving bz prices file: %s", err)
	}

	log.Println("Bazaar prices downloaded and saved successfully!")
}
