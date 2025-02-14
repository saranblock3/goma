main:
	number=500 ; while [[ $$number -le 1000 ]] ; do \
    	go run main.go $$number > out$$number & \
        ((number = number + 1)) ; \
    done
