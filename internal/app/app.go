package app

import (
	"os"
	"walletapp/config"
)

// App represents whole application at its highest level (( of abstraction )).
type App struct {
}

// Run is a shortcut for app.New(cfg).Run(),
// combined with non-zero-code exit if error happens.
func Run(cfg config.Config) {

	app := New(cfg)

	err := app.Run()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}

// # App constructor.
// Almost all "initialisation" and "execution" happens separatly via .Run() method.
//
// Application itself do NOT validate configuration, it assumes that config is valid,
// if something breaks because of unsuitable config - its a .Run() caller responsibility.
//
// Application do not stores global state, so from the application perspective code like:
//
//	app := New(cfg)
//	err := app.Run()
//	if err != nil {/* ... error handle ...  */}
//	app = New(cfg)
//	err = app.Run()
//
// Considered as Safe.
func New(cfg config.Config) *App {
	ap := &App{}

	return ap
}

// Run the app.
// Stages :
//   - a,
//   - a.
//
// Application itself do NOT handling error that can happen in its subsystems.
//
// Application itself do NOT handling panic of its subsystem.
//
// If error or panic happened and subsystem cannot handle it - Run() will perform
// emergency shutdown and return with non-nil error.
//
// Nil error guarantees that all subsystems shutted down successfully.
//
// Nil error do NOT guarantees that external clients that application subsystems are
// interacted with are ready to re-use,
// or have the same "state" as it was on previous Run() call.
func (ap *App) Run() error {
	var err error

	if err != nil {

		return err
	}

	return nil
}
