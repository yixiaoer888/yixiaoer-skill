# 多多视频 视频发布

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [视频发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


## 触发场景 (Trigger)
- **意图辨析**：用户指定在“Duoduoshipin”平台分发视频内容时触发。
- **典型提示词**：
  - “把这个视频发布到Duoduoshipin”
  - “同步视频到Duoduoshipin”

## 执行逻辑 (Logic Flow)
1. **意图确认**：确认目标平台为Duoduoshipin。
2. **参数装配**：识别并填充标题、描述等平台特定字段至 `contentPublishForm`。
3. **指令执行**：调用 `node scripts/api.ts`。


## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| description | string | 否 | 多多视频描述 | - |
| scheduledTime | number | 否 | 定时发布时间戳（13 位 Unix 时间戳，单位：毫秒） | - |
| shopping_cart | object | 否 | 关联商品信息（购物车） | - |

## 2. JSON 示例

```json
{
  "publishType": "video",
  "platforms": ["Duoduoshipin"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "PDD_VIDEO_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1080,
          "height": 1920,
          "duration": 30
        },
        "contentPublishForm": {
          "formType": "task",
          "description": "这是多多视频的商品分享推荐内容。",
          "shopping_cart": {
            "goods_id": "998877",
            "source": "pdd"
          }
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
