package animation

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/prizelobby/reverset-raiders/ui"
)

type Loc struct {
	X, Y int
}

type Moveable interface {
	MoveTo(x, y int)
}

type Anim interface {
	Update()
	IsFinished() bool
}

type CreatureSpriteUpdatePower struct {
	Finished       bool
	CreatureSprite *ui.CreatureSprite
	PowerDelta     int
}

func NewCreatureSpriteUpdatePower(cs *ui.CreatureSprite, v int) *CreatureSpriteUpdatePower {
	return &CreatureSpriteUpdatePower{
		CreatureSprite: cs,
		PowerDelta:     v,
		Finished:       false,
	}
}

func (c *CreatureSpriteUpdatePower) Update() {
	if !c.Finished {
		c.CreatureSprite.Power += c.PowerDelta
		c.Finished = true
	}
}

func (c *CreatureSpriteUpdatePower) IsFinished() bool {
	return c.Finished
}

type SpriteMovement struct {
	CurrentPathIndex int
	Path             []Loc
	Finished         bool
	M                Moveable
}

func NewSpriteMovement(startX, startY, endX, endY int, pathStepCount int, m Moveable) *SpriteMovement {
	stepSizeX := (endX - startX) / pathStepCount
	stepSizeY := (endY - startY) / pathStepCount
	path := make([]Loc, 0, pathStepCount+1)

	// Note: this may not look smooth because of floating point accumulation errors
	// but it should be sufficient for now
	for i := 0; i < pathStepCount; i++ {
		path = append(path, Loc{startX, startY})
		startX += stepSizeX
		startY += stepSizeY
	}
	path = append(path, Loc{endX, endY})

	return &SpriteMovement{
		Path:             path,
		CurrentPathIndex: 0,
		M:                m,
		Finished:         false,
	}
}

func (s *SpriteMovement) Update() {
	if s.Finished {
		return
	}

	s.CurrentPathIndex += 1
	if s.CurrentPathIndex >= len(s.Path) {
		s.Finished = true
		return
	}
	nextLoc := s.Path[s.CurrentPathIndex]
	s.M.MoveTo(nextLoc.X, nextLoc.Y)
}

func (s *SpriteMovement) IsFinished() bool {
	return s.Finished
}

type Animation struct {
	CurrentFrame     int
	LoopImages       bool
	LoopPath         bool
	CurrentPathIndex int
	Path             []Loc
	Frames           []*ebiten.Image
	Finished         bool
}

func NewLinearMovingSprite(startX, startY, endX, endY int, pathStepCount int, img *ebiten.Image) *Animation {
	stepSizeX := endX - startX
	stepSizeY := endY - startY
	path := make([]Loc, pathStepCount+1)

	// Note: this may not look smooth because of floating point accumulation errors
	// but it should be sufficient for now
	for i := 0; i < pathStepCount; i++ {
		path = append(path, Loc{startX, startY})
		startX += stepSizeX
		startY += stepSizeY
	}
	path = append(path, Loc{endX, endY})

	return &Animation{
		LoopImages:       true,
		LoopPath:         false,
		CurrentFrame:     0,
		CurrentPathIndex: 0,
		Path:             path,
		Frames:           []*ebiten.Image{img},
		Finished:         false,
	}
}

func (a *Animation) Update() {
	a.CurrentFrame += 1
	a.CurrentPathIndex += 1
	if a.CurrentFrame >= len(a.Frames) {
		if a.LoopImages {
			a.CurrentFrame = 0
		} else {
			a.Finished = true
		}
	}
	if a.CurrentPathIndex >= len(a.Path) {
		if a.LoopPath {
			a.CurrentPathIndex = 0
		} else {
			a.Finished = true
		}
	}
}

type DeathAnimation struct {
	CurrentFrame   int
	CreatureSprite *ui.CreatureSprite
}

func NewDeathAnimation(cs *ui.CreatureSprite) *DeathAnimation {
	return &DeathAnimation{
		CurrentFrame:   0,
		CreatureSprite: cs,
	}
}

func (d *DeathAnimation) Update() {
	d.CurrentFrame += 1
	d.CreatureSprite.Transp = 1.0 - (float32(d.CurrentFrame) / 20)
	d.CreatureSprite.Y -= 2
	if d.CurrentFrame >= 20 {
		d.CreatureSprite.Removed = true
	}
}

func (d *DeathAnimation) IsFinished() bool {
	return d.CurrentFrame >= 20
}

type SplatAnimation struct {
	CurrentFrame int
	SplateSprite *ui.SplatSprite
}

func NewSplatAnimation(s *ui.SplatSprite) *SplatAnimation {
	return &SplatAnimation{
		CurrentFrame: 0,
		SplateSprite: s,
	}
}

func (s *SplatAnimation) Update() {
	s.CurrentFrame += 1
	s.SplateSprite.Transp = 1.0 - (float32(s.CurrentFrame) / 40)
	s.SplateSprite.Y -= 1
	if s.CurrentFrame >= 40 {
		s.SplateSprite.Removed = true
	}
}

func (s *SplatAnimation) IsFinished() bool {
	return s.CurrentFrame >= 40
}

type TileHighlightAnimation struct {
	CurrentFrame int
	TileSprite1  *ui.TileSprite
	TileSprite2  *ui.TileSprite
}

func NewTileHighlightAnimation(t1, t2 *ui.TileSprite) *TileHighlightAnimation {
	return &TileHighlightAnimation{
		CurrentFrame: 0,
		TileSprite1:  t1,
		TileSprite2:  t2,
	}
}

func (t *TileHighlightAnimation) Update() {
	t.CurrentFrame += 1
	t.TileSprite1.Selected = true
	if t.TileSprite2 != nil {
		t.TileSprite2.Selected = true
	}

	if t.CurrentFrame >= 45 {
		t.TileSprite1.Selected = false
		if t.TileSprite2 != nil {
			t.TileSprite2.Selected = false
		}
	}
}

func (t *TileHighlightAnimation) IsFinished() bool {
	return t.CurrentFrame >= 45
}

type EffectSpriteAnimation struct {
	CurrentFrame int
	S            *ui.EffectSprite
}

func NewEffectSpriteAnimation(s *ui.EffectSprite) *EffectSpriteAnimation {
	return &EffectSpriteAnimation{CurrentFrame: 0, S: s}
}

func (e *EffectSpriteAnimation) Update() {
	e.CurrentFrame += 1
	e.S.I = e.CurrentFrame
}

func (e *EffectSpriteAnimation) IsFinished() bool {
	// note: we want there to be a 16th frame so that the effect sprite becomes empty after
	// the animation is complete
	return e.CurrentFrame > 15
}
