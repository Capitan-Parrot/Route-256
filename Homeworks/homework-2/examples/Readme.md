### CURL запросы для проверки Server
```
curl -X GET "http://localhost:9000?status=goes&how=well"
curl -X DELETE "http://localhost:9000"
curl -X PUT "http://localhost:9000" -d "JokeAboutPhp"
curl -X PUT "http://localhost:9000"
curl -X POST "http://localhost:9000" -H "hw-sum: 95"
curl -X POST "http://localhost:9000" -H "hw-sum: 3.5"
curl -X POST "http://localhost:9000"
```
### CURL запросы для проверки ServerWithData
```
curl -v -X GET "http://localhost:9001?id=1"
curl -v -X DELETE "http://localhost:9001?id=1"
curl -v -X PUT "http://localhost:9001" -d '{ \"id\": 1, \"value\": \"example22222\" }'
curl -v -X POST "http://localhost:9001" -d '{ \"id\": 1, \"value\": \"example1\" }'
curl -v -X GET "http://localhost:9001?id=1"
curl -v -X GET "http://localhost:9001?id=error"
curl -v -X PUT "http://localhost:9001" -d '{ \"id\": 1, \"value\": \"example22\" }'
curl -v -X PUT "http://localhost:9001" -d '{ \"id\": error, \"value\": \"example22\" }'
curl -v -X PUT "http://localhost:9001" -d '{ \"id\": 1 }'
curl -v -X GET "http://localhost:9001?id=1"
curl -v -X DELETE "http://localhost:9001?id=1"
curl -v -X GET "http://localhost:9001?id=1"
curl -v -X POST "http://localhost:9001" -d '{ \"id\": 1, \"value\": \"error\" }'
```