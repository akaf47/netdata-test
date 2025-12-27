package couchdb

import (
	"testing"
)

// Test cdbMetrics.empty() method
func TestCdbMetricsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		metrics  cdbMetrics
		expected bool
	}{
		{
			name:     "empty metrics - all nil and empty",
			metrics:  cdbMetrics{},
			expected: true,
		},
		{
			name: "metrics with ActiveTasks",
			metrics: cdbMetrics{
				ActiveTasks: []cdbActiveTask{},
			},
			expected: false,
		},
		{
			name: "metrics with NodeStats",
			metrics: cdbMetrics{
				NodeStats: &cdbNodeStats{},
			},
			expected: false,
		},
		{
			name: "metrics with NodeSystem",
			metrics: cdbMetrics{
				NodeSystem: &cdbNodeSystem{},
			},
			expected: false,
		},
		{
			name: "metrics with DBStats",
			metrics: cdbMetrics{
				DBStats: []cdbDBStats{},
			},
			expected: false,
		},
		{
			name: "metrics with all fields populated",
			metrics: cdbMetrics{
				ActiveTasks: []cdbActiveTask{{Type: "indexer"}},
				NodeStats:   &cdbNodeStats{},
				NodeSystem:  &cdbNodeSystem{},
				DBStats:     []cdbDBStats{{Key: "test"}},
			},
			expected: false,
		},
		{
			name: "metrics with only NodeStats and DBStats",
			metrics: cdbMetrics{
				NodeStats: &cdbNodeStats{},
				DBStats:   []cdbDBStats{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.metrics.empty()
			if result != tt.expected {
				t.Errorf("empty() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test cdbMetrics.hasActiveTasks() method
func TestCdbMetricsHasActiveTasks(t *testing.T) {
	tests := []struct {
		name     string
		metrics  cdbMetrics
		expected bool
	}{
		{
			name:     "nil ActiveTasks",
			metrics:  cdbMetrics{ActiveTasks: nil},
			expected: false,
		},
		{
			name: "empty ActiveTasks slice",
			metrics: cdbMetrics{
				ActiveTasks: []cdbActiveTask{},
			},
			expected: true,
		},
		{
			name: "non-empty ActiveTasks slice",
			metrics: cdbMetrics{
				ActiveTasks: []cdbActiveTask{
					{Type: "indexer"},
					{Type: "replicator"},
				},
			},
			expected: true,
		},
		{
			name: "single ActiveTask",
			metrics: cdbMetrics{
				ActiveTasks: []cdbActiveTask{{Type: "compaction"}},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.metrics.hasActiveTasks()
			if result != tt.expected {
				t.Errorf("hasActiveTasks() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test cdbMetrics.hasNodeStats() method
func TestCdbMetricsHasNodeStats(t *testing.T) {
	tests := []struct {
		name     string
		metrics  cdbMetrics
		expected bool
	}{
		{
			name:     "nil NodeStats",
			metrics:  cdbMetrics{NodeStats: nil},
			expected: false,
		},
		{
			name: "empty NodeStats",
			metrics: cdbMetrics{
				NodeStats: &cdbNodeStats{},
			},
			expected: true,
		},
		{
			name: "NodeStats with populated CouchDB data",
			metrics: cdbMetrics{
				NodeStats: &cdbNodeStats{
					CouchDB: struct {
						DatabaseReads struct {
							Value float64 `stm:"" json:"value"`
						} `stm:"database_reads" json:"database_reads"`
						DatabaseWrites struct {
							Value float64 `stm:"" json:"value"`
						} `stm:"database_writes" json:"database_writes"`
						HTTPd struct {
							ViewReads struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"view_reads" json:"view_reads"`
						} `stm:"httpd" json:"httpd"`
						HTTPdRequestMethods struct {
							Copy struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"COPY" json:"COPY"`
							Delete struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"DELETE" json:"DELETE"`
							Get struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"GET" json:"GET"`
							Head struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"HEAD" json:"HEAD"`
							Options struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"OPTIONS" json:"OPTIONS"`
							Post struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"POST" json:"POST"`
							Put struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"PUT" json:"PUT"`
						} `stm:"httpd_request_methods" json:"httpd_request_methods"`
						HTTPdStatusCodes struct {
							Code200 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"200" json:"200"`
							Code201 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"201" json:"201"`
							Code202 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"202" json:"202"`
							Code204 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"204" json:"204"`
							Code206 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"206" json:"206"`
							Code301 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"301" json:"301"`
							Code302 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"302" json:"302"`
							Code304 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"304" json:"304"`
							Code400 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"400" json:"400"`
							Code401 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"401" json:"401"`
							Code403 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"403" json:"403"`
							Code404 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"404" json:"404"`
							Code405 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"405" json:"405"`
							Code406 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"406" json:"406"`
							Code409 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"409" json:"409"`
							Code412 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"412" json:"412"`
							Code413 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"413" json:"413"`
							Code414 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"414" json:"414"`
							Code415 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"415" json:"415"`
							Code416 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"416" json:"416"`
							Code417 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"417" json:"417"`
							Code500 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"500" json:"500"`
							Code501 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"501" json:"501"`
							Code503 struct {
								Value float64 `stm:"" json:"value"`
							} `stm:"503" json:"503"`
						} `stm:"httpd_status_codes" json:"httpd_status_codes"`
						OpenOSFiles struct {
							Value float64 `stm:"" json:"value"`
						} `stm:"open_os_files" json:"open_os_files"`
					}{},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.metrics.hasNodeStats()
			if result != tt.expected {
				t.Errorf("hasNodeStats() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test cdbMetrics.hasNodeSystem() method
func TestCdbMetricsHasNodeSystem(t *testing.T) {
	tests := []struct {
		name     string
		metrics  cdbMetrics
		expected bool
	}{
		{
			name:     "nil NodeSystem",
			metrics:  cdbMetrics{NodeSystem: nil},
			expected: false,
		},
		{
			name: "empty NodeSystem",
			metrics: cdbMetrics{
				NodeSystem: &cdbNodeSystem{},
			},
			expected: true,
		},
		{
			name: "NodeSystem with populated data",
			metrics: cdbMetrics{
				NodeSystem: &cdbNodeSystem{
					RunQueue:        10.5,
					EtsTableCount:   15.0,
					ContextSwitches: 100.0,
					Reductions:      1000.0,
					IOInput:         5000.0,
					IOOutput:        2000.0,
					OSProcCount:     4.0,
					ProcessCount:    50.0,
				},
			},
			expected: true,
		},
		{
			name: "NodeSystem with memory info",
			metrics: cdbMetrics{
				NodeSystem: &cdbNodeSystem{
					Memory: struct {
						Other     float64 `stm:"other" json:"other"`
						Atom      float64 `stm:"atom" json:"atom"`
						AtomUsed  float64 `stm:"atom_used" json:"atom_used"`
						Processes float64 `stm:"processes" json:"processes"`
						Binary    float64 `stm:"binary" json:"binary"`
						Code      float64 `stm:"code" json:"code"`
						Ets       float64 `stm:"ets" json:"ets"`
					}{
						Other:     100,
						Atom:      200,
						AtomUsed:  150,
						Processes: 500,
						Binary:    300,
						Code:      400,
						Ets:       250,
					},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.metrics.hasNodeSystem()
			if result != tt.expected {
				t.Errorf("hasNodeSystem() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test cdbMetrics.hasDBStats() method
func TestCdbMetricsHasDBStats(t *testing.T) {
	tests := []struct {
		name     string
		metrics  cdbMetrics
		expected bool
	}{
		{
			name:     "nil DBStats",
			metrics:  cdbMetrics{DBStats: nil},
			expected: false,
		},
		{
			name: "empty DBStats slice",
			metrics: cdbMetrics{
				DBStats: []cdbDBStats{},
			},
			expected: true,
		},
		{
			name: "single DBStat",
			metrics: cdbMetrics{
				DBStats: []cdbDBStats{
					{
						Key:   "db1",
						Error: "",
					},
				},
			},
			expected: true,
		},
		{
			name: "multiple DBStats with data",
			metrics: cdbMetrics{
				DBStats: []cdbDBStats{
					{
						Key:   "db1",
						Error: "",
						Info: struct {
							Sizes struct {
								File     float64 `stm:"file" json:"file"`
								External float64 `stm:"external" json:"external"`
								Active   float64 `stm:"active" json:"active"`
							} `stm:"db_sizes" json:"sizes"`
							DocDelCount float64 `stm:"db_doc_del_counts" json:"doc_del_count"`
							DocCount    float64 `stm:"db_doc_counts" json:"doc_count"`
						}{
							Sizes: struct {
								File     float64 `stm:"file" json:"file"`
								External float64 `stm:"external" json:"external"`
								Active   float64 `stm:"active" json:"active"`
							}{
								File:     1024,
								External: 512,
								Active:   256,
							},
							DocDelCount: 10,
							DocCount:    1000,
						},
					},
					{
						Key: "db2",
						Info: struct {
							Sizes struct {
								File     float64 `stm:"file" json:"file"`
								External float64 `stm:"external" json:"external"`
								Active   float64 `stm:"active" json:"active"`
							} `stm:"db_sizes" json:"sizes"`
							DocDelCount float64 `stm:"db_doc_del_counts" json:"doc_del_count"`
							DocCount    float64 `stm:"db_doc_counts" json:"doc_count"`
						}{
							DocCount: 500,
						},
					},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.metrics.hasDBStats()
			if result != tt.expected {
				t.Errorf("hasDBStats() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test cdbActiveTask structure
func TestCdbActiveTaskStructure(t *testing.T) {
	tests := []struct {
		name string
		task cdbActiveTask
		want string
	}{
		{
			name: "indexer task",
			task: cdbActiveTask{Type: "indexer"},
			want: "indexer",
		},
		{
			name: "replicator task",
			task: cdbActiveTask{Type: "replicator"},
			want: "replicator",
		},
		{
			name: "compaction task",
			task: cdbActiveTask{Type: "compaction"},
			want: "compaction",
		},
		{
			name: "empty type",
			task: cdbActiveTask{Type: ""},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.task.Type != tt.want {
				t.Errorf("Type = %q, want %q", tt.task.Type, tt.want)
			}
		})
	}
}

// Test edge case: empty() method with mixed nil and non-nil fields
func TestCdbMetricsEmptyWithMixedFields(t *testing.T) {
	tests := []struct {
		name     string
		metrics  cdbMetrics
		expected bool
	}{
		{
			name: "only ActiveTasks is nil, rest are not",
			metrics: cdbMetrics{
				ActiveTasks: nil,
				NodeStats:   &cdbNodeStats{},
			},
			expected: false,
		},
		{
			name: "only NodeStats is nil, rest are not",
			metrics: cdbMetrics{
				ActiveTasks: []cdbActiveTask{},
				NodeStats:   nil,
				NodeSystem:  &cdbNodeSystem{},
			},
			expected: false,
		},
		{
			name: "only NodeSystem is nil, rest are not",
			metrics: cdbMetrics{
				ActiveTasks: []cdbActiveTask{},
				NodeStats:   nil,
				NodeSystem:  nil,
				DBStats:     []cdbDBStats{},
			},
			expected: false,
		},
		{
			name: "all fields are nil",
			metrics: cdbMetrics{
				ActiveTasks: nil,
				NodeStats:   nil,
				NodeSystem:  nil,
				DBStats:     nil,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.metrics.empty()
			if result != tt.expected {
				t.Errorf("empty() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test cdbNodeStats structure with all HTTP status codes
func TestCdbNodeStatsHTTPStatusCodes(t *testing.T) {
	stats := &cdbNodeStats{}

	// Test that all status code fields can be set
	stats.CouchDB.HTTPdStatusCodes.Code200.Value = 100
	stats.CouchDB.HTTPdStatusCodes.Code201.Value = 50
	stats.CouchDB.HTTPdStatusCodes.Code202.Value = 25
	stats.CouchDB.HTTPdStatusCodes.Code204.Value = 10
	stats.CouchDB.HTTPdStatusCodes.Code206.Value = 5
	stats.CouchDB.HTTPdStatusCodes.Code301.Value = 1
	stats.CouchDB.HTTPdStatusCodes.Code302.Value = 2
	stats.CouchDB.HTTPdStatusCodes.Code304.Value = 3
	stats.CouchDB.HTTPdStatusCodes.Code400.Value = 4
	stats.CouchDB.HTTPdStatusCodes.Code401.Value = 5
	stats.CouchDB.HTTPdStatusCodes.Code403.Value = 6
	stats.CouchDB.HTTPdStatusCodes.Code404.Value = 7
	stats.CouchDB.HTTPdStatusCodes.Code405.Value = 8
	stats.CouchDB.HTTPdStatusCodes.Code406.Value = 9
	stats.CouchDB.HTTPdStatusCodes.Code409.Value = 11
	stats.CouchDB.HTTPdStatusCodes.Code412.Value = 12
	stats.CouchDB.HTTPdStatusCodes.Code413.Value = 13
	stats.CouchDB.HTTPdStatusCodes.Code414.Value = 14
	stats.CouchDB.HTTPdStatusCodes.Code415.Value = 15
	stats.CouchDB.HTTPdStatusCodes.Code416.Value = 16
	stats.CouchDB.HTTPdStatusCodes.Code417.Value = 17
	stats.CouchDB.HTTPdStatusCodes.Code500.Value = 18
	stats.CouchDB.HTTPdStatusCodes.Code501.Value = 19
	stats.CouchDB.HTTPdStatusCodes.Code503.Value = 20

	if stats.CouchDB.HTTPdStatusCodes.Code200.Value != 100 {
		t.Error("Code200 not set correctly")
	}
	if stats.CouchDB.HTTPdStatusCodes.Code503.Value != 20 {
		t.Error("Code503 not set correctly")
	}
}

// Test cdbNodeStats structure with all HTTP request methods
func TestCdbNodeStatsHTTPRequestMethods(t *testing.T) {
	stats := &cdbNodeStats{}

	// Test that all HTTP method fields can be set
	stats.CouchDB.HTTPdRequestMethods.Copy.Value = 10
	stats.CouchDB.HTTPdRequestMethods.Delete.Value = 20
	stats.CouchDB.HTTPdRequestMethods.Get.Value = 30
	stats.CouchDB.HTTPdRequestMethods.Head.Value = 5
	stats.CouchDB.HTTPdRequestMethods.Options.Value = 1
	stats.CouchDB.HTTPdRequestMethods.Post.Value = 40
	stats.CouchDB.HTTPdRequestMethods.Put.Value = 50

	if stats.CouchDB.HTTPdRequestMethods.Get.Value != 30 {
		t.Error("GET method not set correctly")
	}
	if stats.CouchDB.HTTPdRequestMethods.Post.Value != 40 {
		t.Error("POST method not set correctly")
	}
}

// Test cdbNodeSystem with message queues
func TestCdbNodeSystemMessageQueues(t *testing.T) {
	system := &cdbNodeSystem{
		MessageQueues: map[string]any{
			"queue1": 100,
			"queue2": "active",
			"queue3": true,
		},
	}

	if system.MessageQueues == nil {
		t.Error("MessageQueues is nil")
	}
	if len(system.MessageQueues) != 3 {
		t.Errorf("MessageQueues length = %d, want 3", len(system.MessageQueues))
	}
}

// Test cdbDBStats with error field
func TestCdbDBStatsWithError(t *testing.T) {
	tests := []struct {
		name      string
		dbStats   cdbDBStats
		expectErr bool
	}{
		{
			name: "DBStat with no error",
			dbStats: cdbDBStats{
				Key:   "test_db",
				Error: "",
			},
			expectErr: false,
		},
		{
			name: "DBStat with error",
			dbStats: cdbDBStats{
				Key:   "failed_db",
				Error: "database not found",
			},
			expectErr: true,
		},
		{
			name: "DBStat with empty key but error",
			dbStats: cdbDBStats{
				Key:   "",
				Error: "unknown error",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasError := tt.dbStats.Error != ""
			if hasError != tt.expectErr {
				t.Errorf("hasError = %v, want %v", hasError, tt.expectErr)
			}
		})
	}
}

// Test cdbReplicator jobs structure
func TestCdbNodeStatsReplicatorJobs(t *testing.T) {
	stats := &cdbNodeStats{
		CouchReplicator: struct {
			Jobs struct {
				Running struct {
					Value float64 `stm:"" json:"value"`
				} `stm:"running" json:"running"`
				Pending struct {
					Value float64 `stm:"" json:"value"`
				} `stm:"pending" json:"pending"`
				Crashed struct {
					Value float64 `stm:"" json:"value"`
				} `stm:"crashed" json:"crashed"`
			} `stm:"jobs" json:"jobs"`
		}{
			Jobs: struct {
				Running struct {
					Value float64 `stm:"" json:"value"`
				} `stm:"running" json:"running"`
				Pending struct {
					Value float64 `stm:"" json:"value"`
				} `stm:"pending" json:"pending"`
				Crashed struct {
					Value float64 `stm:"" json:"value"`
				} `stm:"crashed" json:"crashed"`
			}{
				Running: struct {
					Value float64 `stm:"" json:"value"`
				}{Value: 5},
				Pending: struct {
					Value float64 `stm:"" json:"value"`
				}{Value: 3},
				Crashed: struct {
					Value float64 `stm:"" json:"value"`
				}{Value: 1},
			},
		},
	}

	if stats.CouchReplicator.Jobs.Running.Value != 5 {
		t.Error("Running jobs not set correctly")
	}
	if stats.CouchReplicator.Jobs.Pending.Value != 3 {
		t.Error("Pending jobs not set correctly")
	}
	if stats.CouchReplicator.Jobs.Crashed.Value != 1 {
		t.Error("Crashed jobs not set correctly")
	}
}

// Test cdbNodeSystem memory structure
func TestCdbNodeSystemMemory(t *testing.T) {
	system := &cdbNodeSystem{
		Memory: struct {
			Other     float64 `stm:"other" json:"other"`
			Atom      float64 `stm:"atom" json:"atom"`
			AtomUsed  float64 `stm:"atom_used" json:"atom_used"`
			Processes float64 `stm:"processes" json:"processes"`
			Binary    float64 `stm:"binary" json:"binary"`
			Code      float64 `stm:"code" json:"code"`
			Ets       float64 `stm:"ets" json:"ets"`
		}{
			Other:     1000,
			Atom:      2000,
			AtomUsed:  1500,
			Processes: 5000,
			Binary:    3000,
			Code:      4000,
			Ets:       2500,
		},
	}

	if system.Memory.Other != 1000 {
		t.Error("Memory.Other not set correctly")
	}
	if system.Memory.Processes != 5000 {
		t.Error("Memory.Processes not set correctly")
	}
}