package data

import (
	"fmt"
	"github.com/mxschmitt/playwright-go"
	"log"
	"math"
	"strconv"
	"strings"
)

type Product struct {
	Id       int64   `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
	Sum      float64 `json:"sum"`
}

type ProductModel struct {
}

func clearStringPrice(s string) string {
	return strings.Replace(strings.Replace(strings.Trim(s, "₸"), " ", "", -1), ",", ".", -1)
}

func (p *ProductModel) GetBasket(url string, browser playwright.Browser) []Product {
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	if _, err = page.Goto(url); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	//start := time.Now()

	page.WaitForLoadState("networkidle")

	entries, err := page.QuerySelectorAll("app-ticket-header > div.row > div.text-center > p")
	if err != nil {
		log.Fatalf("could not get entries: %v", err)
	}

	for i, entry := range entries {
		title, err := entry.TextContent()
		if err != nil {
			log.Fatalf("could not get text content: %v", err)
		}
		fmt.Printf("%d: %s\n", i+1, title)
	}

	var basket []Product

	products, err := page.QuerySelectorAll("app-ticket-items > div > div > div > div.row")
	if err != nil {
		log.Fatalf("could not get products: %v", err)
	}

	for _, product := range products {
		attributes, err := product.QuerySelectorAll("div")
		if err != nil {
			log.Fatalf("could not get products: %v", err)
		}

		title, _ := attributes[1].TextContent()

		if title != "Коррекция округленияСкидка" {
			quantityStr, _ := attributes[4].TextContent()
			quantity, _ := strconv.ParseFloat(quantityStr, 64)

			_, div := math.Modf(quantity)

			unit := "ШТ"
			if strings.Contains(title, "КГ") && div != 0 {
				unit = "КГ"
			}

			priceSTR, _ := attributes[3].TextContent()
			price, _ := strconv.ParseFloat(clearStringPrice(priceSTR), 64)

			sumSTR, _ := attributes[5].TextContent()
			sum, _ := strconv.ParseFloat(clearStringPrice(sumSTR), 64)

			tmp := Product{
				Name:     title,
				Quantity: quantity,
				Unit:     unit,
				Price:    price,
				Sum:      sum,
			}

			//fmt.Println(tmp)

			basket = append(basket, tmp)
		}

		/*for _, attribute := range attributes {
			content, err := attribute.TextContent()
			if err != nil {
				log.Fatalf("could not get text content: %v", err)
			}

			fmt.Printf("%s\n", content)
		}*/
	}

	/*elapsed := time.Since(start)
	fmt.Println(elapsed)*/

	if err = page.Close(); err != nil {
		log.Fatalf("could not close page: %v", err)
	}

	return basket
}
