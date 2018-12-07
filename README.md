# huang

Huang's Container Manager System.

[Kubernetes](https://kubernetes.io/) is useful, but it's tooooooo complicated. I prefer `simple is better`, so I decide
write my own container manager system. It borrows lots of good designs from Kubernetes, Thanks, Kubernetes!

它包括这些特性（名字基本沿用于Kubernetes，这样可以降低学习复杂度）：

- 内置的扁平化网络（参考flannel）
- namespace
- deployment
- pod（由huang自动管理，用户无法直接操作）
- service discovery & expose
- labels
- 健康检查
- job & cronjob

架构与Kubernetes大致保持一致：

- etcd作为存储
- master控制整个集群
- worker负责实际操作

计划是2018年春节后开始开发，此前的时间用于学习和准备。

---

参考资料：

- [kubernetes官网](https://kubernetes.io/)
- [Borg, Omega, and Kubernetes](https://storage.googleapis.com/pub-tools-public-publication-data/pdf/44843.pdf)
