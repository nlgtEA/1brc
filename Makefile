bench:
	@ go test -bench Evaluate -benchtime=1x -count=10 -cpu=4 | tee ./bench_results/bench_$(shell date +"%Y-%m-%d-%H-%M-%S").txt