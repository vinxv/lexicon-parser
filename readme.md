# 词库解析工具

支持搜狗、百度和qq输入法的词典解析

## 使用说明

1. 下载命令行程序

2. 查看帮助文档

```bash
./lexicon-parser -h
Usage of ./lexicon-parser:
  -i string
    	input file path[required]
  -o string
    	output file path
  -t string
    	kind. eg: qq|sogou|baidu
```

1. 解析词库

输出文件每行包含三个字段，分别为`词条`,`拼音`,`词频`。

```
$ ./lexicon-parser -i '金融机构.scel' 
2022/04/25 22:34:53 metainfo
title: 金融机构
categoty:理工类
description:金融机构
sample:政府银行 城市银行 存款银行 商人银行 实业银行 

阿拉伯货币基金组织      alabohuobijijinzuzhi    16
阿拉伯经济和社会发展基金        alabojingjiheshehuifazhanjijin  15
阿姆斯特丹鹿特丹银行    amusitedanlutedanyinhang        43
巴克莱银行      bakelaiyinhang  64
巴黎国民银行    baliguominyinhang       58
巴黎荷兰金融公司        balihelanjinronggongsi  53
...
2022/04/25 22:34:53 Done. input:金融机构.scel

```
