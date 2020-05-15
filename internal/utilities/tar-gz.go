package util

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

// ExtractTarGz extracts the tar.gz file passed in path parameter
func ExtractTarGz(src string, dest string) ([]string, error) {
	var filenames []string
	gzipStream, err := os.Open(src)
	if err != nil {
		return filenames, err
	}
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return filenames, err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return filenames, err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(dest+header.Name, 0755); err != nil {
				return filenames, err
			}
		case tar.TypeReg:
			outFile, err := os.Create(dest + header.Name)
			if err != nil {
				return filenames, err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return filenames, err
			}
			outFile.Close()
			filenames = append(filenames, dest+header.Name)
		default:
			return filenames, fmt.Errorf("ExtractTarGz: uknown type: %b in %s", header.Typeflag, header.Name)
		}
	}
	// Close the file without defer to close before next iteration of loop
	uncompressedStream.Close()
	gzipStream.Close()
	return filenames, nil
}
