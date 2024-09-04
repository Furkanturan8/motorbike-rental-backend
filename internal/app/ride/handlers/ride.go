package handlers

import (
	"github.com/gofiber/fiber/v2"
	rideService "motorbike-rental-backend/internal/app/ride/services"
	"motorbike-rental-backend/internal/app/ride/viewmodels"
	"motorbike-rental-backend/pkg/app"
)

type RideHandler struct {
	rideService rideService.IRideService
}

func NewRideHandler(s rideService.IRideService) RideHandler {
	return RideHandler{rideService: s}
}

func (h RideHandler) GetAllRides(ctx *app.Ctx) error {
	rides, err := h.rideService.GetAllRides(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Sürüşler getirilirken hata oluştu!"})
	}

	var rideDetails []viewmodels.RideDetailVM

	for _, ride := range *rides {
		vm := viewmodels.RideDetailVM{}
		rideDetail := vm.ToDBModel(ride)
		rideDetails = append(rideDetails, rideDetail)
	}

	return ctx.Status(fiber.StatusOK).JSON(rideDetails)
}
