package controller

import (
	"git.extrasys.it/aldo.daquino/visitors-backend-operator/pkg/controller/visitorsbackend"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, visitorsbackend.Add)
}
