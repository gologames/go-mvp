package networkconfig

import (
	"errors"
	"log/slog"
	"os"
	"strings"
	"testing"

	networkconfig_mock "github.com/gologames/go-mvp/internal/networkconfig/mocks"
	"github.com/stretchr/testify/assert"
	testify_mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const testPath = "somepath.yaml"

func TestLoad_Success(t *testing.T) {
	t.Parallel()

	mock := getFileReaderMock(t, nil, []byte(`
hostname: test-host
interfaces:
  - name: eth0
    address: 192.168.0.10
    mask: 255.255.255.0
    gateway: 192.168.0.1
`), nil)

	cfg, err := Load(testPath, mock)
	require.NoError(t, err)
	assert.Equal(t, "test-host", cfg.Hostname)
	assert.Len(t, cfg.Interfaces, 1)
	assert.Equal(t, "eth0", cfg.Interfaces[0].Name)
}

func TestLoad_InvalidYAML(t *testing.T) {
	t.Parallel()
	mock := getFileReaderMock(t, nil, []byte("::::invalid:::yaml"), nil)
	_, err := Load(testPath, mock)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot unmarshal")
}

func TestLoad_InvalidPath(t *testing.T) {
	t.Parallel()
	const errBadPath = "bad path"
	mock := getFileReaderMock(t, errors.New(errBadPath), nil, nil)
	_, err := Load(testPath, mock)
	assert.EqualError(t, err, errBadPath)
}

func TestLoad_ReadFileError(t *testing.T) {
	t.Parallel()
	const errReadFile = "read error"
	mock := getFileReaderMock(t, nil, nil, errors.New(errReadFile))
	_, err := Load(testPath, mock)
	assert.EqualError(t, err, errReadFile)
}

func TestLoad_ValidateConfigError(t *testing.T) {
	t.Parallel()
	mock := getFileReaderMock(t, nil, []byte("hostname: 1nvalid_config"), nil)
	_, err := Load(testPath, mock)
	require.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "invalid hostname"))
}

func TestSave_Success(t *testing.T) {
	t.Parallel()

	mock := networkconfig_mock.NewMockFileWriter(t)
	mock.EXPECT().
		WriteFile(testPath, testify_mock.MatchedBy(func(data []byte) bool {
			str := string(data)
			return strings.Contains(str, "hostname: host") &&
				strings.Contains(str, "name: eth0")
		}), os.FileMode(0o600)).
		Return(nil)

	cfg := &NetworkConfig{
		Hostname: "host",
		Interfaces: []Interface{
			{
				Name:    "eth0",
				Address: "192.168.0.10",
				Mask:    "255.255.255.0",
				Gateway: "192.168.0.1",
			},
		},
	}

	err := Save(testPath, cfg, nopLogger(), mock)
	require.NoError(t, err)
}

func TestSave_WriteError(t *testing.T) {
	t.Parallel()

	const errWriteFile = "write failed"
	mock := networkconfig_mock.NewMockFileWriter(t)
	mock.EXPECT().
		WriteFile(testPath, testify_mock.Anything, os.FileMode(0o600)).
		Return(errors.New(errWriteFile))

	cfg := &NetworkConfig{Hostname: "valid"}
	err := Save(testPath, cfg, nopLogger(), mock)

	require.Error(t, err)
	assert.EqualError(t, err, errWriteFile)
}

func TestSave_ValidateConfigError(t *testing.T) {
	t.Parallel()

	mock := networkconfig_mock.NewMockFileWriter(t)
	cfg := &NetworkConfig{
		Hostname: "!!!",
	}

	err := Save(testPath, cfg, nopLogger(), mock)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid hostname")
}

func TestOSFileReaderWriter_Coverage(t *testing.T) {
	t.Parallel()

	fr := OSFileReader{}
	_ = fr.ValidatePath("")
	_, _ = fr.ReadFile("")

	fw := OSFileWriter{}
	_ = fw.WriteFile("", nil, 0)
}

func getFileReaderMock(t *testing.T, validateErr error, content []byte, fileReaderErr error) *networkconfig_mock.MockFileReader {
	t.Helper()
	mock := networkconfig_mock.NewMockFileReader(t)

	mock.EXPECT().ValidatePath(testPath).Return(validateErr)
	if content != nil || fileReaderErr != nil {
		mock.EXPECT().ReadFile(testPath).Return(content, fileReaderErr)
	}

	return mock
}

func nopLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}
