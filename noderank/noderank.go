package noderank

import (
	"fmt"
	"encoding/json"
	url2 "net/url"
	"math"
	"sort"
	"strconv"
	"strings"
	"log"
	"github.com/awalterschulze/gographviz"
)

type NodeRank interface{
	AddAttestationInfo()();
	GetRank()();
	CaculateRank()();
	FloatRound()();
	PrintHCGraph()();
}

type message struct {
	TeeNum     int64    `json:"tee_num"`
	TeeContent []teectx `json:"tee_content"`
}

type GetRankRequest struct {
	Blocks   string `json:"blocks"`
	Duration int    `json:"duration"`
}

type teectx struct {
	Attester string  `json:"attester"`
	Attestee string  `json:"attestee"`
	Score    float64 `json:"score"`
	Time     string  `json:"time,omitempty"`
	Nonce    int64   `json:"nonce,omitempty"`
}

type teescore struct {
	Attestee string  `json:"attestee"`
	Score    float64 `json:"score"`
}

// TeeSoreSlice ...
type TeeSoreSlice []teescore

func GetRank(request *GetRankRequest, period int64, numRank int64)([]teescore, []teectx, error) {
	var msgArr []string;
	err := json.Unmarshal([]byte(request.Blocks), &msgArr);
	if err != nil {
		fmt.Println("unmarshal string array error, result.Blocks = ", request.Blocks);
		return nil, nil, err
	}

	graph := NewGraph()

	cm := make(map[string]teectx)

	rArr0 := []teectx{}

	for _, m2 := range msgArr {
		msgT, err := url2.QueryUnescape(m2)
		if err != nil {
			fmt.Println("QueryUnescape error, m2 = ", m2)
			return nil, nil, err
		}
		var msg message
		err = json.Unmarshal([]byte(msgT), &msg)
		if err != nil {
			fmt.Println("unmarshal message error, msgT = ", msgT)
			return nil, nil, err
		}

		rArr := msg.TeeContent

		for _, r := range rArr {
			if math.IsNaN(r.Score) || math.IsInf(r.Score, 0) {
				fmt.Println("un invalid rank param. score : ", r.Score)
			} else {
				if r.Score == 0 {
					fmt.Println("un invalid rank param. score is zero.")
				}
				graph.Link(r.Attester, r.Attestee, r.Score)
				cm[r.Attestee] = teectx{r.Attester, r.Attestee, r.Score, "", 0}
				rArr0 = append(rArr0, r)
			}
		}
	}
	var rst []teescore
	var teectxslice []teectx

	graph.Rank(0.85, 0.0001, func(attestee string, score float64) {
		tee := teescore{attestee, floatRound(score, 8)}
		rst = append(rst, tee)
	})
	sort.Sort(TeeSoreSlice(rst)) // 把计算结果按得分高低排序
	if len(rst) < 1 {
		return nil, nil, nil
	}

	endIdx := int64(len(rst))
	if endIdx > numRank {
		endIdx = numRank
	}

	rst = rst[0:endIdx] // 返回得分较大的 endIdx 个元素
	// 以结果的Attestee作为key
	scoreMap := make(map[string]float64)
	for _, r := range rst {
		scoreMap[r.Attestee] = r.Score
	}

	// 遍历数组，获取前n个排名的被实节点对应的证实交易。
	for _, r := range rArr0 {
		if scoreMap[r.Attestee] != 0 {
			teectxslice = append(teectxslice, r)
		}
	}

	return rst, teectxslice, nil
}

func AddAttestationInfo(addr1 string, url string, info []string) error {
	raw := new(teectx)
	raw.Attester = info[0]
	raw.Attestee = info[1]
	raw.Nonce, _ = strconv.ParseInt(info[3], 10, 64)
	raw.Time = info[4]
	score, err := strconv.ParseFloat(info[2], 64)
	if err != nil {
		return err
	}
	raw.Score = score
	// m := new(message)
	// m.TeeNum = 1
	// m.TeeContent = []teectx{*raw}
	// ms, err := json.Marshal(m)
	if err != nil {
		return err
	}


	//TODO 存数据

	// d := time.Now()
	// ds := d.Format("20060102")
	// data := "{\"command\":\"storeMessage\",\"address\":" + addr1 + ",\"message\":" + url2.QueryEscape(string(ms[:])) + ",\"tag\":\"" + ds + "TEE\"}"
	// _, err = doPost(url, []byte(data))
	// if err != nil {
	// 	return err
	// }
	return nil
}


// PrintHCGraph 辅助方法，用来打印结果
func PrintHCGraph(request *GetRankRequest, period string) error {
	// fmt.Println(request.Duration)
	// fmt.Println(request.Blocks)

	var msgArr []string
	err := json.Unmarshal([]byte(request.Blocks), &msgArr)
	if err != nil {
		log.Panic(err)
	}

	graph := gographviz.NewGraph()

	for _, m2 := range msgArr {
		msgT, err := url2.QueryUnescape(m2)
		if err != nil {
			log.Panicln(err)
		}
		fmt.Println("message : " + msgT)
		var msg message
		err = json.Unmarshal([]byte(msgT), &msg)
		if err != nil {
			log.Panic(err)
		}

		rArr := msg.TeeContent
		for _, r := range rArr {
			//score := strconv.FormatUint(uint64(r.Score), 10) // TODO add this score info
			graph.AddNode("G", r.Attestee, nil)
			graph.AddNode("G", r.Attester, nil)
			graph.AddEdge(r.Attester, r.Attestee, true, nil)
			if err != nil {
				log.Panic(err)
			}
		}
	}

	output := graph.String()
	fmt.Println(output)
	return nil
}


func floatRound(f float64, n int) float64 {
	format := "%." + strconv.Itoa(n) + "f"
	res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return res
}

// Len ...
func (t TeeSoreSlice) Len() int {
	return len(t)
}

// Swap ...
func (t TeeSoreSlice) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// Less ...
func (t TeeSoreSlice) Less(i, j int) bool {
	if t[i].Score != t[j].Score {
		return t[i].Score > t[j].Score
	}
	return strings.Compare(t[j].Attestee, t[i].Attestee) > 0
}
