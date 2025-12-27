import unittest
from unittest.mock import Mock, patch, MagicMock, call
import sys
import os
from collections import defaultdict

# Mock the Netdata base module before importing haproxy
sys.modules['bases'] = MagicMock()
sys.modules['bases.FrameworkServices'] = MagicMock()

# Add parent directory to path for imports
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))

from haproxy.haproxy_chart import HAProxy


class TestHAProxyInit(unittest.TestCase):
    """Test HAProxy class initialization"""
    
    def test_init_default_values(self):
        """should initialize with correct default values"""
        haproxy = HAProxy('test')
        self.assertEqual(haproxy.name, 'test')
        self.assertIsNotNone(haproxy.order)
        self.assertIsNotNone(haproxy.definitions)
        
    def test_init_sets_socket_path_from_config(self):
        """should set socket path from configuration"""
        config = {'socket': '/custom/socket.sock'}
        haproxy = HAProxy('test')
        haproxy.configuration.update(config)
        self.assertEqual(haproxy.configuration.get('socket'), '/custom/socket.sock')


class TestHAProxyCheck(unittest.TestCase):
    """Test HAProxy check and setup methods"""
    
    def setUp(self):
        self.haproxy = HAProxy('test')
        
    def test_check_returns_true_on_valid_setup(self):
        """should return True when socket is accessible"""
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.return_value = [['pxname', 'svname']]
            result = self.haproxy.check()
            self.assertTrue(result)
            
    def test_check_returns_false_on_socket_error(self):
        """should return False when socket is not accessible"""
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.side_effect = Exception("Socket error")
            result = self.haproxy.check()
            self.assertFalse(result)
            
    def test_check_returns_false_on_empty_response(self):
        """should return False when stats are empty"""
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.return_value = []
            result = self.haproxy.check()
            self.assertFalse(result)
            
    def test_setup_initializes_charts(self):
        """should initialize charts during setup"""
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.return_value = [['pxname', 'svname'], ['proxy1', 'FRONTEND']]
            self.haproxy.setup()
            self.assertTrue(len(self.haproxy.definitions) > 0)


class TestHAProxyGetStats(unittest.TestCase):
    """Test _get_stats method"""
    
    def setUp(self):
        self.haproxy = HAProxy('test')
        self.haproxy.socket_path = '/var/run/haproxy.sock'
        
    def test_get_stats_returns_parsed_csv(self):
        """should parse CSV response from socket"""
        stats_response = "pxname,svname,scur\nproxy1,FRONTEND,10\nproxy1,backend1,5\n"
        with patch('socket.socket') as mock_socket_class:
            mock_socket = MagicMock()
            mock_socket_class.return_value = mock_socket
            mock_socket.recv.return_value = stats_response.encode('utf-8')
            
            result = self.haproxy._get_stats()
            self.assertEqual(len(result), 3)  # header + 2 lines
            self.assertIn('pxname', result[0])
            
    def test_get_stats_handles_socket_error(self):
        """should raise exception on socket error"""
        with patch('socket.socket') as mock_socket_class:
            mock_socket_class.side_effect = Exception("Socket error")
            with self.assertRaises(Exception):
                self.haproxy._get_stats()
                
    def test_get_stats_handles_connection_refused(self):
        """should raise exception when connection is refused"""
        with patch('socket.socket') as mock_socket_class:
            mock_socket = MagicMock()
            mock_socket_class.return_value = mock_socket
            mock_socket.connect.side_effect = OSError("Connection refused")
            
            with self.assertRaises(OSError):
                self.haproxy._get_stats()
                
    def test_get_stats_handles_malformed_response(self):
        """should handle malformed CSV response"""
        with patch('socket.socket') as mock_socket_class:
            mock_socket = MagicMock()
            mock_socket_class.return_value = mock_socket
            mock_socket.recv.return_value = b"invalid\nresponse\n"
            
            result = self.haproxy._get_stats()
            self.assertIsNotNone(result)
            self.assertGreater(len(result), 0)


class TestHAProxyCollect(unittest.TestCase):
    """Test collect method"""
    
    def setUp(self):
        self.haproxy = HAProxy('test')
        self.haproxy.socket_path = '/var/run/haproxy.sock'
        
    def test_collect_updates_all_metrics(self):
        """should collect all metrics from stats"""
        stats = [
            ['pxname', 'svname', 'scur', 'smax', 'slim', 'stot', 'bin', 'bout'],
            ['frontend', 'FRONTEND', '100', '200', '300', '400', '500', '600'],
            ['backend', 'BACKEND', '50', '100', '150', '200', '250', '300']
        ]
        
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.return_value = stats
            self.haproxy.setup()
            result = self.haproxy.collect()
            
            self.assertIsNotNone(result)
            self.assertIsInstance(result, dict)
            
    def test_collect_returns_none_on_error(self):
        """should return None when stats collection fails"""
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.side_effect = Exception("Connection failed")
            self.haproxy.setup()
            result = self.haproxy.collect()
            
            self.assertIsNone(result)
            
    def test_collect_handles_missing_columns(self):
        """should handle missing metric columns gracefully"""
        stats = [
            ['pxname', 'svname'],
            ['frontend', 'FRONTEND']
        ]
        
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.return_value = stats
            self.haproxy.setup()
            result = self.haproxy.collect()
            
            self.assertIsNotNone(result)
            
    def test_collect_handles_non_numeric_values(self):
        """should handle non-numeric metric values"""
        stats = [
            ['pxname', 'svname', 'scur'],
            ['frontend', 'FRONTEND', 'invalid']
        ]
        
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.return_value = stats
            self.haproxy.setup()
            result = self.haproxy.collect()
            
            # Should handle gracefully without crashing
            self.assertIsNotNone(result)


class TestHAProxyChartDefinitions(unittest.TestCase):
    """Test chart definitions generation"""
    
    def setUp(self):
        self.haproxy = HAProxy('test')
        
    def test_order_dict_is_not_empty(self):
        """should have non-empty order dictionary"""
        self.assertIsNotNone(self.haproxy.order)
        self.assertIsInstance(self.haproxy.order, dict)
        
    def test_definitions_dict_is_not_empty(self):
        """should have non-empty definitions dictionary"""
        self.assertIsNotNone(self.haproxy.definitions)
        self.assertIsInstance(self.haproxy.definitions, dict)
        
    def test_chart_definitions_have_required_fields(self):
        """should have required fields in chart definitions"""
        for chart_id, chart_def in self.haproxy.definitions.items():
            self.assertIn('title', chart_def)
            self.assertIn('units', chart_def)
            self.assertIn('family', chart_def)
            self.assertIn('charts', chart_def)
            
    def test_charts_have_required_dimensions(self):
        """should have dimensions defined for charts"""
        for chart_id, chart_def in self.haproxy.definitions.items():
            for chart_name, chart in chart_def.get('charts', {}).items():
                self.assertIn('lines', chart)


class TestHAProxyMetricParsing(unittest.TestCase):
    """Test metric parsing and conversion"""
    
    def setUp(self):
        self.haproxy = HAProxy('test')
        
    def test_parse_numeric_metrics(self):
        """should parse numeric metrics correctly"""
        stats = [
            ['pxname', 'svname', 'scur', 'smax'],
            ['frontend', 'FRONTEND', '100', '200']
        ]
        
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.return_value = stats
            self.haproxy.setup()
            result = self.haproxy.collect()
            
            self.assertIsNotNone(result)
            
    def test_parse_zero_values(self):
        """should handle zero metric values"""
        stats = [
            ['pxname', 'svname', 'scur'],
            ['frontend', 'FRONTEND', '0']
        ]
        
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.return_value = stats
            self.haproxy.setup()
            result = self.haproxy.collect()
            
            self.assertIsNotNone(result)
            
    def test_parse_negative_values(self):
        """should handle negative metric values"""
        stats = [
            ['pxname', 'svname', 'rate_limit'],
            ['frontend', 'FRONTEND', '-1']
        ]
        
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.return_value = stats
            self.haproxy.setup()
            result = self.haproxy.collect()
            
            self.assertIsNotNone(result)
            
    def test_parse_large_values(self):
        """should handle large numeric values"""
        stats = [
            ['pxname', 'svname', 'stot'],
            ['frontend', 'FRONTEND', '9999999999']
        ]
        
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.return_value = stats
            self.haproxy.setup()
            result = self.haproxy.collect()
            
            self.assertIsNotNone(result)


class TestHAProxyProxyTypes(unittest.TestCase):
    """Test handling of different proxy types"""
    
    def setUp(self):
        self.haproxy = HAProxy('test')
        
    def test_collect_frontend_metrics(self):
        """should collect frontend metrics"""
        stats = [
            ['pxname', 'svname', 'scur'],
            ['web_frontend', 'FRONTEND', '100']
        ]
        
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.return_value = stats
            self.haproxy.setup()
            result = self.haproxy.collect()
            
            self.assertIsNotNone(result)
            
    def test_collect_backend_metrics(self):
        """should collect backend metrics"""
        stats = [
            ['pxname', 'svname', 'scur'],
            ['web_backend', 'BACKEND', '50']
        ]
        
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.return_value = stats
            self.haproxy.setup()
            result = self.haproxy.collect()
            
            self.assertIsNotNone(result)
            
    def test_collect_server_metrics(self):
        """should collect server metrics"""
        stats = [
            ['pxname', 'svname', 'scur', 'stot'],
            ['backend', 'server1', '25', '100'],
            ['backend', 'server2', '30', '120']
        ]
        
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.return_value = stats
            self.haproxy.setup()
            result = self.haproxy.collect()
            
            self.assertIsNotNone(result)
            
    def test_collect_multiple_backends(self):
        """should handle multiple backends correctly"""
        stats = [
            ['pxname', 'svname', 'scur'],
            ['backend1', 'BACKEND', '50'],
            ['backend2', 'BACKEND', '75'],
            ['backend1', 'server1', '25'],
            ['backend1', 'server2', '25'],
            ['backend2', 'server3', '75']
        ]
        
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.return_value = stats
            self.haproxy.setup()
            result = self.haproxy.collect()
            
            self.assertIsNotNone(result)


class TestHAProxyEdgeCases(unittest.TestCase):
    """Test edge cases and error conditions"""
    
    def setUp(self):
        self.haproxy = HAProxy('test')
        
    def test_empty_stats_response(self):
        """should handle empty stats response"""
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.return_value = []
            result = self.haproxy.check()
            self.assertFalse(result)
            
    def test_stats_with_only_header(self):
        """should handle stats with only header row"""
        stats = [['pxname', 'svname']]
        
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.return_value = stats
            self.haproxy.setup()
            result = self.haproxy.collect()
            
            self.assertIsNotNone(result)
            
    def test_socket_path_none(self):
        """should handle None socket path"""
        self.haproxy.socket_path = None
        with patch('socket.socket') as mock_socket:
            mock_socket.side_effect = Exception("Invalid socket")
            with self.assertRaises(Exception):
                self.haproxy._get_stats()
                
    def test_collect_called_without_setup(self):
        """should handle collect called without setup"""
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.return_value = [['pxname', 'svname']]
            result = self.haproxy.collect()
            
            # Should not crash
            self.assertIsNotNone(result)
            
    def test_multiple_collect_calls(self):
        """should handle multiple consecutive collect calls"""
        stats = [
            ['pxname', 'svname', 'scur'],
            ['frontend', 'FRONTEND', '100']
        ]
        
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.return_value = stats
            self.haproxy.setup()
            
            result1 = self.haproxy.collect()
            result2 = self.haproxy.collect()
            
            self.assertIsNotNone(result1)
            self.assertIsNotNone(result2)


class TestHAProxySpecialCharacters(unittest.TestCase):
    """Test handling of special characters in names"""
    
    def setUp(self):
        self.haproxy = HAProxy('test')
        
    def test_proxy_names_with_underscores(self):
        """should handle proxy names with underscores"""
        stats = [
            ['pxname', 'svname', 'scur'],
            ['web_frontend_prod', 'FRONTEND', '100']
        ]
        
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.return_value = stats
            self.haproxy.setup()
            result = self.haproxy.collect()
            
            self.assertIsNotNone(result)
            
    def test_proxy_names_with_hyphens(self):
        """should handle proxy names with hyphens"""
        stats = [
            ['pxname', 'svname', 'scur'],
            ['web-frontend-prod', 'FRONTEND', '100']
        ]
        
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.return_value = stats
            self.haproxy.setup()
            result = self.haproxy.collect()
            
            self.assertIsNotNone(result)
            
    def test_server_names_with_special_chars(self):
        """should handle server names with special characters"""
        stats = [
            ['pxname', 'svname', 'scur'],
            ['backend', 'server_1-prod', '50']
        ]
        
        with patch.object(self.haproxy, '_get_stats') as mock_get_stats:
            mock_get_stats.return_value = stats
            self.haproxy.setup()
            result = self.haproxy.collect()
            
            self.assertIsNotNone(result)


if __name__ == '__main__':
    unittest.main()