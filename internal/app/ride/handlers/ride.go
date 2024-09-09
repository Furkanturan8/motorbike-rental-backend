package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	motorModel "motorbike-rental-backend/internal/app/motorbike/models"
	motorService "motorbike-rental-backend/internal/app/motorbike/services"
	"motorbike-rental-backend/internal/app/ride/models"
	rideService "motorbike-rental-backend/internal/app/ride/services"
	"motorbike-rental-backend/internal/app/ride/viewmodels"
	"motorbike-rental-backend/pkg/app"
	"motorbike-rental-backend/pkg/errorsx"
	"motorbike-rental-backend/pkg/utils"
	"strconv"
	"time"
)

type RideHandler struct {
	rideService  rideService.IRideService
	motorService motorService.IMotorService
}

func NewRideHandler(s rideService.IRideService, m motorService.IMotorService) RideHandler {
	return RideHandler{rideService: s, motorService: m}
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

	motor, err := h.motorService.GetMotorByID(ctx.Context(), int(ride.MotorbikeID))
	if err != nil {
		if errorsx.Is(err, gorm.ErrRecordNotFound) {
			return errorsx.InternalError(err, "Böyle bir motorsiklet yok! Hatalı bağlantı isteği!")
		}
		return errorsx.InternalError(err, "Bir hata oluştu!")
	}

	if motor == nil {
		return errorsx.InternalError(nil, "Motorbisiklet verisi alınamadı!")
	}

	// Motorbike'ın durumu 'Available' mı kontrol et
	if motor.Status != motorModel.BikeAvailable {
		return errorsx.BadRequestError("Bu Motorbisiklet şu anda müsait değil!")
	}

	motor.Status = motorModel.BikeRented

	if err = h.motorService.UpdateMotor(ctx.Context(), motor); err != nil {
		return errorsx.InternalError(err, "Sürüş oluşturulamadı!")
	}

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

func (h RideHandler) FinishRide(ctx *app.Ctx) error {
	id, err := utils.GetMyParamInt(ctx, "id")
	if err != nil {
		return errorsx.BadRequestError("Hatalı istek!")
	}

	ride, err := h.rideService.GetRideByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Sürüş bulunamadı!"})
	}

	if ride.EndTime == nil {
		now := time.Now().UTC()
		ride.EndTime = &now

		// Sürüş süresini hesapla// Sürüş süresini hesapla (StartTime bir pointer değilse)
		duration := now.Sub(ride.StartTime)
		ride.Duration = strconv.Itoa(int(duration.Seconds())) // Saniye cinsinden süreyi kaydet

		// Süreyi dakika cinsinden hesapla (her dakika için 3 TL)
		minutes := int(duration.Minutes())
		costPerMinute := 3
		ride.Cost = float64(minutes*costPerMinute) + 10
	} else {
		return errorsx.BadRequestError("Zaten sürüş bitirildi!")
	}

	err = h.rideService.UpdateRide(ctx.Context(), ride)
	if err != nil {
		return errorsx.InternalError(err, "Sürüş bitirilemedi!")
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"info": "Sürüş bitirildi!", "cost (TL)": ride.Cost})
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
