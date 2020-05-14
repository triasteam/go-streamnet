### GRPC-demo

&emsp;&emsp;项目地址：https://github.com/triasteam/StreamNet-go
&emsp;&emsp;目录结构：./StreamNet-go/utils/demos/grpc
##### proto文件
&emsp;&emsp;首先定义一个proto文件,proto文件支持跨语言使用，golang中在proto目录下执行如下命令，会自动生成service.pb.go文件，供程序调用。
    
    protoc --go_out=plugins=grpc:. service.proto

``` proto
syntax = "proto3";
package proto;
option java_package = "com.test.rpc";
option java_multiple_files = false;
 
message SayHelloRequest{
    bytes name=1;
}
 
message SayHelloResponse{
    bytes result=1;
}
 
service SayHelloService{
    rpc SayHello(SayHelloRequest) returns (SayHelloResponse);
}

```

#### golang

##### server

&emsp;&emsp;grpc-golang-complete/client/main.go，默认端口41005，可以自行调整

##### client

&emsp;&emsp;grpc-golang-complete/client/main.go，执行后获得返回结果

    hello :Tony 

#### java

&emsp;&emsp;首先新建maven工程，添加相关grpc依赖如下

``` xml
    <dependency>
		<groupId>io.grpc</groupId>
		<artifactId>grpc-netty-shaded</artifactId>
		<version>1.20.0</version>
	</dependency>
	<dependency>
		<groupId>io.grpc</groupId>
		<artifactId>grpc-protobuf</artifactId>
		<version>1.20.0</version>
	</dependency>
	<dependency>
		<groupId>io.grpc</groupId>
		<artifactId>grpc-stub</artifactId>
		<version>1.20.0</version>
	</dependency>

```

&emsp;&emsp;添加build代码，指定proto文件路径

``` xml
    <build>
		<extensions>
			<extension>
				<groupId>kr.motd.maven</groupId>
				<artifactId>os-maven-plugin</artifactId>
				<version>1.5.0.Final</version>
			</extension>
		</extensions>
		<plugins>
			<plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-compiler-plugin</artifactId>
                <version>2.3.2</version>
                <configuration>
                    <source>1.8</source>
                    <target>1.8</target>
                    <encoding>UTF-8</encoding>
                </configuration>
            </plugin>
			<plugin>
				<groupId>org.xolstice.maven.plugins</groupId>
				<artifactId>protobuf-maven-plugin</artifactId>
				<version>0.5.1</version>
				<configuration>
					<protocArtifact>com.google.protobuf:protoc:3.7.1:exe:${os.detected.classifier}</protocArtifact>
					<pluginId>grpc-java</pluginId>
					<pluginArtifact>io.grpc:protoc-gen-grpc-java:1.20.0:exe:${os.detected.classifier}</pluginArtifact>
					<protoSourceRoot>src/main/resources/proto</protoSourceRoot>
				</configuration>
				<executions>
					<execution>
						<goals>
							<goal>compile</goal>
							<goal>compile-custom</goal>
						</goals>
					</execution>
				</executions>
			</plugin>
		</plugins>
	</build>

```

&emsp;&emsp;将开始的service.proto文件放入src/main/resources/proto目录内，执行mvn clean install，会在工程里面生成相关rpc基础代码在目录target/generated-sources/protobuf/java/com/test/rpc和target/generated-sources/protobuf/grpc-java/com/test/rpc内

&emsp;&emsp;保持golang服务开启，执行DemoTest的main方法，获得服务端响应如下:

    hello :倪明

#### python

&emsp;&emsp;首先安装python依赖包:

	pip install grpcio
	pip install protobuf
	pip install grpcio_tools

&emsp;&emsp;使用文章开头定义好的proto文件，生成对应的py代码，包含service_pb2.py和service_pb2_grpc.py两个文件

	python -m grpc_tools.protoc -I ./ --python_out=./ --grpc_python_out=./ service.proto

&emsp;&emsp;保持grpc服务启动，执行main.py程序（python main.py），获得返回结果

	received: hello :Tony Stack