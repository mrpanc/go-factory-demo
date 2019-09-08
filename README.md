# go-factory-demo

> Git: https://github.com/mrpanc/go-factory-demo

## 背景

目前遇到一个需求，支持不同场景下的定制化服务。抽象来看，不同场景下，做的事情都是类似的，只是其中某些接口的业务逻辑有细微的区别。因此，想到了通过工厂模式来创建不同场景下的对象。

## 工厂模式

工厂模式是一种创建型模式，在创建对象时不会对客户端暴露创建逻辑，并且是通过一个公共的接口来指向新创建的对象。因此，当我们提供一个产品类库时，可以不需要暴露内部实现，只显示他们的接口，提高了代码的封装性。通过工厂模式，当一个调用者想创建一个对象，只需要知道其名称即可。并且提高了代码的可扩展性，如果需要适配另一个场景，只需要扩展一个工厂类即可。

## 现实场景

现在我们要适配一个名为OPS的服务，该服务主要功能是向OPS发送心跳、从OPS拉取配置以及通过OPS下发更新包。假设目前有2个场景，一个基础场景，一个特殊场景。


### 接口定义

首先，我们将公共的接口抽出来，定义为一个公共的接口。例如目前的场景下，定义接口`OpsImpl`如下
```go
type OpsImpl interface {
    SendHeartbeat() error
    DoUpdate() error
    DoConfigUpload() error
}
```

### 接口实现

接下来，我们针对基础场景和特殊场景，创建`struct`实现接口。首先是基础场景`baseOps`实现如下：

```go
type baseOps struct {
    postUrl string
}

func (base *baseOps) SendHeartbeat() error {
    fmt.Println("BaseOps: Send heartbeat")
    fmt.Println("BaseOps: send to url: ", base.postUrl)
    return nil
}
func (base *baseOps) DoUpdate() error {
    fmt.Println("BaseOps: DoUpdate")
    return nil
}
func (base *baseOps) DoConfigUpload() error {
    fmt.Println("BaseOps: DoConfigUpload")
    return nil
}
```

接下来是特殊场景，特殊场景是基础场景的子类，我们通过组合的模式实现该逻辑，其中特殊场景的更新逻辑与基础场景保持一致，因此我们可以省略其实现，复用父类的实现，而心跳和配置拉取逻辑不一致，需要重新写业务逻辑，具体实现如下：

```go

type specialOps struct {
    baseOps
    sendConfig bool
}

func (ops *specialOps) SendHeartbeat() error {
    fmt.Println("SpecialOps: SendHeartbeat")
    return nil
}
func (ops *specialOps) DoConfigUpload() error {
    fmt.Println("SpecialOps: DoConfigUpload")
    if ops.sendConfig {
        fmt.Println("SpecialOps: upload config")
    } else {
        fmt.Println("SpecialOps: no need to upload config")
    }
    return nil
}

// not implement DoUpdate, same to baseOps
```

### 工厂函数

接下来，我们为不同的实现工厂方法，返回同样的接口。这些工厂方法接受同样的参数，我们将参数通过`map[string]interface{}`进行传递。实现如下：

```go

type OpsFactory func(conf map[string]interface{}) (OpsImpl, error)

func NewBaseOps(conf map[string]interface{}) (OpsImpl, error) {
    fmt.Println("BaseOps: Create")
    postUrl, ok := conf["PostUrl"]
    if !ok {
        return nil, errors.New("[postUrl] has not been set in config map")
    }
    return &baseOps{
        postUrl: postUrl.(string),
    }, nil
}

func NewSpecialOps(conf map[string]interface{}) (OpsImpl, error) {
    fmt.Println("specialOps: Create")
    sendConfig, ok := conf["SendConfig"]
    if !ok {
        return nil, errors.New("[SendConfig] has not been set in config map")
    }
    return &specialOps{
        sendConfig: sendConfig.(bool),
    }, nil
}
```

### 注册工厂

接下来，我们通过一个公共的函数`RegisterOps`来注册需要用到的工厂，并且通过初始化函数`init`在程序启动前注册这两个工厂。具体实现如下：

```go
var opsFactories = make(map[opsTypeEnum]OpsFactory)

func RegisterOps(opsType opsTypeEnum, factory OpsFactory) {
    if factory == nil {
        panic(fmt.Sprintf("Ops factory %s does not exist", string(opsType)))
    }
    _, ok := opsFactories[opsType]
    if ok {
        fmt.Printf("Ops factory %s has been registered already\n", string(opsType))
    } else {
        fmt.Printf("Register ops factory %s\n", string(opsType))
        opsFactories[opsType] = factory
    }
}

func init() {
    RegisterOps(BaseType, NewBaseOps)
    RegisterOps(SpecialType, NewSpecialOps)
}
```

### 创建工厂

最后我们通过以下函数，就可以方便的创建工厂，返回对应的`Ops`接口。
```go
func CreateOps(conf map[string]interface{}) (OpsImpl, error) {
    opsType, ok := conf["OpsType"]
    if !ok {
        fmt.Println("No ops type in config map. Use base ops type as default.")
        opsType = BaseType
    }
    OpsFactory, ok := opsFactories[opsType.(opsTypeEnum)]
    if !ok {
        availableOps := make([]string, len(opsFactories))
        for k, _ := range opsFactories {
            availableOps = append(availableOps, string(k))
        }
        return nil, errors.New(fmt.Sprintf("Invalid ops type. Must be one of: %s", strings.Join(availableOps, ", ")))
    }
    fmt.Println("Create ops: ", opsType)
    return OpsFactory(conf)
}
```

### 测试函数

最终，我们可以通过以下方式简单的创建不同场景的接口。

```go
func main() {
	baseOps, err := ops.CreateOps(map[string]interface{}{
		"OpsType": ops.BaseType,
		"PostUrl": "http://ops.cloud.com/send_heartbeat",
	})
	if err != nil {
		fmt.Println("create baseOps failed, err: ", err.Error())
		return
	}
	baseOps.DoConfigUpload() // Output: BaseOps: DoConfigUpload
	baseOps.DoUpdate()       // Output: BaseOps: DoUpdate
	baseOps.SendHeartbeat()  // Output: BaseOps: Send heartbeat

	specialOps, err := ops.CreateOps(map[string]interface{}{
		"OpsType":    ops.SpecialType,
		"SendConfig": true,
	})
	if err != nil {
		fmt.Println("create specialOps failed, err: ", err.Error())
		return
	}
	specialOps.DoConfigUpload() // Output: SpecialOps: DoConfigUpload
	specialOps.DoUpdate()       // Output: BaseOps: DoUpdate
	specialOps.SendHeartbeat()  // Output: SpecialOps: SendHeartbeat
}
```