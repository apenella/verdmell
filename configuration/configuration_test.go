/*
Configuration: manage all about configuration.

-Configuration
*/
package configuration

import (
	"testing"

	"github.com/apenella/messageOutput"
)

type test_data struct {
	file string
	folder string
	values map[string] string
	errors map[string] string
}

var tests = []test_data{
	{ 	
		"",
		"../conf.d",
		map[string]string {"name": "apenella", "checks": "../conf.d/checks", "services": "../conf.d/services"},
		nil,
	},
	{ 	
		"config.json",
		"../conf.d",
		map[string]string {"name": "apenella", "checks": "../conf.d/checks", "services": "../conf.d/services"},
		nil,
	},
} 


func TestNewConfiguration(t *testing.T){
	output := message.GetInstance(3)
	
	for _,test := range tests {
		if err, configuration := NewConfiguration(test.file, test.folder, output); err != nil {
			t.Error(err)
		} else {
			if configuration.Name != test.values["name"] {
				t.Error("(environment::NewConfiguration) Property 'name' not loaded")
			}
			if configuration.Checks.Folder != test.values["checks"] {
				t.Error("(environment::NewConfiguration) Property 'checks' not loaded")
			}
		}
	}
}