# my-life
一个生活管理系统，目前支持多人一起记账和查账。
## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=here-tunan/my-life&type=Date)](https://www.star-history.com/#here-tunan/my-life&Date)
## 功能截图
### 登录页面
![image](https://github.com/user-attachments/assets/77175dc7-0d59-4301-b25f-51dd99763921)

### 首页 （个人基本信息和家庭基本信息）
![image](https://github.com/user-attachments/assets/fa22c983-9535-41c6-bf2e-f365d8196104)

### 记账页面
<img width="1920" height="968" alt="image" src="https://github.com/user-attachments/assets/a45769d3-4318-43ec-85d0-0e3909eba88b" />


### 快速手工记账
<img width="1920" height="968" alt="image" src="https://github.com/user-attachments/assets/b4bf9445-d119-4b21-afdb-127b439785aa" />


### 支持微信/支付宝账单导入（且导入时会根据描述信息，基于ES的分词功能和历史记账信息匹配出最优类型）
<img width="1920" height="968" alt="image" src="https://github.com/user-attachments/assets/7538824c-5f59-4680-bbc9-38a521270fed" />

<img width="1920" height="968" alt="image" src="https://github.com/user-attachments/assets/5671c286-d9d6-46fc-a7e8-21bb1e7d902d" />


### 个人收支统计
![image](https://github.com/here-tunan/my-life/assets/40956738/d0a0cc80-eadc-4227-8bb9-b34ba27a27bb)

### 家庭收支统计（可以看到家庭内的其他成员的收支明细）
![image](https://github.com/here-tunan/my-life/assets/40956738/e3c28bef-d02a-474c-b9b1-da3c75a6dac4)

## 运行
### 前提
需要开启 redis、mysql、ElasticSearch服务。

ES 索引构建
```shell
PUT /transaction_index 
{
  "settings": {
    "number_of_shards": 1
  },
  "aliases":{
    "transaction": {}
  },
  "mappings": {
    "properties": {
      "id": {"type": "keyword"},
      "description": {"type": "text"},
      "category": {"type": "keyword"},
      "type": {"type": "keyword"}
    }
  }
}
```

### 本地运行
#### 服务端
设置好dev.yaml中的配置，启动go项目即可。
#### 客户端
```shell
cd ui
npm install
npm run dev
```

### 服务器Linux运行
#### 服务端
在带有go环境的机器上进行编译打包上传到服务器上，
```shell
CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build -o go-my-life main.go
```
或者直接使用仓库中的可执行文件。

在可执行文件的同级/env目录下配置prod.yaml文件，然后配置好环境变量：
```shell
export GO_MY_LIFE_ENV=prod
```

启动后台程序
```shell
nohup ./go-my-life &
```

#### 客户端
通过vue生成的dist静态文件部署。 创建/ui/.env.production环境配置文件，修改对应的服务端地址。在带有npm的环境中执行命令生成dist，上传到服务器的位置。
```shell
npm run build
```
通过nginx做好转发即可。

#### Mac 本地启动
如果你想在本地启动，又正好是Mac系统，那完全可以自己编排好脚本文件，然后使用mac的Automator去实现这个功能，稍后我会出一篇教程来说明如何实现。
<img width="688" height="356" alt="image" src="https://github.com/user-attachments/assets/025e5ee2-a0fb-474e-b23b-009c5eba2100" />

