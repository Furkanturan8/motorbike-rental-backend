package handlers

import (
	"github.com/gofiber/fiber/v2"
	mapService "motorbike-rental-backend/internal/app/map/services"
	"motorbike-rental-backend/internal/app/map/viewmodels"
	"motorbike-rental-backend/pkg/app"
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
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Motorsiklet Konumları getirilemedi!"})
	}

	var mapDetails []viewmodels.MapDetailVM

	for _, _map := range *maps {
		vm := viewmodels.MapDetailVM{}
		mapDetail := vm.ToDBModel(_map)
		mapDetails = append(mapDetails, mapDetail)
	}

	return ctx.SuccessResponse(mapDetails, len(mapDetails))
}
