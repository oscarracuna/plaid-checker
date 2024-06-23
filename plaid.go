package main

import (
	"bufio"

	"context"

	"fmt"

	"os"

	"strings"

	"github.com/joho/godotenv"
	"github.com/plaid/plaid-go/v25/plaid"
)

var (
	query string

	answer string

	reset = "\033[0m"

	red = "\033[31m"

	green = "\033[32m"
)

func main() {

	godotenv.Load()
	saludo := `

 ______   _        ______  _____  _____  
| |  | \ | |      | |  | |  | |  | | \ \  
| |__|_/ | |   _  | |__| |  | |  | |  | |
|_|      |_|__|_| |_|  |_| _|_|_ |_|_/_/  

                                         

    `

	saludo2 := `

_______  _    _  _____  ______   ______  __    _  
  | |   | |  | |  | |  | |  \ \ | | ____ \ \  | |
  | |   | |--| |  | |  | |  | | | |  | |  \_\_| |
  |_|   |_|  |_| _|_|_ |_|  |_| |_|__|_|  ____|_|

                                                 

    `

	fmt.Print(saludo, saludo2, "\n")

	fmt.Println("==================================================")

	fmt.Println("Welcome to the plaid checker for helpdesk.")

	fmt.Println("==================================================")

	reader := bufio.NewReader(os.Stdin)

prompt:

	fmt.Print("Please enter the bank name (max 10 results): ")

	query, _ = reader.ReadString('\n')

	query = strings.TrimSpace(query)

	// ==================================================================

	// This is required for the API to work

	configuration := plaid.NewConfiguration()

	// This is the key and secret, please do not hack me

	client_id := os.Getenv("PLAID_CLIENT_ID")
	secret := os.Getenv("PLAID_SECRET")
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", client_id)

	configuration.AddDefaultHeader("PLAID-SECRET", secret)

	// This states that we're using sandbox, so we're not messing with programmer's API

	// Or real customer data

	configuration.UseEnvironment(plaid.Sandbox)

	ctx := context.Background()

	//===================================================================

	//This section was a bunch of trial and error

	//It 100% can be improved

	countryCodes := []plaid.CountryCode{plaid.COUNTRYCODE_US}

	client := plaid.NewAPIClient(configuration)

	request := plaid.NewInstitutionsSearchRequest(query, countryCodes)

	request.SetProducts([]plaid.Products{plaid.PRODUCTS_AUTH})

	response, _, err := client.PlaidApi.InstitutionsSearch(ctx).InstitutionsSearchRequest(*request).Execute()

	if err != nil {

		fmt.Println("API may be down. Check status on Plaid.", err)

	}

	//We loop over the institutions with the parameters stated in request

	for _, banks := range response.Institutions {

		bank := banks.GetName()

		auth := banks.GetProducts()

		//Since GetProducts() returns an array, and AUTH is always index 1, we check if it's there

		if auth[1] == "auth" {

			fmt.Println(green + "\nAuth found for this bank: " + reset)

			fmt.Printf("Name: %s\n", bank)

		} else {

			fmt.Println(red + "\nThis bank does NOT have auth: " + reset)

			fmt.Println(bank)

			break

		}

	}

	//Loop in case we'd like to double check or look up another bank

	fmt.Println("\nWould you like to look up another bank? (y/n)")

	fmt.Print("Answer: ")

	// Had to switch from fmt.scan to this since it was causing issues due to the

	// \n line skips from the for loop above

	answer, _ = reader.ReadString('\n')

	answer = strings.TrimSpace(answer)

	for {

		if answer == "y" || answer == "yes" {

			goto prompt

		} else {

			fmt.Println("==================================================")

			fmt.Println("Thank you for using this app. Stay safe out there!")

			fmt.Println("==================================================")

			break

		}

	}

	fmt.Print("Press 'Enter' to exit...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')

}
