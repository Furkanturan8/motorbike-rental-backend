package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

	if err := h.bikeService.CreateMotor(ctx.Context(), &motorbike); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Motor oluşturulurken bir hata oluştu."})
	}

	photoModels := bikeCreateVM.ToPhotoModels(int(motorbike.ID))
	if err := h.bikeService.AddPhotosToMotor(ctx.Context(), photoModels); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Fotoğraflar eklenirken bir hata oluştu."})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"info": "Motorsiklet Eklendi!"})
}

func (h MotorHandler) UpdateMotor(ctx *app.Ctx) error {
	var bikeUpdateVM viewmodel.BikeUpdateVM
	if err := ctx.BodyParser(&bikeUpdateVM); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz istek."})
	}

	param := ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz ID."})
	}

	motorbike, err := h.bikeService.GetMotorByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Motor bulunamadı."})
	}

	// Güncelleme verilerini motor modeline uygula
	updatedMotorbike := bikeUpdateVM.ToDBModel(*motorbike)
	if err := h.bikeService.UpdateMotor(ctx.Context(), &updatedMotorbike); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Motor güncellenirken bir hata oluştu."})
	}

	// Eğer fotoğraflar girilmişse güncelle (json da foto değerleri girilmemişse eskiler kalsın)
	if len(bikeUpdateVM.Photos) > 0 {
		photoModels := bikeUpdateVM.ToPhotoModels(id)
		if err := h.bikeService.UpdatePhotosForMotor(ctx.Context(), photoModels, id); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Fotoğraflar güncellenirken bir hata oluştu."})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"info": "Motorsiklet güncellendi!"})
}

func (h MotorHandler) DeleteMotor(ctx *app.Ctx) error {
	param := ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz ID!"})
	}

	err = h.bikeService.DeleteMotor(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Motor bulunamadı!"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Motor silinirken bir hata oluştu!"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"info": "Motor başarıyla silindi!"})
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

func (h MotorHandler) GetAvailableMotors(ctx *app.Ctx) error {
	motors, err := h.bikeService.GetMotorsForStatus(ctx.Context(), string(models.BikeAvailable))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Bir hata oluştu!"})
	}

	if len(*motors) == 0 {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"info": "Müsait motor yok!"})
	} else {
		var motorDetails []viewmodel.BikeDetailVM

		for _, motor := range *motors {
			// Her motorun fotoğraflarını al
			var photos []models.MotorbikePhoto
			err := h.bikeService.GetPhotosByID(ctx.Context(), strconv.FormatInt(motor.ID, 10), &photos)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Motorlar getirilirken bir hata oluştu."})
			}

			// Motor detaylarını oluştur
			motorDetail := viewmodel.NewBikeDetailVM(motor, photos)
			motorDetails = append(motorDetails, motorDetail)
		}

		return ctx.Status(fiber.StatusOK).JSON(motorDetails)
	}
}

func (h MotorHandler) GetMaintenanceMotors(ctx *app.Ctx) error {
	motors, err := h.bikeService.GetMotorsForStatus(ctx.Context(), string(models.BikeInMaintenance))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Bir hata oluştu!"})
	}

	if len(*motors) == 0 {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"info": "bakımda motor yok!"})
	} else {
		var motorDetails []viewmodel.BikeDetailVM

		for _, motor := range *motors {
			// Her motorun fotoğraflarını al
			var photos []models.MotorbikePhoto
			err := h.bikeService.GetPhotosByID(ctx.Context(), strconv.FormatInt(motor.ID, 10), &photos)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Motorlar getirilirken bir hata oluştu."})
			}

			// Motor detaylarını oluştur
			motorDetail := viewmodel.NewBikeDetailVM(motor, photos)
			motorDetails = append(motorDetails, motorDetail)
		}

		return ctx.Status(fiber.StatusOK).JSON(motorDetails)
	}
}

func (h MotorHandler) GetRentedMotors(ctx *app.Ctx) error {
	motors, err := h.bikeService.GetMotorsForStatus(ctx.Context(), string(models.BikeRented))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Bir hata oluştu!"})
	}

	if len(*motors) == 0 {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"info": "kiralanmış motor yok!"})
	} else {
		var motorDetails []viewmodel.BikeDetailVM

		for _, motor := range *motors {
			// Her motorun fotoğraflarını al
			var photos []models.MotorbikePhoto
			err := h.bikeService.GetPhotosByID(ctx.Context(), strconv.FormatInt(motor.ID, 10), &photos)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Motorlar getirilirken bir hata oluştu."})
			}

			// Motor detaylarını oluştur
			motorDetail := viewmodel.NewBikeDetailVM(motor, photos)
			motorDetails = append(motorDetails, motorDetail)
		}

		return ctx.Status(fiber.StatusOK).JSON(motorDetails)
	}
}
