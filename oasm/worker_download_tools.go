package oasm

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	pb "github.com/oasm-platform/open-asm/grpc-client/go/workers"
)

func (c *Client) WorkerDownloadTools(ctx context.Context) error {
	l := NewLogger("Worker.Sync")

	absToolPath, err := filepath.Abs(c.toolPath)
	if err != nil {
		l.ErrorE("Failed to resolve absolute tool path", err)
		return err
	}

	if err := os.MkdirAll(absToolPath, 0o755); err != nil {
		return fmt.Errorf("failed to create tool directory: %w", err)
	}

	registry, err := c.Workers().BuiltinToolRegistry(ctx, &pb.BuiltinToolRegistryRequest{})
	if err != nil {
		l.ErrorE("BuiltinToolRegistry retrieval failed", err)
		return err
	}

	osKey := runtime.GOOS
	if osKey == "darwin" {
		osKey = "macos"
	}

	var osTools []string
	switch osKey {
	case "linux":
		osTools = registry.Linux
	case "windows":
		osTools = registry.Windows
	case "macos":
		osTools = registry.Macos
	default:
		return fmt.Errorf("unsupported OS: %s", osKey)
	}

	statePath := filepath.Join(absToolPath, ".tool_versions.json")
	installedTools := loadToolState(statePath)

	for _, toolUrl := range osTools {
		fileName := filepath.Base(toolUrl)

		if installedTools[fileName] {
			l.Success("Tools cache hit: %s", fileName)
			continue
		}

		l.Info("Downloading tool: %s", fileName)
		if err := c.downloadAndExtractSingleTool(ctx, toolUrl, absToolPath, fileName); err != nil {
			l.ErrorE("Failed to download/extract tool", err, fileName)
			return err
		}

		installedTools[fileName] = true
		if err := saveToolState(statePath, installedTools); err != nil {
			l.ErrorE("Failed to save tool state", err)
		}
	}

	manifest, err := c.Workers().GetManifest(ctx, &pb.GetManifestRequest{})
	if err == nil && len(manifest.InitCommands) > 0 {
		l.Info("Executing %d initialization commands", len(manifest.InitCommands))
		for _, cmdStr := range manifest.InitCommands {
			if err := c.runInitCommand(ctx, cmdStr, absToolPath); err != nil {
				l.ErrorE("Command failed", err, cmdStr)
				return err
			}
		}
		l.Success("All init commands executed successfully")
	}

	return nil
}

func (c *Client) downloadAndExtractSingleTool(ctx context.Context, url string, destDir string, fileName string) error {
	dlLog := NewLogger("Worker.Download")

	stream, err := c.Workers().Storage(ctx, &pb.StorageRequest{
		Path: url,
	})
	if err != nil {
		return fmt.Errorf("failed to start download stream: %w", err)
	}

	tempFile := filepath.Join(destDir, fileName)
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			file.Close()
			os.Remove(tempFile)
			return fmt.Errorf("error receiving stream: %w", err)
		}

		if _, err = file.WriteAt(resp.Chunk, int64(resp.Offset)); err != nil {
			file.Close()
			os.Remove(tempFile)
			return fmt.Errorf("failed to write chunk: %w", err)
		}

		if resp.Eof {
			break
		}
	}
	file.Close()
	dlLog.Success("Download completed: %s", tempFile)

	extLog := NewLogger("Worker.Extract")
	extLog.Info("Extracting %s...", fileName)

	if strings.HasSuffix(fileName, ".zip") {
		err = c.extractZip(tempFile, destDir, extLog)
	} else if strings.HasSuffix(fileName, ".tar.gz") || strings.HasSuffix(fileName, ".tgz") {
		err = c.extractTarGz(tempFile, destDir, extLog)
	} else {
		return fmt.Errorf("unsupported archive format: %s", fileName)
	}

	if err != nil {
		return fmt.Errorf("failed to extract and set permissions: %w", err)
	}

	// Clean up the downloaded archive
	_ = os.Remove(tempFile)
	return nil
}

func isIgnoredFile(fileName string) bool {
	lowerName := strings.ToLower(fileName)
	return strings.HasSuffix(lowerName, ".txt") || strings.HasSuffix(lowerName, ".md") || strings.HasSuffix(lowerName, ".pdf")
}

func (c *Client) extractZip(srcZip string, destDir string, l *LoggerType) error {
	r, err := zip.OpenReader(srcZip)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if isIgnoredFile(f.Name) {
			continue
		}

		target := filepath.Join(destDir, f.Name)
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(destDir)) {
			return fmt.Errorf("illegal file path: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(target, 0o755)
			continue
		}

		os.MkdirAll(filepath.Dir(target), 0o755)

		outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}

		// Ensure binary is executable on Unix-like systems
		if runtime.GOOS != "windows" {
			_ = os.Chmod(target, f.Mode()|0o755)
		}
		l.Verbose("Extracted: %s", f.Name)
	}
	return nil
}

func (c *Client) extractTarGz(srcGzip string, destDir string, l *LoggerType) error {
	file, err := os.Open(srcGzip)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if isIgnoredFile(header.Name) {
			continue
		}

		target := filepath.Join(destDir, header.Name)
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(destDir)) {
			return fmt.Errorf("illegal file path: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}

			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()

			// Ensure binary is executable on Unix-like systems
			if runtime.GOOS != "windows" {
				_ = os.Chmod(target, os.FileMode(header.Mode)|0o755)
			}
			l.Verbose("Extracted: %s", header.Name)
		}
	}
	return nil
}

func (c *Client) runInitCommand(ctx context.Context, cmdStr string, workDir string) error {
	l := NewLogger("Worker.Init")
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return nil
	}

	binaryName := parts[0]
	args := parts[1:]
	fullPath := filepath.Join(workDir, binaryName)

	if _, err := os.Stat(fullPath); err == nil {
		binaryName = fullPath
	}

	l.Debug("Running: %s", cmdStr)
	cmd := exec.CommandContext(ctx, binaryName, args...)
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	pathEnv := os.Getenv("PATH")
	cmd.Env = append(os.Environ(), fmt.Sprintf("PATH=%s%c%s", workDir, os.PathListSeparator, pathEnv))

	return cmd.Run()
}

// State management helpers
func loadToolState(path string) map[string]bool {
	state := make(map[string]bool)
	data, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(data, &state)
	}
	return state
}

func saveToolState(path string, state map[string]bool) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
