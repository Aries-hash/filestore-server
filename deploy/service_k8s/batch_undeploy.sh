#!/bin/bash

kubectl delete -f svc_account.yaml 
kubectl delete -f svc_apigw.yaml
kubectl delete -f svc_dbproxy.yaml 
kubectl delete -f svc_upload.yaml 
kubectl delete -f svc_download.yaml 
kubectl delete -f svc_transfer.yaml 
# 通知去除配置
kubectl delete -f service-ingress.yaml 
