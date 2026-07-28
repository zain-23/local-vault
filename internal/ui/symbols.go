package ui

const (
	iconSuccess = "✓"
	iconError   = "✗"
	iconWarn    = "⚠"
	iconInfo    = "›"
	wordmark    = "LV"
)

var auditGlyphs = map[string]string{
	"add":    "+",
	"update": "~",
	"remove": "-",
	"rotate": "↻",
}

// AuditGlyph returns a refined glyph for a vault audit action.
func AuditGlyph(action string) string {
	if g, ok := auditGlyphs[action]; ok {
		return g
	}
	return "•"
}
