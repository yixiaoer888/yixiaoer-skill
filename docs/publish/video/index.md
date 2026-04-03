# 视频发布 (Video Publish)

> [!CAUTION]
> **阅读规范 (Reading Protocol)**:
> 本文档是 **所有平台** 视频发布的 **唯一入口** 和 **基础 DTO 定义**。
> 在查阅具体的平台文档（如 `douyin.md`）之前，你 **必须** 首先查阅本文档以理解 Payload 的根结构，否则将导致生成的 JSON 无法通过校验。

所有通过 `api.ts`（指定 `action: "publish"`）执行的视频发布任务均遵循以下数据结构。

> [!IMPORTANT]
> **发布合规性要求**:
> 所有的封面 (`cover`)、视频 (`video`) 均**必须**使用通过[资源上传接口](../../upload-resource.md)获得的资源 `key`。
> **严禁**直接填写外部网络 URL 或在该填入 Key 的地方留空。

## 1. 数据结构 (Data Structure)

接口要求传入 `CloudTaskPushRequest` 结构。

### 1.1 基础结构 (Base Structure)

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `publishType` | `string` | **是** | 固定为 `video` | - |
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
| `accountForms` | `Array` | **是** | 账号发布表单列表 | - |

### 1.3 账号表单项 (accountForms Item)

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `platformAccountId` | `string` | **是** | 蚁小二平台账号唯一 ID | - |
| `video` | `Object` | **是** | **VideoFormItem**: 视频对象 (`key`, `width`, `height`, `size`) | - |
| `cover` | `Object` | **是** | **ImageFormItem**: 主封面对象 | - |
| `contentPublishForm`| `Object` | **是** | **透传层**: `{}` | - |
| `coverKey` | `string` | **是** | 账号级封面 Key (必须与 `cover.key` 一致) | - |
| `fps` | `number` | 否 | 视频发布帧率 (海外平台使用) | - |

## 2. 发布示例 (Payload Example)

```json
{
  "action": "publish",
  "publishType": "video",
  "platforms": ["抖音"],
  "coverKey": "video_cover_key",
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "acc_vid_003",
        "video": {
          "key": "video_oss_key",
          "width": 1080,
          "height": 1920,
          "size": 52428800
        },
        "coverKey": "video_cover_key",
        "cover": {
          "key": "video_cover_key",
          "width": 1080,
          "height": 1920,
          "size": 307200
        },
        "contentPublishForm": {
          "formType": "task"
        }
      }
    ]
  }
}
```

## 3. 支持平台列表 (Support Platforms)

以下平台支持通过 `publishType: "video"` 进行发布。
contentPublishForm 中的字段需要从以下文档中获取。

| 平台名称 | 标识符 | 文档链接 |
| :--- | :--- | :--- |
| **头条号** | `Toutiaohao` | [toutiaohao.md](./toutiaohao.md) |
| **哔哩哔哩** | `Bilibili` | [bilibili.md](./bilibili.md) |
| **哔哩哔哩-Open** | `BilibiliOpen` | [bilibili-open.md](./bilibili-open.md) |
| **百家号** | `Baijiahao` | [baijiahao.md](./baijiahao.md) |
| **快手** | `Kuaishou` | [kuaishou.md](./kuaishou.md) |
| **快手-Open** | `KuaishouOpen` | [kuaishou-open.md](./kuaishou-open.md) |
| **新浪微博** | `Xinlangweibo` | [xinlangweibo.md](./xinlangweibo.md) |
| **视频号** | `Shipinghao` | [shipinghao.md](./shipinghao.md) |
| **知乎** | `Zhihu` | [zhihu.md](./zhihu.md) |
| **企鹅号** | `Qiehao` | [qiehao.md](./qiehao.md) |
| **爱奇艺** | `Aiqiyi` | [aiqiyi.md](./aiqiyi.md) |
| **网易号** | `Wangyihao` | [wangyihao.md](./wangyihao.md) |
| **一点号** | `Yidianhao` | [yidianhao.md](./yidianhao.md) |
| **搜狐号** | `Souhuhao` | [souhuhao.md](./souhuhao.md) |
| **腾讯微视** | `Weishi` | [weishi.md](./weishi.md) |
| **搜狐视频** | `Souhushipin` | [souhushipin.md](./souhushipin.md) |
| **皮皮虾** | `Pipixia` | [pipixia.md](./pipixia.md) |
| **腾讯视频** | `TencentVideo` | [tengxunshipin.md](./tengxunshipin.md) |
| **多多视频** | `DuoduoVideo` | [duoduoshipin.md](./duoduoshipin.md) |
| **美拍** | `Meipai` | [meipai.md](./meipai.md) |
| **AcFun** | `AcFun` | [acfun.md](./acfun.md) |
| **大鱼号** | `Dayuhao` | [dayuhao.md](./dayuhao.md) |
| **车家号** | `Chejiahao` | [chejiahao.md](./chejiahao.md) |
| **蜂网** | `Fengwang` | [fengwang.md](./fengwang.md) |
| **得物** | `Dewu` | [dewu.md](./dewu.md) |
| **美柚** | `Meiyou` | [meiyou.md](./meiyou.md) |
| **小红书商家号** | `Xiaohongshushop` | [xiaohongshushop.md](./xiaohongshushop.md) |
| **小红书** | `Xiaohongshu` | [xiaohongshu.md](./xiaohongshu.md) |
| **抖音** | `Douyin` | [douyin.md](./douyin.md) |
| **易车号** | `Yichehao` | [yichehao.md](./yichehao.md) |

> [!TIP]
> 持续增加中... 请参考后端 DTO `*VideoForm` 扩展新平台。

