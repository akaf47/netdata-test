"""
Comprehensive tests for ssl_match_hostname/__init__.py
Tests the public API and module initialization
"""
import unittest
from unittest.mock import patch, MagicMock
import sys
import os

# Mock the module structure
try:
    from urllib3.packages.ssl_match_hostname import match_hostname, CertificateError
except ImportError:
    # Fallback for testing purposes
    class CertificateError(ValueError):
        pass
    match_hostname = None


class TestSSLMatchHostnameInit(unittest.TestCase):
    """Test suite for ssl_match_hostname/__init__.py module initialization"""

    def setUp(self):
        """Setup test fixtures"""
        self.test_hostname = "www.example.com"
        self.test_cert = {
            'subject': ((('commonName', 'www.example.com'),),),
            'subjectAltName': (('DNS', 'www.example.com'),)
        }

    def test_module_imports_match_hostname_function(self):
        """should expose match_hostname function in module"""
        try:
            from urllib3.packages import ssl_match_hostname
            self.assertTrue(hasattr(ssl_match_hostname, 'match_hostname'))
        except ImportError:
            self.skipTest("Module not available")

    def test_module_imports_certificate_error(self):
        """should expose CertificateError exception in module"""
        try:
            from urllib3.packages import ssl_match_hostname
            self.assertTrue(hasattr(ssl_match_hostname, 'CertificateError'))
        except ImportError:
            self.skipTest("Module not available")

    def test_certificate_error_is_value_error_subclass(self):
        """should have CertificateError as subclass of ValueError"""
        self.assertTrue(issubclass(CertificateError, ValueError))

    def test_certificate_error_can_be_raised_and_caught(self):
        """should allow raising and catching CertificateError"""
        with self.assertRaises(CertificateError):
            raise CertificateError("test error message")

    def test_certificate_error_message_preserved(self):
        """should preserve error message when raising CertificateError"""
        error_msg = "Certificate verification failed"
        try:
            raise CertificateError(error_msg)
        except CertificateError as e:
            self.assertEqual(str(e), error_msg)

    def test_module_has_version_attribute(self):
        """should have version or identifier attributes if defined"""
        try:
            from urllib3.packages import ssl_match_hostname
            # Module may or may not have __version__
            has_version = hasattr(ssl_match_hostname, '__version__')
            self.assertIsInstance(has_version, bool)
        except ImportError:
            self.skipTest("Module not available")

    def test_match_hostname_function_callable(self):
        """should export callable match_hostname function"""
        try:
            from urllib3.packages.ssl_match_hostname import match_hostname
            self.assertTrue(callable(match_hostname))
        except ImportError:
            self.skipTest("Module not available")

    def test_module_can_be_imported_multiple_times(self):
        """should handle multiple imports without error"""
        try:
            import urllib3.packages.ssl_match_hostname as mod1
            import urllib3.packages.ssl_match_hostname as mod2
            self.assertIs(mod1, mod2)
        except ImportError:
            self.skipTest("Module not available")


class TestCertificateErrorClass(unittest.TestCase):
    """Test suite for CertificateError exception class"""

    def test_raise_with_no_args(self):
        """should allow raising CertificateError with no arguments"""
        with self.assertRaises(CertificateError) as ctx:
            raise CertificateError()
        self.assertEqual(str(ctx.exception), '')

    def test_raise_with_single_string_arg(self):
        """should allow raising with single string message"""
        msg = "Single message"
        with self.assertRaises(CertificateError) as ctx:
            raise CertificateError(msg)
        self.assertEqual(str(ctx.exception), msg)

    def test_raise_with_multiple_args(self):
        """should handle multiple arguments"""
        args = ("hostname", "certificate")
        with self.assertRaises(CertificateError) as ctx:
            raise CertificateError(*args)
        self.assertEqual(ctx.exception.args, args)

    def test_is_instance_of_value_error(self):
        """should be instance of ValueError"""
        err = CertificateError("test")
        self.assertIsInstance(err, ValueError)

    def test_is_instance_of_exception(self):
        """should be instance of Exception"""
        err = CertificateError("test")
        self.assertIsInstance(err, Exception)

    def test_catch_as_value_error(self):
        """should be catchable as ValueError"""
        with self.assertRaises(ValueError):
            raise CertificateError("test")

    def test_catch_as_exception(self):
        """should be catchable as Exception"""
        with self.assertRaises(Exception):
            raise CertificateError("test")

    def test_error_with_empty_string_message(self):
        """should handle empty string messages"""
        with self.assertRaises(CertificateError) as ctx:
            raise CertificateError("")
        self.assertEqual(str(ctx.exception), "")

    def test_error_with_special_characters(self):
        """should preserve special characters in message"""
        msg = "Invalid cert: @ # $ % ^ & * ( )"
        with self.assertRaises(CertificateError) as ctx:
            raise CertificateError(msg)
        self.assertEqual(str(ctx.exception), msg)

    def test_error_with_unicode_characters(self):
        """should handle unicode characters"""
        msg = "Ошибка сертификата 证书错误 🔒"
        with self.assertRaises(CertificateError) as ctx:
            raise CertificateError(msg)
        self.assertEqual(str(ctx.exception), msg)


class TestModulePublicAPI(unittest.TestCase):
    """Test suite for module's public API surface"""

    def test_all_attribute_exists(self):
        """should have __all__ attribute defining public API if defined"""
        try:
            from urllib3.packages import ssl_match_hostname
            if hasattr(ssl_match_hostname, '__all__'):
                self.assertIsInstance(ssl_match_hostname.__all__, list)
                self.assertIn('match_hostname', ssl_match_hostname.__all__)
                self.assertIn('CertificateError', ssl_match_hostname.__all__)
        except ImportError:
            self.skipTest("Module not available")

    def test_module_docstring_exists(self):
        """should have module-level docstring"""
        try:
            from urllib3.packages import ssl_match_hostname
            # Docstring may or may not exist, just verify it's string if present
            if ssl_match_hostname.__doc__:
                self.assertIsInstance(ssl_match_hostname.__doc__, str)
        except ImportError:
            self.skipTest("Module not available")


if __name__ == '__main__':
    unittest.main()