package vm

const (
	Load  = 0x01
	Store = 0x02
	Add   = 0x03
	Sub   = 0x04
	Halt  = 0xff
)

// Stretch goals
const (
	Addi = 0x05
	Subi = 0x06
	Jump = 0x07
	Beqz = 0x08
)

// Given a 256 byte array of "memory", run the stored program
// to completion, modifying the data in place to reflect the result
//
// The memory format is:
//
// 00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f ... ff
// __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ ... __
// ^==DATA===============^ ^==INSTRUCTIONS==============^
//
func compute(memory []byte) {

	registers := [3]byte{8, 0, 0} // PC, R1 and R2
	var pcincr byte               // Amount to increment the PC each cycle

	// Keep looping, like a physical computer's clock
	for {
		pcincr = 3 // Default instruction length is three bytes

		instr := memory[registers[0] : registers[0]+pcincr] // The full instruction
		op := instr[0]

		// decode and execute
		switch op {
		case Load:
			registers[instr[1]] = memory[instr[2]]
		case Store:
			memory[instr[2]] = registers[instr[1]]
		case Add:
			registers[instr[1]] = registers[instr[1]] + registers[instr[2]]
		case Sub:
			registers[instr[1]] = registers[instr[1]] - registers[instr[2]]
		case Halt:
			return
		case Addi:
			registers[instr[1]] = registers[instr[1]] + instr[2]
		case Subi:
			registers[instr[1]] = registers[instr[1]] - instr[2]
		case Jump:
			pcincr = 0 // Don't increment the PC since we're setting it below
			registers[0] = instr[1]
		case Beqz:
			if registers[instr[1]] == 0 {
				pcincr += instr[2]
			}
		default:
			panic("Unkown Instruction")
		}

		registers[0] += pcincr // Increment PC to next instruction
	}
}
