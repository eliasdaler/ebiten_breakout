package main

import (
	"embed"
	"image/color"
	_ "image/png"
	"log"
	"math"

	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed data
var data embed.FS

const (
	// internal game size (without scaling)
	gameScreenWidth  = 256
	gameScreenHeight = 240

	scale        = 3 // scale 300% in window
	windowWidth  = gameScreenWidth * scale
	windowHeight = gameScreenHeight * scale

	dt              = 1 / 60.0 // assume that delta is fixed and we're always running at 60 FPS
	audioSampleRate = 48000
)

type GameState int

const (
	GameStateMenu GameState = iota
	GameStatePlaying
)

type Game struct {
	bg                 *Sprite
	border             FloatRect
	paddle             *Paddle
	ball               *Ball
	attachBallToPaddle bool
	blocks             []*Block

	font        font.Face
	dialogueBox *Sprite

	state     GameState
	menuState MenuState

	audioContext *audio.Context
	music        *Music

	sounds map[string]*Sound
}

func main() {
	g := &Game{}
	if err := g.Init(); err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Ebiten Breakout")
	ebiten.SetWindowResizable(false)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

	if err := g.Exit(); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return gameScreenWidth, gameScreenHeight
}

func (g *Game) Init() error {
	g.audioContext = audio.NewContext(audioSampleRate)

	music, err := NewMusicFromFile(g.audioContext, "data/music/song.ogg")
	if err != nil {
		return err
	}
	g.music = music

	if err = g.loadSounds(); err != nil {
		return err
	}

	if err = g.loadMenuResources(); err != nil {
		return err
	}

	if err = g.loadObjects(); err != nil {
		return err
	}

	// finally, start going
	g.music.Play()
	g.showMenu(MenuMain)
	g.resetGame()

	return nil
}

func (g *Game) Exit() error {
	if err := g.music.Close(); err != nil {
		return err
	}
	for _, sound := range g.sounds {
		if err := sound.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (g *Game) loadSounds() error {
	g.sounds = make(map[string]*Sound)
	for _, soundName := range []string{
		"ball_collision",
		"ball_collision_paddle",
		"ball_fall",
	} {
		s, err := NewSoundFromFile(g.audioContext, "data/sounds/"+soundName+".wav")
		if err != nil {
			return err
		}
		g.sounds[soundName] = s
	}
	return nil
}

func (g *Game) playSound(soundName string) {
	if s, ok := g.sounds[soundName]; ok {
		s.Play()
	}
}

func (g *Game) loadObjects() error {
	g.border = FloatRect{32, 16, 192, 244}

	// paddle
	paddle, err := makePaddle()
	if err != nil {
		return err
	}
	g.paddle = paddle

	// ball
	ball, err := makeBall()
	if err != nil {
		return err
	}
	g.ball = ball

	err = g.loadBlocks()
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) loadBlocks() error {
	// Here we're just creating numBlocksY rows of numBlocksX blocks in a row
	// Can do something more complicated here like loading a level from some file
	const (
		numBlocksX = 5
		numBlocksY = 4
	)
	blockMargin := Vec2f{48, 48} // position of top-left block

	// reuse the same image for all blocks
	img, _, err := NewImageFromFile("data/images/block.png")
	if err != nil {
		return err
	}

	blockW := img.Bounds().Dx()
	blockH := img.Bounds().Dy()

	for y := 0; y < numBlocksY; y++ {
		for x := 0; x < numBlocksX; x++ {
			blockPos := Vec2f{blockMargin.X + float64(x*blockW), blockMargin.Y + float64(y*blockH)}
			block, err := makeBlock(blockPos, img)
			if err != nil {
				return err
			}

			g.blocks = append(g.blocks, block)
		}
	}
	return nil
}

func (g *Game) showMenu(menuState MenuState) {
	g.state = GameStateMenu
	g.menuState = menuState
}

func (g *Game) resetGame() {
	// reset paddle
	g.paddle.Transform.Position = Vec2f{88, 208}

	// reset ball
	g.ball.Movement.Velocity = Vec2f{}
	g.ball.updateFrame()
	g.attachBallToPaddle = true
	g.setBallOnPaddle()

	// reset blocks
	for _, b := range g.blocks {
		b.IsBroken = false
	}
}

// adjusts ball position to sit in the center of the paddle
func (g *Game) setBallOnPaddle() {
	paddleW := g.paddle.Graphics.Sprite.GetSize().X

	ballW := g.ball.Graphics.Sprite.GetSize().X
	ballH := g.ball.Graphics.Sprite.GetSize().Y

	ballX := (g.paddle.Transform.Position.X + paddleW/2.) - ballW/2.
	ballY := g.paddle.Transform.Position.Y - ballH

	g.ball.Transform.Position = Vec2f{ballX, ballY}
	g.ball.Movement.Velocity = Vec2f{}
}

func (g *Game) Update() error {
	g.processInput()

	if g.state == GameStatePlaying {
		g.updatePaddle()
		g.updateBall()

		if g.checkGameWin() {
			g.onGameWin()
		}
	}

	return nil
}

func (g *Game) processInput() {
	switch g.state {
	case GameStateMenu:
		if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
			g.onStartGame()
		}
	case GameStatePlaying:
		stickState := GetStickState()
		p := g.paddle
		p.Movement.Velocity.X = stickState.X * p.Movement.Speed.X

		if g.attachBallToPaddle {
			if inpututil.IsKeyJustPressed(ebiten.KeyW) {
				g.state = GameStatePlaying
				g.launchBall()
			}
		}
	}
}

func (g *Game) onStartGame() {
	g.resetGame()
	g.state = GameStatePlaying
}

// launches the ball up
func (g *Game) launchBall() {
	g.attachBallToPaddle = false
	g.ball.Movement.Velocity = Vec2f{g.ball.Movement.Speed.X, -g.ball.Movement.Speed.Y}
	g.playSound("ball_collision_paddle")
	g.ball.updateFrame()
}

func (g *Game) updatePaddle() {
	t := g.paddle.Transform
	m := g.paddle.Movement

	t.Position.X += m.Velocity.X * dt

	// handle collision if collided into a wall
	newPos := checkBorderCollision(g.paddle.getAABB(), g.getBorderAABB())
	t.Position = newPos
}

func (g *Game) getBorderAABB() FloatRect {
	return FloatRect{g.border.X, g.border.Y, g.border.W, g.border.H}
}

// If eAABB and borderAABB collide, returns new position for first rectangle to resolve collision
func checkBorderCollision(eAABB, borderAABB FloatRect) Vec2f {
	eL, eU, eR, _ := eAABB.Corners()
	eW := eR - eL

	bL, bU, bR, _ := borderAABB.Corners()

	newPos := Vec2f{eL, eU}
	// left wall
	if eL < bL {
		newPos.X = bL
	}

	// right wall
	if eR > bR {
		newPos.X = bR - eW
	}

	// upper wall
	if eU < bU {
		newPos.Y = bU
	}

	return newPos
}

func checkBallLost(ball *Ball, borderAABB FloatRect) bool {
	if ball.Transform.Position.Y > borderAABB.Y+borderAABB.H {
		// fell down
		return true
	}
	return false
}

func (g *Game) onBallLost() {
	// TODO: lives, etc.
	g.playSound("ball_fall")
	g.showMenu(MenuGameOver)
}

func (g *Game) updateBall() {
	if g.attachBallToPaddle {
		g.setBallOnPaddle()
		return
	}

	if checkBallLost(g.ball, g.getBorderAABB()) {
		g.onBallLost()
		return
	}

	collided := false
	collidedWithPaddle := false
	if handleBallPaddleCollision(g.ball, g.paddle) {
		collided = true
		collidedWithPaddle = true
	}
	if handleBallBlocksCollision(g.ball, g.blocks) {
		collided = true
	}
	if handleBallBorderCollision(g.ball, g.getBorderAABB()) {
		collided = true
	}

	if collided {
		g.ball.updateFrame()
		if collidedWithPaddle {
			g.playSound("ball_collision_paddle")
		} else {
			g.playSound("ball_collision")
		}
	}

	t := g.ball.Transform
	m := g.ball.Movement
	t.Position.X += m.Velocity.X * dt
	t.Position.Y += m.Velocity.Y * dt
}

func handleBallPaddleCollision(ball *Ball, paddle *Paddle) (collided bool) {
	ballAABB := ball.getAABB()
	paddleAABB := paddle.getAABB()

	d := getIntersectionDepth(ballAABB, paddleAABB)
	if d == (Vec2f{}) || // no collision
		ballAABB.Y+ballAABB.H > paddleAABB.Y+paddleAABB.H/2. { // this means that the ball is much lower than paddle and should fall through
		return false
	}

	ball.Transform.Position.Y -= d.Y // resolve collision

	angle := getBallReflectionAngle(ball, paddle)
	ball.Movement.Velocity = getNewVelocity(angle, ball.Movement.Speed)
	return true
}

// Returns an angle at which the ball will be launched from the paddle
// If the ball hits the paddle at leftmost point, the angle will be 3*pi/4
// If the ball hits the rightmost point, the angle will be pi/4
// If the ball hits the center, the angle will be pi/2
func getBallReflectionAngle(ball *Ball, paddle *Paddle) float64 {
	pW := paddle.Graphics.Sprite.GetSize().X
	pRX := paddle.Transform.Position.X + pW                               // paddle rightmost point
	bCX := ball.Transform.Position.X + ball.Graphics.Sprite.GetSize().X/2 // ball center

	f := (pRX - bCX) / pW             // right paddle edge = 0, left paddle edge = 1
	f = math.Max(0., math.Min(f, 1.)) // clamp to [0;1] range

	return math.Pi/4. + f*math.Pi/2 // project [0;1] to [math.Pi/4; 3*math.Pi/4]
}

func getNewVelocity(angle float64, speed Vec2f) Vec2f {
	s := math.Sqrt(speed.X*speed.X + speed.Y*speed.Y) // max allowed speed
	return Vec2f{s * math.Cos(angle), -s * math.Sin(angle)}
}

func handleBallBlocksCollision(ball *Ball, blocks []*Block) (collided bool) {
	m := ball.Movement
	ballAABB := ball.getAABB()

	for _, b := range blocks {
		if b.IsBroken {
			continue
		}

		d := getIntersectionDepth(ballAABB, b.getAABB())
		if d != (Vec2f{}) {
			b.IsBroken = true

			if math.Abs(d.Y) < math.Abs(d.X) { // collided from top or bottom
				m.Velocity.Y = -m.Velocity.Y
			} else { // collided from left or right
				m.Velocity.X = -m.Velocity.X
			}
			return true
		}
	}
	return false
}

func handleBallBorderCollision(ball *Ball, borderAABB FloatRect) (collided bool) {
	ballAABB := ball.getAABB()

	t := ball.Transform
	m := ball.Movement
	if newPos := checkBorderCollision(ballAABB, borderAABB); t.Position != newPos {
		if t.Position.X != newPos.X { // hit vertical wall
			m.Velocity.X = -m.Velocity.X
		}
		if t.Position.Y != newPos.Y { // hit horizontal wall
			m.Velocity.Y = -m.Velocity.Y
		}
		t.Position = newPos
		return true
	}
	return false
}

func (g *Game) checkGameWin() bool {
	for _, b := range g.blocks {
		if !b.IsBroken {
			return false
		}
	}
	return true
}

func (g *Game) onGameWin() {
	g.ball.Movement.Velocity = Vec2f{}
	g.paddle.Movement.Velocity = Vec2f{}

	g.showMenu(MenuWin)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	g.bg.Draw(screen, ebiten.GeoM{})

	for _, b := range g.blocks {
		if !b.IsBroken {
			b.Draw(screen)
		}
	}

	g.paddle.Draw(screen)
	g.ball.Draw(screen)

	if g.state == GameStateMenu {
		g.drawMenu(screen)
	}
}
