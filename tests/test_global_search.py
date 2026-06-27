
import logging
# Omni-Workspace Standard Telemetry
logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)
if not logger.handlers:
    handler = logging.StreamHandler()
    formatter = logging.Formatter('%(asctime)s - [OMNI] - %(name)s - %(levelname)s - %(message)s')
    handler.setFormatter(formatter)
    logger.addHandler(handler)

import os
import tempfile
import unittest
import shutil
from scripts.global_search_and_replace import process_file, search_and_replace_in_dir

class TestGlobalSearchAndReplace(unittest.TestCase):
    def setUp(self):
        self.test_dir = tempfile.mkdtemp()

    def tearDown(self):
        shutil.rmtree(self.test_dir)

    def test_process_file(self):
        test_file = os.path.join(self.test_dir, "test.txt")
        with open(test_file, 'w') as f:
            f.write("This is a old_string test.")

        result = process_file(test_file, "old_string", "new_string")
        self.assertTrue(result)

        with open(test_file, 'r') as f:
            content = f.read()
            self.assertEqual(content, "This is a new_string test.")

    def test_search_and_replace_in_dir(self):
        # Create a nested directory structure
        os.makedirs(os.path.join(self.test_dir, "subdir"))
        file1 = os.path.join(self.test_dir, "file1.txt")
        file2 = os.path.join(self.test_dir, "subdir", "file2.txt")

        with open(file1, 'w') as f:
            f.write("Find me here.")
        with open(file2, 'w') as f:
            f.write("Find me here too.")

        modifications = search_and_replace_in_dir(self.test_dir, "Find me", "Found you")

        self.assertEqual(modifications, 2)

        with open(file1, 'r') as f:
            self.assertEqual(f.read(), "Found you here.")
        with open(file2, 'r') as f:
            self.assertEqual(f.read(), "Found you here too.")

if __name__ == '__main__':
    unittest.main()
