import unittest
import os
import sys
import tempfile
import shutil

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..', 'scripts')))
import telemetry_standardizer

class TestTelemetryStandardizer(unittest.TestCase):

    def setUp(self):
        self.test_dir = tempfile.mkdtemp()
        self.py_file = os.path.join(self.test_dir, "test.py")
        self.js_file = os.path.join(self.test_dir, "test.js")

        with open(self.py_file, 'w') as f:
            f.write("def my_func(): pass\n")

        with open(self.js_file, 'w') as f:
            f.write("function myFunc() {}\n")

    def tearDown(self):
        shutil.rmtree(self.test_dir)

    def test_detect_language(self):
        self.assertEqual(telemetry_standardizer.detect_language("script.py"), "python")
        self.assertEqual(telemetry_standardizer.detect_language("app.js"), "nodejs")
        self.assertEqual(telemetry_standardizer.detect_language("index.ts"), "nodejs")
        self.assertEqual(telemetry_standardizer.detect_language("README.md"), "unknown")

    def test_standardize_python_file(self):
        success, msg = telemetry_standardizer.standardize_file(self.py_file)
        self.assertTrue(success)

        with open(self.py_file, 'r') as f:
            content = f.read()
            self.assertTrue("[OMNI]" in content)
            self.assertTrue("import logging" in content)
            self.assertTrue("def my_func()" in content)

    def test_standardize_nodejs_file(self):
        success, msg = telemetry_standardizer.standardize_file(self.js_file)
        self.assertTrue(success)

        with open(self.js_file, 'r') as f:
            content = f.read()
            self.assertTrue("[OMNI]" in content)
            self.assertTrue("const omniLogger" in content)
            self.assertTrue("function myFunc()" in content)

    def test_skip_already_standardized_file(self):
        # Inject once
        telemetry_standardizer.standardize_file(self.py_file)
        # Attempt second injection
        success, msg = telemetry_standardizer.standardize_file(self.py_file)

        self.assertFalse(success)
        self.assertEqual(msg, "Telemetry already present")

if __name__ == '__main__':
    unittest.main()
