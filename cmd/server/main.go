package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"

	"promptgo/internal/tui"
)

func main() {
	// Get port from env or default to 2222
	port := os.Getenv("PORT")
	if port == "" {
		port = "2222"
	}

	// Get host key path from user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	keyPath := filepath.Join(homeDir, ".ssh", "promptgo_host_key")

	// Create SSH server with Wish
	s, err := wish.NewServer(
		wish.WithAddress(":"+port),
		wish.WithHostKeyPath(keyPath),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Setup graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("🐹 PromptGo SSH server starting")
	log.Printf("   Port: %s", port)
	log.Printf("   Host key: %s", keyPath)
	log.Printf("   Connect with: ssh localhost -p %s", port)

	// Start server in goroutine
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-done
	log.Println("Shutting down server...")

	// Graceful shutdown with 30-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

// teaHandler creates a new Bubble Tea program for each SSH session
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// Log the connection
	log.Printf("New connection from %s", s.RemoteAddr())

	// Create new TUI model for this session
	m := tui.NewModel()

	// Configure program options
	opts := []tea.ProgramOption{
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	}

	return m, opts
}
