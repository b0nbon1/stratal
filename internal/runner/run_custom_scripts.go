package runner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/b0nbon1/stratal/internal/storage/db/dto"
)

// Language configurations for different script types
var languageConfig = map[string]struct {
	interpreter string
	extension   string
	args        []string
}{
	"python":     {"python3", ".py", []string{}},
	"javascript": {"node", ".js", []string{}},
	"typescript": {"tsx", ".ts", []string{}},
	"bash":       {"bash", ".sh", []string{}},
	"sh":         {"sh", ".sh", []string{}},
	"ruby":       {"ruby", ".rb", []string{}},
	"go":         {"go", ".go", []string{"run"}},
	"php":        {"php", ".php", []string{}},
	"perl":       {"perl", ".pl", []string{}},
}

// RunCustomScriptWithOutputs runs a custom script with environment variables containing outputs from previous tasks
func RunCustomScriptWithOutputs(ctx context.Context, script *dto.ScriptConfig, outputs map[string]string) (string, error) {
	if script == nil {
		return "", fmt.Errorf("script configuration is nil")
	}

	if script.Code == "" {
		return "", fmt.Errorf("script code is empty")
	}

	language := strings.ToLower(script.Language)
	config, exists := languageConfig[language]
	if !exists {
		return "", fmt.Errorf("unsupported script language: %s", script.Language)
	}

	// Create temporary directory for script execution
	tempDir, err := os.MkdirTemp("", "stratal-script-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Write script to temporary file
	scriptFile := filepath.Join(tempDir, "script"+config.extension)
	if err := os.WriteFile(scriptFile, []byte(script.Code), 0600); err != nil {
		return "", fmt.Errorf("failed to write script file: %w", err)
	}

	// Prepare command
	args := append(config.args, scriptFile)
	cmd := exec.CommandContext(ctx, config.interpreter, args...)
	cmd.Dir = tempDir

	// Set up output capture
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Set up environment variables with task outputs
	cmd.Env = os.Environ()

	// Add TASK_OUTPUT_ prefix to all outputs
	for taskName, output := range outputs {
		envName := fmt.Sprintf("TASK_OUTPUT_%s", strings.ToUpper(strings.ReplaceAll(taskName, "-", "_")))
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", envName, output))
	}

	// Create a channel to signal completion
	done := make(chan error, 1)

	// Execute script with timeout handling
	go func() {
		done <- cmd.Run()
	}()

	// Wait for completion or timeout
	select {
	case <-ctx.Done():
		// Context cancelled, kill the process
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return "", fmt.Errorf("script execution cancelled: %w", ctx.Err())

	case err := <-done:
		output := stdout.String()

		if err != nil {
			// Include stderr in error message
			errorOutput := stderr.String()
			if errorOutput != "" {
				return output, fmt.Errorf("script execution failed: %v\nstderr: %s", err, errorOutput)
			}
			return output, fmt.Errorf("script execution failed: %v", err)
		}

		// If stderr has content but execution succeeded, append it to output
		if stderrContent := stderr.String(); stderrContent != "" {
			output += "\n--- stderr ---\n" + stderrContent
		}

		return output, nil
	}
}

func RunCustomScript(ctx context.Context, script *dto.ScriptConfig) (string, error) {
	// Call the new function with empty outputs for backward compatibility
	return RunCustomScriptWithOutputs(ctx, script, nil)
}

// ValidateScript checks if a script is valid without executing it
func ValidateScript(script *dto.ScriptConfig) error {
	if script == nil {
		return fmt.Errorf("script configuration is nil")
	}

	if script.Code == "" {
		return fmt.Errorf("script code is empty")
	}

	language := strings.ToLower(script.Language)
	if _, exists := languageConfig[language]; !exists {
		return fmt.Errorf("unsupported script language: %s", script.Language)
	}

	return nil
}

// ExecuteScriptWithTimeout runs a script with a specific timeout
func ExecuteScriptWithTimeout(script *dto.ScriptConfig, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return RunCustomScript(ctx, script)
}

// StreamCustomScript executes a script and streams output in real-time
func StreamCustomScript(ctx context.Context, script *dto.ScriptConfig, outputWriter io.Writer) error {
	if err := ValidateScript(script); err != nil {
		return err
	}

	language := strings.ToLower(script.Language)
	config := languageConfig[language]

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "stratal-script-stream-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Write script to temporary file
	scriptFile := filepath.Join(tempDir, "script"+config.extension)
	if err := os.WriteFile(scriptFile, []byte(script.Code), 0600); err != nil {
		return fmt.Errorf("failed to write script file: %w", err)
	}

	// Prepare command
	args := append(config.args, scriptFile)
	cmd := exec.CommandContext(ctx, config.interpreter, args...)
	cmd.Dir = tempDir

	// Set up pipes for real-time output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start script: %w", err)
	}

	// Stream output
	go io.Copy(outputWriter, stdout)
	go io.Copy(outputWriter, stderr)

	// Wait for completion
	return cmd.Wait()
}
