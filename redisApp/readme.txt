


Start redis:
docker run --rm -p 6379:6379 --name redis redis


Reg user
curl -i -X POST -H "Content-Type: application/json" -d '{"name": "Joe", "e-mail": "joe@yahoo.com", "pass": "pass"}' \
 http://localhost:8080/singup

Check:
curl -i -X GET http://localhost:8080/check?msg=...
