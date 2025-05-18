// pkg/authz/enforcer.go
package authz

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/walterfan/prompt-service/pkg/database"
)

var Enforcer *casbin.Enforcer

// refer to https://github.com/casbin/gorm-adapter , the table named casbin_rules
// Gorm Adapter is the Gorm adapter for Casbin. With this library,
// Casbin can load policy from Gorm supported database or save policy to it.
// The adapter will use the table named "casbin_rule".
// If it doesn't exist, the adapter will create it automatically
func must(err error) {
	if err != nil {
		panic(err)
	}
}
func InitAuthz(modelConfigfile string) {
	adapter, err := gormadapter.NewAdapterByDB(database.DB)
	must(err)

	enforcer, err := casbin.NewEnforcer(modelConfigfile, adapter)
	must(err)

	err = enforcer.LoadPolicy()
	must(err)

	Enforcer = enforcer
	policies, err := Enforcer.GetPolicy()
	must(err)

	if len(policies) == 0 {
		_, err := Enforcer.AddPolicy("admin", "/api/v1/prompts/*", "*")
		must(err)
	}
}

func loadPoliciesFromCSV(enforcer *casbin.Enforcer, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) < 4 || strings.ToLower(strings.TrimSpace(parts[0])) != "p" {
			continue // only process 'p' type policies
		}

		subject := strings.TrimSpace(parts[1])
		object := strings.TrimSpace(parts[2])
		action := strings.TrimSpace(parts[3])

		_, err := enforcer.AddPolicy(subject, object, action)
		if err != nil {
			panic(err)
		}
	}
}
