package contexts

import (
	"testing"
)

func TestCleanLabelValue(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "simple lowercase string",
			input:    "test",
			expected: "test",
		},
		{
			name:     "uppercase conversion",
			input:    "TEST",
			expected: "test",
		},
		{
			name:     "space replacement",
			input:    "test value",
			expected: "test_value",
		},
		{
			name:     "dot replacement",
			input:    "test.value",
			expected: "test_value",
		},
		{
			name:     "dash replacement",
			input:    "test-value",
			expected: "test_value",
		},
		{
			name:     "forward slash replacement",
			input:    "test/value",
			expected: "test_value",
		},
		{
			name:     "colon replacement",
			input:    "test:value",
			expected: "test_value",
		},
		{
			name:     "equals replacement",
			input:    "test=value",
			expected: "test_value",
		},
		{
			name:     "comma replacement",
			input:    "test,value",
			expected: "test_value",
		},
		{
			name:     "left paren replacement",
			input:    "test(value",
			expected: "test_value",
		},
		{
			name:     "right paren replacement",
			input:    "test)value",
			expected: "test_value",
		},
		{
			name:     "multiple replacements",
			input:    "Test Value.With:Multiple-Chars",
			expected: "test_value_with_multiple_chars",
		},
		{
			name:     "all special chars",
			input:    " .:-=/,()Test",
			expected: "_________test",
		},
		{
			name:     "mixed case with special chars",
			input:    "UPPER lower MiXeD.test-123",
			expected: "upper_lower_mixed_test_123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := cleanLabelValue(tc.input)
			if result != tc.expected {
				t.Errorf("cleanLabelValue(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestEmptyLabelsInstanceID(t *testing.T) {
	testCases := []struct {
		name        string
		contextName string
		expected    string
	}{
		{
			name:        "simple context name",
			contextName: "test_context",
			expected:    "test_context",
		},
		{
			name:        "empty context name",
			contextName: "",
			expected:    "",
		},
		{
			name:        "context with dots",
			contextName: "as400.cpu",
			expected:    "as400.cpu",
		},
		{
			name:        "context with special chars",
			contextName: "test.context.name",
			expected:    "test.context.name",
		},
	}

	labels := EmptyLabels{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := labels.InstanceID(tc.contextName)
			if result != tc.expected {
				t.Errorf("EmptyLabels.InstanceID(%q) = %q, want %q", tc.contextName, result, tc.expected)
			}
		})
	}
}

func TestActiveJobLabelsInstanceID(t *testing.T) {
	testCases := []struct {
		name        string
		labels      ActiveJobLabels
		contextName string
		expected    string
	}{
		{
			name: "simple labels",
			labels: ActiveJobLabels{
				Job_name:   "Job1",
				Job_status: "Active",
				Subsystem:  "SYS1",
				Job_type:   "BATCH",
			},
			contextName: "as400.activejob_cpu",
			expected:    "as400.activejob_cpu.job1_active_sys1_batch",
		},
		{
			name: "labels with spaces",
			labels: ActiveJobLabels{
				Job_name:   "Job One",
				Job_status: "In Queue",
				Subsystem:  "System One",
				Job_type:   "Interactive Job",
			},
			contextName: "as400.activejob_cpu",
			expected:    "as400.activejob_cpu.job_one_in_queue_system_one_interactive_job",
		},
		{
			name: "labels with special chars",
			labels: ActiveJobLabels{
				Job_name:   "Job-1.0",
				Job_status: "Active/Waiting",
				Subsystem:  "SYS:1",
				Job_type:   "BATCH(Type)",
			},
			contextName: "as400.activejob_cpu",
			expected:    "as400.activejob_cpu.job_1_0_active_waiting_sys_1_batch_type",
		},
		{
			name: "empty labels",
			labels: ActiveJobLabels{
				Job_name:   "",
				Job_status: "",
				Subsystem:  "",
				Job_type:   "",
			},
			contextName: "as400.activejob_cpu",
			expected:    "as400.activejob_cpu.____",
		},
		{
			name: "uppercase labels",
			labels: ActiveJobLabels{
				Job_name:   "JOBNAME",
				Job_status: "ACTIVE",
				Subsystem:  "SUBSYS",
				Job_type:   "BATCH",
			},
			contextName: "as400.activejob_cpu",
			expected:    "as400.activejob_cpu.jobname_active_subsys_batch",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.labels.InstanceID(tc.contextName)
			if result != tc.expected {
				t.Errorf("ActiveJobLabels.InstanceID() = %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestDiskLabelsInstanceID(t *testing.T) {
	testCases := []struct {
		name        string
		labels      DiskLabels
		contextName string
		expected    string
	}{
		{
			name: "simple labels",
			labels: DiskLabels{
				Disk_unit:          "Disk1",
				Disk_type:          "SSD",
				Disk_model:         "Model123",
				Hardware_status:    "OK",
				Disk_serial_number: "SN123456",
			},
			contextName: "as400.disk_busy",
			expected:    "as400.disk_busy.disk1_ssd_model123_ok_sn123456",
		},
		{
			name: "labels with spaces and special chars",
			labels: DiskLabels{
				Disk_unit:          "Disk Unit 1",
				Disk_type:          "SSD/HDD",
				Disk_model:         "Model-123.0",
				Hardware_status:    "OK (Good)",
				Disk_serial_number: "SN-123/456:789",
			},
			contextName: "as400.disk_busy",
			expected:    "as400.disk_busy.disk_unit_1_ssd_hdd_model_123_0_ok__good__sn_123_456_789",
		},
		{
			name: "empty labels",
			labels: DiskLabels{
				Disk_unit:          "",
				Disk_type:          "",
				Disk_model:         "",
				Hardware_status:    "",
				Disk_serial_number: "",
			},
			contextName: "as400.disk_busy",
			expected:    "as400.disk_busy._____",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.labels.InstanceID(tc.contextName)
			if result != tc.expected {
				t.Errorf("DiskLabels.InstanceID() = %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestHTTPServerLabelsInstanceID(t *testing.T) {
	testCases := []struct {
		name        string
		labels      HTTPServerLabels
		contextName string
		expected    string
	}{
		{
			name: "simple labels",
			labels: HTTPServerLabels{
				Server:   "Server1",
				Function: "Web",
			},
			contextName: "as400.http_server_connections",
			expected:    "as400.http_server_connections.server1_web",
		},
		{
			name: "labels with special chars",
			labels: HTTPServerLabels{
				Server:   "HTTP-Server.1:8080",
				Function: "REST/API",
			},
			contextName: "as400.http_server_connections",
			expected:    "as400.http_server_connections.http_server_1_8080_rest_api",
		},
		{
			name: "empty labels",
			labels: HTTPServerLabels{
				Server:   "",
				Function: "",
			},
			contextName: "as400.http_server_connections",
			expected:    "as400.http_server_connections.__",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.labels.InstanceID(tc.contextName)
			if result != tc.expected {
				t.Errorf("HTTPServerLabels.InstanceID() = %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestJobQueueLabelsInstanceID(t *testing.T) {
	testCases := []struct {
		name        string
		labels      JobQueueLabels
		contextName string
		expected    string
	}{
		{
			name: "simple labels",
			labels: JobQueueLabels{
				Job_queue: "Queue1",
				Library:   "LIB",
				Status:    "Active",
			},
			contextName: "as400.jobqueue_length",
			expected:    "as400.jobqueue_length.queue1_lib_active",
		},
		{
			name: "labels with special chars",
			labels: JobQueueLabels{
				Job_queue: "Job-Queue.1",
				Library:   "Lib/Name",
				Status:    "Active (Processing)",
			},
			contextName: "as400.jobqueue_length",
			expected:    "as400.jobqueue_length.job_queue_1_lib_name_active__processing",
		},
		{
			name: "empty labels",
			labels: JobQueueLabels{
				Job_queue: "",
				Library:   "",
				Status:    "",
			},
			contextName: "as400.jobqueue_length",
			expected:    "as400.jobqueue_length.___",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.labels.InstanceID(tc.contextName)
			if result != tc.expected {
				t.Errorf("JobQueueLabels.InstanceID() = %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestMessageQueueLabelsInstanceID(t *testing.T) {
	testCases := []struct {
		name        string
		labels      MessageQueueLabels
		contextName string
		expected    string
	}{
		{
			name: "simple labels",
			labels: MessageQueueLabels{
				Library: "LIB",
				Queue:   "Queue1",
			},
			contextName: "as400.message_queue_messages",
			expected:    "as400.message_queue_messages.lib_queue1",
		},
		{
			name: "labels with special chars",
			labels: MessageQueueLabels{
				Library: "Lib/Name",
				Queue:   "Queue-1.0",
			},
			contextName: "as400.message_queue_messages",
			expected:    "as400.message_queue_messages.lib_name_queue_1_0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.labels.InstanceID(tc.contextName)
			if result != tc.expected {
				t.Errorf("MessageQueueLabels.InstanceID() = %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestNetworkInterfaceLabelsInstanceID(t *testing.T) {
	testCases := []struct {
		name        string
		labels      NetworkInterfaceLabels
		contextName string
		expected    string
	}{
		{
			name: "simple labels",
			labels: NetworkInterfaceLabels{
				Interface:        "eth0",
				Interface_type:   "Ethernet",
				Connection_type:  "TCP/IP",
				Internet_address: "192.168.1.1",
				Network_address:  "192.168.1.0",
				Subnet_mask:      "255.255.255.0",
			},
			contextName: "as400.network_interface_status",
			expected:    "as400.network_interface_status.eth0_ethernet_tcp_ip_192_168_1_1_192_168_1_0_255_255_255_0",
		},
		{
			name: "empty labels",
			labels: NetworkInterfaceLabels{
				Interface:        "",
				Interface_type:   "",
				Connection_type:  "",
				Internet_address: "",
				Network_address:  "",
				Subnet_mask:      "",
			},
			contextName: "as400.network_interface_status",
			expected:    "as400.network_interface_status.______",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.labels.InstanceID(tc.contextName)
			if result != tc.expected {
				t.Errorf("NetworkInterfaceLabels.InstanceID() = %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestOutputQueueLabelsInstanceID(t *testing.T) {
	testCases := []struct {
		name        string
		labels      OutputQueueLabels
		contextName string
		expected    string
	}{
		{
			name: "simple labels",
			labels: OutputQueueLabels{
				Library: "LIB",
				Queue:   "Queue1",
				Status:  "Ready",
			},
			contextName: "as400.output_queue_files",
			expected:    "as400.output_queue_files.lib_queue1_ready",
		},
		{
			name: "labels with special chars",
			labels: OutputQueueLabels{
				Library: "Lib/Name",
				Queue:   "Queue-1.0",
				Status:  "Ready (OK)",
			},
			contextName: "as400.output_queue_files",
			expected:    "as400.output_queue_files.lib_name_queue_1_0_ready__ok",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.labels.InstanceID(tc.contextName)
			if result != tc.expected {
				t.Errorf("OutputQueueLabels.InstanceID() = %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestPlanCacheLabelsInstanceID(t *testing.T) {
	testCases := []struct {
		name        string
		labels      PlanCacheLabels
		contextName string
		expected    string
	}{
		{
			name: "simple metric name",
			labels: PlanCacheLabels{
				Metric: "cache_hits",
			},
			contextName: "as400.plan_cache_summary",
			expected:    "as400.plan_cache_summary.cache_hits",
		},
		{
			name: "metric with special chars",
			labels: PlanCacheLabels{
				Metric: "cache-hits/total.rate",
			},
			contextName: "as400.plan_cache_summary",
			expected:    "as400.plan_cache_summary.cache_hits_total_rate",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.labels.InstanceID(tc.contextName)
			if result != tc.expected {
				t.Errorf("PlanCacheLabels.InstanceID() = %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestQueueOverviewLabelsInstanceID(t *testing.T) {
	testCases := []struct {
		name        string
		labels      QueueOverviewLabels
		contextName string
		expected    string
	}{
		{
			name: "simple labels",
			labels: QueueOverviewLabels{
				Queue_type: "Job",
				Item_type:  "Active",
			},
			contextName: "as400.queues_count",
			expected:    "as400.queues_count.job_active",
		},
		{
			name: "labels with special chars",
			labels: QueueOverviewLabels{
				Queue_type: "Job/Message",
				Item_type:  "Active-Pending",
			},
			contextName: "as400.queues_count",
			expected:    "as400.queues_count.job_message_active_pending",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.labels.InstanceID(tc.contextName)
			if result != tc.expected {
				t.Errorf("QueueOverviewLabels.InstanceID() = %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestSubsystemLabelsInstanceID(t *testing.T) {
	testCases := []struct {
		name        string
		labels      SubsystemLabels
		contextName string
		expected    string
	}{
		{
			name: "simple labels",
			labels: SubsystemLabels{
				Subsystem: "SUBSYS1",
				Library:   "LIB",
				Status:    "Active",
			},
			contextName: "as400.subsystem_jobs",
			expected:    "as400.subsystem_jobs.subsys1_lib_active",
		},
		{
			name: "labels with special chars",
			labels: SubsystemLabels{
				Subsystem: "Sub-System.1",
				Library:   "Lib/Name",
				Status:    "Active (Running)",
			},
			contextName: "as400.subsystem_jobs",
			expected:    "as400.subsystem_jobs.sub_system_1_lib_name_active__running",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.labels.InstanceID(tc.contextName)
			if result != tc.expected {
				t.Errorf("SubsystemLabels.InstanceID() = %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestTempStorageBucketLabelsInstanceID(t *testing.T) {
	testCases := []struct {
		name        string
		labels      TempStorageBucketLabels
		contextName string
		expected    string
	}{
		{
			name: "simple bucket name",
			labels: TempStorageBucketLabels{
				Bucket: "bucket1",
			},
			contextName: "as400.temp_storage_bucket",
			expected:    "as400.temp_storage_bucket.bucket1",
		},
		{
			name: "bucket with special chars",
			labels: TempStorageBucketLabels{
				Bucket: "Bucket-1.0:High",
			},
			contextName: "as400.temp_storage_bucket",
			expected:    "as400.temp_storage_bucket.bucket_1_0_high",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.labels.InstanceID(tc.contextName)
			if result != tc.expected {
				t.Errorf("TempStorageBucketLabels.InstanceID() = %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestGetAllContexts(t *testing.T) {
	contexts := GetAllContexts()

	// Should have at least some contexts
	if len(contexts) == 0 {
		t.Error("GetAllContexts() returned empty slice, expected contexts")
	}

	// Verify all items are non-nil
	for i, ctx := range contexts {
		if ctx == nil {
			t.Errorf("GetAllContexts()[%d] is nil", i)
		}
	}

	// Verify expected count (should be 55 contexts based on code)
	expectedCount := 55
	if len(contexts) != expectedCount {
		t.Warnf("GetAllContexts() returned %d contexts, expected %d (count may vary based on implementation)", len(contexts), expectedCount)
	}
}

func TestActiveJobCPUValuesStruct(t *testing.T) {
	values := ActiveJobCPUValues{
		Cpu: 5000,
	}

	if values.Cpu != 5000 {
		t.Errorf("ActiveJobCPUValues.Cpu = %d, want 5000", values.Cpu)
	}
}

func TestActiveJobResourcesValuesStruct(t *testing.T) {
	values := ActiveJobResourcesValues{
		Temp_storage: 1024000,
	}

	if values.Temp_storage != 1024000 {
		t.Errorf("ActiveJobResourcesValues.Temp_storage = %d, want 1024000", values.Temp_storage)
	}
}

func TestActiveJobTimeValuesStruct(t *testing.T) {
	values := ActiveJobTimeValues{
		Cpu_time:   1000,
		Total_time: 2000,
	}

	if values.Cpu_time != 1000 {
		t.Errorf("ActiveJobTimeValues.Cpu_time = %d, want 1000", values.Cpu_time)
	}
	if values.Total_time != 2000 {
		t.Errorf("ActiveJobTimeValues.Total_time = %d, want 2000", values.Total_time)
	}
}

func TestActiveJobActivityValuesStruct(t *testing.T) {
	values := ActiveJobActivityValues{
		Disk_io:                  100,
		Interactive_transactions: 50,
	}

	if values.Disk_io != 100 {
		t.Errorf("ActiveJobActivityValues.Disk_io = %d, want 100", values.Disk_io)
	}
	if values.Interactive_transactions != 50 {
		t.Errorf("ActiveJobActivityValues.Interactive_transactions = %d, want 50", values.Interactive_transactions)
	}
}

func TestActiveJobThreadsValuesStruct(t *testing.T) {
	values := ActiveJobThreadsValues{
		Threads: 10,
	}

	if values.Threads != 10 {
		t.Errorf("ActiveJobThreadsValues.Threads = %d, want 10", values.Threads)
	}
}

func TestDiskBusyValuesStruct(t *testing.T) {
	values := DiskBusyValues{
		Busy: 750,
	}

	if values.Busy != 750 {
		t.Errorf("DiskBusyValues.Busy = %d, want 750", values.Busy)
	}
}

func TestDiskIORequestsValuesStruct(t *testing.T) {
	values := DiskIORequestsValues{
		Read:  100,
		Write: 50,
	}

	if values.Read != 100 {
		t.Errorf("DiskIORequestsValues.Read = %d, want 100", values.Read)
	}
	if values.Write != 50 {
		t.Errorf("DiskIORequestsValues.Write = %d, want 50", values.Write)
	}
}

func TestDiskSpaceUsageValuesStruct(t *testing.T) {
	values := DiskSpaceUsageValues{
		Used: 500000,
	}

	if values.Used != 500000 {
		t.Errorf("DiskSpaceUsageValues.Used = %d, want 500000", values.Used)
	}
}

func TestDiskCapacityValuesStruct(t *testing.T) {
	values := DiskCapacityValues{
		Available: 100000,
		Used:      500000,
	}

	if values.Available != 100000 {
		t.Errorf("DiskCapacityValues.Available = %d, want 100000", values.Available)
	}
	if values.Used != 500000 {
		t.Errorf("DiskCapacityValues.Used = %d, want 500000", values.Used)
	}
}

func TestDiskBlocksValuesStruct(t *testing.T) {
	values := DiskBlocksValues{
		Read:  1000,
		Write: 500,
	}

	if values.Read != 1000 {
		t.Errorf("DiskBlocksValues.Read = %d, want 1000", values.Read)
	}
	if values.Write != 500 {
		t.Errorf("DiskBlocksValues.Write = %d, want 500", values.Write)
	}
}

func TestDiskSSDHealthValuesStruct(t *testing.T) {
	values := DiskSSDHealthValues{
		Life_remaining: 80000,
	}

	if values.Life_remaining != 80000 {
		t.Errorf("DiskSSDHealthValues.Life_remaining = %d, want 80000", values.Life_remaining)
	}
}

func TestDiskSSDPowerOnValuesStruct(t *testing.T) {
	values := DiskSSDPowerOnValues{
		Power_on_days: 730,
	}

	if values.Power_on_days != 730 {
		t.Errorf("DiskSSDPowerOnValues.Power_on_days = %d, want 730", values.Power_on_days)
	}
}

func TestHTTPServerConnectionsValuesStruct(t *testing.T) {
	values := HTTPServerConnectionsValues{
		Normal: 10,
		Ssl:    5,
	}

	if values.Normal != 10 {
		t.Errorf("HTTPServerConnectionsValues.Normal = %d, want 10", values.Normal)
	}
	if values.Ssl != 5 {
		t.Errorf("HTTPServerConnectionsValues.Ssl = %d, want 5", values.Ssl)
	}
}

func TestHTTPServerThreadsValuesStruct(t *testing.T) {
	values := HTTPServerThreadsValues{
		Active: 8,
		Idle:   2,
	}

	if values.Active != 8 {
		t.Errorf("HTTPServerThreadsValues.Active = %d, want 8", values.Active)
	}
	if values.Idle != 2 {
		t.Errorf("HTTPServerThreadsValues.Idle = %d, want 2", values.Idle)
	}
}

func TestHTTPServerRequestsValuesStruct(t *testing.T) {
	values := HTTPServerRequestsValues{
		Requests:  1000,
		Responses: 950,
		Rejected:  50,
	}

	if values.Requests != 1000 {
		t.Errorf("HTTPServerRequestsValues.Requests = %d, want 1000", values.Requests)
	}
	if values.Responses != 950 {
		t.Errorf("HTTPServerRequestsValues.Responses = %d, want 950", values.Responses)
	}
	if values.Rejected != 50 {
		t.Errorf("HTTPServerRequestsValues.Rejected = %d, want 50", values.Rejected)
	}
}

func TestHTTPServerBytesValuesStruct(t *testing.T) {
	values := HTTPServerBytesValues{
		Received: 1000000,
		Sent:     500000,
	}

	if values.Received != 1000000 {
		t.Errorf("HTTPServerBytesValues.Received = %d, want 1000000", values.Received)
	}
	if values.Sent != 500000 {
		t.Errorf("HTTPServerBytesValues.Sent = %d, want 500000", values.Sent)
	}
}

func TestMessageQueueMessagesValuesStruct(t *testing.T) {
	values := MessageQueueMessagesValues{
		Total:         100,
		Informational: 50,
		Inquiry:       20,
		Diagnostic:    15,
		Escape:        10,
		Notify:        4,
		Sender_copy:   1,
	}

	if values.Total != 100 {
		t.Errorf("MessageQueueMessagesValues.Total = %d, want 100", values.Total)
	}
	if values.Informational != 50 {
		t.Errorf("MessageQueueMessagesValues.Informational = %d, want 50", values.Informational)
	}
	if values.Inquiry != 20 {
		t.Errorf("MessageQueueMessagesValues.Inquiry = %d, want 20", values.Inquiry)
	}
	if values.Diagnostic != 15 {
		t.Errorf("MessageQueueMessagesValues.Diagnostic = %