# GPU Manager

In k8s can use gpu plugin to manage and schedule gpus among the cluster, but there are no solutions for the swarm cluster. This repo is a plugin can manage and schedule gpu resources among the swarm cluster.

## The parallel calculate in gpu

Two ways are supplied for gpu parallel calculate.
1. Multithreading
2. Multiprocessing

### Multithreding

In Multithreading mode, multi threads share one gpu: the gpu memory and cuda cores are shared by them all. In this way, the gpu really can calculate at the same time. That means, eg. thread A and thread B use the same gpu and calculate simultaneously, if A not use all the gpu memory and cuda cores B can use the rest of them. 

The benefit is the GPU can be 100% utilized, but the weakness also horrible: the crash happened in A thread will cause all the other calculate threads crash.

By the way, if the GPU use MPS mode, the multi processes will be transfered to Multithreading mode. So the MPS mode also has the same problem mentioned upper.

By the way, the multithreading mode can also use "streaming mode" of this repo -> [triton](https://github.com/triton-inference-server/server) 

### Multiprocessing

In Multiprocessing mode, multi processes share one gpu. Eg, process A and process B use the same gpu, in this way the cuda cores only calculate one process at the same time. This mode like "Round-Robin" in CPU, the GPU only do one calculate task at the same time, even though the task not use all the cuda cores or memory.

The weakness of this mode is the GPU not fully be used, but also has benefit, the crash happened in process A won't cause process B crash.

