host0:
	@for i in $$(seq 8000 8049); do \
		go run main.go 10.10.1.1 $$i & \
	done;

host1:
	@for i in $$(seq 8000 8049); do \
		go run main.go 10.10.1.2 $$i & \
	done;

host2:
	@for i in $$(seq 8000 8049); do \
		go run main.go 10.10.1.3 $$i & \
	done;

host3:
	@for i in $$(seq 8000 8049); do \
		go run main.go 10.10.1.4 $$i & \
	done;

host4:
	@for i in $$(seq 8000 8049); do \
		go run main.go 10.10.1.5 $$i & \
	done;
