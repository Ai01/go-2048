package main

import (
	"github.com/nsf/termbox-go"
	"math/rand"
	"fmt"
	"time"
)

var Score int
var step int

// 输出字符
func coverPrintStr(x int, y int, str string, fg, bg termbox.Attribute) error {
	xx := x

	for n, c := range str {
		if c == '\n' {
			y ++
			xx = x - n - 1
		}
		termbox.SetCell(xx+n, y, c, fg, bg)
	}
	termbox.Flush()
	return nil
}

// 游戏状态
type Status uint

const (
	Win  Status = iota
	Lose
	Add
	Max  = 2048
)

// 2048游戏中的16个格子使用4*4二维数组表示
type G2048 [4][4]int

// 检查游戏是否已经胜利，没有胜利的情况下随机将值为0的元素，随机设置为2或4
func (t *G2048) checkWinOrAdd() Status {
	// 判断4*4中是否有元素的值大于等于2048，有则获得胜利
	for _, x := range t {
		for _, y := range x {
			if y >= Max {
				return Win
			}
		}
	}

	// 开始随机设置值元素为2或4
	i := rand.Intn(len(t))
	j := rand.Intn(len(t))
	for x := 0; x < len(t); x++ {
		for y := 0; y < len(t); y++ {
			if t[i%len(t)][j%len(t)] == 0 {
				t[i%len(t)][j%len(t)] = 2 << (rand.Uint32() % 2)
				return Add
			}
			j++
		}
		i++
	}

	// 全部元素都不为0，则失败
	return Lose
}

// 初始化游戏界面
func (t G2048) initialize(ox, oy int) error {
	fg := termbox.ColorYellow
	bg := termbox.ColorBlack
	termbox.Clear(fg, bg)

	str := "Score:" + fmt.Sprint(Score)
	for n, c := range str {
		termbox.SetCell(ox+n, oy-1, c, fg, bg)
	}

	str = "ESC:exit" + "Enter:replay"
	for n, c := range str {
		termbox.SetCell(ox+n, oy-2, c, fg, bg)
	}

	str = "Play with arrow key"
	for n, c := range str {
		termbox.SetCell(ox+n, oy-3, c, fg, bg)
	}

	fg = termbox.ColorBlack
	bg = termbox.ColorGreen
	for i := 0; i <= len(t); i++ {
		for x := 0; x < 5*len(t); x++ {
			termbox.SetCell(ox+x, oy+i*2, 'x', fg, bg)
		}

		for x := 0; x <= 2*len(t); x++ {
			if x%2 == 0 {
				termbox.SetCell(ox+i*5, oy+x, '+', fg, bg)
			} else {
				termbox.SetCell(ox+i*5, oy+x, '|', fg, bg)
			}

		}
	}

	fg = termbox.ColorYellow
	bg = termbox.ColorBlack
	for i := range t {
		for j := range t[i] {
			if t[i][j] > 0 {
				str := fmt.Sprint(t[i][j])

				for n, char := range str {
					termbox.SetCell(ox+j*5+1+n, oy+i*2+1, char, fg, bg)
				}
			}
		}
	}

	return termbox.Flush()
}

// 翻转二维切片
func (t *G2048) mirrorV() {
	tn := new(G2048)
	for i, line := range t {
		for j, num := range line {
			tn[len(t)-i-1][j] = num
		}
	}

	*t = *tn
}

// 向右旋转90度
func (t *G2048) right90() {
	tn := new(G2048)
	for i, line := range t {
		for j, num := range line {
			tn[j][len(t)-i-1] = num
		}
	}

	*t = *tn
}

// 向左旋转90度
func (t *G2048) left90() {
	tn := new(G2048)
	for i, line := range t {
		for j, num := range line {
			tn[len(line)-j-1][i] = num
		}
	}
	*t = *tn
}

func (t *G2048) right180() {
	tn := new(G2048)
	for i, line := range t {
		for j, num := range line {
			tn[len(line)-i-1][len(line)-j-1] = num
		}
	}

	*t = *tn
}

// 向上移动合并
func (t *G2048) mergeUp() bool {
	tl := len(t)
	changed := false
	notfull := false

	for i := 0; i < tl; i++ {

		np := tl
		n := 0

		// 向上移动非零值，如果有零值元素，则用非零元素进行覆盖
		for x := 0; x < np; x ++ {
			if t[x][i] != 0 {
				t[n][i] = t[x][i]
				if n != x {
					changed = true // 表示数组的元素是否有变化
				}
				n++
			}
		}

		if n < tl {
			notfull = true
		}

		np = n

		// 向上合并所有相同的元素
		for x := 0; x < np-1; x++ {
			if t[x][i] == t[x+1][i] {
				t[x][i] *= 2
				t[x+1][i] = 0
				Score += t[x][i] * step
				x++
				changed = true
			}
		}

		// 合并完相同元素后，再次向上移动非零元素
		n = 0
		for x := 0; x < np; x++ {
			if t[x][i] != 0 {
				t[n][i] = t[x][i]
				n++
			}
		}
		for x := n; x < tl; x++ {
			t[x][i] = 0
		}
	}

	return changed || !notfull
}

// 向下移动合并的操作可以转换向上移动合并：
// 1. 向右旋转180度矩阵
// 2. 向上合并
// 3. 再次向右旋转180度矩阵
func (t *G2048) mergeDown() bool {
	t.right180()
	changed := t.mergeUp()
	t.right180()

	return changed
}

// 向左移动合并转换为向上移动合并
func (t *G2048) mergeLeft() bool {
	t.right90()
	changed := t.mergeUp()
	t.left90()

	return changed
}

// 向右移动合并转换为向上移动合并
func (t *G2048) mergeRight() bool {
	t.left90()
	changed := t.mergeUp()
	t.right90()

	return changed;
}

// 检查按键，作出不同的移动动作或者退出程序
func (t *G2048) mergeAndReturnKey() termbox.Key {
	var changed bool

Lable:
	changed = false
	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

	ev := <-event_queue

	switch ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyArrowUp:
			changed = t.mergeUp()
		case termbox.KeyArrowDown:
			changed = t.mergeDown()
		case termbox.KeyArrowLeft:
			changed = t.mergeLeft()
		case termbox.KeyArrowRight:
			changed = t.mergeRight()
		case termbox.KeyEsc, termbox.KeyEnter:
			changed = true
		default:
			changed = false
		}

		if !changed {
			goto Lable
		}
	case termbox.EventResize:
		x, y := termbox.Size()
		t.initialize(x/2-10, y/2-4)
		goto Lable
	case termbox.EventError:
		panic(ev.Err)
	}

	step++
	return ev.Key
}

// 重置
func (b *G2048) clear() {
	next := new(G2048)
	Score = 0
	step = 0
	*b = *next
}

// 开始游戏
func (b *G2048) Run() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	rand.Seed(time.Now().UnixNano())

A:

	b.clear()
	for {
		st := b.checkWinOrAdd()
		x, y := termbox.Size()
		b.initialize(x/2-10, y/2-4)

		switch st {
		case Win:
			str := "Win!!"
			strl := len(str)
			coverPrintStr(x/2-strl/2, y/2, str, termbox.ColorMagenta, termbox.ColorYellow)

		case Lose:
			str := "Lose!!"
			strl := len(str)
			coverPrintStr(x/2-strl/2, y/2, str, termbox.ColorBlack, termbox.ColorRed)
		case Add:
		default:
			fmt.Print("Erf")
		}

		key := b.mergeAndReturnKey()
		if key == termbox.KeyEsc {
			return
		}
		if key == termbox.KeyEnter {
			goto A
		}
	}

}

func main() {
	var game G2048
	game.Run();
}
