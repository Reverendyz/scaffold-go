package create

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type YamlFile struct {
	Apps []struct {
		Name          string `yaml:"name"`
		Watch         bool   `yaml:"watch"`
		Instances     int    `yaml:"instances"`
		Exec_mode     string `yaml:"exec_mode"`
		EnvProduction struct {
			NodeEnv         string `yaml:"NODE_ENV"`
			TypeOrmUsername string `yaml:"TYPEORM_USERNAME"`
			TypeOrmPassword string `yaml:"TYPEORM_PASSWORD"`
			TypeOrmPort     int    `yaml:"TYPEORM_PORT"`
			TypeOrmDatabase string `yaml:"TYPEORM_DATABASE"`
			DatabaseUrl     string `yaml:"DATABASE_URL"`
			SpacesBucket    string `yaml:"SPACES_BUCKET"`
			Port            int    `yaml:"PORT"`
		} `yaml:"env_production"`
	} `yaml:"apps"`
}

var (
	yFile YamlFile
)

func ReadYamlFile(filename string) error {
	err := validateExtension(filename, "yaml")
	if err != nil {
		return err
	}
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(file, &yFile)
	if err != nil {
		return err
	}
	return nil
}

func CreateScaffold(environment string, isProd bool) error {
	file, err := os.ReadFile("scaffold/create/template/template.yaml")
	if err != nil {
		return err
	}

	yaml.Unmarshal(file, &yFile)

	name := yFile.Apps[0].Name

	yFile.Apps[0].Name = strings.Replace(name, "template", environment, 1)
	yFile.Apps[0].EnvProduction.TypeOrmPort = yFile.Apps[0].EnvProduction.TypeOrmPort + 10
	yFile.Apps[0].EnvProduction.Port = yFile.Apps[0].EnvProduction.Port + 1

	databaseUrl := yFile.Apps[0].EnvProduction.DatabaseUrl
	re := regexp.MustCompile(`^postgresql:\/\/(\w+):\w+@(\d+.\d+.\d+.\d+):(\d+)\/(\w+)\?\w+.\w+$`)

	matches := re.FindStringSubmatch(databaseUrl)

	db_user := matches[1]
	db_ip := matches[2]
	db_port, err := strconv.Atoi(matches[3])
	if err != nil {
		return err
	}
	db_port += 10
	db_name := strings.Replace(matches[4], "template", environment, 1)

	dbPassword, err := generateRandomString(12)
	if err != nil {
		return err
	}

	yFile.Apps[0].EnvProduction.TypeOrmPassword = dbPassword

	newDb := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", db_user, dbPassword, db_ip, db_port, db_name)
	yFile.Apps[0].EnvProduction.DatabaseUrl = newDb

	out, err := yaml.Marshal(yFile)
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("pm2.%s.yaml", environment), out, 0644)
	if err != nil {
		return err
	}

	err = updateTempate()
	if err != nil {
		return err
	}

	return nil
}

func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func updateTempate() error {
	filename := "scaffold/create/template/template.yaml"
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	yaml.Unmarshal(file, &yFile)

	yFile.Apps[0].EnvProduction.TypeOrmPort = yFile.Apps[0].EnvProduction.TypeOrmPort + 10
	yFile.Apps[0].EnvProduction.Port = yFile.Apps[0].EnvProduction.Port + 1

	out, err := yaml.Marshal(yFile)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, out, 0644)
	if err != nil {
		return err
	}

	return nil
}

func validateExtension(filename, ext string) error {
	extension := strings.Split(filename, ".")[1]
	if extension != ext {
		return errors.New("Filename does not match " + ext)
	}
	return nil
}
