import unittest
import os
import sys
from unittest.mock import patch

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..', 'scripts')))
import deployment_pipeline

class TestDeploymentPipeline(unittest.TestCase):

    @patch('deployment_pipeline.os.path.exists')
    def test_run_pipeline_invalid_path(self, mock_exists):
        mock_exists.return_value = False
        success, msg = deployment_pipeline.run_pipeline("fake_path")
        self.assertFalse(success)
        self.assertTrue("does not exist" in msg)

    @patch('deployment_pipeline.os.path.exists')
    @patch('deployment_pipeline.build_stage')
    @patch('deployment_pipeline.test_stage')
    @patch('deployment_pipeline.deploy_stage')
    def test_run_pipeline_success(self, mock_deploy, mock_test, mock_build, mock_exists):
        mock_exists.return_value = True
        mock_build.return_value = True
        mock_test.return_value = True
        mock_deploy.return_value = True
<<<<<<< HEAD

=======

>>>>>>> jules-9396211896448288708-4318ead9
        success, msg = deployment_pipeline.run_pipeline("valid_path")
        self.assertTrue(success)
        self.assertTrue("successfully" in msg)

    @patch('deployment_pipeline.os.path.exists')
    @patch('deployment_pipeline.build_stage')
    def test_run_pipeline_build_fail(self, mock_build, mock_exists):
        mock_exists.return_value = True
        mock_build.return_value = False
<<<<<<< HEAD

=======

>>>>>>> jules-9396211896448288708-4318ead9
        success, msg = deployment_pipeline.run_pipeline("valid_path")
        self.assertFalse(success)
        self.assertEqual(msg, "Build stage failed.")

if __name__ == '__main__':
    unittest.main()
