package mock

import (
	"math/rand"
	"time"
	"todolist/models"

	"github.com/brianvoe/gofakeit/v6"
)

func CreateFakeData() []models.Todo {

	gofakeit.Seed(0)

	var fakeList []models.Todo

	generator := createUint8Generator()
	for i := 0; i < 500; i++ {
		fakeTodo := models.Todo{
			ID: generator(),
			// Value:       gofakeit.RandomString([]string{"Electronics", "Furniture", "Books"}),
			Value:       randomChineseString(20),
			IsCompleted: gofakeit.Bool(),
			CreatedAt:   gofakeit.DateRange(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
			UpdatedAt:   gofakeit.DateRange(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
		}
		fakeList = append(fakeList, fakeTodo)
	}

	return fakeList
}

func randomChineseChar() rune {
	var chineseChars = []rune{
		'一', '二', '三', '四', '五', '六', '七', '八', '九', '十',
		'人', '口', '手', '日', '月', '水', '火', '木', '金', '土',
		'天', '地', '山', '川', '风', '雨', '雪', '云', '电', '雷',
		'春', '夏', '秋', '冬', '东', '南', '西', '北', '中', '上',
		'下', '左', '右', '前', '后', '里', '外', '大', '小', '多',
		'少', '好', '坏', '高', '低', '长', '短', '新', '旧', '快',
		'慢', '冷', '热', '明', '暗', '轻', '重', '远', '近', '深',
		'浅', '宽', '窄', '粗', '细', '胖', '瘦', '美', '丑', '强',
		'弱', '真', '假', '对', '错', '是', '非', '有', '无', '来',
		'去', '进', '出', '开', '关', '生', '死', '始', '终', '早',
		'晚', '先', '后', '古', '今', '中', '外', '内', '外', '正',
		'反', '黑', '白', '红', '黄', '蓝', '绿', '紫', '灰', '粉',
		'爱', '恨', '喜', '怒', '哀', '乐', '悲', '欢', '苦', '甜',
		'酸', '辣', '咸', '淡', '香', '臭', '软', '硬', '干', '湿',
		'国', '家', '人', '民', '党', '政', '军', '法', '学', '校',
		'师', '生', '工', '农', '商', '兵', '官', '吏', '君', '臣',
		'父', '母', '子', '女', '兄', '弟', '姐', '妹', '夫', '妻',
		'男', '女', '老', '少', '长', '幼', '亲', '友', '朋', '伴',
		'书', '画', '诗', '词', '文', '章', '字', '句', '段', '篇',
		'笔', '墨', '纸', '砚', '琴', '棋', '书', '画', '歌', '舞',
		'吃', '喝', '穿', '住', '行', '走', '跑', '跳', '坐', '站',
		'睡', '醒', '看', '听', '说', '笑', '哭', '唱', '读', '写',
		'想', '做', '学', '教', '问', '答', '知', '道', '会', '能',
		'可', '要', '应', '该', '得', '给', '让', '叫', '请', '谢',
		'对', '不', '起', '没', '关', '系', '再', '见', '您', '好',
	}

	return chineseChars[rand.Intn(len(chineseChars))]
}

// 生成随机中文字符串
func randomChineseString(length int) string {
	result := make([]rune, length)
	for i := range result {
		result[i] = randomChineseChar()
	}
	return string(result)
}

func createUint8Generator() func() uint {
	var counter uint = 6
	return func() uint {
		current := counter
		counter++
		return current
	}
}
