import unittest
import os
import sys
import tempfile
import shutil

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..', 'scripts')))
import ui_scaffold_generator

class TestUIScaffoldGenerator(unittest.TestCase):

    def setUp(self):
        self.test_dir = tempfile.mkdtemp()

    def tearDown(self):
        shutil.rmtree(self.test_dir)

    def test_generate_component_name(self):
        self.assertEqual(ui_scaffold_generator.generate_component_name("user_auth.py"), "UserAuth")
        self.assertEqual(ui_scaffold_generator.generate_component_name("data-processor.rs"), "DataProcessor")
        self.assertEqual(ui_scaffold_generator.generate_component_name("main.go"), "Main")

    def test_create_scaffold_success(self):
        success, msg = ui_scaffold_generator.create_scaffold(self.test_dir, "test_backend.py")
        self.assertTrue(success)
        self.assertTrue("Successfully created" in msg)

        # Verify file exists and has content
        file_path = os.path.join(self.test_dir, "TestBackend.jsx")
        self.assertTrue(os.path.exists(file_path))
        with open(file_path, "r") as f:
            content = f.read()
            self.assertTrue("const TestBackend = () =>" in content)
            self.assertTrue("test_backend-container" in content)

    def test_create_scaffold_already_exists(self):
        ui_scaffold_generator.create_scaffold(self.test_dir, "test_backend.py")
        success, msg = ui_scaffold_generator.create_scaffold(self.test_dir, "test_backend.py")
        self.assertFalse(success)
        self.assertTrue("already exists" in msg)

if __name__ == '__main__':
    unittest.main()
