package handlers

import (
	"errors"
	"gorm.io/gorm"
	mapService "motorbike-rental-backend/internal/app/map/services"
	"motorbike-rental-backend/internal/app/map/viewmodels"
	"motorbike-rental-backend/pkg/app"
	"motorbike-rental-backend/pkg/errorsx"
	"strconv"
)

type MapHandler struct {
	mapService mapService.IMapService
}

func NewMapHandler(s mapService.IMapService) MapHandler {
	return MapHandler{mapService: s}
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
