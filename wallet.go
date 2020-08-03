package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/templates"
	"google.golang.org/grpc"
	"io/ioutil"
	"os"
	"time"
)

type Account struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
	SigAlgo    string `json:"sigAlgorithm"`
	HashAlgo   string `json:"hashAlgorithm"`
}

type wallet struct {
	Accounts struct {
		Seller      Account
		Buyer       Account
		Ft          Account
		Nft         Account
		Marketplace Account
		Service     Account
	}
}

func walletAccounts() wallet {
	f, err := os.Open("./wallet/accounts.json")
	if err != nil {
		panic(err)
	}

	d := json.NewDecoder(f)

	var accountsInWallet wallet

	err = d.Decode(&accountsInWallet)
	if err != nil {
		panic(err)
	}

	return accountsInWallet
}

func ReadFile(path string) []byte {
	contents, err := ioutil.ReadFile(path)
	Handle(err)
	return contents
}

func Handle(err error) {
	if err != nil {
		panic(err)
	}
}

func AccountInfo(account Account) (crypto.PrivateKey, crypto.SignatureAlgorithm, crypto.HashAlgorithm) {
	sigAlgo := crypto.StringToSignatureAlgorithm(account.SigAlgo)
	hashAlgo := crypto.StringToHashAlgorithm(account.HashAlgo)
	privateKey, err := crypto.DecodePrivateKeyHex(sigAlgo, account.PrivateKey)
	Handle(err)

	return privateKey, sigAlgo, hashAlgo
}

func CreateAccount(node string, user Account, service Account, code []byte) string {
	ctx := context.Background()

	// User account
	privateKey, sigAlgo, hashAlgo := AccountInfo(user)
	publicKey := privateKey.PublicKey()

	accountKey := flow.NewAccountKey().
		SetPublicKey(publicKey).
		SetSigAlgo(sigAlgo).
		SetHashAlgo(hashAlgo).
		SetWeight(flow.AccountKeyWeightThreshold)

	// Service account
	servicePrivateKey, _, _ := AccountInfo(service)
	serviceAddress := flow.HexToAddress(service.Address)

	c, err := client.New(node, grpc.WithInsecure())
	Handle(err)

	serviceAccount, err := c.GetAccountAtLatestBlock(ctx, serviceAddress)
	if err != nil {
		panic(err)
	}
	serviceAccountKey := serviceAccount.Keys[0]
	serviceSigner := crypto.NewInMemorySigner(servicePrivateKey, serviceAccountKey.HashAlgo)

	tx := templates.CreateAccount([]*flow.AccountKey{accountKey}, code, serviceAddress)
	tx.SetProposalKey(serviceAddress, serviceAccountKey.ID, serviceAccountKey.SequenceNumber)
	tx.SetPayer(serviceAddress)
	tx.SetGasLimit(uint64(100))

	err = tx.SignEnvelope(serviceAddress, serviceAccountKey.ID, serviceSigner)
	Handle(err)

	err = c.SendTransaction(ctx, *tx)
	Handle(err)

	blockTime := 10 * time.Second
	time.Sleep(blockTime)

	result, err := c.GetTransactionResult(ctx, tx.ID())
	Handle(err)

	var address flow.Address

	if result.Status == flow.TransactionStatusSealed {
		for _, event := range result.Events {
			if event.Type == flow.EventAccountCreated {
				accountCreatedEvent := flow.AccountCreatedEvent(event)
				address = accountCreatedEvent.Address()
			}
		}
	}

	return address.Hex()
}

func CreateTokenAccounts() {
	node := "127.0.0.1:3569"
	wallet := walletAccounts()

	ftContract := ReadFile("./contracts/FT.cdc")
	ftAddresss := CreateAccount(node, wallet.Accounts.Ft, wallet.Accounts.Service, ftContract)
	fmt.Println("Fungible token deployed at address:", ftAddresss)

	nftContract := ReadFile("./contracts/NFT.cdc")
	nftAddresss := CreateAccount(node, wallet.Accounts.Nft, wallet.Accounts.Service, nftContract)
	fmt.Println("Non-fungible token deployed at address:", nftAddresss)
}

func CreateMarketplaceAccount() {
	node := "127.0.0.1:3569"
	wallet := walletAccounts()
	mktContract := ReadFile("./contracts/Marketplace.cdc")
	mktAddress := CreateAccount(node, wallet.Accounts.Marketplace, wallet.Accounts.Service, mktContract)
	fmt.Println("Marketplace deployed at address:", mktAddress)
}

func CreateBuyerAndSellerAccounts() {
	node := "127.0.0.1:3569"
	wallet := walletAccounts()

	buyerAddress := CreateAccount(node, wallet.Accounts.Buyer, wallet.Accounts.Service, nil)
	fmt.Println("Buyer account address:", buyerAddress)

	sellerAddress := CreateAccount(node, wallet.Accounts.Seller, wallet.Accounts.Service, nil)
	fmt.Println("Seller account address:", sellerAddress)
}

func ShowAccoundCode(address string) {
	node := "127.0.0.1:3569"

	c, err := client.New(node, grpc.WithInsecure())
	Handle(err)

	ctx := context.Background()

	account, err := c.GetAccountAtLatestBlock(ctx, flow.HexToAddress(address))
	Handle(err)

	fmt.Println(string(account.Code))
}

func ExecuteTransaction(user Account, code []byte) {
	node := "127.0.0.1:3569"

	c, err := client.New(node, grpc.WithInsecure())
	Handle(err)

	ctx := context.Background()

	account, err := c.GetAccount(ctx, flow.HexToAddress(user.Address))
	Handle(err)

	key := account.Keys[0]

	tx := flow.NewTransaction().
		SetScript(code).
		SetGasLimit(100).
		SetProposalKey(account.Address, key.ID, key.SequenceNumber).
		SetPayer(account.Address).
		AddAuthorizer(account.Address)

	privateKey, _, _ := AccountInfo(user)
	signer := crypto.NewInMemorySigner(privateKey, key.HashAlgo)
	err = tx.SignEnvelope(account.Address, key.ID, signer)
	Handle(err)

	err = c.SendTransaction(ctx, *tx)
	Handle(err)

	blockTime := 10 * time.Second
	time.Sleep(blockTime)

	result, err := c.GetTransactionResult(ctx, tx.ID())
	Handle(err)

	fmt.Println("Transaction status: ", result.Status.String())
}

func ExecuteScript(code []byte) {
	node := "127.0.0.1:3569"

	c, err := client.New(node, grpc.WithInsecure())
	Handle(err)

	ctx := context.Background()
	result, err := c.ExecuteScriptAtLatestBlock(ctx, code, nil)
	Handle(err)

	fmt.Println("Script result: ", result)
}

func SetupBuyerAccount() {
	wallet := walletAccounts()
	code := ReadFile("./transactions/BuyerSetup.cdc")
	ExecuteTransaction(wallet.Accounts.Buyer, code)
}

func DepositFTsIntoBuyersAccount() {
	wallet := walletAccounts()
	code := ReadFile("./transactions/DepositFT.cdc")
	ExecuteTransaction(wallet.Accounts.Ft, code)
}

func SetupSellerAccount() {
	wallet := walletAccounts()
	code := ReadFile("./transactions/SellerSetup.cdc")
	ExecuteTransaction(wallet.Accounts.Seller, code)
}

func DepositNFTIntoSellersAccount() {
	wallet := walletAccounts()
	code := ReadFile("./transactions/DepositNFT.cdc")
	ExecuteTransaction(wallet.Accounts.Nft, code)
}

func ListNFTForSale() {
	wallet := walletAccounts()
	code := ReadFile("./transactions/ListNFTForSale.cdc")
	ExecuteTransaction(wallet.Accounts.Seller, code)
}

func PurchaseNFT() {
	wallet := walletAccounts()
	code := ReadFile("./transactions/PurchaseNFT.cdc")
	ExecuteTransaction(wallet.Accounts.Buyer, code)
}

func CheckAccounts() {
	code := ReadFile("./scripts/CheckAccounts.cdc")
	ExecuteScript(code)
}

func main() {
	// CreateTokenAccounts()
	// CreateMarketplaceAccount()
	// CreateBuyerAndSellerAccounts()
	// SetupBuyerAccount()
	// DepositFTsIntoBuyersAccount()
	// SetupSellerAccount()
	// DepositNFTIntoSellersAccount()
	// ListNFTForSale()
	// PurchaseNFT()
	CheckAccounts()
	//ShowAccoundCode("ff8975b2fe6fb6f1")

}
