main:
	@for i in $$(seq 0 102); do \
		go run main.go $$i >> latencies_$$i & \
		python3 monitor_cpu.py 58585 > cpu_usage_$$i; \
	done;
