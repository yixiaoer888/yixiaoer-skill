# 📄 平台定义 索引 参数 (Platform Definitions Index)

本节定义了蚁小二 (YiXiaoEr) 支持的所有自媒体平台标识。在调用 API 的 `platforms` 数组或 `platform` 字段时，**必须**使用本列表中的“平台中文名”。

> [!TIP]
> 详细的发布接口结构请参考: [通用内容发布基础结构](./publish/index.md)

## 1. 平台标识枚举 (Platform Codes)

| 平台标识 (Code) | 平台中文名 (Name) | 说明 |
| :--- | :--- | :--- |
| `DouYin` | `抖音` | |
| `KuaiShou` | `快手` | |
| `ShiPinHao` | `视频号` | |
| `BiLiBiLi` | `哔哩哔哩` | |
| `XiaoHongShu` | `小红书` | |
| `BaiJiaHao` | `百家号` | |
| `TouTiaoHao` | `头条号` | |
| `XiGuaShiPin` | `西瓜视频` | |
| `ZhiHu` | `知乎` | |
| `QiEHao` | `企鹅号` | |
| `XinLangWeiBo` | `新浪微博` | |
| `SouHuHao` | `搜狐号` | |
| `YiDianHao` | `一点号` | |
| `DaYuHao` | `大鱼号` | |
| `WangYiHao` | `网易号` | |
| `AiQiYi` | `爱奇艺` | |
| `TengXunWeiShi` | `腾讯微视` | |
| `WeiXinGongZhongHao` | `微信公众号` | |
| `WeiXin` | `微信` | |
| `Tiktok` | `TikTok` | |
| `Youtube` | `Youtube` | |
| `Twitter` | `X` | |
| `Facebook` | `Facebook` | |
| `Instagram` | `Instagram` | |
| `Other` | `其他账号` | |
| `SouHuShiPin` | `搜狐视频` | |
| `PiPiXia` | `皮皮虾` | |
| `TengXunShiPin` | `腾讯视频` | |
| `DuoDuoShiPin` | `多多视频` | |
| `MeiPai` | `美拍` | |
| `AcFun` | `AcFun` | |
| `KuaiChuanHao` | `快传号` | |
| `XueQiuHao` | `雪球号` | |
| `CheJiaHao` | `车家号` | |
| `XiaoHongShuShangJiaHao` | `小红书商家号` | |
| `YiCheHao` | `易车号` | |
| `FengWang` | `蜂网` | |
| `DouBan` | `豆瓣` | |
| `CSDN` | `CSDN` | |
| `DeWu` | `得物` | |
| `JianShu` | `简书` | |
| `MeiYou` | `美柚` | |
| `WeiXinGongZhongHao-Open` | `微信公众号-Open` | |
| `BiLiBiLi-Open` | `哔哩哔哩-Open` | |
| `KuaiShou-Open` | `快手-Open` | |

## 2. 平台类型枚举 (PlatformType)

> [!NOTE]
> 用于区分平台的底层协议类型，Agent 在处理不同类型的账号时可能需要不同的交互策略。

| 值 (Value) | 定义 | 说明 |
| :--- | :--- | :--- |
| `0` | `Crawler` | 模拟人工/爬虫账号 |
| `1` | `OpenPlatform` | 开放平台 API 账号 |
| `2` | `Overseas` | 海外平台账号 |

---
> [!CAUTION]
> **填值准确性**：API 交互中请务必使用“中文名”。例如，调用抖音发布时，应传入 `"platform": "抖音"` 而非 `"DouYin"`。

