import psutil
import sys

if len(sys.argv) < 2:
    print("Usage: python monitor_cpu.py <pid>")
    sys.exit(1)

pid = int(sys.argv[1])

try:
    process = psutil.Process(pid)
    while True:
        cpu_percent = process.cpu_percent(interval=0.1)
        print(f"{cpu_percent}")
except psutil.NoSuchProcess:
    print(f"Process with PID {pid} not found.")
except Exception as e:
    print(f"An error occurred: {e}")
