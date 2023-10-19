package src

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.sistematis.com.ar/OC/be/common/db"
	helpers "gitlab.sistematis.com.ar/OPEN-SOURCE/GO/api-helpers"
)

func GET(router *gin.RouterGroup) {

	router.GET("/:machineIdUnico", func(c *gin.Context) {
		ocDB := db.MysqlDb("TNZ_LEGACY")
		log := helpers.Log(c)
		IdUnico := c.Param("machineIdUnico")

		log.Debug("Machine ID >>>>>>>>>>>>>> ", IdUnico)

		type QueryParams struct {
			Time db.Time
		}

		var params QueryParams
		if err := c.ShouldBindQuery(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		myTime := params.Time.AsTime()

		cursor, closeCursor, queryErr, sqlLog := helpers.SqlQuery(log, ocDB, `SELECT x0.hora_registro, ed.descripcion, c.coinIn, c.coinOut, c.jackpot, c.hpCancelCredit, c.gamePlayed, c.gameWon, d.denominacion FROM (SELECT * FROM (SELECT e.id_evento, e.hora_registro, e.id_contador, e.clave FROM eventos e WHERE e.id_maquina = (select id_maquina from maquinas where idUnico=?) AND e.hora_registro > subdate(?, interval 1 hour) ORDER BY e.hora_registro DESC LIMIT 20) j0 UNION SELECT * FROM (SELECT e.id_evento, e.hora_registro, e.id_contador, e.clave FROM eventos e WHERE e.id_maquina = (select id_maquina from maquinas where idUnico=?) AND e.hora_registro > adddate(?, interval 10 minute) ORDER BY e.hora_registro ASC LIMIT 10) j1) x0 INNER JOIN contadores c ON x0.id_contador = c.id_contador INNER JOIN eventosdetalle ed ON x0.clave = ed.clave INNER JOIN denominacion d ON c.denomination = d.codigo ORDER BY x0.hora_registro DESC`, IdUnico, myTime, IdUnico, myTime)

		if queryErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Query": queryErr.Error()})
			return
		}

		defer closeCursor()

		var (
			horario        time.Time
			descripcion    string
			coinIn         int
			coinOut        int
			jackpot        int
			hpCancelCredit int
			gamePlayed     int
			gameWon        int
			denominacion   string
		)

		meters := make([]EGMMeter, 0)

		for cursor.Next() {
			err := cursor.Scan(
				&horario,
				&descripcion,
				&coinIn,
				&coinOut,
				&jackpot,
				&hpCancelCredit,
				&gamePlayed,
				&gameWon,
				&denominacion,
			)

			if err != nil {
				sqlLog.Errorf("ERROR: %s", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"DB": err.Error()})
				return
			}

			egmmeter := EGMMeter{
				TotalIn:  coinIn,
				TotalOut: coinOut,
				Wins:     gameWon,
				Bets:     gamePlayed,
				Jackpot:  jackpot,
				Handpay:  hpCancelCredit,
				Events:   descripcion,
				Date:     db.Time(horario.UnixMilli()),
			}
			meters = append(meters, egmmeter)
		}
		c.JSON(http.StatusOK, meters)
	})

	router.GET("/by-date/:machineIdUnico", func(c *gin.Context) {
		ocDB := db.MysqlDb("TNZ_LEGACY")
		log := helpers.Log(c)
		IdUnico := c.Param("machineIdUnico")

		log.Debug("Machine ID >>>>>>>>>>>>>> ", IdUnico)

		type QueryParams struct {
			Time1 db.Time `json:"time1"`
			Time2 db.Time `json:"time2"`
			Limit int64   `json:"limit"`
		}

		params := QueryParams{
			Limit: 5,
		}

		if err := c.ShouldBindQuery(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		myTime1 := params.Time1.AsTime()
		myTime2 := params.Time2.AsTime()
		limit := params.Limit

		log.Debug("LIMIT >>>>>>>>>>> ", limit)

		cursor, closeCursor, queryErr, sqlLog := helpers.SqlQuery(log, ocDB, `SELECT x0.hora_registro, ed.descripcion, c.coinIn, c.coinOut, c.jackpot, c.gamePlayed, c.gameWon, c.hpCancelCredit, d.denominacion
		FROM (SELECT e.id_evento, e.hora_registro, e.id_contador, e.clave
			FROM eventos e WHERE e.id_maquina = (select id_maquina from maquinas where idUnico=?) AND e.hora_registro BETWEEN ? AND ?) x0 INNER JOIN contadores c ON x0.id_contador = c.id_contador INNER JOIN denominacion d ON c.denomination = d.codigo INNER JOIN eventosdetalle ed ON x0.clave = ed.clave ORDER BY x0.hora_registro DESC
			LIMIT ?`, IdUnico, myTime1, myTime2, limit)

		if queryErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Query": queryErr.Error()})
			return
		}

		defer closeCursor()

		var (
			horario        time.Time
			descripcion    string
			coinIn         int
			coinOut        int
			jackpot        int
			hpCancelCredit int
			gamePlayed     int
			gameWon        int
			denominacion   string
		)

		meters := make([]EGMMeter, 0)

		for cursor.Next() {
			err := cursor.Scan(
				&horario,
				&descripcion,
				&coinIn,
				&coinOut,
				&jackpot,
				&hpCancelCredit,
				&gamePlayed,
				&gameWon,
				&denominacion,
			)

			if err != nil {
				sqlLog.Errorf("ERROR: %s", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"DB": err.Error()})
				return
			}

			egmmeter := EGMMeter{
				TotalIn:  coinIn,
				TotalOut: coinOut,
				Wins:     gameWon,
				Bets:     gamePlayed,
				Jackpot:  jackpot,
				Handpay:  hpCancelCredit,
				Events:   descripcion,
				Date:     db.Time(horario.UnixMilli()),
			}
			meters = append(meters, egmmeter)
		}
		c.JSON(http.StatusOK, meters)
	})

	router.GET("/get_all_meters/:machineIdUnico", func(c *gin.Context) {
		ocDB := db.MysqlDb("TNZ_LEGACY")
		log := helpers.Log(c)
		IdUnico := c.Param("machineIdUnico")

		log.Debug("Machine ID >>>>>>>>>>>>>> ", IdUnico)

		type QueryParams struct {
			Time db.Time
		}

		var params QueryParams
		if err := c.ShouldBindQuery(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		myTime := params.Time.AsTime()

		cursor, closeCursor, queryErr, sqlLog := helpers.SqlQuery(log, ocDB, `SELECT x0.hora_registro, ed.descripcion, c.coinIn, c.coinOut, c.mdrop, c.jackpot, c.hpCancelCredit, c.gamePlayed, c.gameWon, c.doorOpen, c.powerReset, c.basePorcent, c.phyCoinIn, c.phyCoinOut, c.cancelCredits, c.billIn, c.machineProgresive, c.totalImpuesto, c.cantImpuesto, c.currentCredits, c.attendantProgresive, d.denominacion FROM (
			SELECT * FROM (
				SELECT e.id_evento, e.hora_registro, e.id_contador, e.clave
				 FROM eventos e
				 WHERE e.id_maquina = (select id_maquina from maquinas where idUnico=?)
				 AND e.hora_registro > subdate(?, interval 1 hour) 
				 ORDER BY e.hora_registro DESC
				 LIMIT 20
			) j0
			UNION
			SELECT * FROM (
				SELECT e.id_evento, e.hora_registro, e.id_contador, e.clave
				 FROM eventos e
				 WHERE e.id_maquina = (select id_maquina from maquinas where idUnico='22300307')
				 AND e.hora_registro > adddate(?, interval 10 minute) 
				 ORDER BY e.hora_registro ASC
				 LIMIT 10
			) j1
		) x0
		INNER JOIN contadores c ON x0.id_contador = c.id_contador
		INNER JOIN eventosdetalle ed ON x0.clave = ed.clave
		INNER JOIN denominacion d ON c.denomination = d.codigo
		ORDER BY x0.hora_registro DESC`, IdUnico, myTime, myTime)

		if queryErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Query": queryErr.Error()})
			return
		}

		defer closeCursor()

		var (
			horario             time.Time
			descripcion         string
			coinIn              int
			coinOut             int
			mdrop               int
			jackpot             int
			hpCancelCredit      int
			gamePlayed          int
			gameWon             int
			doorOpen            int
			powerReset          int
			basePorcent         string
			phyCoinIn           int
			phyCoinOut          int
			cancelCredits       int
			billIn              int
			machineProgresive   int
			totalImpuesto       int
			cantImpuesto        int
			currentCredits      int
			attendantProgresive int
			denominacion        string
		)

		meters := make([]EGMMeter, 0)

		for cursor.Next() {
			err := cursor.Scan(
				&horario,
				&descripcion,
				&coinIn,
				&coinOut,
				&mdrop,
				&jackpot,
				&hpCancelCredit,
				&gamePlayed,
				&gameWon,
				&doorOpen,
				&powerReset,
				&basePorcent,
				&phyCoinIn,
				&phyCoinOut,
				&cancelCredits,
				&billIn,
				&machineProgresive,
				&totalImpuesto,
				&cantImpuesto,
				&currentCredits,
				&attendantProgresive,
				&denominacion,
			)

			if err != nil {
				sqlLog.Errorf("ERROR: %s", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"DB": err.Error()})
				return
			}

			egmmeter := EGMMeter{
				TotalIn:             coinIn,
				TotalOut:            coinOut,
				Wins:                gameWon,
				Bets:                gamePlayed,
				Jackpot:             jackpot,
				Handpay:             hpCancelCredit,
				Date:                db.Time(horario.UnixMilli()),
				MachineDrop:         mdrop,
				DoorOpen:            doorOpen,
				PowerReset:          powerReset,
				BasePorcent:         basePorcent,
				PhysicalCoinIn:      phyCoinIn,
				PhysicalCoinOut:     phyCoinOut,
				CancelCredits:       cancelCredits,
				BillIn:              billIn,
				MachineProgresive:   machineProgresive,
				TotalTaxes:          totalImpuesto,
				TaxesQuantity:       cantImpuesto,
				CurrentCredits:      currentCredits,
				AttendantProgresive: attendantProgresive,
				Events:              descripcion,
			}

			meters = append(meters, egmmeter)
		}
		c.JSON(http.StatusOK, meters)
	})

	router.GET("/profit", func(c *gin.Context) {
		ocDB := db.MysqlDb("TNZ_MDQ")

		log := helpers.Log(c)

		type QueryParams struct {
			Time db.Time
		}

		var params QueryParams
		if err := c.ShouldBindQuery(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		time1 := params.Time.AsTime().Format("2006-01-02")

		log.Debug("TIME >>>>>>>>>>>> ", time1)

		//RECORDAR SACAR EL LIMIT A LA QUERY
		//EN TEST LA TABLA TECNOAZARMDQ ES SIN LA T MAYUSCULA

		cursor, closeCursor, queryErr, sqlLog := helpers.SqlQuery(log, ocDB, `SELECT m.idUnico, c.input_time, c.coinIn, c.coinOut, c.mdrop, c.jackpot, c.gamePlayed, c.gameWon, c.doorOpen, c.powerReset, c.currentCredits, c.cancelCredits, 
        c.billFive, c.billTen, c.billTwenty, c.billFifty, c.billOneHundred, c.billTwoHundred, c.billFiveHundred, c.billOneThousand,
	 (SELECT d.denominacion FROM denominacion d WHERE c.denomination = d.codigo) denomination, IF(f.id_maquina is null, "no_data", "")
		FROM maquinas m
		LEFT JOIN facturacion f ON m.id_maquina = f.id_maquina
		LEFT JOIN contadores c ON (f.id_contador_inicial = c.id_contador OR f.id_contador_final = c.id_contador)
		WHERE
		(m.fechaBaja IS NULL OR m.fechaBaja = '0000-00-00' OR m.fechaBaja > ?) 
		AND f.fecha = ?
		AND c.input_time BETWEEN CONCAT(?, ' 07:00:00') AND CONCAT(DATE_ADD(?, INTERVAL 1 DAY), ' 07:00:00') 
	UNION
	SELECT m.idUnico, c.input_time, c.coinIn, c.coinOut, c.mdrop, c.jackpot, c.gamePlayed, c.gameWon, c.doorOpen, c.powerReset, c.currentCredits, 
    c.cancelCredits, c.billFive, c.billTen, c.billTwenty, c.billFifty, c.billOneHundred, c.billTwoHundred, c.billFiveHundred, c.billOneThousand,
		0 denomination, IFNULL(IF ( c.linea = 1 AND c.tipo is null,	"rollover",	IF ( c.tipo = "CR", "clear_full_ram", IF( c.tipo = "CP", "clear_partial_ram",""))),"")
		FROM maquinas m
		JOIN anomalias c ON m.id_maquina = c.id_maquina
		WHERE
		(m.fechaBaja IS NULL OR m.fechaBaja = '0000-00-00' OR m.fechaBaja > ?) 
		AND c.input_time BETWEEN CONCAT(?, ' 07:00:00') AND CONCAT(DATE_ADD(?, INTERVAL 1 DAY), ' 07:00:00') 
		AND (c.linea = 1 or c.linea = 2)
		ORDER BY 1, 2
		`, time1, time1, time1, time1, time1, time1, time1)

		if queryErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Query": queryErr.Error()})
			return
		}

		defer closeCursor()

		fullReport := make([]Report, 0)

		var report *Report
		var machineID int

		for cursor.Next() {

			var (
				idUnico      int
				input_time   time.Time
				actualMeters Meters
				denomination float32
				anomaly      string
			)

			err := cursor.Scan(
				&idUnico,
				&input_time,
				&actualMeters.TotalIn,
				&actualMeters.TotalOut,
				&actualMeters.MachineDrop,
				&actualMeters.Jackpot,
				&actualMeters.Games,
				&actualMeters.GameWon,
				&actualMeters.DoorOpen,
				&actualMeters.PowerReset,
				&actualMeters.CurrentCredits,
				&actualMeters.CancelCredits,
				&actualMeters.BillFive,
				&actualMeters.BillTen,
				&actualMeters.BillTwenty,
				&actualMeters.BillFifty,
				&actualMeters.BillOneHundred,
				&actualMeters.BillTwoHundred,
				&actualMeters.BillFiveHundred,
				&actualMeters.BillOneThousand,
				&denomination,
				&anomaly,
			)

			if err != nil {
				sqlLog.Errorf("ERROR: %s", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"DB": err.Error()})
				return
			}

			if idUnico != machineID {
				machineID = idUnico

				if report != nil {
					report.Periods[len(report.Periods)-1].CalcDelta()
					fullReport = append(fullReport, *report)
				}
				report = &Report{
					MachineID:         idUnico,
					DenominationValue: denomination,
					Periods:           make([]Period, 1),
				}
				report.Periods[0].MetersFrom.At = db.Time(input_time.UnixMilli())
				report.Periods[0].MetersFrom.Meters = actualMeters.Clone()

			} else {
				idx := len(report.Periods) - 1
				if report.Periods[idx].MetersTo.At != 0 {
					period := Period{
						MetersFrom: MetersWithTime{
							At:     report.Periods[idx].MetersTo.At,
							Meters: report.Periods[idx].MetersTo.Clone(),
						},
					}
					report.Periods[idx].CalcDelta()
					report.Periods = append(report.Periods, period)
					idx = idx + 1
				}

				if report.Periods[idx].MetersFrom.Meters.CheckIfRollover(actualMeters) {

					report.Periods[idx].MetersTo.At = report.Periods[idx].MetersFrom.At + 1
					report.Periods[idx].MetersTo.Meters.CalcRollover(report.Periods[idx].MetersFrom.Meters, actualMeters, 100000000)

					period := Period{
						MetersTo: MetersWithTime{
							At:     db.Time(input_time.UnixMilli()),
							Meters: actualMeters,
						},
					}
					period.Type = PeriodTypeRollover

					period.MetersFrom.At = report.Periods[idx].MetersFrom.At + 2
					period.MetersFrom.Meters.CalcRollover(report.Periods[idx].MetersFrom.Meters, actualMeters, 0)

					report.Periods[idx].CalcDelta()
					report.Periods = append(report.Periods, period)
				} else {
					report.Periods[idx].MetersTo.At = db.Time(input_time.UnixMilli())
					report.Periods[idx].MetersTo.Meters = actualMeters.Clone()
				}
			}
		}

		if report == nil {
			res := make([]Report, 0)
			c.JSON(http.StatusAccepted, res)

		}

		report.Periods[len(report.Periods)-1].CalcDelta()
		fullReport = append(fullReport, *report)

		c.JSON(http.StatusOK, fullReport)
	})
	//REPORTE PARA SABER AQUELLAS MAQUINAS QUE DEBERIAN HABER REPORTADO EL IMPUESTO AFT Y NO LO HICIERON
	router.GET("/reportingMachines", func(c *gin.Context) {
		ocDB := db.MysqlDb("TNZ_LEGACY")
		log := helpers.Log(c)

		cursor, closeCursor, queryErr, sqlLog := helpers.SqlQuery(log, ocDB, `SELECT m.IdUnico, count(*) qty, SUM(i.delta_billIn) AS total_delta_billIn, min(i.input_time), max(i.input_time)
		FROM tecnoazarMdq.impuesto AS i
		LEFT JOIN tecnoazarMdq.maquinas AS m ON m.id_maquina = i.id_maquina
		WHERE i.input_time >= DATE_SUB(NOW(), INTERVAL 24 HOUR) AND i.descuento = 0 AND i.delta_billIn > 0
		GROUP BY m.IdUnico
		ORDER BY i.input_time DESC;`)

		if queryErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Query": queryErr.Error()})
			return
		}

		defer closeCursor()

		reportingMachines := make([]ReportingMachine, 0)

		for cursor.Next() {
			var (
				id_maquina   int
				delta_billIn int
				descuento    float32
				input_time   time.Time
			)
			err := cursor.Scan(
				&id_maquina,
				&delta_billIn,
				&descuento,
				&input_time,
			)

			if err != nil {
				sqlLog.Errorf("ERROR: %s", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"DB": err.Error()})
				return
			}

			reportingMachine := ReportingMachine{
				MachineID:   id_maquina,
				DeltaBillIn: delta_billIn,
				Discount:    descuento,
				Date:        db.Time(input_time.UnixMilli()),
			}

			reportingMachines = append(reportingMachines, reportingMachine)
		}
		c.JSON(http.StatusOK, reportingMachines)
	})
}
