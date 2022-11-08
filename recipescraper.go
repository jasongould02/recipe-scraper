package main

import (
	"os"
	"fmt"
	"strings"
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

type RecipeMeta struct {
	Servings		string	`json:"servings"`
	ServingsUnit	string	`json:"servingsUnit"`
	Cuisine			string	`json:"cuisine"`
	Course			string	`json:"course"`
	Author			string	`json:"author"`
	PrepTime		string	`json:"prepTime"`
	PrepTimeUnit	string	`json:"prepTimeUnit"`
	CookTime		string	`json:"cookTime"`
	CookTimeUnit	string	`json:"cookTimeUnit"`
	TotalTime		string	`json:"totalTime"`
	Summary			string	`json:"summary"`
}

func scrapeMeta(url string) RecipeMeta {

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

	servings := doc.Find(".wprm-recipe-servings").Text()
	servingsUnit := doc.Find(".wprm-recipe-servings-unit").Text()
	fmt.Println("servings: ", servings, "\tservings unit: ", servingsUnit)

	cuisine := doc.Find(".wprm-recipe-cuisine").Text()
	fmt.Println("cuisine: ", cuisine)

	course := doc.Find(".wprm-recipe-course").Text()
	fmt.Println("course: ", course)

	author := doc.Find(".wprm-recipe-author").Text()
	fmt.Println("author: ", author)

	prepTime := doc.Find(".wprm-recipe-prep_time").Text()
	prepTimeUnit := doc.Find(".wprm-recipe-prep_time-unit").Text()

	cookTime := doc.Find(".wprm-recipe-cook_time").Text()
	cookTimeUnit := doc.Find(".wprm-recipe-cook_time-unit").Text()

	fmt.Printf("PrepTime:%s %s\tCookTime:%s %s\n", prepTime, prepTimeUnit, cookTime, cookTimeUnit)

	totalTime := ""
	doc.Find(".wprm-recipe-total_time, .wprm-recipe-total_time-unit").Each(func(i int, selection *goquery.Selection) {
		totalTime += " " + selection.Text()
	})
	totalTime = strings.TrimSpace(totalTime)

	summary := doc.Find(".wprm-recipe-summary").Text()
	fmt.Println("summary:", summary)

	recipeMeta := RecipeMeta{Servings: servings, ServingsUnit: servingsUnit, Cuisine: cuisine, Course: course,
							 Author: author, PrepTime: prepTime, PrepTimeUnit: prepTimeUnit,
							 CookTime: cookTime, CookTimeUnit: cookTimeUnit, TotalTime: totalTime, Summary: summary}

	return recipeMeta
}

type Recipe struct {
	Metadata RecipeMeta						`json:"Metadata"`
	NutritionList []*RecipeNutrition		`json:"Nutrition"`
	InstructionList []*RecipeInstruction	`json:"Instruction"`
	IngredientList []*RecipeIngredient		`json:"Ingredient"`
}

func (r *Recipe) EncodeRecipe() []byte {
	b, err := json.Marshal(r)
	if err != nil {
		log.Printf("Error: %s\t", err)
	}
	fmt.Println("\n\n\n\n\n\n")
	fmt.Println(string(b))
	fmt.Println("\n\n\n\n\n\n")
	return b
}

func (r *Recipe) ScrapeRecipe(url string) { // Puts all recipe data into
	r.IngredientList = scrapeIngredients(url)
	r.InstructionList = scrapeInstructions(url)
	r.NutritionList = scrapeNutrition(url)
	r.Metadata = scrapeMeta(url)
}

func (r Recipe) encodeMeta(meta RecipeMeta) []byte {
	b, err := json.Marshal(meta)
	if err != nil {
		log.Printf("Error: %s\n", err)
	}
	return b
}

func (r Recipe) encodeNutritionList(list []*RecipeNutrition) []byte {
	b, err := json.Marshal(list)
	if err != nil {
		log.Printf("Error: %s\n", err)
	}
	return b
}

func (r Recipe) encodeInstructionList(list []*RecipeInstruction) []byte {
	b, err := json.Marshal(list)
	if err != nil {
		log.Printf("Error: %s\n", err)
	}
	return b
}

func (r Recipe) encodeIngredientList(list []*RecipeIngredient) []byte {
	b, err := json.Marshal(list)
	if err != nil {
		log.Printf("Error: %s\n", err)
	}
	return b
}

func main() {
	args := os.Args[1:]
	if args == nil || args[0] == "" {
		return
	}
	var data Recipe

	data.ScrapeRecipe(args[0])
	data.EncodeRecipe()
}
