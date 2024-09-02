package handlers

import (
	"motorbike-rental-backend/internal/app/motorbike/models"
	bikeService "motorbike-rental-backend/internal/app/motorbike/services"
	viewmodel "motorbike-rental-backend/internal/app/motorbike/viewmodels"
	"motorbike-rental-backend/pkg/app"
	"motorbike-rental-backend/pkg/errorsx"
)

type MotorHandler struct {
	bikeService bikeService.IMotorService
}

func NewMotorHandler(s bikeService.IMotorService) MotorHandler {
	return MotorHandler{bikeService: s}
}

func (h MotorHandler) CreateMotor(ctx *app.Ctx) error {
	var vm viewmodel.BikeCreateVM
	if err := ctx.BodyParseValidate(&vm); err != nil {
		return errorsx.ValidationError(err)
	}

	bike := vm.ToDBModel(models.Motorbike{})

	err := h.bikeService.CreateMotorbike(&bike)
	if err != nil {
		return err
	}

	return nil
}
