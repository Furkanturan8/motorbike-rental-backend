package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	connService "motorbike-rental-backend/internal/app/bluetooth-connection/services"
	"motorbike-rental-backend/internal/app/bluetooth-connection/viewmodels"
	motorModel "motorbike-rental-backend/internal/app/motorbike/models"
	motorService "motorbike-rental-backend/internal/app/motorbike/services"
	"motorbike-rental-backend/pkg/app"
	"motorbike-rental-backend/pkg/errorsx"
	"strconv"
	"time"
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
func (h ConnHandler) Connect(ctx *app.Ctx) error {
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
	if motor.Status != motorModel.BikeAvailable {
		return errorsx.BadRequestError("Bu Motorbisiklet şu anda müsait değil!")
	}

	if err = h.connService.CreateConn(ctx.Context(), &connection); err != nil {
		return errorsx.InternalError(err, "Bağlantı kurulurken hata oluştu!")
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"info": "Bağlantı kuruldu!"})
}

// when disconnect the motor status = available. but the lock status does not change. The lock status will be checked when the user sends a photo!
func (h ConnHandler) Disconnect(ctx *app.Ctx) error {
	param := ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return errorsx.BadRequestError("Hatalı istek!")
	}

	connection, err := h.connService.GetConnByID(ctx.Context(), id)
	if err != nil {
		if errorsx.Is(err, gorm.ErrRecordNotFound) {
			return errorsx.NotFoundError("Böyle bir bağlantı yok!")
		}
		return errorsx.InternalError(err, "Bir hata oluştu!")
	}

	// zaten bağlantı koptuysa..
	if connection.DisconnectedAt != nil {
		return errorsx.BadRequestError("Zaten bağlantı kopmuş!")
	}

	now := time.Now()
	connection.DisconnectedAt = &now

	if err = h.connService.UpdateConn(ctx.Context(), connection); err != nil {
		return errorsx.InternalError(err, "Bağlantı kesilemedi!")
	}

	motor, err := h.motorService.GetMotorByID(ctx.Context(), int(connection.MotorbikeID))
	if err != nil {
		return errorsx.InternalError(err, "Bir hata oluştu!")
	}

	motor.Status = motorModel.BikeAvailable

	err = h.motorService.UpdateMotor(ctx.Context(), motor)
	if err != nil {
		return errorsx.InternalError(err, "Motor status güncellenirken hata oluştu!")
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"info": "Bağlantı kesildi!"})
}

// why i need this func? maybe admin wanna delete history connection? idk bro! but maybe they needs to use this func!
func (h ConnHandler) DeleteConn(ctx *app.Ctx) error {
	param := ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return errorsx.BadRequestError("Hatalı istek!")
	}

	if err = h.connService.DeleteConn(ctx.Context(), id); err != nil {
		if errorsx.Is(err, gorm.ErrRecordNotFound) {
			return errorsx.NotFoundError("Böyle bir bağlantı zaten yok!")
		}
		return errorsx.InternalError(err, "Bağlantı silinirken bir hata oluştu!")
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"info": "Bağlantı başarılı bir şekilde silindi!"})
}
