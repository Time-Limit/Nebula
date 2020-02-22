package main

import (
	"net/http"
	"io/ioutil"
	"github.com/Comdex/imgo"
	"fmt"
	"encoding/xml"
	"time"
	"os/exec"
	"sync"
	"strings"
	"errors"
	"strconv"
	"database/sql"
	_ "github.com/lib/pq"
	"regexp"
)
type rgbMatrix struct {
	h int
	w int
	matrix [][]uint8
}

type matchResult struct {
	value rune
	rate  float64
	posx  int
	posy  int
}

type workStatus int

const (
	waitingImage workStatus = iota
	waitingComment
	working
)

type workStatusManager struct {
	data map[string]workStatus
	lock sync.Mutex
}

var workStatusMan workStatusManager

func init() {
	workStatusMan.data = make(map[string]workStatus)
}

var typeAlias []string
func init() {
	typeAlias = []string {
		"无效",
		"餐饮",
		"日用",
		"交通",
		"通讯",
		"服饰",
		"美容",
		"住房",
		"医疗",
		"书籍",
		"亲友",
		"工资",
		"理财",
		"其他",
	}
}

type matchResultList []matchResult 

func (t matchResultList) Len() int {
	return len(t)
}

func (t matchResultList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t matchResultList) Less(i, j int) bool {
	return t[i].posy < t[j].posy
}

func getRGBMatrix(filePath string, clip bool) rgbMatrix {
	img, err := imgo.DecodeImage(filePath)  // 获取 图片 image.Image 对象
	if err != nil {
		fmt.Println(err)
	}

	height := imgo.GetImageHeight(img)
	width := imgo.GetImageWidth(img)

	imgMatrix := imgo.MustRead(filePath)
	matrix := make([][]uint8, height)
	for starty := 0; starty < height; starty++ {
		matrix[starty] = make([]uint8, width)
		for startx := 0; startx < width; startx++ {
			R := imgMatrix[starty][startx][0]
			G := imgMatrix[starty][startx][1]
			B := imgMatrix[starty][startx][2]
			diff := func(v uint8) uint8 {
				if v > uint8(128) {
					return uint8(255)
				}
				return uint8(0)
			}
			R = diff(uint8(R))
			G = diff(uint8(G))
			B = diff(uint8(B))
			if int(R) + int(G) + int(B) == 0 {
				matrix[starty][startx] = 0
			} else {
				matrix[starty][startx] = 255
			}
		}
	}
	h, w := len(matrix), 0
	if h > 0 {
		w = len(matrix[0])
	}

	res := rgbMatrix {
		h: h,
		w: w,
		matrix: matrix,
	}

	if clip {
		return clipRgbMatrix(res, 0, 0, res.w-1, res.h-1)
	}

	return res
}

func loadCharLib() map[rune][]rgbMatrix {
	result := make(map[rune][]rgbMatrix)
	var m rgbMatrix
	m = getRGBMatrix("./charlib/iphone7.-.strong.png", true)
	result['-'] = append(result['-'], m)
	m = getRGBMatrix("./charlib/iphone7.-.small.png", true)
	result['-'] = append(result['-'], m)
	m = getRGBMatrix("./charlib/iphone7.-.strong.2.png", true)
	result['-'] = append(result['-'], m)
	m = getRGBMatrix("./charlib/iphone7.-.small.zhifubao.png", true)
	result['-'] = append(result['-'], m)
	
	m = getRGBMatrix("./charlib/iphone7.+.strong.png", true)
	result['+'] = append(result['+'], m)
	
	m = getRGBMatrix("./charlib/iphone7.p.strong.png", true)
	result['.'] = append(result['.'], m)
	m = getRGBMatrix("./charlib/iphoneMax.p.strong.1.png", true)
	result['.'] = append(result['.'], m)
	
	m = getRGBMatrix("./charlib/iphone7.0.strong.png", true)
	result['0'] = append(result['0'], m)
	m = getRGBMatrix("./charlib/iphone7.0.small.png", true)
	result['0'] = append(result['0'], m)
	m = getRGBMatrix("./charlib/iphone7.0.small.1.png", true)
	result['0'] = append(result['0'], m)
	
	m = getRGBMatrix("./charlib/iphone7.1.strong.png", true)
	result['1'] = append(result['1'], m)
	m = getRGBMatrix("./charlib/iphone7.1.small.png", true)
	result['1'] = append(result['1'], m)
	
	m = getRGBMatrix("./charlib/iphone7.2.strong.png", true)
	result['2'] = append(result['2'], m)
	m = getRGBMatrix("./charlib/iphone7.2.small.png", true)
	result['2'] = append(result['2'], m)
	
	m = getRGBMatrix("./charlib/iphone7.3.strong.png", true)
	result['3'] = append(result['3'], m)
	m = getRGBMatrix("./charlib/iphoneMax.3.strong.1.png", true)
	result['3'] = append(result['3'], m)
	m = getRGBMatrix("./charlib/iphone7.3.small.png", true)
	result['3'] = append(result['3'], m)
	
	m = getRGBMatrix("./charlib/iphone7.4.strong.png", true)
	result['4'] = append(result['4'], m)
	m = getRGBMatrix("./charlib/iphone7.4.small.png", true)
	result['4'] = append(result['4'], m)
	
	m = getRGBMatrix("./charlib/iphone7.5.strong.png", true)
	result['5'] = append(result['5'], m)
	m = getRGBMatrix("./charlib/iphone7.5.small.png", true)
	result['5'] = append(result['5'], m)
	
	m = getRGBMatrix("./charlib/iphone7.6.strong.png", true)
	result['6'] = append(result['6'], m)
	m = getRGBMatrix("./charlib/iphone7.6.strong.1.png", true)
	result['6'] = append(result['6'], m)
	m = getRGBMatrix("./charlib/iphone7.6.small.png", true)
	result['6'] = append(result['6'], m)
	
	m = getRGBMatrix("./charlib/iphone7.7.strong.png", true)
	result['7'] = append(result['7'], m)
	m = getRGBMatrix("./charlib/iphone7.7.small.png", true)
	result['7'] = append(result['7'], m)
	
	m = getRGBMatrix("./charlib/iphone7.8.strong.png", true)
	result['8'] = append(result['8'], m)
	m = getRGBMatrix("./charlib/iphone7.8.small.png", true)
	result['8'] = append(result['8'], m)
	
	m = getRGBMatrix("./charlib/iphone7.9.strong.png", true)
	result['9'] = append(result['9'], m)
	m = getRGBMatrix("./charlib/iphone7.9.small.png", true)
	result['9'] = append(result['9'], m)
	m = getRGBMatrix("./charlib/iphone7.9.zhifubao.strong.png", true)
	result['9'] = append(result['9'], m)
	
	m = getRGBMatrix("./charlib/iphone7.s.small.png", true)
	result[':'] = append(result[':'], m)
	m = getRGBMatrix("./charlib/iphone7.s.strong.zhifubao.png", true)
	result[':'] = append(result[':'], m)

	return result
}

func needSkip(x, y int, mark [][]int) bool {
	for _, v := range mark {
		if v[0] <= x && x <= v[1] && v[2] <= y && y <= v[3] {
			return true
		}
	}
	return false
}

func extractFeature(a rgbMatrix) [40]float64 {
	var res [40]float64
	mw, mh := a.w/5 + 1, a.h/5 + 1
	getPos := func (i, j, h, w int) int {
		res := i/(h/3+1)*4 + j/(w/4+1)
		return res
	}
	for i := 0; i < a.h; i++ {
		for j := 0; j < a.w; j++ {
			if a.matrix[i][j] == 0 {
				res[i/mh] += 0
				res[j/mw + 10] += 0
				res[getPos(i, j, a.h, a.w) + 20] += 1
			}
		}
	}
	var tmpc float64
	for i := 0; i < 10; i++ {
		tmpc += res[i]
	}
	for i := 0; i < 10; i++ {
		if tmpc < 0.1 {
			res[i] = 10000.0
		} else {
			res[i] /= tmpc
		}
	}

	tmpc = 0
	for i := 10; i < 20; i++ {
		tmpc += res[i]
	}
	for i := 10; i < 20; i++ {
		if tmpc < 0.1 {
			res[i] = 10000.0
		} else {
			res[i] /= tmpc
		}
	}
	tmpc = 0
	for i := 20; i < 40; i++ {
		tmpc += res[i]
	}
	for i := 20; i < 40; i++ {
		if tmpc < 0.1 {
			res[i] = 10000.0
		} else {
			res[i] /= tmpc
		}
	}
	return res
}

func calDistance(a, b rgbMatrix) float64 {
	fa, fb := extractFeature(a), extractFeature(b)
	var c float64
	for i := 0; i < len(fa); i++ {
		c += (fa[i]-fb[i])*(fb[i]-fa[i])
	}
	return c
}

func tryMatch(matrix rgbMatrix, charMap map[rune][]rgbMatrix) (rune, float64, int, int) {
	var res rune
	var dis float64 = 0.1
	var h, w int
	for char, matrixList := range charMap {
		for _, m := range matrixList {
			tmpdis := calDistance(matrix, m)
			if tmpdis > dis || dis > 0.01 {
				res, dis, h, w = char, tmpdis, m.h, m.w
			}
			//fmt.Printf("-%c, %f, %d, %d\n", char, tmpdis, m.h, m.w)
		}
	}
	return res, dis, h, w
}

type Message struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
	MsgId        int
	PicUrl       string
}


var charMap map[rune][]rgbMatrix
func init() {
	charMap = loadCharLib()
}

func isValidChar(m rgbMatrix, L, T, R, B int, x, y int) bool {
	mark := make(map[int]bool)
	q := make([][]int, 0)
	q = append(q, []int{x, y})
	dx := []int{-1,  0,  1,  0, -1, -1,  1,  1}
	dy := []int{ 0, -1,  0,  1, -1,  1, -1,  1}
	//fmt.Println("aaa", m.h, m.w, L, R, T, B)
	for i := 0; i < len(q); i++ {
		for j := 0; j < 8; j++ {
			tx, ty := q[i][0] + dx[j], q[i][1] + dy[j]
			//fmt.Println(tx, ty)
			if ty < L || ty > R || tx < T || tx > B || m.matrix[tx][ty] != 0 || mark[tx*10000 + ty] == true {
				continue
			}
			q = append(q, []int{tx, ty})
			mark[tx*10000 + ty] = true
		}
		if len(q) >= 10 {
			return true
		}
	}
	return false;
}

func clipRgbMatrix(m rgbMatrix, L, T, R, B int) rgbMatrix {
	getT := func () int {
		for i := T; i <= B; i++ {
			for j := L; j <= R; j++ {
				if m.matrix[i][j] == 0 && isValidChar(m, L, T, R, B, i, j) {
					return i
				}
			}
		}
		return -1
	}
	getB := func () int {
		for i := B; i >= T; i-- {
			for j := L; j <= R; j++ {
				if m.matrix[i][j] == 0 && isValidChar(m, L, T, R, B, i, j) {
					return i
				}
			}
		}
		return -1
	}
	getL := func () int {
		for j := L; j <= R; j++ {
			for i := T; i <= B; i++ {
				if m.matrix[i][j] == 0 && isValidChar(m, L, T, R, B, i, j) {
					return j
				}
			}
		}
		return -1
	}
	getR := func () int {
		for j := R; j >= L; j-- {
			for i := T; i <= B; i++ {
				if m.matrix[i][j] == 0 && isValidChar(m, L, T, R, B, i, j) {
					return j
				}
			}
		}
		return -1
	}
	l := getL()
	r := getR()
	t := getT()
	b := getB()
	if l >= r || t >= b {
		return rgbMatrix{}
	}

	var tmp [][]uint8

	var blank []uint8
	for i := 0; i < r-l+1+4; i++ {
		blank = append(blank, 255)
	}

	tmp = append(tmp, blank)
	tmp = append(tmp, blank)

	for i := t; i <= b; i++ {
		var part []uint8
		part = append(part, 255)
		part = append(part, 255)
		part = append(part, m.matrix[i][l:r+1]...)
		part = append(part, 255)
		part = append(part, 255)
		tmp = append(tmp, part)
	}

	tmp = append(tmp, blank)
	tmp = append(tmp, blank)

	return rgbMatrix{
		h : b-t+1 + 4,
		w : r-l+1 + 4,
		matrix: tmp,
	}
}

func (this rgbMatrix) Output() {
	fmt.Println(this.h, this.w)
	for i := 0; i < this.h; i++ {
		for j := 0; j < this.w; j++ {
			if this.matrix[i][j] == 255 {
				fmt.Printf("*")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}
}

func extractInfo(image string, T, L, B, R int) string {
	input := getRGBMatrix(image, false)
	min := func (l, r int) int {
		if l < r  {
			return l
		}
		return r
	}
	max := func (l, r int) int {
		if l > r  {
			return l
		}
		return r
	}
	L = max(0, L)
	T = max(0, T)
	R = min(R, input.w - 1)
	B = min(B, input.h - 1)
	input = clipRgbMatrix(input, L, T, R, B)
	L, T, R, B = 0, 0, input.w-1, input.h-1
	//fmt.Println(input.h, input.w)
	//for i := T; i <= B; i++ {
	//	for j := L; j <= R; j++ {
	//		if input.matrix[i][j] == 255 {
	//			fmt.Printf("*")
	//		} else {
	//			fmt.Printf(".")
	//		}
	//	}
	//	fmt.Println()
	//}
	var str string
	for i, pre := L, 0; i <= R; i++ {
		flag := false
		for j := T; j <= B && flag == false; j++ {
			if input.matrix[j][i] == 0 {
				flag = true
			}
		}
		if flag == false {
			if pre + 10 < i {
				part := clipRgbMatrix(input, pre, T, i, B)
				char, _, _, _ := tryMatch(part, charMap)
				//fmt.Printf("%c, %v\n", char, dis)
				str = fmt.Sprintf("%s%c", str, char)
			}
			pre = i
		}
	}
	return str
}

func extractData(imageName string) (float64, time.Time, error) {
	input := getRGBMatrix(imageName, false)
	L, T, R, B := 0, 0, input.w-1, input.h-1
	rect := [][]int{}
	for i, pre := T, 0; i <= B; i++ {
		flag := false
		for j := L; j <= R && flag == false; j++ {
			if input.matrix[i][j] == 0 {
				flag = true
			}
		}
		if flag == false {
			if pre + 10 < i {
				rect = append(rect, []int{pre, i})
				//part := clipRgbMatrix(input, L, pre, R, i)
				//part.Output()
				//fmt.Println("------------------------")
			}
			pre = i
		}
	}

	var money float64 = 0.0
	var moneyFlag bool = false
	var timeVal time.Time
	var timeValFlag bool = false
	for _, r := range rect {
		L, T, R, B := 0, r[0], input.w-1, r[1]
		str := ""
		tmp := clipRgbMatrix(input, L, T, R, B)
		L, T, R, B = 0, 0, tmp.w-1, tmp.h-1
		//tmp.Output()
		for i, pre := L, 0; i <= R; i++ {
			flag := false
			for j := T; j <= B && flag == false; j++ {
				if tmp.matrix[j][i] == 0 {
					flag = true
				}
			}
			if flag == false {
				if pre + 1 < i {
					part := clipRgbMatrix(tmp, pre, T, i, B)
					char, _, _, _ := tryMatch(part, charMap)
					str = fmt.Sprintf("%s%c", str, char)
			//		part.Output()
				}
				pre = i
			}
		}
		fmt.Println(str)
		match, err := regexp.MatchString("^[-+][0-9]+\\.[0-9][0-9]$", str)
		if err != nil {
			return money, timeVal, err
		}
		if match {
			fmt.Sscanf(str, "%f", &money)
			moneyFlag = true
		}
		match, err = regexp.MatchString("^20[0-9][0-9]-[0-9][0-9]-[0-9][0-9][0-9][0-9]:[0-6][0-9]", str)
		if err != nil {
			return money, timeVal, err
		}
		if match {
			var year, month, day, hour, minute int
			fmt.Sscanf(str, "%4d-%2d-%2d%2d:%2d", &year, &month, &day, &hour, &minute)
			hour, minute = 0, 0

			stringTime := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:00", year, month, day, hour, minute)
			//fmt.Println(stringTime)
			loc, err := time.LoadLocation("Asia/Shanghai")
			if err != nil {
				return money, timeVal, err
			}

			timeVal, err = time.ParseInLocation("2006-01-02 15:04:05", stringTime, loc)
			if err != nil {
				return money, timeVal, err
			}
			timeValFlag = true
		}
		if moneyFlag && timeValFlag {
			return money, timeVal, nil
		}
	}
	fmt.Println(money, timeVal)
	return 0.0, time.Time{}, errors.New("not found")
}

func extractSpecInfo(imageName string, pos int) string {
	res := []string{}
	if pos == 0 || pos == 1{
		str1 := extractInfo(imageName, 370, 200, 500, 500)
		str1 += "   "
		str1 +=  extractInfo(imageName, 760, 200, 800, 550)
		res = append(res, "1 | " + str1)
	}

	if pos == 0 || pos == 2 {
		str2 := extractInfo(imageName, 220, 100, 320, 600)
		str2 += "   "
		str2 +=  extractInfo(imageName, 600, 380, 670, 740)
		res = append(res, "2 | " + str2)
	}

	if pos == 0  || pos == 3 {
		str3 := extractInfo(imageName, 410, 200, 550, 900)
		str3 += "   "
		str3 +=  extractInfo(imageName, 1090, 740, 1170, 1210) + ":00"
		res = append(res, "3 | " + str3)
	}

	if pos == 0  || pos == 4 {
		str3 := extractInfo(imageName, 230, 200, 310, 700)
		str3 += "   "
		str3 +=  extractInfo(imageName, 820, 450, 880, 730) + ":00"
		res = append(res, "4 | " + str3)
	}

	return strings.Join(res, "\n")
}

func downloadImage(imageName, picUrl string) error {
	cmd := exec.Command("/bin/bash", "-c", "wget -O " + imageName + " \"" + picUrl + "\"")
	_, err := cmd.Output()
	return err
}

func imageHandler(msg Message) (string, error) {
	err := downloadImage(msg.FromUserName, msg.PicUrl)
	if err != nil {
		return "", err
	}
	//res := extractSpecInfo(msg.FromUserName, 0)
	money, timeVal, err := extractData(msg.FromUserName)
	if err != nil {
		return "", err
	}
	var typeStr string
	for k, v := range typeAlias {
		if k == 0 {
			continue
		}
		typeStr = fmt.Sprintf("%s%02d: %v\n", typeStr, k, v)
	}
	return typeStr + fmt.Sprintf("	%v %v", money, timeVal), nil
}

func commentHandler(msg Message) (string, error) {
	resp := msg.Content
	if len(resp) < 1 {
		return "", errors.New("empty response")
	}
	if resp[0] == '0' {
		return "", errors.New("drop this transcation")
	}

	numStr := strings.Split(resp, " ")
	if len(numStr) < 1 {
		return "", errors.New("invalid response : \"" + resp + "\"")
	}

	typeId, err := strconv.Atoi(numStr[0])
	if err != nil {
		return "", errors.New("invalid response: get type failed, " + err.Error())
	}
	if typeId < 0 || typeId >= len(typeAlias) {
		return "", errors.New("invalid response: invalid data id, " + numStr[1])
	}
	money, timeVal, err := extractData(msg.FromUserName)
	if err != nil {
		return "", err
	}
	info := typeAlias[typeId]
	comment := "无"
	if len(numStr) >= 2 {
		comment = strings.Join(numStr[1:], " ")
	}

	_, err = db.Exec(fmt.Sprintf("INSERT INTO bill (money, date, type, comment) VALUES(%d, '%v', '%s', '%s')", int(money*100), timeVal.Format("2006-01-02 15:04:05"), info, comment))
	if err != nil {
		return "", err
	}

	return "store: " + fmt.Sprintf("%f %v", money, timeVal) + "	" + info, nil
}

func (this *workStatusManager) getContent(msg Message) string {
	this.lock.Lock()
	var curStatus workStatus
	if _, ok := this.data[msg.FromUserName]; !ok {
		this.data[msg.FromUserName] = waitingImage
	}
	curStatus = this.data[msg.FromUserName]
	if curStatus == working {
		this.lock.Unlock()
		return "working"
	}
	this.data[msg.FromUserName] = working
	this.lock.Unlock()

	var content string
	var err error
	defer func () {
		fmt.Println(content, err, curStatus)
		var nextStatus workStatus
		if err != nil {
			this.lock.Lock()
			this.data[msg.FromUserName] = waitingImage
			nextStatus = waitingImage
			this.lock.Unlock()
		} else {
			this.lock.Lock()
			if curStatus == waitingImage {
				this.data[msg.FromUserName] = waitingComment
				nextStatus = waitingComment
			} else if curStatus == waitingComment {
				this.data[msg.FromUserName] = waitingImage
				nextStatus = waitingImage
			}
			this.lock.Unlock()
		}
		if nextStatus == waitingImage {
			if err == nil {
				cmd := exec.Command("/bin/bash", "-c", "mv " + msg.FromUserName + " " + fmt.Sprintf("./backup/%v.png", time.Now().UnixNano()))
				output, err := cmd.Output()
				if err != nil {
					fmt.Println(output, err)
				}
			} else {
				cmd := exec.Command("/bin/bash", "-c", "mv " + msg.FromUserName + " " + fmt.Sprintf("./debug/%v.png", time.Now().UnixNano()))
				output, err := cmd.Output()
				if err != nil {
					fmt.Println(output, err)
				}
			}
		}
	}()
	switch curStatus {
		case waitingImage: {
			content, err = imageHandler(msg)
		}
		case waitingComment: {
			content, err = commentHandler(msg)
		}
		default: {
			panic("invalid curStatus")
		}
	}
	if err != nil {
		return err.Error()
	}
	return content
}

func waitingMessage(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var msg Message
	err = xml.Unmarshal(body, &msg)
	if err != nil {
		resp := msg
		resp.FromUserName, resp.ToUserName = resp.ToUserName, resp.FromUserName
		resp.Content = "unmarshal failed"
		resp.MsgType = "text"
		resp.CreateTime = time.Duration(time.Now().Unix())
		c, _ := xml.MarshalIndent(resp, "	", "	")
		fmt.Fprintf(w, string(c))
		return
	}

	content := workStatusMan.getContent(msg)
	resp := msg
	resp.FromUserName, resp.ToUserName = resp.ToUserName, resp.FromUserName
	resp.Content = content
	resp.MsgType = "text"
	resp.CreateTime = time.Duration(time.Now().Unix())
	c, _ := xml.MarshalIndent(resp, "	", "	")
	fmt.Fprintf(w, string(c))
	return
}

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("postgres", "port=5432 user=bill password=bill dbname=bill sslmode=disable")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	db.Exec(`
		CREATE TABLE bill(
		money integer NOT NULL,
		date TIMESTAMPTZ NOT NULL,
		type text NOT NULL,
		comment text NOT NULL)`)
}

func main(){
	//tmp := getRGBMatrix("./backup/1580807029539510230.png", false)
	//tmp.Output()
	//m, t, e := extractData("./backup/1580807029539510230.png")
	//fmt.Println(m, t, e)
	//return
	InitDB()
	http.HandleFunc("/hellonebula", waitingMessage)
	http.ListenAndServe(":80", nil)
}
