package entity

import "slices"

type CategoryType string

const (
	CategoryTypeElectronics   CategoryType = "electronics"
	CategoryTypeEntertainment CategoryType = "entertainment"
	CategoryTypeEducation     CategoryType = "education"
	CategoryTypeClothing      CategoryType = "clothing"
	CategoryTypeWork          CategoryType = "work"
	CategoryTypeSports        CategoryType = "sports"
)

var AllCategoryTypes = []CategoryType{
	CategoryTypeElectronics,
	CategoryTypeEntertainment,
	CategoryTypeEducation,
	CategoryTypeClothing,
	CategoryTypeWork,
	CategoryTypeSports,
}

func (c CategoryType) IsValid() bool {
	return slices.Contains(AllCategoryTypes, c)
}
