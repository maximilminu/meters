package src

import "gitlab.sistematis.com.ar/OC/be/common/db"

type EGMMeter struct {
	TotalIn             int
	TotalOut            int
	Bets                int
	Wins                int
	Jackpot             int
	Handpay             int
	Events              string
	Date                db.Time
	MachineDrop         int    `json:",omitempty"`
	DoorOpen            int    `json:",omitempty"`
	PowerReset          int    `json:",omitempty"`
	BasePorcent         string `json:",omitempty"`
	PhysicalCoinIn      int    `json:",omitempty"`
	PhysicalCoinOut     int    `json:",omitempty"`
	CancelCredits       int    `json:",omitempty"`
	BillIn              int    `json:",omitempty"`
	MachineProgresive   int    `json:",omitempty"`
	TotalTaxes          int    `json:",omitempty"`
	TaxesQuantity       int    `json:",omitempty"`
	CurrentCredits      int    `json:",omitempty"`
	AttendantProgresive int    `json:",omitempty"`
}

type Report struct {
	MachineID         int
	DenominationValue float32
	Periods           []Period
}

const PeriodTypeRegular = ""
const PeriodTypeRollover = "rollover"
const PeriodTypeClearFullRam = "clear_full_ram"
const PeriodTypeClearPartialRam = "clear_partial_ram"

type Period struct {
	MetersFrom MetersWithTime
	MetersTo   MetersWithTime
	Deltas     Meters
	Type       string
}

type MetersWithTime struct {
	At db.Time
	Meters
}

type Meter int

type Meters struct {
	Games           Meter
	TotalIn         Meter
	TotalOut        Meter
	Jackpot         Meter
	MachineDrop     Meter
	GameWon         Meter
	DoorOpen        Meter
	PowerReset      Meter
	CurrentCredits  Meter
	CancelCredits   Meter
	BillFive        Meter
	BillTen         Meter
	BillTwenty      Meter
	BillFifty       Meter
	BillOneHundred  Meter
	BillTwoHundred  Meter
	BillFiveHundred Meter
	BillOneThousand Meter
}

type ReportingMachine struct {
	MachineID   int
	DeltaBillIn int
	Discount    float32
	Date        db.Time
}

func (p *Period) CalcDelta() {
	p.Deltas.Games = p.MetersTo.Games - p.MetersFrom.Games
	p.Deltas.TotalIn = p.MetersTo.TotalIn - p.MetersFrom.TotalIn
	p.Deltas.TotalOut = p.MetersTo.TotalOut - p.MetersFrom.TotalOut
	p.Deltas.Jackpot = p.MetersTo.Jackpot - p.MetersFrom.Jackpot
	p.Deltas.MachineDrop = p.MetersTo.MachineDrop - p.MetersFrom.MachineDrop
	p.Deltas.GameWon = p.MetersTo.GameWon - p.MetersFrom.GameWon
	p.Deltas.DoorOpen = p.MetersTo.DoorOpen - p.MetersFrom.DoorOpen
	p.Deltas.PowerReset = p.MetersTo.PowerReset - p.MetersFrom.PowerReset
	p.Deltas.CurrentCredits = p.MetersTo.CurrentCredits - p.MetersFrom.CurrentCredits
	p.Deltas.CancelCredits = p.MetersTo.CancelCredits - p.MetersFrom.CancelCredits
	p.Deltas.BillFive = p.MetersTo.BillFive - p.MetersFrom.BillFive
	p.Deltas.BillTen = p.MetersTo.BillTen - p.MetersFrom.BillTen
	p.Deltas.BillTwenty = p.MetersTo.BillTwenty - p.MetersFrom.BillTwenty
	p.Deltas.BillFifty = p.MetersTo.BillFifty - p.MetersFrom.BillFifty
	p.Deltas.BillOneHundred = p.MetersTo.BillOneHundred - p.MetersFrom.BillOneHundred
	p.Deltas.BillTwoHundred = p.MetersTo.BillTwoHundred - p.MetersFrom.BillTwoHundred
	p.Deltas.BillFiveHundred = p.MetersTo.BillFiveHundred - p.MetersFrom.BillFiveHundred
	p.Deltas.BillOneThousand = p.MetersTo.BillOneThousand - p.MetersFrom.BillOneThousand
}

func (mts *Meters) CalcRollover(oldValue Meters, newValue Meters, def Meter) {

	if oldValue.Games > newValue.Games {
		mts.Games = def
	} else {
		mts.Games = oldValue.Games
	}

	if oldValue.TotalIn > newValue.TotalIn {
		mts.TotalIn = def
	} else {
		mts.TotalIn = oldValue.TotalIn
	}
	if oldValue.TotalOut > newValue.TotalOut {
		mts.TotalOut = def
	} else {
		mts.TotalOut = oldValue.TotalOut
	}
	if oldValue.Jackpot > newValue.Jackpot {
		mts.Jackpot = def
	} else {
		mts.Jackpot = oldValue.Jackpot
	}
	if oldValue.MachineDrop > newValue.MachineDrop {
		mts.MachineDrop = def
	} else {
		mts.MachineDrop = oldValue.MachineDrop
	}
	if oldValue.GameWon > newValue.GameWon {
		mts.GameWon = def
	} else {
		mts.GameWon = oldValue.GameWon
	}
	if oldValue.DoorOpen > newValue.DoorOpen {
		mts.DoorOpen = def
	} else {
		mts.DoorOpen = oldValue.DoorOpen
	}
	if oldValue.PowerReset > newValue.PowerReset {
		mts.PowerReset = def
	} else {
		mts.PowerReset = oldValue.PowerReset
	}
	if oldValue.CurrentCredits > newValue.CurrentCredits {
		mts.CurrentCredits = def
	} else {
		mts.CurrentCredits = oldValue.CurrentCredits
	}
	if oldValue.CancelCredits > newValue.CancelCredits {
		mts.CancelCredits = def
	} else {
		mts.CancelCredits = oldValue.CancelCredits
	}
	if oldValue.BillFive > newValue.BillFive {
		mts.BillFive = def
	} else {
		mts.BillFive = oldValue.BillFive
	}
	if oldValue.BillTen > newValue.BillTen {
		mts.BillTen = def
	} else {
		mts.BillTen = oldValue.BillTen
	}
	if oldValue.BillTwenty > newValue.BillTwenty {
		mts.BillTwenty = def
	} else {
		mts.BillTwenty = oldValue.BillTwenty
	}
	if oldValue.BillFifty > newValue.BillFifty {
		mts.BillFifty = def
	} else {
		mts.BillFifty = oldValue.BillFifty
	}
	if oldValue.BillOneHundred > newValue.BillOneHundred {
		mts.BillOneHundred = def
	} else {
		mts.BillOneHundred = oldValue.BillOneHundred
	}
	if oldValue.BillTwoHundred > newValue.BillTwoHundred {
		mts.BillTwoHundred = def
	} else {
		mts.BillTwoHundred = oldValue.BillTwoHundred
	}
	if oldValue.BillFiveHundred > newValue.BillFiveHundred {
		mts.BillFiveHundred = def
	} else {
		mts.BillFiveHundred = oldValue.BillFiveHundred
	}
	if oldValue.BillOneThousand > newValue.BillOneThousand {
		mts.BillOneThousand = def
	} else {
		mts.BillOneThousand = oldValue.BillOneThousand
	}
}

func (mts Meters) Clone() Meters {
	m := Meters{
		Games:           mts.Games,
		TotalIn:         mts.TotalIn,
		TotalOut:        mts.TotalOut,
		Jackpot:         mts.Jackpot,
		MachineDrop:     mts.MachineDrop,
		GameWon:         mts.GameWon,
		DoorOpen:        mts.DoorOpen,
		PowerReset:      mts.PowerReset,
		CurrentCredits:  mts.CurrentCredits,
		CancelCredits:   mts.CancelCredits,
		BillFive:        mts.BillFive,
		BillTen:         mts.BillTen,
		BillTwenty:      mts.BillTwenty,
		BillFifty:       mts.BillFifty,
		BillOneHundred:  mts.BillOneHundred,
		BillTwoHundred:  mts.BillTwoHundred,
		BillFiveHundred: mts.BillFiveHundred,
		BillOneThousand: mts.BillOneThousand,
	}
	return m
}

func (mts Meters) CheckIfRollover(newMeters Meters) bool {
	return mts.Games > newMeters.Games ||
		mts.TotalIn > newMeters.TotalIn ||
		mts.TotalOut > newMeters.TotalOut ||
		mts.Jackpot > newMeters.Jackpot ||
		mts.MachineDrop > newMeters.MachineDrop ||
		mts.GameWon > newMeters.GameWon ||
		mts.DoorOpen > newMeters.DoorOpen ||
		mts.PowerReset > newMeters.PowerReset ||
		mts.CurrentCredits > newMeters.CurrentCredits ||
		mts.CancelCredits > newMeters.CancelCredits ||
		mts.BillFive > newMeters.BillFive ||
		mts.BillTen > newMeters.BillTen ||
		mts.BillTwenty > newMeters.BillTwenty ||
		mts.BillFifty > newMeters.BillFifty ||
		mts.BillOneHundred > newMeters.BillOneHundred ||
		mts.BillTwoHundred > newMeters.BillTwoHundred ||
		mts.BillFiveHundred > newMeters.BillFiveHundred ||
		mts.BillOneThousand > newMeters.BillOneThousand
}
