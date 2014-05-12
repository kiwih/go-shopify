package main

import (
	"fmt"
	"github.com/hammond-bones/go-shopify/shopify"
	"strconv"
)

func main() {
	fmt.Printf("Hello!\n")
	shop := shopify.NewClient("mcglynn-quitzon-and-windler8990", "98bfc43d6aa771567326d76e6395173d")

	shop.LoadProducts()

	fmt.Printf("%v\n", shop.Products[0].Id)

	prod := shop.GetLiveProduct(strconv.Itoa(shop.Products[0].Id))

	fmt.Printf("%v\n", prod.Handle)
}
