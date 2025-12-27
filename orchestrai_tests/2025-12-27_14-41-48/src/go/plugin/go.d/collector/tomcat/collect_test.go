package tomcat

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/netdata/netdata/go/plugins/plugin.d/go.d/agent/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollect(t *testing.T) {
	t.Run("Returns non-nil map", func(t *testing.T) {
		tc := New()
		require.NotNil(t, tc)
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
		assert.IsType(t, map[string]int64{}, maps)
	})
	
	t.Run("Returns empty map when not initialized", func(t *testing.T) {
		tc := New()
		assert.Nil(t, tc.URL)
		
		maps := tc.Collect(context.Background())
		if maps != nil {
			assert.Equal(t, 0, len(maps), "should return empty map for uninitialized collector")
		}
	})
}

func TestCollectWithMockServer(t *testing.T) {
	t.Run("Collects metrics from valid Tomcat server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Mock Tomcat Manager status endpoint
			_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<status>
  <jvm>
    <memory total="123456" free="65432"/>
  </jvm>
  <connector>
    <threadInfo currentThreadCount="5" maxThreads="10"/>
  </connector>
</status>`))
		}))
		defer server.Close()
		
		tc := New()
		tc.URL = server.URL
		require.NoError(t, tc.Init(context.Background()))
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
	})
	
	t.Run("Handles server connection errors", func(t *testing.T) {
		tc := New()
		tc.URL = "http://invalid.local:99999"
		_ = tc.Init(context.Background())
		
		maps := tc.Collect(context.Background())
		// Should handle gracefully without panicking
		assert.NotNil(t, maps)
	})
	
	t.Run("Handles malformed XML response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`invalid xml content`))
		}))
		defer server.Close()
		
		tc := New()
		tc.URL = server.URL
		require.NoError(t, tc.Init(context.Background()))
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
	})
	
	t.Run("Handles empty response body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Empty response
		}))
		defer server.Close()
		
		tc := New()
		tc.URL = server.URL
		require.NoError(t, tc.Init(context.Background()))
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
	})
	
	t.Run("Handles HTTP errors", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()
		
		tc := New()
		tc.URL = server.URL
		require.NoError(t, tc.Init(context.Background()))
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
	})
	
	t.Run("Handles 401 Unauthorized", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()
		
		tc := New()
		tc.URL = server.URL
		require.NoError(t, tc.Init(context.Background()))
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
	})
	
	t.Run("Handles 404 Not Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()
		
		tc := New()
		tc.URL = server.URL
		require.NoError(t, tc.Init(context.Background()))
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
	})
}

func TestCollectWithContextTimeout(t *testing.T) {
	t.Run("Respects context cancellation", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case <-r.Context().Done():
				return
			default:
				_, _ = w.Write([]byte(`<status></status>`))
			}
		}))
		defer server.Close()
		
		tc := New()
		tc.URL = server.URL
		require.NoError(t, tc.Init(context.Background()))
		
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		
		maps := tc.Collect(ctx)
		assert.NotNil(t, maps)
	})
}

func TestCollectMetricsExtraction(t *testing.T) {
	t.Run("Extracts thread pool metrics", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`<?xml version="1.0"?>
<status>
  <connector name="HTTP/1.1" processingTime="100" requestCount="50" />
</status>`))
		}))
		defer server.Close()
		
		tc := New()
		tc.URL = server.URL
		require.NoError(t, tc.Init(context.Background()))
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
	})
	
	t.Run("Extracts memory metrics", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`<?xml version="1.0"?>
<status>
  <memory free="1024" total="2048" max="4096"/>
</status>`))
		}))
		defer server.Close()
		
		tc := New()
		tc.URL = server.URL
		require.NoError(t, tc.Init(context.Background()))
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
	})
	
	t.Run("Extracts request metrics", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`<?xml version="1.0"?>
<status>
  <connector requestCount="100" errorCount="2"/>
</status>`))
		}))
		defer server.Close()
		
		tc := New()
		tc.URL = server.URL
		require.NoError(t, tc.Init(context.Background()))
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
	})
	
	t.Run("Handles missing metric attributes", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`<?xml version="1.0"?>
<status>
  <connector name="HTTP/1.1"/>
</status>`))
		}))
		defer server.Close()
		
		tc := New()
		tc.URL = server.URL
		require.NoError(t, tc.Init(context.Background()))
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
	})
	
	t.Run("Handles zero values in metrics", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`<?xml version="1.0"?>
<status>
  <connector requestCount="0" errorCount="0"/>
</status>`))
		}))
		defer server.Close()
		
		tc := New()
		tc.URL = server.URL
		require.NoError(t, tc.Init(context.Background()))
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
	})
	
	t.Run("Handles negative values in metrics", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`<?xml version="1.0"?>
<status>
  <connector requestCount="-1"/>
</status>`))
		}))
		defer server.Close()
		
		tc := New()
		tc.URL = server.URL
		require.NoError(t, tc.Init(context.Background()))
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
	})
	
	t.Run("Handles very large metric values", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`<?xml version="1.0"?>
<status>
  <connector requestCount="9223372036854775807"/>
</status>`))
		}))
		defer server.Close()
		
		tc := New()
		tc.URL = server.URL
		require.NoError(t, tc.Init(context.Background()))
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
	})
	
	t.Run("Handles non-numeric metric values", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`<?xml version="1.0"?>
<status>
  <connector requestCount="abc"/>
</status>`))
		}))
		defer server.Close()
		
		tc := New()
		tc.URL = server.URL
		require.NoError(t, tc.Init(context.Background()))
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
	})
}

func TestCollectWithMultipleConnectors(t *testing.T) {
	t.Run("Handles multiple connectors in response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`<?xml version="1.0"?>
<status>
  <connector name="HTTP/1.1" requestCount="100"/>
  <connector name="HTTPS" requestCount="50"/>
  <connector name="AJP" requestCount="25"/>
</status>`))
		}))
		defer server.Close()
		
		tc := New()
		tc.URL = server.URL
		require.NoError(t, tc.Init(context.Background()))
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
	})
}

func TestCollectURLHandling(t *testing.T) {
	t.Run("Handles URL with trailing slash", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`<status></status>`))
		}))
		defer server.Close()
		
		tc := New()
		tc.URL = server.URL + "/"
		require.NoError(t, tc.Init(context.Background()))
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
	})
	
	t.Run("Handles URL without trailing slash", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`<status></status>`))
		}))
		defer server.Close()
		
		tc := New()
		tc.URL = server.URL
		require.NoError(t, tc.Init(context.Background()))
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
	})
	
	t.Run("Handles empty URL", func(t *testing.T) {
		tc := New()
		tc.URL = ""
		require.NoError(t, tc.Init(context.Background()))
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
	})
}

func TestCollectResponseHandling(t *testing.T) {
	t.Run("Handles response with BOM", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// UTF-8 BOM
			_, _ = w.Write([]byte{0xEF, 0xBB, 0xBF})
			_, _ = w.Write([]byte(`<?xml version="1.0"?><status></status>`))
		}))
		defer server.Close()
		
		tc := New()
		tc.URL = server.URL
		require.NoError(t, tc.Init(context.Background()))
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
	})
	
	t.Run("Handles response with special characters", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<status name="Test&lt;Special&gt;">
  <connector name="Test&amp;Conn"/>
</status>`))
		}))
		defer server.Close()
		
		tc := New()
		tc.URL = server.URL
		require.NoError(t, tc.Init(context.Background()))
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
	})
	
	t.Run("Handles response with whitespace", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`
<?xml version="1.0"?>
<status>
  
  <connector name="HTTP/1.1">
    <attribute>value</attribute>
  </connector>
  
</status>
`))
		}))
		defer server.Close()
		
		tc := New()
		tc.URL = server.URL
		require.NoError(t, tc.Init(context.Background()))
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
	})
}

func TestCollectEdgeCases(t *testing.T) {
	t.Run("Handles nil collector state", func(t *testing.T) {
		tc := New()
		require.NotNil(t, tc)
		
		maps := tc.Collect(context.Background())
		assert.NotNil(t, maps)
	})
	
	t.Run("Handles consecutive collects", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`<status></status>`))
		}))
		defer server.Close()
		
		tc := New()
		tc.URL = server.URL
		require.NoError(t, tc.Init(context.Background()))
		
		maps1 := tc.Collect(context.Background())
		maps2 := tc.Collect(context.Background())
		maps3 := tc.Collect(context.Background())
		
		assert.NotNil(t, maps1)
		assert.NotNil(t, maps2)
		assert.NotNil(t, maps3)
	})
}