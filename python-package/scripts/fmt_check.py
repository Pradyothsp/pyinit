import subprocess
import sys

# ANSI escape codes for colors
CYAN = "\033[96m"
GREEN = "\033[92m"
RESET = "\033[0m"


def main():
    try:
        print(f"{CYAN}Running:{RESET} {GREEN}ruff check .{RESET}")
        subprocess.run(["ruff", "check", "."], check=True)

        print(f"{CYAN}Running:{RESET} {GREEN}pyright{RESET}")
        subprocess.run(["pyright"], check=True)

        print(f"{CYAN}Running:{RESET} {GREEN}ruff format --check --diff .{RESET}")
        subprocess.run(["ruff", "format", "--check", "--diff", "."], check=True)

    except subprocess.CalledProcessError as e:
        sys.exit(e.returncode)


if __name__ == "__main__":
    main()
