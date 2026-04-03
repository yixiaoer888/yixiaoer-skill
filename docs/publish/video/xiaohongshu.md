# 小红书 视频发布

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [视频发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 否 | 视频标题 | - |
| description | string | 否 | 视频描述 | - |
| declaration | number | 否 | 内容类型申明：1-虚构演绎仅供娱乐, 2-笔记含 AI 合成内容 | - |
| type | number | 否 | 创作类型：1-原创, 0-不申明 | 0 |
| location | object | 否 | 视频位置，使用 `PlatformDataItem` 结构 | - |
| scheduledTime | number | 否 | 定时发布时间戳（单位：秒） | - |
| visibleType | number | 是 | 可见类型：0-公开, 1-私密, 3-好友可见 | 0 |
| collection | object | 否 | 合集信息，使用 `Collection` 结构 | - |
| group | object | 否 | 群聊信息，使用 `Group` 结构 | - |
| bind_live_info | object | 否 | 直播预告信息，使用 `LiveInfo` 结构 | - |
| shopping_cart | object[] | 否 | 关联商品信息，使用 `ShoppingCartItem` 结构数组 | - |

## 2. 复杂对象结构

### PlatformDataItem
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| id | string | 是 | ID |
| text | string | 是 | 文本内容 |
| raw | object | 是 | 平台原始数据 |

### Collection
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| yixiaoerId | string | 是 | 合集 ID |
| yixiaoerName | string | 是 | 合集名称 |
| child | object[] | 否 | 子级合集列表 |

### Group
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| yixiaoerId | string | 是 | 群聊 ID |
| yixiaoerName | string | 是 | 群聊标题 |
| yixiaoerDesc | string | 否 | 群聊描述 |
| yixiaoerImageUrl | string | 否 | 群聊头像 URL |

### LiveInfo
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| yixiaoerId | string | 是 | 直播预告 ID |
| yixiaoerName | string | 是 | 直播预告标题 |

### ShoppingCartItem
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| yixiaoerId | string | 是 | 商品 ID |
| yixiaoerName | string | 是 | 商品名称 |
| yixiaoerDesc | string | 否 | 商品规格说明 |
| yixiaoerImageUrl | string | 否 | 商品图片 URL |
| price | number | 否 | 商品价格（单位：分） |
| earnPrice | number | 否 | 预估佣金（单位：分） |

### 数据获取途径

| 目标字段 | 对应 Action | 文档参考 |
| :--- | :--- | :--- |
| `location` | `locations` | [获取位置信息](../get-locations.md) |
| `collection` | `collections` | [获取合集列表](../../get-collections.md) |
| `group` | `groups` | [获取群聊列表](../../get-groups.md) |
| `shopping_cart` | `goods` | [获取商品列表](../../get-goods.md) |
| `declaration` | - | 请按表格中的枚举值直接填写（1 或 2） |

## 3. JSON 示例

```json
{
  "publishType": "video",
  "platforms": ["Xiaohongshu"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "XHS_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1080,
          "height": 1440,
          "duration": 30
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "小红书笔记标题",
          "description": "这是在小红书分享的一段精彩视频 #好物分享 #生活",
          "type": 1,
          "visibleType": 0,
          "declaration": 2
        }
      }
    ]
  }
}
```
