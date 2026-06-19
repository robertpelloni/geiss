import unittest
import os
import sys
from unittest.mock import patch

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..', 'scripts')))
import conflict_resolution

class TestConflictResolution(unittest.TestCase):

    def test_extract_conflict_blocks_no_file(self):
        result = conflict_resolution.extract_conflict_blocks("nonexistent.txt")
        self.assertEqual(result, [])

    def test_analyze_conflict_accept_remote(self):
        block = {"local": ["\n"], "remote": ["new feature\n"]}
        strategy = conflict_resolution.analyze_conflict(block)
        self.assertEqual(strategy, "Accept Remote (Local is empty)")

    def test_analyze_conflict_accept_local(self):
        block = {"local": ["my local fix\n"], "remote": ["\n", " "]}
        strategy = conflict_resolution.analyze_conflict(block)
        self.assertEqual(strategy, "Accept Local (Remote is empty)")

    def test_analyze_conflict_intelligent_merge(self):
        block = {"local": ["print('hello')\n"], "remote": ["print('world')\n"]}
        strategy = conflict_resolution.analyze_conflict(block)
        self.assertEqual(strategy, "Requires Intelligent Merge (Both contain features)")

if __name__ == '__main__':
    unittest.main()
