/*
	Copyright 2015 Arduino LLC (http://www.arduino.cc/)

	This file is part of go-shopify.

	go-shopify is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, version 3 of the License,
	go-shopify is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with go-shopify.  If not, see <http://www.gnu.org/licenses/>.
*/

package shopify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	jww "github.com/spf13/jwalterweatherman"
)

const (
	baseURLString = ".myshopify.com/"
)

// Shopify models the shopify client parameters
type Shopify struct {
	shopifyDomain      string
	shopifySecretToken string
	// we need the public Store URL because shipping_rates.json doesn't work with Redirect
	shopifyPublicURL string

	Products []Product
}

// NewClient inits shopify client
func NewClient(domain string, secrettoken string, publicURL string) Shopify {
	shop := Shopify{shopifyDomain: domain, shopifySecretToken: secrettoken, shopifyPublicURL: publicURL}
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
			jww.INFO.Printf("[LoadProducts] - Shopify: Loaded page %v, that's %v products!\n", page, len(shopifyClient.Products))
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
	jww.INFO.Printf("[GetLiveProduct] -  Product ID: %s\n", strconv.Itoa(shopifyResponse.SingleProduct.ID))

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

	jww.INFO.Printf("[GetOrder] - Order id: %s\n", strconv.Itoa(shopifyResponse.SingleOrder.ID))

	return shopifyResponse.SingleOrder, nil
}

// GetOrderByName gets order by Order Name
func (shopifyClient *Shopify) GetOrderByName(shopifyOrderName string) (Order, error) {
	// adding "&status=any" allows to retrieve canceled orders also
	urlStr := "admin/orders.json?name=" + shopifyOrderName + "&status=any"
	var shopifyResponse = new(OrderResponse)

	err := shopifyClient.makeRequest("GET", urlStr, shopifyResponse, "")
	if err != nil {
		return shopifyResponse.SingleOrder, err
	}

	if len(shopifyResponse.Orders) > 0 {
		jww.INFO.Printf("[GetOrderByName] - Order id: %s\n", strconv.Itoa(shopifyResponse.Orders[0].ID))

		return shopifyResponse.Orders[0], nil
	}

	jww.INFO.Printf("[GetOrderByName] - No active order with name '%s'", shopifyOrderName)
	return shopifyResponse.SingleOrder, err

}

// CancelOrder deletes order by ID
func (shopifyClient *Shopify) CancelOrder(shopifyID string) (Order, error) {
	urlStr := "admin/orders/" + shopifyID + "/cancel.json"
	var shopifyResponse = new(OrderResponse)

	err := shopifyClient.makeRequest("POST", urlStr, shopifyResponse, "")
	if err != nil {
		return shopifyResponse.SingleOrder, err
	}

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
	jww.INFO.Printf("[ShippingOptions] - Request URL: %s", completeURL)
	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: cookieJar,
	}

	r, err := http.NewRequest("GET", completeURL, nil)
	resp, err := client.Do(r)

	defer resp.Body.Close()

	if err != nil {
		jww.ERROR.Printf("[ShippingOptions] - Error executing request : %s", err)
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

	completeURL = fmt.Sprintf("%s%s", shopifyClient.shopifyPublicURL, urlStr)
	jww.INFO.Printf("[ShippingOptions] - Request URL: %s", completeURL)

	r, err = http.NewRequest("GET", completeURL, nil)
	resp, err = client.Do(r)

	defer resp.Body.Close()

	if err != nil {
		jww.ERROR.Printf("[ShippingOptions] - Error executing request : %s", err)
		return shopifyResponse.ShippingRates, err
	}

	err = json.NewDecoder(resp.Body).Decode(shopifyResponse)

	if err != nil {
		jww.ERROR.Printf("[ShippingOptions] - Decoding error: %#v", err)
		jww.ERROR.Printf("[ShippingOptions] - Response: %#v", resp.Body)
		return shopifyResponse.ShippingRates, err
	}

	if shopifyResponse.Error != nil {
		genericError := errors.New(strings.Join(shopifyResponse.Error, ", "))
		return shopifyResponse.ShippingRates, genericError
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
	jww.INFO.Printf("[makeRequest] - Request URL: %s", url)
	client := &http.Client{}
	buf := new(bytes.Buffer)

	if payload != "" {
		buf = bytes.NewBuffer([]byte(payload))
	}
	r, err := http.NewRequest(method, url, buf)

	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	r.Header.Add("X-Shopify-Access-Token", shopifyClient.shopifySecretToken)

	resp, err := client.Do(r)

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		jww.ERROR.Printf("[makeRequest] - 404 on executing request: %s\n", url)
	} else if resp.StatusCode == 429 {
		jww.ERROR.Printf("[makeRequest] - Rate limited!\n")
		rateLimitErr := errors.New("API rate limit exceeded")
		return rateLimitErr
	} else if resp.StatusCode == 422 {
		message, _ := ioutil.ReadAll(resp.Body)
		message = bytes.Replace(message, []byte("\\u003e"), []byte(">"), -1)
		jww.ERROR.Printf("[makeRequest] - %s\n", message)
		inventoryErr := errors.New(string(message))
		return inventoryErr
	}
	if err != nil {
		jww.ERROR.Printf("[makeRequest] - Error executing request : %s", err)
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(body)

	if err != nil {
		jww.ERROR.Printf("[makeRequest] - Decoding error: %#v", err)
		jww.ERROR.Printf("[makeRequest] - Response: %#v", resp.Body)
		return err
	}

	return nil
}
