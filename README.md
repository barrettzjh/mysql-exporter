# mysql-exporter
prometheus exporter for mysql


```
# listen port
port: "8024"

# endpoint
endpoint: "/metrics"

#mysql configure
mysqlUrl: "USER:PASSWD@tcp(172.18.39.26:3306)/mysql;root:123456@tcp(172.18.39.9:3306)/mysql"


docker
cd docker
docker build -t mysql-exporter:1.0 .

docker run -e port="8024" -e endpoint="/metrics" -e mysqlUrl="USER:PASSWD@tcp(172.18.39.26:3306)/mysql;root:123456@tcp(172.18.39.9:3306)/mysql" mysql-exporter:1.0 -d


k8s
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mysql-exporter
  namespace: monitoring
  labels:
    app: mysql-exporter
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mysql-exporter
  template:
    metadata:
      labels:
        app: mysql-exporter
    spec:
      imagePullSecrets:
      - name: standard
      containers:
      - name: mysql-exporter
        image: mysql-exporter:1.1
        imagePullPolicy: "Always"
        ports:
        - containerPort: 8024
        env:
        - name: port
          value: "8024"
        - name: endpoint
          value: "/metrics"
        - name: mysqlUrl
          value: "USER:PASSWD@tcp(172.18.39.26:3306)/mysql;root:123456@tcp(172.18.39.9:3306)/mysql"
---
apiVersion: v1
kind: Service
metadata:
  name: mysql-exporter
  namespace: monitoring
  labels:
    prometheus: mysql
spec:
  ports:
    - name: http
      port: 8024
  selector:
    app: mysql-exporter

#prometheus-operator
---
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  labels:
    app: mysql-exporter
  name: mysql
  namespace: monitoring
spec:
  endpoints:
  - interval: 30s
    port: http
  selector:
    matchLabels:
      prometheus: mysql
```
