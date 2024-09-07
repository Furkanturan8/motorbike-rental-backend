package handlers

import (
	connService "motorbike-rental-backend/internal/app/bluetooth-connection/services"
	"motorbike-rental-backend/internal/app/bluetooth-connection/viewmodels"
	"motorbike-rental-backend/pkg/app"
	"motorbike-rental-backend/pkg/errorsx"
)

type ConnHandler struct {
	connService connService.IConnService
}

func NewConnHandler(s connService.IConnService) ConnHandler {
	return ConnHandler{connService: s}
}

func (h ConnHandler) GetAllConnections(ctx *app.Ctx) error {
	connections, err := h.connService.GetAllConnections(ctx.Context())
	if err != nil {
		return errorsx.InternalError(err, "Bağlantılar getirilemedi!")
	}

	var connDetails []viewmodels.BluetoothConnectionDetailVM

	for _, conn := range *connections {
		vm := viewmodels.BluetoothConnectionDetailVM{}
		connDetail := vm.ToViewModel(conn)
		connDetails = append(connDetails, connDetail)
	}

	return ctx.SuccessResponse(connDetails, len(connDetails))
}
