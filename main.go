package main

import (
	"fmt"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/vault/api"
)

func main() {
	config := api.DefaultConfig()
	config.Address = "https://dev-vault.fromnovatech.xyz/"

	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Unable to initialize Vault client: %v", err)
	}

	client.SetToken("example-token")

	tokenParams := map[string]any{
		"policies":     []string{"default", "my-app-policy"}, 
		"display_name": "my-app-token",
		"ttl":          "1h", 
		"renewable":    true,
	}

	secret, err := client.Logical().Write("auth/token/create", tokenParams)
	if err != nil {
		log.Fatalf("Error creating token: %v", err)
	}

	if secret == nil || secret.Auth == nil {
		log.Fatalf("No token information returned")
	}

	newToken := secret.Auth.ClientToken
	fmt.Printf("New token created: %s\n", newToken)
	fmt.Printf("Token Lease Duration: %d seconds\n", secret.Auth.LeaseDuration)

	router := gin.Default()

	router.Use(vaultAuthMiddleware(client))

	router.GET("/my-protected-route", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Access granted to protected route!",
		})
	})

	router.Run(":8080")
}

func vaultAuthMiddleware(client *api.Client) gin.HandlerFunc{
	return func (c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		client.SetToken(token)
		secret, err := client.Auth().Token().LookupSelf()

		if err != nil || secret == nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}