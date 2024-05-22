package cpu

import (
	"fmt"
	"math/rand"
	"os"
)

var fontSet = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, //0
	0x20, 0x60, 0x20, 0x20, 0x70, //1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, //2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, //3
	0x90, 0x90, 0xF0, 0x10, 0x10, //4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, //5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, //6
	0xF0, 0x10, 0x20, 0x40, 0x40, //7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, //8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, //9
	0xF0, 0x90, 0xF0, 0x90, 0x90, //A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, //B
	0xF0, 0x80, 0x80, 0x80, 0xF0, //C
	0xE0, 0x90, 0x90, 0x90, 0xE0, //D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, //E
	0xF0, 0x80, 0xF0, 0x80, 0x80, //F
}

type Chip8 struct {
	display [32][64]uint8

	memory [4096]uint8 // Memory
	vx     [16]uint8   // Registers
	key    [16]uint8   // Input
	stack  [16]uint16  // Stack

	opCode uint16 // Current opcode
	pc     uint16 // Program counter
	sp     uint16 // Stack pointer
	iv     uint16 // Index register

	delayTimer uint8 // Delay timer
	soundTimer uint8 // Sound Timer

	shouldDraw bool
}

func Init() Chip8 {
	c := Chip8{
		shouldDraw: true,
		pc:         0x200,
	}

	copy(c.memory[:], fontSet)

	return c
}

func (c *Chip8) Buffer() [32][64]uint8 {
	return c.display
}

func (c *Chip8) Draw() bool {
	sd := c.shouldDraw
	c.shouldDraw = false
	return sd
}

func (c *Chip8) KeyPress(num uint8) {
	c.key[num] = 1
}

func (c *Chip8) KeyRelease(num uint8) {
	c.key[num] = 0
}

func (c *Chip8) Cycle() {
	c.opCode = uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])

	fmt.Printf("OPCODE: %#x\n", c.opCode&0xFFFF)

	switch c.opCode & 0xF000 {
	case 0x0000:
		switch c.opCode & 0x00FF {
		// CLR
		case 0x00E0:
			for i := 0; i < 32; i++ {
				for j := 0; j < 64; j++ {
					c.display[i][j] = 0x0
				}
			}
			c.shouldDraw = true
			c.pc += 2

		// RET
		case 0x00EE:
			c.sp--
			c.pc = c.stack[c.sp]
			c.pc += 2

		default:
			panic("Invalid OP CODE 0x0FFF")
		}

	// 1nnn - jump to nnn
	case 0x1000:
		c.pc = c.opCode & 0x0FFF

	// 0nnn - call subroutine
	case 0x2000:
		c.stack[c.sp] = c.pc
		c.sp++
		c.pc = c.opCode & 0x0FFF

	// 3xkk - skip if vx = kk
	case 0x3000:
		if uint16(c.vx[(c.opCode&0x0F00)>>8]) == (c.opCode & 0x00FF) {
			c.pc += 4
		} else {
			c.pc += 2
		}

	// 4xkk - skip if vx != kk
	case 0x4000:
		if c.vx[(c.opCode&0x0F00)>>8] != uint8(c.opCode&0x00FF) {
			c.pc += 4
		} else {
			c.pc += 2
		}

	// 5xy0 - skip if vx = vy
	case 0x5000:
		if c.vx[(c.opCode&0x0F00)>>8] == c.vx[(c.opCode&0x00F0)>>4] {
			c.pc += 4
		} else {
			c.pc += 2
		}

	// 6xkk - set vx = kk
	case 0x6000:
		c.vx[(c.opCode&0x0F00)>>8] = uint8(c.opCode & 0x00FF)
		c.pc += 2

	// 7xkk - add kk to vx
	case 0x7000:
		c.vx[(c.opCode&0x0F00)>>8] += uint8(c.opCode & 0x00FF)
		c.pc += 2

	case 0x8000:
		switch c.opCode & 0x000F {
		// 8xy0 - set vx = vy
		case 0x0000:
			c.vx[(c.opCode&0x0F00)>>8] = c.vx[(c.opCode&0x00F0)>>4]
			c.pc += 2

		// 8xy1 - vx = vx OR vy
		case 0x0001:
			c.vx[(c.opCode&0x0F00)>>8] |= c.vx[(c.opCode&0x00F0)>>4]
			c.pc += 2

		// 8xy2 - vx = vx AND vy
		case 0x0002:
			c.vx[(c.opCode&0x0F00)>>8] &= c.vx[(c.opCode&0x00F0)>>4]
			c.pc += 2

		// 8xy3 - vx = vx XOR vy
		case 0x0003:
			c.vx[(c.opCode&0x0F00)>>8] ^= c.vx[(c.opCode&0x00F0)>>4]
			c.pc += 2

		// 8xy4 - add vx vy
		case 0x0004:
			var res uint16 = uint16(c.vx[(c.opCode&0x0F00)>>8] + c.vx[(c.opCode&0x00F0)>>4])
			c.vx[(c.opCode&0x0F00)>>8] = uint8(res & 0x00FF)
			if res < 256 {
				c.vx[0xF] = 0
			} else {
				c.vx[0xF] = 1
			}
			c.pc += 2

		// 8xy5 - sub vx vy
		case 0x0005:
			if c.vx[(c.opCode&0x0F00)>>8] > c.vx[(c.opCode&0x00F0)>>4] {
				c.vx[0xF] = 1
				c.vx[(c.opCode&0x0F00)>>8] -= c.vx[(c.opCode&0x00F0)>>4]
			} else {
				c.vx[0xF] = 0
				c.vx[(c.opCode&0x0F00)>>8] = c.vx[(c.opCode&0x00F0)>>4] - c.vx[(c.opCode&0x0F00)>>8]
			}
			c.pc += 2

		// 8xy6 - vx = vx shr 1
		case 0x0006:
			c.vx[0xF] = c.vx[(c.opCode&0x0F00)>>8] & 0x1
			c.vx[(c.opCode&0x0F00)>>8] /= 2
			c.pc += 2

		// 8xy7 - subn vx vy
		case 0x0007:
			if c.vx[(c.opCode&0x0F00)>>8] < c.vx[(c.opCode&0x00F0)>>4] {
				c.vx[0xF] = 1
			} else {
				c.vx[0xF] = 0
			}
			c.vx[(c.opCode&0x0F00)>>8] = c.vx[(c.opCode&0x00F0)>>4] - c.vx[(c.opCode&0x0F00)>>8]
			c.pc += 2

		// SHL vx {, vy}
		case 0x000E:
			c.vx[0xF] = c.vx[(c.opCode&0x0F00)>>8] >> 7
			c.vx[(c.opCode&0x0F00)>>8] *= 2
			c.pc += 2

		default:
			fmt.Println("Invalid OP CODE")
		}

	case 0x9000:
		if c.vx[(c.opCode&0x0F00)>>8] != c.vx[(c.opCode&0x00F0)>>4] {
			c.pc += 4
		} else {
			c.pc += 2
		}

	case 0xA000:
		c.iv = c.opCode & 0x0FFF
		c.pc += 2

	case 0xB000:
		c.pc = c.opCode&0x0FFF + uint16(c.vx[0x0])

	case 0xC000:
		c.vx[(c.opCode&0x0F00)>>8] = uint8(rand.Intn(256)) & uint8(c.opCode&0x00FF)
		c.pc += 2

	case 0xD000:
		x := c.vx[(c.opCode&0x0F00)>>8]
		y := c.vx[(c.opCode&0x00F0)>>4]
		n := c.opCode & 0x000F

		c.vx[0xF] = 0

		var i uint16 = 0
		var j uint16 = 0

		for j = 0; j < uint16(n); j++ {
			pixel := c.memory[c.iv+j]
			for i = 0; i < 8; i++ {
				if pixel&(0x80>>i) != 0 {
					displayX := (x + uint8(i)) % 64
					displayY := (y + uint8(j)) % 32
					if c.display[displayY][displayX] == 1 {
						c.vx[0xF] = 1
					}
					c.display[displayY][displayX] ^= 1
				}
			}
		}

		c.shouldDraw = true
		c.pc += 2

	case 0xE000:
		switch c.opCode & 0x00FF {
		case 0x009E:
			if c.key[c.vx[(c.opCode&0x0F00)>>8]] == 1 {
				c.pc += 4
			} else {
				c.pc += 2
			}

		case 0x00A1:
			if c.key[c.vx[(c.opCode&0x0F00)>>8]] == 0 {
				c.pc += 4
			} else {
				c.pc += 2
			}

		default:
			fmt.Println("Invalid OP CODE")

		}

	case 0xF000:
		switch c.opCode & 0x00FF {
		case 0x0007:
			c.vx[(c.opCode&0x0F00)>>8] = c.delayTimer
			c.pc += 2

		case 0x000A:
			pressed := false

			for i := 0; i < len(c.key); i++ {
				if c.key[i] != 0 {
					c.vx[(c.opCode&0x0F00)>>8] = uint8(i)
					pressed = true

				}
			}

			if !pressed {
				return
			}
			c.pc += 2

		case 0x0015:
			c.delayTimer = c.vx[(c.opCode&0x0F00)>>8]
			c.pc += 2

		case 0x0018:
			c.soundTimer = c.vx[(c.opCode&0x0F00)>>8]
			c.pc += 2

		case 0x001E:
			c.iv += uint16(c.vx[(c.opCode&0x0F00)>>8])
			c.pc += 2

		case 0x0029:
			c.iv = uint16(c.vx[(c.opCode&0x0F00)>>8] * 0x5)
			c.pc += 2

		case 0x0033:
			c.memory[c.iv] = c.vx[(c.opCode&0x0F00)>>8] / 100
			c.memory[c.iv+1] = (c.vx[(c.opCode&0x0F00)>>8] / 10) % 10
			c.memory[c.iv+2] = (c.vx[(c.opCode&0x0F00)>>8] % 100) % 10
			c.pc += 2

		case 0x0055: // 0xFX55 Stores V0 to VX (including VX) in memory starting at address I. I is increased by 1 for each value written
			x := (c.opCode & 0x0F00) >> 8
			for i := uint16(0); i <= x; i++ {
				c.memory[c.iv+i] = c.vx[i]
			}
			c.pc += 2

		case 0x0065:
			x := (c.opCode & 0x0F00) >> 8
			for i := uint16(0); i <= x; i++ {
				c.vx[i] = c.memory[c.iv+i]
			}
			c.pc += 2

		default:
			// for i := 0x200; i < 0x200+80; i++ {
			// 	fmt.Printf("%#x ", c.memory[i]&0xFF)
			// }

			os.Exit(1)

		}

	default:
		fmt.Printf("Invalid OP CODE")
	}

	if c.delayTimer > 0 {
		c.delayTimer--
	}

}

func (c *Chip8) LoadProgram(fileName string) error {
	file, fileErr := os.OpenFile(fileName, os.O_RDONLY, 0777)
	if fileErr != nil {
		return fileErr
	}
	defer file.Close()

	fStat, fStatErr := file.Stat()
	if fStatErr != nil {
		return fStatErr
	}
	if int64(len(c.memory)-0x200) < fStat.Size() { // program is loaded at 0x200
		return fmt.Errorf("Program size bigger than memory")
	}

	buffer := make([]byte, fStat.Size())
	if _, readErr := file.Read(buffer); readErr != nil {
		return readErr
	}

	for i := 0; i < len(buffer); i++ {
		c.memory[i+0x200] = buffer[i]
	}

	return nil

}
