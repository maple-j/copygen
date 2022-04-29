// Code generated by github.com/switchupcb/copygen
// DO NOT EDIT.

// Package domain contains the setup information for copygen generated code.
package domain

import (
	c "strconv"

	"github.com/switchupcb/copygen/examples/_tests/import/models"
)

/* Define the function and field this converter is applied to using regex. */
// Itoa converts an integer to an ascii value.
func Itoa(i int) string {
	return c.Itoa(i)
}

// ModelsToDomain copies a *models.Account, *models.User to a *Account.
func ModelsToDomain(tA *Account, fA *models.Account, fU *models.User) {
	// *Account fields
	tA.ID = fA.ID
	tA.UserID = Itoa(fU.UserID)
	tA.Name = fA.Name
}
