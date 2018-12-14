/*
 Package configuration manage all about configuration.
*/

package configuration

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

//
// TestNewConfiguration tests configuration constructor
func TestNewConfiguration(t *testing.T) {
	tests := []struct {
		desc   string
		file   string
		folder string
		c      *Configuration
		err    error
	}{
		{
			desc:   "New configuration with file and folder",
			file:   "config.json",
			folder: "../test/conf.d",
			c: &Configuration{
				Name:     "apenella",
				IP:       "0.0.0.0",
				Port:     5497,
				Cluster:  []string{"1.1.1.1:5497", "2.2.2.2:5497"},
				Checks:   &ChecksConfiguration{},
				Services: &ServicesConfiguration{},
			},
			err: nil,
		},
		{
			desc:   "New configuration with default file",
			file:   "",
			folder: "../test/conf.d",
			c: &Configuration{
				Name:     "apenella",
				IP:       "0.0.0.0",
				Port:     5497,
				Cluster:  []string{"1.1.1.1:5497", "2.2.2.2:5497"},
				Checks:   &ChecksConfiguration{},
				Services: &ServicesConfiguration{},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		configuration, err := NewConfiguration(test.file, test.folder)
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err, err)
		} else {
			assert.Equal(t, configuration.Name, test.c.Name, "Unexpected node name.")
		}
	}
}

//
// TestValidateConfiguration test a valid configuration
func TestValidateConfiguration(t *testing.T) {
	tests := []struct {
		desc string
		c    *Configuration
		err  error
	}{
		{
			desc: "Valid configuration",
			c: &Configuration{
				Name:    "apenella",
				IP:      "0.0.0.0",
				Port:    5497,
				Cluster: []string{"1.1.1.1:5497", "2.2.2.2:5497"},
				Checks: &ChecksConfiguration{
					Folder:      "../test/conf.d/checks",
					MinInterval: 10,
					LoadTimeout: 60,
				},
				Services: &ServicesConfiguration{
					Folder:      "../test/conf.d/services",
					LoadTimeout: 60,
				},
			},
			err: nil,
		},
		{
			desc: "Configuration with no node name",
			c: &Configuration{
				Name:    "",
				IP:      "0.0.0.0",
				Port:    5497,
				Cluster: []string{"1.1.1.1:5497", "2.2.2.2:5497"},
				Checks: &ChecksConfiguration{
					Folder:      "../test/conf.d/checks",
					MinInterval: 10,
					LoadTimeout: 60,
				},
				Services: &ServicesConfiguration{
					Folder:      "../test/conf.d/services",
					LoadTimeout: 60,
				},
			},
			err: errors.New("(Configuration::ValidateConfiguration) Undefined name properties on configuration file"),
		},
		{
			desc: "Port lower than 0",
			c: &Configuration{
				Name:    "apenella",
				IP:      "0.0.0.0",
				Port:    -1,
				Cluster: []string{"1.1.1.1:5497", "2.2.2.2:5497"},
				Checks: &ChecksConfiguration{
					Folder:      "../test/conf.d/checks",
					MinInterval: 10,
					LoadTimeout: 60,
				},
				Services: &ServicesConfiguration{
					Folder:      "../test/conf.d/services",
					LoadTimeout: 60,
				},
			},
			err: errors.New("(Configuration::ValidateConfiguration) Invalid port definition"),
		},
		{
			desc: "Port greater than 65535",
			c: &Configuration{
				Name:    "apenella",
				IP:      "0.0.0.0",
				Port:    70000,
				Cluster: []string{"1.1.1.1:5497", "2.2.2.2:5497"},
				Checks: &ChecksConfiguration{
					Folder:      "../test/conf.d/checks",
					MinInterval: 10,
					LoadTimeout: 60,
				},
				Services: &ServicesConfiguration{
					Folder:      "../test/conf.d/services",
					LoadTimeout: 60,
				},
			},
			err: errors.New("(Configuration::ValidateConfiguration) Invalid port definition"),
		},
		{
			desc: "Configuration with no checks configuration",
			c: &Configuration{
				Name:    "apenella",
				IP:      "0.0.0.0",
				Port:    5497,
				Cluster: []string{"1.1.1.1:5497", "2.2.2.2:5497"},
				Checks: &ChecksConfiguration{
					Folder:      "../test/conf.d/checks",
					MinInterval: 10,
					LoadTimeout: 60,
				},
				Services: nil,
			},
			err: errors.New("(Configuration::ValidateConfiguration) (Configuration::validateServices) Undefined services properties on configuration file"),
		},
		{
			desc: "Configuration with no services configuration",
			c: &Configuration{
				Name:    "apenella",
				IP:      "0.0.0.0",
				Port:    5497,
				Cluster: []string{"1.1.1.1:5497", "2.2.2.2:5497"},
				Checks:  nil,
				Services: &ServicesConfiguration{
					Folder:      "../test/conf.d/services",
					LoadTimeout: 60,
				},
			},
			err: errors.New("(Configuration::ValidateConfiguration) (Configuration::validateChecks) Undefined checks properties on configuration file"),
		},
		{
			desc: "Invalid service load timeout",
			c: &Configuration{
				Name:    "apenella",
				IP:      "0.0.0.0",
				Port:    5497,
				Cluster: []string{"1.1.1.1:5497", "2.2.2.2:5497"},
				Checks: &ChecksConfiguration{
					Folder:      "../test/conf.d/checks",
					MinInterval: 10,
					LoadTimeout: 60,
				},
				Services: &ServicesConfiguration{
					Folder:      "../test/conf.d/services",
					LoadTimeout: -1,
				},
			},
			err: errors.New("(Configuration::ValidateConfiguration) (Configuration::validateServices) Invalid load timeout"),
		},
		{
			desc: "Invalid checks load timeout",
			c: &Configuration{
				Name:    "apenella",
				IP:      "0.0.0.0",
				Port:    5497,
				Cluster: []string{"1.1.1.1:5497", "2.2.2.2:5497"},
				Checks: &ChecksConfiguration{
					Folder:      "../test/conf.d/checks",
					MinInterval: 10,
					LoadTimeout: -1,
				},
				Services: &ServicesConfiguration{
					Folder:      "../test/conf.d/services",
					LoadTimeout: 60,
				},
			},
			err: errors.New("(Configuration::ValidateConfiguration) (Configuration::validateChecks) Invalid load timeout"),
		},
		{
			desc: "Invalid checks minimum interval",
			c: &Configuration{
				Name:    "apenella",
				IP:      "0.0.0.0",
				Port:    5497,
				Cluster: []string{"1.1.1.1:5497", "2.2.2.2:5497"},
				Checks: &ChecksConfiguration{
					Folder:      "../test/conf.d/checks",
					MinInterval: -1,
					LoadTimeout: 60,
				},
				Services: &ServicesConfiguration{
					Folder:      "../test/conf.d/services",
					LoadTimeout: 60,
				},
			},
			err: errors.New("(Configuration::ValidateConfiguration) (Configuration::validateChecks) Invalid minimum interval"),
		},
		{
			desc: "Unexistent checks folder",
			c: &Configuration{
				Name:    "apenella",
				IP:      "0.0.0.0",
				Port:    5497,
				Cluster: []string{"1.1.1.1:5497", "2.2.2.2:5497"},
				Checks: &ChecksConfiguration{
					Folder:      "../unexistent/checks",
					MinInterval: 10,
					LoadTimeout: 60,
				},
				Services: &ServicesConfiguration{
					Folder:      "../test/conf.d/services",
					LoadTimeout: 60,
				},
			},
			err: errors.New("(Configuration::ValidateConfiguration) (Configuration::validateChecks) Folder '../unexistent/checks' does not exist"),
		},
		{
			desc: "Unexistent checks folder",
			c: &Configuration{
				Name:    "apenella",
				IP:      "0.0.0.0",
				Port:    5497,
				Cluster: []string{"1.1.1.1:5497", "2.2.2.2:5497"},
				Checks: &ChecksConfiguration{
					Folder:      "../test/conf.d/checks",
					MinInterval: 10,
					LoadTimeout: 60,
				},
				Services: &ServicesConfiguration{
					Folder:      "../unexistent/services",
					LoadTimeout: 60,
				},
			},
			err: errors.New("(Configuration::ValidateConfiguration) (Configuration::validateServices) Folder '../unexistent/services' does not exist"),
		},
	}

	for _, test := range tests {
		t.Log(test.desc)
		err := test.c.ValidateConfiguration()
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err, err)
		}
	}
}
