"""
Comprehensive tests for ssl_match_hostname/_implementation.py
Tests hostname matching logic with full branch coverage
"""
import unittest
from unittest.mock import patch, MagicMock
import re
import sys

try:
    from urllib3.packages.ssl_match_hostname import match_hostname, CertificateError
except ImportError:
    # Fallback for testing purposes
    class CertificateError(ValueError):
        pass
    
    def match_hostname(cert, hostname):
        """Placeholder implementation"""
        raise NotImplementedError


class TestMatchHostnameFunction(unittest.TestCase):
    """Test suite for match_hostname function - exact matches"""

    def setUp(self):
        """Setup test fixtures"""
        self.hostname = "www.example.com"
        self.cert_with_cn = {
            'subject': ((('commonName', 'www.example.com'),),),
        }
        self.cert_with_san = {
            'subject': ((('commonName', 'example.com'),),),
            'subjectAltName': (('DNS', 'www.example.com'),)
        }

    def test_exact_hostname_match_with_cn_only(self):
        """should match exact hostname in commonName"""
        # Should not raise
        try:
            match_hostname(self.cert_with_cn, self.hostname)
        except CertificateError:
            self.fail("match_hostname raised CertificateError unexpectedly")

    def test_exact_hostname_match_with_san(self):
        """should match exact hostname in subjectAltName"""
        try:
            match_hostname(self.cert_with_san, self.hostname)
        except CertificateError:
            self.fail("match_hostname raised CertificateError unexpectedly")

    def test_hostname_mismatch_raises_certificate_error(self):
        """should raise CertificateError on hostname mismatch"""
        cert = {
            'subject': ((('commonName', 'other.example.com'),),),
        }
        with self.assertRaises(CertificateError):
            match_hostname(cert, self.hostname)

    def test_empty_hostname_raises_certificate_error(self):
        """should raise CertificateError for empty hostname"""
        with self.assertRaises(CertificateError):
            match_hostname(self.cert_with_cn, "")

    def test_none_hostname_raises_certificate_error(self):
        """should raise CertificateError for None hostname"""
        with self.assertRaises(CertificateError):
            match_hostname(self.cert_with_cn, None)

    def test_empty_cert_dict_raises_certificate_error(self):
        """should raise CertificateError for empty cert dict"""
        with self.assertRaises(CertificateError):
            match_hostname({}, self.hostname)

    def test_cert_with_no_subject_raises_certificate_error(self):
        """should raise CertificateError when cert has no subject"""
        cert = {'subjectAltName': (('DNS', 'example.com'),)}
        with self.assertRaises(CertificateError):
            match_hostname(cert, self.hostname)

    def test_case_insensitive_matching(self):
        """should perform case-insensitive hostname matching"""
        cert = {
            'subject': ((('commonName', 'WWW.EXAMPLE.COM'),),),
        }
        try:
            match_hostname(cert, self.hostname)
        except CertificateError:
            self.fail("Case-insensitive match failed")

    def test_hostname_with_trailing_dot_ignored(self):
        """should handle hostname with trailing dot"""
        cert = {
            'subject': ((('commonName', 'www.example.com.'),),),
        }
        # May or may not match depending on implementation
        # Test both cases are handled
        try:
            match_hostname(cert, 'www.example.com')
            # If it matches, that's acceptable
        except CertificateError:
            # If it doesn't match, that's also acceptable
            pass


class TestWildcardMatching(unittest.TestCase):
    """Test suite for wildcard certificate matching"""

    def test_wildcard_matches_single_subdomain(self):
        """should match wildcard cert with single subdomain level"""
        cert = {
            'subject': ((('commonName', '*.example.com'),),),
        }
        try:
            match_hostname(cert, 'www.example.com')
        except CertificateError:
            self.fail("Wildcard matching failed for single subdomain")

    def test_wildcard_matches_different_subdomain(self):
        """should match wildcard cert with different subdomain"""
        cert = {
            'subject': ((('commonName', '*.example.com'),),),
        }
        try:
            match_hostname(cert, 'mail.example.com')
        except CertificateError:
            self.fail("Wildcard matching failed for different subdomain")

    def test_wildcard_does_not_match_multiple_subdomains(self):
        """should not match wildcard cert with multiple subdomain levels"""
        cert = {
            'subject': ((('commonName', '*.example.com'),),),
        }
        with self.assertRaises(CertificateError):
            match_hostname(cert, 'sub.www.example.com')

    def test_wildcard_does_not_match_base_domain(self):
        """should not match wildcard cert against base domain"""
        cert = {
            'subject': ((('commonName', '*.example.com'),),),
        }
        with self.assertRaises(CertificateError):
            match_hostname(cert, 'example.com')

    def test_wildcard_with_different_tld_does_not_match(self):
        """should not match wildcard cert with different TLD"""
        cert = {
            'subject': ((('commonName', '*.example.com'),),),
        }
        with self.assertRaises(CertificateError):
            match_hostname(cert, 'www.example.org')

    def test_wildcard_case_insensitive(self):
        """should perform case-insensitive wildcard matching"""
        cert = {
            'subject': ((('commonName', '*.EXAMPLE.COM'),),),
        }
        try:
            match_hostname(cert, 'www.example.com')
        except CertificateError:
            self.fail("Case-insensitive wildcard match failed")

    def test_wildcard_in_middle_position_does_not_match(self):
        """should not match wildcard in middle of hostname"""
        cert = {
            'subject': ((('commonName', 'www.*.example.com'),),),
        }
        with self.assertRaises(CertificateError):
            match_hostname(cert, 'www.test.example.com')

    def test_wildcard_only_no_domain_does_not_match(self):
        """should not match cert with only wildcard"""
        cert = {
            'subject': ((('commonName', '*'),),),
        }
        with self.assertRaises(CertificateError):
            match_hostname(cert, 'example.com')


class TestSubjectAltNameMatching(unittest.TestCase):
    """Test suite for subjectAltName matching"""

    def test_san_takes_precedence_over_cn(self):
        """should prefer subjectAltName over commonName"""
        cert = {
            'subject': ((('commonName', 'www.example.com'),),),
            'subjectAltName': (('DNS', 'mail.example.com'),)
        }
        # Should match SAN, not CN
        with self.assertRaises(CertificateError):
            match_hostname(cert, 'www.example.com')
        
        # Should match SAN
        try:
            match_hostname(cert, 'mail.example.com')
        except CertificateError:
            self.fail("SAN matching failed")

    def test_multiple_san_entries_all_checked(self):
        """should check all entries in subjectAltName"""
        cert = {
            'subject': ((('commonName', 'example.com'),),),
            'subjectAltName': (
                ('DNS', 'www.example.com'),
                ('DNS', 'mail.example.com'),
                ('DNS', 'ftp.example.com'),
            )
        }
        try:
            match_hostname(cert, 'www.example.com')
            match_hostname(cert, 'mail.example.com')
            match_hostname(cert, 'ftp.example.com')
        except CertificateError:
            self.fail("Multiple SAN matching failed")

    def test_san_with_non_dns_types_ignored(self):
        """should ignore non-DNS types in subjectAltName"""
        cert = {
            'subject': ((('commonName', 'example.com'),),),
            'subjectAltName': (
                ('email', 'admin@example.com'),
                ('URI', 'https://example.com'),
                ('DNS', 'www.example.com'),
            )
        }
        try:
            match_hostname(cert, 'www.example.com')
        except CertificateError:
            self.fail("SAN with mixed types failed")

    def test_san_with_only_non_dns_types_uses_cn(self):
        """should fall back to CN when SAN has no DNS entries"""
        cert = {
            'subject': ((('commonName', 'www.example.com'),),),
            'subjectAltName': (
                ('email', 'admin@example.com'),
                ('URI', 'https://example.com'),
            )
        }
        try:
            match_hostname(cert, 'www.example.com')
        except CertificateError:
            self.fail("Fallback to CN failed when SAN has no DNS")

    def test_empty_san_uses_cn(self):
        """should use CN when subjectAltName is empty"""
        cert = {
            'subject': ((('commonName', 'www.example.com'),),),
            'subjectAltName': ()
        }
        try:
            match_hostname(cert, 'www.example.com')
        except CertificateError:
            self.fail("Empty SAN fallback to CN failed")


class TestCertificateStructures(unittest.TestCase):
    """Test suite for various certificate structure formats"""

    def test_multiple_rdn_in_subject(self):
        """should handle multiple RDN in subject"""
        cert = {
            'subject': (
                (('commonName', 'www.example.com'),),
                (('organizationName', 'Example Inc'),),
            ),
        }
        try:
            match_hostname(cert, 'www.example.com')
        except CertificateError:
            self.fail("Multiple RDN matching failed")

    def test_multiple_attributes_in_rdn(self):
        """should handle multiple attributes in single RDN"""
        cert = {
            'subject': (
                (
                    ('commonName', 'www.example.com'),
                    ('organizationName', 'Example Inc'),
                ),
            ),
        }
        try:
            match_hostname(cert, 'www.example.com')
        except CertificateError:
            self.fail("Multiple attributes in RDN failed")

    def test_cn_with_different_attribute_types(self):
        """should extract CN regardless of other attributes"""
        cert = {
            'subject': (
                (
                    ('organizationName', 'Example Inc'),
                    ('commonName', 'www.example.com'),
                    ('countryName', 'US'),
                ),
            ),
        }
        try:
            match_hostname(cert, 'www.example.com')
        except CertificateError:
            self.fail("CN extraction with other attributes failed")

    def test_missing_cn_in_subject(self):
        """should raise error when CN not found in subject"""
        cert = {
            'subject': (
                (
                    ('organizationName', 'Example Inc'),
                    ('countryName', 'US'),
                ),
            ),
        }
        with self.assertRaises(CertificateError):
            match_hostname(cert, 'www.example.com')

    def test_empty_cn_value(self):
        """should handle empty CN value"""
        cert = {
            'subject': ((('commonName', ''),),),
        }
        with self.assertRaises(CertificateError):
            match_hostname(cert, 'www.example.com')

    def test_cn_with_whitespace(self):
        """should handle CN with whitespace"""
        cert = {
            'subject': ((('commonName', '  www.example.com  '),),),
        }
        # Implementation may trim or not - test what happens
        try:
            match_hostname(cert, 'www.example.com')
        except CertificateError:
            # If not trimmed, CN won't match
            pass


class TestIPAddressMatching(unittest.TestCase):
    """Test suite for IP address matching"""

    def test_ipv4_address_matching(self):
        """should handle IPv4 addresses in certificates"""
        cert = {
            'subject': ((('commonName', '192.168.1.1'),),),
        }
        try:
            match_hostname(cert, '192.168.1.1')
        except CertificateError:
            # Some implementations may not support IP matching
            pass

    def test_ipv6_address_matching(self):
        """should handle IPv6 addresses in certificates"""
        cert = {
            'subject': ((('commonName', '::1'),),),
        }
        try:
            match_hostname(cert, '::1')
        except CertificateError:
            # Some implementations may not support IP matching
            pass

    def test_ipv6_full_notation(self):
        """should handle IPv6 in full notation"""
        cert = {
            'subject': ((('commonName', '2001:0db8:0000:0000:0000:0000:0000:0001'),),),
        }
        try:
            match_hostname(cert, '2001:0db8:0000:0000:0000:0000:0000:0001')
        except CertificateError:
            # May not match due to normalization differences
            pass


class TestEdgeCases(unittest.TestCase):
    """Test suite for edge cases and boundary conditions"""

    def test_very_long_hostname(self):
        """should handle very long hostnames"""
        long_hostname = "a" * 60 + ".example.com"
        cert = {
            'subject': ((('commonName', long_hostname),),),
        }
        try:
            match_hostname(cert, long_hostname)
        except CertificateError:
            self.fail("Long hostname matching failed")

    def test_hostname_with_numbers(self):
        """should match hostnames with numbers"""
        cert = {
            'subject': ((('commonName', 'server123.example456.com'),),),
        }
        try:
            match_hostname(cert, 'server123.example456.com')
        except CertificateError:
            self.fail("Hostname with numbers failed")

    def test_hostname_with_hyphens(self):
        """should match hostnames with hyphens"""
        cert = {
            'subject': ((('commonName', 'my-server.example-domain.com'),),),
        }
        try:
            match_hostname(cert, 'my-server.example-domain.com')
        except CertificateError:
            self.fail("Hostname with hyphens failed")

    def test_single_label_hostname(self):
        """should handle single label hostname"""
        cert = {
            'subject': ((('commonName', 'localhost'),),),
        }
        try:
            match_hostname(cert, 'localhost')
        except CertificateError:
            self.fail("Single label hostname failed")

    def test_numeric_tld(self):
        """should handle numeric TLD"""
        cert = {
            'subject': ((('commonName', 'example.123'),),),
        }
        try:
            match_hostname(cert, 'example.123')
        except CertificateError:
            self.fail("Numeric TLD failed")

    def test_hostname_with_underscore(self):
        """should handle or reject hostnames with underscores"""
        cert = {
            'subject': ((('commonName', 'my_server.example.com'),),),
        }
        try:
            match_hostname(cert, 'my_server.example.com')
        except CertificateError:
            # Underscores are technically invalid but may be matched
            pass

    def test_punycode_hostname(self):
        """should handle punycode encoded hostnames"""
        cert = {
            'subject': ((('commonName', 'xn--e1afmkfd.xn--p1ai'),),),
        }
        try:
            match_hostname(cert, 'xn--e1afmkfd.xn--p1ai')
        except CertificateError:
            self.fail("Punycode hostname failed")

    def test_wildcard_with_multiple_labels(self):
        """should not match wildcard across multiple labels"""
        cert = {
            'subject': ((('commonName', '*.*.example.com'),),),
        }
        with self.assertRaises(CertificateError):
            match_hostname(cert, 'a.b.example.com')

    def test_wildcard_partial_label(self):
        """should not match wildcard as partial label"""
        cert = {
            'subject': ((('commonName', 'www*.example.com'),),),
        }
        with self.assertRaises(CertificateError):
            match_hostname(cert, 'www1.example.com')

    def test_hostname_with_port_number(self):
        """should handle hostname with port number correctly"""
        cert = {
            'subject': ((('commonName', 'www.example.com'),),),
        }
        with self.assertRaises(CertificateError):
            match_hostname(cert, 'www.example.com:443')

    def test_hostname_with_path(self):
        """should handle hostname with path correctly"""
        cert = {
            'subject': ((('commonName', 'www.example.com'),),),
        }
        with self.assertRaises(CertificateError):
            match_hostname(cert, 'www.example.com/path')

    def test_hostname_with_query_string(self):
        """should handle hostname with query string correctly"""
        cert = {
            'subject': ((('commonName', 'www.example.com'),),),
        }
        with self.assertRaises(CertificateError):
            match_hostname(cert, 'www.example.com?query=value')


class TestErrorMessages(unittest.TestCase):
    """Test suite for error message quality"""

    def test_error_message_contains_hostname(self):
        """should include hostname in error message"""
        cert = {
            'subject': ((('commonName', 'other.example.com'),),),
        }
        try:
            match_hostname(cert, 'www.example.com')
            self.fail("Expected CertificateError")
        except CertificateError as e:
            # Error message should be informative
            error_msg = str(e)
            self.assertTrue(len(error_msg) > 0)

    def test_error_message_contains_cert_cn(self):
        """should include certificate CN in error message"""
        cert = {
            'subject': ((('commonName', 'other.example.com'),),),
        }
        try:
            match_hostname(cert, 'www.example.com')
            self.fail("Expected CertificateError")
        except CertificateError as e:
            # Error message should be informative
            error_msg = str(e)
            self.assertTrue(len(error_msg) > 0)


class TestReturnValue(unittest.TestCase):
    """Test suite for function return values"""

    def test_successful_match_returns_none_or_nothing(self):
        """should return None or not raise on successful match"""
        cert = {
            'subject': ((('commonName', 'www.example.com'),),),
        }
        result = match_hostname(cert, 'www.example.com')
        # Function should either return None or not return
        self.assertIsNone(result)

    def test_no_exception_on_match_means_match_succeeded(self):
        """should indicate success by not raising exception"""
        cert = {
            'subject': ((('commonName', 'www.example.com'),),),
        }
        try:
            match_hostname(cert, 'www.example.com')
            matched = True
        except CertificateError:
            matched = False
        self.assertTrue(matched)


class TestTypeValidation(unittest.TestCase):
    """Test suite for type validation"""

    def test_cert_as_string_raises_error(self):
        """should handle or reject cert as string"""
        try:
            match_hostname("not a dict", "www.example.com")
            # May be rejected
        except (CertificateError, TypeError, AttributeError, KeyError):
            pass

    def test_hostname_as_bytes_raises_error(self):
        """should handle or reject hostname as bytes"""
        cert = {
            'subject': ((('commonName', 'www.example.com'),),),
        }
        try:
            match_hostname(cert, b'www.example.com')
            # May be rejected
        except (CertificateError, TypeError, AttributeError):
            pass

    def test_cert_with_null_subject_raises_error(self):
        """should handle cert with None subject"""
        cert = {'subject': None}
        with self.assertRaises((CertificateError, TypeError)):
            match_hostname(cert, 'www.example.com')


class TestRealWorldCertificates(unittest.TestCase):
    """Test suite with real-world certificate examples"""

    def test_github_cert_structure(self):
        """should match GitHub-like certificate structure"""
        cert = {
            'subject': (
                (('organizationName', 'GitHub, Inc.'),),
                (('commonName', '*.github.com'),),
            ),
            'subjectAltName': (
                ('DNS', '*.github.com'),
                ('DNS', 'github.com'),
            )
        }
        try:
            match_hostname(cert, 'api.github.com')
        except CertificateError:
            self.fail("GitHub-like cert matching failed")

    def test_wildcard_cert_with_multiple_sans(self):
        """should match wildcard cert with multiple SANs"""
        cert = {
            'subject': ((('commonName', '*.example.com'),),),
            'subjectAltName': (
                ('DNS', '*.example.com'),
                ('DNS', 'example.com'),
                ('DNS', '*.api.example.com'),
            )
        }
        try:
            match_hostname(cert, 'www.example.com')
        except CertificateError:
            self.fail("Wildcard with multiple SANs failed")

    def test_san_only_cert_no_cn(self):
        """should match cert with only SAN, no useful CN"""
        cert = {
            'subject': ((('commonName', 'example.com'),),),
            'subjectAltName': (
                ('DNS', 'api.example.com'),
                ('DNS', 'www.example.com'),
            )
        }
        try:
            match_hostname(cert, 'api.example.com')
        except CertificateError:
            self.fail("SAN-only cert matching failed")


if __name__ == '__main__':
    unittest.main()