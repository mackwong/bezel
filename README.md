# Bezel

此工程是为 diamond-on-edge 提供集群配置管理的工具

[![build status](http://gitlab.bj.sensetime.com/diamond/service-providers/bezel/badges/master/pipeline.svg)](http://gitlab.bj.sensetime.com/diamond/service-providers/bezel/commits/master)
[![coverage report](http://gitlab.bj.sensetime.com/diamond/service-providers/bezel/badges/master/coverage.svg)](http://gitlab.bj.sensetime.com/diamond/service-providers/bezel/commits/master)

## 构建

```bash
make linux #or darwin/windows
```

支持构建不同架构的版本，包含 linux(amd64/arm64), macOS(darwin), windows.

## 使用方法

1. 创建 edge 设备全局配置文件

   全局配置文件可以通过下面 2 种方式生成：

   - 通过交互方式产生

     ```bash
     bezel create --output ./
     ```

     此种方式适用于节点数量比较少的情况，客户端通过逐行提示的方式引导用户配置信息：

     ```bash
     $ bezel create --output ./
     ==> please configure Name:
         demo
     ==> please configure Arranger:
         k3s
     ==> please configure UpstreamDNS:
         1.1.1.1
     ==> please configure DockerRegistry:
         2.2.2.2
     ==> please configure MachineNum:
         3
     ==> please configure MasterNum:
         1
     ==> please configure K8sMasterIP:
         3.3.3.3
     ```

    - 通过配置文件方式产生

      首先通过下面命令生成示例配置文件：

      ```bash
      $ ./bezel gen
      INFO demo.yaml is generated successfully, Please modify it and use `./bezel create -c demo.yaml` to generate edge configs
      ```

      示例配置文件格式如下：

      ```yaml
      #demo.yaml
      name: diamond-edge-ha
      machine-num: 4
      master-num: 3
      arranger: edgesite
      upstream-dns: 114.114.114.114
      docker-registry: 10.5.49.73
      k8sMaster-ip: 10.4.72.231
      ip-range:
        - ipRange: 10.4.72.1/24
          gatewayIP: 10.4.72.254
          netmask: 255.255.255.0
        - ipRange: 10.4.73.1/32
          gatewayIP: 10.4.73.254
          netmask: 255.255.255.255
      master-ip:
        - 10.4.72.1
        - 10.4.72.2
        - 10.4.73.1
      name-format: node-{{.Role}}-{{.Index}}
      hostname-format: ubuntu-{{.Role}}-{{.Index}}
      ```

      需要说明的是：

      1.  如果 master-ip 不需要特别指定，这个字段可以不写。master 会从 IP 段中自动分配
      2.  `name` 和 `hostname` 目前支持 `Role`,`Index`,`IP`三个字段。

      用户在此配置文件进行相应修改后，执行下面命令即可生成配置文件：

      ```bash
      $ bezel create -c demo.yaml
      INFO edge-config.yaml generated successfully
      INFO Sub config will write to sub/sub-edge-config-master-10.4.72.1.yaml
      INFO Sub config will write to sub/sub-edge-config-master-10.4.72.2.yaml
      INFO Sub config will write to sub/sub-edge-config-master-10.4.73.1.yaml
      INFO Sub config will write to sub/sub-edge-config-worker-10.4.72.3.yaml
      INFO sub files ./ generated successfully
      ```

2. 渲染模板文件

   使用渲染功能时，请执行：

   ```bash
   $ bezel parse  --output ./ --source ./sub/sub-edge-config-master-10.4.72.1.yaml -t /tmp/templates
   INFO parse templates successfully
   ```

   执行成功之后，`output`指定目录下会产生渲染文件。

## 测试相关

[覆盖率](http://diamond.pages.gitlab.bj.sensetime.com/service-providers/bezel/)

## Binary

[下载链接](https://gitlab.bj.sensetime.com/diamond/service-providers/bezel/-/tags/)