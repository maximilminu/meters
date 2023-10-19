package main

import (
	"os"

	"gitlab.sistematis.com.ar/OC/be/api/egm-meter/src"
	"gitlab.sistematis.com.ar/OC/be/common/audit"
	"gitlab.sistematis.com.ar/OC/be/common/auth"
	"gitlab.sistematis.com.ar/OC/be/common/db"
	"gitlab.sistematis.com.ar/OC/be/common/microservice"
	"gitlab.sistematis.com.ar/OC/be/common/session"
	helpers "gitlab.sistematis.com.ar/OPEN-SOURCE/GO/api-helpers"
)

func main() {

	router := microservice.Init("", true)
	audit.Init("egmmeter")

	db.ArangoDbInit("OC")
	microservice.AddReadinessCheck("arangodb_OC", db.ArangoDbGetReadinessCheck("OC"))
	session.Init(router)
	microservice.AddReadinessCheck("redis_session", session.ReadinessCheck)
	auth.Init(router)
	microservice.AddReadinessCheck("redis_auth", auth.ReadinessCheck)

	if helpers.GetEnv("TNZ_TITO_MYSQL_USER", "") != "" {
		db.MysqlDbInit("TNZ_TITO")
		// db.MysqlDbInit("TNZ_TITO")
		microservice.AddReadinessCheck("TNZ_TITO", db.MysqlDbGetReadinessCheck("TNZ_TITO"))
	}

	if helpers.GetEnv("TNZ_LEGACY_MYSQL_USER", "") != "" {
		db.MysqlDbInit("TNZ_LEGACY")
		// db.MysqlDbInit("TNZ_LEGACY")
		microservice.AddReadinessCheck("TNZ_LEGACY", db.MysqlDbGetReadinessCheck("TNZ_LEGACY"))
	}

	if helpers.GetEnv("TNZ_MDQ_MYSQL_USER", "") != "" {
		db.MysqlDbInit("TNZ_MDQ")
		// db.MysqlDbInit("TNZ_TITO")
		// microservice.AddReadinessCheck("TNZ_TITO", db.MysqlDbGetReadinessCheck("TNZ_TITO"))
	}

	ret := microservice.ReadinessCheck()
	if len(ret) > 0 {
		for name, err := range ret {
			helpers.Logger().Errorf("Problems with Readiness test %s: %s", name, err)
		}
		os.Exit(1)
		return
	}

	// router.Use(auth.AuthOrFailWithAccessTo)

	src.GET(router)

	microservice.Bind()

	helpers.Logger().Info("Server exiting")
	if helpers.GetEnv("TNZ_TITO_MYSQL_USER", "") != "" {
		db.MysqlDbShutdown("TNZ_TITO")
	}
}
