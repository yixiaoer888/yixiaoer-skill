# 发布图文 (Publish Post Engine)

图文动态、微头条、微博类内容的统一分布引擎。

## 场景描述 (Usage)
- "发布这条动态到我的小红书和微博。"
- "带上 3 张图片 (Key: xxx, yyy, zzz) 发布到微头条。"

## 支持平台 (Supported Platforms)

目前基座支持以下 8 个主流图文/动态平台：
- **抖音** (DouYin): 支持描述、位置、公开/隐私
- **小红书** (XiaoHongShu): 支持标题、多图描述
- **快手** (KuaiShou): 支持描述、位置
- **新浪微博** (XinLangWeiBo): 支持 140 字及以上动态
- **视频号** (WeiXinShiPinHao): 微信视频号图文
- **百家号** (BaiJiaHao): 百家号动态
- **头条号** (TouTiaoHao): 微头条
- **知乎** (ZhiHu): 知乎想法/动态

## 参数定义 (Parameters)
| 参数名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `content` | `string` | 是 | 动态文本内容 |
| `image_keys` | `string[]` | 否 | 图片素材 Key 列表 |
| `platforms` | `string[]` | 是 | 目标平台列表 |
| `account_ids` | `string[]` | 否 | 目标账号列表 |

## 脚本逻辑 (Backend)
- **脚本路径**: `../scripts/publish-post.ts`
- **核心逻辑**: 资源预处理 -> 构建图文 DTO -> 提交任务集。
