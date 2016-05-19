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

/*
{
	"token": "f04c38b0b44ba72d222034452ce723ff",
	"note": null,
	"attributes": {},
	"total_price": 500,
	"total_discount": 0,
	"total_weight": 3997,
	"item_count": 10,
	"items": [{
		"id": 2171133763,
		"properties": null,
		"quantity": 10,
		"variant_id": 2171133763,
		"title": "Logistic Services for AutoDesk Arduino Basic Kit",
		"price": 50,
		"line_price": 500,
		"total_discount": 0,
		"discounts": [],
		"sku": "AKX00001",
		"grams": 400,
		"vendor": "Arduino LLC",
		"product_id": 769682435,
		"gift_card": false,
		"url": "\/products\/logistic-services-for-autodesk-arduino-basic-kit?variant=2171133763",
		"image": "https:\/\/cdn.shopify.com\/s\/files\/1\/0870\/7082\/products\/AKX00001_Bk_3_Front.jpg?v=1433317209",
		"handle": "logistic-services-for-autodesk-arduino-basic-kit",
		"requires_shipping": true,
		"product_type": "",
		"product_title": "Logistic Services for AutoDesk Arduino Basic Kit",
		"product_description": "This kit contains all the components you need to build simple projects and learn how to turn an idea into reality using Arduino. At arduino.cc\/basicstarterkit you will find the step-by-step tutorials to realise 15 simple projects.",
		"variant_title": null,
		"variant_options": ["Default Title"]
	}],
	"requires_shipping": true
}
*/

package shopify

// CartResponse models the shopify API response for cart
type CartResponse struct {
	Token string `json:"token"`
	Note  string `json:"note,omitempty"`
}
