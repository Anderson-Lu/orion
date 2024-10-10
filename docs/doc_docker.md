# 容器性能优化

在Golang的GPM模型中P的数量决定了并行的G的数量,而P的数量又是间接依赖M的数量,M的数量又由CPU的核心数量决定.在Golang中,通过`runtime.NumCPU()`获取宿主机的CPU核数,再通过`runtime.GOMAXPROCS(n)`设置GOMAXPROCS的值.当服务部署在容器中时,每个容器拿到的都是宿主机的核数,因此,当容器数量过多时,会导致产生过多的P,从而导致频繁的线程切换,最终导致服务性能下降.

因此,我们需要通过动态限制容器中GOMAXPROCS的值来避免上述问题.容器化是通过cgroup机制来限制容器能使用的cpu核心数的,因此,通过读取虚拟化为容器分配的cpu核数来为golang程序动态设置GOMAXPROCS的值. 这里引入`go.uber.org/automaxprocs`来实现容器内部署的性能优化.