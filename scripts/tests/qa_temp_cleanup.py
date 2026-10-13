"""Bounded cleanup for QA-owned TemporaryDirectory instances."""
import time


def cleanup_temporary_directory(directory):
    """Retry only Windows sharing violations; every other failure propagates."""
    delays = (0.1, 0.2, 0.4, 0.8, 1.0, 1.0, 1.0)
    for attempt in range(len(delays) + 1):
        try:
            directory.cleanup()
            return
        except OSError as error:
            if getattr(error, 'winerror', None) != 32 or attempt == len(delays):
                raise
            time.sleep(delays[attempt])
