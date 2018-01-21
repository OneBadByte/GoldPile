package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"encoding/json"
	"io/ioutil"

    _ "github.com/gin-gonic/gin"
	"strconv"
)


type Pile struct{
	Name string `json:"name"`
	Total float64 `json:"total"`
	FileName string `json:"fileName"`
	Accounts []Account `json:"accounts"`
	BankAccounts []BankAccount `json:"bankAccounts"`
}

func CreatePlan(name string) Pile {

	pile := Pile{
		Name: name,
		FileName: name + ".json",
		Total: 0,
		Accounts: []Account{},
		BankAccounts: []BankAccount{},
	}
	pile.SavePlan()
	return pile
}

func (pile *Pile) LoadPlan(fileName string){
	fileData, _ := ioutil.ReadFile(fileName)
	err := json.Unmarshal(fileData, pile)
	if err != nil{
		fmt.Println(err)
	}
	pile.GetTotal()
}

func (pile *Pile) SavePlan(){
	jsonData, err := json.Marshal(pile)
	err = ioutil.WriteFile(pile.FileName, jsonData, 0777)
	if err != nil{
		fmt.Println(err)
	}
}

func (pile *Pile) AddAccount(account Account){
	if pile.CheckIfAccountExists(account.Name){
		return
	}
	pile.Accounts = append(pile.Accounts, account)
}

func (pile *Pile) DeleteAccount(accountName string){
	accountToRemove := pile.GetAccountsLocation(accountName)
	if accountToRemove == -1{
		fmt.Println("couldn't delete ", accountName)
		return
	}
	if accountToRemove == 0{
		pile.Accounts = pile.Accounts[1:]
	}else if accountToRemove == len(pile.Accounts)-1{
		pile.Accounts = pile.Accounts[0:len(pile.Accounts)-1]
	}else{
		pile.Accounts = append(pile.Accounts[:accountToRemove], pile.Accounts[accountToRemove+1:]...)
	}
}

func (pile Pile) CheckIfAccountExists(accountName string) bool{
	for _, account := range pile.Accounts{
		if account.Name == accountName{
			return true
		}
	}
	return false
}

func (pile Pile) GetAccountsLocation(accountName string) int{
	for key, account := range pile.Accounts{
		if account.Name == accountName{
			return key
		}
	}
	return -1
}

func (pile Pile) PrintOutAccounts(){
	for _, account := range pile.Accounts{
		text := fmt.Sprintf("the account name is: %s, Total: %v", account.Name, account.Total)
		fmt.Println(text)
		account.PrintOutCategories()
	}
}

func (pile *Pile) GetTotal(){
	pile.Total = 0
	for _, bankAccount := range pile.BankAccounts{
		pile.Total += bankAccount.Amount
	}
}


///////////////////////////////////////////////////////////////////////////////////////////

type BankAccount struct{
	Name string `json:"name"`
	Amount float64 `json:"amount"`
}

func CreateBankAccount(name string, amount float64) BankAccount{
	return BankAccount{name, amount}
}


func (pile *Pile) AddBankAccount(bankAccount BankAccount){
	if pile.CheckIfBankAccountExists(bankAccount.Name){
		return
	}
	pile.BankAccounts = append(pile.BankAccounts, bankAccount)
}

func (pile *Pile) DeleteBankAccount(bankAccountName string){
	bankAccountToRemove := pile.GetBankAccountsLocation(bankAccountName)
	if bankAccountToRemove == -1{
		fmt.Println("couldn't delete ", bankAccountName)
		return
	}
	if bankAccountToRemove == 0{
		pile.BankAccounts = pile.BankAccounts[1:]
	}else if bankAccountToRemove == len(pile.BankAccounts)-1{
		pile.BankAccounts = pile.BankAccounts[0:len(pile.BankAccounts)-1]
	}else{
		pile.BankAccounts = append(pile.BankAccounts[:bankAccountToRemove], pile.BankAccounts[bankAccountToRemove+1:]...)
	}
}

func (pile Pile) CheckIfBankAccountExists(bankAccountName string) bool{
	for _, bankAccount := range pile.BankAccounts{
		if bankAccount.Name == bankAccountName{
			return true
		}
	}
	return false
}

func (pile Pile) GetBankAccountsLocation(bankAccountName string) int{
	for key, bankAccount := range pile.BankAccounts{
		if bankAccount.Name == bankAccountName{
			return key
		}
	}
	return -1
}

func (pile Pile) PrintOutBankAccounts(){
	for _, bankAccount := range pile.BankAccounts{
		text := fmt.Sprintf("the Bank account name is: %s, Total: %v", bankAccount.Name, bankAccount.Amount)
		fmt.Println(text)
	}
}

///////////////////////////////////////////////////////////////////////////////////////////

type Account struct{
	Name       string `json:"name"`
	Total float64 `json:"total"`
	Categories []Category `json:"categories"`
}

func CreateAccount(name string) Account{
	return Account{name, 0, []Category{}}
}

func (account *Account) AddCategory(category Category){
	account.Categories = append(account.Categories, category)
}

func (account *Account) DeleteCategory(categoryName string){
	categoryToRemove := account.GetCategoryLocation(categoryName)
	if categoryToRemove == -1{
		fmt.Println("couldn't delete ", categoryName)
		return
	}
	fmt.Println("hit at ", categoryToRemove)
	if categoryToRemove == 0{
		account.Categories = account.Categories[1:]
	}else if categoryToRemove == len(account.Categories)-1{
		account.Categories = account.Categories[0:len(account.Categories)-1]
	}else{
		account.Categories = append(account.Categories[:categoryToRemove], account.Categories[categoryToRemove+1:]...)
	}
}

func (account Account) CheckIfCategoryExists(categoryName string) bool{
	for _, category := range account.Categories{
		if category.Name == categoryName{
			return true
		}
	}
	return false
}

func (account Account) GetCategoryLocation(categoryName string) int{
	for key, category := range account.Categories{
		if category.Name == categoryName{
			return key
		}
	}
	return -1
}

func (account Account) PrintOutCategories(){
	for _, category := range account.Categories{
		text := fmt.Sprintf("	the Categories name is: %s, Amount: %v", category.Name, category.Amount)
		fmt.Println(text)
	}
}

///////////////////////////////////////////////////////////////////////////////////////////

type Category struct{
	Name string `json:"name"`
	Amount float64 `json:"amount"`
}

func CreateCategory(name string) Category{
	return Category{name, 0}
}

func (category *Category) UpdateAmount(amount float64){
	category.Amount += amount
}

func GetInput() string{
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func GetAmountInput() float64{
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	amount, _ := strconv.ParseFloat(scanner.Text(), 64)
	return amount
}

func GetListOfCommands() []string{
	commands := GetInput()
	return strings.Split(commands, " ")

}

func PlanMenu() Pile {
	var plan Pile
	Loop:
	for {
		fmt.Println("Welcome!")
		fmt.Println("Press l to load a plan")
		fmt.Println("Press c to create a plan")
		input := GetInput()
		switch input {
		case "load":
			fallthrough
		case "l":
			fmt.Println("whats the name of the file? ")
			fileName := GetInput()
			plan.LoadPlan(fileName)
			break Loop

		case "create":
			fallthrough
		case "c":
			fmt.Print("whats the name a your plan? ")
			name := GetInput()
			newPlan := CreatePlan(name)
			return newPlan
			break Loop
		default:
			continue
		}
	}
	return plan
}

func (pile *Pile) MainMenu(){
	Loop:
	for {
		fmt.Println("Gold: $", pile.Total)
		pile.PrintOutAccounts()
		fmt.Println("")
		input := GetInput()
		switch input {
		case "add account":
			fmt.Print("whats the name of the new account? ")
			name := GetInput()
			pile.AddAccount(CreateAccount(name))
			pile.SavePlan()

		case "delete account":
			fmt.Print("what account do you want to delete? ")
			name := GetInput()
			pile.DeleteAccount(name)
			pile.SavePlan()

		case "add bank account":
			pile.PrintOutBankAccounts()
			fmt.Print("whats the name of the new bank account? ")
			name := GetInput()
			fmt.Print("how much is in there? ")
			amount := GetAmountInput()
			pile.AddBankAccount(CreateBankAccount(name, amount))
			pile.GetTotal()
			pile.SavePlan()

		case "delete bank account":
			pile.PrintOutBankAccounts()
			fmt.Print("what bank account do you want to delete? ")
			name := GetInput()
			pile.DeleteBankAccount(name)
			pile.GetTotal()
			pile.SavePlan()

		case "update bank account":
			fmt.Print("whats the name of the bank account to update? ")
			name := GetInput()
			fmt.Print("how much is in there? ")
			amount := GetAmountInput()
			pile.AddBankAccount(CreateBankAccount(name, amount))
			pile.SavePlan()


		case "add category":
			fmt.Print("whats the name of the account you want to add it to? ")
			accountName := GetInput()
			fmt.Print("whats the name a your category? ")
			name := GetInput()
			accountLocation := pile.GetAccountsLocation(accountName)
			pile.Accounts[accountLocation].AddCategory(CreateCategory(name))
			pile.Accounts[accountLocation].PrintOutCategories()
			pile.SavePlan()

		case "delete category":
			fmt.Print("whats the name of the account you want to delete from? ")
			accountName := GetInput()
			fmt.Print("whats the name a your category? ")
			name := GetInput()
			accountLocation := pile.GetAccountsLocation(accountName)
			pile.Accounts[accountLocation].DeleteCategory(name)
			pile.SavePlan()

		case "print categories":
			fmt.Print("what Accounts categories do you want printed? ")
			accountName := GetInput()
			fmt.Println("")
			accountLocation := pile.GetAccountsLocation(accountName)
			pile.Accounts[accountLocation].PrintOutCategories()

		case "help":
			fmt.Println("type add account to add a new account")
			fmt.Println("type delete account to delete an account")
			fmt.Println("type add category to create a category")
			fmt.Println("type delete category to delete a category")
			fmt.Println("type print categories to print categories")
			fmt.Println("Press c to create a pile")
			fmt.Println("type quit to quit")

		case "quit":
			break Loop
		}
	}


}

func main() {
	plan := PlanMenu()
	plan.MainMenu()

}
