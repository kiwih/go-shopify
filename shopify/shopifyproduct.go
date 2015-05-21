/*

{
"body_html": null,
"created_at": "2013-12-10T03:40:05-05:00",
"handle": "cheese",
"id": 191302485,
"product_type": "General",
"published_at": "2013-12-10T21:29:21-05:00",
"published_scope": "global",
"template_suffix": null,
"title": "Cheese",
"updated_at": "2014-02-11T19:37:59-05:00",
"vendor": "Not specified",
"tags": "",
"variants": [
{
"barcode": null,
"compare_at_price": null,
"created_at": "2013-12-10T03:40:05-05:00",
"fulfillment_service": "manual",
"grams": 0,
"id": 437938881,
"inventory_management": "shopify",
"inventory_policy": "deny",
"option1": "Default",
"option2": null,
"option3": null,
"position": 1,
"price": "23.00",
"product_id": 191302485,
"requires_shipping": true,
"sku": "10061",
"taxable": true,
"title": "Default",
"updated_at": "2013-12-10T03:40:05-05:00",
"inventory_quantity": 20,
"old_inventory_quantity": 20
}
],
"options": [
{
"id": 228536785,
"name": "Title",
"position": 1,
"product_id": 191302485
}
],
"images": [
{
"created_at": "2013-12-10T03:40:12-05:00",
"id": 391973353,
"position": 1,
"product_id": 191302505,
"updated_at": "2013-12-10T03:40:12-05:00",
"src": "http://cdn.shopify.com/s/files/1/0199/1250/products/df4331ec645329971ea45cdfb77315fead191d32.jpeg?v=1386664812"
}
],
"image": {
"created_at": "2013-12-10T03:40:12-05:00",
"id": 391973353,
"position": 1,
"product_id": 191302505,
"updated_at": "2013-12-10T03:40:12-05:00",
"src": "http://cdn.shopify.com/s/files/1/0199/1250/products/df4331ec645329971ea45cdfb77315fead191d32.jpeg?v=1386664812"
}
},
*/

package shopify

// Product models the shopify product model
type Product struct {
	BodyHTML       string           `json:"body_html,omitempty"`
	CreatedAt      string           `json:"created_at,omitempty"`
	Handle         string           `json:"handle,omitempty"`
	ID             int              `json:"id,omitempty"`
	ProductType    string           `json:"product_type,omitempty"`
	PublishedAt    string           `json:"published_at,omitempty"`
	PublishedScope string           `json:"published_scope,omitempty"`
	TemplateSuffix string           `json:"template_suffix,omitempty"`
	Title          string           `json:"title,omitempty"`
	UpdatedAt      string           `json:"updated_at,omitempty"`
	Vendor         string           `json:"vendor,omitempty"`
	Tags           string           `json:"tags,omitempty"`
	Variants       []ProductVariant `json:"variants,omitempty"`
	Options        []ProductOption  `json:"options,omitempty"`
	Images         []ProductImage   `json:"images,omitempty"`
	Image          ProductImage     `json:"image,omitempty"`
}

// ProductVariant models the shopify product variant
type ProductVariant struct {
	Barcode              string  `json:"barcode,omitempty"`
	CompareAtPrice       string  `json:"compare_at_price,omitempty"`
	CreatedAt            string  `json:"created_at,omitempty"`
	FulfillmentService   string  `json:"fulfillment_service,omitempty"`
	Grams                int     `json:"grams,omitempty"`
	ID                   int     `json:"id,omitempty"`
	InventoryManagement  string  `json:"inventory_management,omitempty"`
	InventoryPolicy      string  `json:"inventory_policy,omitempty"`
	Option1              string  `json:"option1,omitempty"`
	Option2              string  `json:"option2,omitempty"`
	Option3              string  `json:"option3,omitempty"`
	Position             int     `json:"position,omitempty"`
	Price                string  `json:"price,omitempty"`
	ProductID            int     `json:"product_id,omitempty"`
	RequiresShipping     bool    `json:"requires_shopping,omitempty"`
	Sku                  string  `json:"sku,omitempty"`
	Taxable              bool    `json:"taxable,omitempty"`
	Title                string  `json:"title,omitempty"`
	UpdatedAt            string  `json:"updated_at,omitempty"`
	InventoryQuantity    float64 `json:"inventory_quantity,omitempty"`
	OldInventoryQuantity int     `json:"old_inventory_quantity,omitempty"`
}

// ProductOption models the shopify product option
type ProductOption struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Position  int    `json:"position,omitempty"`
	ProductID int    `json:"product_id,omitempty"`
}

// ProductImage models the shopify product image object
type ProductImage struct {
	CreatedAt string `json:"created_at,omitempty"`
	ID        int    `json:"id,omitempty"`
	Position  int    `json:"position,omitempty"`
	ProductID int    `json:"product_id,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	Src       string `json:"src,omitempty"`
}

// Response models the shopify API response
type productResponse struct {
	Products      []Product `json:"products,omitempty"`
	SingleProduct Product   `json:"product,omitempty"`
}
