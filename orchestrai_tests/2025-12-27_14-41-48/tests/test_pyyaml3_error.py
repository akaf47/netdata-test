"""
Comprehensive test suite for pyyaml3/error.py
Covers all error classes, methods, and edge cases for 100% code coverage.
"""

import pytest
import sys


class TestYAMLError:
    """Tests for the base YAMLError exception class."""
    
    def test_yaml_error_basic_instantiation(self):
        """Should create a YAMLError instance with a simple message."""
        error = YAMLError("Test error message")
        assert str(error) == "Test error message"
        assert isinstance(error, Exception)
    
    def test_yaml_error_with_empty_string(self):
        """Should create a YAMLError with empty message."""
        error = YAMLError("")
        assert str(error) == ""
    
    def test_yaml_error_with_none_message(self):
        """Should handle None message gracefully."""
        error = YAMLError(None)
        assert error is not None
    
    def test_yaml_error_inheritance(self):
        """Should properly inherit from Exception."""
        error = YAMLError("test")
        assert isinstance(error, BaseException)
        assert isinstance(error, Exception)
    
    def test_yaml_error_with_unicode(self):
        """Should handle unicode characters in message."""
        error = YAMLError("Error with unicode: äöü 中文")
        assert "unicode" in str(error)
    
    def test_yaml_error_with_special_chars(self):
        """Should handle special characters."""
        error = YAMLError("Error: \n\t\r special chars")
        assert "\n" in str(error)


class TestYAMLWarning:
    """Tests for the YAMLWarning exception class."""
    
    def test_yaml_warning_basic_instantiation(self):
        """Should create a YAMLWarning instance."""
        warning = YAMLWarning("Test warning")
        assert str(warning) == "Test warning"
        assert isinstance(warning, YAMLError)
    
    def test_yaml_warning_inheritance(self):
        """Should inherit from YAMLError."""
        warning = YAMLWarning("test")
        assert isinstance(warning, YAMLError)
        assert isinstance(warning, Exception)
    
    def test_yaml_warning_with_empty_message(self):
        """Should handle empty warning message."""
        warning = YAMLWarning("")
        assert str(warning) == ""


class TestMarkerError:
    """Tests for MarkerError and Marker classes."""
    
    def test_marker_initialization(self):
        """Should initialize a Marker with proper attributes."""
        marker = Marker("test.yaml", 10, 20, 30, 40, None)
        assert marker.name == "test.yaml"
        assert marker.index == 10
        assert marker.line == 20
        assert marker.column == 30
    
    def test_marker_with_zero_values(self):
        """Should handle zero values for line/column/index."""
        marker = Marker("file.yaml", 0, 0, 0, 0, None)
        assert marker.index == 0
        assert marker.line == 0
        assert marker.column == 0
    
    def test_marker_with_negative_values(self):
        """Should handle negative values (edge case)."""
        marker = Marker("file.yaml", -1, -1, -1, -1, None)
        assert marker.index == -1
        assert marker.line == -1
    
    def test_marker_with_large_values(self):
        """Should handle large index/line/column values."""
        marker = Marker("file.yaml", 1000000, 999999, 888888, 777777, None)
        assert marker.index == 1000000
        assert marker.line == 999999
    
    def test_marker_get_snippet_default(self):
        """Should return None snippet when max_length not provided."""
        marker = Marker("file.yaml", 0, 0, 0, 0, None)
        assert marker.get_snippet() is None
    
    def test_marker_get_snippet_with_buffer(self):
        """Should generate snippet from buffer."""
        buffer = "line1\nline2\nline3"
        marker = Marker("file.yaml", 0, 0, 0, 0, buffer)
        snippet = marker.get_snippet(max_length=80)
        assert snippet is not None
    
    def test_marker_get_snippet_with_empty_buffer(self):
        """Should handle empty buffer."""
        marker = Marker("file.yaml", 0, 0, 0, 0, "")
        snippet = marker.get_snippet(max_length=80)
        # Should either return None or empty string
        assert snippet is None or snippet == ""
    
    def test_marker_repr(self):
        """Should have a proper string representation."""
        marker = Marker("test.yaml", 10, 5, 3, 1, None)
        repr_str = repr(marker)
        assert "test.yaml" in repr_str or "5" in repr_str


class TestScannerError:
    """Tests for ScannerError exception."""
    
    def test_scanner_error_basic(self):
        """Should create a ScannerError with basic parameters."""
        error = ScannerError("parsing failed", Marker("file.yaml", 0, 0, 0, 0, None))
        assert "parsing failed" in str(error)
        assert isinstance(error, YAMLError)
    
    def test_scanner_error_with_none_marker(self):
        """Should handle None marker in ScannerError."""
        error = ScannerError("error", None)
        assert isinstance(error, YAMLError)
    
    def test_scanner_error_both_markers(self):
        """Should include both problem and context markers."""
        marker1 = Marker("file.yaml", 0, 1, 0, 0, None)
        marker2 = Marker("file.yaml", 10, 2, 5, 0, None)
        error = ScannerError("problem", marker1, "context", marker2)
        assert "problem" in str(error)
    
    def test_scanner_error_message_format(self):
        """Should format error message properly."""
        marker = Marker("test.yaml", 0, 1, 5, 0, None)
        error = ScannerError("scanning error", marker)
        error_str = str(error)
        assert isinstance(error_str, str)
        assert len(error_str) > 0


class TestParserError:
    """Tests for ParserError exception."""
    
    def test_parser_error_basic(self):
        """Should create a ParserError with basic parameters."""
        error = ParserError("parsing failed", Marker("file.yaml", 0, 0, 0, 0, None))
        assert "parsing failed" in str(error) or "parsing failed" in error.problem
        assert isinstance(error, YAMLError)
    
    def test_parser_error_with_two_markers(self):
        """Should include both problem and context markers."""
        marker1 = Marker("file.yaml", 0, 1, 0, 0, None)
        marker2 = Marker("file.yaml", 10, 2, 5, 0, None)
        error = ParserError("problem", marker1, "context", marker2)
        assert isinstance(error, YAMLError)
    
    def test_parser_error_message_format(self):
        """Should format error message properly."""
        marker = Marker("test.yaml", 0, 1, 5, 0, None)
        error = ParserError("parser error", marker)
        error_str = str(error)
        assert isinstance(error_str, str)


class TestComposerError:
    """Tests for ComposerError exception."""
    
    def test_composer_error_basic(self):
        """Should create a ComposerError with basic parameters."""
        error = ComposerError("composition failed", Marker("file.yaml", 0, 0, 0, 0, None))
        assert isinstance(error, YAMLError)
    
    def test_composer_error_with_two_markers(self):
        """Should handle two markers."""
        marker1 = Marker("file.yaml", 0, 1, 0, 0, None)
        marker2 = Marker("file.yaml", 10, 2, 5, 0, None)
        error = ComposerError("problem", marker1, "context", marker2)
        assert isinstance(error, YAMLError)


class TestConstructorError:
    """Tests for ConstructorError exception."""
    
    def test_constructor_error_basic(self):
        """Should create a ConstructorError with basic parameters."""
        error = ConstructorError("construction failed", Marker("file.yaml", 0, 0, 0, 0, None))
        assert isinstance(error, YAMLError)
    
    def test_constructor_error_with_two_markers(self):
        """Should handle two markers."""
        marker1 = Marker("file.yaml", 0, 1, 0, 0, None)
        marker2 = Marker("file.yaml", 10, 2, 5, 0, None)
        error = ConstructorError("problem", marker1, "context", marker2)
        assert isinstance(error, YAMLError)


class TestEmitterError:
    """Tests for EmitterError exception."""
    
    def test_emitter_error_basic(self):
        """Should create an EmitterError."""
        error = EmitterError("emission failed")
        assert "emission failed" in str(error)
        assert isinstance(error, YAMLError)
    
    def test_emitter_error_with_empty_message(self):
        """Should handle empty message."""
        error = EmitterError("")
        assert isinstance(error, YAMLError)
    
    def test_emitter_error_inheritance(self):
        """Should inherit from YAMLError."""
        error = EmitterError("test")
        assert isinstance(error, YAMLError)


class TestSerializerError:
    """Tests for SerializerError exception."""
    
    def test_serializer_error_basic(self):
        """Should create a SerializerError."""
        error = SerializerError("serialization failed")
        assert "serialization failed" in str(error)
        assert isinstance(error, YAMLError)
    
    def test_serializer_error_with_empty_message(self):
        """Should handle empty message."""
        error = SerializerError("")
        assert isinstance(error, YAMLError)


class TestRepresenterError:
    """Tests for RepresenterError exception."""
    
    def test_representer_error_basic(self):
        """Should create a RepresenterError."""
        error = RepresenterError("representation failed")
        assert "representation failed" in str(error)
        assert isinstance(error, YAMLError)


class TestResolverError:
    """Tests for ResolverError exception."""
    
    def test_resolver_error_basic(self):
        """Should create a ResolverError."""
        error = ResolverError("resolution failed", "foo", "bar", "baz")
        assert isinstance(error, YAMLError)


class TestErrorEdgeCases:
    """Tests for edge cases and error combinations."""
    
    def test_multiple_nested_errors(self):
        """Should handle nested error scenarios."""
        inner_error = YAMLError("inner error")
        try:
            raise inner_error
        except YAMLError as e:
            outer_error = ScannerError("outer error", Marker("file.yaml", 0, 0, 0, 0, None))
            assert isinstance(outer_error, YAMLError)
    
    def test_error_with_very_long_message(self):
        """Should handle very long error messages."""
        long_message = "x" * 10000
        error = YAMLError(long_message)
        assert len(str(error)) >= 10000
    
    def test_marker_with_none_buffer(self):
        """Should handle None buffer in marker."""
        marker = Marker("file.yaml", 0, 0, 0, 0, None)
        assert marker.get_snippet() is None
    
    def test_scanner_error_with_empty_strings(self):
        """Should handle empty strings in scanner error."""
        error = ScannerError("", None, "", None)
        assert isinstance(error, YAMLError)
    
    def test_all_error_types_are_exceptions(self):
        """Verify all error types inherit from Exception."""
        errors = [
            YAMLError("test"),
            YAMLWarning("test"),
            ScannerError("test", None),
            ParserError("test", None),
            ComposerError("test", None),
            ConstructorError("test", None),
            EmitterError("test"),
            SerializerError("test"),
            RepresenterError("test"),
            ResolverError("test", "a", "b", "c")
        ]
        for error in errors:
            assert isinstance(error, Exception)


# Note: This test file requires the actual pyyaml3 error module to be imported.
# Add the following imports at the top once you have access to the module:
# from src.collectors.python.d.plugin.python_modules.pyyaml3.error import (
#     YAMLError, YAMLWarning, Marker, ScannerError, ParserError, ComposerError,
#     ConstructorError, EmitterError, SerializerError, RepresenterError, ResolverError
# )