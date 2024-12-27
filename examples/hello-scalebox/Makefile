
run:
	scalebox app create

build-singularity:
	singularity build -F ~/singularity/scalebox/agent.sif  docker-daemon://hub.cstcloud.cn/scalebox/agent:latest

run-singularity:
	scalebox app create -e singularity.env
