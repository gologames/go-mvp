// Functions to load, validate and save network configuration in YAML format.
package networkconfig

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"regexp"

	pathcheck "github.com/gologames/go-mvp/internal/security"
	"gopkg.in/yaml.v3"
)

// NetworkConfig represents network configuration structure.
type NetworkConfig struct {
	Hostname   string      `yaml:"hostname"`
	Interfaces []Interface `yaml:"interfaces"`
}

// Interface describes a single network interface with address settings.
type Interface struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
	Mask    string `yaml:"mask"`
	Gateway string `yaml:"gateway"`
}

// FileReader defines an interface for validating paths and reading files.
type FileReader interface {
	ValidatePath(path string) error
	ReadFile(path string) ([]byte, error)
}

// OSFileReader implements FileReader using OS functions.
type OSFileReader struct{}

// ValidatePath uses pathcheck to validate the given path.
func (OSFileReader) ValidatePath(path string) error {
	return pathcheck.ValidatePath(path)
}

// ReadFile reads the file from disk using os.ReadFile.
func (OSFileReader) ReadFile(path string) ([]byte, error) {
	//nolint:gosec
	return os.ReadFile(path)
}

// FileWriter defines an interface for writing files.
type FileWriter interface {
	WriteFile(path string, data []byte, perm os.FileMode) error
}

// OSFileWriter implements FileWriter using os.WriteFile.
type OSFileWriter struct{}

// WriteFile writes data to a file with the given permissions.
func (OSFileWriter) WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

// Load reads and parses a network configuration file.
func Load(path string, fr FileReader) (*NetworkConfig, error) {
	if err := fr.ValidatePath(path); err != nil {
		return nil, err
	}

	data, err := fr.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg NetworkConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	err = validate(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

const maxInterfaceCount = 2

// Save validates and writes a network configuration file.
func Save(path string, cfg *NetworkConfig, logger *slog.Logger, writer FileWriter) error {
	logger.Debug("saving network config", slog.String("path", path))

	if err := validate(cfg); err != nil {
		return err
	}

	cfgCopy := *cfg
	for len(cfgCopy.Interfaces) < maxInterfaceCount {
		cfgCopy.Interfaces = append(cfgCopy.Interfaces, Interface{})
	}

	data, err := yaml.Marshal(cfgCopy)
	if err != nil {
		return err
	}

	return writer.WriteFile(path, data, 0o600)
}

func validate(cfg *NetworkConfig) error {
	err := ValidateIdentifier(cfg.Hostname)
	if err != nil {
		return fmt.Errorf("invalid hostname: %w", err)
	}

	count := len(cfg.Interfaces)
	if count > maxInterfaceCount {
		return fmt.Errorf("number of interfaces must not exceed %d: got %d", maxInterfaceCount, count)
	}

	for i := range count {
		err = validateInterface(&cfg.Interfaces[i])
		if err != nil {
			return fmt.Errorf("invalid interface%d: %w", i+1, err)
		}
	}

	return nil
}

func validateInterface(iface *Interface) error {
	err := ValidateIdentifier(iface.Name)
	if err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	validateIP := func(label, ip string) error {
		if err := ValidateIP(ip); err != nil {
			return fmt.Errorf("invalid %s: %w", label, err)
		}
		return nil
	}

	if err := validateIP("address", iface.Address); err != nil {
		return err
	}
	if err := validateIP("mask", iface.Mask); err != nil {
		return err
	}
	if err := validateIP("gateway", iface.Gateway); err != nil {
		return err
	}

	return nil
}

// ValidateIdentifier checks if a string is a valid identifier or empty.
func ValidateIdentifier(identifier string) error {
	const identifierPattern = `^$|^[-\p{L}_.]+[-\p{L}\p{N}_.]*$`

	match, err := regexp.MatchString(identifierPattern, identifier)
	if err != nil {
		return err
	}

	if !match {
		return fmt.Errorf("invalid identifier %q", identifier)
	}
	return nil
}

// ValidateIP checks if a string is a valid IP address or empty.
func ValidateIP(ip string) error {
	if ip != "" && net.ParseIP(ip) == nil {
		return fmt.Errorf("invalid ip %q", ip)
	}
	return nil
}
