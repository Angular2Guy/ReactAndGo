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
	go build	

full-build: frontend-build backend-build
