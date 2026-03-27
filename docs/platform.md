# 平台定义 (Platform Definitions)

本文件定义了系统中支持的所有平台名称、内部标识符及相关枚举值。此参数在多个接口（如查询账号、发布作品、获取分类等）中广泛使用。

## 平台名称与标识 (Platforms)

在调用接口时，`platform` 字段通常建议使用**中文名称**（如 `抖音`），或者在支持的情况下使用其**内部 ID**。

| 平台中文名 | 内部标识 (PlatformID) | 说明 |
| :--- | :--- | :--- |
| `抖音` | `DouYin` | |
| `卷手` | `KuaiShou` | |
| `视频号` | `ShiPinHao` | |
| `哔哩哔哩` | `BiLiBiLi` | 亦称 B站 |
| `小红书` | `XiaoHongShu` | |
| `百家号` | `BaiJiaHao` | 百度旗下 |
| `头条号` | `TouTiaoHao` | 字节跳动旗下 |
| `西瓜视频` | `XiGuaShiPin` | |
| `知乎` | `ZhiHu` | |
| `企鹅号` | `QiEHao` | 腾讯旗下 |
| `新浪微博` | `XinLangWeiBo` | |
| `搜狐号` | `SouHuHao` | |
| `一点号` | `YiDianHao` | |
| `大鱼号` | `DaYuHao` | |
| `网易号` | `WangYiHao` | |
| `爱奇艺` | `AiQiYi` | |
| `腾讯微视` | `TengXunWeiShi` | |
| `微信公众号` | `WeiXinGongZhongHao` | |
| `微信` | `WeiXin` | 通用微信标识 |
| `TikTok` | `Tiktok` | 海外版抖音 |
| `Youtube` | `Youtube` | 海外视频平台 |
| `X` | `Twitter` | 原 Twitter |
| `Facebook` | `Facebook` | |
| `Instagram` | `Instagram` | |
| `CSDN` | `CSDN` | 技术社区 |
| `得物` | `DeWu` | |
| `简书` | `JianShu` | |
| `豆瓣` | `DouBan` | |

> [!TIP]
> 列表未穷举所有长尾平台，可通过 `query-accounts.ts` 获取当前租户实际支持的所有平台。

## 平台类型枚举 (PlatformType)

| 值 (Value) | 定义 (Key) | 说明 |
| :--- | :--- | :--- |
| `0` | `Crawler` | **其他/爬虫账号** (通过浏览器插件或模拟登录获取授权) |
| `1` | `OpenPlatform` | **开放平台** (通过官方 API 接口进行 OAuth 授权) |
| `2` | `Overseas` | **海外平台** |

## 参数使用规范

- **单选**: 参数名通常为 `platform`，类型为 `string`，值为平台中文名或内部 ID。
- **多选**: 参数名通常为 `platforms[]` 或以逗号分隔的字符串。
- **类型过滤**: 参数名通常为 `platformType`，类型为 `number`。
