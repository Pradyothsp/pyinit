#!/usr/bin/env python3
"""Main CLI entry point for pyinit."""

import subprocess
import sys

from .downloader import ensure_binary


def main():
    """Main entry point for the pyinit CLI."""
    try:
        # Ensure the binary is downloaded and available
        binary_path = ensure_binary()

        # Execute the binary with all arguments passed through
        result = subprocess.run([str(binary_path), *sys.argv[1:]])
        sys.exit(result.returncode)

    except KeyboardInterrupt:
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
