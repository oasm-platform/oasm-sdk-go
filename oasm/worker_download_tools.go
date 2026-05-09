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
	absToolPath, err := filepath.Abs(c.toolPath)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute tool path: %w", err)
	}

	manifest, err := c.Workers().GetManifest(ctx, &pb.GetManifestRequest{})
	if err != nil {
		return fmt.Errorf("failed to get manifest: %w", err)
	}

	entries, err := os.ReadDir(absToolPath)
	isDirEmpty := err != nil || len(entries) == 0

	if isDirEmpty {
		downloadURL := manifest.DownloadToolsUrl
		if downloadURL == "" {
			return fmt.Errorf("download URL is empty")
		}

		stream, err := c.Workers().DownloadTools(ctx, &pb.DownloadToolsRequest{
			Url: downloadURL,
		})
		if err != nil {
			return fmt.Errorf("failed to start download stream: %w", err)
		}

		if err := os.MkdirAll(absToolPath, 0o755); err != nil {
			return fmt.Errorf("failed to create tool directory: %w", err)
		}

		tempGzip := filepath.Join(absToolPath, "tools_download.tar.gz")
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

			_, err = file.WriteAt(resp.Chunk, int64(resp.Offset))
			if err != nil {
				file.Close()
				return fmt.Errorf("failed to write chunk: %w", err)
			}

			if resp.Eof {
				break
			}
		}
		file.Close()

		fmt.Printf("Extracting tools to %s...\n", absToolPath)
		if err := c.extractAndChmod(tempGzip, absToolPath); err != nil {
			return fmt.Errorf("failed to extract and set permissions: %w", err)
		}
		_ = os.Remove(tempGzip)
	} else {
		Logger("Sync").Log(fmt.Sprintf("Tools exist in %s, skipping download.\n", absToolPath))
	}

	if len(manifest.InitCommands) > 0 {
		for _, cmdStr := range manifest.InitCommands {
			parts := strings.Fields(cmdStr)
			if len(parts) == 0 {
				continue
			}

			binaryName := parts[0]
			args := parts[1:]

			fullPath := filepath.Join(absToolPath, binaryName)
			if runtime.GOOS == "windows" && !strings.HasSuffix(fullPath, ".exe") {
				if _, err := os.Stat(fullPath + ".exe"); err == nil {
					fullPath += ".exe"
				}
			}

			if _, err := os.Stat(fullPath); err == nil {
				binaryName = fullPath
			}

			cmd := exec.CommandContext(ctx, binaryName, args...)
			cmd.Dir = absToolPath
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			pathEnv := os.Getenv("PATH")
			cmd.Env = append(os.Environ(), fmt.Sprintf("PATH=%s%c%s", absToolPath, os.PathListSeparator, pathEnv))

			fmt.Printf("Running init command: %s\n", cmdStr)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to execute init command '%s': %w", cmdStr, err)
			}
		}
	}

	fmt.Println("All steps completed successfully!")
	return nil
}

func (c *Client) extractAndChmod(srcGzip string, destDir string) error {
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

		// Prevent Zip Slip
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
		}
	}
	return nil
}
