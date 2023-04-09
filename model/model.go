package model

type Cocktail struct {
	Name        string
	Ingredients []Ingredient
	Steps       []string
}

type Ingredient struct {
	Name string
	Size string
}
