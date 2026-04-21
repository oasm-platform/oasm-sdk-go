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
	manifest, err := c.Workers().GetManifest(ctx, &pb.GetManifestRequest{})
	if err != nil {
		return fmt.Errorf("failed to get manifest: %w", err)
	}

	downloadUrl := manifest.DownloadToolsUrl
	if downloadUrl == "" {
		return fmt.Errorf("download URL is empty")
	}

	stream, err := c.Workers().DownloadTools(ctx, &pb.DownloadToolsRequest{
		Url: downloadUrl,
	})
	if err != nil {
		return fmt.Errorf("failed to start download stream: %w", err)
	}

	absToolPath, err := filepath.Abs(c.toolPath)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute tool path: %w", err)
	}

	if err := os.MkdirAll(absToolPath, 0o755); err != nil {
		return fmt.Errorf("failed to create tool directory: %w", err)
	}

	tempGzip := filepath.Join(absToolPath, "tools_download.tar.gz")
	file, err := os.Create(tempGzip)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer file.Close()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error receiving stream: %w", err)
		}

		_, err = file.WriteAt(resp.Chunk, int64(resp.Offset))
		if err != nil {
			return fmt.Errorf("failed to write chunk at offset %d: %w", resp.Offset, err)
		}

		if resp.Eof {
			break
		}
	}

	fmt.Printf("Extracting tools to %s...\n", absToolPath)

	if err := c.extractAndChmod(tempGzip, absToolPath); err != nil {
		return fmt.Errorf("failed to extract and set permissions: %w", err)
	}

	if err := os.Remove(tempGzip); err != nil {
		fmt.Printf("Warning: failed to remove temporary file %s: %v\n", tempGzip, err)
	}

	fmt.Println("Tools updated and execution permissions granted successfully!")
	if len(manifest.InitCommands) > 0 {
		fmt.Println("Executing initialization commands...")
		for _, cmdStr := range manifest.InitCommands {
			parts := strings.Fields(cmdStr)
			if len(parts) == 0 {
				continue
			}

			binaryName := parts[0]
			args := parts[1:]

			fullPath := filepath.Join(absToolPath, binaryName)

			// if runtime.GOOS == "windows" && !strings.HasSuffix(fullPath, ".exe") {
			// 	fullPath += ".exe"
			// }

			if _, err := os.Stat(fullPath); err == nil {
				binaryName = fullPath
			}

			cmd := exec.CommandContext(ctx, binaryName, args...)
			cmd.Dir = absToolPath

			pathEnv := os.Getenv("PATH")
			cmd.Env = append(os.Environ(), fmt.Sprintf("PATH=%s%c%s", absToolPath, os.PathListSeparator, pathEnv))

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			fmt.Printf("Running: %s\n", cmdStr)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to execute init command '%s': %w", cmdStr, err)
			}
		}
	}
	return nil
}

func (c *Client) extractAndChmod(srcGzip string, destDir string) error {
	file, err := os.Open(srcGzip)
	if err != nil {
		return fmt.Errorf("failed to open source gzip: %w", err)
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	absDestDir, err := filepath.Abs(destDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for destDir: %w", err)
	}

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading tar archive: %w", err)
		}

		if header.Name == "" {
			continue
		}

		target := filepath.Join(absDestDir, header.Name)

		absTarget, err := filepath.Abs(target)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for target %s: %w", target, err)
		}

		if !strings.HasPrefix(absTarget, absDestDir+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path (Zip Slip detected): %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(absTarget, 0o755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", absTarget, err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(absTarget), 0o755); err != nil {
				return fmt.Errorf("failed to create parent directory for %s: %w", absTarget, err)
			}

			f, err := os.OpenFile(absTarget, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
			defer f.Close()
			if err != nil {
				return fmt.Errorf("failed to open file %s: %w", absTarget, err)
			}

			if _, err := io.Copy(f, tr); err != nil {
				return fmt.Errorf("failed to extract file %s: %w", absTarget, err)
			}

			if runtime.GOOS != "windows" {
				if err := os.Chmod(absTarget, 0o755); err != nil {
					return fmt.Errorf("failed to set execution permission for %s: %w", absTarget, err)
				}
			}
		}
	}
	return nil
}
