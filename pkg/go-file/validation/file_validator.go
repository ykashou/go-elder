package validation

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileValidator struct {
	Rules       []ValidationRule
	Checksums   map[string]string
	MaxFileSize int64
}

type ValidationRule struct {
	Name      string
	Validator func(string) error
	Required  bool
}

type ValidationResult struct {
	Valid    bool
	Errors   []string
	Warnings []string
	Checksum string
}

func NewFileValidator() *FileValidator {
	fv := &FileValidator{
		Rules:       make([]ValidationRule, 0),
		Checksums:   make(map[string]string),
		MaxFileSize: 100 * 1024 * 1024, // 100MB
	}
	
	fv.addDefaultRules()
	return fv
}

func (fv *FileValidator) addDefaultRules() {
	fv.Rules = append(fv.Rules, ValidationRule{
		Name:      "file_exists",
		Validator: fv.validateFileExists,
		Required:  true,
	})
	
	fv.Rules = append(fv.Rules, ValidationRule{
		Name:      "file_size",
		Validator: fv.validateFileSize,
		Required:  true,
	})
	
	fv.Rules = append(fv.Rules, ValidationRule{
		Name:      "file_extension",
		Validator: fv.validateFileExtension,
		Required:  false,
	})
}

func (fv *FileValidator) Validate(filepath string) ValidationResult {
	result := ValidationResult{
		Valid:    true,
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
	}
	
	checksum, err := fv.computeChecksum(filepath)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to compute checksum: %v", err))
		result.Valid = false
	} else {
		result.Checksum = checksum
	}
	
	for _, rule := range fv.Rules {
		err := rule.Validator(filepath)
		if err != nil {
			if rule.Required {
				result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", rule.Name, err))
				result.Valid = false
			} else {
				result.Warnings = append(result.Warnings, fmt.Sprintf("%s: %v", rule.Name, err))
			}
		}
	}
	
	return result
}

func (fv *FileValidator) validateFileExists(filepath string) error {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist")
	}
	return nil
}

func (fv *FileValidator) validateFileSize(filepath string) error {
	info, err := os.Stat(filepath)
	if err != nil {
		return err
	}
	
	if info.Size() > fv.MaxFileSize {
		return fmt.Errorf("file size %d exceeds maximum %d", info.Size(), fv.MaxFileSize)
	}
	
	return nil
}

func (fv *FileValidator) validateFileExtension(filepath string) error {
	ext := filepath.Ext(filepath)
	allowedExtensions := []string{".go", ".json", ".txt", ".md", ".yaml", ".yml"}
	
	for _, allowed := range allowedExtensions {
		if ext == allowed {
			return nil
		}
	}
	
	return fmt.Errorf("file extension %s not in allowed list", ext)
}

func (fv *FileValidator) computeChecksum(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func (fv *FileValidator) AddRule(name string, validator func(string) error, required bool) {
	rule := ValidationRule{
		Name:      name,
		Validator: validator,
		Required:  required,
	}
	fv.Rules = append(fv.Rules, rule)
}

func (fv *FileValidator) ValidateChecksum(filepath, expectedChecksum string) bool {
	actualChecksum, err := fv.computeChecksum(filepath)
	if err != nil {
		return false
	}
	
	return actualChecksum == expectedChecksum
}

func (fv *FileValidator) StoreChecksum(filepath string) error {
	checksum, err := fv.computeChecksum(filepath)
	if err != nil {
		return err
	}
	
	fv.Checksums[filepath] = checksum
	return nil
}

func (fv *FileValidator) VerifyIntegrity(filepath string) bool {
	if expectedChecksum, exists := fv.Checksums[filepath]; exists {
		return fv.ValidateChecksum(filepath, expectedChecksum)
	}
	return false
}
