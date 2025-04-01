import psutil
import sys
import signal

def alarm_handler(signum, frame):
    sys.exit(1)


signal.signal(signal.SIGALRM, alarm_handler)

if len(sys.argv) < 2:
    print("Usage: python monitor_cpu.py <pid>")
    sys.exit(1)

pid = int(sys.argv[1])

signal.alarm(20)

try:
    process1 = psutil.Process(pid)
    while True:
        cpu_percent = process1.cpu_percent(interval=0.1)
        print(f"{cpu_percent}")
except psutil.NoSuchProcess:
    print("Process with PID not found.")
except Exception as e:
    print(f"An error occurred: {e}")
