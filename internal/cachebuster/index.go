package cachebuster

import (
	"fmt"
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
		uiStaticDirectoryPath:  "../ui/static",
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

	uiFullDirPath := fmt.Sprintf("%s%s", cb.uiStaticDirectoryPath, dirPath)
	binFullDirPath := fmt.Sprintf("%s%s", cb.binStaticDirectoryPath, dirPath)

	err := filepath.Walk(binFullDirPath, func(path string, info os.FileInfo, err error) error {
		uiFullFilePath := fmt.Sprintf("%s/%s", uiFullDirPath, path)
		binFullFilePath := fmt.Sprintf("%s/%s", binFullDirPath, path)

		if err != nil {
			return err
		}

		// we only want files
		if info.IsDir() {
			return nil
		}

		if !cb.isFileHashable(path) {
			fmt.Printf("File not hashable: %s\n", path)
			fmt.Println()
			cb.cache[path] = binFullFilePath
			return nil
		}

		originalInfo, err := os.Stat(fmt.Sprintf("%s/%s", cb.uiStaticDirectoryPath, path))

		fmt.Printf("Name: %s\n", info.Name())
		fmt.Printf("ModTime: %s\n", info.ModTime())
		fmt.Printf("Size: %d bytes\n", info.Size())
		fmt.Println()

		hash := cb.generateHash(info.Size())
		cb.cache[path] = fmt.Sprintf(
			"/%s/%s.%s%s",
			rootDirectory,
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

	fmt.Printf("Total files: %d\n", count)
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

func (cb *CacheBuster) generateHash(sizeBytes int64) string {
	return fmt.Sprintf("%d", sizeBytes)
}

func (cb *CacheBuster) printCache() error {
	builder := ""
	for key, value := range cb.cache {
		builder += fmt.Sprintf("%s\t=>\t%s\n", key, value)
	}

	err := ioutil.WriteFile("./cache-buster.txt", []byte(builder), 0644)
	if err != nil {
		return err
	}
	return nil
}
