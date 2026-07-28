package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/zain-23/local-vault/apps/server/internal/common/apperror"
)

// MembershipChecker is anything that can report a user's role in a workspace.
// *workspace.Store satisfies it — defining the interface here avoids an import cycle.
type MembershipChecker interface {
	RoleOf(ctx context.Context, workspaceID, userID string) (string, error)
}

// RequireRole returns middleware that allows the request only if the caller's role
// in the workspace (read from the `param` route param, e.g. "id" or "wid") is in allowed.
// Must run AFTER Auth middleware, which populates the user in Locals.
func RequireRole(checker MembershipChecker, param string, allowed ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := GetUser(c)              // caller identity from the JWT (set by Auth)
		workspaceID := c.Params(param) // e.g. /workspaces/:id -> c.Params("id")

		role, err := checker.RoleOf(c.UserContext(), workspaceID, user.ID)
		if err != nil {
			return apperror.ErrInternal
		}
		if role == "" {
			// Not a member — hide the workspace's existence entirely
			return apperror.New(404, "workspace not found")
		}

		// Allow if the caller's role matches any of the allowed roles
		for _, r := range allowed {
			if role == r {
				c.Locals("role", role) // stash for the handler if it wants it
				return c.Next()         // pass control to the next handler
			}
		}
		return apperror.ErrForbidden // member, but not high enough role
	}
}
