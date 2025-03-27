# Ball Physics Simulation

A 2D physics simulation of bouncing balls with collision detection, written in Go using the Ebiten game engine.

## Features

- Realistic 2D ball physics with gravity
- Elastic collision detection between multiple balls
- Wall collision handling
- Clean visualization with FPS counter
- Vector math operations (addition, subtraction, dot/cross products)
- Unit vector calculations and angle measurements

## How It Works

The simulation implements:
- Position and velocity vectors for each ball
- Gravity force applied to all objects
- Momentum conservation during collisions
- Perfectly elastic collisions (energy conserved)
- Boundary detection and response

## Installation

1. Make sure you have Go installed (version 1.16+ recommended)
2. Install Ebiten:
   ```bash
   go get github.com/hajimehoshi/ebiten/v2
   ```
3. Clone or download this repository
4. Run the simulation:
   ```bash
   go run main.go
   ```

## Controls

- The simulation runs automatically
- Close the window to exit

## Technical Details

- Uses vector mathematics for all physics calculations
- Implements proper collision normal calculations
- Handles multiple simultaneous collisions
- Optimized for smooth performance

## Customization

You can easily modify:
- Number of balls
- Starting positions and velocities
- Gravity strength
- Ball radius
- Screen dimensions

## Requirements

- Go 1.16+
- Ebiten v2
