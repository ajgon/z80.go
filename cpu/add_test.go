package cpu

import (
	"testing"
	"z80/dma"
	"z80/memory"
)

var addTruthTable [36][9]uint8 = [36][9]uint8{
	// a, r, a+r, C, N, PV, H, Z, S
	[9]uint8{0, 0, 0, 0, 0, 0, 0, 1, 0},
	[9]uint8{0, 1, 1, 0, 0, 0, 0, 0, 0},
	[9]uint8{0, 127, 127, 0, 0, 0, 0, 0, 0},
	[9]uint8{0, 128, 128, 0, 0, 0, 0, 0, 1},
	[9]uint8{0, 129, 129, 0, 0, 0, 0, 0, 1},
	[9]uint8{0, 255, 255, 0, 0, 0, 0, 0, 1},
	[9]uint8{1, 0, 1, 0, 0, 0, 0, 0, 0},
	[9]uint8{1, 1, 2, 0, 0, 0, 0, 0, 0},
	[9]uint8{1, 127, 128, 0, 0, 1, 1, 0, 1},
	[9]uint8{1, 128, 129, 0, 0, 0, 0, 0, 1},
	[9]uint8{1, 129, 130, 0, 0, 0, 0, 0, 1},
	[9]uint8{1, 255, 0, 1, 0, 0, 1, 1, 0},
	[9]uint8{127, 0, 127, 0, 0, 0, 0, 0, 0},
	[9]uint8{127, 1, 128, 0, 0, 1, 1, 0, 1},
	[9]uint8{127, 127, 254, 0, 0, 1, 1, 0, 1},
	[9]uint8{127, 128, 255, 0, 0, 0, 0, 0, 1},
	[9]uint8{127, 129, 0, 1, 0, 0, 1, 1, 0},
	[9]uint8{127, 255, 126, 1, 0, 0, 1, 0, 0},
	[9]uint8{128, 0, 128, 0, 0, 0, 0, 0, 1},
	[9]uint8{128, 1, 129, 0, 0, 0, 0, 0, 1},
	[9]uint8{128, 127, 255, 0, 0, 0, 0, 0, 1},
	[9]uint8{128, 128, 0, 1, 0, 1, 0, 1, 0},
	[9]uint8{128, 129, 1, 1, 0, 1, 0, 0, 0},
	[9]uint8{128, 255, 127, 1, 0, 1, 0, 0, 0},
	[9]uint8{129, 0, 129, 0, 0, 0, 0, 0, 1},
	[9]uint8{129, 1, 130, 0, 0, 0, 0, 0, 1},
	[9]uint8{129, 127, 0, 1, 0, 0, 1, 1, 0},
	[9]uint8{129, 128, 1, 1, 0, 1, 0, 0, 0},
	[9]uint8{129, 129, 2, 1, 0, 1, 0, 0, 0},
	[9]uint8{129, 255, 128, 1, 0, 0, 1, 0, 1},
	[9]uint8{255, 0, 255, 0, 0, 0, 0, 0, 1},
	[9]uint8{255, 1, 0, 1, 0, 0, 1, 1, 0},
	[9]uint8{255, 127, 126, 1, 0, 0, 1, 0, 0},
	[9]uint8{255, 128, 127, 1, 0, 1, 0, 0, 0},
	[9]uint8{255, 129, 128, 1, 0, 0, 1, 0, 1},
	[9]uint8{255, 255, 254, 1, 0, 0, 1, 0, 1},
}

func TestAddRegister(t *testing.T) {
	var mem = memory.MemoryNew()
	var dmaX = dma.DMANew(mem)
	var cpu = CPUNew(dmaX)

	for _, row := range addTruthTable {
		for _, register := range [7]byte{'B', 'C', 'D', 'E', 'H', 'L', 'A'} {
			if register == 'A' {
				if row[0] != row[1] {
					continue
				}
			}

			cpu.PC = 0
			cpu.setAcc(row[0])
			cpu.BC = (uint16(row[1]) << 8) | uint16(row[1])
			cpu.DE = (uint16(row[1]) << 8) | uint16(row[1])
			cpu.HL = (uint16(row[1]) << 8) | uint16(row[1])
			tstates := cpu.addAR(register)()

			if cpu.getAcc() != row[2] || cpu.getC() != (row[3] == 1) || cpu.getN() != (row[4] == 1) || cpu.getPV() != (row[5] == 1) || cpu.getH() != (row[6] == 1) || cpu.getZ() != (row[7] == 1) || cpu.getS() != (row[8] == 1) {
				t.Errorf(
					"\ngot:  A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t\nwant: A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t for (%d + %d)",
					cpu.getAcc(), cpu.getC(), cpu.getN(), cpu.getPV(), cpu.getH(), cpu.getZ(), cpu.getS(),
					row[2], row[3] == 1, row[4] == 1, row[5] == 1, row[6] == 1, row[7] == 1, row[8] == 1, row[0], row[1],
				)
			}

			if cpu.PC != 1 || tstates != 4 {
				t.Errorf("got PC=%d, %d T-states, want PC=%d, %d T-states", cpu.PC, tstates, 1, 4)
			}
		}
	}
}

func TestAdd_Hl_(t *testing.T) {
	var mem = memory.MemoryNew()
	var dmaX = dma.DMANew(mem)
	var cpu = CPUNew(dmaX)
	cpu.HL = 0x1234

	for _, row := range addTruthTable {
		cpu.PC = 0
		cpu.setAcc(row[0])
		dmaX.SetMemoryByte(cpu.HL, row[1])
		tstates := cpu.addA_Hl_()

		if cpu.getAcc() != row[2] || cpu.getC() != (row[3] == 1) || cpu.getN() != (row[4] == 1) || cpu.getPV() != (row[5] == 1) || cpu.getH() != (row[6] == 1) || cpu.getZ() != (row[7] == 1) || cpu.getS() != (row[8] == 1) {
			t.Errorf(
				"\ngot:  A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t\nwant: A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t for (%d + %d)",
				cpu.getAcc(), cpu.getC(), cpu.getN(), cpu.getPV(), cpu.getH(), cpu.getZ(), cpu.getS(),
				row[2], row[3] == 1, row[4] == 1, row[5] == 1, row[6] == 1, row[7] == 1, row[8] == 1, row[0], row[1],
			)
		}

		if cpu.PC != 1 || tstates != 7 {
			t.Errorf("got PC=%d, %d T-states, want PC=%d, %d T-states", cpu.PC, tstates, 1, 7)
		}
	}
}
