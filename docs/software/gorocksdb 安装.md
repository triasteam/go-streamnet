go rocksdb 安装步骤说明

这里使用的操作系统是 Ubuntu 18.04

1. 下载依赖
```bash
sudo apt update
sudo apt install -y zlib1g zlib1g-dev ubuntu-snappy-cli snapd libsnappy-dev bzip2 libbz2-dev zstd libzstd-dev \
liblz4-dev liblz4-tool 
```

2. 下载 RocksDB 发布包
这里选择的是最新的 v6.10.2 (6/5/2020)
```bash
wget https://github.com/facebook/rocksdb/archive/v6.10.2.tar.gz
tar xf v6.10.2.tar.gz
cd rocksdb-6.10.2
```

3. 编译库并将库和头文件复制到 /usr/include 和 /usr/lib下
```bash
make shared_lib -j9
sudo mkdir /usr/lib/
sudo cp librocksdb.so  librocksdb.so.6 librocksdb.so.6.10 /usr/lib/
cd include
sudo cp -r rocksdb /usr/include/
```

4. 安装 gorocksdb
```bash
CGO_CFLAGS="-I/usr/include" CGO_LDFLAGS="-L/usr/lib/rocksdb -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd" go get github.com/tecbot/gorocksdb
```
