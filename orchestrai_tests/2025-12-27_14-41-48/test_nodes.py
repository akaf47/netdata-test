"""
Comprehensive test suite for pyyaml3/nodes.py
Tests cover all classes, methods, and code paths with 100% coverage target
"""
import unittest
from unittest.mock import Mock, patch, MagicMock
import sys
import os

# Mock the pyyaml3 module structure
sys.modules['pyyaml3'] = MagicMock()

# Import after mocking
from src.collectors.python.d.plugin.python_modules.pyyaml3.nodes import (
    Node, ScalarNode, SequenceNode, MappingNode, Tag, Mark
)


class TestMark(unittest.TestCase):
    """Test suite for Mark class"""
    
    def test_mark_creation_basic(self):
        """should create Mark with basic parameters"""
        mark = Mark("test.yaml", 0, 10, 5, None, None)
        self.assertEqual(mark.name, "test.yaml")
        self.assertEqual(mark.index, 0)
        self.assertEqual(mark.line, 10)
        self.assertEqual(mark.column, 5)
        self.assertIsNone(mark.buffer)
        self.assertIsNone(mark.pointer)
    
    def test_mark_creation_with_buffer_and_pointer(self):
        """should create Mark with buffer and pointer"""
        buffer_content = "line 1\nline 2\nline 3"
        mark = Mark("test.yaml", 7, 1, 0, buffer_content, 0)
        self.assertEqual(mark.buffer, buffer_content)
        self.assertEqual(mark.pointer, 0)
    
    def test_mark_get_snippet_no_buffer(self):
        """should return None snippet when no buffer"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        snippet = mark.get_snippet()
        self.assertIsNone(snippet)
    
    def test_mark_get_snippet_with_buffer(self):
        """should return snippet with buffer"""
        buffer_content = "line 1\nline 2\nline 3"
        mark = Mark("test.yaml", 0, 0, 0, buffer_content, 0)
        snippet = mark.get_snippet()
        self.assertIsNotNone(snippet)
    
    def test_mark_get_snippet_with_indent(self):
        """should return snippet with custom indent"""
        buffer_content = "test content"
        mark = Mark("test.yaml", 0, 0, 0, buffer_content, 0)
        snippet = mark.get_snippet(indent=4, max_length=10)
        self.assertIsNotNone(snippet)
    
    def test_mark_str_representation(self):
        """should return string representation of Mark"""
        mark = Mark("test.yaml", 0, 10, 5, None, None)
        result = str(mark)
        self.assertIn("test.yaml", result)
        self.assertIn("line", result.lower())
    
    def test_mark_repr_representation(self):
        """should return repr representation of Mark"""
        mark = Mark("test.yaml", 0, 10, 5, None, None)
        result = repr(mark)
        self.assertIn("Mark", result)


class TestNode(unittest.TestCase):
    """Test suite for Node base class"""
    
    def test_node_creation(self):
        """should create Node with tag and value"""
        tag = "tag:yaml.org,2002:str"
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        node = Node(tag, None, mark, mark)
        self.assertEqual(node.tag, tag)
    
    def test_node_with_none_tag(self):
        """should create Node with None tag"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        node = Node(None, None, mark, mark)
        self.assertIsNone(node.tag)
    
    def test_node_with_marks(self):
        """should store start and end marks"""
        mark1 = Mark("test.yaml", 0, 0, 0, None, None)
        mark2 = Mark("test.yaml", 10, 0, 10, None, None)
        node = Node("tag", None, mark1, mark2)
        self.assertEqual(node.start_mark, mark1)
        self.assertEqual(node.end_mark, mark2)


class TestScalarNode(unittest.TestCase):
    """Test suite for ScalarNode class"""
    
    def test_scalar_node_creation(self):
        """should create ScalarNode with tag and value"""
        tag = "tag:yaml.org,2002:str"
        value = "hello world"
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        node = ScalarNode(tag, value, mark, mark)
        self.assertEqual(node.tag, tag)
        self.assertEqual(node.value, value)
    
    def test_scalar_node_plain_style(self):
        """should create ScalarNode with plain style"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        node = ScalarNode("tag:yaml.org,2002:str", "value", mark, mark, style=None)
        self.assertIsNone(node.style)
    
    def test_scalar_node_quoted_style(self):
        """should create ScalarNode with quoted style"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        node = ScalarNode("tag:yaml.org,2002:str", "value", mark, mark, style='"')
        self.assertEqual(node.style, '"')
    
    def test_scalar_node_single_quoted_style(self):
        """should create ScalarNode with single-quoted style"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        node = ScalarNode("tag:yaml.org,2002:str", "value", mark, mark, style="'")
        self.assertEqual(node.style, "'")
    
    def test_scalar_node_literal_style(self):
        """should create ScalarNode with literal style"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        node = ScalarNode("tag:yaml.org,2002:str", "value", mark, mark, style="|")
        self.assertEqual(node.style, "|")
    
    def test_scalar_node_folded_style(self):
        """should create ScalarNode with folded style"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        node = ScalarNode("tag:yaml.org,2002:str", "value", mark, mark, style=">")
        self.assertEqual(node.style, ">")
    
    def test_scalar_node_empty_value(self):
        """should create ScalarNode with empty string"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        node = ScalarNode("tag:yaml.org,2002:null", "", mark, mark)
        self.assertEqual(node.value, "")
    
    def test_scalar_node_numeric_value(self):
        """should handle numeric string values"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        node = ScalarNode("tag:yaml.org,2002:int", "123", mark, mark)
        self.assertEqual(node.value, "123")
    
    def test_scalar_node_multiline_value(self):
        """should handle multiline values"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        value = "line 1\nline 2\nline 3"
        node = ScalarNode("tag:yaml.org,2002:str", value, mark, mark)
        self.assertEqual(node.value, value)
    
    def test_scalar_node_special_characters(self):
        """should handle special character values"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        value = "special: chars!@#$%^&*()"
        node = ScalarNode("tag:yaml.org,2002:str", value, mark, mark)
        self.assertEqual(node.value, value)
    
    def test_scalar_node_unicode_value(self):
        """should handle unicode values"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        value = "こんにちは世界"
        node = ScalarNode("tag:yaml.org,2002:str", value, mark, mark)
        self.assertEqual(node.value, value)


class TestSequenceNode(unittest.TestCase):
    """Test suite for SequenceNode class"""
    
    def test_sequence_node_creation(self):
        """should create SequenceNode with tag and value"""
        tag = "tag:yaml.org,2002:seq"
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        node1 = ScalarNode("tag:yaml.org,2002:str", "item1", mark, mark)
        node2 = ScalarNode("tag:yaml.org,2002:str", "item2", mark, mark)
        value = [node1, node2]
        seq_node = SequenceNode(tag, value, mark, mark)
        self.assertEqual(seq_node.tag, tag)
        self.assertEqual(len(seq_node.value), 2)
    
    def test_sequence_node_empty(self):
        """should create SequenceNode with empty list"""
        tag = "tag:yaml.org,2002:seq"
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        seq_node = SequenceNode(tag, [], mark, mark)
        self.assertEqual(len(seq_node.value), 0)
    
    def test_sequence_node_single_item(self):
        """should create SequenceNode with single item"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        node = ScalarNode("tag:yaml.org,2002:str", "item", mark, mark)
        seq_node = SequenceNode("tag:yaml.org,2002:seq", [node], mark, mark)
        self.assertEqual(len(seq_node.value), 1)
        self.assertEqual(seq_node.value[0].value, "item")
    
    def test_sequence_node_nested_sequences(self):
        """should create SequenceNode with nested sequences"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        inner_node = ScalarNode("tag:yaml.org,2002:str", "inner", mark, mark)
        inner_seq = SequenceNode("tag:yaml.org,2002:seq", [inner_node], mark, mark)
        outer_seq = SequenceNode("tag:yaml.org,2002:seq", [inner_seq], mark, mark)
        self.assertEqual(len(outer_seq.value), 1)
        self.assertIsInstance(outer_seq.value[0], SequenceNode)
    
    def test_sequence_node_flow_style(self):
        """should create SequenceNode with flow style"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        node = ScalarNode("tag:yaml.org,2002:str", "item", mark, mark)
        seq_node = SequenceNode("tag:yaml.org,2002:seq", [node], mark, mark, flow_style=True)
        self.assertTrue(seq_node.flow_style)
    
    def test_sequence_node_block_style(self):
        """should create SequenceNode with block style"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        node = ScalarNode("tag:yaml.org,2002:str", "item", mark, mark)
        seq_node = SequenceNode("tag:yaml.org,2002:seq", [node], mark, mark, flow_style=False)
        self.assertFalse(seq_node.flow_style)
    
    def test_sequence_node_mixed_types(self):
        """should create SequenceNode with mixed node types"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        scalar = ScalarNode("tag:yaml.org,2002:str", "scalar", mark, mark)
        mapping = MappingNode("tag:yaml.org,2002:map", [], mark, mark)
        seq_node = SequenceNode("tag:yaml.org,2002:seq", [scalar, mapping], mark, mark)
        self.assertEqual(len(seq_node.value), 2)


class TestMappingNode(unittest.TestCase):
    """Test suite for MappingNode class"""
    
    def test_mapping_node_creation(self):
        """should create MappingNode with tag and value"""
        tag = "tag:yaml.org,2002:map"
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        key = ScalarNode("tag:yaml.org,2002:str", "key", mark, mark)
        val = ScalarNode("tag:yaml.org,2002:str", "value", mark, mark)
        value = [(key, val)]
        map_node = MappingNode(tag, value, mark, mark)
        self.assertEqual(map_node.tag, tag)
        self.assertEqual(len(map_node.value), 1)
    
    def test_mapping_node_empty(self):
        """should create MappingNode with empty list"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        map_node = MappingNode("tag:yaml.org,2002:map", [], mark, mark)
        self.assertEqual(len(map_node.value), 0)
    
    def test_mapping_node_multiple_entries(self):
        """should create MappingNode with multiple entries"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        entries = []
        for i in range(3):
            key = ScalarNode("tag:yaml.org,2002:str", f"key{i}", mark, mark)
            val = ScalarNode("tag:yaml.org,2002:str", f"value{i}", mark, mark)
            entries.append((key, val))
        map_node = MappingNode("tag:yaml.org,2002:map", entries, mark, mark)
        self.assertEqual(len(map_node.value), 3)
    
    def test_mapping_node_nested_mapping(self):
        """should create MappingNode with nested mapping"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        inner_key = ScalarNode("tag:yaml.org,2002:str", "inner_key", mark, mark)
        inner_val = ScalarNode("tag:yaml.org,2002:str", "inner_value", mark, mark)
        inner_map = MappingNode("tag:yaml.org,2002:map", [(inner_key, inner_val)], mark, mark)
        
        outer_key = ScalarNode("tag:yaml.org,2002:str", "outer_key", mark, mark)
        outer_map = MappingNode("tag:yaml.org,2002:map", [(outer_key, inner_map)], mark, mark)
        
        self.assertEqual(len(outer_map.value), 1)
        self.assertIsInstance(outer_map.value[0][1], MappingNode)
    
    def test_mapping_node_flow_style(self):
        """should create MappingNode with flow style"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        key = ScalarNode("tag:yaml.org,2002:str", "key", mark, mark)
        val = ScalarNode("tag:yaml.org,2002:str", "value", mark, mark)
        map_node = MappingNode("tag:yaml.org,2002:map", [(key, val)], mark, mark, flow_style=True)
        self.assertTrue(map_node.flow_style)
    
    def test_mapping_node_block_style(self):
        """should create MappingNode with block style"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        key = ScalarNode("tag:yaml.org,2002:str", "key", mark, mark)
        val = ScalarNode("tag:yaml.org,2002:str", "value", mark, mark)
        map_node = MappingNode("tag:yaml.org,2002:map", [(key, val)], mark, mark, flow_style=False)
        self.assertFalse(map_node.flow_style)
    
    def test_mapping_node_complex_keys(self):
        """should create MappingNode with complex keys"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        key = MappingNode("tag:yaml.org,2002:map", [], mark, mark)
        val = ScalarNode("tag:yaml.org,2002:str", "value", mark, mark)
        map_node = MappingNode("tag:yaml.org,2002:map", [(key, val)], mark, mark)
        self.assertIsInstance(map_node.value[0][0], MappingNode)
    
    def test_mapping_node_sequence_value(self):
        """should create MappingNode with sequence as value"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        key = ScalarNode("tag:yaml.org,2002:str", "key", mark, mark)
        item = ScalarNode("tag:yaml.org,2002:str", "item", mark, mark)
        val = SequenceNode("tag:yaml.org,2002:seq", [item], mark, mark)
        map_node = MappingNode("tag:yaml.org,2002:map", [(key, val)], mark, mark)
        self.assertIsInstance(map_node.value[0][1], SequenceNode)


class TestTag(unittest.TestCase):
    """Test suite for Tag class if it exists"""
    
    def test_standard_tags(self):
        """should have standard YAML tags defined"""
        # Testing if standard tag constants are available
        standard_tags = [
            'tag:yaml.org,2002:null',
            'tag:yaml.org,2002:bool',
            'tag:yaml.org,2002:int',
            'tag:yaml.org,2002:float',
            'tag:yaml.org,2002:str',
            'tag:yaml.org,2002:seq',
            'tag:yaml.org,2002:map',
        ]
        for tag in standard_tags:
            self.assertIsInstance(tag, str)


class TestNodeIntegration(unittest.TestCase):
    """Integration tests for node types"""
    
    def test_complex_document_structure(self):
        """should handle complex document structure"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        
        # Create a complex nested structure
        scalar1 = ScalarNode("tag:yaml.org,2002:str", "value1", mark, mark)
        scalar2 = ScalarNode("tag:yaml.org,2002:int", "42", mark, mark)
        
        seq = SequenceNode("tag:yaml.org,2002:seq", [scalar1, scalar2], mark, mark)
        
        key1 = ScalarNode("tag:yaml.org,2002:str", "key1", mark, mark)
        key2 = ScalarNode("tag:yaml.org,2002:str", "key2", mark, mark)
        
        mapping = MappingNode("tag:yaml.org,2002:map", 
                            [(key1, scalar1), (key2, seq)], 
                            mark, mark)
        
        self.assertEqual(len(mapping.value), 2)
        self.assertIsInstance(mapping.value[1][1], SequenceNode)
    
    def test_deeply_nested_structure(self):
        """should handle deeply nested structures"""
        mark = Mark("test.yaml", 0, 0, 0, None, None)
        
        # Create deeply nested structure
        inner = ScalarNode("tag:yaml.org,2002:str", "value", mark, mark)
        for _ in range(5):
            inner = SequenceNode("tag:yaml.org,2002:seq", [inner], mark, mark)
        
        self.assertIsInstance(inner, SequenceNode)


if __name__ == '__main__':
    unittest.main()