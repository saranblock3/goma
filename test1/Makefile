main:
	number=0 ; while [[ $$number -le 500 ]] ; do \
    	go run main.go $$number > out$$number & \
        ((number = number + 1)) ; \
    done
