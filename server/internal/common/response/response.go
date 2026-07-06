package response

import "github.com/gofiber/fiber/v2"

// Envelope writes a success response shape for every endpoint
type Envelope struct {
	Data	any 	`json:"data"`	// payload on success, null on error
	Success	bool	`json:"success"`
	Message	string	`json:"message"`
	Status	int		`json:"status"`
}

// Success writes a success envelope with the given status code
func Success(c *fiber.Ctx, data any, status int, message string) error {
	return c.Status(status).JSON(Envelope{
		Data: data,
		Success: true,
		Message: message,
		Status: status,
	})
}


// Error writes an error envelope â called by the central ErrorHandler
func Error(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(Envelope{
		Data: nil,
		Success: false,
		Message: message,
		Status: status,
	})
}