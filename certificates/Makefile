server-key:
	@rm -rf server.key
	@openssl genrsa -out server.key 2048

certificate: 
	@rm -rf server.crt
	@openssl req -new -x509 -newkey rsa:2048 -sha256 -nodes -keyout server.key -days 3560 -out server.crt -config certificate.conf