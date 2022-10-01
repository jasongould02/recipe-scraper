package main

import (
	"os"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type RecipeIngredient struct {
	name string;
	amount string;
	unit string;
}


// only for recipes on noracooks.com for now
func ingredientScraper(url string) []*RecipeIngredient {
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
	// Place ingredients into a JSON file to be sent to db
		func(i int, selection *goquery.Selection) {
			amount := selection.Find(".wprm-recipe-ingredient-amount").Text()
			unit := selection.Find(".wprm-recipe-ingredient-unit").Text()
			name := selection.Find(".wprm-recipe-ingredient-name").Text()

			if amount == "" { // probably an optional ingredient
				fmt.Println("Finish check for optional ingredient")
			}

			fmt.Printf("Ingredient: %s\t Amount: %s\t Unit: %s\n", name, amount, unit)
			ingredientList = append(ingredientList, &RecipeIngredient{name: name, amount: amount, unit: unit})
		});

	return ingredientList
}

/*doc.Find("script").Each(func(i int, s *goquery.Selection) {
		recipe_198 := s.Find("var wprmpuc_recipe_198").Text()
		fmt.Println("script:", recipe_198)
	})
}*/


func main() {
	args := os.Args[1:]
	if args == nil || args[0] == "" {
		return
	}
	list := ingredientScraper(args[0])
	for _, e := range list {
		fmt.Printf("%+v\n", e)
	}
}

