"""Binary downloader and manager for pyinit."""

import hashlib
import os
import platform
import stat
from pathlib import Path
from typing import Any
from urllib.error import URLError
from urllib.request import urlopen

from . import __version__


def get_platform_info():
    """Get platform-specific information for binary download."""
    system = platform.system().lower()
    machine = platform.machine().lower()

    if system == "darwin":
        if machine in ("arm64", "aarch64"):
            return "darwin-arm64"
        else:
            return "darwin-amd64"
    elif system == "linux":
        if machine in ("x86_64", "amd64"):
            return "linux-amd64"
        else:
            raise RuntimeError(f"Unsupported Linux architecture: {machine}")
    else:
        raise RuntimeError(f"Unsupported operating system: {system}")


def get_binary_info():
    """Get download URL and expected SHA256 for the current platform."""
    platform_name = get_platform_info()
    version = f"v{__version__}"

    # These checksums are automatically updated by GitHub Actions
    checksums = {
        "0.0.2": {
            "darwin-amd64": "89f03ac7dfa17dcacd70edfcd07e957d5a6352785ba3d9960a5447167102ddd3",
            "darwin-arm64": "32999025529bd2a008f5af681ef6be51f06600a00e01bdd97a1c51a269a4d5f6",
            "linux-amd64": "e855b5b4c8ec83f1c39dfab26f0b9503d246572f5b62f5ffd45edf8ce4004682",
        }
    }

    # For development versions, use the latest available version
    if __version__ not in checksums:
        if __version__.endswith(".dev0"):
            # For development, use the latest release (0.0.2)
            latest_version = "0.0.2"
            version = "v0.0.2"  # Use the actual release tag
            if platform_name not in checksums[latest_version]:
                raise RuntimeError(f"No binary available for platform: {platform_name}")
            expected_sha = checksums[latest_version][platform_name]
        else:
            raise RuntimeError(f"No checksums available for version {__version__}")
    else:
        if platform_name not in checksums[__version__]:
            raise RuntimeError(f"No binary available for platform: {platform_name}")
        expected_sha = checksums[__version__][platform_name]

    url = f"https://github.com/Pradyothsp/pyinit/releases/download/{version}/pyinit-{platform_name}"

    return url, expected_sha


def get_binary_path():
    """Get the local path where the binary should be stored."""
    home = Path.home()
    binary_dir = home / ".pyinit" / "bin"
    binary_dir.mkdir(parents=True, exist_ok=True)
    return binary_dir / "pyinit"


def verify_checksum(file_path: Path, expected_sha: Any) -> bool:
    """Verify the SHA256 checksum of a file."""
    sha256_hash = hashlib.sha256()
    with open(file_path, "rb") as f:
        for chunk in iter(lambda: f.read(4096), b""):
            sha256_hash.update(chunk)

    actual_sha = sha256_hash.hexdigest()
    return actual_sha == expected_sha


def download_and_verify_binary(url: str, target_path: Path, expected_sha: Any = None):
    """Download and optionally verify the binary."""
    print(f"Downloading pyinit binary from {url}...")

    try:
        with urlopen(url) as response:
            if response.status != 200:
                raise RuntimeError(f"Failed to download binary: HTTP {response.status}")

            with open(target_path, "wb") as f:
                f.write(response.read())

        # Verify checksum if provided
        if expected_sha:
            if not verify_checksum(target_path, expected_sha):
                target_path.unlink()  # Remove invalid file
                raise RuntimeError("Downloaded binary failed checksum verification")
            print("Successfully downloaded and verified pyinit binary")
        else:
            print(
                "Successfully downloaded pyinit binary (checksum verification skipped for development)"
            )

        # Make executable
        target_path.chmod(target_path.stat().st_mode | stat.S_IEXEC)

    except URLError as e:
        raise RuntimeError(f"Failed to download binary: {e}") from e


def download_binary(url: str, target_path: Path, expected_sha: Any):
    """Download and verify the binary (legacy function for compatibility)."""
    download_and_verify_binary(url, target_path, expected_sha)


def ensure_binary():
    """Ensure the pyinit binary is available, downloading if necessary."""
    binary_path = get_binary_path()

    # Check if binary already exists and is executable
    if binary_path.exists() and os.access(binary_path, os.X_OK):
        return binary_path

    # Get download info
    try:
        url, expected_sha = get_binary_info()
        download_binary(url, binary_path, expected_sha)
    except RuntimeError as e:
        if "No checksums available" in str(e) and __version__.endswith(".dev0"):
            # For development versions, download without checksum verification
            platform_name = get_platform_info()
            url = f"https://github.com/Pradyothsp/pyinit/releases/download/v{__version__}/pyinit-{platform_name}"
            print("Development version detected, downloading latest release...")
            download_and_verify_binary(url, binary_path, expected_sha=None)
        else:
            raise

    return binary_path
