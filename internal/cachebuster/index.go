package cachebuster

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type CacheBuster struct {
	// cache is a mapping of original static asset file name to one with an added unique, deterministic hash
	cache     map[string]string
	cacheKeys []string

	// staticAssetsRootDirectoryPath is the path to the bin directory where the static assets are copied to
	staticAssetsRootDirectoryPath string

	// staticAssetsSubdirectoryPaths is a list of subdirectories to loop through when generating hashed static assets
	staticAssetsSubdirectoryPaths []string

	Debug bool
}

func NewCacheBuster(staticAssetsRootDirectoryPath string, staticAssetsSubdirectoryPaths []string) (*CacheBuster, error) {
	cb := &CacheBuster{
		cache:                         make(map[string]string),
		cacheKeys:                     []string{},
		staticAssetsRootDirectoryPath: staticAssetsRootDirectoryPath,
		staticAssetsSubdirectoryPaths: staticAssetsSubdirectoryPaths,
		Debug:                         false,
	}

	err := cb.hashStaticAssets()
	if err != nil {
		return &CacheBuster{}, err
	}

	if cb.Debug {
		fmt.Printf("Total files hashed for Cache-Busting: %d\n", len(cb.cache))
	}
	return cb, nil
}

// Get takes a path from root domain to a static asset (as it would be called from a browser, so with a leading slash)
// and returns the version of the filepath that contains a unique hash.
//
// e.g. "/static/image/favicon/favicon.ico" -> "/static/image/favicon/favicon.66189abc248d80832e458ee37e93c9e8.ico"
func (cb *CacheBuster) Get(path string) string {
	if val, ok := cb.cache[path]; ok {
		return val
	}

	if cb.Debug {
		fmt.Printf("file not found in CacheBuster: `%s`\n", path)
	}
	return ""
}

// Add takes a path from root domain to a static asset (as it would be called from a browser, so with a leading slash),
// generates a unique hash for that file, renames it on disk, and stores the uniquely-hashed filepaths in a cache for
// lookup later.
//
// e.g. "/static/image/favicon/favicon.ico" -> "/static/image/favicon/favicon.66189abc248d80832e458ee37e93c9e8.ico"
func (cb *CacheBuster) Add(nonHashedFilepath string) error {
	// generate unique hash of the file
	hash, err := cb.generateHash(nonHashedFilepath)
	if err != nil {
		return err
	}

	// generate new name with the unique hash
	hashedFilepath := fmt.Sprintf(
		"%s.%s%s",
		strings.TrimSuffix(nonHashedFilepath, filepath.Ext(nonHashedFilepath)),
		hash,
		filepath.Ext(nonHashedFilepath),
	)

	// store the hashed name in the cache
	cb.cacheKeys = append(cb.cacheKeys, nonHashedFilepath)
	cb.cache[nonHashedFilepath] = hashedFilepath

	// rename the file to the new name
	// the files can't be prepended by slashes, as that would point to the root directory of the computer as opposed to
	// finding the file from the current directory
	err = os.Rename(strings.TrimPrefix(nonHashedFilepath, "/"), strings.TrimPrefix(hashedFilepath, "/"))
	if err != nil {
		return err
	}

	return nil
}

// PrintToFile prints the hashed filepaths in the cache to a file.
//
// This is useful in the case that you want to check in a file to help diff what static asset hashes are modified
// in-between commits.
func (cb *CacheBuster) PrintToFile(outputFilepath string) error {
	// make sure to sort cacheKeys to ensure same-order read-out every access
	sort.Strings(cb.cacheKeys)

	builder := ""
	for _, key := range cb.cacheKeys {
		builder += fmt.Sprintf("%s\n", cb.cache[key])
	}

	err := ioutil.WriteFile(outputFilepath, []byte(builder), 0644)
	if err != nil {
		return err
	}
	return nil
}

// hashStaticAssets loops over every subdirectory of the static assets directory in order to generate unique hashes for
// all the files contained within.
func (cb *CacheBuster) hashStaticAssets() error {
	for _, subdirectory := range cb.staticAssetsSubdirectoryPaths {
		err := cb.walk(subdirectory)
		if err != nil {
			return err
		}
	}

	return nil
}

// walk the given directory, renames each file(path) by adding a unique hash, and saving the results for lookup
// later.
func (cb *CacheBuster) walk(dirPath string) error {
	binFullDirPath := fmt.Sprintf("%s%s", cb.staticAssetsRootDirectoryPath, dirPath)
	err := filepath.Walk(binFullDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// we only want files
		if info.IsDir() {
			return nil
		}

		nonHashedFilepath := fmt.Sprintf("/%s", path)
		err = cb.Add(nonHashedFilepath)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// generateHash generates a unique hash for the given file(path).
//
// This implementation generates a hash of the file by creating an MD5 hash of the file contents.
func (cb *CacheBuster) generateHash(filepath string) (string, error) {
	file, err := os.Open(strings.TrimPrefix(filepath, "/"))
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
