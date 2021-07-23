package archive

import (
	"bufio"
	"bytes"
	"compress/flate"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/mholt/archiver/v3"
	"github.com/sirupsen/logrus"
)

const delimeter = "__checksum__"

type Archiver interface {
	String() string
	Archive([]string, string) error
	Unarchive(string, string) error
	Write(archiver.File) error
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
	// TODO: Allow user to specify compression level
	switch v := f.(type) {
	case *archiver.Tar:
		return mytar, nil
	case *archiver.TarGz:
		v.Tar = mytar
		v.CompressionLevel = flate.DefaultCompression
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

	// Change directory to root and archive. Chdir needed for symlinks.
	os.Chdir(rootDir)
	if err := a.Archive([]string{"."}, filepath.Join(cwd, output)); err != nil {
		return fmt.Errorf("failed to archive bundle %s: %v", output, err)
	}
	return nil
}

// CreateSplitAchrive will create multiple tar archives from source directory
func CreateSplitArchive(a Archiver, destDir, prefix string, maxSplitSize int64, sourceDir string) error {

	// Declare split variables
	splitNum := 0
	splitSize := int64(0)
	splitPath := fmt.Sprintf("%s/%s_%06d.%s", destDir, prefix, splitNum, a.String())

	// Create first split tar archive
	splitFile, err := os.Create(splitPath)

	if err != nil {
		return fmt.Errorf("creating %s: %v", splitPath, err)
	}

	// Create a new tar archive for writing
	logrus.Infof("Creating archive %s", splitPath)
	if a.Create(splitFile); err != nil {
		return fmt.Errorf("creating archive %s: %v", splitPath, err)
	}

	sourceInfo, err := os.Stat(sourceDir)

	if err != nil {
		return fmt.Errorf("%s: stat: %v", sourceDir, err)
	}

	filepath.Walk(sourceDir, func(fpath string, info os.FileInfo, err error) error {

		if err != nil {
			return fmt.Errorf("traversing %s: %v", fpath, err)
		}
		if info == nil {
			return fmt.Errorf("no file info")
		}

		// build the name to be used within the archive
		nameInArchive, err := archiver.NameInArchive(sourceInfo, sourceDir, fpath)
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

		// If the file is too large create a new one
		if info.Size()+splitSize > maxSplitSize {

			// Current current tar archive
			a.Close()
			splitFile.Close()

			// Increment split number and reset splitSize
			splitNum += 1
			splitSize = int64(0)
			splitPath = fmt.Sprintf("%s/%s_%06d.%s", destDir, prefix, splitNum, a.String())

			// Create a new tar archive for writing
			logrus.Infof("Creating archive %s", splitPath)

			splitFile, err = os.Create(splitPath)

			if err != nil {
				return fmt.Errorf("creating %s: %v", splitPath, err)
			}

			if err := a.Create(splitFile); err != nil {
				return fmt.Errorf("creating archive %s: %v", splitPath, err)
			}

		}

		// Write file to current archive file
		if err = a.Write(f); err != nil {
			return fmt.Errorf("%s: writing: %s", fpath, err)
		}

		splitSize += info.Size()

		return nil
	})

	a.Close()
	splitFile.Close()

	return nil
}

// ExtractArchive will unpack the archive at the specified directory
func ExtractArchive(a Archiver, src, dest string) error {
	return a.Unarchive(src, dest)
}

// VerifyArchive will verify the contents of the archive against the provided metadata file
// TODO: add more verification actions
func VerifyArchive(a Archiver, src string) error {

	scanForChecksum(src)
	return a.Walk(src, func(f archiver.File) error {
		fmt.Println("Filename:", f.Name())
		return nil
	})
}

// getHash is a helper function to get the checksum of a bundle
func generateCheckSum(input *os.File) []byte {

	hash := sha256.New()
	if _, err := io.Copy(hash, input); err != nil {
		return nil
	}

	return hash.Sum(nil)
}

// appendChecksum will conca the checksum of the
//archive to the archive
func appendChecksum(input string, sum []byte) error {

	var buf bytes.Buffer

	// append checksum to file
	b, err := ioutil.ReadFile(input)

	if err != nil {
		return err
	}

	buf.Write(b)
	buf.Write([]byte(delimeter))
	buf.Write(sum)

	if err = ioutil.WriteFile(input, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("error writing file %s: %v", input, err)
	}

	return nil
}

// scanForChecksum
func scanForChecksum(input string) (string, error) {

	// Specify delimiter
	r, _ := regexp.Compile(delimeter)

	// Open input file
	file, err := os.Open(input)
	if err != nil {
		return "", fmt.Errorf("error opening file %s: %v", input, err)
	}

	defer file.Close()

	// Create file scanner
	scanner := bufio.NewScanner(file)

	// Scan each line and return subsequent string if
	// delimeter is found
	for scanner.Scan() {

		if r.MatchString(scanner.Text()) {
			fmt.Println(scanner.Text())
			scanner.Scan()
			return scanner.Text(), nil
		}

	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("scanner error %v", err)
	}

	return "", fmt.Errorf("checksum is not available in file %s", input)
}
