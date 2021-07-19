package archive

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"compress/flate"

	"github.com/mholt/archiver/v3"
	"github.com/sirupsen/logrus"
)

type Archiver interface {
	String() string
	Archive([]string, string) error
	Unarchive(string, string) error
	Read() (archiver.File, error)
	Write(archiver.File) error
	Open(io.Reader, int64) error
	Create(io.Writer) error
	Close() error
	Walk(string, archiver.WalkFunc) error
}

// NewArchiver create a new archiver for tar archive manipultation
func NewArchiver(ext string) (Archiver, error) {

	// Check extestion of target file
	f, err := archiver.ByExtension(ext)

	if err != nil {
		return nil, fmt.Errorf("error parsing type %s for format: %v", ext, err)
	}

	// Create tar
	mytar := &archiver.Tar{
		OverwriteExisting:      true,
		MkdirAll:               true,
		ImplicitTopLevelFolder: false,
		StripComponents:        0,
		ContinueOnError:        false,
	}

	// Check compression type (if using)
	switch v := f.(type) {
	case *archiver.Tar:
		return mytar, nil
	case *archiver.TarBz2:
		v.Tar = mytar
		v.CompressionLevel = flate.DefaultCompression
		return v, nil
	case *archiver.TarGz:
		v.Tar = mytar
		v.CompressionLevel = flate.DefaultCompression
		return v, nil
	case *archiver.TarZstd:
		v.Tar = mytar
		return v, nil
	default:
		return nil, fmt.Errorf("format does not support customization: %s", f)
	}
}

// CreateArchive will create a tar archive from source directory
func CreateArchive(a Archiver, rootDir, output string) error {

	// Get current directory
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	logrus.Infof("Building bundle %s with directory %s", output, rootDir)

	// Change directory to root and archive. ChDir needed for symlinks.
	os.Chdir(rootDir)
	if err := a.Archive([]string{"."}, filepath.Join(cwd, output)); err != nil {
		logrus.Errorf("failed to archive bundle %s: %v", output, err)
		return err
	}
	return nil
}

// CreateSplitAchrive will create multiple tar archives from source directory
func CreateSplitArchive(a Archiver, destDir, prefix string, maxSplitSize int64, sourceDir string) error {

	// Declare split variables
	splitNum := 0
	splitSize := int64(0)
	splitPath := fmt.Sprintf("%s/%s_%06d.%s", destDir, prefix, splitNum, a.String())

	// Open first split tar archive
	splitFile, err := os.Create(splitPath)

	if err != nil {
		logrus.Errorf("error creating archive path %s: %v", splitPath, err)
		return fmt.Errorf("creating %s: %v", splitPath, err)
	}

	logrus.Infof("Creating archive %s", splitPath)
	// Create a new tar archive for writing
	if a.Create(splitFile); err != nil {
		logrus.Errorf("error creating archive: %v", err)
		return fmt.Errorf("creating archive %s: %v", splitPath, err)
	}

	sourceInfo, err := os.Stat(sourceDir)

	if err != nil {
		return fmt.Errorf("%s: stat: %v", sourceDir, err)
	}

	return filepath.Walk(sourceDir, func(fpath string, info os.FileInfo, err error) error {

		if err != nil {
			return fmt.Errorf("traversing %s: %v", fpath, err)
		}
		if info == nil {
			return fmt.Errorf("no file info")
		}

		// build the name to be used within the archive
		nameInArchive, err := makeNameInArchive(sourceInfo, sourceDir, sourceDir, fpath)
		if err != nil {
			return fmt.Errorf("creating %s: %v", nameInArchive, err)
		}

		var file io.ReadCloser
		if info.Mode().IsRegular() {
			file, err = os.Open(fpath)
			if err != nil {
				return fmt.Errorf("%s: opening: %v", fpath, err)
			}
			defer file.Close()
		}

		f := archiver.File{
			FileInfo: archiver.FileInfo{
				FileInfo:   info,
				CustomName: nameInArchive,
			},
			ReadCloser: file,
		}

		// If the file is too large close the current tar archive and file
		if info.Size()+splitSize > maxSplitSize {

			a.Close()
			splitFile.Close()

			// Create new tar archive for writing
			splitNum += 1
			splitSize = int64(0)
			splitPath = fmt.Sprintf("%s/%s_%06d.%s", destDir, prefix, splitNum, a.String())

			splitFile, err = os.Create(splitPath)

			if err != nil {
				return err
			}

			logrus.Infof("Creating archive %s", splitPath)
			if err := a.Create(splitFile); err != nil {
				return err
			}

		}

		// Write file to current archive file
		if err = a.Write(f); err != nil {
			return fmt.Errorf("%s: writing: %s", fpath, err)
		}

		splitSize += info.Size()

		return nil
	})
}

// makeNameInArchive is a helper function pulled from mholt library
func makeNameInArchive(sourceInfo os.FileInfo, source, baseDir, fpath string) (string, error) {
	name := filepath.Base(fpath) // start with the file or dir name
	if sourceInfo.IsDir() {
		// preserve internal directory structure; that's the path components
		// between the source directory's leaf and this file's leaf
		dir, err := filepath.Rel(filepath.Dir(source), filepath.Dir(fpath))
		if err != nil {
			return "", err
		}
		// prepend the internal directory structure to the leaf name,
		// and convert path separators to forward slashes as per spec
		name = path.Join(filepath.ToSlash(dir), name)
	}
	return path.Join(baseDir, name), nil // prepend the base directory
}
