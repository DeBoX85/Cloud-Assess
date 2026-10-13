#!/usr/bin/env python3
"""Exercise QA cleanup failure classification, bounds and a native file lock."""
import ctypes
import os
from pathlib import Path
import tempfile
import unittest
from unittest.mock import Mock, patch

from qa_temp_cleanup import cleanup_temporary_directory


def sharing_violation():
    error = PermissionError('controlled sharing violation')
    error.winerror = 32
    return error


class CleanupTests(unittest.TestCase):
    def test_immediate_success(self):
        directory = Mock()
        with patch('qa_temp_cleanup.time.sleep') as pause:
            cleanup_temporary_directory(directory)
        self.assertEqual(directory.cleanup.call_count, 1)
        pause.assert_not_called()

    def test_transient_sharing_violation(self):
        directory = Mock()
        directory.cleanup.side_effect = [sharing_violation(), sharing_violation(), None]
        with patch('qa_temp_cleanup.time.sleep') as pause:
            cleanup_temporary_directory(directory)
        self.assertEqual(directory.cleanup.call_count, 3)
        self.assertEqual([call.args[0] for call in pause.call_args_list], [0.1, 0.2])

    def test_persistent_sharing_violation_fails_with_original_error(self):
        error = sharing_violation()
        directory = Mock()
        directory.cleanup.side_effect = error
        with patch('qa_temp_cleanup.time.sleep') as pause:
            with self.assertRaises(PermissionError) as raised:
                cleanup_temporary_directory(directory)
        self.assertIs(raised.exception, error)
        self.assertEqual(directory.cleanup.call_count, 8)
        self.assertEqual(pause.call_count, 7)
        self.assertLessEqual(sum(call.args[0] for call in pause.call_args_list), 5)

    def test_other_failures_are_not_retried(self):
        for error in (PermissionError('ordinary denied access'), RuntimeError('QA failure')):
            with self.subTest(error=type(error).__name__):
                directory = Mock()
                directory.cleanup.side_effect = error
                with patch('qa_temp_cleanup.time.sleep') as pause:
                    with self.assertRaises(type(error)) as raised:
                        cleanup_temporary_directory(directory)
                self.assertIs(raised.exception, error)
                self.assertEqual(directory.cleanup.call_count, 1)
                pause.assert_not_called()
        error = PermissionError('Windows access denied')
        error.winerror = 5
        directory = Mock()
        directory.cleanup.side_effect = error
        with patch('qa_temp_cleanup.time.sleep') as pause:
            with self.assertRaises(PermissionError):
                cleanup_temporary_directory(directory)
        self.assertEqual(directory.cleanup.call_count, 1)
        pause.assert_not_called()

    @unittest.skipUnless(os.name == 'nt', 'native Windows sharing lock')
    def test_native_windows_lock_released_between_cleanup_attempts(self):
        from ctypes import wintypes
        kernel = ctypes.WinDLL('kernel32', use_last_error=True)
        kernel.CreateFileW.argtypes = [
            wintypes.LPCWSTR, wintypes.DWORD, wintypes.DWORD, ctypes.c_void_p,
            wintypes.DWORD, wintypes.DWORD, wintypes.HANDLE]
        kernel.CreateFileW.restype = wintypes.HANDLE
        kernel.CloseHandle.argtypes = [wintypes.HANDLE]
        kernel.CloseHandle.restype = wintypes.BOOL
        directory = tempfile.TemporaryDirectory(prefix='cloud-assess-cleanup-lock-')
        self.addCleanup(cleanup_temporary_directory, directory)
        locked = Path(directory.name) / 'owned locked file.txt'
        locked.write_text('controlled lock')
        handle = kernel.CreateFileW(str(locked), 0x80000000, 1, None, 3, 0x80, None)
        self.assertNotEqual(handle, ctypes.c_void_p(-1).value, ctypes.get_last_error())
        released = False

        def release_after_observed_lock(_delay):
            nonlocal released
            self.assertFalse(released)
            self.assertTrue(locked.exists())
            self.assertTrue(kernel.CloseHandle(handle), ctypes.get_last_error())
            released = True

        try:
            with patch('qa_temp_cleanup.time.sleep', side_effect=release_after_observed_lock) as pause:
                cleanup_temporary_directory(directory)
            self.assertTrue(released, 'first cleanup must encounter the native sharing violation')
            self.assertEqual(pause.call_count, 1)
            self.assertFalse(Path(directory.name).exists())
        finally:
            if not released:
                self.assertTrue(kernel.CloseHandle(handle), ctypes.get_last_error())


if __name__ == '__main__':
    unittest.main()
