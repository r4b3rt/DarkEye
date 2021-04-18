# 基本的增删改查

## 插入节点。插入一个Person类别的节点，且这个节点有一个属性name，属性值为Andres

  CREATE (n:Person {name : 'Andres'});
## 插入边。插入一条a到b的有向边，且边的类别为Follow

  MATCH (a:Person),(b:Person)
  WHERE a.name = 'Node A' AND b.name = 'Node B'
  CREATE (a)-[r:Follow]->(b);
## 更新节点。更新一个Person类别的节点，设置新的name。

  MATCH (n:Person { name: 'Andres' })
  SET n.name = 'Taylor';
## 删除节点。Neo4j中如果一个节点有边相连，是不能单单删除这个节点的。

  MATCH (n:Person { name:'Taylor' })
  DETACH DELETE n;
## 删除边。

  MATCH (a:Person)-[r:Follow]->(b:Person)
  WHERE a.name = 'Node A' AND b.name = 'Node B'
  DELETE r;
## 查询最短路径。

  MATCH (ms:Person { name:'Node A' }),(cs:Person { name:'Node B' }), p = shortestPath((ms)-[r:Follow]-(cs)) RETURN p;
## 查询两个节点之间的关系。

  MATCH (a:Person { name:'Node A' })-[r]->(b:Person { name:'Node B' })
  RETURN type(r);
## 查询一个节点的所有Follower。

  MATCH (:Person { name:'Taylor' })-[r:Follow]->(Person)
  RETURN Person.name;
