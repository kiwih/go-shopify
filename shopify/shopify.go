package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	baseURLString = ".myshopify.com/"
)

// Shopify models the shopify client parameters
type Shopify struct {
	shopifyDomain      string
	shopifySecretToken string

	Products []Product
}

// NewClient inits shopify client
func NewClient(domain string, secrettoken string) Shopify {
	shop := Shopify{shopifyDomain: domain, shopifySecretToken: secrettoken}
	return shop
}

// LoadProducts returns first block of 250
func (shopifyClient *Shopify) LoadProducts() {

	urlStr := "admin/products.json?limit=250&page=1"
	var page = 1

	var shopifyResponse = new(productResponse)
	shopifyClient.makeRequest("GET", urlStr, shopifyResponse, "")

	shopifyClient.Products = shopifyResponse.Products[:]

	//if response was == 250 products
	if len(shopifyResponse.Products) == 250 {
		//load every successive block
		var done = false
		//var lastId int = shopifyResponse.Products[249].Id

		for done == false {
			log.Printf("[LoadProducts] - Shopify: Loaded page %v, that's %v products!\n", page, len(shopifyClient.Products))
			page++
			//get this thread to wait 0.5 seconds
			time.Sleep(time.Second / 2)
			urlStr = fmt.Sprintf("admin/products.json?limit=250&page=%v", page)
			shopifyResponse = new(productResponse)
			shopifyClient.makeRequest("GET", urlStr, shopifyResponse, "")
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

// GetLiveProduct gets product by ID
func (shopifyClient *Shopify) GetLiveProduct(shopifyID string) Product {
	urlStr := "admin/products/" + shopifyID + ".json"
	var shopifyResponse = new(productResponse)

	shopifyClient.makeRequest("GET", urlStr, shopifyResponse, "")

	fmt.Printf("[GetLiveProduct] -  Product ID: %s\n", strconv.Itoa(shopifyResponse.SingleProduct.ID))

	return shopifyResponse.SingleProduct
}

// GetOrder gets order by ID
func (shopifyClient *Shopify) GetOrder(shopifyID string) Order {
	urlStr := "admin/orders/" + shopifyID + ".json"
	var shopifyResponse = new(OrderResponse)

	shopifyClient.makeRequest("GET", urlStr, shopifyResponse, "")

	fmt.Printf("[GetOrder] - Order id: %s\n", strconv.Itoa(shopifyResponse.SingleOrder.ID))

	return shopifyResponse.SingleOrder
}

// CancelOrder deletes order by ID
func (shopifyClient *Shopify) CancelOrder(shopifyID string) Order {
	urlStr := "admin/orders/" + shopifyID + "/cancel.json"
	var shopifyResponse = new(OrderResponse)

	shopifyClient.makeRequest("POST", urlStr, shopifyResponse, "")

	//fmt.Printf("[CancelOrder] - Order: %v\n", shopifyResponse)

	return shopifyResponse.SingleOrder
}

// PlaceOrder creates a new order
func (shopifyClient *Shopify) PlaceOrder(order OrderResponse) Order {
	urlStr := "admin/orders.json"
	var shopifyResponse = new(OrderResponse)

	orderString, _ := json.Marshal(order)

	shopifyClient.makeRequest("POST", urlStr, shopifyResponse, string(orderString))

	fmt.Printf("[PlaceOrder] - Order: %v\n", shopifyResponse.SingleOrder)

	return shopifyResponse.SingleOrder
}

func (shopifyClient *Shopify) makeRequest(method string, urlStr string, body interface{}, payload string) {
	url := fmt.Sprintf("https://%s%s%s", shopifyClient.shopifyDomain, baseURLString, urlStr)
	log.Printf("[makeRequest] - Request URL: %s", url)
	client := &http.Client{}
	buf := new(bytes.Buffer)

	if payload != "" {
		//fmt.Printf("\n\n\nPAYLOAD string: %#v", payload)
		buf = bytes.NewBuffer([]byte(payload))
	}
	r, err := http.NewRequest(method, url, buf)

	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	r.Header.Add("X-Shopify-Access-Token", shopifyClient.shopifySecretToken)

	resp, err := client.Do(r)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("[makeRequest] - 404 on executing request: %s", url)
	} else if resp.StatusCode == 429 {
		fmt.Printf("[makeRequest] - Rate limited!")
	}
	if err != nil {
		fmt.Printf("[makeRequest] - Error executing request : %s", err)
	}

	// bodyResp, _ := ioutil.ReadAll(resp.Body)
	// fmt.Printf("\n\nRESPONSE BODY: %#v", string(bodyResp))

	err = json.NewDecoder(resp.Body).Decode(body)

	if err != nil {
		fmt.Printf("\n[makeRequest] - Decoding error: %#v", err)
		fmt.Printf("\n[makeRequest] - Response: %#v", resp.Body)
	}
}
