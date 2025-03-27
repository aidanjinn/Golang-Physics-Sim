package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type vector struct {
	x float64
	y float64
	z float64
}

func newVector(x float64, y float64, z float64) *vector {
	return &vector{x: x, y: y, z: z}
}

func (v *vector) magnitude() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

func (v *vector) angles() (float64, float64, float64) {
	// calculate angle between vector and x axis
	i := newVector(1, 0, 0)
	x_angle := math.Acos((v.x*i.x + v.y*i.y + v.z*i.z) / (v.magnitude() * i.magnitude()))

	// calculate angle between vector and y axis
	j := newVector(0, 1, 0)
	y_angle := math.Acos((v.x*j.x + v.y*j.y + v.z*j.z) / (v.magnitude() * j.magnitude()))

	// calculate angle between vector and z axis
	k := newVector(0, 0, 1)
	z_angle := math.Acos((v.x*k.x + v.y*k.y + v.z*k.z) / (v.magnitude() * k.magnitude()))

	return x_angle * 180 / math.Pi, y_angle * 180 / math.Pi, z_angle * 180 / math.Pi
}

func (v *vector) toString() string {
	return fmt.Sprintf("(%.2f,%.2f,%.2f)", v.x, v.y, v.z)
}

func add(vect1 vector, vect2 vector) vector {
	return vector{vect1.x + vect2.x, vect1.y + vect2.y, vect1.z + vect2.z}
}

func subtract(vect1 vector, vect2 vector) vector {
	return vector{vect1.x - vect2.x, vect1.y - vect2.y, vect1.z - vect2.z}
}

func scalar_mult(vect vector, scalar float64) vector {
	return vector{vect.x * scalar, vect.y * scalar, vect.z * scalar}
}

func cross_product(vect1 vector, vect2 vector) vector {
	return vector{
		vect1.y*vect2.z - vect1.z*vect2.y,
		vect1.z*vect2.x - vect1.x*vect2.z,
		vect1.x*vect2.y - vect1.y*vect2.x,
	}
}

func unit_vector(v vector) vector {
	magnitude := v.magnitude()
	return vector{v.x / magnitude, v.y / magnitude, v.z / magnitude}
}

func dot_product(vect1 vector, vect2 vector) float64 {
	return vect1.x*vect2.x + vect1.y*vect2.y + vect1.z*vect2.z
}

func angle_between_vectors(vect1 vector, vect2 vector) float64 {
	dot := dot_product(vect1, vect2)
	magnitude_product := vect1.magnitude() * vect2.magnitude()
	return math.Acos(dot/magnitude_product) * 180 / math.Pi
}

func projection(vect1 vector, vect2 vector) vector {
	dot := dot_product(vect1, vect2)
	magnitude_squared := dot_product(vect2, vect2)
	scale := dot / magnitude_squared
	return scalar_mult(vect2, scale)
}

func reflect(vect vector, normal vector) vector {
	dot := dot_product(vect, normal)
	return subtract(vect, scalar_mult(normal, 2*dot))
}

type Ball struct {
	ballPosition vector
	ballVelocity vector
}

type Game struct {
	objects []Ball
	gravity vector
}

const (
	screenWidth  = 640
	screenHeight = 480
	ballRadius   = 20
)

func (g *Game) Update() error {

	for i := range g.objects {

		currBall := &g.objects[i]
		currBall.ballVelocity = add(currBall.ballVelocity, g.gravity)
		currBall.ballPosition = add(currBall.ballPosition, currBall.ballVelocity)

		for j := i + 1; j < len(g.objects); j++ {
			otherBall := &g.objects[j]

			// Calculate distance between balls
			distanceVector := subtract(currBall.ballPosition, otherBall.ballPosition)
			distance := distanceVector.magnitude()

			// Check if balls are colliding
			if distance < 2*ballRadius {

				// Calculate collision normal (unit vector between centers)
				collisionNormal := unit_vector(distanceVector)

				// Calculate relative velocity
				relativeVelocity := subtract(currBall.ballVelocity, otherBall.ballVelocity)

				// Calculate velocity along the normal
				velocityAlongNormal := dot_product(relativeVelocity, collisionNormal)

				// Only proceed if balls are moving towards each other
				if velocityAlongNormal > 0 {
					continue
				}

				// Calculate impulse scalar (perfectly elastic collision)
				impulse := -(1 + 1.0) * velocityAlongNormal
				impulse /= 2 // Since both balls have equal mass in this case

				// Apply impulse
				impulseVector := scalar_mult(collisionNormal, impulse)

				// Update velocities
				currBall.ballVelocity = add(currBall.ballVelocity, impulseVector)
				otherBall.ballVelocity = subtract(otherBall.ballVelocity, impulseVector)

				// Separate balls to prevent sticking
				overlap := 2*ballRadius - distance
				separationVector := scalar_mult(collisionNormal, overlap/2)
				currBall.ballPosition = add(currBall.ballPosition, separationVector)
				otherBall.ballPosition = subtract(otherBall.ballPosition, separationVector)
			}
		}

		// If we are out of bounds left side
		if currBall.ballPosition.x-ballRadius < 0 {
			currBall.ballPosition.x = ballRadius
			currBall.ballVelocity.x *= -1

			// If we are out bounds right side
		} else if currBall.ballPosition.x+ballRadius > screenWidth {
			currBall.ballPosition.x = screenWidth - ballRadius
			currBall.ballVelocity.x *= -1
		}

		// If we are out bounds Bottom Side
		if currBall.ballPosition.y-ballRadius < 0 {
			currBall.ballPosition.y = ballRadius
			currBall.ballVelocity.y *= -1

			// If We are out of bounds Top Side
		} else if currBall.ballPosition.y+ballRadius > screenHeight {
			currBall.ballPosition.y = screenHeight - ballRadius
			currBall.ballVelocity.y *= -1
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	for _, ball := range g.objects {
		ebitenutil.DrawCircle(screen, float64(ball.ballPosition.x), float64(ball.ballPosition.y), ballRadius, color.White)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Bouncing Balls")

	game := &Game{
		objects: []Ball{
			{
				ballPosition: vector{x: 100, y: 100},
				ballVelocity: vector{x: 2, y: 3},
			},
			{
				ballPosition: vector{x: 300, y: 200},
				ballVelocity: vector{x: -1, y: -2},
			},
			{
				ballPosition: vector{x: 10, y: 150},
				ballVelocity: vector{x: 2, y: 3},
			},
			{
				ballPosition: vector{x: 20, y: 20},
				ballVelocity: vector{x: -1, y: -2},
			},
			{
				ballPosition: vector{x: 200, y: 100},
				ballVelocity: vector{x: 2, y: 3},
			},
			{
				ballPosition: vector{x: 30, y: 200},
				ballVelocity: vector{x: -1, y: -2},
			},
			{
				ballPosition: vector{x: 100, y: 100},
				ballVelocity: vector{x: 2, y: 3},
			},
			{
				ballPosition: vector{x: 300, y: 200},
				ballVelocity: vector{x: -1, y: -2},
			},
		},
		gravity: vector{x: 0, y: .3},
	}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
