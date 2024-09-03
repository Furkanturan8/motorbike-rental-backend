package handlers

import (
	"github.com/gofiber/fiber/v2"
	"motorbike-rental-backend/internal/app/motorbike/models"
	bikeService "motorbike-rental-backend/internal/app/motorbike/services"
	viewmodel "motorbike-rental-backend/internal/app/motorbike/viewmodels"
	"motorbike-rental-backend/pkg/app"
	"strconv"
)

type MotorHandler struct {
	bikeService bikeService.IMotorService
}

func NewMotorHandler(s bikeService.IMotorService) MotorHandler {
	return MotorHandler{bikeService: s}
}

func (h MotorHandler) CreateMotor(ctx *app.Ctx) error {
	var bikeCreateVM viewmodel.BikeCreateVM
	if err := ctx.BodyParser(&bikeCreateVM); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz istek."})
	}

	motorbike := bikeCreateVM.ToDBModel()

	if err := h.bikeService.CreateMotorbike(ctx.Context(), &motorbike); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Motor oluşturulurken bir hata oluştu."})
	}

	photoModels := bikeCreateVM.ToPhotoModels(int(motorbike.ID))
	if err := h.bikeService.AddPhotosToMotorbike(ctx.Context(), photoModels); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Fotoğraflar eklenirken bir hata oluştu."})
	}

	return ctx.Status(fiber.StatusCreated).JSON(motorbike)
}

func (h MotorHandler) GetAllMotors(ctx *app.Ctx) error {
	// Motorları al
	motors, err := h.bikeService.GetAllMotors(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Motorlar getirilirken bir hata oluştu."})
	}

	var motorDetails []viewmodel.BikeDetailVM

	for _, motor := range *motors {
		// Her motorun fotoğraflarını al
		var photos []models.MotorbikePhoto
		err := h.bikeService.GetPhotosByID(ctx.Context(), strconv.FormatInt(motor.ID, 10), &photos)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Fotoğraflar getirilirken bir hata oluştu."})
		}

		// Motor detaylarını oluştur
		motorDetail := viewmodel.NewBikeDetailVM(motor, photos)
		motorDetails = append(motorDetails, motorDetail)
	}

	return ctx.Status(fiber.StatusOK).JSON(motorDetails)
}

func (h MotorHandler) GetPhotosByID(ctx *app.Ctx) error {
	motorbikeID := ctx.Params("id")
	var photos []models.MotorbikePhoto

	// fotoları motor id sine göre getiriyoruz
	if err := h.bikeService.GetPhotosByID(ctx.Context(), motorbikeID, &photos); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Fotoğraflar getirilirken bir hata oluştu."})
	}

	var photoDetailVMs []viewmodel.PhotoDetailVM
	for _, photo := range photos {
		photoDetailVMs = append(photoDetailVMs, viewmodel.PhotoDetailVM{
			ID:          int(photo.ID),
			MotorbikeID: photo.MotorbikeID,
			PhotoURL:    photo.PhotoURL,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(photoDetailVMs)
}

func (h MotorHandler) GetMotorByID(ctx *app.Ctx) error {
	param := ctx.Params("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz id isteği!"})
	}

	motor, err := h.bikeService.GetMotorByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Motor bulunamadı!"})
	}

	var photos []models.MotorbikePhoto

	err = h.bikeService.GetPhotosByID(ctx.Context(), strconv.Itoa(id), &photos)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Fotoğraflar getirilirken bir hata oluştu."})
	}

	motorDetail := viewmodel.NewBikeDetailVM(*motor, photos)

	return ctx.Status(fiber.StatusOK).JSON(motorDetail)
}
