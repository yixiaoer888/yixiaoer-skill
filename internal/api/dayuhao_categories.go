package api

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

// dayuhaoVideoCategories is the fixed catalogue used by the 蚁小二 Web 大鱼号
// video form. It is embedded into the CLI binary so category lookup never
// depends on a local yixiaoer-universal checkout or the unsupported account
// categories endpoint.
//
//go:embed dayuhao_video_categories.json
var dayuhaoVideoCategories []byte

// CategoriesForAccount resolves the account platform before querying dynamic
// categories. 大鱼号 video categories are a static Web contract, so this
// path returns the embedded catalogue and deliberately does not call the
// account categories endpoint.
func (c *Client) CategoriesForAccount(accountID, publishType string) (interface{}, error) {
	if publishType == "video" {
		accounts, err := c.Accounts("")
		if err != nil {
			return nil, err
		}
		for _, account := range accounts {
			if AccountID(account) != accountID {
				continue
			}
			platformName := stringField(account, "platformName")
			if platformName == "" {
				platformName = stringField(account, "platform")
			}
			if platformutil.CanonicalKey(platformName) == "dayuhao" {
				return dayuhaoVideoCategoriesCatalog()
			}
			break
		}
	}
	return c.Categories(accountID, publishType)
}

func dayuhaoVideoCategoriesCatalog() (interface{}, error) {
	var source []map[string]interface{}
	if err := json.Unmarshal(dayuhaoVideoCategories, &source); err != nil {
		return nil, yxerrors.Internal("大鱼号分类目录无效", err.Error()).
			WithCategory("dayuhao_category_catalog")
	}
	items, err := normalizeDayuhaoCategoryItems(source)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"dataList": items}, nil
}

func normalizeDayuhaoCategoryItems(source []map[string]interface{}) ([]interface{}, error) {
	items := make([]interface{}, 0, len(source))
	for index, raw := range source {
		id := strings.TrimSpace(fmt.Sprint(raw["yixiaoerId"]))
		name := strings.TrimSpace(fmt.Sprint(raw["yixiaoerName"]))
		if id == "" || name == "" || id == "<nil>" || name == "<nil>" {
			return nil, yxerrors.Internal("大鱼号分类目录无效", map[string]interface{}{
				"index": index,
				"item":  raw,
			}).WithCategory("dayuhao_category_catalog")
		}
		item := map[string]interface{}{
			"yixiaoerId":   id,
			"yixiaoerName": name,
			"raw":          cloneInterfaceMap(raw),
		}
		if child, ok := raw["child"].([]interface{}); ok {
			children, err := normalizeDayuhaoCategoryChildren(child)
			if err != nil {
				return nil, err
			}
			item["child"] = children
		}
		if children, ok := raw["children"].([]interface{}); ok {
			childItems, err := normalizeDayuhaoCategoryChildren(children)
			if err != nil {
				return nil, err
			}
			item["children"] = childItems
		}
		items = append(items, item)
	}
	return items, nil
}

func normalizeDayuhaoCategoryChildren(source []interface{}) ([]interface{}, error) {
	children := make([]map[string]interface{}, 0, len(source))
	for index, item := range source {
		child, ok := item.(map[string]interface{})
		if !ok {
			return nil, yxerrors.Internal("大鱼号分类目录无效", map[string]interface{}{
				"index": index,
				"item":  item,
			}).WithCategory("dayuhao_category_catalog")
		}
		children = append(children, child)
	}
	return normalizeDayuhaoCategoryItems(children)
}
