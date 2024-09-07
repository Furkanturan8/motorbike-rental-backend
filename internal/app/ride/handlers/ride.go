package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"motorbike-rental-backend/internal/app/ride/models"
	rideService "motorbike-rental-backend/internal/app/ride/services"
	"motorbike-rental-backend/internal/app/ride/viewmodels"
	"motorbike-rental-backend/pkg/app"
	"motorbike-rental-backend/pkg/errorsx"
	"strconv"
	"time"
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
		rideDetail := vm.ToViewModel(ride)
		rideDetails = append(rideDetails, rideDetail)
	}

	return ctx.SuccessResponse(rideDetails, len(rideDetails))
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
			return errorsx.NotFoundError("Sürüş bulunamadı!")
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Sürüş detayları getirilirken hata oluştu!"})
	}

	var vm viewmodels.RideDetailVM
	rideDetail := vm.ToViewModel(*data)

	return ctx.SuccessResponse(rideDetail, 1)
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

func (h RideHandler) GetRidesByUserID(ctx *app.Ctx) error {
	param := ctx.Params("userID")
	id, err := strconv.Atoi(param)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Hatalı istek!"})
	}

	rides, err := h.rideService.GetRidesByUserID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Bu kullanıcı bulunamadı!"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Sürüş detayları getirilirken hata oluştu!"})
	}

	var rideDetails []viewmodels.RideDetailVM

	for _, ride := range *rides {
		vm := viewmodels.RideDetailVM{}
		rideDetail := vm.ToViewModel(ride)
		rideDetails = append(rideDetails, rideDetail)
	}

	return ctx.SuccessResponse(rideDetails, len(rideDetails))
}

func (h RideHandler) GetRideByUserID(ctx *app.Ctx) error {
	param1 := ctx.Params("userID")
	userID, err := strconv.Atoi(param1)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Hatalı istek!"})
	}

	param2 := ctx.Params("rideID")
	rideID, err := strconv.Atoi(param2)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Hatalı istek!"})
	}

	ride, err := h.rideService.GetRideByUserID(ctx.Context(), userID, rideID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Bu kullanıcının ilgili sürüşü bulunamadı!"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Sürüş detayları getirilirken hata oluştu!"})
	}

	var vm viewmodels.RideDetailVM
	rideDetail := vm.ToViewModel(*ride)

	return ctx.SuccessResponse(rideDetail, 1) // pkg/errorsx.go sınıfından
}

func (h RideHandler) GetRidesByBikeID(ctx *app.Ctx) error {
	param := ctx.Params("bikeID")
	bikeID, err := strconv.Atoi(param)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Hatalı istek!"})
	}

	rides, err := h.rideService.GetRidesByBikeID(ctx.Context(), bikeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Bu motor yok!"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Bu motora ait sürüş getirilirken hata oluştu!"})
	}

	var rideDetails []viewmodels.RideDetailVM

	for _, ride := range *rides {
		vm := viewmodels.RideDetailVM{}
		rideDetail := vm.ToViewModel(ride)
		rideDetails = append(rideDetails, rideDetail)
	}

	return ctx.SuccessResponse(rideDetails, len(rideDetails))
}

func (h RideHandler) UpdateRideByID(ctx *app.Ctx) error {
	var rideUpdateVM viewmodels.RideUpdateVM
	if err := ctx.BodyParser(&rideUpdateVM); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz istek!"})
	}

	param := ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Hatalı istek!"})
	}

	ride, err := h.rideService.GetRideByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Sürüş bulunamadı!"})
	}

	updatedRide := rideUpdateVM.ToDBModel(*ride)
	if err := h.rideService.UpdateRide(ctx.Context(), &updatedRide); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Sürüş güncellenirken bir hata oluştu!"})
	}

	/*var vm viewmodels.RideDetailVM
	rideDetail := vm.ToDBModel(updatedRide)*/ // eğer güncellediğimiz veriyi listelemek istersek rideDetail i gönder!

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"info": "Sürüş başarıyla güncellendi!"})
}

func (h RideHandler) DeleteRide(ctx *app.Ctx) error {
	param := ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Hatalı istek!"})
	}

	err = h.rideService.DeleteRide(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Sürüş bulunamadı!"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Sürüş silinirken hata oluştu!"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"info": "Sürüş başarıyla silindi!"})
}

// (adminler için) belirli tarih aralıklarındaki sürüşleri getirir -> /filtered-rides?start_time=2024-09-04&end_time=2024-09-05
func (h RideHandler) GetRidesByDateRange(ctx *app.Ctx) error {
	startTimeStr := ctx.Query("start_time")
	endTimeStr := ctx.Query("end_time")

	if startTimeStr == "" || endTimeStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Lütfen geçerli bir start_time ve end_time parametresi girin!",
		})
	}

	// Tarihleri parse et
	startTime, err := time.Parse("2006-01-02", startTimeStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz start_time formatı!"})
	}

	endTime, err := time.Parse("2006-01-02", endTimeStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz end_time formatı!"})
	}

	rides, err := h.rideService.GetRidesByDateRange(ctx.Context(), startTime, endTime)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Sürüşler listelenirken bir hata oluştu!"})
	}

	var ridesDetails []viewmodels.RideDetailVM

	for _, ride := range *rides {
		vm := viewmodels.RideDetailVM{}
		rideDetail := vm.ToViewModel(ride)
		ridesDetails = append(ridesDetails, rideDetail)
	}

	return ctx.SuccessResponse(ridesDetails, len(ridesDetails))
}

// (kullanıcılar için) userID ye göre belirli tarihler arasında getirir -> /rides/user/:userID/filter?start_time=2024-09-01&end_time=2024-09-09
func (h RideHandler) GetRidesByUserAndDate(ctx *app.Ctx) error {
	param := ctx.Params("userID")
	id, err := strconv.Atoi(param)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Hatalı istek!"})
	}

	// start_time ve end_time parametrelerini al
	startTimeStr := ctx.Query("start_time")
	endTimeStr := ctx.Query("end_time")

	// Zaman formatını kontrol et
	startTime, err := time.Parse("2006-01-02", startTimeStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz start_time formatı!"})
	}

	endTime, err := time.Parse("2006-01-02", endTimeStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz end_time formatı!"})
	}

	// Önce mevcut fonksiyon ile kullanıcıya ait tüm sürüşleri getir (yukarıdaki func kullandık)
	rides, err := h.rideService.GetRidesByUserID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Bu kullanıcıya ait sürüş bulunamadı!"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Sürüş detayları getirilirken hata oluştu!"})
	}

	// Tarih aralığına göre filtreleme yapar
	var filteredRides []models.Ride
	for _, ride := range *rides {
		if ride.StartTime.After(startTime) && ride.EndTime.Before(endTime) {
			filteredRides = append(filteredRides, ride)
		}
	}

	// Sonuçları ViewModel'e dönüştür
	var rideDetails []viewmodels.RideDetailVM
	for _, ride := range filteredRides {
		vm := viewmodels.RideDetailVM{}
		rideDetail := vm.ToViewModel(ride)
		rideDetails = append(rideDetails, rideDetail)
	}

	return ctx.SuccessResponse(rideDetails, len(rideDetails))
}
