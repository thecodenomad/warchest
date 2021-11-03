package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
	"warchest/src/auth"
	// "warchest/src/config"
	"warchest/src/query"
)

// FailedLoadConfigRC Return code for failing to load the Warchest configuration
const FailedLoadConfigRC = 2

// FailedRetrievingData Return code for failing the retrieval of coin data
const FailedRetrievingData = 3

// FailedCalculatingWallet Return code for failing calculation of wallet
const FailedCalculatingWallet = 3

//
// Env Variables
////////////////////

// CbAPIKey is the api key established via your coinbase profile
const CbAPIKey = "CB_API_KEY"

// CbAPISecret is the secret associated with the cbAPIKey
const CbAPISecret = "CB_API_SECRET"

// WarchestStaticPath is the env var that defines where public/static files are being served from
const WarchestStaticPath = "WARCHEST_STATIC_PATH"

// WarchestConfigEnv is the environment variable that will point to coin transactions used by Warchest
const WarchestConfigEnv = "WARCHEST_CONFIG"

var (
	once           sync.Once
	warchestWallet *query.Wallet
)

// GetWalletSingleton will retrieve the wallet singleton used by the application
// TODO: this should take in a new flag to specify whether or not to use local config for the transaction
//       base
func GetWalletSingleton() *query.Wallet {

	// TODO: Could this be a singleton?
	client := http.Client{
		Timeout: time.Second * 10,
	}
	var absClient query.HTTPClient
	absClient = &client
	apiKey := os.Getenv(CbAPIKey)
	apiSecret := os.Getenv(CbAPISecret)
	cbAuth := auth.CBAuth{apiKey, apiSecret}

	// Instantiate the object since it doesn't exist
	once.Do(func() {
		log.Printf("Wallet is being instantiated now")
		warchestWallet = &query.Wallet{Coins: map[string]query.WarchestCoin{}, NetProfit: 0.0}

		// Retreive coins for account
		coins, err := query.GetWarchestCoins(cbAuth, absClient)
		if err != nil {
			log.Printf("Failed to retrieve Warchest Coins: %s\n", err)
			return
		}

		log.Printf("There are %d coins in this wallet", len(coins))

		warchestWallet.Coins = coins
	})
	warchestWallet.UpdateNetProfit(cbAuth, absClient)
	return warchestWallet
}

// GetWallet API Endpoint to retrieve a wallet
func GetWallet(c *gin.Context) {
	warchestWallet := GetWalletSingleton()
	if warchestWallet == nil {
		log.Printf("Warchest wallet must be instantiated at runtime prior to this call!")
	}

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	c.IndentedJSON(http.StatusOK, warchestWallet)
}

func main() {

	// Args
	serverPtr := flag.Bool("server", false, "whether or not to start server (default port: 8080)")
	savePtr := flag.Bool("save", true, "whether or not to save a list of transactions")
	transactionTypePtr := flag.String("transaction-type", "all", "the type of coin to parse transactions against")

	// Parse the argument flags
	flag.Parse()

	//// Helper vars
	//supportedCoins := getSupportedCoins()

	client := http.Client{
		Timeout: time.Second * 10,
	}
	var absClient query.HTTPClient
	absClient = &client

	fmt.Println("Server enabled:", *serverPtr)
	fmt.Println("Save enabled:", *savePtr)
	fmt.Println("Transaction type:", *transactionTypePtr)

	apiKey := os.Getenv(CbAPIKey)
	apiSecret := os.Getenv(CbAPISecret)
	cbAuth := auth.CBAuth{apiKey, apiSecret}

	// Retrieve all available wallets for the account associated with the provided API Key
	accountsResp, _ := query.CBRetrieveAccounts(cbAuth, absClient)

	// Filter out non-zero amount of individual coin types
	// Restrict to supported coins only
	//valueMap := map[string]query.WarchestCoin{}
	//for _, account := range accountsResp.Accounts {
	//	if IsSupportedCoin(account.Currency.Code, supportedCoins) {
	//		if account.Balance.Amount > 0.0 {
	//			warchestCoin := query.WarchestCoin{AccountID: account.ID, Symbol: account.Currency.Code}
	//
	//			// Update the coin's rates, profit, and cost
	//			warchestCoin.Update(cbAuth, absClient)
	//			valueMap[account.Currency.Code] = warchestCoin
	//		}
	//	}
	//}

	log.Printf("There are %d types of coins.\n", len(accountsResp.Accounts))
	//log.Printf("You have %d type(s) of coin(s) in your wallet:\n", len(valueMap))
	//for coinSymbol, wcCoin := range valueMap {
	//	log.Printf("\t%s\n", coinSymbol)
	//	log.Printf("\t\tNum Coins: %.14f\n", wcCoin.Amount)
	//	log.Printf("\t\tCost: %.14f\n", wcCoin.Cost)
	//	log.Printf("\t\tProfit: %.14f\n", wcCoin.Profit)
	//	log.Printf("\t\tCurrent Value: %.14f\n", wcCoin.Amount*wcCoin.Rates.USD)
	//}

	// Setup server
	if *serverPtr {

		// Establish the static path, defaulting to public folder in current execution path
		staticPath, ok := os.LookupEnv(WarchestStaticPath)
		if !ok {
			log.Printf("WARCHEST_STATIC_PATH not set, using default ./public for serving static files")
			staticPath = "./public"
		}

		router := gin.Default()
		router.Static("/static", staticPath)
		router.Static("/js", staticPath+"/js")
		router.Static("/css", staticPath+"/css")
		//router.StaticFile("/favicon.ico", staticPath+"/favicon.ico")

		// Simple API to check if server is working correctly
		router.GET("/api/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		// Setup redirect for static files
		router.GET("/", func(c *gin.Context) {
			c.Redirect(301, "/static/index.html")
		})

		// Setup Basic call to retrieve wallet
		router.GET("/api/wallet", GetWallet)

		router.Run()
	}
}
