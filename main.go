package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

var ACCESSTOKEN = os.Getenv("MONZOACCESS")

type MyClient struct {
	http.Client
}

func (client *MyClient) DoWithHeader(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+ACCESSTOKEN)
	return client.Do(req)
}

var client = &MyClient {
	http.Client {
		Timeout: time.Duration(10 * time.Second),
	},
}

type Account struct {
	ID          string    `json:"id"`
	Created     time.Time `json:"created"`
	Description string    `json:"description"`
}

type AccountsResponse struct {
	Accounts []Account `json:"accounts"`
}

type Transaction struct {
	ID          string      `json:"id"`
	Created     time.Time   `json:"created"`
	Description string      `json:"description"`
	Amount      int         `json:"amount"`
	Currency    string      `json:"currency"`
	Merchant    interface{} `json:"merchant"`
	Notes       string      `json:"notes"`
	Metadata    struct {
		IsTopup string `json:"is_topup"`
	} `json:"metadata"`
	AccountBalance int           `json:"account_balance"`
	Attachments    []interface{} `json:"attachments"`
	Category       string        `json:"category"`
	IsLoad         bool          `json:"is_load"`
	Settled        string        `json:"settled"`
	LocalAmount    int           `json:"local_amount"`
	LocalCurrency  string        `json:"local_currency"`
	Updated        string        `json:"updated"`
	AccountID      string        `json:"account_id"`
	Counterparty   struct {
	} `json:"counterparty"`
	Scheme            string `json:"scheme"`
	DedupeID          string `json:"dedupe_id"`
	Originator        bool   `json:"originator"`
	IncludeInSpending bool   `json:"include_in_spending"`
	DeclineReason     string `json:"decline_reason,omitempty"`
}

type TransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
}

type BalanceResponse struct {
	Balance int `json:"balance"`
	Currency string `json:"currency"`
	SpendToday int `json:"spend_today"`
	LocalCurrency string `json:"local_currency"`
	LocalExchangeRate int `json:"local_exchange_rate"`
	LocalSpend []struct {
		SpendToday int `json:"spend_today"`
		Currency string `json:"currency"`
	} `json:"local_spend"`
}

func GetAccounts() []Account {
	req, err := http.NewRequest("GET", "https://api.monzo.com/accounts", nil)

	res, err := client.DoWithHeader(req)
	if err != nil {
		fmt.Println(err)
		return []Account{}
	}

	defer res.Body.Close()

	body := AccountsResponse{}

	json.NewDecoder(res.Body).Decode(&body)
	return body.Accounts
}

func GetBalance(accountID string) int {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.monzo.com/balance?account_id=%s", accountID), nil)

	res, err := client.DoWithHeader(req)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	defer res.Body.Close()

	body := BalanceResponse{}

	json.NewDecoder(res.Body).Decode(&body)
	return body.Balance
}

func GetTransactions(accountID string) []Transaction {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.monzo.com/transactions?account_id=%s", accountID), nil)

	res, err := client.DoWithHeader(req)

	if err != nil {
		fmt.Println(err)
		return []Transaction{}
	}

	defer res.Body.Close()

	body := TransactionsResponse{}

	json.NewDecoder(res.Body).Decode(&body)
	return body.Transactions
}

func main() {
	accounts := GetAccounts()
	for i := 0; i < len(accounts); i++ {
		fmt.Printf("[%d] %s\t%d\n", i, accounts[i].Description, GetBalance(accounts[i].ID))
	}

	fmt.Println("Enter number...")
	var accountNumber int = 0
	fmt.Scanln(&accountNumber)
	transactions := GetTransactions(accounts[accountNumber].ID)
	for i := 0; i < len(transactions); i++ {
		fmt.Printf("[%d]\t%d\t%s\n", i, transactions[i].Amount, transactions[i].Description)
	}
}
