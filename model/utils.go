package model

import (
	"strings"
)

func ParseCocktailsMarkdown(markdown string) ([]Cocktail, error) {
	lines := strings.Split(markdown, "\n")
	var cocktails []Cocktail
	var currentCocktail *Cocktail

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if strings.HasPrefix(line, "# ") {
			if currentCocktail != nil {
				cocktails = append(cocktails, *currentCocktail)
			}
			currentCocktail = &Cocktail{Name: strings.TrimSpace(line[1:])}
		} else if strings.HasPrefix(line, "## Ingredients") {
			i++
			for i < len(lines) && strings.HasPrefix(lines[i], "*") {
				ingredientLine := strings.TrimSpace(lines[i][1:])
				ingredientParts := strings.SplitN(ingredientLine, " - ", 2)
				ingredient := Ingredient{}
				if len(ingredientParts) < 2 {
					ingredient = Ingredient{
						Name: strings.TrimSpace(ingredientParts[0]),
						Size: "",
					}
				} else {
					ingredient = Ingredient{
						Name: strings.TrimSpace(ingredientParts[0]),
						Size: strings.TrimSpace(ingredientParts[1]),
					}
				}
				currentCocktail.Ingredients = append(currentCocktail.Ingredients, ingredient)
				i++
			}
		} else if strings.HasPrefix(line, "## Steps") {
			i++
			for i < len(lines) && strings.HasPrefix(lines[i], "*") {
				step := strings.TrimSpace(lines[i][1:])
				currentCocktail.Steps = append(currentCocktail.Steps, step)
				i++
			}
		}
	}

	if currentCocktail != nil {
		cocktails = append(cocktails, *currentCocktail)
	}

	return cocktails, nil
}
