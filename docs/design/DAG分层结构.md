# 分层思路
考虑到以后的扩展需求，将DAG与具体应用进行结构。应用层将数据的hash字符串传递给DAG保存上链，仅通过hash进行交互可以最大成都扩展DAG的组件的功能。
分层后由应用层校验数据是否被篡改

## 应用层：
应用层是go-streamnet的一个基础组件，可适配不同的业务场景。

## DAG层：
是go-streamnet的核心组件，负责将应用层数据存储上链。

# 应用层结构：
## Persistable 数据内容持久化接口：
    byte[] bytes();
    void read(byte[] bytes);
    byte[] metadata();
    void readMetadata(byte[] bytes);
    boolean merge(); //预留
## Indexable key持久化接口：
    byte[] bytes();
    void read(byte[] bytes);
    Indexable incremented();
    Indexable decremented();

## Hash implement Indexable Hash持久化对象:
  sha256(内容+时间戳+nonce)

## ApprovalTransaction implement Persistable：证实交易持久化对象
```
 ApprovalTransaction{
    attestee
    attester
    score
  }
```
# DAG结构：

DAG 的 graph 结构可表示为: ```G = (B, g, P, E)```

B 表示DAG中的所有区块。目前一个区块只含有一笔交易。
g 表示genesis节点；
P 表示父引用关系；
E 表示边引用关系；
由上可得出DAG关键结构如下：

```
DAG {
  // graph from genesis to tips
  private Map<Hash, List<Hash>> graph;
  // graph from tip block to previous node, if tip is genesis, previous block is null.
  private Map<Hash, List<Hash>> revGraph;
  
  // graph of parent reference relation, from genesis to tips
  private Map<Hash, List<Hash>> parentGraph;
  // graph of parent reference relation, from tip to genesis
  private Map<Hash, Hash> parentGraph;
  
  //score used for calcultate parent block and reference block
  Map<Hash, Double> score;
  Map<Hash, Double> parentScore;
}

Hash {
    string sha256;
}
```
