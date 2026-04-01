# 文章发布 (Article Publish)

所有通过 `api.ts`（指定 `action: "publish"`）执行的文章发布任务均遵循以下数据结构。

> [!IMPORTANT]
> **发布合规性要求**:
> 所有的封面 (`cover`) 均**必须**使用通过[资源上传接口](../../upload-resource.md)获得的资源 `key`。
> **严禁**直接填写外部网络 URL 或在该填入 Key 的地方留空。

## 1. 数据结构 (Data Structure)

接口要求传入 `CloudTaskPushRequest` 结构。

### 1.1 基础结构 (Base Structure)

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `publishType` | `string` | **是** | 固定为 `article` | - |
| `platforms` | `string[]` | **是** | 目标平台枚举数组，详见下方平台列表 | - |
| `coverKey` | `string` | **是** | 任务封面资源 Key | - |
| `publishArgs` | `Object` | **是** | 发布参数核心容器 | - |
| `taskSetId` | `string` | 否 | 任务集唯一标识 (草稿发布时必填) | - |
| `desc` | `string` | 否 | 任务描述/摘要 | - |
| `publishChannel` | `string` | 否 | `cloud` (云端) 或 `local` (本机) | `local` |
| `clientId` | `string` | 否 | 客户端连接 ID (`local` 发布时必填) | - |
| `isDraft` | `boolean` | 否 | 是否仅保存为草稿 | `false` |

### 1.2 发布参数 (publishArgs)

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `content` | `string` | **是** | **文章正文**: HTML 格式字符串 | - |
| `accountForms` | `Array` | **是** | 账号发布表单列表 | - |

### 1.3 账号表单项 (accountForms Item)

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `platformAccountId` | `string` | **是** | 蚁小二平台账号唯一 ID | - |
| `cover` | `Object` | **是** | **ImageFormItem**: 主封面对象 (`key`, `width`, `height`, `size`) | - |
| `contentPublishForm`| `Object` | **是** | **透传层**: `{}` | - |
| `coverKey` | `string` | 否 | 账号级封面 Key (通常与 `cover.key` 一致) | - |

## 2. 发布示例 (Payload Example)

```json
{
  "action": "publish",
  "publishType": "article",
  "platforms": ["微信公众号"],
  "publishArgs": {
    "content": "<h1>演示文章标题</h1><p>这是一个演示文章的正文内容，支持 HTML 格式。</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_art_001",
        "cover": {
          "key": "article_cover_key",
          "width": 900,
          "height": 500,
          "size": 150000
        },
        "coverKey": "article_cover_key",
        "contentPublishForm": {
          "formType": "task"
        }
      }
    ]
  }
}
```

## 3. 支持平台列表 (Support Platforms)

以下平台支持通过 `publishType: "article"` 进行发布。
contentPublishForm 中的字段需要从以下文档中获取。

| 平台名称 | 标识符 | 文档链接 |
| :--- | :--- | :--- |
| **抖音** | `抖音`, `DouYin` | [douyin.md](./douyin.md) |
| **头条号** | `头条号`, `TouTiaoHao` | [toutiaohao.md](./toutiaohao.md) |
| **百家号** | `百家号`, `BaiJiaHao` | [baijiahao.md](./baijiahao.md) |
| **企鹅号** | `企鹅号`, `QiEHao` | [qiehao.md](./qiehao.md) |
| **搜狐号** | `搜狐号`, `SouHuHao` | [souhuhao.md](./souhuhao.md) |
| **一点号** | `一点号`, `YiDianHao` | [yidianhao.md](./yidianhao.md) |
| **大鱼号** | `大鱼号`, `DaYuHao` | [dayuhao.md](./dayuhao.md) |
| **网易号** | `网易号`, `WangYiHao` | [wangyihao.md](./wangyihao.md) |
| **知乎** | `知乎`, `ZhiHu` | [zhihu.md](./zhihu.md) |
| **爱奇艺** | `爱奇艺`, `AiQiYi` | [aiqiyi.md](./aiqiyi.md) |
| **新浪微博** | `新浪微博`, `XinLangWeiBo` | [xinlangweibo.md](./xinlangweibo.md) |
| **哔哩哔哩** | `哔哩哔哩`, `BiLiBiLi` | [bilibili.md](./bilibili.md) |
| **雪球号** | `雪球号`, `XueQiuHao` | [xueqiuhao.md](./xueqiuhao.md) |
| **快传号** | `快传号`, `KuaiChuanHao` | [kuaichuanhao.md](./kuaichuanhao.md) |
| **豆瓣** | `豆瓣`, `DouBan` | [douban.md](./douban.md) |
| **CSDN** | `CSDN`, `CSDN` | [csdn.md](./csdn.md) |
| **车家号** | `车家号`, `Chejiahao` | [chejiahao.md](./chejiahao.md) |
| **简书** | `简书`, `JianShu` | [jianshu.md](./jianshu.md) |
| **WiFi万能钥匙** | `WiFi万能钥匙`, `WifiWanNeng` | [wifiwanneng.md](./wifiwanneng.md) |
| **AcFun** | `AcFun`, `AcFun` | [acfun.md](./acfun.md) |
| **易车号** | `易车号`, `YiCheHao` | [yichehao.md](./yichehao.md) |
| **微信公众号** | `微信公众号`, `WeiXinGongZhongHao` | [weixingongzhonghao.md](./weixingongzhonghao.md) |

> [!TIP]
> 持续增加中... 请参考后端 DTO `*ArticleForm` 扩展新平台。

