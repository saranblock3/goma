main:
	@for i in $$(seq 0 102); do \
		go run main.go host0 $$i >> latencies_$$i & \
		bin/python3 monitor_cpu.py $$pid > cpu_usage_$$i; \
	done;

host0:
	@for i in $$(seq 100 349); do \
		go run main.go host0 $$i & \
	done;

host1:
	@for i in $$(seq 100 349); do \
		go run main.go host1 $$i & \
	done;

host2:
	@for i in $$(seq 100 349); do \
		go run main.go host2 $$i & \
    done;

host3:
	@for i in $$(seq 100 349); do \
		go run main.go host3 $$i & \
    done;

host4:
	@for i in $$(seq 100 349); do \
		go run main.go host4 $$i & \
    done;
