# gocommerce

Go bindings for [Squarespace Commerce APIs v1](https://developers.squarespace.com/commerce-apis/overview)

## Requirements

[Go 1.22](https://go.dev/doc/install) or higher

## Available APIs

The package provides bindings for the following Squarespace Commerce APIs:

- [Inventory](https://developers.squarespace.com/commerce-apis/inventory-overview)
- [Orders](https://developers.squarespace.com/commerce-apis/orders-overview)
- [Products](https://developers.squarespace.com/commerce-apis/products-overview)
- [Profiles](https://developers.squarespace.com/commerce-apis/profiles-overview)
- [Transactions](https://developers.squarespace.com/commerce-apis/transactions-overview)
- [Webhook Subscriptions](https://developers.squarespace.com/commerce-apis/webhook-subscriptions-overview)

## Installation

From your project root, run:

```
go get github.com/j-low/gocommerce@v1.0.1
```

## Usage

```
import (
  "context"
  "fmt"
  "net/http"

  "github.com/j-low/gocommerce/products"
)

config := common.Config{
  APIKey:      "my_api_key-999",
  UserAgent:   "my_user-agent-999",
  Client:      http.DefaultClient,
}

resp, err := products.DeleteProduct(ctx, &config, "some-product-id-999")
if err != nil {
  return fmt.Println(err)
}
```

## License

MIT License
