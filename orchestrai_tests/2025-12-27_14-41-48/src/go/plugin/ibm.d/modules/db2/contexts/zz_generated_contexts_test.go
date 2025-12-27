package contexts

import (
	"strings"
	"testing"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/module"
	"github.com/netdata/netdata/go/plugins/plugin/ibm.d/framework"
)

// TestCleanLabelValue tests the cleanLabelValue function with various inputs
func TestCleanLabelValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic space replacement
		{
			name:     "single space",
			input:    "hello world",
			expected: "hello_world",
		},
		// Dots
		{
			name:     "single dot",
			input:    "example.com",
			expected: "example_com",
		},
		// Dashes
		{
			name:     "single dash",
			input:    "some-value",
			expected: "some_value",
		},
		// Forward slashes
		{
			name:     "single slash",
			input:    "path/to/resource",
			expected: "path_to_resource",
		},
		// Colons
		{
			name:     "single colon",
			input:    "key:value",
			expected: "key_value",
		},
		// Equals signs
		{
			name:     "single equals",
			input:    "key=value",
			expected: "key_value",
		},
		// Commas
		{
			name:     "single comma",
			input:    "value1,value2",
			expected: "value1_value2",
		},
		// Parentheses
		{
			name:     "parentheses",
			input:    "test(value)",
			expected: "test_value_",
		},
		// Multiple special characters
		{
			name:     "multiple special chars",
			input:    "test.value-123/path:name=key,val(1)",
			expected: "test_value_123_path_name_key_val_1_",
		},
		// Uppercase conversion
		{
			name:     "uppercase",
			input:    "HELLO WORLD",
			expected: "hello_world",
		},
		// Mixed case and special chars
		{
			name:     "mixed case with special",
			input:    "Hello.World-Test",
			expected: "hello_world_test",
		},
		// Empty string
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		// Only special characters
		{
			name:     "only special chars",
			input:    ".-/:=,()",
			expected: "__________",
		},
		// Already clean
		{
			name:     "already clean",
			input:    "clean_value_123",
			expected: "clean_value_123",
		},
		// Unicode should be lowercased
		{
			name:     "unicode chars",
			input:    "Café",
			expected: "café",
		},
		// Numbers only
		{
			name:     "numbers only",
			input:    "12345",
			expected: "12345",
		},
		// Leading and trailing special chars
		{
			name:     "leading trailing special",
			input:    "  test  ",
			expected: "__test__",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanLabelValue(tt.input)
			if result != tt.expected {
				t.Errorf("cleanLabelValue(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestEmptyLabelsInstanceID tests EmptyLabels.InstanceID method
func TestEmptyLabelsInstanceID(t *testing.T) {
	tests := []struct {
		name        string
		contextName string
		expected    string
	}{
		{
			name:        "simple context name",
			contextName: "test.context",
			expected:    "test.context",
		},
		{
			name:        "empty context name",
			contextName: "",
			expected:    "",
		},
		{
			name:        "context with special chars",
			contextName: "db2.bufferpool_hit_ratio",
			expected:    "db2.bufferpool_hit_ratio",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labels := EmptyLabels{}
			result := labels.InstanceID(tt.contextName)
			if result != tt.expected {
				t.Errorf("EmptyLabels.InstanceID(%q) = %q, expected %q", tt.contextName, result, tt.expected)
			}
		})
	}
}

// TestBufferpoolLabelsInstanceID tests BufferpoolLabels.InstanceID method
func TestBufferpoolLabelsInstanceID(t *testing.T) {
	tests := []struct {
		name        string
		contextName string
		bufferpool  string
		pageSize    string
		expected    string
	}{
		{
			name:        "basic labels",
			contextName: "db2.bufferpool",
			bufferpool:  "IBMDEFAULTBP",
			pageSize:    "4096",
			expected:    "db2.bufferpool.ibmdefaultbp_4096",
		},
		{
			name:        "labels with spaces",
			contextName: "test.context",
			bufferpool:  "My Buffer Pool",
			pageSize:    "Page Size 4",
			expected:    "test.context.my_buffer_pool_page_size_4",
		},
		{
			name:        "labels with special chars",
			contextName: "test",
			bufferpool:  "BP-1.0",
			pageSize:    "4K/8K",
			expected:    "test.bp_1_0_4k_8k",
		},
		{
			name:        "empty bufferpool",
			contextName: "ctx",
			bufferpool:  "",
			pageSize:    "4096",
			expected:    "ctx._4096",
		},
		{
			name:        "empty page size",
			contextName: "ctx",
			bufferpool:  "BP1",
			pageSize:    "",
			expected:    "ctx.bp1_",
		},
		{
			name:        "both empty",
			contextName: "ctx",
			bufferpool:  "",
			pageSize:    "",
			expected:    "ctx.__",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labels := BufferpoolLabels{
				Bufferpool: tt.bufferpool,
				Page_size:  tt.pageSize,
			}
			result := labels.InstanceID(tt.contextName)
			if result != tt.expected {
				t.Errorf("BufferpoolLabels.InstanceID() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestBufferpoolGroupLabelsInstanceID tests BufferpoolGroupLabels.InstanceID method
func TestBufferpoolGroupLabelsInstanceID(t *testing.T) {
	tests := []struct {
		name        string
		contextName string
		group       string
		expected    string
	}{
		{
			name:        "simple group",
			contextName: "db2.test",
			group:       "group1",
			expected:    "db2.test.group1",
		},
		{
			name:        "group with spaces",
			contextName: "ctx",
			group:       "Group Name",
			expected:    "ctx.group_name",
		},
		{
			name:        "group with special chars",
			contextName: "test",
			group:       "G-1.0",
			expected:    "test.g_1_0",
		},
		{
			name:        "empty group",
			contextName: "ctx",
			group:       "",
			expected:    "ctx.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labels := BufferpoolGroupLabels{Group: tt.group}
			result := labels.InstanceID(tt.contextName)
			if result != tt.expected {
				t.Errorf("BufferpoolGroupLabels.InstanceID() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestConnectionLabelsInstanceID tests ConnectionLabels.InstanceID method
func TestConnectionLabelsInstanceID(t *testing.T) {
	tests := []struct {
		name            string
		contextName     string
		appID           string
		appName         string
		clientHostname  string
		clientIP        string
		clientUser      string
		state           string
		expected        string
	}{
		{
			name:            "all labels populated",
			contextName:     "db2.connection",
			appID:           "1",
			appName:         "app1",
			clientHostname:  "host1",
			clientIP:        "192.168.1.1",
			clientUser:      "user1",
			state:           "active",
			expected:        "db2.connection.1_app1_host1_192_168_1_1_user1_active",
		},
		{
			name:            "labels with spaces",
			contextName:     "ctx",
			appID:           "App ID 1",
			appName:         "My App",
			clientHostname:  "my host",
			clientIP:        "10.0.0.1",
			clientUser:      "my user",
			state:           "in progress",
			expected:        "ctx.app_id_1_my_app_my_host_10_0_0_1_my_user_in_progress",
		},
		{
			name:            "empty labels",
			contextName:     "ctx",
			appID:           "",
			appName:         "",
			clientHostname:  "",
			clientIP:        "",
			clientUser:      "",
			state:           "",
			expected:        "ctx.______",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labels := ConnectionLabels{
				Application_id:   tt.appID,
				Application_name: tt.appName,
				Client_hostname:  tt.clientHostname,
				Client_ip:        tt.clientIP,
				Client_user:      tt.clientUser,
				State:            tt.state,
			}
			result := labels.InstanceID(tt.contextName)
			if result != tt.expected {
				t.Errorf("ConnectionLabels.InstanceID() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestConnectionGroupLabelsInstanceID tests ConnectionGroupLabels.InstanceID method
func TestConnectionGroupLabelsInstanceID(t *testing.T) {
	tests := []struct {
		name        string
		contextName string
		group       string
		expected    string
	}{
		{
			name:        "simple group",
			contextName: "db2.test",
			group:       "admin",
			expected:    "db2.test.admin",
		},
		{
			name:        "group with special chars",
			contextName: "ctx",
			group:       "Group-1.0",
			expected:    "ctx.group_1_0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labels := ConnectionGroupLabels{Group: tt.group}
			result := labels.InstanceID(tt.contextName)
			if result != tt.expected {
				t.Errorf("ConnectionGroupLabels.InstanceID() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestDatabaseLabelsInstanceID tests DatabaseLabels.InstanceID method
func TestDatabaseLabelsInstanceID(t *testing.T) {
	tests := []struct {
		name        string
		contextName string
		database    string
		status      string
		expected    string
	}{
		{
			name:        "simple labels",
			contextName: "db2.db",
			database:    "TESTDB",
			status:      "active",
			expected:    "db2.db.testdb_active",
		},
		{
			name:        "database with spaces",
			contextName: "ctx",
			database:    "My Database",
			status:      "In Progress",
			expected:    "ctx.my_database_in_progress",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labels := DatabaseLabels{
				Database: tt.database,
				Status:   tt.status,
			}
			result := labels.InstanceID(tt.contextName)
			if result != tt.expected {
				t.Errorf("DatabaseLabels.InstanceID() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestIndexLabelsInstanceID tests IndexLabels.InstanceID method
func TestIndexLabelsInstanceID(t *testing.T) {
	tests := []struct {
		name        string
		contextName string
		index       string
		expected    string
	}{
		{
			name:        "simple index",
			contextName: "db2.index",
			index:       "idx_1",
			expected:    "db2.index.idx_1",
		},
		{
			name:        "index with special chars",
			contextName: "ctx",
			index:       "Index.Name-1",
			expected:    "ctx.index_name_1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labels := IndexLabels{Index: tt.index}
			result := labels.InstanceID(tt.contextName)
			if result != tt.expected {
				t.Errorf("IndexLabels.InstanceID() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestMemoryPoolLabelsInstanceID tests MemoryPoolLabels.InstanceID method
func TestMemoryPoolLabelsInstanceID(t *testing.T) {
	tests := []struct {
		name        string
		contextName string
		poolType    string
		expected    string
	}{
		{
			name:        "simple pool type",
			contextName: "db2.mem",
			poolType:    "database",
			expected:    "db2.mem.database",
		},
		{
			name:        "pool type with special chars",
			contextName: "ctx",
			poolType:    "Type-1.0",
			expected:    "ctx.type_1_0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labels := MemoryPoolLabels{Pool_type: tt.poolType}
			result := labels.InstanceID(tt.contextName)
			if result != tt.expected {
				t.Errorf("MemoryPoolLabels.InstanceID() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestMemorySetLabelsInstanceID tests MemorySetLabels.InstanceID method
func TestMemorySetLabelsInstanceID(t *testing.T) {
	tests := []struct {
		name        string
		contextName string
		host        string
		database    string
		setType     string
		member      string
		expected    string
	}{
		{
			name:        "all labels populated",
			contextName: "db2.memset",
			host:        "host1",
			database:    "db1",
			setType:     "type1",
			member:      "member1",
			expected:    "db2.memset.host1_db1_type1_member1",
		},
		{
			name:        "empty labels",
			contextName: "ctx",
			host:        "",
			database:    "",
			setType:     "",
			member:      "",
			expected:    "ctx.____",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labels := MemorySetLabels{
				Host:     tt.host,
				Database: tt.database,
				Set_type: tt.setType,
				Member:   tt.member,
			}
			result := labels.InstanceID(tt.contextName)
			if result != tt.expected {
				t.Errorf("MemorySetLabels.InstanceID() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestPrefetcherLabelsInstanceID tests PrefetcherLabels.InstanceID method
func TestPrefetcherLabelsInstanceID(t *testing.T) {
	tests := []struct {
		name        string
		contextName string
		bufferpool  string
		expected    string
	}{
		{
			name:        "simple bufferpool",
			contextName: "db2.pref",
			bufferpool:  "BP1",
			expected:    "db2.pref.bp1",
		},
		{
			name:        "bufferpool with special chars",
			contextName: "ctx",
			bufferpool:  "BP-1.0",
			expected:    "ctx.bp_1_0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labels := PrefetcherLabels{Bufferpool: tt.bufferpool}
			result := labels.InstanceID(tt.contextName)
			if result != tt.expected {
				t.Errorf("PrefetcherLabels.InstanceID() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestTableLabelsInstanceID tests TableLabels.InstanceID method
func TestTableLabelsInstanceID(t *testing.T) {
	tests := []struct {
		name        string
		contextName string
		table       string
		expected    string
	}{
		{
			name:        "simple table name",
			contextName: "db2.table",
			table:       "TABLE1",
			expected:    "db2.table.table1",
		},
		{
			name:        "table with special chars",
			contextName: "ctx",
			table:       "Table-Name.1",
			expected:    "ctx.table_name_1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labels := TableLabels{Table: tt.table}
			result := labels.InstanceID(tt.contextName)
			if result != tt.expected {
				t.Errorf("TableLabels.InstanceID() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestTableIOLabelsInstanceID tests TableIOLabels.InstanceID method
func TestTableIOLabelsInstanceID(t *testing.T) {
	tests := []struct {
		name        string
		contextName string
		table       string
		expected    string
	}{
		{
			name:        "simple table",
			contextName: "db2.tableio",
			table:       "T1",
			expected:    "db2.tableio.t1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labels := TableIOLabels{Table: tt.table}
			result := labels.InstanceID(tt.contextName)
			if result != tt.expected {
				t.Errorf("TableIOLabels.InstanceID() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestTablespaceLabelsInstanceID tests TablespaceLabels.InstanceID method
func TestTablespaceLabelsInstanceID(t *testing.T) {
	tests := []struct {
		name        string
		contextName string
		tablespace  string
		typeVal     string
		contentType string
		state       string
		expected    string
	}{
		{
			name:        "all labels populated",
			contextName: "db2.ts",
			tablespace:  "TS1",
			typeVal:     "SMS",
			contentType: "data",
			state:       "active",
			expected:    "db2.ts.ts1_sms_data_active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labels := TablespaceLabels{
				Tablespace:   tt.tablespace,
				Type:         tt.typeVal,
				Content_type: tt.contentType,
				State:        tt.state,
			}
			result := labels.InstanceID(tt.contextName)
			if result != tt.expected {
				t.Errorf("TablespaceLabels.InstanceID() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestTablespaceGroupLabelsInstanceID tests TablespaceGroupLabels.InstanceID method
func TestTablespaceGroupLabelsInstanceID(t *testing.T) {
	tests := []struct {
		name        string
		contextName string
		group       string
		expected    string
	}{
		{
			name:        "simple group",
			contextName: "db2.tsg",
			group:       "G1",
			expected:    "db2.tsg.g1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labels := TablespaceGroupLabels{Group: tt.group}
			result := labels.InstanceID(tt.contextName)
			if result != tt.expected {
				t.Errorf("TablespaceGroupLabels.InstanceID() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestBufferpoolContexts tests Bufferpool context Set methods
func TestBufferpoolContexts(t *testing.T) {
	t.Run("BufferpoolHitRatio.Set", func(t *testing.T) {
		state := &framework.CollectorState{}
		labels := BufferpoolLabels{Bufferpool: "BP1", Page_size: "4096"}
		values := BufferpoolHitRatioValues{Overall: 95000}

		ctx := Bufferpool.HitRatio
		ctx.Set(state, labels, values)

		// Verify the context properties
		if ctx.Context.Name != "db2.bufferpool_instance_hit_ratio" {
			t.Errorf("Expected context name 'db2.bufferpool_instance_hit_ratio', got %s", ctx.Context.Name)
		}
		if ctx.Context.Units != "percentage" {
			t.Errorf("Expected units 'percentage', got %s", ctx.Context.Units)
		}
	})

	t.Run("BufferpoolDetailedHitRatio.Set", func(t *testing.T) {
		state := &framework.CollectorState{}
		labels := BufferpoolLabels{Bufferpool: "BP1", Page_size: "4096"}
		values := BufferpoolDetailedHitRatioValues{Data: 90000, Index: 85000, Xda: 80000, Column: 75000}

		ctx := Bufferpool.DetailedHitRatio
		ctx.Set(state, labels, values)

		if ctx.Context.Name != "db2.bufferpool_instance_detailed_hit_ratio" {
			t.Errorf("Expected context name 'db2.bufferpool_instance_detailed_hit_ratio'")
		}
	})

	t.Run("BufferpoolReads.Set", func(t *testing.T) {
		state := &framework.CollectorState{}
		labels := BufferpoolLabels{Bufferpool: "BP1", Page_size: "4096"}
		values := BufferpoolReadsValues{Logical: 1000, Physical: 500}

		ctx := Bufferpool.Reads
		ctx.Set(state, labels, values)

		if ctx.Context.Name != "db2.bufferpool_instance_reads" {
			t.Errorf("Expected context name 'db2.bufferpool_instance_reads'")
		}
	})

	t.Run("BufferpoolPages.Set", func(t *testing.T) {
		state := &framework.CollectorState{}
		labels := BufferpoolLabels{Bufferpool: "BP1", Page_size: "4096"}
		values := BufferpoolPagesValues{Used: 5000, Total: 10000}

		ctx := Bufferpool.Pages
		ctx.Set(state, labels, values)

		if ctx.Context.Name != "db2.bufferpool_instance_pages" {
			t.Errorf("Expected context name 'db2.bufferpool_instance_pages'")
		}
	})

	t.Run("BufferpoolWrites.Set", func(t *testing.T) {
		state := &framework.CollectorState{}
		labels := BufferpoolLabels{Bufferpool: "BP1", Page_size: "4096"}
		values := BufferpoolWritesValues{Writes: 2000}

		ctx := Bufferpool.Writes
		ctx.Set(state, labels, values)

		if ctx.Context.Name != "db2.bufferpool_instance_writes" {
			t.Errorf("Expected context name 'db2.bufferpool_instance_writes'")
		}
	})
}

// TestBufferpoolSetUpdateEvery tests SetUpdateEvery for Bufferpool contexts
func TestBufferpoolSetUpdateEvery(t *testing.T) {
	tests := []struct {
		name     string
		updateIn int
	}{
		{name: "zero update", updateIn: 0},
		{name: "small update", updateIn: 1},
		{name: "large update", updateIn: 60},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &framework.CollectorState{}
			labels := BufferpoolLabels{Bufferpool: "BP1", Page_size: "4096"}

			ctx := Bufferpool.HitRatio
			ctx.SetUpdateEvery(state, labels, tt.updateIn)

			// Verify no panic occurred and context is still valid
			if ctx.Context.Name == "" {
				t.Error("Context name should not be empty after SetUpdateEvery")
			}
		})
	}
}

// TestConnectionContexts tests Connection context Set methods
func TestConnectionContexts(t *testing.T) {
	t.Run("ConnectionState.Set", func(t *testing.T) {
		state := &framework.CollectorState{}
		labels := ConnectionLabels{
			Application_id:   "1",
			Application_name: "app1",
			Client_hostname:  "host1",
			Client_ip:        "192.168.1.1",
			Client_user:      "user1",
			State:            "active",
		}
		values := ConnectionStateValues{State: 1}

		ctx := Connection.State
		ctx.Set(state, labels, values)

		if ctx.Context.Name != "db2.connection_state" {
			t.Errorf("Expected context name 'db2.connection_state'")
		}
	})

	t.Run("ConnectionActivity.Set", func(t *testing.T) {
		state := &framework.CollectorState{}
		labels := ConnectionLabels{
			Application_id:   "1",
			Application_name: "app1",
			Client_hostname:  "host1",
			Client_ip:        "192.168.1.1",
			Client_user:      "user1",
			State:            "active",
		}
		values := ConnectionActivityValues{Read: 1000, Written: 500}

		ctx := Connection.Activity
		ctx.Set(state, labels, values)

		if ctx.Context.Name != "db2.connection_activity" {
			t.Errorf("Expected context name 'db2.connection_activity'")
		}
	})

	t.Run("ConnectionWaitTime.Set", func(t *testing.T) {
		state := &framework.CollectorState{}
		labels := ConnectionLabels{
			Application_id:   "1",
			Application_name: "app1",
			Client_hostname:  "host1",
			Client_ip:        "192.168.1.1",
			Client_user:      "user1",
			State:            "active",
		}
		values := ConnectionWaitTimeValues{
			Lock:         100,
			Log_disk:     200,
			Log_buffer:   300,
			Pool_read:    400,
			Pool_write:   500,
			Direct_read:  600,
			Direct_write: 700,
			Fcm_recv:     800,
			Fcm_send:     900,
		}

		ctx := Connection.WaitTime
		ctx.Set(state, labels, values)

		if ctx.Context.Name != "db2.connection_wait_time" {
			t.Errorf("Expected context name 'db2.connection_wait_time'")
		}
	})

	t.Run("ConnectionProcessingTime.Set", func(t *testing.T) {
		state := &framework.CollectorState{}
		labels := ConnectionLabels{
			Application_id:   "1",
			Application_name: "app1",
			Client_hostname:  "host1",
			Client_ip:        "192.168.1.1",
			Client_user:      "user1",
			State:            "active",
		}
		values := ConnectionProcessingTimeValues{
			Routine:  100,
			Compile:  200,
			Section:  300,
			Commit:   400,
			Rollback: 500,
		}

		ctx := Connection.ProcessingTime
		ctx.Set(state, labels, values)

		if ctx.Context.Name != "db2.connection_processing_time" {
			t.Errorf("Expected context name 'db2.connection_processing_time'")
		}
	})
}

// TestConnectionGroupContexts tests ConnectionGroup context Set methods
func TestConnectionGroupContexts(t *testing.T) {
	t.Run("ConnectionGroupCount.Set", func(t *testing.T) {
		state := &framework.CollectorState{}
		labels := ConnectionGroupLabels{Group: "admin"}
		values := ConnectionGroupCountValues{Count: 5}

		ctx := ConnectionGroup.Count
		ctx.Set(state, labels, values)

		if ctx.Context.Name != "db2.connection_group.count" {
			t.Errorf("Expected context name 'db2.connection_group.count'")
		}
	})
}

// TestDatabaseContexts tests Database context Set methods
func TestDatabaseContexts(t *testing.T) {
	t.Run("DatabaseStatus.Set", func(t *testing.T) {
		state := &framework.CollectorState{}
		labels := DatabaseLabels{Database: "TESTDB", Status: "active"}
		values := DatabaseStatusValues{Status: 1}

		ctx := Database.Status
		ctx.Set(state, labels, values)

		if ctx.Context.Name != "db2.database_instance_status" {
			t.Errorf("Expected context name 'db2.database_instance_status'")
		}
	})

	t.Run("DatabaseApplications.Set", func(t *testing.T) {
		state := &framework.CollectorState{}
		labels := DatabaseLabels{Database: "TESTDB", Status: "active"}
		values := DatabaseApplicationsValues{Applications: 10}

		ctx := Database.Applications
		ctx.Set(