"""Python wrapper for the pyinit Go binary."""

try:
    from importlib.metadata import version

    __version__ = version("pyinit-cli")
except ImportError:
    # Fallback for older Python versions or when not installed
    __version__ = "unknown"

__author__ = "Pradyoth S P"
__email__ = "contact@pradyoth-sp.me"
