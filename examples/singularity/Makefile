build-sin-hello:
	mkdir -p  ~/singularity/scalebox/
	rm -f ~/singularity/scalebox/hello-scalebox.sif
	singularity build ~/singularity/scalebox/hello-scalebox.sif docker://hub.cstcloud.cn/scalebox/hello-scalebox
	ssh login-node mkdir -p singularity/scalebox
	scp ~/singularity/scalebox/hello-scalebox.sif login-node:singularity/scalebox

build-sin-calc:
	mkdir -p  ~/singularity/app-primes/
	rm -f ~/singularity/app-primes/calc.sif
	singularity build ~/singularity/app-primes/calc.sif docker-daemon://app-primes/calc:latest
	ssh login mkdir -p singularity/scalebox
	scp ~/singularity/app-primes/calc.sif login:singularity/app-primes
