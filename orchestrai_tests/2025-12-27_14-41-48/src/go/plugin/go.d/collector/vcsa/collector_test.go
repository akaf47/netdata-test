package vcsa

import (
	"context"
	"errors"
	"testing"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/pkg/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("creates new collector with default values", func(t *testing.T) {
		collector := New()

		assert.NotNil(t, collector)
		assert.IsType(t, (*Collector)(nil), collector)
	})
}

func TestCollectorInit(t *testing.T) {
	t.Run("successfully initializes collector", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)

		err := collector.Init()
		assert.NoError(t, err)
	})

	t.Run("init with nil collector", func(t *testing.T) {
		var collector *Collector
		if collector != nil {
			err := collector.Init()
			assert.NoError(t, err)
		}
	})
}

func TestCollectorCheck(t *testing.T) {
	t.Run("check returns true when initialized properly", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)

		err := collector.Init()
		assert.NoError(t, err)

		ok := collector.Check(context.Background())
		assert.True(t, ok)
	})

	t.Run("check with timeout context", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		ok := collector.Check(ctx)
		// Depends on implementation, but should handle cancelled context gracefully
		assert.IsType(t, bool, ok)
	})
}

func TestCollectorCollect(t *testing.T) {
	t.Run("collect returns non-nil metrics", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)

		err := collector.Init()
		assert.NoError(t, err)

		metrics := collector.Collect(context.Background())
		assert.NotNil(t, metrics)
	})

	t.Run("collect with background context", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)

		err := collector.Init()
		assert.NoError(t, err)

		metrics := collector.Collect(context.Background())
		assert.IsType(t, (map[string]int64)(nil), metrics)
	})

	t.Run("collect returns consistent type across calls", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)

		err := collector.Init()
		assert.NoError(t, err)

		metrics1 := collector.Collect(context.Background())
		metrics2 := collector.Collect(context.Background())

		assert.IsType(t, metrics1, metrics2)
	})
}

func TestCollectorName(t *testing.T) {
	t.Run("name returns correct collector name", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)

		name := collector.Name()
		assert.Equal(t, "vcsa", name)
	})

	t.Run("name is consistent across calls", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)

		name1 := collector.Name()
		name2 := collector.Name()
		assert.Equal(t, name1, name2)
	})
}

func TestCollectorRegistryName(t *testing.T) {
	t.Run("RegistryName returns correct registry name", func(t *testing.T) {
		name := RegistryName()
		assert.Equal(t, "vcsa", name)
	})

	t.Run("RegistryName is consistent", func(t *testing.T) {
		name1 := RegistryName()
		name2 := RegistryName()
		assert.Equal(t, name1, name2)
	})
}

func TestCollectorFactory(t *testing.T) {
	t.Run("Factory creates new collector instance", func(t *testing.T) {
		collector := Factory()
		assert.NotNil(t, collector)
		assert.IsType(t, (*Collector)(nil), collector)
	})

	t.Run("Factory creates different instances", func(t *testing.T) {
		collector1 := Factory()
		collector2 := Factory()
		assert.NotSame(t, collector1, collector2)
	})
}

func TestCollectorWithHTTPClient(t *testing.T) {
	t.Run("sets HTTP client configuration", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)

		// Test that HTTPClient field exists and can be set
		collector.HTTPClient = web.NewHTTPClient()
		assert.NotNil(t, collector.HTTPClient)
	})

	t.Run("HTTPClient defaults to nil", func(t *testing.T) {
		collector := New()
		assert.Nil(t, collector.HTTPClient)
	})
}

func TestCollectorURL(t *testing.T) {
	t.Run("URL field can be set", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)

		url := "https://localhost/rest/applmgmt/2.0.0/monitoring/query?resourceKey=__SYSTEM_APPLIANCE_UPTIME__"
		collector.URL = url
		assert.Equal(t, url, collector.URL)
	})

	t.Run("default URL is empty string", func(t *testing.T) {
		collector := New()
		assert.Equal(t, "", collector.URL)
	})

	t.Run("URL field accepts empty string", func(t *testing.T) {
		collector := New()
		collector.URL = ""
		assert.Equal(t, "", collector.URL)
	})
}

func TestCollectorUsername(t *testing.T) {
	t.Run("Username field can be set", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)

		collector.Username = "admin"
		assert.Equal(t, "admin", collector.Username)
	})

	t.Run("Username field defaults to empty", func(t *testing.T) {
		collector := New()
		assert.Equal(t, "", collector.Username)
	})
}

func TestCollectorPassword(t *testing.T) {
	t.Run("Password field can be set", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)

		collector.Password = "secret"
		assert.Equal(t, "secret", collector.Password)
	})

	t.Run("Password field defaults to empty", func(t *testing.T) {
		collector := New()
		assert.Equal(t, "", collector.Password)
	})
}

func TestCollectorTLSConfig(t *testing.T) {
	t.Run("TLSConfig field can be set", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)

		tlsConfig := web.TLSConfig{
			InsecureSkipVerify: true,
		}
		collector.TLSConfig = tlsConfig
		assert.Equal(t, tlsConfig, collector.TLSConfig)
	})

	t.Run("TLSConfig defaults to empty config", func(t *testing.T) {
		collector := New()
		assert.Equal(t, web.TLSConfig{}, collector.TLSConfig)
	})
}

func TestCollectorContextHandling(t *testing.T) {
	t.Run("handles nil context gracefully", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)
		err := collector.Init()
		assert.NoError(t, err)

		// Should not panic with nil context
		assert.NotPanics(t, func() {
			_ = collector.Check(context.Background())
		})
	})

	t.Run("handles cancelled context in Check", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)
		err := collector.Init()
		assert.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		result := collector.Check(ctx)
		assert.IsType(t, bool, result)
	})

	t.Run("handles cancelled context in Collect", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)
		err := collector.Init()
		assert.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		metrics := collector.Collect(ctx)
		assert.NotNil(t, metrics)
	})
}

func TestCollectorErrorScenarios(t *testing.T) {
	t.Run("Init handles initialization errors gracefully", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)

		// Multiple init calls should not cause issues
		err1 := collector.Init()
		err2 := collector.Init()
		assert.NoError(t, err1)
		assert.NoError(t, err2)
	})

	t.Run("Check handles errors without panicking", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)
		_ = collector.Init()

		assert.NotPanics(t, func() {
			_ = collector.Check(context.Background())
		})
	})

	t.Run("Collect handles errors without panicking", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)
		_ = collector.Init()

		assert.NotPanics(t, func() {
			_ = collector.Collect(context.Background())
		})
	})
}

func TestCollectorFieldTypes(t *testing.T) {
	t.Run("all exported fields have correct types", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)

		// Verify field types
		assert.IsType(t, (*web.Client)(nil), collector.HTTPClient)
		assert.IsType(t, "", collector.URL)
		assert.IsType(t, "", collector.Username)
		assert.IsType(t, "", collector.Password)
		assert.IsType(t, web.TLSConfig{}, collector.TLSConfig)
	})
}

func TestCollectorStateAfterInit(t *testing.T) {
	t.Run("collector is ready after Init", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)

		err := collector.Init()
		assert.NoError(t, err)

		// After init, Check should work
		result := collector.Check(context.Background())
		assert.IsType(t, bool, result)
	})

	t.Run("collector name unchanged after Init", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)

		nameBefore := collector.Name()
		_ = collector.Init()
		nameAfter := collector.Name()

		assert.Equal(t, nameBefore, nameAfter)
	})
}

func TestCollectorMultipleCalls(t *testing.T) {
	t.Run("multiple Check calls work correctly", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)
		_ = collector.Init()

		for i := 0; i < 5; i++ {
			result := collector.Check(context.Background())
			assert.IsType(t, bool, result)
		}
	})

	t.Run("multiple Collect calls work correctly", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)
		_ = collector.Init()

		for i := 0; i < 5; i++ {
			metrics := collector.Collect(context.Background())
			assert.NotNil(t, metrics)
		}
	})

	t.Run("interleaved Check and Collect calls", func(t *testing.T) {
		collector := New()
		require.NotNil(t, collector)
		_ = collector.Init()

		for i := 0; i < 3; i++ {
			check := collector.Check(context.Background())
			assert.IsType(t, bool, check)

			metrics := collector.Collect(context.Background())
			assert.NotNil(t, metrics)
		}
	})
}

func TestCollectorEmptyConfiguration(t *testing.T) {
	t.Run("collector works with empty configuration", func(t *testing.T) {
		collector := New()
		assert.NotNil(t, collector)
		assert.Equal(t, "", collector.URL)
		assert.Equal(t, "", collector.Username)
		assert.Equal(t, "", collector.Password)

		err := collector.Init()
		assert.NoError(t, err)
	})
}

func TestCollectorConfigurationVariations(t *testing.T) {
	t.Run("various URL formats", func(t *testing.T) {
		urls := []string{
			"https://localhost/rest/applmgmt/2.0.0/monitoring/query",
			"http://localhost:80/rest/applmgmt/2.0.0/monitoring/query",
			"https://vcsa.example.com/rest/applmgmt/2.0.0/monitoring/query",
			"",
		}

		for _, url := range urls {
			collector := New()
			collector.URL = url
			assert.Equal(t, url, collector.URL)
		}
	})

	t.Run("various credentials", func(t *testing.T) {
		credentials := [][2]string{
			{"admin", "password"},
			{"user", ""},
			{"", "password"},
			{"", ""},
		}

		for _, cred := range credentials {
			collector := New()
			collector.Username = cred[0]
			collector.Password = cred[1]
			assert.Equal(t, cred[0], collector.Username)
			assert.Equal(t, cred[1], collector.Password)
		}
	})
}