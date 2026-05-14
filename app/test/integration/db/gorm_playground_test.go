package db

import (
	"example-wikipedia-scraper/internal/db"
	"example-wikipedia-scraper/test/integration"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFuelTypeEnum(t *testing.T) {
	_, err := integration.InitTest()
	var fuelTypes []string
	assert.NoError(t, err, "Failed to initialize test environment")
	db.DB.Raw("SELECT unnest(enum_range(NULL::fuel_type))").Scan(&fuelTypes)
	assert.NotEmpty(t, fuelTypes, "Fuel types should not be empty")
}
