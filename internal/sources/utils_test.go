package sources

import (
	"testing"

	"github.com/spf13/viper"
)

func Test_validateConfigSet(t *testing.T) {
	type args struct {
		flag    string
		configs []string
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"feature enabled with set configs", args{"true_bool", []string{"set_config", "set_config"}}, true, false},
		{"feature enabled with missing configs", args{"true_bool", []string{"set_config", "empty_config"}}, false, true},
		{"feature enabled with no configs", args{"true_bool", []string{}}, true, false},
		{"feature disabled with set configs", args{"false_bool", []string{"set_config", "set_config"}}, false, false},
		{"feature disabled with missing configs", args{"false_bool", []string{"set_config", "empty_config"}}, false, false},
		{"feature disabled with no configs", args{"false_bool", []string{}}, false, false},
	}

	setupViper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateConfigSet(tt.args.flag, tt.args.configs)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfigSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateConfigSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func setupViper() {
	viper.SetDefault("true_bool", true)
	viper.SetDefault("false_bool", false)
	viper.SetDefault("empty_config", "")
	viper.SetDefault("set_config", "testing")
}
