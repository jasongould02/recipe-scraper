package main

import (
	"fmt"
	"strings"
	"log"
	"net/http"
	"net/http/httputil"
	"encoding/json"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
)

type RecipeIngredient struct {
	Name	string	`json:"name"`
	Amount	string	`json:"amount"`
	Unit	string	`json:"unit"`
}

type RecipeInstruction struct {
	Instruction string	`json:"instruction"`
	Number		int		`json:"number"`
}

type RecipeNutrition struct {
	Name	string	`json:"name"`
	Amount	string	`json:"amount"`
	Unit	string	`json:"unit"`
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
	Title			string  `json:"title"`
}

func scrapeIngredients(url string) []*RecipeIngredient {
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

	var ingredientList []*RecipeIngredient
	doc.Find(".wprm-recipe-ingredient-group .wprm-recipe-ingredient").Each(
		func(i int, selection *goquery.Selection) {
		   amount := selection.Find(".wprm-recipe-ingredient-amount").Text()
           unit := selection.Find(".wprm-recipe-ingredient-unit").Text()
           name := selection.Find(".wprm-recipe-ingredient-name").Text()
           if amount == "" { // probably an optional ingredient
			   log.Println("Ingredient found but has no amount, finish check for optional ingredients")
           }
		   log.Printf("Ingredient:%s\tAmount:%s\tUnit:%s\n", name, amount, unit)
           ingredientList = append(ingredientList, &RecipeIngredient{Name: name, Amount: amount, Unit: unit})
	});
	return ingredientList
}

func scrapeInstructions(url string) []*RecipeInstruction {
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

	var instructionList []*RecipeInstruction
	doc.Find(".wprm-recipe-instruction-group .wprm-recipe-instructions .wprm-recipe-instruction").Each(
		func(i int, selection *goquery.Selection) {
			instruction := selection.Find(".wprm-recipe-instruction-text").Text()
			instructionList = append(instructionList, &RecipeInstruction{Instruction: instruction, Number: i})
	});
	return instructionList
}

func scrapeNutrition(url string) []*RecipeNutrition {
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

	var nutritionList []*RecipeNutrition
	doc.Find(".wprm-nutrition-label-text-nutrition-container").Each(func(i int, selection *goquery.Selection) {
		amount := selection.Find(".wprm-nutrition-label-text-nutrition-value").Text()
		unit   := selection.Find(".wprm-nutrition-label-text-nutrition-unit").Text()
		name   := selection.Find(".wprm-nutrition-label-text-nutrition-label").Text()

		log.Printf("Nutrition Label: %s\tAmount:%s\tUnit:%s\n", name, amount, unit)
		nutritionList = append(nutritionList, &RecipeNutrition{Name: name, Amount: amount, Unit: unit})
	});
	return nutritionList
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
	cuisine := doc.Find(".wprm-recipe-cuisine").Text()
	course := doc.Find(".wprm-recipe-course").Text()
	author := doc.Find(".wprm-recipe-author").Text()
	log.Printf("Servings:%s %s\tCuisine:%s\tCourse:%s\tAuthor:%s\n", servings, servingsUnit, cuisine, course, author)

	prepTime := doc.Find(".wprm-recipe-prep_time").Text()
	prepTimeUnit := doc.Find(".wprm-recipe-prep_time-unit").Text()

	cookTime := doc.Find(".wprm-recipe-cook_time").Text()
	cookTimeUnit := doc.Find(".wprm-recipe-cook_time-unit").Text()

	totalTime := ""
	doc.Find(".wprm-recipe-total_time, .wprm-recipe-total_time-unit").Each(func(i int, selection *goquery.Selection) {
		totalTime += " " + selection.Text()
	})
	totalTime = strings.TrimSpace(totalTime)
	log.Printf("Prep Time:%s %s\tCook Time:%s %s\tTotal Time:%s\n", prepTime, prepTimeUnit, cookTime, cookTimeUnit, totalTime)

	summary := doc.Find(".wprm-recipe-summary").Text()
	log.Printf("Summary:%s\n", summary)

	title := doc.Find(".breadcrumb_last").Text()
	log.Printf("Title:%s\n", title)

	recipeMeta := RecipeMeta{Servings: servings, ServingsUnit: servingsUnit, Cuisine: cuisine, Course: course,
							 Author: author, PrepTime: prepTime, PrepTimeUnit: prepTimeUnit,
							 CookTime: cookTime, CookTimeUnit: cookTimeUnit, TotalTime: totalTime, Summary: summary, Title: title}

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
	return b
}

func (r *Recipe) ScrapeRecipe(url string) { // Puts all recipe data into
	log.Printf("Initiating scraping on URL:%s\n", url)
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

type GetRecipeRequest struct {
	Url string	`json:"URL"`
}

func RecipeHandler(w http.ResponseWriter, r *http.Request) {
	var recipeRequest GetRecipeRequest
	var data Recipe
	//var decodedString string

	fmt.Println("----------------------------")
	temp, _ := httputil.DumpRequest(r, true)
	fmt.Println(string(temp))
	fmt.Println("----------------------------")


	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	//err := decoder.Decode(&decodedString)
	err := decoder.Decode(&recipeRequest)
	fmt.Println("recipeRequest is now:", recipeRequest.Url)
	if err != nil {
		log.Printf("Error decoding JSON: %s\n", err)
	}
	/*fmt.Println("Remove this print statement")
	err = json.Unmarshal([]byte(decodedString), &recipeRequest)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %s\n", err)
	}*/

	log.Printf("Received recipe URL is: %s\n", recipeRequest.Url)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data.ScrapeRecipe(recipeRequest.Url)

	w.WriteHeader(http.StatusOK)
	w.Write(data.EncodeRecipe())
}

func main() {
	//args := os.Args[1:]
	//fmt.Println(args)
	//args := os.Args[1:]
	//if args == nil || args[0] == "" {
		//return
	//}

	fmt.Println("Starting Recipe-Scraper Server.")
	fmt.Println("Waiting for connection...")
	mux := mux.NewRouter()
	mux.HandleFunc("/new", RecipeHandler) // Receive new recipe URLs on this route

	err := http.ListenAndServe("localhost:4000", mux)
	//log.Fatal(err)
	if err != nil {
		log.Printf("Error: %s\n", err);
	}
}
