Windows PowerShell
(C) Корпорация Майкрософт (Microsoft Corporation). Все права защищены.

Попробуйте новую кроссплатформенную оболочку PowerShell (https://aka.ms/pscore6)

PS C:\Users\Андрей\GolandProjects\Ozon\Homework2> rm alias:curl
PS C:\Users\Андрей\GolandProjects\Ozon\Homework2> curl -v -X GET "http://localhost:9001?id=1"                        
Note: Unnecessary use of -X or --request, GET is already inferred.
*   Trying 127.0.0.1:9001...                                      
* Connected to localhost (127.0.0.1) port 9001 (#0)               
> GET /?id=1 HTTP/1.1                                             
> Host: localhost:9001                                            
> User-Agent: curl/7.83.1                                         
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 404 Not Found
< Date: Thu, 09 Mar 2023 10:36:26 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact
PS C:\Users\Андрей\GolandProjects\Ozon\Homework2> curl -v -X DELETE "http://localhost:9001?id=1"
*   Trying 127.0.0.1:9001...
* Connected to localhost (127.0.0.1) port 9001 (#0)
> DELETE /?id=1 HTTP/1.1
> Host: localhost:9001
> User-Agent: curl/7.83.1
> Accept: */*
>
< HTTP/1.1 404 Not Found
< Date: Thu, 09 Mar 2023 10:36:26 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
PS C:\Users\Андрей\GolandProjects\Ozon\Homework2> curl -v -X PUT "http://localhost:9001" -d '{ \"id\": 1, \"value\": \"example22222\" }'
*   Trying 127.0.0.1:9001...
* Connected to localhost (127.0.0.1) port 9001 (#0)
> PUT / HTTP/1.1
> Host: localhost:9001
> User-Agent: curl/7.83.1
> Accept: */*
> Content-Length: 36
> Content-Type: application/x-www-form-urlencoded
>
* Mark bundle as not supporting multiuse
< Date: Thu, 09 Mar 2023 10:36:26 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
PS C:\Users\Андрей\GolandProjects\Ozon\Homework2> curl -v -X POST "http://localhost:9001" -d '{ \"id\": 1, \"value\": \"example1\" }'
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:9001...
* Connected to localhost (127.0.0.1) port 9001 (#0)
> POST / HTTP/1.1
> Host: localhost:9001
> User-Agent: curl/7.83.1
> Accept: */*
> Content-Length: 32
> Content-Type: application/x-www-form-urlencoded
>
< HTTP/1.1 200 OK
< Date: Thu, 09 Mar 2023 10:36:26 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
PS C:\Users\Андрей\GolandProjects\Ozon\Homework2> curl -v -X GET "http://localhost:9001?id=1"
Note: Unnecessary use of -X or --request, GET is already inferred.
*   Trying 127.0.0.1:9001...
* Connected to localhost (127.0.0.1) port 9001 (#0)
> GET /?id=1 HTTP/1.1
> Host: localhost:9001
> User-Agent: curl/7.83.1
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Thu, 09 Mar 2023 10:36:26 GMT
< Content-Length: 8
< Content-Type: text/plain; charset=utf-8
<
example1* Connection #0 to host localhost left intact
PS C:\Users\Андрей\GolandProjects\Ozon\Homework2> curl -v -X GET "http://localhost:9001?id=error"
Note: Unnecessary use of -X or --request, GET is already inferred.
*   Trying 127.0.0.1:9001...
* Connected to localhost (127.0.0.1) port 9001 (#0)
> GET /?id=error HTTP/1.1
> Host: localhost:9001
> User-Agent: curl/7.83.1
> Accept: */*
>
< HTTP/1.1 500 Internal Server Error
< Date: Thu, 09 Mar 2023 10:36:26 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
PS C:\Users\Андрей\GolandProjects\Ozon\Homework2> curl -v -X PUT "http://localhost:9001" -d '{ \"id\": 1, \"value\": \"example22\" }'
*   Trying 127.0.0.1:9001...
* Connected to localhost (127.0.0.1) port 9001 (#0)
> PUT / HTTP/1.1
> Host: localhost:9001
> User-Agent: curl/7.83.1
> Accept: */*
> Content-Length: 33
> Content-Type: application/x-www-form-urlencoded
>
< HTTP/1.1 200 OK
< Date: Thu, 09 Mar 2023 10:36:26 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
PS C:\Users\Андрей\GolandProjects\Ozon\Homework2> curl -v -X PUT "http://localhost:9001" -d '{ \"id\": error, \"value\": \"example22\" }'
*   Trying 127.0.0.1:9001...
* Connected to localhost (127.0.0.1) port 9001 (#0)
> PUT / HTTP/1.1
> Host: localhost:9001
> User-Agent: curl/7.83.1
> Accept: */*
> Content-Length: 37
> Content-Type: application/x-www-form-urlencoded
>
< HTTP/1.1 500 Internal Server Error
< Date: Thu, 09 Mar 2023 10:36:26 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
PS C:\Users\Андрей\GolandProjects\Ozon\Homework2> curl -v -X PUT "http://localhost:9001" -d '{ \"id\": 1 }'
*   Trying 127.0.0.1:9001...
* Connected to localhost (127.0.0.1) port 9001 (#0)
> PUT / HTTP/1.1
> Host: localhost:9001
> User-Agent: curl/7.83.1
> Accept: */*
> Content-Length: 11
> Content-Type: application/x-www-form-urlencoded
>
< HTTP/1.1 400 Bad Request
< Date: Thu, 09 Mar 2023 10:36:26 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
PS C:\Users\Андрей\GolandProjects\Ozon\Homework2> curl -v -X GET "http://localhost:9001?id=1"
Note: Unnecessary use of -X or --request, GET is already inferred.
*   Trying 127.0.0.1:9001...
* Connected to localhost (127.0.0.1) port 9001 (#0)
> GET /?id=1 HTTP/1.1
> Host: localhost:9001
> User-Agent: curl/7.83.1
> Accept: */*
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Thu, 09 Mar 2023 10:36:26 GMT
< Content-Length: 9
< Content-Type: text/plain; charset=utf-8
<
example22* Connection #0 to host localhost left intact
PS C:\Users\Андрей\GolandProjects\Ozon\Homework2> curl -v -X DELETE "http://localhost:9001?id=1"
*   Trying 127.0.0.1:9001...
* Connected to localhost (127.0.0.1) port 9001 (#0)
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Thu, 09 Mar 2023 10:36:26 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
PS C:\Users\Андрей\GolandProjects\Ozon\Homework2> curl -v -X GET "http://localhost:9001?id=1"
Note: Unnecessary use of -X or --request, GET is already inferred.
*   Trying 127.0.0.1:9001...
* Connected to localhost (127.0.0.1) port 9001 (#0)
> GET /?id=1 HTTP/1.1
> Host: localhost:9001
> User-Agent: curl/7.83.1
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 404 Not Found
< Date: Thu, 09 Mar 2023 10:36:27 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
PS C:\Users\Андрей\GolandProjects\Ozon\Homework2> curl -v -X POST "http://localhost:9001" -d '{ \"id\": 1, \"value\": \"error\" }'





