package main

import (
	"os"
	"fmt"
	"log"
	"net/http"
	"encoding/json"

	"github.com/PuerkitoBio/goquery"
)

type RecipeIngredient struct {
	Name	string	`json:"name"`
	Amount	string	`json:"amount"`
	Unit	string	`json:"unit"`
}

func scrapeIngredients(url string) []*RecipeIngredient {
	var ingredientList []*RecipeIngredient
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Error Code: %d \t Status:%s\n", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".wprm-recipe-ingredient-group .wprm-recipe-ingredient").Each(
		func(i int, selection *goquery.Selection) {
		   amount := selection.Find(".wprm-recipe-ingredient-amount").Text()
           unit := selection.Find(".wprm-recipe-ingredient-unit").Text()
           name := selection.Find(".wprm-recipe-ingredient-name").Text()
           if amount == "" { // probably an optional ingredient
			   fmt.Println("Finish check for optional ingredient")
           }
           fmt.Printf("Ingredient: %s\t Amount: %s\t Unit: %s\n", name, amount, unit)
           ingredientList = append(ingredientList, &RecipeIngredient{Name: name, Amount: amount, Unit: unit})
	});
	return ingredientList
}

func encodeIngredientList(list []*RecipeIngredient) {

	b, err := json.Marshal(list)

	if err != nil {
		log.Printf("Error: %s\n", err)
	}
	//fmt.Printf("Type:%T \t V:%v\n\n\n", string(b), string(b))
	fmt.Println("{")
	fmt.Println(string(b))
	fmt.Println("}")
}

func encodeInstructionList(list []*RecipeInstruction) {
	b, err := json.Marshal(list)

	if err != nil {
		log.Printf("Error: %s\n", err)
	}

	fmt.Println("{")
	fmt.Println(string(b))
	fmt.Println("}")
}


type RecipeInstruction struct {
	Instruction string	`json:"instruction"`
	Number		int		`json:"number"`
}

func scrapeInstructions(url string) []*RecipeInstruction {
	var instructionList []*RecipeInstruction
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Error Code: %d \t Status:%s\n", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".wprm-recipe-instruction-group .wprm-recipe-instructions .wprm-recipe-instruction").Each(
		func(i int, selection *goquery.Selection) {
			instruction := selection.Find(".wprm-recipe-instruction-text").Text()
			instructionList = append(instructionList, &RecipeInstruction{Instruction: instruction, Number: i})
	});
	return instructionList
}

type RecipeNutrition struct {
	Name	string	`json:"name"`
	Amount	string	`json:"amount"`
	Unit	string	`json:"unit"`
}

func encodeNutritionList(list []*RecipeNutrition) {
	b, err := json.Marshal(list)

	if err != nil {
		log.Printf("Error: %s\n", err)
	}

	fmt.Println("{")
	fmt.Println(string(b))
	fmt.Println("}")
}

func scrapeNutrition(url string) []*RecipeNutrition {
	var nutritionList []*RecipeNutrition
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Error Code: %d \t Status:%s\n", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".wprm-nutrition-label-text-nutrition-container").Each(func(i int, selection *goquery.Selection) {
		amount := selection.Find(".wprm-nutrition-label-text-nutrition-value").Text()
		unit   := selection.Find(".wprm-nutrition-label-text-nutrition-unit").Text()
		name   := selection.Find(".wprm-nutrition-label-text-nutrition-label").Text()

		fmt.Printf("Nutrition Label: %s\t Amount: %s\t Unit: %s\n", name, amount, unit)
		nutritionList = append(nutritionList, &RecipeNutrition{Name: name, Amount: amount, Unit: unit})
	});
	return nutritionList
}

func main() {
	args := os.Args[1:]
	if args == nil || args[0] == "" {
		return
	}
	ingredientList := scrapeIngredients(args[0])
	for _, e := range ingredientList {
		fmt.Printf("%+v\n", e)
	}
	fmt.Println("--------")
	encodeIngredientList(ingredientList)

	fmt.Println("--------")
	instructionList := scrapeInstructions(args[0])
	for _, e := range instructionList {
		fmt.Printf("%+v\n", e)
	}
	encodeInstructionList(instructionList)

	nutritionList := scrapeNutrition(args[0])
	for _, e := range nutritionList {
		fmt.Printf("%+v\n", e)
	}
	encodeNutritionList(nutritionList)
}

