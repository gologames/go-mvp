import os
import platform
import re
import subprocess
import sys

COVERAGE_FILE = "__coverage"
THRESHOLD = 80.0
FUNC_COVERAGE_RE = re.compile(r"(.+):\s+(\w+)\s+([\d.]+)%")

def run_go_test():
    result = subprocess.run(["go", "test", f"-coverprofile={COVERAGE_FILE}", "./..."])
    if result.returncode != 0:
        sys.exit("go test failed")

def filter_coverage():
    with open(COVERAGE_FILE, "r") as f:
        lines = f.readlines()

    with open(COVERAGE_FILE, "w") as f:
        for line in lines:
            if ".gen.go:" in line or ".mock.go:" in line:
                continue
            if platform.system().lower() != "linux" and "pathcheck.go:" in line:
                continue

            f.write(line)

def check_func_coverage():
    result = subprocess.run(["go", "tool", "cover", f"-func={COVERAGE_FILE}"], capture_output=True, text=True)
    if result.returncode != 0:
        sys.exit("cover tool failed")

    print(result.stdout, flush=True)

    for line in result.stdout.splitlines():
        match = FUNC_COVERAGE_RE.match(line.strip())
        if not match:
            continue

        location = match.group(1)
        func = match.group(2)
        percent = float(match.group(3))
        if percent < THRESHOLD:
            print(f"::error ::Test coverage too low: {percent}% (threshold: {THRESHOLD}%) in {location}::{func}", flush=True)
            sys.exit(1)
    
    print(f"::notice ::Test coverage OK: all targets >= {THRESHOLD}% threshold", flush=True)

def show_html_report():
    subprocess.run(["go", "tool", "cover", f"-html={COVERAGE_FILE}"])

def main():
    if len(sys.argv) != 2 or sys.argv[1] not in ("func", "html"):
        sys.exit("Usage: check_coverage.py [func|html]")

    try:
        run_go_test()
        filter_coverage()

        if sys.argv[1] == "func":
            check_func_coverage()
        elif sys.argv[1] == "html":
            show_html_report()
    finally:
        if os.path.exists(COVERAGE_FILE):
            os.remove(COVERAGE_FILE)

if __name__ == "__main__":
    main()
