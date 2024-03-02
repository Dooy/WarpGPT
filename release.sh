docker build -t hts999/warp-gpt:latest . 

# 修改镜像标签为当前日期时间
time=$(date "+%Y%m%d%H%M%S")
# 获取当前git commit id 
commit_id=$(git rev-parse HEAD)
docker tag hts999/warp-gpt:latest hts999/warp-gpt:$time-$commit_id
# 推送镜像到docker hub
docker push hts999/warp-gpt:$time-$commit_id
docker push hts999/warp-gpt:latest

#debug
# docker pull hts999/warp-gpt:latest && docker rm -f warpgpt
# docker rm -f warpgpt && docker run --name warpgpt -d -p 6050:5000 -e proxy=http://107.181.187.120:10076 -e port=5000 -e host=0.0.0.0 xiaoguaishou92/warpgpt
