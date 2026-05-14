package scraper

import (
	"example-wikipedia-scraper/internal/model"
	"reflect"
)

// PageComparator porównuje zmiany w stronach (DRY - jedna odpowiedzialność)
type PageComparator struct{}

func NewPageComparator() *PageComparator {
	return &PageComparator{}
}

func (pc *PageComparator) HasChanged(current, updated *model.Page) bool {
	if current == nil || updated == nil {
		return true
	}

	return pc.compareFields(current, updated)
}

func (pc *PageComparator) compareFields(current, updated *model.Page) bool {
	fields := []struct {
		current, updated interface{}
	}{
		{current.Title, updated.Title},
		{current.TextField1, updated.TextField1},
		{current.TextField2, updated.TextField2},
		{current.TextField3, updated.TextField3},
	}

	for _, field := range fields {
		if !pc.valuesEqual(field.current, field.updated) {
			return true
		}
	}
	return false
}

// valuesEqual porównuje dwie wartości, obsługując slice'y
func (pc *PageComparator) valuesEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}
