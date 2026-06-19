import unittest
import os
import sys
from unittest.mock import patch

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..', 'scripts')))
import drift_detection_daemon

class TestDriftDetection(unittest.TestCase):

    @patch('drift_detection_daemon.os.path.exists')
    def test_uninitialized_submodule(self, mock_exists):
        mock_exists.return_value = False
        result = drift_detection_daemon.check_drift("fake/path")
        self.assertEqual(result, "fake/path: Not initialized")

    @patch('drift_detection_daemon.os.path.exists')
    @patch('drift_detection_daemon.run_cmd')
    def test_sync_submodule(self, mock_run_cmd, mock_exists):
        mock_exists.return_value = True

        # Mock responses for git fetch, local hash, branch name, and remote hash
        def side_effect(cmd, **kwargs):
            if "fetch" in cmd: return ""
            if "HEAD" in cmd and "rev-parse" in cmd: return "abcdef1"
            if "@{u}" in cmd: return "origin/main"
            if "origin/main" in cmd: return "abcdef1"
            return ""

        mock_run_cmd.side_effect = side_effect

        result = drift_detection_daemon.check_drift("fake/path")
        self.assertEqual(result, "fake/path: In sync")

    @patch('drift_detection_daemon.os.path.exists')
    @patch('drift_detection_daemon.run_cmd')
    def test_drifted_submodule(self, mock_run_cmd, mock_exists):
        mock_exists.return_value = True

        # Mock responses indicating a mismatch
        def side_effect(cmd, **kwargs):
            if "fetch" in cmd: return ""
            if "HEAD" in cmd and "rev-parse" in cmd: return "abcdef1"
            if "@{u}" in cmd: return "origin/main"
            if "origin/main" in cmd: return "1234567"
            return ""

        mock_run_cmd.side_effect = side_effect

        result = drift_detection_daemon.check_drift("fake/path")
        self.assertTrue("DRIFT DETECTED" in result)

if __name__ == '__main__':
    unittest.main()
