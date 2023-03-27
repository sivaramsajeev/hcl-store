default: run

build: fmt
	go build -o terraform-provider-student

fmt:
	go fmt ./... 

run: build move_plugin tf_version start_server start_redis 
	cd tf && terraform init && terraform plan 

start_server: fmt
	go run api/main.go &

start_redis:
	sudo docker run -p 6379:6379 redis &

tf_version:
	curl -L https://raw.githubusercontent.com/warrensbox/terraform-switcher/release/install.sh | sudo bash && tfswitch 0.11.15 && export PATH=$PATH:~/bin/

move_plugin:
	mkdir -p ~/.terraform.d/plugins/ && mv terraform-provider-student ~/.terraform.d/plugins/


