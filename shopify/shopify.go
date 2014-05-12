package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	baseUrlString = ".myshopify.com/"
)

type Shopify struct {
	shopifyDomain      string
	shopifySecretToken string

	Products []ShopifyProduct
}

func NewClient(domain string, secrettoken string) Shopify {
	shop := Shopify{shopifyDomain: domain, shopifySecretToken: secrettoken}
	return shop
}

func (shopifyClient *Shopify) LoadProducts() {
	//load first block of 250
	urlStr := "admin/products.json?limit=250&page=1"
	var page int = 1

	var shopifyResponse = new(ShopifyResponse)
	shopifyClient.MakeRequest("GET", urlStr, shopifyResponse)

	shopifyClient.Products = shopifyResponse.Products[:]

	//if response was == 250 products
	if len(shopifyResponse.Products) == 250 {
		//load every successive block
		var done bool = false
		//var lastId int = shopifyResponse.Products[249].Id

		for done == false {
			log.Printf("Shopify: Loaded page %v, that's %v products!\n", page, len(shopifyClient.Products))
			page++
			//get this thread to wait 0.5 seconds
			time.Sleep(time.Second / 2)
			urlStr = fmt.Sprintf("admin/products.json?limit=250&page=%v", page)
			shopifyResponse = new(ShopifyResponse)
			shopifyClient.MakeRequest("GET", urlStr, shopifyResponse)
			if len(shopifyResponse.Products) > 0 {
				shopifyClient.Products = append(shopifyClient.Products, shopifyResponse.Products[:]...)
				//lastId = shopifyResponse.Products[len(shopifyResponse.Products)-1].Id
			} else {
				done = true
			}

		}

	}
	return

}

func (shopifyClient *Shopify) GetLiveProduct(shopifyId string) ShopifyProduct {
	urlStr := "admin/products/" + shopifyId + ".json"
	var shopifyResponse = new(ShopifyResponse)

	shopifyClient.MakeRequest("GET", urlStr, shopifyResponse)

	fmt.Printf("%v\n", shopifyResponse.Product.Id)

	return shopifyResponse.Product
}

func (shopifyClient *Shopify) MakeRequest(method string, urlStr string, body interface{}) {
	url := fmt.Sprintf("https://%s%s%s", shopifyClient.shopifyDomain, baseUrlString, urlStr)
	client := &http.Client{}
	buf := new(bytes.Buffer)
	r, err := http.NewRequest(method, url, buf)
	r.Header.Add("X-Shopify-Access-Token", shopifyClient.shopifySecretToken)
	resp, err := client.Do(r)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("404 on executing request: %s", url)
	} else if resp.StatusCode == 429 {
		fmt.Printf("Rate limited!")
	}
	if err != nil {
		fmt.Printf("Error executing request : %s", err)
	}

	err = json.NewDecoder(resp.Body).Decode(body)
	if err != nil {
		fmt.Print(err)
		fmt.Print(resp.Body)
	}
}
