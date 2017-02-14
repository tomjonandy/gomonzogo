package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

var ACCESSTOKEN = os.Getenv("MONZOACCESS")

var client = &http.Client{
	Timeout: time.Duration(10 * time.Second),
}

type Account struct {
	ID          string    `json:"id"`
	Created     time.Time `json:"created"`
	Description string    `json:"description"`
}

type AccountsResponse struct {
	Accounts []Account `json:"accounts"`
}

type Transactions struct {
	Transactions []struct {
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
	} `json:"transactions"`
}

func GetAccounts() []Account {
	req, err := http.NewRequest("GET", "https://api.monzo.com/accounts", nil)
	req.Header.Add("Authorization", "Bearer "+ACCESSTOKEN)

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return []Account{}
	}

	defer res.Body.Close()

	body := AccountsResponse{}

	json.NewDecoder(res.Body).Decode(&body)
	return body.Accounts
}

func GetTransactions(accountID string) Transactions {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.monzo.com/transactions?account_id=%s", accountID), nil)
	req.Header.Add("Authorization", "Bearer "+ACCESSTOKEN)

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return Transactions{}
	}

	defer res.Body.Close()

	body := Transactions{}

	json.NewDecoder(res.Body).Decode(&body)
	return body
}

func main() {
	accounts := GetAccounts()
	for i := 0; i < len(accounts); i++ {
		fmt.Printf("[%d] %s\n", i, accounts[i].Description)
	}

	fmt.Println("Enter number...")
	var accountNumber int = 0
	fmt.Scanln(&accountNumber)
	transactions := GetTransactions(accounts[accountNumber].ID)
	for i := 0; i < len(transactions.Transactions); i++ {
		fmt.Printf("[%d]\t%d\t%s\n", i, transactions.Transactions[i].Amount, transactions.Transactions[i].Description)
	}
}
