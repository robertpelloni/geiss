import unittest
import os
import sys
from unittest.mock import patch

# Adjust path to import from the scripts directory
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..', 'scripts')))

import update_repos_v5
import generate_dashboard
import prune_broken_submodules

class TestOrchestrationScripts(unittest.TestCase):

    @patch('update_repos_v5.run_cmd')
    def test_update_repos_get_submodules_empty(self, mock_run_cmd):
        # Test handling of an empty or missing .gitmodules file
        mock_run_cmd.return_value = ""
        result = update_repos_v5.get_submodules()
        self.assertEqual(result, [])

    @patch('update_repos_v5.run_cmd')
    def test_update_repos_get_submodules_exists(self, mock_run_cmd):
        # Test parsing a valid .gitmodules output
        mock_output = "submodule.path/one.path path/one\nsubmodule.path/two.path path/two"
        mock_run_cmd.return_value = mock_output
        result = update_repos_v5.get_submodules()
        self.assertEqual(result, ["path/one", "path/two"])

    @patch('generate_dashboard.run_cmd')
    def test_generate_dashboard_get_status_uninitialized(self, mock_run_cmd):
        # We need to mock os.path.exists to simulate an uninitialized submodule
        with patch('os.path.exists', return_value=False):
            branch, commit, status = generate_dashboard.get_repo_status("fake/path")
            self.assertEqual(branch, "Not Initialized")
            self.assertEqual(commit, "N/A")
            self.assertEqual(status, "N/A")

    @patch('prune_broken_submodules.run_cmd')
    def test_prune_get_submodules_handles_none(self, mock_run_cmd):
        # Ensure it safely handles a failed command execution (returns None)
        mock_run_cmd.return_value = None
        result = prune_broken_submodules.get_submodules()
        self.assertEqual(result, [])

if __name__ == '__main__':
    unittest.main()
