# WiFi万能钥匙文章发布 (Publish WiFiWanNengYaoShi Article)

该指令用于通过文章引擎向 WiFi万能钥匙（连尚）分发内容。

## DTO 溯源 (Knowledge from WifiWanNengArticleForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/wifiwanneng.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `article`** | 业务模态识别 |
| `--title` | string | 是 | 文章标题 | 不可为空 |
| `--content` | string | 是 | 文章内容 (HTML) | 不可为空 |
| `--cover_url` | string | 是 | 封面图 | 引擎自动上传并映射为 `covers` 数组 |
| `--category` | array | 否 | 文章分类 | 对象列表 |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=article \
  --title="WiFi万能钥匙平台发布测试" \
  --content="<p>文章内容展示...</p>" \
  --platforms="WiFi万能钥匙" \
  --account_ids="wifi_acc_001" \
  --cover_url="https://example.com/cover.jpg"
```
