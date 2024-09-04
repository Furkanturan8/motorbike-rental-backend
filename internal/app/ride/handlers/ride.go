package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	rideService "motorbike-rental-backend/internal/app/ride/services"
	"motorbike-rental-backend/internal/app/ride/viewmodels"
	"motorbike-rental-backend/pkg/app"
	"strconv"
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

func (h RideHandler) GetRideByID(ctx *app.Ctx) error {
	param := ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Hatalı istek!"})
	}

	data, err := h.rideService.GetRideByID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Sürüş bulunamadı!"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Sürüş detayları getirilirken hata oluştu!"})
	}

	var vm viewmodels.RideDetailVM
	rideDetail := vm.ToDBModel(*data)

	return ctx.Status(fiber.StatusOK).JSON(rideDetail)
}

func (h RideHandler) CreateRide(ctx *app.Ctx) error {
	var rideCreateVM viewmodels.RideCreateVM

	if err := ctx.BodyParser(&rideCreateVM); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz istek!"})
	}

	ride := rideCreateVM.ToDBModel()

	if err := h.rideService.CreateRide(ctx.Context(), &ride); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Sürüş oluşturulurken hata oluştu!"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"info": "Sürüş eklendi!"})
}
