package audit

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/zain-23/local-vault/apps/server/internal/common/apperror"
	"github.com/zain-23/local-vault/apps/server/internal/common/response"
	"github.com/zain-23/local-vault/apps/server/internal/common/validate"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// List — GET /api/v1/workspaces/:wid/audit
func (h *Handler) List(c *fiber.Ctx) error {
	var q ListQuery
	if err := c.QueryParser(&q); err != nil {
		return apperror.ErrInvalidBody
	}
	if msg := validate.Struct(q); msg != "" { // validates page/limit bounds
		return apperror.New(400, msg)
	}
	q.Normalize() // fill page/limit defaults

	resp, err := h.svc.List(c.UserContext(), c.Params("wid"), q)
	if err != nil {
		return err
	}
	return response.Success(c, resp, fiber.StatusOK, "audit events retrieved")
}

// Export — GET /api/v1/workspaces/:wid/audit/export (owner/admin only).
func (h *Handler) Export(c *fiber.Ctx) error {
	var q ListQuery
	if err := c.QueryParser(&q); err != nil {
		return apperror.ErrInvalidBody
	}
	// No pagination for export — same filters apply.
	events, err := h.svc.Export(c.UserContext(), c.Params("wid"), q)
	if err != nil {
		return err
	}

	// Build the CSV in memory (audit volume is low).
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"created_at", "action", "actor_id", "actor_name", "target_type", "target_id", "target_name", "device_id", "ip", "details"})
	for _, e := range events {
		details := ""
		if e.Details != nil {
			b, _ := json.Marshal(e.Details) // details column as compact JSON
			details = string(b)
		}
		_ = w.Write([]string{
			e.CreatedAt.Format(time.RFC3339),
			e.Action, e.ActorID, e.ActorName,
			e.TargetType, e.TargetID, e.TargetName,
			e.DeviceID, e.IP, details,
		})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return apperror.ErrInternal
	}

	// Attachment headers so the browser downloads a dated file.
	filename := fmt.Sprintf("audit-log-%s.csv", time.Now().Format("2006-01-02"))
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	return c.Send(buf.Bytes())
}
