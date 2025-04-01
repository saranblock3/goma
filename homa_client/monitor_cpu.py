import psutil
import sys

if len(sys.argv) < 3:
    print("Usage: python monitor_cpu.py <pid1> <pid2>")
    sys.exit(1)

pid1 = int(sys.argv[1])
pid2 = int(sys.argv[2])

try:
    process1 = psutil.Process(pid1)
    process2 = psutil.Process(pid2)
    while True:
        cpu_percent = process1.cpu_percent(interval=0.1)
        print(f"{cpu_percent}")
except psutil.NoSuchProcess:
    print("Process with PID not found.")
except Exception as e:
    print(f"An error occurred: {e}")
