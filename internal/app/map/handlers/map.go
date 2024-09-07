package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	mapService "motorbike-rental-backend/internal/app/map/services"
	"motorbike-rental-backend/internal/app/map/viewmodels"
	bikeService "motorbike-rental-backend/internal/app/motorbike/services"
	"motorbike-rental-backend/pkg/app"
	"motorbike-rental-backend/pkg/errorsx"
	"strconv"
)

type MapHandler struct {
	mapService  mapService.IMapService
	bikeService bikeService.IMotorService
}

func NewMapHandler(s mapService.IMapService, m bikeService.IMotorService) MapHandler {
	return MapHandler{mapService: s, bikeService: m}
}

func (h MapHandler) GetAllMaps(ctx *app.Ctx) error {
	maps, err := h.mapService.GetAllMaps(ctx.Context())
	if err != nil {
		return errorsx.InternalError(err, "motorsiklet konumları getirilemedi!")
	}

	var mapDetails []viewmodels.MapDetailVM

	for _, _map := range *maps {
		vm := viewmodels.MapDetailVM{}
		mapDetail := vm.ToDBModel(_map)
		mapDetails = append(mapDetails, mapDetail)
	}

	return ctx.SuccessResponse(mapDetails, len(mapDetails))
}

func (h MapHandler) GetMapByID(ctx *app.Ctx) error {
	param := ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return errorsx.BadRequestError("Hatalı istek!")
	}

	data, err := h.mapService.GetMapByID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorsx.NotFoundError("Bu adres bulunamadı!")
		}
		return errorsx.InternalError(err, "Bu adres getirilirken hata oluştu!")
	}

	var vm viewmodels.MapDetailVM

	mapDetail := vm.ToDBModel(*data)

	return ctx.SuccessResponse(mapDetail, 1)
}

func (h MapHandler) GetMapByMotorID(ctx *app.Ctx) error {
	param := ctx.Params("motorbikeID")
	id, err := strconv.Atoi(param)
	if err != nil {
		return errorsx.BadRequestError("Hatalı istek!")
	}

	_, err = h.bikeService.GetMotorByID(ctx.Context(), id)
	if err != nil {
		if errorsx.Is(err, gorm.ErrRecordNotFound) {
			return errorsx.NotFoundError("Böyle bir motorsiklet zaten yok! Var olan motoru seçiniz!")
		}
		return errorsx.InternalError(err, "Motorsiklet sorgulama sırasında bir hata oluştu!")
	}

	data, err := h.mapService.GetMapByMotorbikeID(ctx.Context(), id)
	if err != nil && !errorsx.Is(err, gorm.ErrRecordNotFound) {
		return errorsx.InternalError(err, "Adres kontrol edilirken hata oluştu!")
	}

	vm := viewmodels.MapDetailVM{}
	mapDetail := vm.ToDBModel(*data)

	return ctx.SuccessResponse(mapDetail, 1)
}

func (h MapHandler) CreateMap(ctx *app.Ctx) error {
	mapCreateVM := viewmodels.MapCreateVM{}
	if err := ctx.BodyParser(&mapCreateVM); err != nil {
		return errorsx.BadRequestError("Geçersiz istek!")
	}

	_map := mapCreateVM.ToDBModel()

	// Motorbike ID ile zaten kayıtlı bir adres var mı kontrol et
	existingMap, err := h.mapService.GetMapByMotorbikeID(ctx.Context(), int(_map.MotorbikeID))
	if err == nil && existingMap != nil {
		// Eğer motorbike'a ait bir harita zaten varsa, ekleme işlemi yapılmaz
		return errorsx.ConflictError("Bu motorsiklet için zaten bir adres kaydı bulunmakta!")
	} else if err != nil && !errorsx.Is(err, gorm.ErrRecordNotFound) {
		return errorsx.InternalError(err, "Adres kontrol edilirken hata oluştu!")
	}

	// var olmayan bir motor ile motor ve motorun yeri olan map ilişkisi hakkında önce motor var mı diye kontrol ediyoruz!
	_, err = h.bikeService.GetMotorByID(ctx.Context(), int(_map.MotorbikeID))
	if err != nil {
		if errorsx.Is(err, gorm.ErrRecordNotFound) {
			return errorsx.NotFoundError("Böyle bir motorsiklet zaten yok! Var olan motoru seçiniz, veya yeni bir motor oluşturun!")
		}
		return errorsx.InternalError(err, "Motorsiklet sorgulama sırasında bir hata oluştu!")
	}

	if err := h.mapService.CreateMap(ctx.Context(), &_map); err != nil {
		return errorsx.InternalError(err, "Adres oluşturulurken hata oluştu!")
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"info": "Adres oluşturuldu!"})
}

func (h MapHandler) DeleteMap(ctx *app.Ctx) error {
	param := ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return errorsx.BadRequestError("Hatalı istek!")
	}

	if err = h.mapService.DeleteMap(ctx.Context(), id); err != nil {
		if errorsx.Is(err, gorm.ErrRecordNotFound) {
			return errorsx.NotFoundError("Böyle bir adres zaten yok!")
		}
		return errorsx.InternalError(err, "Adres silinirken bir hata oluştu!")
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"info": "Adres başarılı bir şekilde silindi!"})
}

func (h MapHandler) UpdateMap(ctx *app.Ctx) error {
	var mapUpdateVM viewmodels.MapUpdateVM
	if err := ctx.BodyParser(&mapUpdateVM); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz istek."})
	}

	param := ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return errorsx.BadRequestError("Hatalı istek!")
	}

	data, err := h.mapService.GetMapByID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorsx.NotFoundError("Bu adres bulunamadı!")
		}
		return errorsx.InternalError(err, "Bu adres getirilirken hata oluştu!")
	}

	updatedMap := mapUpdateVM.ToDBModel(*data)
	if err = h.mapService.UpdateMap(ctx.Context(), &updatedMap); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Adres güncellenirken bir hata oluştu!"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"info": "Adres güncellendi!"})
}

func (h MapHandler) UpdateMapByMotorID(ctx *app.Ctx) error {
	var mapUpdateVM viewmodels.MapUpdateVM
	if err := ctx.BodyParser(&mapUpdateVM); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz istek."})
	}

	param := ctx.Params("motorbikeID")
	motorbikeID, err := strconv.Atoi(param)
	if err != nil {
		return errorsx.BadRequestError("Hatalı istek!")
	}

	_, err = h.bikeService.GetMotorByID(ctx.Context(), motorbikeID)
	if err != nil {
		if errorsx.Is(err, gorm.ErrRecordNotFound) {
			return errorsx.NotFoundError("Böyle bir motorsiklet zaten yok! Var olan motoru seçiniz!")
		}
		return errorsx.InternalError(err, "Motorsiklet sorgulama sırasında bir hata oluştu!")
	}

	data, err := h.mapService.GetMapByMotorbikeID(ctx.Context(), motorbikeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorsx.NotFoundError("Bu adres bulunamadı!")
		}
		return errorsx.InternalError(err, "Bu adres getirilirken hata oluştu!")
	}

	updatedMap := mapUpdateVM.ToDBModel(*data)
	if err = h.mapService.UpdateMap(ctx.Context(), &updatedMap); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Adres güncellenirken bir hata oluştu!"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"info": "Adres güncellendi!"})
}
