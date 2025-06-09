package networkconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateIdentifier_Correct(t *testing.T) {
	t.Parallel()

	identifiers := []string{
		"",
		"name",
		"a",
		"dev1",
		"host.local",
		"o000000000",
		"-",
		"-0-0-0-0-0-",
		"a.0.-",
		".",
		"..",
		"用户",
		"____test____",
	}

	for _, identifier := range identifiers {
		assert.NoError(t, ValidateIdentifier(identifier))
	}
}

func TestValidateIdentifier_Incorrect(t *testing.T) {
	t.Parallel()

	identifiers := []string{
		"1",
		"5guys",
		"123",
		"user!",
		"@mail",
		"a b",
		"a/b",
		"a\tb",
		"a\nb",
	}

	for _, identifier := range identifiers {
		assert.Error(t, ValidateIdentifier(identifier))
	}
}

func TestValidateIP_Correct(t *testing.T) {
	t.Parallel()

	validIPs := []string{
		"",
		"127.0.0.1",
		"::1",
		"192.168.1.100",
		"2001:db8::68",
		"2001:0db8:0000:0000:0000:0000:0000:0068",
	}

	for _, ip := range validIPs {
		assert.NoError(t, ValidateIP(ip))
	}
}

func TestValidateIP_Incorrect(t *testing.T) {
	t.Parallel()

	invalidIPs := []string{
		"not_an_ip",
		"256.256.256.256",
		"192.168.1.",
		"2001:db8:::68",
		"gggg::1",
		"1234:5678:90ab:cdef:ghij:klmn:opqr:stuv",
		"2001:db8::68::1",
	}

	for _, ip := range invalidIPs {
		assert.Error(t, ValidateIP(ip))
	}
}

func TestValidateInterface(t *testing.T) {
	t.Parallel()

	type testCase struct {
		iface   Interface
		wantErr bool
	}

	tests := []testCase{
		{
			iface: Interface{
				Name:    "eth0",
				Address: "192.168.1.10",
				Mask:    "255.255.255.0",
				Gateway: "192.168.1.1",
			},
			wantErr: false,
		},
		{
			iface: Interface{
				Name:    "123", // invalid name
				Address: "192.168.1.10",
				Mask:    "255.255.255.0",
				Gateway: "192.168.1.1",
			},
			wantErr: true,
		},
		{
			iface: Interface{
				Name:    "eth0",
				Address: "999.999.999.999", // invalid address
				Mask:    "255.255.255.0",
				Gateway: "192.168.1.1",
			},
			wantErr: true,
		},
		{
			iface: Interface{
				Name:    "eth0",
				Address: "192.168.1.10",
				Mask:    "not_a_mask", // invalid mask
				Gateway: "192.168.1.1",
			},
			wantErr: true,
		},
		{
			iface: Interface{
				Name:    "eth0",
				Address: "192.168.1.10",
				Mask:    "255.255.255.0",
				Gateway: "bad_gateway", // invalid gateway
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		err := validateInterface(&tt.iface)
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestValidateConfig(t *testing.T) {
	t.Parallel()

	type testCase struct {
		cfg     NetworkConfig
		wantErr bool
	}

	tests := []testCase{
		{
			cfg: NetworkConfig{
				Hostname:   "host1",
				Interfaces: nil,
			},
			wantErr: false,
		},
		{
			cfg: NetworkConfig{
				Hostname: "host1",
				Interfaces: []Interface{
					{
						Name:    "eth0",
						Address: "192.168.1.10",
						Mask:    "255.255.255.0",
						Gateway: "192.168.1.1",
					},
				},
			},
			wantErr: false,
		},
		{
			cfg: NetworkConfig{
				Hostname: "host1",
				Interfaces: []Interface{
					{
						Name:    "eth0",
						Address: "192.168.1.10",
						Mask:    "255.255.255.0",
						Gateway: "192.168.1.1",
					},
					{
						Name:    "eth1",
						Address: "192.168.1.11",
						Mask:    "255.255.255.0",
						Gateway: "192.168.1.1",
					},
				},
			},
			wantErr: false,
		},
		{
			cfg: NetworkConfig{
				Hostname: "bad host", // invalid hostname
				Interfaces: []Interface{
					{
						Name:    "eth0",
						Address: "192.168.1.10",
						Mask:    "255.255.255.0",
						Gateway: "192.168.1.1",
					},
				},
			},
			wantErr: true,
		},
		{
			cfg: NetworkConfig{
				Hostname: "host1",
				Interfaces: []Interface{
					{
						Name:    "eth0",
						Address: "192.168.1.10",
						Mask:    "255.255.255.0",
						Gateway: "192.168.1.1",
					},
					{
						Name:    "eth1",
						Address: "192.168.1.11",
						Mask:    "255.255.255.0",
						Gateway: "192.168.1.1",
					},
					{
						Name:    "eth2",
						Address: "192.168.1.12",
						Mask:    "255.255.255.0",
						Gateway: "192.168.1.1",
					},
				},
			},
			wantErr: true, // more than maxInterfaceCount
		},
		{
			cfg: NetworkConfig{
				Hostname: "host1",
				Interfaces: []Interface{
					{
						Name:    "123", // invalid interface
						Address: "192.168.1.10",
						Mask:    "255.255.255.0",
						Gateway: "192.168.1.1",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		err := validate(&tt.cfg)
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
