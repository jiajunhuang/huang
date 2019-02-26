# huang

Huang's Container Manager System.

k8s太复杂了，所以自己做一个简单版玩玩。

它包括这些特性（名字基本沿用于Kubernetes，这样可以降低学习复杂度）：

- namespace
- deployment
- pod（由huang自动管理，用户无法直接操作）
- labels

架构与Kubernetes大致保持一致：

- master控制整个集群
- worker负责实际操作

---

参考资料：

- [kubernetes官网](https://kubernetes.io/)
- [Borg, Omega, and Kubernetes](https://storage.googleapis.com/pub-tools-public-publication-data/pdf/44843.pdf)
