**Пример запросов для тестирования Student:**
```
curl -v -X POST "http://localhost:9000/student" -d "{ \"name\": \"Andrew\", \"courseProgram\": \"Go\" }"
curl -v -X GET "http://localhost:9000/student?studentId=1"
curl -v -X PUT "http://localhost:9000/student" -d "{ \"id\": 1, \"name\": \"Andrew\", \"courseProgram\": \"Go Junior\" }"
curl -v -X GET "http://localhost:9000/student?studentId=1"
```

**Пример запросов для тестирования Task:**
```
curl -v -X GET "http://localhost:9000/task"
```

**Пример запросов для тестирования Solution:**
```
curl -v -X POST "http://localhost:9000/solution" -d "{ \"studentId\": 1, \"taskId\": 1}"
curl -v -X GET "http://localhost:9000/solution?id=1"
```