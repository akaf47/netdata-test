package dnsquery

import (
	"strings"
	"testing"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDNSServerCharts(t *testing.T) {
	tests := []struct {
		name                string
		server              string
		network             string
		rtype               string
		expectedChartCount  int
		expectedDimsInChart int
	}{
		{
			name:                "basic dns query",
			server:              "8.8.8.8",
			network:             "udp",
			rtype:               "A",
			expectedChartCount:  2,
			expectedDimsInChart: 3,
		},
		{
			name:                "with dots in server",
			server:              "1.1.1.1",
			network:             "tcp",
			rtype:               "AAAA",
			expectedChartCount:  2,
			expectedDimsInChart: 3,
		},
		{
			name:                "ipv6 address",
			server:              "2001:4860:4860::8888",
			network:             "udp",
			rtype:               "MX",
			expectedChartCount:  2,
			expectedDimsInChart: 3,
		},
		{
			name:                "localhost",
			server:              "localhost",
			network:             "udp",
			rtype:               "NS",
			expectedChartCount:  2,
			expectedDimsInChart: 3,
		},
		{
			name:                "short server name",
			server:              "a",
			network:             "tcp",
			rtype:               "TXT",
			expectedChartCount:  2,
			expectedDimsInChart: 3,
		},
		{
			name:                "empty strings",
			server:              "",
			network:             "",
			rtype:               "",
			expectedChartCount:  2,
			expectedDimsInChart: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charts := newDNSServerCharts(tt.server, tt.network, tt.rtype)

			assert.NotNil(t, charts)
			assert.Len(t, *charts, tt.expectedChartCount)

			// Validate first chart (status chart)
			statusChart := (*charts)[0]
			assert.Contains(t, statusChart.ID, strings.ReplaceAll(tt.server, ".", "_"))
			assert.Contains(t, statusChart.ID, tt.rtype)
			assert.Equal(t, "DNS Query Status", statusChart.Title)
			assert.Equal(t, "status", statusChart.Units)
			assert.Equal(t, "query status", statusChart.Fam)
			assert.Equal(t, "dns_query.query_status", statusChart.Ctx)
			assert.Len(t, statusChart.Dims, tt.expectedDimsInChart)

			// Validate labels
			assert.Len(t, statusChart.Labels, 3)
			assert.Equal(t, "server", statusChart.Labels[0].Key)
			assert.Equal(t, tt.server, statusChart.Labels[0].Value)
			assert.Equal(t, "network", statusChart.Labels[1].Key)
			assert.Equal(t, tt.network, statusChart.Labels[1].Value)
			assert.Equal(t, "record_type", statusChart.Labels[2].Key)
			assert.Equal(t, tt.rtype, statusChart.Labels[2].Value)

			// Validate status chart dimensions
			assert.Equal(t, "success", statusChart.Dims[0].Name)
			assert.Equal(t, "network_error", statusChart.Dims[1].Name)
			assert.Equal(t, "dns_error", statusChart.Dims[2].Name)

			// Validate second chart (time chart)
			timeChart := (*charts)[1]
			assert.Contains(t, timeChart.ID, strings.ReplaceAll(tt.server, ".", "_"))
			assert.Contains(t, timeChart.ID, tt.rtype)
			assert.Equal(t, "DNS Query Time", timeChart.Title)
			assert.Equal(t, "seconds", timeChart.Units)
			assert.Equal(t, "query time", timeChart.Fam)
			assert.Equal(t, "dns_query.query_time", timeChart.Ctx)
			assert.Len(t, timeChart.Dims, 1)
			assert.Equal(t, "query_time", timeChart.Dims[0].Name)
			assert.Equal(t, int64(1e9), timeChart.Dims[0].Div)

			// Validate time chart labels
			assert.Len(t, timeChart.Labels, 3)
			assert.Equal(t, "server", timeChart.Labels[0].Key)
			assert.Equal(t, tt.server, timeChart.Labels[0].Value)
		})
	}
}

func TestNewDNSServerChartsChartIDs(t *testing.T) {
	tests := []struct {
		name               string
		server             string
		rtype              string
		expectedStatusID   string
		expectedTimeID     string
	}{
		{
			name:               "simple values",
			server:             "8.8.8.8",
			rtype:              "A",
			expectedStatusID:   "server_8_8_8_8_record_A_query_status",
			expectedTimeID:     "server_8_8_8_8_record_A_query_time",
		},
		{
			name:               "complex server",
			server:             "dns.google.com",
			rtype:              "AAAA",
			expectedStatusID:   "server_dns_google_com_record_AAAA_query_status",
			expectedTimeID:     "server_dns_google_com_record_AAAA_query_time",
		},
		{
			name:               "localhost",
			server:             "localhost",
			rtype:              "MX",
			expectedStatusID:   "server_localhost_record_MX_query_status",
			expectedTimeID:     "server_localhost_record_MX_query_time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charts := newDNSServerCharts(tt.server, "", tt.rtype)
			require.NotNil(t, charts)
			require.Len(t, *charts, 2)

			assert.Equal(t, tt.expectedStatusID, (*charts)[0].ID)
			assert.Equal(t, tt.expectedTimeID, (*charts)[1].ID)
		})
	}
}

func TestNewDNSServerChartsDimensionIDs(t *testing.T) {
	server := "8.8.8.8"
	rtype := "A"
	
	charts := newDNSServerCharts(server, "udp", rtype)
	require.NotNil(t, charts)
	require.Len(t, *charts, 2)

	statusChart := (*charts)[0]
	// Verify dimension IDs are properly formatted
	assert.Contains(t, statusChart.Dims[0].ID, server)
	assert.Contains(t, statusChart.Dims[0].ID, rtype)
	assert.Contains(t, statusChart.Dims[1].ID, server)
	assert.Contains(t, statusChart.Dims[1].ID, rtype)
	assert.Contains(t, statusChart.Dims[2].ID, server)
	assert.Contains(t, statusChart.Dims[2].ID, rtype)

	timeChart := (*charts)[1]
	assert.Contains(t, timeChart.Dims[0].ID, server)
	assert.Contains(t, timeChart.Dims[0].ID, rtype)
}

func TestNewDNSServerChartsNotModifyingTemplate(t *testing.T) {
	// Test that calling newDNSServerCharts multiple times doesn't modify the template
	originalStatusID := dnsQueryStatusChartTmpl.ID
	originalTimeID := dnsQueryTimeChartTmpl.ID

	_ = newDNSServerCharts("1.1.1.1", "udp", "A")
	_ = newDNSServerCharts("8.8.8.8", "tcp", "AAAA")
	_ = newDNSServerCharts("dns.google.com", "udp", "MX")

	assert.Equal(t, originalStatusID, dnsQueryStatusChartTmpl.ID)
	assert.Equal(t, originalTimeID, dnsQueryTimeChartTmpl.ID)
}

func TestNewDNSServerChartsWithSpecialCharacters(t *testing.T) {
	// Test with server names containing various characters
	server := "ns1.example.co.uk"
	charts := newDNSServerCharts(server, "udp", "NS")
	
	require.NotNil(t, charts)
	require.Len(t, *charts, 2)

	// All dots should be replaced with underscores in chart ID
	assert.NotContains(t, (*charts)[0].ID, ".")
	assert.NotContains(t, (*charts)[1].ID, ".")
	assert.Contains(t, (*charts)[0].ID, "ns1_example_co_uk")
}

func TestChartsPriority(t *testing.T) {
	charts := newDNSServerCharts("8.8.8.8", "udp", "A")
	require.NotNil(t, charts)
	require.Len(t, *charts, 2)

	statusChart := (*charts)[0]
	timeChart := (*charts)[1]

	// Status chart should have higher priority
	assert.Less(t, statusChart.Priority, timeChart.Priority)
}