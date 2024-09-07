package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	connService "motorbike-rental-backend/internal/app/bluetooth-connection/services"
	"motorbike-rental-backend/internal/app/bluetooth-connection/viewmodels"
	motorService "motorbike-rental-backend/internal/app/motorbike/services"
	"motorbike-rental-backend/pkg/app"
	"motorbike-rental-backend/pkg/errorsx"
	"strconv"
)

type ConnHandler struct {
	connService  connService.IConnService
	motorService motorService.IMotorService
}

func NewConnHandler(s connService.IConnService, m motorService.IMotorService) ConnHandler {
	return ConnHandler{connService: s, motorService: m}
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

func (h ConnHandler) GetConnByID(ctx *app.Ctx) error {
	param := ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return errorsx.BadRequestError("Hatalı istek!")
	}

	data, err := h.connService.GetConnByID(ctx.Context(), id)
	if err != nil {
		if errorsx.Is(err, gorm.ErrRecordNotFound) {
			return errorsx.NotFoundError("Böyle bir bağlantı yok!")
		}
		return errorsx.InternalError(err, "Bir hata oluştu!")
	}

	var vm viewmodels.BluetoothConnectionDetailVM
	connDetail := vm.ToViewModel(*data)

	return ctx.SuccessResponse(connDetail, 1)
}

// important logical thing : when i create a new connection, i need to check motorbikeID.
// Because maybe, the motor are using by somebody. so we'll be careful to this when we add it: the status of the motorbike must be available. not rented or maintained
func (h ConnHandler) CreateConn(ctx *app.Ctx) error {
	var connVM viewmodels.BluetoothConnectionCreateVM
	if err := ctx.BodyParser(&connVM); err != nil {
		return errorsx.BadRequestError("Geçersiz istek!")
	}

	connection := connVM.ToDBModel()

	motor, err := h.motorService.GetMotorByID(ctx.Context(), int(connVM.MotorbikeID))
	if err != nil {
		if errorsx.Is(err, gorm.ErrRecordNotFound) {
			return errorsx.InternalError(err, "Böyle bir motorsiklet yok! Hatalı bağlantı isteği!")
		}
		return errorsx.InternalError(err, "Bir hata oluştu!")
	}

	// Motorbike'ın durumu 'Available' mı kontrol et
	if motor.Status != "Available" {
		return errorsx.BadRequestError("Bu Motorbisiklet şu anda müsait değil!")
	}

	if err = h.connService.CreateConn(ctx.Context(), &connection); err != nil {
		return errorsx.InternalError(err, "Bağlantı kurulurken hata oluştu!")
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"info": "Bağlantı kuruldu!"})
}
