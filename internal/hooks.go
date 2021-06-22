package internal

// Hooks is a collection of func.
type Hooks []func()

// Add adds a hook.
func (h *Hooks) Add(g func()) {
	*h = append(*h, g)
}

// Server hooks
var (
	OnServerRun Hooks // Hooks executes after the server starts.
	OnShutdown  Hooks // Hooks executes when the server is shutdown gracefully.
)
