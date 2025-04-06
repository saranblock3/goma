import psutil
import sys
import signal

def alarm_handler(signum, frame):
    sys.exit(1)


signal.signal(signal.SIGALRM, alarm_handler)

signal.alarm(5)

try:
    while True:
        cpu_percent = psutil.cpu_percent(interval=0.1)
        print(f"{cpu_percent}")
except psutil.NoSuchProcess:
    print("Process with PID not found.")
except Exception as e:
    print(f"An error occurred: {e}")
