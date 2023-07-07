
# on work
```
## 构建
go build -o ./bin/cmd/autoOn.exe
## 注册服务
sc.exe create autoOnWork binPath= "D:\Project\emolve\autoclick\bin\cmd\autoOn.exe"
## 启动服务
sc.exe start autoOnWork
## 停止服务
sc.exe stop autoOnWork
```

# off work
```
## 构建
go build -o ./bin/cmd/autoOff.exe
## 注册服务
sc.exe create autoOffWork binPath= "D:\Project\emolve\autoclick\bin\cmd\autoOff.exe"
## 启动服务
sc.exe start autoOffWork
## 停止服务
sc.exe stop autoOffWork
## 删除服务
sc.exe delete autoOffWork
```


# test work
```
## 构建
go build -o ./bin/cmd/autoTest.exe
## 注册服务
sc.exe create autoTest binPath= "D:\Project\emolve\autoclick\bin\cmd\autoTest.exe"
## 启动服务
sc.exe start autoTest
## 停止服务
sc.exe stop autoTest
## 删除服务
sc.exe delete autoTest
```

