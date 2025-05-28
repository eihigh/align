package align_test

import (
	"fmt"

	"github.com/eihigh/align"
)

const (
	mid    = 0.5
	top    = 0.0
	bottom = 1.0
	left   = 0.0
	right  = 1.0
)

func Example() {
	txtNode := func(s string) *align.Node[int] {
		// dummy function to create a node
		return align.WH(10, 10)
	}

	screen := align.WH(150, 150)
	menuBounds := screen.Inset(50)
	off := align.XY(10, 0)
	txtNewGame := txtNode("New Game")
	txtContinue := txtNode("Continue").StackY(txtNewGame, 0, 1).Offset(off)
	txtExit := txtNode("Exit").StackY(txtContinue, 0, 1).Offset(off)

	menu := align.Union(txtNewGame, txtContinue, txtExit).Nest(menuBounds, 1, 1)
	fmt.Println(txtNewGame, txtContinue, txtExit, menu)

	// Output:
	// Pos: (70,70), Size: (10,10) Pos: (80,80), Size: (10,10) Pos: (90,90), Size: (10,10) Pos: (70,70), Size: (30,30)
}

func Example_form() {
	screen := align.WH(800, 600)
	formBounds := screen.Inset(24)
	rowBounds := align.WH(formBounds.Size.X, 0)
	label := align.WH(80, 20).Nest(rowBounds, 0, .5)
	input := align.WH(200, 25).Nest(rowBounds, 1, .5)
	row := align.Union(label, input).Nest(formBounds, 0, 0)
	fmt.Println("Row:", row)

	// Output:
	// Row: Pos: (24,24), Size: (752,25)
}

func ExampleWrapper() {
	stack := func(a, b *align.Node[int]) {
		a.StackX(b, 1, 0)
		a.Pos.X += 8
	}
	wrap := func(a, b *align.Node[int]) {
		a.StackY(b, 0, 1)
		a.Pos.Y += 8
	}
	bounds := align.WH(100, 100)
	w := align.NewWrapper(bounds, 0, 0, stack, wrap)
	for _, d := range []int{40, 40, 40, 40, 40, 40} {
		n := align.WH(d, d)
		if !w.Add(n) {
			break
		}
	}
	for _, n := range w.Nodes() {
		fmt.Println(n)
	}

	// Output:
	// Pos: (0,0), Size: (40,40)
	// Pos: (48,0), Size: (40,40)
	// Pos: (0,48), Size: (40,40)
	// Pos: (48,48), Size: (40,40)
}

// Example_inventory demonstrates a game inventory grid layout
func Example_inventory() {
	// ゲーム画面とインベントリウィンドウ
	screen := align.WH(800, 600)
	inventoryWindow := align.WH(320, 240).Nest(screen, 0.5, 0.5)

	// インベントリグリッド（8x6のアイテムスロット）
	slotSize := 32
	padding := 4
	slots := make([]*align.Node[int], 0, 48)

	for row := range 6 {
		for col := range 8 {
			slot := align.WH(slotSize, slotSize)
			if row == 0 && col == 0 {
				// 最初のスロットをインベントリウィンドウの左上に配置
				slot.Nest(inventoryWindow.Inset(16), 0, 0)
			} else if col == 0 {
				// 各行の最初のスロットは前の行の最初のスロットの下に配置
				slot.StackY(slots[(row-1)*8], 0, 1)
				slot.Pos.Y += padding
			} else {
				// それ以外は左のスロットの右に配置
				slot.StackX(slots[len(slots)-1], 1, 0)
				slot.Pos.X += padding
			}
			slots = append(slots, slot)
		}
	}

	// 最初と最後のスロットの位置を表示
	fmt.Printf("First slot: %s\n", slots[0])
	fmt.Printf("Last slot: %s\n", slots[47])
	fmt.Printf("Inventory window: %s\n", inventoryWindow)

	// Output:
	// First slot: Pos: (256,196), Size: (32,32)
	// Last slot: Pos: (508,376), Size: (32,32)
	// Inventory window: Pos: (240,180), Size: (320,240)
}

// Example_dialog demonstrates a game dialog box layout
func Example_dialog() {
	// ゲーム画面とダイアログボックス
	screen := align.WH(800, 600)
	dialogBounds := align.WH(600, 200).Nest(screen, 0.5, 1).Offset(align.XY(0, -40))

	// ダイアログ内の要素
	contentBounds := dialogBounds.Inset(20)
	portrait := align.WH(120, 120).Nest(contentBounds, 0, 0.5)
	textArea := align.WH(400, 120).StackX(portrait, 1, 0.5).Offset(align.XY(20, 0))

	// 選択肢ボタン用の領域
	buttonW, buttonH := 150, 40
	// ダイアログ内のボタン配置ガイド（左右40、下30の余白）
	// buttonGuide := dialogBounds.InsetLTRB(40, 0, 40, 30)
	buttonGuide := dialogBounds.InsetXY(40, 30)
	choice1 := align.WH(buttonW, buttonH).Nest(buttonGuide, left, bottom)
	choice2 := align.WH(buttonW, buttonH).Nest(buttonGuide, right, bottom)

	fmt.Printf("Dialog box: %s\n", dialogBounds)
	fmt.Printf("Portrait: %s\n", portrait)
	fmt.Printf("Text area: %s\n", textArea)
	fmt.Printf("Choice 1: %s\n", choice1)
	fmt.Printf("Choice 2: %s\n", choice2)

	// Output:
	// Dialog box: Pos: (100,360), Size: (600,200)
	// Portrait: Pos: (120,400), Size: (120,120)
	// Text area: Pos: (260,400), Size: (400,120)
	// Choice 1: Pos: (140,490), Size: (150,40)
	// Choice 2: Pos: (510,490), Size: (150,40)
}

// Example_hud demonstrates a game HUD (Heads-Up Display) layout
func Example_hud() {
	// ゲーム画面
	screen := align.WH(800, 600)

	// 上部のステータスバー
	topBar := align.WH(screen.Size.X, 60).Nest(screen, mid, top)

	// HP/MPバー用の領域
	barWidth, barHeight := 200, 20
	barsGuide := topBar.InsetXY(60, 12)
	hpBar := align.WH(barWidth, barHeight).Nest(barsGuide, left, top)
	mpBar := align.WH(barWidth, barHeight).StackY(hpBar, left, bottom).Offset(align.XY(0, 4))

	// レベルと経験値
	levelDisplay := align.WH(80, 40).CenterOf(topBar)
	expBar := align.WH(150, 10).StackY(levelDisplay, mid, bottom).Offset(align.XY(0, 4))

	// ミニマップ
	minimapGuide := screen.Inset(20)
	minimap := align.WH(150, 150).Nest(minimapGuide, right, top)

	// スキルバー（画面下部中央）
	skillBarGuide := screen.Inset(20)
	skillBar := align.WH(400, 60).Nest(skillBarGuide, 0.5, 1)
	skillSlots := make([]*align.Node[int], 8)
	for i := range 8 {
		slot := align.WH(40, 40)
		if i == 0 {
			slot.Nest(skillBar, left, mid).Offset(align.XY(20, 0))
		} else {
			slot.StackX(skillSlots[i-1], right, mid).Offset(align.XY(10, 0))
		}
		skillSlots[i] = slot
	}

	fmt.Printf("HP bar: %s\n", hpBar)
	fmt.Printf("MP bar: %s\n", mpBar)
	fmt.Printf("Level display: %s\n", levelDisplay)
	fmt.Printf("EXP bar: %s\n", expBar)
	fmt.Printf("Minimap: %s\n", minimap)
	fmt.Printf("Skill bar: %s\n", skillBar)
	fmt.Printf("First skill: %s\n", skillSlots[0])
	fmt.Printf("Last skill: %s\n", skillSlots[7])

	// Output:
	// HP bar: Pos: (60,12), Size: (200,20)
	// MP bar: Pos: (60,36), Size: (200,20)
	// Level display: Pos: (360,10), Size: (80,40)
	// EXP bar: Pos: (325,54), Size: (150,10)
	// Minimap: Pos: (630,20), Size: (150,150)
	// Skill bar: Pos: (200,520), Size: (400,60)
	// First skill: Pos: (220,530), Size: (40,40)
	// Last skill: Pos: (570,530), Size: (40,40)
}
