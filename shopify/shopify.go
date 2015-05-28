package shopify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
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
func (shopifyClient *Shopify) GetLiveProduct(shopifyID string) (Product, error) {
	urlStr := "admin/products/" + shopifyID + ".json"
	var shopifyResponse = new(productResponse)

	err := shopifyClient.makeRequest("GET", urlStr, shopifyResponse, "")
	if err != nil {
		return shopifyResponse.SingleProduct, err
	}
	fmt.Printf("[GetLiveProduct] -  Product ID: %s\n", strconv.Itoa(shopifyResponse.SingleProduct.ID))

	return shopifyResponse.SingleProduct, nil
}

// GetOrder gets order by ID
func (shopifyClient *Shopify) GetOrder(shopifyID string) (Order, error) {
	urlStr := "admin/orders/" + shopifyID + ".json"
	var shopifyResponse = new(OrderResponse)

	err := shopifyClient.makeRequest("GET", urlStr, shopifyResponse, "")
	if err != nil {
		return shopifyResponse.SingleOrder, err
	}

	fmt.Printf("[GetOrder] - Order id: %s\n", strconv.Itoa(shopifyResponse.SingleOrder.ID))

	return shopifyResponse.SingleOrder, nil
}

// CancelOrder deletes order by ID
func (shopifyClient *Shopify) CancelOrder(shopifyID string) (Order, error) {
	urlStr := "admin/orders/" + shopifyID + "/cancel.json"
	var shopifyResponse = new(OrderResponse)

	err := shopifyClient.makeRequest("POST", urlStr, shopifyResponse, "")
	if err != nil {
		return shopifyResponse.SingleOrder, err
	}

	//fmt.Printf("[CancelOrder] - Order: %v\n", shopifyResponse)

	return shopifyResponse.SingleOrder, nil
}

// PlaceOrder creates a new order
func (shopifyClient *Shopify) PlaceOrder(order OrderResponse) (Order, error) {
	urlStr := "admin/orders.json"
	var shopifyResponse = new(OrderResponse)

	orderString, _ := json.Marshal(order)

	err := shopifyClient.makeRequest("POST", urlStr, shopifyResponse, string(orderString))
	if err != nil {
		return shopifyResponse.SingleOrder, err
	}

	//fmt.Printf("[PlaceOrder] - Order: %v\n", shopifyResponse.SingleOrder)

	return shopifyResponse.SingleOrder, nil
}

// ShippingOptions returns shipping options and rates for a given shipping address
func (shopifyClient *Shopify) ShippingOptions(order Order) ([]ShippingRate, error) {

	var itemsInfo []string
	var shopifyResponse = new(ShippingRatesResponse)

	cartURLStr := "cart/"

	// Create CART
	for _, itemObj := range order.Items {
		itemsInfo = append(itemsInfo, fmt.Sprintf("%d:%d", itemObj.VariantID, itemObj.Quantity))
	}

	itemsInCartURLStr := cartURLStr + strings.Join(itemsInfo, ",")

	completeURL := fmt.Sprintf("https://%s%s%s", shopifyClient.shopifyDomain, baseURLString, itemsInCartURLStr)
	log.Printf("[ShippingOptions] - Request URL: %s", completeURL)
	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: cookieJar,
	}

	r, err := http.NewRequest("GET", completeURL, nil)
	resp, err := client.Do(r)

	defer resp.Body.Close()

	if err != nil {
		fmt.Printf("[ShippingOptions] - Error executing request : %s", err)
		return shopifyResponse.ShippingRates, err
	}

	// GET shipping Options given the cart (cookies used)
	address := order.ShippingAddress

	urlStr := cartURLStr + "shipping_rates.json?"

	v := url.Values{}
	v.Set("shipping_address[zip]", address.PostalCode)
	v.Add("shipping_address[country]", address.CountryCode)
	v.Add("shipping_address[province]", address.State)
	v.Encode()

	urlStr = urlStr + v.Encode()

	completeURL = fmt.Sprintf("https://%s%s%s", shopifyClient.shopifyDomain, baseURLString, urlStr)
	log.Printf("\n\n[ShippingOptions] - Request URL: %s", completeURL)

	r, err = http.NewRequest("GET", completeURL, nil)

	resp, err = client.Do(r)

	defer resp.Body.Close()

	if err != nil {
		fmt.Printf("[ShippingOptions] - Error executing request : %s", err)
		return shopifyResponse.ShippingRates, err
	}

	//bodyResp, _ := ioutil.ReadAll(resp.Body)
	//fmt.Printf("\n\n *****RESPONSE: %#v\n", string(bodyResp))
	err = json.NewDecoder(resp.Body).Decode(shopifyResponse)

	if err != nil {
		fmt.Printf("\n[ShippingOptions] - Decoding error: %#v", err)
		fmt.Printf("\n[ShippingOptions] - Response: %#v", resp.Body)
		return shopifyResponse.ShippingRates, err
	}

	// Address not supported error handling
	if shopifyResponse.Country != nil || shopifyResponse.Zip != nil || shopifyResponse.Province != nil {
		var errorsArray []string
		errorMessageStr := "Address Not supported: "

		if shopifyResponse.Country != nil {
			errorsArray = append(errorsArray, address.CountryCode+" "+shopifyResponse.Country[0])
		}
		if shopifyResponse.Zip != nil {
			errorsArray = append(errorsArray, address.PostalCode+" "+shopifyResponse.Zip[0])
		}
		if shopifyResponse.Province != nil {
			errorsArray = append(errorsArray, address.State+" "+shopifyResponse.Province[0])
		}
		errorMessageStr = errorMessageStr + strings.Join(errorsArray, ", ")
		addressNotSupported := errors.New(errorMessageStr)
		return shopifyResponse.ShippingRates, addressNotSupported
	}

	return shopifyResponse.ShippingRates, nil
}

func (shopifyClient *Shopify) makeRequest(method string, urlStr string, body interface{}, payload string) error {
	url := fmt.Sprintf("https://%s%s%s", shopifyClient.shopifyDomain, baseURLString, urlStr)
	log.Printf("\n\n[makeRequest] - Request URL: %s", url)
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
		fmt.Printf("[makeRequest] - 404 on executing request: %s\n", url)
	} else if resp.StatusCode == 429 {
		fmt.Printf("[makeRequest] - Rate limited!\n")
		rateLimitErr := errors.New("API rate limit exceeded")
		return rateLimitErr
	}
	if err != nil {
		fmt.Printf("[makeRequest] - Error executing request : %s", err)
		return err
	}

	//bodyResp, _ := ioutil.ReadAll(resp.Body)
	//fmt.Printf("\n\n *****RESPONSE: %#v\n", string(bodyResp))

	err = json.NewDecoder(resp.Body).Decode(body)

	if err != nil {
		fmt.Printf("\n[makeRequest] - Decoding error: %#v", err)
		fmt.Printf("\n[makeRequest] - Response: %#v", resp.Body)
		return err
	}

	return nil
}
