package pipelines_test

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/grafana/grafana-build/pipelines"
)

type TestCLIContext struct {
	Data map[string]interface{}
}

func (t *TestCLIContext) Bool(key string) bool {
	return t.Data[key].(bool)
}

func (t *TestCLIContext) String(key string) string {
	log.Println("getting key", key)
	return t.Data[key].(string)
}

func (t *TestCLIContext) Set(key string, val string) error {
	t.Data[key] = val

	return nil
}

func (t *TestCLIContext) StringSlice(key string) []string {
	return t.Data[key].([]string)
}

func (t *TestCLIContext) Path(key string) string {
	return t.Data[key].(string)
}

func TestPipelineArgsFromContext(t *testing.T) {
	enterpriseDir, err := os.MkdirTemp("", "grafana-enterprise-*")
	if err != nil {
		t.Fatal(err)
	}

	validData := map[string]interface{}{
		"v":              true,
		"version":        "v1.0.0",
		"grafana":        true,
		"grafana-dir":    "/grafana",
		"grafana-ref":    "asdf",
		"enterprise":     true,
		"enterprise-dir": enterpriseDir,
		"enterprise-ref": "1234",
		"build-id":       "build-1234",
		"github-token":   "",
	}

	t.Run("It should return a PipelineArgs object if there are no errors", func(t *testing.T) {
		args, err := pipelines.PipelineArgsFromContext(context.Background(), &TestCLIContext{
			Data: validData,
		})
		if err != nil {
			t.Fatal(err)
		}

		if args.Verbose != true {
			t.Error("args.Verbose should be true")
		}

		if args.Version != "v1.0.0" {
			t.Error("args.Version should be v1.0.0")
		}

		if args.BuildGrafana != true {
			t.Error("args.BuildGrafana should be true")
		}

		if args.GrafanaDir != "/grafana" {
			t.Error("args.GrafanaDir should be /grafana")
		}

		if args.GrafanaRef != "asdf" {
			t.Error("args.GrafanaRef should be asdf")
		}

		if args.BuildEnterprise != true {
			t.Error("args.Enterprise should be true")
		}

		if args.EnterpriseDir != enterpriseDir {
			t.Errorf("args.EnterpriseDir should be %s", enterpriseDir)
		}

		if args.EnterpriseRef != "1234" {
			t.Error("args.EnterpriseRef should be 1234")
		}
	})

	t.Run("If no build ID is provided, a random 12-character string should be given", func(t *testing.T) {
		data := validData
		data["build-id"] = ""
		args, err := pipelines.PipelineArgsFromContext(context.Background(), &TestCLIContext{
			Data: data,
		})
		if err != nil {
			t.Fatal(err)
		}

		if args.BuildID == "" {
			t.Fatal("BuildID should not be empty")
		}
		if len(args.BuildID) != 12 {
			t.Fatal("BuildID should be a 12-character string")
		}
	})

	t.Run("If the --enterprise-ref is set to a non-default value, it should set the enterprise flag to true", func(t *testing.T) {
		data := validData
		data["enterprise"] = false
		data["enterprise-ref"] = "ref-1234"

		args, err := pipelines.PipelineArgsFromContext(context.Background(), &TestCLIContext{
			Data: data,
		})
		if err != nil {
			t.Fatal(err)
		}

		if args.BuildEnterprise != true {
			t.Fatal("args.BuildEnterprise should be true")
		}
	})

	t.Run("If the --enterprise-ref is set to a non-default value, it should set the enterprise flag to true", func(t *testing.T) {
		data := validData
		data["enterprise"] = false
		data["enterprise-ref"] = ""
		data["enterprise-dir"] = filepath.Join(enterpriseDir, "does-not-exist")

		_, err := pipelines.PipelineArgsFromContext(context.Background(), &TestCLIContext{
			Data: data,
		})
		if err == nil {
			t.Fatal("error should not be empty")
		}
	})
}
