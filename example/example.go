package main

import (
	"fmt"
	"strconv"

	"github.com/arduino/go-shopify/shopify"
)

func main() {
	fmt.Printf("Hello!\n")
	shop := shopify.NewClient("mcglynn-quitzon-and-windler8990", "98bfc43d6aa771567326d76e6395173d")

	shop.LoadProducts()

	fmt.Printf("%v\n", shop.Products[0].ID)

	prod := shop.GetLiveProduct(strconv.Itoa(shop.Products[0].ID))

	order := shop.GetOrder("123456")

	fmt.Printf("%v\n", prod.Handle)
	fmt.Printf("%v\n", order)

	removedOrder := shop.CancelOrder("123456")

	fmt.Printf("%v\n", removedOrder)

	// shop.PlaceOrder(order)
	// shop.ShippingOptions(order)
}
