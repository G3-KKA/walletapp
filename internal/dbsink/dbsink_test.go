package dbsink

import (
	"context"
	"encoding/json"
	"testing"
	"time"
	"walletapp/config"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (

	// In-Code skip hook, should be replaced with proper IE test.
	_IgnoreDBIntegrationTest = false
)

type TS struct {
	TMM time.Time `json:"TMM"`
	STR string    `json:"STR"`
}

// SelfUnmarshal should be used when you have to compare
// basicaly two exactly the same objects, but one got processed by
// Marshal=>Unmarshal, and its [time.Time] field(s) got corrupted,
// resulting in False Negative comparation via reflect.DeepEqual or assert.Equal.
//
//nolint:all // Insane function by itself.
func SelfUnmarshal[V any](t *testing.T, v *V) {
	t.Helper()
	bts, err := json.Marshal(v)
	require.NoError(t, err)
	err = json.Unmarshal(bts, v)
	require.NoError(t, err)
}

func TestNewConnect(t *testing.T) {
	if _IgnoreDBIntegrationTest {
		t.SkipNow()
	}

	// Linter is angry, and he is right on point, test invalid.
	cfg := config.Database{
		DBUser: config.DBUser{
			DBRole:     "postgres", // Default.
			DBPassword: "postgres", // Default.
		},
		DBAddress: "В ДУШЕ НЕ ЧАЮ", // TODO: Нужно резолвить имя контейнера.
		DBPort:    "5432",
		DBName:    "postgres", // Default.
	}

	// TODO: Основная задача - параллельное интеграционное тестирование
	/*
		step1: Заставить местный докер поднять анонимную базу данных
		UPDATE: возможно testcontainers всё таки стартует всё сам,
		это было бы очень очень хорошо.


		stepN-1: в документации к редису https://golang.testcontainers.org/quickstart/
		написано много непонятного про рандомные порты и t.Parallel, ещё раз перечитай
		скорее всего случайный порт маппится сам и мне не нужно выдумать колесо.


		stepN: testcontainer

	*/

	// TODO: подумать как лучше.
	/*
	 	testcontainers.WithLogger(zerolog.Ctx(context.TODO()))
	   	testcontainers.WithLogger(testcontainers.TestLogger(t))
	*/
	// TODO: по хорошему для параллельного тестирования
	// необходимо вмапить 5432 на контейнере на свободный 25256+ порт
	// функция для этого у меня уже есть.

	// BTW очень инетерсная штука, API docker'a это офигенно.

	//c, err := postgres.Run()

	// TODO: >>>>>> ОБЯЗАТЕЛЬНО ПРОАНАЛИЗИРОВАТЬ ВНУТРЯНКУ ЭТОЙ ФУНКЦИИ.
	//	postgres.BasicWaitStrategies()
	// <<<<<<
	req := testcontainers.ContainerRequest{
		FromDockerfile:    testcontainers.FromDockerfile{},
		HostAccessPorts:   []int{},
		Image:             "postgres:17.0",
		ImageSubstitutors: []testcontainers.ImageSubstitutor{},
		Entrypoint:        []string{},
		Env: map[string]string{
			"POSTGRES_USER":     cfg.DBRole,
			"POSTGRES_PASSWORD": cfg.DBPassword,
			"POSTGRES_DB":       cfg.DBName, // defaults to the user name
		},
		ExposedPorts: []string{
			"5432/tcp",
		},
		// Cmd:             []string{},
		// Без этого разве не будет дефолтная cmd самого контейнера?
		Cmd:             []string{"postgres", "-c", "fsync=off"},
		Labels:          map[string]string{},
		Mounts:          []testcontainers.ContainerMount{},
		Tmpfs:           map[string]string{},
		RegistryCred:    "",
		WaitingFor:      nil,
		Name:            "",
		Hostname:        "",
		WorkingDir:      "",
		ExtraHosts:      []string{},
		Privileged:      false,
		Networks:        []string{},
		NetworkAliases:  map[string][]string{},
		NetworkMode:     "",
		Resources:       container.Resources{},
		Files:           []testcontainers.ContainerFile{},
		User:            "",
		SkipReaper:      false,
		ReaperImage:     "",
		ReaperOptions:   []testcontainers.ContainerOption{},
		AutoRemove:      false,
		AlwaysPullImage: false,
		ImagePlatform:   "",
		Binds:           []string{},
		ShmSize:         0,
		CapAdd:          []string{},
		CapDrop:         []string{},
		ConfigModifier: func(*container.Config) {
		},
		HostConfigModifier: func(*container.HostConfig) {
		},
		EnpointSettingsModifier: func(m map[string]*network.EndpointSettings) {},
		LifecycleHooks:          []testcontainers.ContainerLifecycleHooks{},
		LogConsumerCfg:          &testcontainers.LogConsumerConfig{},
	}

	// TODO:	testcontainers.CleanupContainer(TESTING_T, CONTAINER)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	genericR := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		ProviderType:     0,
		Logger:           testcontainers.TestLogger(t),
		Reuse:            false,
	}

	genericC, err := testcontainers.GenericContainer(ctx, genericR)

	testcontainers.CleanupContainer(t, genericC)

	assert.NoError(t, err)

	err = genericC.Start(ctx)

	assert.NoError(t, err)

	cfg.DBAddress, err = genericC.ContainerIP(ctx)

	assert.NoError(t, err)

	port, err := genericC.MappedPort(ctx, nat.Port("5432"))

	t.Log(port)

	assert.NoError(t, err)

	cfg.DBPort = port.Port()

	l := zerolog.New(zerolog.NewTestWriter(t))

	db := New(&l, cfg)

	err = db.Ping(context.TODO())

	require.NoError(t, err)

	err = db.Shutdown(context.TODO())

	// Useless, no error will be returned anyway.
	assert.NoError(t, err)
}

func TestThroughWrapper(t *testing.T) {
	if _IgnoreDBIntegrationTest {
		t.SkipNow()
	}

	// Linter is angry, and he is right on point, test invalid.
	cfg := config.Database{
		DBUser: config.DBUser{
			DBRole:     "postgres", // Default.
			DBPassword: "postgres", // Default.
		},
		DBAddress: "В ДУШЕ НЕ ЧАЮ", // TODO: Нужно резолвить имя контейнера.
		DBPort:    "5432",
		DBName:    "testdb", // Default.
	}

	ctx := context.Background()

	// 1. Start the postgres ctr and run any migrations on it
	ctr, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(cfg.DBName),
		postgres.WithUsername(cfg.DBRole),
		postgres.WithPassword(cfg.DBPassword),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	testcontainers.CleanupContainer(t, ctr)
	require.NoError(t, err)

	// Run any migrations on the database
	_, _, err = ctr.Exec(ctx, []string{"psql", "-U", cfg.DBRole, "-d", cfg.DBName, "-c", "CREATE TABLE users (id SERIAL, name TEXT NOT NULL, age INT NOT NULL)"})
	require.NoError(t, err)

	// 2. Create a snapshot of the database to restore later
	// tt.options comes the test case, it can be specified as e.g. `postgres.WithSnapshotName("custom-snapshot")` or omitted, to use default name
	//err = ctr.Snapshot(ctx, tt.options...)
	require.NoError(t, err)

	dbURL, err := ctr.ConnectionString(ctx)
	require.NoError(t, err)

	t.Run("Test inserting a user", func(t *testing.T) {
		t.Cleanup(func() {
			// 3. In each test, reset the DB to its snapshot state.
			err = ctr.Restore(ctx)
			require.NoError(t, err)
		})

		conn, err := pgx.Connect(context.Background(), dbURL)
		require.NoError(t, err)
		defer conn.Close(context.Background())

		_, err = conn.Exec(ctx, "INSERT INTO users(name, age) VALUES ($1, $2)", "test", 42)
		require.NoError(t, err)

		var name string
		var age int64
		err = conn.QueryRow(context.Background(), "SELECT name, age FROM users LIMIT 1").Scan(&name, &age)
		require.NoError(t, err)

		require.Equal(t, "test", name)
		require.EqualValues(t, 42, age)
	})

	// 4. Run as many tests as you need, they will each get a clean database
	t.Run("Test querying empty DB", func(t *testing.T) {
		t.Cleanup(func() {
			err = ctr.Restore(ctx)
			require.NoError(t, err)
		})

		conn, err := pgx.Connect(context.Background(), dbURL)
		require.NoError(t, err)
		defer conn.Close(context.Background())

		var name string
		var age int64
		err = conn.QueryRow(context.Background(), "SELECT name, age FROM users LIMIT 1").Scan(&name, &age)

		// Errors is NIL cuz in cleanup we are trying to drop 'postgres' database, which is impossible
		// even if we are postgres (superuser-owner) user himself.
		require.ErrorIs(t, err, pgx.ErrNoRows)
	})

	l := zerolog.New(zerolog.NewTestWriter(t))

	db := New(&l, cfg)

	err = db.Ping(context.TODO())

	require.NoError(t, err)

	err = db.Shutdown(context.TODO())

	// Useless, no error will be returned anyway.
	assert.NoError(t, err)

}
