package businesslogic

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

var MaxWrittenBytes = int64(3 * 1024 * 1024 * 1024)

func untar(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return errors.New("ExtractTarGz: Opening tar file failed")
	}

	uncompressedStream, err := gzip.NewReader(file)
	if err != nil {
		return errors.New("ExtractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("ExtractTarGz: Next() failed: %s", err.Error())
		}

		dst := filepath.Dir(fileName)
		dst = path.Join(dst, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(dst, 0755); err != nil {
				return fmt.Errorf("ExtractTarGz: Mkdir() failed: %s", err.Error())
			}
		case tar.TypeReg:
			outFile, err := os.Create(dst)
			if err != nil {
				return fmt.Errorf("ExtractTarGz: Create() failed: %s", err.Error())
			}

			limited := io.LimitReader(tarReader, MaxWrittenBytes)

			if _, err = io.Copy(outFile, limited); err != nil {
				return fmt.Errorf("ExtractTarGz: Copy() failed: %s", err.Error())
			}

			outFile.Close()
		case tar.TypeSymlink:
			if err := os.Symlink(header.Linkname, dst); err != nil {
				return err
			}
		default:
			return fmt.Errorf(
				"ExtractTarGz: uknown type: %x in %s",
				header.Typeflag,
				header.Name)
		}
	}
	return nil
}
