package cli

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	CoreV1 "k8s.io/api/core/v1"
	MetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	KRT "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	utiltesting "k8s.io/client-go/util/testing"
)

func testServerEnv(t *testing.T, statusCode int) (*httptest.Server, *utiltesting.FakeHandler, *MetaV1.Status) {
	status := &MetaV1.Status{
		TypeMeta: MetaV1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Status",
		},
		Status: fmt.Sprintf("%s", MetaV1.StatusSuccess),
	}
	expectedBody, _ := runtime.Encode(scheme.Codecs.LegacyCodec(CoreV1.SchemeGroupVersion), status)
	fakeHandler := utiltesting.FakeHandler{
		StatusCode:   statusCode,
		ResponseBody: string(expectedBody),
		T:            t,
	}
	testServer := httptest.NewServer(&fakeHandler)
	return testServer, &fakeHandler, status
}

func testTableServerEnv(t *testing.T, statusCode int) *httptest.Server {
	obj := &MetaV1.Table{
		TypeMeta: MetaV1.TypeMeta{
			APIVersion: "meta.k8s.io/v1",
			Kind:       "Table",
		},
		ColumnDefinitions: []MetaV1.TableColumnDefinition{
			{
				Name:        "SMS",
				Type:        "string",
				Format:      "",
				Description: "Session Management Service Name",
				Priority:    0,
			},
			{
				Name:        "SESS-COUNT",
				Type:        "uint32",
				Format:      "",
				Description: "Number of sessions",
				Priority:    0,
			},
			{
				Name:        "TERMINAL-COUNT",
				Type:        "int64",
				Format:      "",
				Description: "Number of terminals",
				Priority:    0,
			},
		},
		Rows: []MetaV1.TableRow{
			{
				Cells: []any{"sms1", 0, 0, 0},
			},
			{
				Cells: []any{"sms2", 1, 1, 1},
			},
		},
	}
	gv := schema.GroupVersion{Group: "meta.k8s.io", Version: "v1"}
	_scheme := KRT.NewScheme()
	_schemeBuilder := KRT.NewSchemeBuilder(func(scheme *KRT.Scheme) error {
		scheme.AddKnownTypes(gv, &MetaV1.Table{}, &MetaV1.Table{})
		MetaV1.AddToGroupVersion(scheme, gv)
		return nil
	})
	_ = _schemeBuilder.AddToScheme(_scheme)

	expectedBody, err := runtime.Encode(serializer.NewCodecFactory(_scheme).LegacyCodec(gv), obj)
	if err != nil {

	}
	fakeHandler := utiltesting.FakeHandler{
		StatusCode:   statusCode,
		ResponseBody: string(expectedBody),
		T:            t,
	}
	testServer := httptest.NewServer(&fakeHandler)
	return testServer
}

func testJsonServerEnv(t *testing.T, statusCode int) *httptest.Server {
	fakeHandler := utiltesting.FakeHandler{
		StatusCode:   statusCode,
		ResponseBody: "{\"kind\":\"SmsInfoList\",\"apiVersion\":\"toybox.benz.site/v1alpha1\",\"metadata\":{},\"items\":[{\"metadata\":{\"namespace\":\"cluster1\",\"creationTimestamp\":null},\"spec\":{\"smSvc\":\"sms1\"},\"status\":{\"terminalCount\":0}}]}",
		T:            t,
	}
	testServer := httptest.NewServer(&fakeHandler)
	return testServer
}

func testJsonTimeoutServerEnv(t *testing.T, statusCode int) *httptest.Server {
	fakeHandler := utiltesting.FakeHandler{
		StatusCode:   statusCode,
		ResponseBody: "{\"kind\":\"SmsInfoList\",\"apiVersion\":\"toybox.benz.site/v1alpha1\",\"metadata\":{},\"items\":[{\"metadata\":{\"namespace\":\"cluster1\",\"creationTimestamp\":null},\"spec\":{\"smSvc\":\"sms1\"},\"status\":{\"terminalCount\":0}}]}",
		T:            t,
		SkipRequestFn: func(verb string, url url.URL) bool {
			time.Sleep(6 * time.Second)
			return true
		},
	}
	testServer := httptest.NewServer(&fakeHandler)
	return testServer
}

func testNoResultServerEnv(t *testing.T, statusCode int) *httptest.Server {
	fakeHandler := utiltesting.FakeHandler{
		StatusCode:   statusCode,
		ResponseBody: "",
		T:            t,
	}
	testServer := httptest.NewServer(&fakeHandler)
	return testServer
}

func restClient(testServer *httptest.Server) (*rest.RESTClient, error) {
	gv := schema.GroupVersion{Group: "toybox.benz.site", Version: "v1alpha1"}
	c, err := rest.RESTClientFor(&rest.Config{
		Host:    testServer.URL,
		APIPath: "apis",
		ContentConfig: rest.ContentConfig{
			GroupVersion:         &gv,
			NegotiatedSerializer: scheme.Codecs.WithoutConversion(),
		},
		//Username: "user",
		//Password: "pass",
	})
	return c, err
}

func TestCMDParse(t *testing.T) {
	asserter := assert.New(t)

	testServer, _, _ := testServerEnv(t, 200)
	defer testServer.Close()

	c, err := restClient(testServer)
	asserter.NoError(err)

	mg := &sync.Map{}
	mg.Store("smsinfos", "smsinfos")

	commands := []string{"kubectl", "-h"}
	commander := NewDecoratedCommander(commands[0], rest.NewRequest(c).Namespace("default"))
	commander.SetArguments(commands[1:])
	res, err := commander.Exec()
	asserter.Nil(err)
	if len(res) > 0 {
		t.Log(res)
	}

	commands = []string{"kubectl", "describe", "smsinfos"}
	commander = NewDecoratedCommander(commands[0], rest.NewRequest(c).Namespace("default"))
	commander.SetArguments(commands[1:])
	res, err = commander.Exec()
	asserter.Error(err)
	if len(res) > 0 {
		t.Log(res)
	}

	commands = []string{"kubectl", "get", "-h"}
	commander = NewDecoratedCommander(commands[0], rest.NewRequest(c).Namespace("default"))
	commander.SetArguments(commands[1:], &CliKubeResourceMapsDTO{
		GetResourcesMap: mg,
	})
	res, err = commander.Exec()
	asserter.Nil(err)
	if len(res) > 0 {
		t.Log(res)
	}

	commands = []string{"kubectl", "delete", "--help"}
	commander = NewDecoratedCommander(commands[0], rest.NewRequest(c).Namespace("default"))
	commander.SetArguments(commands[1:], &CliKubeResourceMapsDTO{
		DelResourcesMap: mg,
	})
	res, err = commander.Exec()
	asserter.Nil(err)
	if len(res) > 0 {
		t.Log(res)
	}

	commands = []string{"kubectl", "get", ""}
	commander = NewDecoratedCommander(commands[0], rest.NewRequest(c).Namespace("default"))
	commander.SetArguments(commands[1:])
	res, err = commander.Exec()
	asserter.Error(err)
	if len(res) > 0 {
		t.Log(res)
	}

	commands = []string{"kubectl", "get", "--h"}
	commander = NewDecoratedCommander(commands[0], rest.NewRequest(c).Namespace("default"))
	commander.SetArguments(commands[1:])
	res, err = commander.Exec()
	asserter.Error(err)
	if len(res) > 0 {
		t.Log(res)
	}

	commands = []string{"kubectl", "get", "smsinfos"}
	commander = NewDecoratedCommander(commands[0], rest.NewRequest(c).Namespace("default"))
	commander.SetArguments(commands[1:], &CliKubeResourceMapsDTO{
		GetResourcesMap: mg,
	})
	res, err = commander.Exec()
	asserter.Nil(err)
	if len(res) > 0 {
		t.Log(res)
	}
}

func TestOutAsTable(t *testing.T) {
	asserter := assert.New(t)

	testServer := testTableServerEnv(t, 200)
	defer testServer.Close()

	c, err := restClient(testServer)
	asserter.NoError(err)

	mg := &sync.Map{}
	mg.Store("smsinfos", "smsinfos")
	commands := []string{"kubectl", "get", "smsinfos"}
	commander := NewDecoratedCommander(commands[0], rest.NewRequest(c).Namespace("default"))
	commander.SetArguments(commands[1:], &CliKubeResourceMapsDTO{
		GetResourcesMap: mg,
	})
	res, err := commander.Exec()
	asserter.Nil(err)
	if len(res) > 0 {
		t.Log(res)
	}
}

func TestOutAsJSON(t *testing.T) {
	asserter := assert.New(t)

	testServer := testJsonServerEnv(t, 200)
	defer testServer.Close()

	c, err := restClient(testServer)
	asserter.NoError(err)

	mg := &sync.Map{}
	mg.Store("smsinfos", "smsinfos")

	t.Parallel()
	t.Run("out json flag without eq", func(tt *testing.T) {
		commands := []string{"kubectl", "get", "smsinfos", "-o", "json"}
		commander := NewDecoratedCommander(commands[0], rest.NewRequest(c).Namespace("default"))
		commander.SetArguments(commands[1:], &CliKubeResourceMapsDTO{
			GetResourcesMap: mg,
		})
		res, err := commander.Exec()
		asserter.Nil(err)
		if len(res) > 0 {
			tt.Log(res)
		}
	})

	t.Run("out json flag without eq", func(tt *testing.T) {
		commands := []string{"kubectl", "get", "smsinfos", "--output", "json"}
		commander := NewDecoratedCommander(commands[0], rest.NewRequest(c).Namespace("default"))
		commander.SetArguments(commands[1:], &CliKubeResourceMapsDTO{
			GetResourcesMap: mg,
		})
		res, err := commander.Exec()
		asserter.Nil(err)
		if len(res) > 0 {
			tt.Log(res)
		}
	})

	t.Run("out json flag with eq", func(tt *testing.T) {
		commands := []string{"kubectl", "get", "smsinfos", "-o=json"}
		commander := NewDecoratedCommander(commands[0], rest.NewRequest(c).Namespace("default"))
		commander.SetArguments(commands[1:], &CliKubeResourceMapsDTO{
			GetResourcesMap: mg,
		})
		res, err := commander.Exec()
		asserter.Nil(err)
		if len(res) > 0 {
			tt.Log(res)
		}
	})

	t.Run("out json long flag with eq", func(tt *testing.T) {
		commands := []string{"kubectl", "get", "smsinfos", "--output=json"}
		commander := NewDecoratedCommander(commands[0], rest.NewRequest(c).Namespace("default"))
		commander.SetArguments(commands[1:], &CliKubeResourceMapsDTO{
			GetResourcesMap: mg,
		})
		res, err := commander.Exec()
		asserter.Nil(err)
		if len(res) > 0 {
			tt.Log(res)
		}
	})
}

func TestOutAsYAML(t *testing.T) {
	asserter := assert.New(t)

	testServer := testJsonServerEnv(t, 200)
	defer testServer.Close()

	c, err := restClient(testServer)
	asserter.NoError(err)

	mg := &sync.Map{}
	mg.Store("smsinfos", "smsinfos")

	t.Parallel()
	t.Run("out yaml flag without eq", func(tt *testing.T) {
		commands := []string{"kubectl", "get", "smsinfos", "-o", "yaml"}
		commander := NewDecoratedCommander(commands[0], rest.NewRequest(c).Namespace("default"))
		commander.SetArguments(commands[1:], &CliKubeResourceMapsDTO{
			GetResourcesMap: mg,
		})
		res, err := commander.Exec()
		asserter.Nil(err)
		if len(res) > 0 {
			tt.Log(res)
		}
	})

	t.Run("out yaml flag without eq", func(tt *testing.T) {
		commands := []string{"kubectl", "get", "smsinfos", "--output", "yaml"}
		commander := NewDecoratedCommander(commands[0], rest.NewRequest(c).Namespace("default"))
		commander.SetArguments(commands[1:], &CliKubeResourceMapsDTO{
			GetResourcesMap: mg,
		})
		res, err := commander.Exec()
		asserter.Nil(err)
		if len(res) > 0 {
			tt.Log(res)
		}
	})

	t.Run("out yaml flag with eq", func(tt *testing.T) {
		commands := []string{"kubectl", "get", "smsinfos", "-o=yaml"}
		commander := NewDecoratedCommander(commands[0], rest.NewRequest(c).Namespace("default"))
		commander.SetArguments(commands[1:], &CliKubeResourceMapsDTO{
			GetResourcesMap: mg,
		})
		res, err := commander.Exec()
		asserter.Nil(err)
		if len(res) > 0 {
			tt.Log(res)
		}
	})

	t.Run("out json long flag with eq", func(tt *testing.T) {
		commands := []string{"kubectl", "get", "smsinfos", "--output=yaml"}
		commander := NewDecoratedCommander(commands[0], rest.NewRequest(c).Namespace("default"))
		commander.SetArguments(commands[1:], &CliKubeResourceMapsDTO{
			GetResourcesMap: mg,
		})
		res, err := commander.Exec()
		asserter.Nil(err)
		if len(res) > 0 {
			tt.Log(res)
		}
	})
}

func TestOutAsJSONTimeout(t *testing.T) {
	asserter := assert.New(t)

	testServer := testJsonTimeoutServerEnv(t, 200)
	defer testServer.Close()

	c, err := restClient(testServer)
	asserter.NoError(err)

	mg := &sync.Map{}
	mg.Store("smsinfos", "smsinfos")

	commands := []string{"kubectl", "get", "smsinfos", "-o", "json"}
	commander := NewDecoratedCommander(commands[0], rest.NewRequest(c).Namespace("default"))
	commander.SetArguments(commands[1:], &CliKubeResourceMapsDTO{
		GetResourcesMap: mg,
	})
	res, err := commander.Exec()
	asserter.Nil(err)
	if len(res) > 0 {
		t.Log(res)
	}
}

func TestOutAsNoResult(t *testing.T) {
	asserter := assert.New(t)

	testServer := testNoResultServerEnv(t, 200)
	defer testServer.Close()

	c, err := restClient(testServer)
	asserter.NoError(err)

	mg := &sync.Map{}
	mg.Store("smsinfos", "smsinfos")

	commands := []string{"kubectl", "get", "smsinfos", "-o", "json"}
	commander := NewDecoratedCommander(commands[0], rest.NewRequest(c).Namespace("default"))
	commander.SetArguments(commands[1:], &CliKubeResourceMapsDTO{
		GetResourcesMap: mg,
	})
	res, err := commander.Exec()
	asserter.True(errors.Is(err, errCliExecEmptyResult))
	if len(res) > 0 {
		t.Log(res)
	}
}
