/*
 * Copyright 2018 Comcast Cable Communications Management, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package routing

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/tricksterproxy/trickster/cmd/trickster/config"
	"github.com/tricksterproxy/trickster/pkg/backends"
	"github.com/tricksterproxy/trickster/pkg/backends/healthcheck"
	bo "github.com/tricksterproxy/trickster/pkg/backends/options"
	"github.com/tricksterproxy/trickster/pkg/backends/reverseproxycache"
	"github.com/tricksterproxy/trickster/pkg/backends/rule"
	"github.com/tricksterproxy/trickster/pkg/cache/registration"
	tl "github.com/tricksterproxy/trickster/pkg/observability/logging"
	"github.com/tricksterproxy/trickster/pkg/observability/tracing"
	"github.com/tricksterproxy/trickster/pkg/observability/tracing/exporters/zipkin"
	to "github.com/tricksterproxy/trickster/pkg/observability/tracing/options"
	po "github.com/tricksterproxy/trickster/pkg/proxy/paths/options"
	tlstest "github.com/tricksterproxy/trickster/pkg/util/testing/tls"

	"github.com/gorilla/mux"
)

func TestRegisterPprofRoutes(t *testing.T) {
	router := http.NewServeMux()
	log := tl.ConsoleLogger("info")
	RegisterPprofRoutes("test", router, log)
	r, _ := http.NewRequest("GET", "http://0/debug/pprof", nil)
	_, p := router.Handler(r)
	if p != "/debug/pprof/" {
		t.Error("expected pprof route path")
	}
}

func TestRegisterHealthHandler(t *testing.T) {
	router := http.NewServeMux()
	path := "/test"
	hc := healthcheck.New()
	RegisterHealthHandler(router, path, hc)
}

func TestRegisterProxyRoutes(t *testing.T) {

	var proxyClients backends.Backends

	log := tl.ConsoleLogger("info")
	conf, _, err := config.Load("trickster", "test",
		[]string{"-log-level", "debug", "-origin-url", "http://1", "-provider", "prometheus"})
	if err != nil {
		t.Fatalf("Could not load configuration: %s", err.Error())
	}
	caches := registration.LoadCachesFromConfig(conf, tl.ConsoleLogger("error"))
	defer registration.CloseCaches(caches)
	proxyClients, err = RegisterProxyRoutes(conf, mux.NewRouter(), http.NewServeMux(), caches, nil, log, false)
	if err != nil {
		t.Error(err)
	}
	z, err := zipkin.NewTracer(&to.Options{ServiceName: "test", CollectorURL: "http://1.2.3.4/"})
	if err != nil {
		t.Error(err)
	}
	tr := map[string]*tracing.Tracer{"test": z}
	o := conf.Backends["default"]
	o.TracingConfigName = "test"

	o.Hosts = []string{"test", "test2"}

	registration.LoadCachesFromConfig(conf, tl.ConsoleLogger("error"))
	RegisterProxyRoutes(conf, mux.NewRouter(), http.NewServeMux(), caches, tr, log, false)

	if len(proxyClients) == 0 {
		t.Errorf("expected %d got %d", 1, 0)
	}

	conf.Backends["default"] = bo.New()

	// Test Too Many Defaults
	o1 := conf.Backends["default"]
	o2 := bo.New()

	o1.IsDefault = true
	o2.IsDefault = true

	o1.Provider = "rpc"
	o2.Provider = "rpc"

	conf.Backends["2"] = o2

	router := mux.NewRouter()
	_, err = RegisterProxyRoutes(conf, router, http.NewServeMux(), caches, tr, log, false)
	if err == nil {
		t.Error("Expected error for too many default backends.")
	}

	o1.IsDefault = false
	o1.CacheName = "invalid"
	_, err = RegisterProxyRoutes(conf, router, http.NewServeMux(), caches, tr, log, false)
	if err == nil {
		t.Errorf("Expected error for invalid cache name")
	}

	o1.CacheName = o2.CacheName
	_, err = RegisterProxyRoutes(conf, router, http.NewServeMux(), caches, tr, log, false)
	if err != nil {
		t.Error(err)
	}

	o2.IsDefault = false
	o2.CacheName = "invalid"
	_, err = RegisterProxyRoutes(conf, router, http.NewServeMux(), caches, tr, log, false)
	if err == nil {
		t.Errorf("Expected error for invalid cache name")
	}

	o2.CacheName = "default"
	_, err = RegisterProxyRoutes(conf, router, http.NewServeMux(), caches, tr, log, false)
	if err != nil {
		t.Error(err)
	}

	// test the condition where no backends are IsDefault true,
	// and no backends are named default

	o1.IsDefault = false
	o2.IsDefault = false
	conf.Backends["1"] = o1
	delete(conf.Backends, "default")

	o1.Paths["/-0000000011"].Methods = nil

	_, err = RegisterProxyRoutes(conf, router, http.NewServeMux(), caches, tr, log, false)
	if err != nil {
		t.Error(err)
	}

}

func TestRegisterProxyRoutesInflux(t *testing.T) {
	conf, _, err := config.Load("trickster", "test",
		[]string{"-log-level", "debug", "-origin-url", "http://1", "-provider", "influxdb"})
	if err != nil {
		t.Fatalf("Could not load configuration: %s", err.Error())
	}

	caches := registration.LoadCachesFromConfig(conf, tl.ConsoleLogger("error"))
	defer registration.CloseCaches(caches)
	proxyClients, err := RegisterProxyRoutes(conf, mux.NewRouter(), http.NewServeMux(), caches,
		nil, tl.ConsoleLogger("info"), false)
	if err != nil {
		t.Error(err)
	}

	if len(proxyClients) == 0 {
		t.Errorf("expected %d got %d", 1, 0)
	}

}

func TestRegisterProxyRoutesReverseProxy(t *testing.T) {
	conf, _, err := config.Load("trickster", "test",
		[]string{"-log-level", "debug", "-origin-url", "http://1", "-provider", "rp"})
	if err != nil {
		t.Fatalf("Could not load configuration: %s", err.Error())
	}

	caches := registration.LoadCachesFromConfig(conf, tl.ConsoleLogger("error"))
	defer registration.CloseCaches(caches)
	proxyClients, err := RegisterProxyRoutes(conf, mux.NewRouter(), http.NewServeMux(), caches,
		nil, tl.ConsoleLogger("info"), false)
	if err != nil {
		t.Error(err)
	}

	if len(proxyClients) == 0 {
		t.Errorf("expected %d got %d", 1, 0)
	}

}

func TestRegisterProxyRoutesClickHouse(t *testing.T) {

	conf, _, err := config.Load("trickster", "test",
		[]string{"-log-level", "debug", "-origin-url", "http://1", "-provider", "clickhouse"})
	if err != nil {
		t.Fatalf("Could not load configuration: %s", err.Error())
	}

	caches := registration.LoadCachesFromConfig(conf, tl.ConsoleLogger("error"))
	defer registration.CloseCaches(caches)
	proxyClients, err := RegisterProxyRoutes(conf, mux.NewRouter(), http.NewServeMux(), caches,
		nil, tl.ConsoleLogger("info"), false)
	if err != nil {
		t.Error(err)
	}

	if len(proxyClients) == 0 {
		t.Errorf("expected %d got %d", 1, 0)
	}

}

func TestRegisterProxyRoutesIRONdb(t *testing.T) {

	conf, _, err := config.Load("trickster", "test",
		[]string{"-origin-url", "http://example.com", "-provider", "irondb", "-log-level", "debug"})
	if err != nil {
		t.Fatalf("Could not load configuration: %s", err.Error())
	}

	caches := registration.LoadCachesFromConfig(conf, tl.ConsoleLogger("error"))
	defer registration.CloseCaches(caches)
	proxyClients, err := RegisterProxyRoutes(conf, mux.NewRouter(), http.NewServeMux(), caches,
		nil, tl.ConsoleLogger("info"), false)
	if err != nil {
		t.Error(err)
	}

	if len(proxyClients) == 0 {
		t.Errorf("expected %d got %d", 1, 0)
	}
}

func TestRegisterProxyRoutesWithReqRewriters(t *testing.T) {

	conf, _, err := config.Load("trickster", "test",
		[]string{"-config", "../../testdata/test.routing.req_rewriter.conf"})
	if err != nil {
		t.Fatalf("Could not load configuration: %s", err.Error())
	}

	tpo := po.New()
	tpo.ReqRewriterName = "path"
	conf.Backends["test"].Paths["test"] = tpo

	caches := registration.LoadCachesFromConfig(conf, tl.ConsoleLogger("error"))
	defer registration.CloseCaches(caches)
	proxyClients, err := RegisterProxyRoutes(conf, mux.NewRouter(), http.NewServeMux(), caches,
		nil, tl.ConsoleLogger("info"), false)
	if err != nil {
		t.Error(err)
	}

	if len(proxyClients) != 2 {
		t.Errorf("expected %d got %d", 1, len(proxyClients))
	}
}

func TestRegisterProxyRoutesMultipleDefaults(t *testing.T) {
	expected1 := "only one backend can be marked as default. Found both test and test2"
	expected2 := "only one backend can be marked as default. Found both test2 and test"

	a := []string{"-config", "../../testdata/test.too_many_defaults.conf"}
	conf, _, err := config.Load("trickster", "test", a)
	if err != nil {
		t.Fatalf("Could not load configuration: %s", err.Error())
	}
	caches := registration.LoadCachesFromConfig(conf, tl.ConsoleLogger("error"))
	defer registration.CloseCaches(caches)
	_, err = RegisterProxyRoutes(conf, mux.NewRouter(), http.NewServeMux(), caches,
		nil, tl.ConsoleLogger("info"), false)
	if err == nil {
		t.Errorf("expected error `%s` got nothing", expected1)
	} else if err.Error() != expected1 && err.Error() != expected2 {
		t.Errorf("expected error `%s` got `%s`", expected1, err.Error())
	}
}

func TestRegisterProxyRoutesInvalidCert(t *testing.T) {
	expected := "tls: failed to find any PEM data in certificate input"

	kb, _, _ := tlstest.GetTestKeyAndCert(false)

	td := t.TempDir()

	certfile := td + "/cert.pem"
	keyfile := td + "/key.pem"
	confFile := td + "/trickster_test_config.conf"

	err := os.WriteFile(certfile, []byte{}, 0600)
	if err != nil {
		t.Error(err)
	}
	err = os.WriteFile(keyfile, kb, 0600)
	if err != nil {
		t.Error(err)
	}

	b, err := os.ReadFile("../../testdata/test.bad_tls_cert.routes.conf")
	b = []byte(strings.ReplaceAll(string(b), `../../testdata/test.06.`, td+"/"))

	err = os.WriteFile(confFile, b, 0600)
	if err != nil {
		t.Error(err)
	}

	a := []string{"-config", confFile}
	conf, _, err := config.Load("trickster", "test", a)
	if err != nil {
		t.Fatalf("Could not load configuration: %s", err.Error())
	}
	caches := registration.LoadCachesFromConfig(conf, tl.ConsoleLogger("error"))
	defer registration.CloseCaches(caches)
	_, err = RegisterProxyRoutes(conf, mux.NewRouter(), http.NewServeMux(), caches,
		nil, tl.ConsoleLogger("info"), false)
	if err == nil {
		t.Errorf("expected error: %s", expected)
	}
	if err != nil && err.Error() != expected {
		t.Errorf("expected error: %s, got: %s", expected, err.Error())
	}
}

func TestRegisterProxyRoutesBadCacheName(t *testing.T) {
	expected := `invalid cache name "test2" provided in backend options "test"`
	a := []string{"-config", "../../testdata/test.bad_cache_name.conf"}
	_, _, err := config.Load("trickster", "test", a)
	if err == nil {
		t.Errorf("expected error `%s` got nothing", expected)
	} else if err.Error() != expected {
		t.Errorf("expected error `%s` got `%s`", expected, err.Error())
	}
}

func TestRegisterProxyRoutesBadProvider(t *testing.T) {
	expected := "unknown backend provider in backend options. backendName: test, backendProvider: foo"
	a := []string{"-config", "../../testdata/test.unknown_backend_provider.conf"}
	conf, _, err := config.Load("trickster", "test", a)
	if err != nil {
		t.Fatalf("Could not load configuration: %s", err.Error())
	}
	caches := registration.LoadCachesFromConfig(conf, tl.ConsoleLogger("error"))
	defer registration.CloseCaches(caches)
	_, err = RegisterProxyRoutes(conf, mux.NewRouter(), http.NewServeMux(), caches,
		nil, tl.ConsoleLogger("info"), false)
	if err == nil {
		t.Errorf("expected error `%s` got nothing", expected)
	} else if err.Error() != expected {
		t.Errorf("expected error `%s` got `%s`", expected, err.Error())
	}
}

func TestRegisterMultipleBackends(t *testing.T) {
	a := []string{"-config", "../../testdata/test.multiple_backends.conf"}
	conf, _, err := config.Load("trickster", "test", a)
	if err != nil {
		t.Fatalf("Could not load configuration: %s", err.Error())
	}
	caches := registration.LoadCachesFromConfig(conf, tl.ConsoleLogger("error"))
	defer registration.CloseCaches(caches)
	_, err = RegisterProxyRoutes(conf, mux.NewRouter(), http.NewServeMux(), caches,
		nil, tl.ConsoleLogger("info"), false)
	if err != nil {
		t.Error(err)
	}
}

func TestRegisterMultipleBackendsPlusDefault(t *testing.T) {
	a := []string{"-config", "../../testdata/test.multiple_backends_plus_default.conf"}
	conf, _, err := config.Load("trickster", "test", a)
	if err != nil {
		t.Fatalf("Could not load configuration: %s", err.Error())
	}
	caches := registration.LoadCachesFromConfig(conf, tl.ConsoleLogger("error"))
	defer registration.CloseCaches(caches)
	_, err = RegisterProxyRoutes(conf, mux.NewRouter(), http.NewServeMux(), caches,
		nil, tl.ConsoleLogger("info"), false)
	if err != nil {
		t.Error(err)
	}
	if !conf.Backends["default"].IsDefault {
		t.Errorf("expected backend %s.IsDefault to be true", "default")
	}
}

func TestRegisterPathRoutes(t *testing.T) {
	p := map[string]*po.Options{"test": {}}
	RegisterPathRoutes(nil, nil, nil, nil, nil, p, nil, "", nil)

	conf, _, err := config.Load("trickster", "test",
		[]string{"-log-level", "debug", "-origin-url", "http://1", "-provider", "rpc"})
	if err != nil {
		t.Fatalf("Could not load configuration: %s", err.Error())
	}

	oo := conf.Backends["default"]
	rpc, _ := reverseproxycache.NewClient("test", oo, mux.NewRouter(), nil)
	dpc := rpc.DefaultPathConfigs(oo)
	dpc["/-GET-HEAD"].Methods = nil
	RegisterPathRoutes(nil, nil, rpc, oo, nil, dpc, nil, "", tl.ConsoleLogger("INFO"))

}

func TestValidateRuleClients(t *testing.T) {

	c, err := rule.NewClient("test", nil, nil, nil)
	if err != nil {
		t.Error(err)
	}

	var cl = backends.Backends{"test": c}
	rule.ValidateOptions(cl, nil)

	conf, _, err := config.Load("trickster", "test",
		[]string{"-log-level", "debug", "-origin-url", "http://1", "-provider", "rpc"})
	if err != nil {
		t.Fatalf("Could not load configuration: %s", err.Error())
	}

	caches := registration.LoadCachesFromConfig(conf, tl.ConsoleLogger("error"))
	defer registration.CloseCaches(caches)

	o := conf.Backends["default"]
	o.Provider = "rule"

	_, err = RegisterProxyRoutes(conf, mux.NewRouter(), http.NewServeMux(), caches,
		nil, tl.ConsoleLogger("info"), false)
	if err == nil {
		t.Error("expected error")
	}

}
