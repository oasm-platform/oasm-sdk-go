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

	statePath := filepath.Join(absToolPath, ".tool_versions.json")

	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		l.Info("Tool cache not found, cleaning up directory contents for fresh download")

		entries, err := os.ReadDir(absToolPath)
		if err == nil {
			for _, entry := range entries {
				removePath := filepath.Join(absToolPath, entry.Name())
				if err := os.RemoveAll(removePath); err != nil {
					return fmt.Errorf("failed to remove item %s: %w", removePath, err)
				}
			}
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("failed to read tool directory for cleanup: %w", err)
		}
	}

	if err := os.MkdirAll(absToolPath, 0o755); err != nil {
		return fmt.Errorf("failed to create tool directory: %w", err)
	}

	osKey := runtime.GOOS
	if osKey == "darwin" {
		osKey = "macos"
	}

	archKey := runtime.GOARCH

	registry, err := c.Workers().BuiltinToolRegistry(ctx, &pb.BuiltinToolRegistryRequest{Os: osKey, Arch: archKey})
	if err != nil {
		l.ErrorE("BuiltinToolRegistry retrieval failed", err)
		return err
	}

	oldState := loadToolState(statePath)
	newState := make(map[string][]string)

	for _, toolUrl := range registry.ToolPaths {
		fileName := filepath.Base(toolUrl)

		if extractedFiles, exists := oldState[fileName]; exists {
			l.Success("Tools cache hit: %s", fileName)
			newState[fileName] = extractedFiles
			continue
		}

		l.Info("Downloading tool: %s", fileName)
		extractedFiles, err := c.downloadAndExtractSingleTool(ctx, toolUrl, absToolPath, fileName)
		if err != nil {
			l.ErrorE("Failed to download/extract tool", err, fileName)
			return err
		}

		newState[fileName] = extractedFiles
	}

	activeFiles := make(map[string]bool)
	for _, files := range newState {
		for _, f := range files {
			activeFiles[f] = true
		}
	}

	for oldFileName, oldExtractedFiles := range oldState {
		if _, stillExists := newState[oldFileName]; !stillExists {
			l.Info("Cleaning up obsolete tool: %s", oldFileName)
			for _, file := range oldExtractedFiles {
				if !activeFiles[file] {
					fullPath := filepath.Join(absToolPath, file)
					_ = os.Remove(fullPath)
					l.Verbose("Deleted unused file: %s", file)
				}
			}
		}
	}

	if err := saveToolState(statePath, newState); err != nil {
		l.ErrorE("Failed to save tool state", err)
	}

	manifest, err := c.Workers().GetManifest(ctx, &pb.GetManifestRequest{})
	if err != nil {
		l.ErrorE("Failed to retrieve GetManifest for init commands", err)
	} else if len(manifest.InitCommands) > 0 {
		l.Info("Executing %d initialization commands", len(manifest.InitCommands))
		for _, cmdStr := range manifest.InitCommands {
			if err := c.runInitCommand(ctx, cmdStr, absToolPath); err != nil {
				l.ErrorE("Init command failed", err, cmdStr)
				return err
			}
		}
		l.Success("All init commands executed successfully")
	} else {
		l.Debug("GetManifest success, but no init commands to execute")
	}

	return nil
}

func (c *Client) downloadAndExtractSingleTool(ctx context.Context, url string, destDir string, fileName string) ([]string, error) {
	dlLog := NewLogger("Worker.Download")
	var extractedFiles []string

	stream, err := c.Workers().Storage(ctx, &pb.StorageRequest{Path: url})
	if err != nil {
		return nil, fmt.Errorf("failed to start download stream: %w", err)
	}

	tempFile := filepath.Join(destDir, fileName)
	file, err := os.Create(tempFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			file.Close()
			os.Remove(tempFile)
			return nil, fmt.Errorf("error receiving stream: %w", err)
		}
		if _, err = file.WriteAt(resp.Chunk, int64(resp.Offset)); err != nil {
			file.Close()
			os.Remove(tempFile)
			return nil, fmt.Errorf("failed to write chunk: %w", err)
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
		extractedFiles, err = c.extractZip(tempFile, destDir, extLog)
	} else if strings.HasSuffix(fileName, ".tar.gz") || strings.HasSuffix(fileName, ".tgz") {
		extractedFiles, err = c.extractTarGz(tempFile, destDir, extLog)
	} else {
		return nil, fmt.Errorf("unsupported archive format: %s", fileName)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to extract and set permissions: %w", err)
	}

	_ = os.Remove(tempFile)
	return extractedFiles, nil
}

func (c *Client) extractZip(srcZip string, destDir string, l *LoggerType) ([]string, error) {
	var extractedFiles []string
	r, err := zip.OpenReader(srcZip)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for _, f := range r.File {
		if isIgnoredFile(f.Name) {
			continue
		}

		target := filepath.Join(destDir, f.Name)
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(destDir)) {
			return nil, fmt.Errorf("illegal file path: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(target, 0o755)
			continue
		}

		os.MkdirAll(filepath.Dir(target), 0o755)

		outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, f.Mode())
		if err != nil {
			return nil, err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return nil, err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return nil, err
		}

		if runtime.GOOS != "windows" {
			_ = os.Chmod(target, f.Mode()|0o755)
		}

		extractedFiles = append(extractedFiles, f.Name) // Lưu lại path tương đối
		l.Verbose("Extracted: %s", f.Name)
	}
	return extractedFiles, nil
}

func (c *Client) extractTarGz(srcGzip string, destDir string, l *LoggerType) ([]string, error) {
	var extractedFiles []string
	file, err := os.Open(srcGzip)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if isIgnoredFile(header.Name) {
			continue
		}

		target := filepath.Join(destDir, header.Name)
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(destDir)) {
			return nil, fmt.Errorf("illegal file path: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return nil, err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return nil, err
			}

			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return nil, err
			}

			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return nil, err
			}
			f.Close()

			if runtime.GOOS != "windows" {
				_ = os.Chmod(target, os.FileMode(header.Mode)|0o755)
			}

			extractedFiles = append(extractedFiles, header.Name)
			l.Verbose("Extracted: %s", header.Name)
		}
	}
	return extractedFiles, nil
}

func isIgnoredFile(fileName string) bool {
	lowerName := strings.ToLower(fileName)
	return strings.HasSuffix(lowerName, ".txt") || strings.HasSuffix(lowerName, ".md") || strings.HasSuffix(lowerName, ".pdf")
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

	if runtime.GOOS == "windows" && !strings.HasSuffix(strings.ToLower(fullPath), ".exe") {
		if _, err := os.Stat(fullPath + ".exe"); err == nil {
			fullPath += ".exe"
		}
	}

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
func loadToolState(path string) map[string][]string {
	state := make(map[string][]string)
	data, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(data, &state)
	}
	return state
}

func saveToolState(path string, state map[string][]string) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
