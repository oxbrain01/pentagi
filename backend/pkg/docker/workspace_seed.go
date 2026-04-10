package docker

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/sirupsen/logrus"
)

func (dc *dockerClient) maybeSeedFlowWorkspaceBind(logger *logrus.Entry, bindHostPath string) error {
	if dc.flowWorkspaceSeedDir == "" || bindHostPath == "" {
		return nil
	}
	src, err := filepath.Abs(dc.flowWorkspaceSeedDir)
	if err != nil {
		return fmt.Errorf("flow workspace seed: resolve source: %w", err)
	}
	fi, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("flow workspace seed: stat source %q: %w", src, err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("flow workspace seed: source is not a directory: %s", src)
	}
	if err := os.MkdirAll(bindHostPath, 0o755); err != nil {
		return fmt.Errorf("flow workspace seed: mkdir bind target: %w", err)
	}
	empty, err := isDirEmpty(bindHostPath)
	if err != nil {
		return fmt.Errorf("flow workspace seed: %w", err)
	}
	if !empty {
		logger.WithField("bind_host_path", bindHostPath).Debug("flow workspace seed skipped: bind target not empty")
		return nil
	}
	logger.WithFields(logrus.Fields{
		"seed_src": src,
		"seed_dst": bindHostPath,
	}).Info("flow workspace seed: copying into bind mount path")
	if err := copyDirTree(src, bindHostPath); err != nil {
		return fmt.Errorf("flow workspace seed: copy into bind mount: %w", err)
	}
	return nil
}

func (dc *dockerClient) maybeSeedFlowWorkspaceVolume(ctx context.Context, logger *logrus.Entry, containerID string) error {
	if dc.flowWorkspaceSeedDir == "" {
		return nil
	}
	src, err := filepath.Abs(dc.flowWorkspaceSeedDir)
	if err != nil {
		return fmt.Errorf("flow workspace seed: resolve source: %w", err)
	}
	fi, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("flow workspace seed: stat source %q: %w", src, err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("flow workspace seed: source is not a directory: %s", src)
	}
	empty, err := dc.isContainerWorkDirEmpty(ctx, containerID)
	if err != nil {
		return fmt.Errorf("flow workspace seed: check /work empty: %w", err)
	}
	if !empty {
		logger.Debug("flow workspace seed skipped: /work not empty")
		return nil
	}
	logger.WithFields(logrus.Fields{
		"seed_src":     src,
		"container_id": containerID,
	}).Info("flow workspace seed: copying into container volume at /work")

	r, w := io.Pipe()
	go func() {
		err := writeTarFromDir(src, w)
		_ = w.CloseWithError(err)
	}()
	defer func() { _ = r.Close() }()

	if err := dc.CopyToContainer(ctx, containerID, WorkFolderPathInContainer, r, container.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	}); err != nil {
		return fmt.Errorf("flow workspace seed: copy into volume: %w", err)
	}
	return nil
}

func isDirEmpty(dir string) (bool, error) {
	f, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer f.Close()
	names, err := f.Readdirnames(1)
	if err != nil {
		return false, err
	}
	return len(names) == 0, nil
}

func copyDirTree(srcRoot, dstRoot string) error {
	return filepath.WalkDir(srcRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcRoot, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		target := filepath.Join(dstRoot, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if !d.Type().IsRegular() {
			return nil
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return copyRegularFile(path, target)
	})
}

func copyRegularFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	info, err := in.Stat()
	if err != nil {
		return err
	}
	tmp := dst + ".pentagi-seed.tmp"
	out, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode()&0o777)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := out.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, dst)
}

func writeTarFromDir(srcRoot string, w io.Writer) error {
	tw := tar.NewWriter(w)
	defer tw.Close()
	return filepath.WalkDir(srcRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcRoot, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		hdr.Name = filepath.ToSlash(rel)
		if d.IsDir() {
			hdr.Name += "/"
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if d.IsDir() || !d.Type().IsRegular() {
			return nil
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		if _, err := io.Copy(tw, in); err != nil {
			_ = in.Close()
			return err
		}
		return in.Close()
	})
}

func (dc *dockerClient) isContainerWorkDirEmpty(ctx context.Context, containerID string) (bool, error) {
	const script = `test -z "$(ls -A /work 2>/dev/null)"`
	execResp, err := dc.client.ContainerExecCreate(ctx, containerID, container.ExecOptions{
		Cmd:          []string{"/bin/sh", "-c", script},
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return false, fmt.Errorf("exec create: %w", err)
	}
	attach, err := dc.client.ContainerExecAttach(ctx, execResp.ID, container.ExecAttachOptions{})
	if err != nil {
		return false, fmt.Errorf("exec attach: %w", err)
	}
	_, _ = io.Copy(io.Discard, attach.Reader)
	attach.Close()
	inspect, err := dc.client.ContainerExecInspect(ctx, execResp.ID)
	if err != nil {
		return false, fmt.Errorf("exec inspect: %w", err)
	}
	return inspect.ExitCode == 0, nil
}
