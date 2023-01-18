generate-rsa-private-key-access:
	openssl genrsa -out access-private.pem 2048
generate-rsa-public-key-access:
	openssl rsa -in access-private.pem -outform PEM -pubout -out access-public.pem
generate-rsa-private-key-refresh:
	openssl genrsa -out refresh-private.pem 2048
generate-rsa-public-key-refresh:
	openssl rsa -in refresh-private.pem -outform PEM -pubout -out refresh-public.pem

generate-new-rsa-keys:
	make generate-rsa-private-key-access
	make generate-rsa-private-key-refresh
	make generate-rsa-public-key-access
	make generate-rsa-public-key-refresh