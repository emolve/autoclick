

## 注册服务
sc.exe create autoClickService binPath= "D:\Project\emolve\autoclick\autoclick.exe"
## 启动服务
sc.exe start autoClickService
## 停止服务
sc.exe stop autoClickService



# off work
## 构建
go build -o autoOff.exe

## 注册服务
sc.exe create autoOffWork binPath= "D:\Project\emolve\autoclick\autoOff.exe"
## 启动服务
sc.exe start autoOffWork
## 停止服务
sc.exe stop autoOffWork