# 引言
## 编写目的
StreamNet版DAG实现是内存级别的实现，不涉及持久化，每次重启需要从RocksDB中把历史数据加载到内存中。且整个StreamNet不仅在网络层，逻辑层也依赖Iota的实现。因此，新的优化版本将业务数据持久化与DAG持久化解藕，业务数据不需要关心trunk、branch等DAG的实现，而DAG也只关心业务层数据的hash结果。且不再依赖Iota的模型及网络。
## 范围
本文档重点描述DAG的实现和持久化，不关心业务层情况。不实现pow和broadcast。

## 定义
- DAG：有向无环图，以一个被称为genesis的节点作为起点，每个加入到该图重的节点必须指定一个parent节点和reference节点。
- DAG持久化：
## 参考资料
- StreamNet工程源代码

# 总体设计：
DAG主要分为两个部分，内存组件和持久化组件。score、totalOrder等数据不持久化只存在于内存中，DAG的结构体（含hash、trunk、branch）持久化到本地。
## DAG内存组件
内存组件与StreamNet实现基本无差异

## DAG持久化组件
当节点被确认增加DAG中时，需要将对应的节点持久化到本地。不限定持久化层使用rocksDB或是其他KV数据库甚至Mysql，只久化层主要以接口形式提供持久化和查询功能，可以有多种实现。

# 数据结构：
  参照[DAG分层结构](./DAG分层结构.md)

# 接口设计：
## Add
### param 
- hash

### return
- void

## TotalOrder
### param
无

### return 
- List<Hash>

# 实现原理：
同StreamNet，参见[StreamNet](https://github.com/triasteam/StreamNet/blob/dev/document/yellow_paper/StreamNet/StreamNet.pdf)
