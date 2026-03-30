# 通用内容发布基础结构 (Universal Publish Base)

所有通过 `publish.ts` 执行的内容发布任务均遵循此通用结构。

## 1. 基础结构 (Base Structure)

接口要求传入 `CloudTaskPushRequest` 结构。

| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `publishType` | `string` | **是** | 发布类型: `article` (文章), `imageText` (图文), `video` (视频) |
| `platforms` | `string[]` | **是** | 目标平台枚举数组 (如 `["抖音", "微信公众号"]`) |
| `publishArgs` | `Object` | **是** | 发布参数核心容器 |
| `taskSetId` | `string` | 否 | 任务集唯一标识 (草稿发布时必填) |
| `desc` | `string` | 否 | 任务描述/摘要 |
| `publishChannel` | `string` | 否 | `cloud` (云端) 或 `local` (本机)，默认 `local` |
| `clientId` | `string` | 否 | 客户端连接 ID (`local` 发布时必填) |
| `isDraft` | `boolean` | 否 | 是否仅保存为草稿，默认 `false` |
| `intervalConfig` | `Object` | 否 | 间隔发布配置 (包含 `enable`, `interval`, `timeUnit` 等) |

## 2. 发布参数 (publishArgs)

| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `accountForms` | `Array` | **是** | 账号发布表单列表 |
| `content` | `string` | 条件 | **内容正文**: 文章或图文类型下为**必填** (文章为 HTML，图文为纯文本) |

## 3. 账号表单项 (accountForms Item)

| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `platformAccountId` | `string` | **是** | 蚁小二平台账号唯一 ID |
| `cover` | `Object` | **是** | **ImageFormItem**: 主封面对象 (`key`, `width`, `height`, `size`) |
| `contentPublishForm`| `Object` | **是** | **透传层**: `{formType: "task", ...}`，其他更多属性需要参考对应平台文档 |
| `video` | `Object` | 条件 | **VideoFormItem**: 视频对象 (**视频发布必填**) |
| `images` | `Array` | 条件 | **ImageFormItem[]**: 图文图片列表 (**图文发布必填**) |
| `mediaId` | `string` | 否 | 第三方库素材 ID |
| `fps` | `number` | 否 | 视频发布帧率 (海外平台使用) |

---

### [文章发布平台列表 (Article)](./article/)
各文章平台 (WeChat, Baijia, etc.) 的详细 DTO 见子目录。

### [图文发布平台列表 (Image-Text)](./image-text/)
各图文/动态平台 (Xiaohongshu, Kuaishou, etc.) 的详细 DTO 见子目录。

### [视频发布平台列表 (Video)](./video/)
各视频平台 (Douyin, Bilibili, etc.) 的详细 DTO 见子目录。

> [!TIP]
> 平台特定的 DTO 字段请查阅对应模态子目录下的各平台文档，将其放入 `contentPublishForm` 中。
