# 小红书图文发布参数 (Xiaohongshu Image-Text)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [图文发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。

## 触发场景 (Trigger)
- **意图辨析**：用户指定在“小红书”平台发布图文笔记，且需要配置如“标记话题”、“地点挂载”、“多级分类挂载”或“定时发布”等功能时触发。
- **典型提示词**：
  - “帮我发一篇小红书笔记，带上 #穿搭 话题”
  - “小红书发布，地点选在上海东方明珠”
  - “把这两张图存为小红书草稿，设置仅好友可见”
  - “查询小红书的分类并设置”

## 执行逻辑 (Logic Flow)
1. **内容识别**：识别笔记标题、正文及内嵌话题（小红书正文支持 HTML 话题标签）。
2. **辅助检索**：
   - 话题：调用 `challenges` 获取标准话题 DTO。
   - 地点：调用 `locations` 获取 POI 数据。
   - 音乐：若需要，调用 `music` 获取。
3. **参数装配**：将处理后的字段填入 `accountForms[i].contentPublishForm`。
4. **状态执行**：调用 `node scripts/api.ts`。

## 1. contentPublishForm 参数 definition

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | 否 | 标题 (笔记标题，最多 20 字) | - |
| `description` | `string` | **是** | 笔记描述，支持 HTML (`<p>`, `<topic>`)。最多 1000 字符。 | - |
| `images` | `Array` | **是** | 图片数组 (`OldImage[]`) | - |
| `location` | `Object` | 否 | 位置对象 (`PlatformDataItem`) | - |
| `music` | `Object` | 否 | 音乐对象 (`MusicItem`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (13 位 Unix 时间戳，单位: 毫秒) | - |
| `collection` | `Object` | 否 | 合集信息，使用 `Collection` 结构 | - |
| `visibleType` | `number` | **是** | 可见类型: 0-公开, 1-私密, 3-好友可见 | 0 |

## 2. 复杂对象结构说明

### OldImage
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `width` | `number` | **是** | 图片宽度 |
| `height` | `number` | **是** | 图片高度 |
| `size` | `number` | **是** | 文件大小 (Bytes) |
| `key` | `string` | **是** | 资源 Key (通过上传接口获取) |
| `format` | `string` | **是** | 文件格式 (e.g., `jpg`, `png`) |

### PlatformDataItem (基础结构)
包含 `yixiaoerId`, `yixiaoerName`, `raw`。

### MusicItem (音乐)
包含 `yixiaoerId`, `yixiaoerName`, `duration`, `playUrl`, `raw` 等。

## 3. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "imageText",
  "platforms": ["小红书"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "XHS_ACC_ID",
        "contentPublishForm": {
          "formType": "task",
          "title": "今日穿搭",
          "description": "<p>今日穿搭分享 <topic text='穿搭' raw='{\"id\":\"xxx\",\"name\":\"穿搭\"}'>#穿搭</topic></p>",
          "images": [
            { "key": "img_xhs_01", "size": 1024, "width": 1080, "height": 1440, "format": "jpg" }
          ],
          "visibleType": 0
        }
      }
    ]
  }
}
```

## 相关接口

| 目标字段 | 对应 Action | 文档参考 |
| :--- | :--- | :--- |
| `location` | `locations` | [获取位置信息](../../get-locations.md) |
| `collection` | `collections` | [获取合集列表](../../get-collections.md) |
| `tags/topic`| `challenges` | [获取话题/挑战](../../get-challenges.md) |
| `images.key`| `upload` | [资源上传](../../upload-resource.md) |
