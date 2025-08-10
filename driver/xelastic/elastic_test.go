package xelastic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/HaleyLeoZhang/go-component/driver/xlog"
	v7 "github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"
)

const (
	// 测试用的 es 对应 mapping ---- 这个 Mapping 需要有 ik 分词器
	// --- 安装 ik 分词插件 elasticsearch-plugin install https://github.com/medcl/elasticsearch-analysis-ik/releases/download/v7.5.1/elasticsearch-analysis-ik-7.5.1.zip
	// --- ElasticSearch7.x 字段数据类型  https://www.cnblogs.com/xiaofeng-fu/p/13596725.html
	// --- ElasticSearch 八股文 - 学习阶段 https://zhuanlan.zhihu.com/p/102500311
	TestMapping = `
{
	"settings":{
		"number_of_shards":3,
		"number_of_replicas":2
	},
	"mappings":{
	    "properties":{
	        "id":{
	            "type":"integer"
	        },
	        "title":{
	            "type":"text",
	            "analyzer":"ik_smart",
	            "search_analyzer":"ik_smart"
	        },
	        "describe":{
	            "type":"text",
	            "analyzer":"ik_smart",
	            "search_analyzer":"ik_smart"
	        },
	        "category":{
	            "type":"keyword"
	        }
	    }
	}
}
`
)

// 测试用的 ES 存储结构
type BlogEsModel struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Describe string `json:"describe"`
	Category string `json:"category"`
}

// 测试用的 ES 对应 index 名字
func (BlogEsModel) GetIndex() string {
	return "blog_front_search_v1"
}

// 方便查询/删除 用格式化后的ID
func (b *BlogEsModel) GetIdString() string {
	return fmt.Sprintf("%v", b.Id)
}

type TestConfig struct {
	Es *Config `yaml:"elastic"`
}

var (
	cfg = &TestConfig{}
	ctx = context.Background()

	index = BlogEsModel{}.GetIndex()

	instance *v7.Client

	item = &BlogEsModel{
		Id:       74,
		Title:    "Go pprof性能调优",
		Describe: "在计算机性能调试领域里，profiling 是指对应用程序的画像，画像就是应用程序使用 CPU 和内存的情况，Go语言是一个对性能特别看重的语言，因此语言中自带了 profiling 的库，这篇文章就要讲解怎么在 golang 中做 profiling",
		Category: "Golang",
	}
)

func TestRun(t *testing.T) {
	err := loadConfig()
	if err != nil {
		t.Fatalf("Err(%+v)", err)
	}
	TestIniInstance(t) // 初始化ES相关测试数据
	TestIniIndex(t)    // 初始化测试用的 index
	// CURD 测试
	TestDoUpsert(t)
	<-time.After(2 * time.Second) // --- 理论上，允许 2 以内生成索引的延迟
	TestDoSearch(t)
	//TestDoDelete(t)
	// -----
}

func loadConfig() (err error) {
	var yamlFile string
	yamlFile, err = filepath.Abs("./app.yml")
	if err != nil {
		return
	}
	yamlRead, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(yamlRead, cfg)
	if err != nil {
		return
	}
	return
}

func TestIniInstance(t *testing.T) {
	var (
		err error
	)
	//instance, err = NewV7(cfg.Es)
	instance, err = NewV7(cfg.Es, v7.SetTraceLog(new(TraceLog))) // 同时打印请求日志
	if err != nil {
		msg := fmt.Sprintf("NewClient err(%+v)", err)
		t.Fatal(msg)
		return
	}

}

func TestIniIndex(t *testing.T) {
	t.Log("--------------正在 初始化Mapping--------------")
	b, err := instance.IndexExists(index).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !b {
		createIndex, err := instance.CreateIndex(index).Body(TestMapping).Do(ctx)
		if err != nil {
			msg := fmt.Sprintf("NewClient err(%+v)", err)
			t.Fatal(msg)
		}
		if createIndex == nil {
			msg := fmt.Sprintf("Expected result to be != nil")
			t.Fatal(msg)
		}
	}
	t.Log("-------------- 初始化Mapping Done --------------")
}

func TestDoUpsert(t *testing.T) {
	t.Log("--------------正在 Upsert--------------")
	_, err := instance.Update().Index(index).Id(item.GetIdString()).Doc(item).DocAsUpsert(true).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("-------------- Upsert Done --------------")
}

func TestDoSearch(t *testing.T) {
	t.Log("--------------正在 Search--------------")
	var (
		list  []*BlogEsModel
		total int64
		err   error
	)
	_, err = instance.Update().Index(index).Id(item.GetIdString()).Doc(item).DocAsUpsert(true).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// 拉取数据
	// - 设置需要的字段
	fields := []string{"id", "title", "describe", "category"}
	search := v7.NewBoolQuery()
	//search.Must(v7.NewMatchQuery("describe", "画像")) // 中文分词直接用

	shouldConditionTwoShape1 := v7.NewRangeQuery("id").Gte(50).Lte(500)                                  // 指代Id必须 >= 50 且 < 500
	shouldConditionTwoShape2 := v7.NewMatchQuery("title", "沐临风")                                         // 分词查询
	shouldConditionTwoShape3 := v7.NewTermsQuery("comment", v7.NewMatchQuery("comment.username", "沐临风")) // 精确查询
	mustShould := v7.NewBoolQuery().Should(shouldConditionTwoShape1, shouldConditionTwoShape2, shouldConditionTwoShape3)
	search.Must(mustShould)
	/**
	// 需要避坑的： should 和 must 在同一层级的时候 must 会生效 但是 should 不会

	// 场景：满足其中之一条件
	search := v7.NewBoolQuery()
	search.Must(v7.NewTermsQuery("id", 75, 5641))

	// 场景：必须同时满足多个条件
	search := v7.NewBoolQuery()
	search.Must(v7.NewTermsQuery("id", 75, 5641))
	//search.Must(v7.NewTermsQuery("title", "goland")) // 精确匹配
	search.Must(v7.NewMatchQuery("title", "goland")) // 分词匹配


	// 场景：只要满足其中一个条件即可
	search := v7.NewBoolQuery()
	search.Should(v7.NewTermsQuery("id", 75, 5641))
	search.Should(v7.NewTermsQuery("title", "goland"))

	// 场景：范围区间搜素
	search := v7.NewBoolQuery()
	shouldConditionTwoShape1 := v7.NewRangeQuery("id").Gte(50) // 指代Id必须 >= 50
	search.Must(shouldConditionTwoShape1)


	// 场景：必须满足至少一个条件 -- 时间查询
	示例 mapping 结构
	{
	    "mappings":{
	        "properties":{
	            "online_time":{
	                "type":"date",
	                "format":"yyyy-MM-dd HH:mm:ss"
	            }
	        }
	    }
	}

	search := v7.NewBoolQuery()
	shouldConditionTwoShape1 := v7.NewRangeQuery("online_time").Gte("2022-03-01 18:35:05") // 指代时间必须大于等于 2022-03-01 18:35:05
	shouldConditionTwoShape2 := v7.NewRangeQuery("online_time").Lt("2022-04-01 18:35:05") //  指代时间必须小于 2022-04-01 18:35:05
	mustShould := v7.NewBoolQuery().Should(shouldConditionTwoShape1, shouldConditionTwoShape2)
	search.Must(mustShould)



	// 场景: nested 结构搜索 - 假设 目前文章评论信息是 comment 参数，其下有个 评论人 参数名是 username
	示例 mapping 结构
	{
	    "mappings":{
	        "properties":{
	            "comment":{
	                "type":"nested",
	                "properties":{
	                    "id":{
	                        "type":"integer"
	                    },
	                    "username":{
	                        "type":"text",
	                        "analyzer":"ik_smart",
	                        "search_analyzer":"ik_smart"
	                    },
	                    "content":{
	                        "type":"keyword"
	                    }
	                }
	            }
	        }
	    }
	}
	// - ES Nested结构 嵌套查询 https://czjxy881.github.io/elasticsearch/%E4%B8%80%E8%B5%B7%E6%9D%A5%E5%AD%A6ES-%E6%B5%85%E8%B0%88Nested%E7%BB%93%E6%9E%84/
	shouldCond1 := v7.NewNestedQuery("comment", v7.NewMatchQuery("comment.username","沐临风")) // 分词
	shouldCond2 := v7.NewNestedQuery("comment", v7.NewTermsQuery("comment.username", "沐临风")) // 完全匹配
	search.Must(v7.NewBoolQuery().Should(shouldCond1, shouldCond2))

	*/
	// - 计算分页
	page := 1
	limit := 10
	offset := (page - 1) * limit

	result, err := instance.Search().Index(index).
		Sort("id", false). // 依据Id
		From(offset).Size(limit). // 取数据区间
		Query(search).
		FetchSourceContext(v7.NewFetchSourceContext(true).Include(fields...)).
		// 关于 dfs 搜索方式 官方文档 https://www.elastic.co/cn/blog/understanding-query-then-fetch-vs-dfs-query-then-fetch
		SearchType("query_then_fetch"). // 前台系统用默认的就可以了。后台系统可以用dfs，但是很少有场景需要dfs  https://blog.csdn.net/HuoqilinHeiqiji/article/details/103460430
		Do(ctx)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	// 解析数据
	list, total = handleSearchData(ctx, result)
	// ------ Print
	t.Logf("total(%v)", total)
	if total == 0 {
		return
	}
	for _, one := range list {
		bytes, _ := json.Marshal(one)
		t.Logf("one(%v)", string(bytes))
	}
	t.Log("-------------- Search Done --------------")
	return
}

// 处理返回结果
func handleSearchData(ctx context.Context, result *v7.SearchResult) (list []*BlogEsModel, totalInt64 int64) {
	list = make([]*BlogEsModel, 0)

	totalInt64 = result.TotalHits()
	if totalInt64 == 0 {
		return
	}
	for _, hit := range result.Hits.Hits {
		d := &BlogEsModel{}
		buf, _ := hit.Source.MarshalJSON()
		err := json.Unmarshal(buf, &d)
		if err != nil {
			err = errors.WithStack(err)
			xlog.Warnf(ctx, "Warning value(%+v) err(%+v)", string(buf), err)
			continue
		}
		// 固定解析数据类型返回
		list = append(list, d)
	}

	return
}

func TestDoDelete(t *testing.T) {
	t.Log("--------------正在 Delete--------------")
	_, err := instance.Delete().Index(index).Id(item.GetIdString()).Do(ctx)
	if err != nil {
		if v7.IsNotFound(err) || v7.IsConflict(err) {
			err = nil
		} else {
			xlog.Errorf(ctx, "DeleteSingleGoods Err(%+v) id(%v)", err, item.Id)
		}
		return
	}
	t.Log("-------------- Delete Done --------------")
}
