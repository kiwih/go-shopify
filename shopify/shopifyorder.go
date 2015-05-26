/*

{
   "order":{
      "buyer_accepts_marketing":false,
      "cancel_reason":null,
      "cancelled_at":null,
      "cart_token":"aa28851379959e88db2b2c7b44d73ca4",
      "checkout_token":"70e2f8eea5509ba6f0292e1f8b8f0e19",
      "closed_at":null,
      "confirmed":true,
      "created_at":"2015-05-21T16:32:02+02:00",
      "currency":"USD",
      "device_id":null,
      "email":"sRossi@hotmail.it",
      "financial_status":"paid",
      "gateway":"shopify_payments",
      "id":544322051,
      "landing_site":"\/",
      "location_id":null,
      "name":"#1001",
      "note":null,
      "number":1,
      "processed_at":"2015-05-21T16:32:02+02:00",
      "reference":null,
      "referring_site":"https:\/\/arduino-us-dev-store.myshopify.com\/products\/arduino-starter-kit",
      "source_identifier":null,
      "source_url":null,
      "subtotal_price":"89.90",
      "taxes_included":false,
      "test":true,
      "token":"793311bb486b943ccbed5960c96a5c87",
      "total_discounts":"0.00",
      "total_line_items_price":"89.90",
      "total_price":"107.90",
      "total_price_usd":"107.90",
      "total_tax":"0.00",
      "total_weight":862,
      "updated_at":"2015-05-21T16:32:03+02:00",
      "user_id":null,
      "browser_ip":"82.112.219.187",
      "landing_site_ref":null,
      "order_number":1001,
      "discount_codes":[

      ],
      "note_attributes":[

      ],
      "processing_method":"direct",
      "source":"checkout_next",
      "checkout_id":871376195,
      "source_name":"web",
      "fulfillment_status":null,
      "tax_lines":[

      ],
      "tags":"",
      "line_items":[
         {
            "fulfillment_service":"manual",
            "fulfillment_status":null,
            "gift_card":false,
            "grams":862,
            "id":996468995,
            "price":"89.90",
            "product_id":591343875,
            "quantity":1,
            "requires_shipping":true,
            "sku":"7640152111310",
            "taxable":true,
            "title":"Arduino Starter Kit",
            "variant_id":1670212291,
            "variant_title":"",
            "vendor":"Arduino Store USA",
            "name":"Arduino Starter Kit",
            "variant_inventory_management":"shopify",
            "properties":[

            ],
            "product_exists":true,
            "fulfillable_quantity":1,
            "total_discount":"0.00",
            "tax_lines":[

            ]
         }
      ],
      "shipping_lines":[
         {
            "code":"International Shipping",
            "price":"18.00",
            "source":"shopify",
            "title":"International Shipping",
            "tax_lines":[

            ]
         }
      ],
      "billing_address":{
         "address1":"Via Roma 14",
         "address2":"",
         "city":"Torino",
         "company":"Arduino",
         "country":"Italy",
         "first_name":"Mario",
         "last_name":"Rossi",
         "latitude":45.071789,
         "longitude":7.64311,
         "phone":"",
         "province":"Torino",
         "zip":"1039",
         "name":"Mario Rossi",
         "country_code":"IT",
         "province_code":"TO"
      },
      "shipping_address":{
         "address1":"Via Roma 14",
         "address2":"",
         "city":"Torino",
         "company":"Arduino",
         "country":"Italy",
         "first_name":"Mario",
         "last_name":"Rossi",
         "latitude":45.071789,
         "longitude":7.64311,
         "phone":"",
         "province":"Torino",
         "zip":"1039",
         "name":"Mario Rossi",
         "country_code":"IT",
         "province_code":"TO"
      },
      "fulfillments":[

      ],
      "client_details":{
         "accept_language":"it-IT,it;q=0.8,en-US;q=0.6,en;q=0.4",
         "browser_height":969,
         "browser_ip":"82.112.219.187",
         "browser_width":1855,
         "session_hash":"c9486220e82a534bd2daa46d5fd65e1b90b9294abf25c8760f00814b4a1c7ae2",
         "user_agent":"Mozilla\/5.0 (X11; Linux x86_64) AppleWebKit\/537.36 (KHTML, like Gecko) Ubuntu Chromium\/41.0.2272.76 Chrome\/41.0.2272.76 Safari\/537.36"
      },
      "refunds":[

      ],
      "payment_details":{
         "avs_result_code":"Y",
         "credit_card_bin":"424242",
         "cvv_result_code":"M",
         "credit_card_number":"•••• •••• •••• 4242",
         "credit_card_company":"Visa"
      },
      "customer":{
         "accepts_marketing":false,
         "created_at":"2015-05-21T16:31:28+02:00",
         "email":"sRossi@hotmail.it",
         "first_name":"Mario",
         "id":520734019,
         "last_name":"Rossi",
         "last_order_id":null,
         "multipass_identifier":null,
         "note":null,
         "orders_count":0,
         "state":"disabled",
         "tax_exempt":false,
         "total_spent":"0.00",
         "updated_at":"2015-05-21T16:32:03+02:00",
         "verified_email":true,
         "tags":"",
         "last_order_name":null,
         "default_address":{
            "address1":"Via Roma 14",
            "address2":"",
            "city":"Torino",
            "company":"Arduino",
            "country":"Italy",
            "first_name":"Mario",
            "id":610183107,
            "last_name":"Rossi",
            "phone":"",
            "province":"Torino",
            "zip":"1039",
            "name":"Mario Rossi",
            "province_code":"TO",
            "country_code":"IT",
            "country_name":"Italy",
            "default":true
         }
      }
   }
}

*/

package shopify

// Address maps the shopify Order Customer Address
type Address struct {
	Name        string `json:"name"`
	Company     string `json:"company,omitempty"`
	Address1    string `json:"address1"`
	Address2    string `json:"address2,omitempty"`
	City        string `json:"city"`
	State       string `json:"province,omitempty"`
	PostalCode  string `json:"zip"`
	CountryCode string `json:"country_code,omitempty"`
	Telephone   string `json:"phone,omitempty"`
}

// Customer maps the shopify Order Customer
type Customer struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// Tax maps the shopify Order Shipping Taxes
type Tax struct {
	Rate  float64 `json:"rate"`
	Price string  `json:"price"`
}

// Shipping maps the shopify Order Shipping method
type Shipping struct {
	Shipper  string `json:"source"`
	Method   string `json:"title"`
	Price    string `json:"price"`
	TaxLines []Tax  `json:"tax_lines"`
}

// Item maps the shopify Order Items
type Item struct {
	Sku      string `json:"sku"`
	Quantity int    `json:"quantity"`
	//Data     string `json:"data"`
	Name      string `json:"name"`
	Price     string `json:"price"`
	TaxLines  []Tax  `json:"tax_lines"`
	VariantID int    `json:"variant_id"`
}

// Fulfillment maps the shopify Order Fulfillments
type Fulfillment struct {
	TrackingNumber string `json:"tracking_number,omitempty"`
	TrackingURL    string `json:"tracking_url,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"` //"2015-05-18T19:09:32-04:00"
}

// Order maps the shopify Order
type Order struct {
	ID                 int           `json:"id,omitempty"`
	Customer           Customer      `json:"customer"`
	Email              string        `json:"email,omitempty"`
	BillingAddress     Address       `json:"billing_address"`
	ShippingAddress    Address       `json:"shipping_address"`
	Shipping           []Shipping    `json:"shipping_lines"`
	DiscountCodes      []string      `json:"discount_codes,omitempty"`
	Note               string        `json:"note,omitempty"` //AutoDesk pspReferenceNumber here
	Items              []Item        `json:"line_items"`
	Fulfillments       []Fulfillment `json:"fulfillments,omitempty"`
	FulfillmentStatus  string        `json:"fulfillment_status,omitempty"`
	TotalPrice         string        `json:"total_price,omitempty"`
	TotalTax           string        `json:"total_tax,omitempty"`
	CreatedAt          string        `json:"created_at,omitempty"`   //"2015-05-18T19:09:32-04:00"
	CancelledAt        string        `json:"cancelled_at,omitempty"` //"2015-05-18T19:09:32-04:00"
	InventoryBehaviour string        `json:"inventory_behaviour,omitempty"`
	FinancialStatus    string        `json:"financial_status,omitempty"`
}

// OrderResponse models the shopify API response for order
type OrderResponse struct {
	SingleOrder Order `json:"order"`
}
