all:
	docker build -f Dockerfile -t faryon93/crony:latest .

push:
	docker push faryon93/crony:latest