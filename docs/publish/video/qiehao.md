# 企鹅号 视频发布

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [视频发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定为 `task` | `task` |
| `title` | `string` | **是** | 视频标题 | - |
| `description` | `string` | **是** | 视频描述 | - |
| `tags` | `string[]` | **是** | 视频标签 | - |
| `category` | `object[]` | **是** | 视频分类，使用 `CascadingPlatformDataItem[]` 结构 | - |
| `declaration` | `number` | 否 | 创作申明：0-暂不申明, 1-该内容由AI生成, 2-个人观点仅供参考, 3-剧情演绎仅供娱乐, 4-取材网络谨慎甄别, 5-旧闻 | 0 |
| `pubType` | `number` | **是** | 发布类型: 0-草稿, 1-直接发布 | 1 |
| `scheduledTime` | `number` | 否 | 定时发布时间戳（单位：秒） | - |

## 2. 复杂对象结构

### 3.1 CascadingPlatformDataItem (级联分类对象)

> [!IMPORTANT]
> **规则 (Rule)**:
> 所有的 `raw` 数据必须透传通过分类接口获取的原始对象。

| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | 蚁小二分类 ID |
| `yixiaoerName` | `string` | **是** | 蚁小二分类名称 |
| `children` | `object[]` | 否 | 子级分类列表 (`CascadingPlatformDataItem[]`) |
| `raw` | `object` | **是** | 平台原始数据 (必须透传) |

## 3. JSON 示例

```json
{
  "publishType": "video",
  "platforms": ["企鹅号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "acc_qh_001",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1920,
          "height": 1080,
          "duration": 60
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "企鹅号视频标题示例",
          "description": "这是关于企鹅号内容的详细描述文字。",
          "tags": ["综合", "生活"],
          "category": [
            {
              "yixiaoerId": "cat_001",
              "yixiaoerName": "生活",
              "raw": { "id": "1", "name": "生活" }
            }
          ],
          "declaration": 0,
          "pubType": 1
        }
      }
    ]
  }
}
```

## 相关接口

| 目标数据 | 对应 Action | 相关文档 |
| :--- | :--- | :--- |
| `video.key` | `upload` | [资源上传](../../upload-resource.md) |
| `category` | `category` | [获取分类](../../platform-category.md) |
