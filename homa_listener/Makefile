main:
	number=500 ; while [[ $$number -le 700 ]] ; do \
    	go run main.go $$number > out$$number & \
        ((number = number + 1)) ; \
    done
