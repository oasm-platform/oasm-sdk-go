package oasm

import (
	"archive/tar"
	"compress/gzip"
	"context"
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

	manifest, err := c.Workers().GetManifest(ctx, &pb.GetManifestRequest{})
	if err != nil {
		l.ErrorE("Manifest retrieval failed", err)
		return err
	}

	entries, err := os.ReadDir(absToolPath)
	if err != nil && !os.IsNotExist(err) {
		l.ErrorE("Failed to read tool directory", err)
	}

	if err != nil || len(entries) == 0 {
		l.Info("Local tools missing. Starting synchronization...")

		if manifest.DownloadToolsUrl == "" {
			return fmt.Errorf("manifest contains empty download URL")
		}

		if err := c.downloadAndExtractTools(ctx, manifest.DownloadToolsUrl, absToolPath); err != nil {
			return err
		}

		l.Success("Tools synchronized successfully")
	} else {
		l.Success("Tools cache hit: %s", absToolPath)
	}

	if len(manifest.InitCommands) > 0 {
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

func (c *Client) downloadAndExtractTools(ctx context.Context, url string, destDir string) error {
	dlLog := NewLogger("Worker.Download")
	dlLog.Info("Starting download from: %s", url)

	stream, err := c.Workers().DownloadTools(ctx, &pb.DownloadToolsRequest{
		Url: url,
	})
	if err != nil {
		return fmt.Errorf("failed to start download stream: %w", err)
	}

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return fmt.Errorf("failed to create tool directory: %w", err)
	}

	tempGzip := filepath.Join(destDir, "tools_download.tar.gz")
	file, err := os.Create(tempGzip)
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
			return fmt.Errorf("error receiving stream: %w", err)
		}

		if _, err = file.WriteAt(resp.Chunk, int64(resp.Offset)); err != nil {
			file.Close()
			return fmt.Errorf("failed to write chunk: %w", err)
		}

		if resp.Eof {
			break
		}
	}
	file.Close()
	dlLog.Success("Download completed: %s", tempGzip)

	extLog := NewLogger("Worker.Extract")
	extLog.Info("Extracting tools to %s...", destDir)

	if err := c.extractAndChmod(tempGzip, destDir, extLog); err != nil {
		return fmt.Errorf("failed to extract and set permissions: %w", err)
	}

	_ = os.Remove(tempGzip)
	return nil
}

func (c *Client) extractAndChmod(srcGzip string, destDir string, l *LoggerType) error {
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

			if runtime.GOOS != "windows" {
				_ = os.Chmod(target, 0o755)
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
	if runtime.GOOS == "windows" && !strings.HasSuffix(fullPath, ".exe") {
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
