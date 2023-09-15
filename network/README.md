Ether Type Refs:

> http://en.wikipedia.org/wiki/Ethertype

IP Protocol Refs:
> http://en.wikipedia.org/wiki/List_of_IP_protocol_numbers

TUN/TAP Repository Refs:
> https://pkg.go.dev/golang.zx2c4.com/wireguard/tun
> https://github.com/songgao/water

TUN/TAP Article Refs:
> https://www.kernel.org/doc/html/latest/networking/tuntap.html

TUN/TAP 原理
在Linux内核中添加了一个TUN/TAP虚拟网络设备的驱动程序和一个与之相关连的字符设备/dev/net/tun，字符设备tun作为用户空间和内核空间交换数据的接口。
当内核将数据包发送到虚拟网络设备时，数据包被保存在设备相关的一个队列中，直到用户空间程序通过打开的字符设备tun的描述符读取时，它才会被拷贝到用户空间的缓冲区中，其效果就相当于，数据包直接发送到了用户空间。通过系统调用write发送数据包时其原理与此类似。
一次read系统调用，有且只有一个数据包被传送到用户空间，并且当用户空间的缓冲区比较小时，数据包将被截断，剩余部分将永久地消失，write系统调用与read类似，每次只发送一个数据包。所以在编写此类程序的时候，请用足够大的缓冲区，直接调用系统调用read/write，避免采用C语言的带缓存的IO函数。

虚拟网卡驱动 = netdev + driver，把数据包通过 driver 发送出去。netdev 就是几个 hook 函数。

运行 TUN/TAP 设备之后，就会在 kernel space 添加一个 miscdevice (tun: /dev/net/tun, tap: /dev/tap0)。
从功能上看，tun设备驱动主要应该包括两个部分，一是虚拟网卡驱动，其实就是虚拟网卡中对skb进行封装解封装等操作；二是字符设备驱动，用于内核空间与用户空间的交互。

```bash
# 查看驱动
modinfo tun
name:           tun
filename:       (builtin)
alias:          devname:net/tun
alias:          char-major-10-200
license:        GPL
author:         (C) 1999-2004 Max Krasnyansky <maxk@qualcomm.com>
description:    Universal TUN/TAP device driver

# 加载 TUN
modprobe tun

# 查看挂载
lsmod | grep tun
ip6_udp_tunnel         16384  1 vxlan
udp_tunnel             16384  1 vxlan
```

OpenVSwitch Repository Refs:
> https://github.com/openvswitch/ovs