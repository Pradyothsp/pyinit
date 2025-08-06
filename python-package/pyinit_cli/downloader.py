"""Binary downloader and manager for pyinit."""

import hashlib
import os
import platform
import stat
from pathlib import Path
from urllib.error import URLError
from urllib.request import urlopen

from . import __version__


def get_platform_info() -> str:
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
    elif system == "windows":
        if machine in ("x86_64", "amd64"):
            return "windows-amd64"
        else:
            raise RuntimeError(f"Unsupported Windows architecture: {machine}")
    else:
        raise RuntimeError(f"Unsupported operating system: {system}")


def get_binary_info() -> tuple[str, str | None]:
    """Get download URL and expected SHA256 for the current platform."""
    platform_name = get_platform_info()
    version = f"v{__version__}"

    # These checksums are automatically updated by GitHub Actions during build
    checksums: dict[str, dict[str, str]] = {
        # Checksums will be added here automatically during the release process
    }

    # Determine binary filename (Windows needs .exe extension)
    binary_filename = f"pyinit-{platform_name}"
    if platform_name.startswith("windows"):
        binary_filename += ".exe"

    if __version__ not in checksums:
        # For any version not in checksums (including development),
        # return None for checksum to skip verification
        url = f"https://github.com/Pradyothsp/pyinit/releases/download/v0.0.3/{binary_filename}"
        return url, None

    if platform_name not in checksums[__version__]:
        raise RuntimeError(f"No binary available for platform: {platform_name}")

    url = f"https://github.com/Pradyothsp/pyinit/releases/download/{version}/{binary_filename}"
    expected_sha: str = checksums[__version__][platform_name]

    return url, expected_sha


def get_binary_path() -> Path:
    """Get the local path where the binary should be stored."""
    home = Path.home()
    binary_dir = home / ".pyinit" / "bin"
    binary_dir.mkdir(parents=True, exist_ok=True)

    # Windows binaries need .exe extension
    platform_name = get_platform_info()
    binary_name = "pyinit"
    if platform_name.startswith("windows"):
        binary_name += ".exe"

    return binary_dir / binary_name


def verify_checksum(file_path: Path, expected_sha: str) -> bool:
    """Verify the SHA256 checksum of a file."""
    sha256_hash = hashlib.sha256()
    with open(file_path, "rb") as f:
        for chunk in iter(lambda: f.read(4096), b""):
            sha256_hash.update(chunk)

    actual_sha = sha256_hash.hexdigest()
    return actual_sha == expected_sha


def download_and_verify_binary(
    url: str, target_path: Path, expected_sha: str | None = None
) -> None:
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


def download_binary(url: str, target_path: Path, expected_sha: str) -> None:
    """Download and verify the binary (legacy function for compatibility)."""
    download_and_verify_binary(url, target_path, expected_sha)


def ensure_binary() -> Path:
    """Ensure the pyinit binary is available, downloading if necessary."""
    binary_path = get_binary_path()

    # Check if binary already exists and is executable
    if binary_path.exists() and os.access(binary_path, os.X_OK):
        return binary_path

    # Get download info
    url, expected_sha = get_binary_info()

    if expected_sha:
        # Use checksum verification for releases
        download_binary(url, binary_path, expected_sha)
    else:
        # Skip checksum verification for development/unknown versions
        print(
            f"No checksums available for version {__version__}, downloading without verification..."
        )
        download_and_verify_binary(url, binary_path, expected_sha=None)

    return binary_path
