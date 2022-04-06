package cachebuster

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type CacheBuster struct {
	// cache is a mapping of original static asset file name to one with an added unique, deterministic hash
	cache map[string]string

	// uiStaticDirectoryPath is the path to the ui directory where the static assets are originally from
	uiStaticDirectoryPath string

	// binStaticDirectoryPath is the path to the bin directory where the static assets are copied to
	binStaticDirectoryPath string

	// allowedFileExtensions are file extensions that are automatically known to need a hash value for cache busting
	allowedFileExtensions []string
}

func NewCacheBuster() (*CacheBuster, error) {
	cb := &CacheBuster{
		cache:                  make(map[string]string),
		binStaticDirectoryPath: "static",
		allowedFileExtensions: []string{
			".pdf",
			".webp",
			".ico",
			".svg",
			".js",
			".css",
			".mp4",
		},
	}

	count, err := cb.hashStaticAssets()
	if err != nil {
		return &CacheBuster{}, err
	}

	err = cb.printCache()
	if err != nil {
		return &CacheBuster{}, err
	}
	fmt.Printf("Total files hashed for Cache-Busting: %d\n", count)

	return cb, nil
}

func (cb *CacheBuster) hashStaticAssets() (int, error) {
	totalCount := 0

	count, err := cb.walk("/file")
	if err != nil {
		return totalCount, err
	}
	totalCount += count

	count, err = cb.walk("/image/bot")
	if err != nil {
		return totalCount, err
	}
	totalCount += count

	count, err = cb.walk("/image/favicon")
	if err != nil {
		return totalCount, err
	}
	totalCount += count

	count, err = cb.walk("/image/logo")
	if err != nil {
		return totalCount, err
	}
	totalCount += count

	count, err = cb.walk("/image/purdoobah")
	if err != nil {
		return totalCount, err
	}
	totalCount += count

	count, err = cb.walk("/image/section")
	if err != nil {
		return totalCount, err
	}
	totalCount += count

	count, err = cb.walk("/image/socials")
	if err != nil {
		return totalCount, err
	}
	totalCount += count

	count, err = cb.walk("/image/tradition")
	if err != nil {
		return totalCount, err
	}
	totalCount += count

	count, err = cb.walk("/script")
	if err != nil {
		return totalCount, err
	}
	totalCount += count

	count, err = cb.walk("/stylesheet")
	if err != nil {
		return totalCount, err
	}
	totalCount += count

	count, err = cb.walk("/video")
	if err != nil {
		return totalCount, err
	}
	totalCount += count

	return totalCount, nil
}

func (cb *CacheBuster) walk(dirPath string) (int, error) {
	count := 0

	binFullDirPath := fmt.Sprintf("%s%s", cb.binStaticDirectoryPath, dirPath)
	err := filepath.Walk(binFullDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// we only want files
		if info.IsDir() {
			return nil
		}

		// only hash files we pre-select to want to hash
		if !cb.isFileHashable(path) {
			fmt.Printf("File not hashable: %s\n", path)
			fmt.Println()
			cb.cache[path] = path
			return nil
		}

		// generate unique hash of the file
		hash, err := cb.generateHash(path)
		if err != nil {
			return err
		}

		// generate new name
		cb.cache[path] = fmt.Sprintf(
			"/%s.%s%s",
			strings.TrimSuffix(path, filepath.Ext(path)),
			hash,
			filepath.Ext(path),
		)
		count += 1

		return nil
	})
	if err != nil {
		return count, err
	}

	return count, nil
}

func (cb *CacheBuster) isFileHashable(filename string) bool {
	// check if filename has immediate allow-listed file extension
	// for _, extension := range cb.allowedFileExtensions {
	// 	if filepath.Ext(filename) == extension {
	// 		return true
	// 	}
	// }

	// return false
	return true
}

func (cb *CacheBuster) generateHash(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func (cb *CacheBuster) printCache() error {
	builder := ""
	for _, value := range cb.cache {
		builder += fmt.Sprintf("%s\n", value)
	}

	err := ioutil.WriteFile("../cache-buster.txt", []byte(builder), 0644)
	if err != nil {
		return err
	}
	return nil
}
