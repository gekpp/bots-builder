package tests

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func ConnectTestContainers(bindMounts map[string]string) (*sql.DB, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15.1",
		ExposedPorts: []string{"5432/tcp"},
		Mounts:       makeMounts(bindMounts),
		Env: map[string]string{
			"POSTGRES_PASSWORD": "admin",
		},
		// BindMounts:   bindMounts,
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}

	time.Sleep(time.Second)

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	// Get the port mapped to 5432 and set as ENV
	p, err := postgresC.MappedPort(ctx, "5432")
	if err != nil {
		panic(err)
	}
	return connect("localhost", p.Port(), "postgres", "postgres", "admin"), func() {
		_ = postgresC.Terminate(ctx)
	}
}

func makeMounts(src2dst map[string]string) testcontainers.ContainerMounts {
	res := testcontainers.ContainerMounts{}
	for src, dst := range src2dst {
		res = append(res, testcontainers.ContainerMount{
			Source: testcontainers.GenericBindMountSource{HostPath: src},
			Target: testcontainers.ContainerMountTarget(dst),
		})
	}
	return res
}

func connect(host, port, dbName, user, password string) *sql.DB {

	conn, err := pq.NewConnector(fmt.Sprintf("user='%s' password='%s' host=%v port=%v dbname=%s sslmode=disable",
		user, password, host, port, dbName))
	if err != nil {
		panic(fmt.Sprintf("can't connect to db: %v", err))
	}
	db := sql.OpenDB(conn)
	time.Sleep(2 * time.Second)

	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("can't ping db: %v", err))
	}
	return db
}
