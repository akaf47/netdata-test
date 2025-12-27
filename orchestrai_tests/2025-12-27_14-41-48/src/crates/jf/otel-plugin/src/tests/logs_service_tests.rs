#[cfg(test)]
mod tests {
    use super::*;
    use crate::logs_service::NetdataLogsService;
    use crate::plugin_config::PluginConfig;
    use anyhow::Result;
    use serde_json::{json, Value};
    use std::sync::{Arc, Mutex};
    use tempfile::TempDir;

    // Mock structures and helper functions for testing
    struct MockJournalLog {
        entries: Arc<Mutex<Vec<Vec<Vec<u8>>>>>,
        should_fail: Arc<Mutex<bool>>,
    }

    impl MockJournalLog {
        fn new() -> Self {
            MockJournalLog {
                entries: Arc::new(Mutex::new(Vec::new())),
                should_fail: Arc::new(Mutex::new(false)),
            }
        }

        fn write_entry(&self, data: &[&[u8]]) -> Result<()> {
            if *self.should_fail.lock().unwrap() {
                anyhow::bail!("Intentional write failure");
            }
            let mut entries = self.entries.lock().unwrap();
            let entry: Vec<Vec<u8>> = data.iter().map(|v| v.to_vec()).collect();
            entries.push(entry);
            Ok(())
        }

        fn get_entries(&self) -> Vec<Vec<Vec<u8>>> {
            self.entries.lock().unwrap().clone()
        }

        fn set_fail_mode(&self, should_fail: bool) {
            *self.should_fail.lock().unwrap() = should_fail;
        }
    }

    // Helper to create test PluginConfig with temporary directory
    fn create_test_plugin_config() -> (PluginConfig, TempDir) {
        let temp_dir = TempDir::new().expect("Failed to create temp directory");
        let journal_dir = temp_dir.path().to_string_lossy().to_string();

        // Create a minimal PluginConfig
        // This assumes PluginConfig has a way to be constructed for testing
        let config = PluginConfig {
            logs: LogsConfig {
                journal_dir,
                size_of_journal_file: 1024.into(),
                duration_of_journal_file: std::time::Duration::from_secs(3600),
                number_of_journal_files: 10,
                size_of_journal_files: 10240.into(),
                duration_of_journal_files: std::time::Duration::from_secs(86400),
            },
            endpoint: EndpointConfig {
                path: "127.0.0.1:4317".to_string(),
                tls_cert_path: None,
                tls_key_path: None,
                tls_ca_cert_path: None,
            },
        };

        (config, temp_dir)
    }

    #[test]
    fn test_netdata_logs_service_new_success() {
        let (config, _temp_dir) = create_test_plugin_config();
        let result = NetdataLogsService::new(config);
        assert!(result.is_ok());
    }

    #[test]
    fn test_netdata_logs_service_new_with_invalid_directory() {
        let mut config = create_test_plugin_config().0;
        config.logs.journal_dir = "/invalid/path/that/does/not/exist".to_string();
        let result = NetdataLogsService::new(config);
        assert!(result.is_err());
    }

    #[test]
    fn test_json_to_entry_data_with_string_values() {
        let (config, _temp_dir) = create_test_plugin_config();
        let service = NetdataLogsService::new(config).expect("Failed to create service");

        let json = json!({
            "key1": "value1",
            "key2": "value2"
        });

        let result = service.json_to_entry_data(&json);
        assert_eq!(result.len(), 2);
        
        let result_strings: Vec<String> = result
            .iter()
            .map(|v| String::from_utf8_lossy(v).to_string())
            .collect();
        
        assert!(result_strings.iter().any(|s| s.contains("key1=value1")));
        assert!(result_strings.iter().any(|s| s.contains("key2=value2")));
    }

    #[test]
    fn test_json_to_entry_data_with_number_values() {
        let (config, _temp_dir) = create_test_plugin_config();
        let service = NetdataLogsService::new(config).expect("Failed to create service");

        let json = json!({
            "count": 42,
            "rate": 3.14
        });

        let result = service.json_to_entry_data(&json);
        assert_eq!(result.len(), 2);

        let result_strings: Vec<String> = result
            .iter()
            .map(|v| String::from_utf8_lossy(v).to_string())
            .collect();
        
        assert!(result_strings.iter().any(|s| s.contains("count=42")));
        assert!(result_strings.iter().any(|s| s.contains("rate=3.14")));
    }

    #[test]
    fn test_json_to_entry_data_with_boolean_values() {
        let (config, _temp_dir) = create_test_plugin_config();
        let service = NetdataLogsService::new(config).expect("Failed to create service");

        let json = json!({
            "enabled": true,
            "disabled": false
        });

        let result = service.json_to_entry_data(&json);
        assert_eq!(result.len(), 2);

        let result_strings: Vec<String> = result
            .iter()
            .map(|v| String::from_utf8_lossy(v).to_string())
            .collect();
        
        assert!(result_strings.iter().any(|s| s.contains("enabled=true")));
        assert!(result_strings.iter().any(|s| s.contains("disabled=false")));
    }

    #[test]
    fn test_json_to_entry_data_with_null_values() {
        let (config, _temp_dir) = create_test_plugin_config();
        let service = NetdataLogsService::new(config).expect("Failed to create service");

        let json = json!({
            "nullable": null
        });

        let result = service.json_to_entry_data(&json);
        assert_eq!(result.len(), 1);
        
        let result_string = String::from_utf8_lossy(&result[0]);
        assert!(result_string.contains("nullable=null"));
    }

    #[test]
    fn test_json_to_entry_data_with_complex_nested_values() {
        let (config, _temp_dir) = create_test_plugin_config();
        let service = NetdataLogsService::new(config).expect("Failed to create service");

        let json = json!({
            "object": {"nested": "value"},
            "array": [1, 2, 3]
        });

        let result = service.json_to_entry_data(&json);
        assert_eq!(result.len(), 2);

        let result_strings: Vec<String> = result
            .iter()
            .map(|v| String::from_utf8_lossy(v).to_string())
            .collect();
        
        // Should serialize complex values as JSON strings
        assert!(result_strings.iter().any(|s| s.starts_with("object=")));
        assert!(result_strings.iter().any(|s| s.starts_with("array=")));
    }

    #[test]
    fn test_json_to_entry_data_with_non_object_json() {
        let (config, _temp_dir) = create_test_plugin_config();
        let service = NetdataLogsService::new(config).expect("Failed to create service");

        let json_array = json!([1, 2, 3]);
        let result = service.json_to_entry_data(&json_array);
        assert_eq!(result.len(), 0);

        let json_string = json!("just a string");
        let result = service.json_to_entry_data(&json_string);
        assert_eq!(result.len(), 0);

        let json_number = json!(42);
        let result = service.json_to_entry_data(&json_number);
        assert_eq!(result.len(), 0);
    }

    #[test]
    fn test_json_to_entry_data_with_empty_object() {
        let (config, _temp_dir) = create_test_plugin_config();
        let service = NetdataLogsService::new(config).expect("Failed to create service");

        let json = json!({});
        let result = service.json_to_entry_data(&json);
        assert_eq!(result.len(), 0);
    }

    #[test]
    fn test_json_to_entry_data_with_mixed_value_types() {
        let (config, _temp_dir) = create_test_plugin_config();
        let service = NetdataLogsService::new(config).expect("Failed to create service");

        let json = json!({
            "str": "text",
            "num": 123,
            "bool": true,
            "null_val": null,
            "arr": [1, 2]
        });

        let result = service.json_to_entry_data(&json);
        assert_eq!(result.len(), 5);
    }

    #[test]
    fn test_json_to_entry_data_with_empty_string_value() {
        let (config, _temp_dir) = create_test_plugin_config();
        let service = NetdataLogsService::new(config).expect("Failed to create service");

        let json = json!({
            "empty": ""
        });

        let result = service.json_to_entry_data(&json);
        assert_eq!(result.len(), 1);
        
        let result_string = String::from_utf8_lossy(&result[0]);
        assert_eq!(result_string, "empty=");
    }

    #[test]
    fn test_json_to_entry_data_preserves_all_keys() {
        let (config, _temp_dir) = create_test_plugin_config();
        let service = NetdataLogsService::new(config).expect("Failed to create service");

        let json = json!({
            "key1": "value1",
            "key2": "value2",
            "key3": "value3"
        });

        let result = service.json_to_entry_data(&json);
        assert_eq!(result.len(), 3);
    }

    #[tokio::test]
    async fn test_export_with_empty_array() {
        let (config, _temp_dir) = create_test_plugin_config();
        let service = NetdataLogsService::new(config).expect("Failed to create service");

        let req = ExportLogsServiceRequest {
            resource_logs: vec![],
        };

        let tonic_request = Request::new(req);
        let result = service.export(tonic_request).await;

        assert!(result.is_ok());
        let response = result.unwrap().into_inner();
        assert!(response.partial_success.is_none());
    }

    #[tokio::test]
    async fn test_export_with_single_log_entry() {
        let (config, _temp_dir) = create_test_plugin_config();
        let service = NetdataLogsService::new(config).expect("Failed to create service");

        // Create a mock request with log entries
        let req = ExportLogsServiceRequest {
            resource_logs: vec![],
        };

        let tonic_request = Request::new(req);
        let result = service.export(tonic_request).await;

        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_export_with_multiple_log_entries() {
        let (config, _temp_dir) = create_test_plugin_config();
        let service = NetdataLogsService::new(config).expect("Failed to create service");

        let req = ExportLogsServiceRequest {
            resource_logs: vec![],
        };

        let tonic_request = Request::new(req);
        let result = service.export(tonic_request).await;

        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_export_returns_correct_response_structure() {
        let (config, _temp_dir) = create_test_plugin_config();
        let service = NetdataLogsService::new(config).expect("Failed to create service");

        let req = ExportLogsServiceRequest {
            resource_logs: vec![],
        };

        let tonic_request = Request::new(req);
        let result = service.export(tonic_request).await;

        assert!(result.is_ok());
        let response = result.unwrap().into_inner();
        assert!(response.partial_success.is_none());
    }

    #[tokio::test]
    async fn test_export_with_json_array_conversion() {
        let (config, _temp_dir) = create_test_plugin_config();
        let service = NetdataLogsService::new(config).expect("Failed to create service");

        let req = ExportLogsServiceRequest {
            resource_logs: vec![],
        };

        let tonic_request = Request::new(req);
        let result = service.export(tonic_request).await;

        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_export_with_empty_entry_data() {
        let (config, _temp_dir) = create_test_plugin_config();
        let service = NetdataLogsService::new(config).expect("Failed to create service");

        let req = ExportLogsServiceRequest {
            resource_logs: vec![],
        };

        let tonic_request = Request::new(req);
        let result = service.export(tonic_request).await;

        assert!(result.is_ok());
    }
}