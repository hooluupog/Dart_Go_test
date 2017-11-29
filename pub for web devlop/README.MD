### 使用pub来开发web应用：  
`pub build`　　生成可部署的web应用  
`pub serve`　　提供web应用本地服务端localhost  
`pub build`　　默认为release模式，使用dart2js编译器编译最小化js输出；  
`pub serve`　　默认为debug模式。  
可以手动指定模式：    
`pub build/serve --mode=debug/release`  
要指定使用的编译器，有两种办法：  
* 在puspec.yaml文件中指定，例如：
  ```yaml
  web:
    compiler:
      debug: dartdevc
  ```
* 使用命令行参数：  
  `pub build/serve --web-compiler=dartdevc/dart2js`
