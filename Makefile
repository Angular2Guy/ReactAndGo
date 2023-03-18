hello: 
	echo "Hello"

.ONESHELL:
frontend-build: 
	cd frontend	
	npm install
	npm run build

.ONESHELL:
backend-build: 
	cd backend
	export GOGC=off
	export GOMEMLIMIT=96MiB
	#to support differen libc versions
	export CGO_ENABLED=0
	go build	

full-build: frontend-build backend-build
