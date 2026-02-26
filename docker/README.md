### 构造镜像
- 在根目录（Makefile所在目录）执行`make build_img`

### docker部署
- 进入`docker`目录
- 修改`run.sh`脚本中的环境变量 如：IP，数据库用户名、密码等
- 执行`sh run.sh`

### k8s部署
- 进入`docker`目录
- 修改`default.env`脚本中的环境变量 如：IP，数据库用户名、密码等
- 创建secret `bash create_env_secret.sh`
- 执行`kubectl apply -f pa.yaml`
- 端口映射`bash pm_port_forward.sh`

### 首次部署
- 进入`docker`目录
- 执行 `bash upload_image.sh`
- 执行 `bash upload_deploy.sh`
- 进入云服务器
- 执行 `bash create_env_secret.sh`
- 执行 `kubectl apply -f pa.yaml`
- 执行 `bash pm_port_forward.sh` 推荐screen内执行