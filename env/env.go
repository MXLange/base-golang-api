package env

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	e "github.com/MXLange/go-model/internal/errors"
	"github.com/MXLange/go-model/internal/logger"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type Env struct {
	Environment        string `mapstructure:"ENVIRONMENT"`
	Port               string `mapstructure:"PORT"`
	DBDriverName       string `mapstructure:"DB_DRIVER_NAME"`
	DBConnectionString string `mapstructure:"DB_CONNECTION_STRING"`
}

func New(ctx context.Context, logger logger.LoggerIF) (*Env, error) {

	if logger == nil {
		return nil, e.ErrNilLogger
	}

	v := viper.New()

	if err := mergeConfigFile(v, ".env"); err != nil {
		return nil, err
	}

	v.AutomaticEnv()

	if err := validateRequiredKeys(v); err != nil {
		return nil, err
	}

	var cfg Env
	if err := v.Unmarshal(&cfg, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "mapstructure"
	}); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func mergeConfigFile(v *viper.Viper, path string) error {
	v.SetConfigFile(path)
	v.SetConfigType("env")

	if err := v.MergeInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}

	return nil
}

func validateRequiredKeys(v *viper.Viper) error {
	envType := reflect.TypeOf(Env{})
	missing := make([]string, 0, envType.NumField())

	for i := range envType.NumField() {
		field := envType.Field(i)
		key := field.Tag.Get("mapstructure")
		if key == "" {
			continue
		}

		if !v.IsSet(key) {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}

	return nil
}
