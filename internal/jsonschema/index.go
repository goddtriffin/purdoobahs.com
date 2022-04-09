package jsonschema

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/purdoobahs/purdoobahs.com/internal/logger"

	"github.com/xeipuuv/gojsonschema"
)

func ValidateJsonSchema(logger logger.ILogger) (bool, error) {
	invalidFiles, err := validatePurdoobahJsonSchema(logger)
	if err != nil {
		return invalidFiles, err
	}

	invalidFiles, err = validateTraditionJsonSchema(logger)
	if err != nil {
		return invalidFiles, err
	}

	// no invalid files
	return false, nil
}

func validatePurdoobahJsonSchema(logger logger.ILogger) (bool, error) {
	// read in the Purdoobah JSON Schema
	purdoobahJSONSchemaFilepath := "./assets/purdoobahs/_purdoobah.schema.json"
	b, err := ioutil.ReadFile(purdoobahJSONSchemaFilepath)
	if err != nil {
		logger.Error(fmt.Sprintf(
			"error reading file: %s",
			"./assets/purdoobahs/_purdoobah.schema.json"),
		)
		return true, err
	}
	schema := gojsonschema.NewStringLoader(string(b))

	// find all the individual Purdoobah files
	filepaths, err := walkMatch("./assets/purdoobahs/", `*.json`)
	if err != nil {
		logger.Error("error parsing Purdoobah assets directory")
		return true, err
	}

	// loop through each file
	invalidFiles := false
	for _, path := range filepaths {
		// ignore _purdoobah.schema.json and _template.json
		if strings.Contains(path, "_") {
			continue
		}

		// read in the Purdoobah JSON document
		b, err := ioutil.ReadFile(path)
		if err != nil {
			logger.Error(fmt.Sprintf("error reading file: %s", path))
			return true, err
		}
		document := gojsonschema.NewStringLoader(string(b))

		// validate the document against the schema
		result, err := gojsonschema.Validate(schema, document)
		if err != nil {
			logger.Error(fmt.Sprintf("error validating file: %s", path))
			return true, err
		}

		// if not valid, print errors
		if !result.Valid() {
			invalidFiles = true
			for _, desc := range result.Errors() {
				logger.Error(fmt.Sprintf("validation error (%s): %s", path, desc))
			}
		}
	}

	return invalidFiles, nil
}

func validateTraditionJsonSchema(logger logger.ILogger) (bool, error) {
	// read in the Tradition JSON Schema
	traditionJSONSchemaFilepath := "./assets/traditions/_tradition.schema.json"
	b, err := ioutil.ReadFile(traditionJSONSchemaFilepath)
	if err != nil {
		logger.Error(fmt.Sprintf(
			"error reading file: %s",
			"./assets/traditions/_tradition.schema.json"),
		)
		return true, err
	}
	schema := gojsonschema.NewStringLoader(string(b))

	// find all the individual Tradition files
	filepaths, err := walkMatch("./assets/traditions/", `*.json`)
	if err != nil {
		logger.Error("error parsing Tradition assets directory")
		return true, err
	}

	// loop through each file
	invalidFiles := false
	for _, path := range filepaths {
		// ignore _tradition.schema.json and _template.json
		if strings.Contains(path, "_") {
			continue
		}

		// read in the Tradition JSON document
		b, err := ioutil.ReadFile(path)
		if err != nil {
			logger.Error(fmt.Sprintf("error reading file: %s", path))
			return true, err
		}
		document := gojsonschema.NewStringLoader(string(b))

		// validate the document against the schema
		result, err := gojsonschema.Validate(schema, document)
		if err != nil {
			logger.Error(fmt.Sprintf("error validating file: %s", path))
			return true, err
		}

		// if not valid, print errors
		if !result.Valid() {
			invalidFiles = true
			for _, desc := range result.Errors() {
				logger.Error(fmt.Sprintf("validation error (%s): %s", path, desc))
			}
		}
	}

	return invalidFiles, nil
}

func walkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
