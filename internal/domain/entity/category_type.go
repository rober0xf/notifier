package entity

type CategoryType string

const (
	CategoryTypeElectronics   CategoryType = "electronics"
	CategoryTypeEntertainment CategoryType = "entertainment"
	CategoryTypeEducation     CategoryType = "education"
	CategoryTypeClothing      CategoryType = "clothing"
	CategoryTypeWork          CategoryType = "work"
	CategoryTypeSports        CategoryType = "sports"
)

func (c CategoryType) IsValid() bool {
	switch c {
	case CategoryTypeElectronics, CategoryTypeEntertainment, CategoryTypeEducation, CategoryTypeClothing, CategoryTypeWork, CategoryTypeSports:
		return true
	}
	return false
}
