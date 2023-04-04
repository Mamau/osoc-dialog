#### Proto
Для работы с профилировщиком надо:
* Выставить на http сервере время WriteTimeout больше чем собираешься собирать данные
```go
package serviceprovider

import (
	"net/http"
	"time"
)

func NewHttp() *http.Server {
	return http.New(
		http.WriteTimeout(31*time.Second),
	)
}
```
Далее в роутер надо добавить
```go
package serviceprovider

import (
	"osoc-dialog/pkg/router"
	"github.com/gin-gonic/gin"
)

func NewBaseRouter() *gin.Engine {
    return router.New(
		    router.Pprof(true),
		)
}
```
* Запустить нагрузку (например wrk):
```shell
wrk -t8 -c16 -d30s 'http://localhost:8081/external/api/v1/user/1?test=123'
```
Удобный способ собрать данные
```shell
make report type=heap
```
Или можно руками собирать [нужные нам данные](https://pkg.go.dev/net/http/pprof)
```shell
go tool pprof http://localhost:8081/debug/pprof/heap
```
После этого мы попадем в интерактивный режим где нам предложат ввести команду, напишем **top**  
так, мы увидим сколько и какие функции потребляют памяти.  
Мы можем также вывести эти данные в виде изображения:
```shell
go tool pprof -png http://localhost:8081/debug/pprof/heap > out.png
```
Более удобно работать так
```shell
curl -s http://127.0.0.1:8081/debug/pprof/profile?seconds=30 > ./cpu.out
go tool pprof -http=:8080 ./cpu.out
```
Поднимется сервер на 8080 порту где можно будет в разных видах смотреть стату.  
Также можно смотреть heap|trace
#### Важно!
/debug/pprof/trace?seconds=5 выводит файл не являющийся профилем pprof.  
Трассировку можно рассмотреть с помощью
```shell
go tool trace
```