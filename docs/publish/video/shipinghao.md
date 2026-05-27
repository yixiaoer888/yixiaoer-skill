# 视频号视频发布参数 (WeiXin ShiPingHao Video)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [视频发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。

## 触发场景 (Trigger)
- **意图辨析**：用户指定在“微信视频号”平台分发视频内容，且需要执行如“地点标记”、“关联网店商品”、“参加活动”或“存为视频号草稿”等微信生态功能时触发。
- **典型提示词**：
  - “帮我把这个视频发到我的视频号”
  - “视频号发布，带上我在广州的位置”
  - “在视频号推文里挂载这个商品链接”
  - “参加视频号最新的创作激励活动”

## 执行逻辑 (Logic Flow)
1. **意图解析**：识别是否需要挂载商品 (Goods)、活动 (Activity) 或位置 (Location)。
2. **多维辅助检索**：
   - 位置：调用 `locations` 获取 POI。
   - 活动：调用 `activities` 获取活动 ID。
   - 商品：调用 `goods` 获取带货商品信息。
3. **参数装配**：将获取的 `raw` 结构与 `yixiaoerId` 组装进 `accountForms[i].contentPublishForm`。
4. **指令执行**：调用 `node scripts/api.ts`。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定为 `task` | `task` |
| `title` | `string` | 否 | 视频标题 (最多 80 字) | - |
| `short_title` | `string` | 否 | 视频短标题 | - |
| `description` | `string` | 否 | 视频描述，支持 HTML 格式和 `@` 好友/话题标签 | - |
| `horizontalCover` | `object` | 否 | 视频横板封面，使用 `OldCover` 结构 | - |
| `createType` | `number` | **是** | 创建类型：1-草稿，2-直接发布 | 2 |
| `pubType` | `number` | **是** | 发布类型：0-草稿，1-直接发布 | 1 |
| `location` | `object` | 否 | 视频位置，使用 `PlatformDataItem` 结构 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (13 位 Unix 时间戳，单位: 毫秒) | - |
| `shoppingCart` | `object` | 否 | 关联商品信息 (`yixiaoerId`, `yixiaoerName`, `raw`) | - |
| `collection` | `object` | 否 | 合集信息 (`yixiaoerId`, `yixiaoerName`, `raw`) | - |
| `activity` | `object` | 否 | 活动信息 (`yixiaoerId`, `yixiaoerName`, `raw`) | - |

## 2. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Shipinghao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "SPH_ACC_ID",
        "video": { "key": "v_key", "size": 1024000, "width": 1080, "height": 1920, "duration": 15 },
        "contentPublishForm": {
          "formType": "task",
          "title": "记录美好瞬间",
          "description": "<p>这是我的视频号首发 #生活 #记录</p>",
          "createType": 2,
          "pubType": 1,
          "location": {
            "yixiaoerId": "loc_001",
            "yixiaoerName": "广州市",
            "raw": {}
          }
        }
      }
    ]
  }
}
```

## 3. 复杂对象结构说明

### 3.1 OldCover
包含 `key`, `size`, `width`, `height`。

### 3.2 PlatformDataItem (位置/商品/合集/活动)
包含 `yixiaoerId`, `yixiaoerName`, `raw` (必须完整透传)。

## 相关接口

| 目标数据 | 对应 Action | 文档参考 |
| :--- | :--- | :--- |
| `location`  | `locations` | [获取位置信息](../../get-locations.md) |
| `activity`  | `activities` | [获取活动列表](../../get-publish-activities.md) |
| `shoppingCart`| `goods`   | [获取商品列表](../../get-goods.md) |
| `video.key` | `upload`    | [资源上传](../../upload-resource.md) |
