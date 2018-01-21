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


type Plan struct{
	Name string `json:"name"`
	Total float64 `json:"total"`
	FileName string `json:"fileName"`
	Accounts []Account `json:"accounts"`
	BankAccounts []BankAccount `json:"bankAccounts"`
}

func (plan *Plan) LoadPlan(fileName string){
	fileData, _ := ioutil.ReadFile(fileName)
	err := json.Unmarshal(fileData, plan)
	if err != nil{
		fmt.Println(err)
	}
}

func (plan *Plan) SavePlan(){
	jsonData, err := json.Marshal(plan)
	err = ioutil.WriteFile(plan.FileName, jsonData, 0777)
	if err != nil{
		fmt.Println(err)
	}
}

func (plan *Plan) AddAccount(account Account){
	if plan.CheckIfAccountExists(account.Name){
		return
	}
	plan.Accounts = append(plan.Accounts, account)
}

func (plan *Plan) DeleteAccount(accountName string){
	accountToRemove := plan.GetAccountsLocation(accountName)
	if accountToRemove == -1{
		fmt.Println("couldn't delete ", accountName)
		return
	}
	fmt.Println("hit at ", accountToRemove)
	if accountToRemove == 0{
		plan.Accounts = plan.Accounts[1:]
	}else if accountToRemove == len(plan.Accounts)-1{
		plan.Accounts = plan.Accounts[0:len(plan.Accounts)-1]
	}else{
		plan.Accounts = append(plan.Accounts[:accountToRemove], plan.Accounts[accountToRemove+1:]...)
	}
}

func (plan Plan) CheckIfAccountExists(accountName string) bool{
	for _, account := range plan.Accounts{
		if account.Name == accountName{
			return true
		}
	}
	return false
}

func (plan Plan) GetAccountsLocation(accountName string) int{
	for key, account := range plan.Accounts{
		if account.Name == accountName{
			return key
		}
	}
	return -1
}

func (plan Plan) PrintOutAccounts(){
	for _, account := range plan.Accounts{
		text := fmt.Sprintf("the account name is: %s, Total: %v", account.Name, account.Total)
		fmt.Println(text)
	}
}

///////////////////////////////////////////////////////////////////////////////////////////

type BankAccount struct{
	Name string `json:"name"`
	Amount string `json:"amount"`
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

func (account Account) PrintOutACategories(){
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

func PlanMenu() Plan{
	var plan Plan
	fmt.Println("Welcome!")
	fmt.Println("Press l to load a plan")
	fmt.Println("Press c to create a plan")
	input := GetInput()
	switch input{
	case "l":
		fmt.Println("whats the name of the file? ")
		fileName := GetInput()
		plan.LoadPlan(fileName)

	case "c":
		fmt.Print("whats the name a your plan? ")
		plan.Name = GetInput()
		plan.FileName = plan.Name + ".json"
		plan.Total = 0
		plan.Accounts = []Account{}
		plan.BankAccounts = []BankAccount{}
		plan.SavePlan()
	}
	return plan
}

func (plan *Plan) MainMenu(){
	Loop:
	for {
		plan.PrintOutAccounts()
		fmt.Println("")
		fmt.Println("type add account to add a new account")
		fmt.Println("type delete account to delete an account")
		fmt.Println("type add category to create a category")
		fmt.Println("type delete category to create a category")
		fmt.Println("type print categories to print categories")
		fmt.Println("Press c to create a plan")
		fmt.Println("type quit to quit")
		input := GetInput()
		switch input {
		case "add account":
			fmt.Print("whats the name of the new account? ")
			name := GetInput()
			plan.AddAccount(CreateAccount(name))
			plan.SavePlan()

		case "delete account":
			fmt.Print("what account do you want to delete? ")
			name := GetInput()
			plan.AddAccount(CreateAccount(name))
			plan.SavePlan()

		case "add category":
			fmt.Print("whats the name of the account you want to add it to? ")
			accountName := GetInput()
			fmt.Print("whats the name a your category? ")
			name := GetInput()
			accountLocation := plan.GetAccountsLocation(accountName)
			plan.Accounts[accountLocation].AddCategory(CreateCategory(name))
			plan.Accounts[accountLocation].PrintOutACategories()
			plan.SavePlan()

		case "print categories":
			fmt.Print("what Accounts categories do you want printed? ")
			accountName := GetInput()
			fmt.Println("")
			accountLocation := plan.GetAccountsLocation(accountName)
			plan.Accounts[accountLocation].PrintOutACategories()



		case "quit":
			break Loop
		}
	}


}

func main() {
	plan := PlanMenu()
	plan.MainMenu()

}
