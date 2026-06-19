import unittest
import os
import sys
from unittest.mock import patch

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..', 'scripts')))
import ui_auditor

class TestUIAuditor(unittest.TestCase):

    @patch('ui_auditor.scan_directory_for_extensions')
    def test_audit_submodule_ok(self, mock_scan):
        # Simulate balanced frontend and backend
        mock_scan.side_effect = [5, 5]
        result = ui_auditor.audit_submodule("mock_repo")
        self.assertEqual(result["status"], "OK")

    @patch('ui_auditor.scan_directory_for_extensions')
    def test_audit_submodule_zero_frontend(self, mock_scan):
        # Simulate backend files but zero frontend
        mock_scan.side_effect = [5, 0]
        result = ui_auditor.audit_submodule("mock_repo")
        self.assertTrue("ZERO frontend files" in result["status"])

    @patch('ui_auditor.scan_directory_for_extensions')
    def test_audit_submodule_high_ratio(self, mock_scan):
        # Simulate many backend files and very few frontend
        mock_scan.side_effect = [15, 1]
        result = ui_auditor.audit_submodule("mock_repo")
        self.assertTrue("High backend-to-frontend ratio" in result["status"])

if __name__ == '__main__':
    unittest.main()
