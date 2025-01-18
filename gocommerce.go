// Package gocommerce provides Go bindings for the SquareSpace v1 Commerce APIs .
//
// This package serves as the root module for gocommerce. Users can access
// specific functionality by importing subpackages such as inventory, orders,
// products, profiles, transactions, and webhooks.
//
// Example:
//  import (
//  	"context"
//  	"fmt"
//  	"net/http"

//  	"github.com/j-low/gocommerce/products"
// 	)

// 	config := common.Config{
//  	APIKey:      "my_api_key-999",
//  	UserAgent:   "my_user-agent-999",
//  	Client:      http.DefaultClient,
// 	}

//		resp, err := products.DeleteProduct(ctx, &config, "some-product-id-999")
//		if err != nil {
//	   return fmt.Println(err)
//		}
package gocommerce