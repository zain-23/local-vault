package device

import (
	"github.com/gofiber/fiber/v2"

	"github.com/zain-23/local-vault/server/internal/common/apperror"
	"github.com/zain-23/local-vault/server/internal/common/middleware"
	"github.com/zain-23/local-vault/server/internal/common/response"
	"github.com/zain-23/local-vault/server/internal/common/validate"
)

// Handler turns HTTP requests into service calls.
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Authorize — POST /api/v1/device/authorize (no auth; the CLI has no credentials yet).
func (h *Handler) Authorize(c *fiber.Ctx) error {
	var req AuthorizeRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.ErrInvalidBody
	}
	// validate.Struct returns a message string ("" = valid), not an error.
	if msg := validate.Struct(req); msg != "" {
		return apperror.New(400, msg)
	}
	// c.IP() is the client address — recorded so the approval screen can show
	// "from 192.168.1.1" and the user can spot a request that isn't theirs.
	res, err := h.svc.Authorize(c.UserContext(), req, c.IP())
	if err != nil {
		return err
	}
	return response.Success(c, res, 201, "device authorization started")
}

// ApprovalDetails — GET /api/v1/device/authorize/:userCode (auth required).
func (h *Handler) ApprovalDetails(c *fiber.Ctx) error {
	res, err := h.svc.ApprovalDetails(c.UserContext(), c.Params("userCode"))
	if err != nil {
		return err
	}
	return response.Success(c, res, 200, "device authorization request")
}

// Decide — PUT /api/v1/device/authorize/:userCode (auth required).
func (h *Handler) Decide(c *fiber.Ctx) error {
	var req DecisionRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.ErrInvalidBody
	}
	if msg := validate.Struct(req); msg != "" {
		return apperror.New(400, msg)
	}
	user := middleware.GetUser(c) // identity from the JWT — set by the Auth middleware
	if err := h.svc.Decide(c.UserContext(), c.Params("userCode"), req, user.ID); err != nil {
		return err
	}
	return response.Success(c, nil, 200, "device authorization "+req.Action+"d")
}

// Poll — POST /api/v1/device/authorize/poll (no auth; the device_code is the credential).
func (h *Handler) Poll(c *fiber.Ctx) error {
	var req PollRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.ErrInvalidBody
	}
	if msg := validate.Struct(req); msg != "" {
		return apperror.New(400, msg)
	}
	res, err := h.svc.Poll(c.UserContext(), req, c.IP(), c.Get("User-Agent"))
	if err != nil {
		return err
	}
	return response.Success(c, res, 200, "device authorization status")
}

// List — GET /api/v1/device?workspace_id=... (auth required).
func (h *Handler) List(c *fiber.Ctx) error {
	var q ListDevicesQuery
	if err := c.QueryParser(&q); err != nil { // QueryParser reads the query string, not the body
		return apperror.ErrInvalidBody
	}
	if msg := validate.Struct(q); msg != "" {
		return apperror.New(400, msg)
	}
	user := middleware.GetUser(c)
	devices, err := h.svc.ListDevices(c.UserContext(), q, user.ID)
	if err != nil {
		return err
	}
	return response.Success(c, devices, 200, "devices retrieved")
}

// Revoke — DELETE /api/v1/device/:id (auth required).
func (h *Handler) Revoke(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	if err := h.svc.RevokeDevice(c.UserContext(), c.Params("id"), user.ID); err != nil {
		return err
	}
	return response.Success(c, nil, 200, "device revoked")
}
